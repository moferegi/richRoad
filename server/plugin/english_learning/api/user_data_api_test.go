package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	requestModel "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userDataAPIResp struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

type watchHistoryListRespItem struct {
	EpisodeID    uint    `json:"episodeId"`
	EpisodeName  string  `json:"episodeName"`
	SeriesID     uint    `json:"seriesId"`
	SeriesName   string  `json:"seriesName"`
	CoverURL     string  `json:"coverUrl"`
	ProgressSecs float64 `json:"progressSecs"`
}

type watchHistoryListRespData struct {
	List     []watchHistoryListRespItem `json:"list"`
	Total    int64                      `json:"total"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"pageSize"`
}

func setupUserDataAPITest(t *testing.T) (token string, episode1ID uint, episode2ID uint) {
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
	global.GVA_CONFIG.JWT.SigningKey = "user-data-api-test-sign-key"
	global.GVA_CONFIG.JWT.ExpiresTime = "2h"
	global.GVA_CONFIG.JWT.BufferTime = "10m"
	global.GVA_CONFIG.JWT.Issuer = "user-data-api-test"

	dsn := fmt.Sprintf("file:user_data_api_test_%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	if err = db.AutoMigrate(&model.VideoSeries{}, &model.VideoEpisode{}, &model.UserWatchHistory{}); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	series1 := model.VideoSeries{Name: `{"zh":"英语剧集A","en":"Series A"}`, CoverUrl: "https://example.com/cover-a.jpg"}
	series2 := model.VideoSeries{Name: `{"zh":"英语剧集B","en":"Series B"}`, CoverUrl: "https://example.com/cover-b.jpg"}
	if err = db.Create(&series1).Error; err != nil {
		t.Fatalf("seed series1 failed: %v", err)
	}
	if err = db.Create(&series2).Error; err != nil {
		t.Fatalf("seed series2 failed: %v", err)
	}

	episode1 := model.VideoEpisode{SeriesID: series1.ID, Name: `{"zh":"第1集","en":"EP1"}`, VideoUrl: "https://example.com/video-1.m3u8"}
	episode2 := model.VideoEpisode{SeriesID: series2.ID, Name: `{"zh":"第2集","en":"EP2"}`, VideoUrl: "https://example.com/video-2.m3u8"}
	if err = db.Create(&episode1).Error; err != nil {
		t.Fatalf("seed episode1 failed: %v", err)
	}
	if err = db.Create(&episode2).Error; err != nil {
		t.Fatalf("seed episode2 failed: %v", err)
	}

	userID := uint(9101)
	history1 := model.UserWatchHistory{UserID: userID, EpisodeID: episode1.ID, ProgressSecs: 35.5}
	history2 := model.UserWatchHistory{UserID: userID, EpisodeID: episode2.ID, ProgressSecs: 82.0}
	otherUserHistory := model.UserWatchHistory{UserID: 9102, EpisodeID: episode2.ID, ProgressSecs: 12.0}

	if err = db.Create(&history1).Error; err != nil {
		t.Fatalf("seed history1 failed: %v", err)
	}
	if err = db.Create(&history2).Error; err != nil {
		t.Fatalf("seed history2 failed: %v", err)
	}
	if err = db.Create(&otherUserHistory).Error; err != nil {
		t.Fatalf("seed otherUserHistory failed: %v", err)
	}

	if err = db.Model(&history1).Update("updated_at", time.Now().Add(-2*time.Hour)).Error; err != nil {
		t.Fatalf("update history1 updated_at failed: %v", err)
	}
	if err = db.Model(&history2).Update("updated_at", time.Now().Add(-10*time.Minute)).Error; err != nil {
		t.Fatalf("update history2 updated_at failed: %v", err)
	}

	global.GVA_DB = db

	j := utils.NewJWT()
	claims := j.CreateClaims(requestModel.BaseClaims{
		UUID:        uuid.New(),
		ID:          userID,
		Username:    "watch_history_user",
		NickName:    "watch_history_user",
		AuthorityId: 8080,
	})
	token, err = j.CreateToken(claims)
	if err != nil {
		t.Fatalf("create token failed: %v", err)
	}

	return token, episode1.ID, episode2.ID
}

func performUserDataAPIRequest(t *testing.T, method, target, token string, body []byte, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("x-token", token)
	if method == http.MethodPost || method == http.MethodDelete {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req
	handler(c)
	return recorder
}

func decodeUserDataAPIResp(t *testing.T, recorder *httptest.ResponseRecorder) userDataAPIResp {
	t.Helper()

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected http 200, got %d, body=%s", recorder.Code, recorder.Body.String())
	}

	var resp userDataAPIResp
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response failed: %v, body=%s", err, recorder.Body.String())
	}
	return resp
}

