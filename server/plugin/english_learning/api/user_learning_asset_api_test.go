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

type userLearningAssetAPIResp struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

type freeTimeRecordRespItem struct {
	ID                 uint      `json:"ID"`
	UserID             uint      `json:"userId"`
	ChangeType         string    `json:"changeType"`
	MinuteChange       int       `json:"minuteChange"`
	OperationType      string    `json:"operationType"`
	Reason             string    `json:"reason"`
	CurrentFreeMinutes int       `json:"currentFreeMinutes"`
	RelatedEpisodeID   uint      `json:"relatedEpisodeId"`
	CreatedAt          time.Time `json:"CreatedAt"`
}

type freeTimeRecordRespData struct {
	List     []freeTimeRecordRespItem `json:"list"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"pageSize"`
}

func setupUserLearningAssetAPITest(t *testing.T) string {
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
	global.GVA_CONFIG.JWT.SigningKey = "user-learning-asset-api-test-sign-key"
	global.GVA_CONFIG.JWT.ExpiresTime = "2h"
	global.GVA_CONFIG.JWT.BufferTime = "10m"
	global.GVA_CONFIG.JWT.Issuer = "user-learning-asset-api-test"

	dsn := fmt.Sprintf("file:user_learning_asset_api_test_%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	if err = db.AutoMigrate(&model.EnglishFreeTimeRecord{}); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	userID := uint(9501)
	otherUserID := uint(9502)
	now := time.Now()
	records := []model.EnglishFreeTimeRecord{
		{
			UserID:             userID,
			ChangeType:         "increase",
			MinuteChange:       2,
			OperationType:      "point_exchange",
			Reason:             "reason_learningExchange",
			CurrentFreeMinutes: 2,
			RelatedEpisodeID:   0,
			Remark:             "points:200",
			GVA_MODEL:          global.GVA_MODEL{CreatedAt: now.Add(-2 * time.Hour)},
		},
		{
			UserID:             userID,
			ChangeType:         "decrease",
			MinuteChange:       -1,
			OperationType:      "heartbeat_consume",
			Reason:             "reason_learningWatchConsume",
			CurrentFreeMinutes: 1,
			RelatedEpisodeID:   811,
			Remark:             "episode:811,usingSecs:68",
			GVA_MODEL:          global.GVA_MODEL{CreatedAt: now.Add(-1 * time.Hour)},
		},
		{
			UserID:             otherUserID,
			ChangeType:         "increase",
			MinuteChange:       3,
			OperationType:      "point_exchange",
			Reason:             "reason_learningExchange",
			CurrentFreeMinutes: 3,
			RelatedEpisodeID:   0,
			Remark:             "points:300",
			GVA_MODEL:          global.GVA_MODEL{CreatedAt: now.Add(-30 * time.Minute)},
		},
	}
	if err = db.Create(&records).Error; err != nil {
		t.Fatalf("seed records failed: %v", err)
	}

	global.GVA_DB = db

	j := utils.NewJWT()
	claims := j.CreateClaims(requestModel.BaseClaims{
		UUID:        uuid.New(),
		ID:          userID,
		Username:    "asset_api_user",
		NickName:    "asset_api_user",
		AuthorityId: 8080,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		t.Fatalf("create token failed: %v", err)
	}

	return token
}

func performUserLearningAssetAPIRequest(t *testing.T, method, target, token string, body []byte, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("x-token", token)
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req
	handler(c)
	return recorder
}

func decodeUserLearningAssetAPIResp(t *testing.T, recorder *httptest.ResponseRecorder) userLearningAssetAPIResp {
	t.Helper()

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected http 200, got %d, body=%s", recorder.Code, recorder.Body.String())
	}

	var resp userLearningAssetAPIResp
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response failed: %v, body=%s", err, recorder.Body.String())
	}
	return resp
}

func TestUserLearningAssetAPI_GetFreeTimeRecordList_Success(t *testing.T) {
	token := setupUserLearningAssetAPITest(t)
	api := UserLearningAssetApi{}

	resp := performUserLearningAssetAPIRequest(t, http.MethodGet, "/englishLearning/asset/getFreeTimeRecordList?page=1&pageSize=10", token, nil, api.GetFreeTimeRecordList)
	decoded := decodeUserLearningAssetAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var data freeTimeRecordRespData
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

	if data.List[0].OperationType != "heartbeat_consume" {
		t.Fatalf("expected first operationType=heartbeat_consume, got=%s", data.List[0].OperationType)
	}
	if data.List[1].OperationType != "point_exchange" {
		t.Fatalf("expected second operationType=point_exchange, got=%s", data.List[1].OperationType)
	}
}

func TestUserLearningAssetAPI_GetFreeTimeRecordList_ParamError(t *testing.T) {
	token := setupUserLearningAssetAPITest(t)
	api := UserLearningAssetApi{}

	resp := performUserLearningAssetAPIRequest(t, http.MethodGet, "/englishLearning/asset/getFreeTimeRecordList?page=bad&pageSize=10", token, nil, api.GetFreeTimeRecordList)
	decoded := decodeUserLearningAssetAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected parameter failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "参数错误") {
		t.Fatalf("expected parameter error message, got=%s", decoded.Msg)
	}
}

func TestUserLearningAssetAPI_GetFreeTimeRecordList_AuthFailure(t *testing.T) {
	_ = setupUserLearningAssetAPITest(t)
	api := UserLearningAssetApi{}

	resp := performUserLearningAssetAPIRequest(t, http.MethodGet, "/englishLearning/asset/getFreeTimeRecordList?page=1&pageSize=10", "", nil, api.GetFreeTimeRecordList)
	decoded := decodeUserLearningAssetAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected auth failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "获取用户信息失败") {
		t.Fatalf("expected auth failure message, got=%s", decoded.Msg)
	}
}

func TestUserLearningAssetAPI_GetFreeTimeRecordList_NormalizedPage(t *testing.T) {
	token := setupUserLearningAssetAPITest(t)
	api := UserLearningAssetApi{}

	resp := performUserLearningAssetAPIRequest(t, http.MethodGet, "/englishLearning/asset/getFreeTimeRecordList", token, nil, api.GetFreeTimeRecordList)
	decoded := decodeUserLearningAssetAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var data freeTimeRecordRespData
	if err := json.Unmarshal(decoded.Data, &data); err != nil {
		t.Fatalf("decode data failed: %v raw=%s", err, string(decoded.Data))
	}
	if data.Page != 1 || data.PageSize != 10 {
		t.Fatalf("expected normalized page/pageSize=1/10, got=%d/%d", data.Page, data.PageSize)
	}
}
