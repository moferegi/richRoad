package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	requestModel "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type englishWordAPIResp struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

type englishWordErrorListItem struct {
	WordID           uint   `json:"wordId"`
	WrongCount       int    `json:"wrongCount"`
	LastExpectedChar string `json:"lastExpectedChar"`
	LastInputChar    string `json:"lastInputChar"`
}

type englishWordErrorListData struct {
	List     []englishWordErrorListItem `json:"list"`
	Total    int64                      `json:"total"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"pageSize"`
}

type englishWordListData struct {
	List     []model.EnglishWord `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}

type englishWordPreflightData struct {
	Word       string `json:"word"`
	VoiceUS    string `json:"voiceUs"`
	VoiceUK    string `json:"voiceUk"`
	AudioUS    string `json:"audioUs"`
	AudioUK    string `json:"audioUk"`
	DurationMS int64  `json:"durationMs"`
}

func setupEnglishWordAPITest(t *testing.T) (wordID uint, token string) {
	t.Helper()

	originalDB := global.GVA_DB
	originalLog := global.GVA_LOG
	originalJWT := global.GVA_CONFIG.JWT

	t.Cleanup(func() {
		global.GVA_DB = originalDB
		global.GVA_LOG = originalLog
		global.GVA_CONFIG.JWT = originalJWT
	})

	gin.SetMode(gin.TestMode)
	global.GVA_LOG = zap.NewNop()
	global.GVA_CONFIG.JWT.SigningKey = "english-word-api-test-sign-key"
	global.GVA_CONFIG.JWT.ExpiresTime = "2h"
	global.GVA_CONFIG.JWT.BufferTime = "10m"
	global.GVA_CONFIG.JWT.Issuer = "english-word-api-test"

	db, err := gorm.Open(sqlite.Open("file:english_word_api_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err = db.AutoMigrate(
		&model.EnglishWord{},
		&model.EnglishWordErrorLog{},
		&model.EnglishCategoryWord{},
		&model.EnglishWordSentence{},
		&model.UserCollection{},
		&clientModel.SysConfig{},
	); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	word := model.EnglishWord{Word: "api_report_error_word_" + uuid.NewString(), Explanation: "{}"}
	if err = db.Create(&word).Error; err != nil {
		t.Fatalf("seed word failed: %v", err)
	}

	global.GVA_DB = db

	j := utils.NewJWT()
	claims := j.CreateClaims(requestModel.BaseClaims{
		UUID:        uuid.New(),
		ID:          7001,
		Username:    "api_test_user",
		NickName:    "api_test_user",
		AuthorityId: 8080,
	})
	token, err = j.CreateToken(claims)
	if err != nil {
		t.Fatalf("create token failed: %v", err)
	}

	return word.ID, token
}

func performEnglishWordAPIRequest(t *testing.T, method, target, token string, body []byte, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("x-token", token)
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req

	handler(c)
	return recorder
}

func decodeEnglishWordAPIResp(t *testing.T, recorder *httptest.ResponseRecorder) englishWordAPIResp {
	t.Helper()

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected http 200, got %d, body=%s", recorder.Code, recorder.Body.String())
	}

	var resp englishWordAPIResp
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response failed: %v, body=%s", err, recorder.Body.String())
	}
	return resp
}

func findWordErrorItemByWordID(items []englishWordErrorListItem, wordID uint) (englishWordErrorListItem, bool) {
	for _, item := range items {
		if item.WordID == wordID {
			return item, true
		}
	}
	return englishWordErrorListItem{}, false
}

func upsertEnglishWordSysConfig(t *testing.T, key, value string) {
	t.Helper()

	var cfg clientModel.SysConfig
	err := global.GVA_DB.Where("config_key = ?", key).First(&cfg).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("query sys config failed: %v", err)
		}

		cfg = clientModel.SysConfig{
			ConfigKey:   key,
			ConfigValue: value,
			ConfigName:  key,
			ConfigGroup: "english_learning",
		}
		if createErr := global.GVA_DB.Create(&cfg).Error; createErr != nil {
			t.Fatalf("create sys config failed: %v", createErr)
		}
		return
	}

	if updateErr := global.GVA_DB.Model(&clientModel.SysConfig{}).Where("id = ?", cfg.ID).Update("config_value", value).Error; updateErr != nil {
		t.Fatalf("update sys config failed: %v", updateErr)
	}
}

func TestEnglishWordAPI_ReportThenGetErrorLogList(t *testing.T) {
	wordID, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	postBody1 := []byte(fmt.Sprintf(`{"wordId":%d,"wrongIndex":1,"expectedChar":"ab","inputChar":"xy"}`, wordID))
	postResp1 := performEnglishWordAPIRequest(t, http.MethodPost, "/englishLearning/word/reportError", token, postBody1, api.ReportWordError)
	decodedPost1 := decodeEnglishWordAPIResp(t, postResp1)
	if decodedPost1.Code != 0 {
		t.Fatalf("first report should succeed, code=%d msg=%s", decodedPost1.Code, decodedPost1.Msg)
	}

	getResp1 := performEnglishWordAPIRequest(t, http.MethodGet, "/englishLearning/word/getErrorLogList?page=1&pageSize=10", token, nil, api.GetWordErrorLogList)
	decodedGet1 := decodeEnglishWordAPIResp(t, getResp1)
	if decodedGet1.Code != 0 {
		t.Fatalf("first get list should succeed, code=%d msg=%s", decodedGet1.Code, decodedGet1.Msg)
	}

	var listData1 englishWordErrorListData
	if err := json.Unmarshal(decodedGet1.Data, &listData1); err != nil {
		t.Fatalf("decode list data failed: %v, raw=%s", err, string(decodedGet1.Data))
	}
	item1, ok := findWordErrorItemByWordID(listData1.List, wordID)
	if !ok {
		t.Fatalf("expected wordId=%d in list, list=%+v", wordID, listData1.List)
	}
	if item1.WrongCount != 1 {
		t.Fatalf("expected wrongCount=1 after first report, got=%d", item1.WrongCount)
	}
	if item1.LastExpectedChar != "a" || item1.LastInputChar != "x" {
		t.Fatalf("expected sanitized chars a/x, got %s/%s", item1.LastExpectedChar, item1.LastInputChar)
	}

	postBody2 := []byte(fmt.Sprintf(`{"wordId":%d,"wrongIndex":3,"expectedChar":"c","inputChar":"z"}`, wordID))
	postResp2 := performEnglishWordAPIRequest(t, http.MethodPost, "/englishLearning/word/reportError", token, postBody2, api.ReportWordError)
	decodedPost2 := decodeEnglishWordAPIResp(t, postResp2)
	if decodedPost2.Code != 0 {
		t.Fatalf("second report should succeed, code=%d msg=%s", decodedPost2.Code, decodedPost2.Msg)
	}

	getResp2 := performEnglishWordAPIRequest(t, http.MethodGet, "/englishLearning/word/getErrorLogList?page=1&pageSize=10", token, nil, api.GetWordErrorLogList)
	decodedGet2 := decodeEnglishWordAPIResp(t, getResp2)
	if decodedGet2.Code != 0 {
		t.Fatalf("second get list should succeed, code=%d msg=%s", decodedGet2.Code, decodedGet2.Msg)
	}

	var listData2 englishWordErrorListData
	if err := json.Unmarshal(decodedGet2.Data, &listData2); err != nil {
		t.Fatalf("decode second list data failed: %v, raw=%s", err, string(decodedGet2.Data))
	}
	item2, ok := findWordErrorItemByWordID(listData2.List, wordID)
	if !ok {
		t.Fatalf("expected wordId=%d in second list, list=%+v", wordID, listData2.List)
	}
	if item2.WrongCount != 2 {
		t.Fatalf("expected wrongCount=2 after second report, got=%d", item2.WrongCount)
	}
	if item2.LastExpectedChar != "c" || item2.LastInputChar != "z" {
		t.Fatalf("expected updated chars c/z, got %s/%s", item2.LastExpectedChar, item2.LastInputChar)
	}
}

func TestEnglishWordAPI_ReportWordError_AuthFailure(t *testing.T) {
	wordID, validToken := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}
	body := []byte(fmt.Sprintf(`{"wordId":%d,"wrongIndex":1,"expectedChar":"a","inputChar":"b"}`, wordID))

	testCases := []struct {
		name  string
		token string
	}{
		{name: "missing token", token: ""},
		{name: "invalid token", token: validToken + ".broken"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := performEnglishWordAPIRequest(t, http.MethodPost, "/englishLearning/word/reportError", tc.token, body, api.ReportWordError)
			decoded := decodeEnglishWordAPIResp(t, resp)
			if decoded.Code == 0 {
				t.Fatalf("expected auth failure, got success msg=%s", decoded.Msg)
			}
			if !strings.Contains(decoded.Msg, "获取用户信息失败") {
				t.Fatalf("expected auth failure message, got=%s", decoded.Msg)
			}
		})
	}
}

func TestEnglishWordAPI_GetWordErrorLogList_AuthFailure(t *testing.T) {
	_, validToken := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	testCases := []struct {
		name  string
		token string
	}{
		{name: "missing token", token: ""},
		{name: "invalid token", token: validToken + ".broken"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := performEnglishWordAPIRequest(t, http.MethodGet, "/englishLearning/word/getErrorLogList?page=1&pageSize=10", tc.token, nil, api.GetWordErrorLogList)
			decoded := decodeEnglishWordAPIResp(t, resp)
			if decoded.Code == 0 {
				t.Fatalf("expected auth failure, got success msg=%s", decoded.Msg)
			}
			if !strings.Contains(decoded.Msg, "获取用户信息失败") {
				t.Fatalf("expected auth failure message, got=%s", decoded.Msg)
			}
		})
	}
}

func TestEnglishWordAPI_ReportWordError_ParamError(t *testing.T) {
	wordID, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	testCases := []struct {
		name string
		body []byte
	}{
		{name: "missing required wordId", body: []byte(`{"wrongIndex":1,"expectedChar":"a","inputChar":"b"}`)},
		{name: "invalid json", body: []byte(`{"wordId":`)},
		{name: "wordId as wrong type", body: []byte(fmt.Sprintf(`{"wordId":"%d","wrongIndex":1}`, wordID))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := performEnglishWordAPIRequest(t, http.MethodPost, "/englishLearning/word/reportError", token, tc.body, api.ReportWordError)
			decoded := decodeEnglishWordAPIResp(t, resp)
			if decoded.Code == 0 {
				t.Fatalf("expected parameter failure, got success msg=%s", decoded.Msg)
			}
			if !strings.Contains(decoded.Msg, "参数错误") {
				t.Fatalf("expected parameter error message, got=%s", decoded.Msg)
			}
		})
	}
}

func TestEnglishWordAPI_ReportWordError_WordNotFound(t *testing.T) {
	_, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}
	body := []byte(`{"wordId":999999,"wrongIndex":2,"expectedChar":"d","inputChar":"x"}`)

	resp := performEnglishWordAPIRequest(t, http.MethodPost, "/englishLearning/word/reportError", token, body, api.ReportWordError)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected report failure for missing word, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "上报失败") {
		t.Fatalf("expected report failure message, got=%s", decoded.Msg)
	}
}

func TestEnglishWordAPI_GetWordErrorLogList_ParamError(t *testing.T) {
	_, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	resp := performEnglishWordAPIRequest(t, http.MethodGet, "/englishLearning/word/getErrorLogList?page=not-a-number&pageSize=10", token, nil, api.GetWordErrorLogList)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected query parameter failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "参数错误") {
		t.Fatalf("expected parameter error message, got=%s", decoded.Msg)
	}
}

func TestEnglishWordAPI_FindEnglishWord_SupportsLowercaseID(t *testing.T) {
	wordID, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	resp := performEnglishWordAPIRequest(t, http.MethodGet, fmt.Sprintf("/englishLearning/word/findWord?id=%d", wordID), token, nil, api.FindEnglishWord)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected findWord success with lowercase id, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var word model.EnglishWord
	if err := json.Unmarshal(decoded.Data, &word); err != nil {
		t.Fatalf("decode word data failed: %v raw=%s", err, string(decoded.Data))
	}
	if word.ID != wordID {
		t.Fatalf("expected word id=%d, got=%d", wordID, word.ID)
	}
}

func TestEnglishWordAPI_FindEnglishWord_SupportsUppercaseID(t *testing.T) {
	wordID, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	resp := performEnglishWordAPIRequest(t, http.MethodGet, fmt.Sprintf("/englishLearning/word/findWord?ID=%d", wordID), token, nil, api.FindEnglishWord)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected findWord success with uppercase ID, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var word model.EnglishWord
	if err := json.Unmarshal(decoded.Data, &word); err != nil {
		t.Fatalf("decode word data failed: %v raw=%s", err, string(decoded.Data))
	}
	if word.ID != wordID {
		t.Fatalf("expected word id=%d, got=%d", wordID, word.ID)
	}
}

func TestEnglishWordAPI_GetWordList_NormalizedPage(t *testing.T) {
	_, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	resp := performEnglishWordAPIRequest(t, http.MethodGet, "/englishLearning/word/getWordList", token, nil, api.GetWordList)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected getWordList success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var listData englishWordListData
	if err := json.Unmarshal(decoded.Data, &listData); err != nil {
		t.Fatalf("decode list data failed: %v raw=%s", err, string(decoded.Data))
	}
	if listData.Page != 1 || listData.PageSize != 10 {
		t.Fatalf("expected normalized page/pageSize=1/10, got=%d/%d", listData.Page, listData.PageSize)
	}
}

func TestEnglishWordAPI_GetWordList_ChapterFilterAndDistinct(t *testing.T) {
	seedWordID, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	wordInChapter := model.EnglishWord{Word: "api_word_in_chapter_" + uuid.NewString(), Explanation: "{}"}
	wordOtherChapter := model.EnglishWord{Word: "api_word_other_chapter_" + uuid.NewString(), Explanation: "{}"}
	if err := global.GVA_DB.Create(&wordInChapter).Error; err != nil {
		t.Fatalf("seed wordInChapter failed: %v", err)
	}
	if err := global.GVA_DB.Create(&wordOtherChapter).Error; err != nil {
		t.Fatalf("seed wordOtherChapter failed: %v", err)
	}

	chapterA := uint(1001)
	chapterB := uint(1002)
	assocs := []model.EnglishCategoryWord{
		{ChapterID: chapterA, WordID: seedWordID, Sort: 3},
		{ChapterID: chapterA, WordID: seedWordID, Sort: 1}, // duplicate association for same word
		{ChapterID: chapterA, WordID: wordInChapter.ID, Sort: 2},
		{ChapterID: chapterB, WordID: wordOtherChapter.ID, Sort: 1},
	}
	if err := global.GVA_DB.Create(&assocs).Error; err != nil {
		t.Fatalf("seed chapter associations failed: %v", err)
	}

	resp := performEnglishWordAPIRequest(
		t,
		http.MethodGet,
		fmt.Sprintf("/englishLearning/word/getWordList?chapterId=%d&page=1&pageSize=1", chapterA),
		token,
		nil,
		api.GetWordList,
	)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected getWordList success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var listData englishWordListData
	if err := json.Unmarshal(decoded.Data, &listData); err != nil {
		t.Fatalf("decode list data failed: %v raw=%s", err, string(decoded.Data))
	}

	if listData.Total != 2 {
		t.Fatalf("expected total=2 distinct words in chapterA, got=%d", listData.Total)
	}
	if len(listData.List) != 1 {
		t.Fatalf("expected pageSize=1 result length=1, got=%d", len(listData.List))
	}
	if listData.List[0].ID != seedWordID {
		t.Fatalf("expected first row to be seedWord by MIN(sort) ordering, got id=%d", listData.List[0].ID)
	}
}

func TestEnglishWordAPI_GetWordErrorLogList_NormalizedPage(t *testing.T) {
	_, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	resp := performEnglishWordAPIRequest(t, http.MethodGet, "/englishLearning/word/getErrorLogList", token, nil, api.GetWordErrorLogList)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected getErrorLogList success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var listData englishWordErrorListData
	if err := json.Unmarshal(decoded.Data, &listData); err != nil {
		t.Fatalf("decode error list data failed: %v raw=%s", err, string(decoded.Data))
	}
	if listData.Page != 1 || listData.PageSize != 10 {
		t.Fatalf("expected normalized page/pageSize=1/10, got=%d/%d", listData.Page, listData.PageSize)
	}
}

func TestEnglishWordAPI_FindEnglishWord_ParamError(t *testing.T) {
	_, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	testCases := []struct {
		name   string
		target string
	}{
		{name: "missing id", target: "/englishLearning/word/findWord"},
		{name: "invalid id value", target: "/englishLearning/word/findWord?ID=abc"},
		{name: "zero id", target: "/englishLearning/word/findWord?ID=0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := performEnglishWordAPIRequest(t, http.MethodGet, tc.target, token, nil, api.FindEnglishWord)
			decoded := decodeEnglishWordAPIResp(t, resp)
			if decoded.Code == 0 {
				t.Fatalf("expected parameter failure, got success msg=%s", decoded.Msg)
			}
			if !strings.Contains(decoded.Msg, "ID参数错误") {
				t.Fatalf("expected ID parameter error message, got=%s", decoded.Msg)
			}
		})
	}
}

func TestEnglishWordAPI_FindEnglishWord_NotFound(t *testing.T) {
	_, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	resp := performEnglishWordAPIRequest(t, http.MethodGet, "/englishLearning/word/findWord?ID=999999", token, nil, api.FindEnglishWord)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected findWord failure for missing record, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "查询失败") {
		t.Fatalf("expected query failure message, got=%s", decoded.Msg)
	}
}

func TestEnglishWordAPI_UpdateEnglishWord_Success(t *testing.T) {
	wordID, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	body := []byte(fmt.Sprintf(`{"ID":%d,"word":"updated_word_%s","phoneticUs":"/u/","phoneticUk":"/k/","explanation":"{}"}`,
		wordID,
		uuid.NewString(),
	))

	resp := performEnglishWordAPIRequest(t, http.MethodPut, "/englishLearning/word/update", token, body, api.UpdateEnglishWord)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected update success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	findResp := performEnglishWordAPIRequest(t, http.MethodGet, fmt.Sprintf("/englishLearning/word/findWord?ID=%d", wordID), token, nil, api.FindEnglishWord)
	findDecoded := decodeEnglishWordAPIResp(t, findResp)
	if findDecoded.Code != 0 {
		t.Fatalf("expected find after update success, code=%d msg=%s", findDecoded.Code, findDecoded.Msg)
	}

	var word model.EnglishWord
	if err := json.Unmarshal(findDecoded.Data, &word); err != nil {
		t.Fatalf("decode word data failed: %v raw=%s", err, string(findDecoded.Data))
	}
	if !strings.HasPrefix(word.Word, "updated_word_") {
		t.Fatalf("expected updated word prefix, got=%s", word.Word)
	}
}

func TestEnglishWordAPI_DeleteEnglishWord_Success(t *testing.T) {
	wordID, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	resp := performEnglishWordAPIRequest(t, http.MethodDelete, fmt.Sprintf("/englishLearning/word/delete?ID=%d", wordID), token, nil, api.DeleteEnglishWord)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected delete success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	findResp := performEnglishWordAPIRequest(t, http.MethodGet, fmt.Sprintf("/englishLearning/word/findWord?ID=%d", wordID), token, nil, api.FindEnglishWord)
	findDecoded := decodeEnglishWordAPIResp(t, findResp)
	if findDecoded.Code == 0 {
		t.Fatalf("expected find failure after delete, got success msg=%s", findDecoded.Msg)
	}
}

func TestEnglishWordAPI_RegenerateWordAudio_ParamError(t *testing.T) {
	_, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	resp := performEnglishWordAPIRequest(t, http.MethodPost, "/englishLearning/word/regenerateAudio", token, []byte(`{"word":"x"}`), api.RegenerateWordAudio)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected parameter failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "参数错误") {
		t.Fatalf("expected parameter error message, got=%s", decoded.Msg)
	}
}

func TestEnglishWordAPI_RegenerateWordAudio_Success(t *testing.T) {
	wordID, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	ttsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"audioUrl":"https://tts.example.com/audio.mp3"}}`))
	}))
	defer ttsServer.Close()

	upsertEnglishWordSysConfig(t, "learning_tts_enabled", "true")
	upsertEnglishWordSysConfig(t, "learning_tts_provider_url", ttsServer.URL)
	upsertEnglishWordSysConfig(t, "learning_tts_timeout_ms", "3000")

	body := []byte(fmt.Sprintf(`{"ID":%d}`, wordID))
	resp := performEnglishWordAPIRequest(t, http.MethodPost, "/englishLearning/word/regenerateAudio", token, body, api.RegenerateWordAudio)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected regenerate success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	findResp := performEnglishWordAPIRequest(t, http.MethodGet, fmt.Sprintf("/englishLearning/word/findWord?ID=%d", wordID), token, nil, api.FindEnglishWord)
	findDecoded := decodeEnglishWordAPIResp(t, findResp)
	if findDecoded.Code != 0 {
		t.Fatalf("expected find success after regenerate, code=%d msg=%s", findDecoded.Code, findDecoded.Msg)
	}

	var word model.EnglishWord
	if err := json.Unmarshal(findDecoded.Data, &word); err != nil {
		t.Fatalf("decode word data failed: %v raw=%s", err, string(findDecoded.Data))
	}
	if word.AudioUS == "" || word.AudioUK == "" {
		t.Fatalf("expected both audio fields to be generated, got us=%s uk=%s", word.AudioUS, word.AudioUK)
	}
}