func TestUserDataAPI_GetWatchHistoryList_Success(t *testing.T) {
	token, episode1ID, episode2ID := setupUserDataAPITest(t)
	api := UserDataApi{}

	resp := performUserDataAPIRequest(t, http.MethodGet, "/englishLearning/userData/getWatchHistoryList?page=1&pageSize=10", token, nil, api.GetWatchHistoryList)
	decoded := decodeUserDataAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var data watchHistoryListRespData
	if err := json.Unmarshal(decoded.Data, &data); err != nil {
		t.Fatalf("decode data failed: %v raw=%s", err, string(decoded.Data))
	}

	if data.Total != 2 {
		t.Fatalf("expected total=2, got=%d", data.Total)
	}
	if len(data.List) != 2 {
		t.Fatalf("expected list size=2, got=%d", len(data.List))
	}
	if data.Page != 1 || data.PageSize != 10 {
		t.Fatalf("expected page/pageSize=1/10, got=%d/%d", data.Page, data.PageSize)
	}

	if data.List[0].EpisodeID != episode2ID {
		t.Fatalf("expected first item episodeId=%d by updatedAt desc, got=%d", episode2ID, data.List[0].EpisodeID)
	}
	if data.List[1].EpisodeID != episode1ID {
		t.Fatalf("expected second item episodeId=%d, got=%d", episode1ID, data.List[1].EpisodeID)
	}

	if data.List[0].SeriesID == 0 || strings.TrimSpace(data.List[0].SeriesName) == "" || strings.TrimSpace(data.List[0].EpisodeName) == "" {
		t.Fatalf("expected series/episode details in first item, got=%+v", data.List[0])
	}
}

func TestUserDataAPI_GetWatchHistoryList_ParamError(t *testing.T) {
	token, _, _ := setupUserDataAPITest(t)
	api := UserDataApi{}

	resp := performUserDataAPIRequest(t, http.MethodGet, "/englishLearning/userData/getWatchHistoryList?page=invalid&pageSize=10", token, nil, api.GetWatchHistoryList)
	decoded := decodeUserDataAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected parameter failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "参数错误") {
		t.Fatalf("expected parameter error message, got=%s", decoded.Msg)
	}
}

func TestUserDataAPI_GetWatchHistoryList_AuthFailure(t *testing.T) {
	_, _, _ = setupUserDataAPITest(t)
	api := UserDataApi{}

	resp := performUserDataAPIRequest(t, http.MethodGet, "/englishLearning/userData/getWatchHistoryList?page=1&pageSize=10", "", nil, api.GetWatchHistoryList)
	decoded := decodeUserDataAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected auth failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "获取用户信息失败") {
		t.Fatalf("expected auth failure message, got=%s", decoded.Msg)
	}
}

func TestUserDataAPI_GetWatchHistoryList_NormalizedPage(t *testing.T) {
	token, _, _ := setupUserDataAPITest(t)
	api := UserDataApi{}

	resp := performUserDataAPIRequest(t, http.MethodGet, "/englishLearning/userData/getWatchHistoryList", token, nil, api.GetWatchHistoryList)
	decoded := decodeUserDataAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var data watchHistoryListRespData
	if err := json.Unmarshal(decoded.Data, &data); err != nil {
		t.Fatalf("decode data failed: %v raw=%s", err, string(decoded.Data))
	}

	if data.Page != 1 || data.PageSize != 10 {
		t.Fatalf("expected normalized page/pageSize=1/10, got=%d/%d", data.Page, data.PageSize)
	}
}