func TestEnglishWordAPI_RegenerateWordAudio_BothDisabled(t *testing.T) {
	wordID, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	body := []byte(fmt.Sprintf(`{"ID":%d,"regenerateUs":false,"regenerateUk":false}`, wordID))
	resp := performEnglishWordAPIRequest(t, http.MethodPost, "/englishLearning/word/regenerateAudio", token, body, api.RegenerateWordAudio)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected regenerate failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "重生成失败") {
		t.Fatalf("expected regenerate failure message, got=%s", decoded.Msg)
	}
}

func TestEnglishWordAPI_PreflightTTS_Success(t *testing.T) {
	_, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	ttsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"audioUrl":"https://tts.example.com/preflight.mp3"}}`))
	}))
	defer ttsServer.Close()

	upsertEnglishWordSysConfig(t, "learning_tts_enabled", "true")
	upsertEnglishWordSysConfig(t, "learning_tts_provider_url", ttsServer.URL)
	upsertEnglishWordSysConfig(t, "learning_tts_timeout_ms", "3000")

	body := []byte(`{"word":"future"}`)
	resp := performEnglishWordAPIRequest(t, http.MethodPost, "/englishLearning/word/preflightTTS", token, body, api.PreflightTTS)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected preflight success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var data englishWordPreflightData
	if err := json.Unmarshal(decoded.Data, &data); err != nil {
		t.Fatalf("decode preflight data failed: %v raw=%s", err, string(decoded.Data))
	}
	if data.Word != "future" {
		t.Fatalf("expected word=future, got=%s", data.Word)
	}
	if data.AudioUS == "" || data.AudioUK == "" {
		t.Fatalf("expected both audio urls, got us=%s uk=%s", data.AudioUS, data.AudioUK)
	}
}

func TestEnglishWordAPI_PreflightTTS_BothDisabled(t *testing.T) {
	_, token := setupEnglishWordAPITest(t)
	api := EnglishWordApi{}

	body := []byte(`{"checkUs":false,"checkUk":false}`)
	resp := performEnglishWordAPIRequest(t, http.MethodPost, "/englishLearning/word/preflightTTS", token, body, api.PreflightTTS)
	decoded := decodeEnglishWordAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected preflight failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "预检失败") {
		t.Fatalf("expected preflight failure message, got=%s", decoded.Msg)
	}
}
