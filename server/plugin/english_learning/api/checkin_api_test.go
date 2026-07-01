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

type checkinAPIResp struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

type pointRecordRespItem struct {
	ID            uint      `json:"ID"`
	UserID        uint      `json:"userId"`
	ChangeType    string    `json:"changeType"`
	PointChange   int       `json:"pointChange"`
	OperationType string    `json:"operationType"`
	Reason        string    `json:"reason"`
	CurrentPoints int       `json:"currentPoints"`
	CreatedAt     time.Time `json:"CreatedAt"`
}

type pointRecordRespData struct {
	List     []pointRecordRespItem `json:"list"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

type checkinRecordRespItem struct {
	ID          uint      `json:"ID"`
	UserID      uint      `json:"userId"`
	CheckinDate time.Time `json:"checkinDate"`
	PointAward  int       `json:"pointAward"`
}

type checkinRecordRespData struct {
	List     []checkinRecordRespItem `json:"list"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"pageSize"`
}

func setupCheckinAPITest(t *testing.T) string {
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
	global.GVA_CONFIG.JWT.SigningKey = "checkin-api-test-sign-key"
	global.GVA_CONFIG.JWT.ExpiresTime = "2h"
	global.GVA_CONFIG.JWT.BufferTime = "10m"
	global.GVA_CONFIG.JWT.Issuer = "checkin-api-test"

	dsn := fmt.Sprintf("file:checkin_api_test_%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	if err = db.AutoMigrate(&model.EnglishPointRecord{}, &model.CheckinRecord{}); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	userID := uint(9201)
	otherUserID := uint(9202)
	now := time.Now()
	records := []model.EnglishPointRecord{
		{
			UserID:        userID,
			ChangeType:    "increase",
			PointChange:   100,
			OperationType: "checkin_reward",
			Reason:        "reason_signInReward",
			CurrentPoints: 300,
			Remark:        "",
			GVA_MODEL: global.GVA_MODEL{
				CreatedAt: now.Add(-2 * time.Hour),
			},
		},
		{
			UserID:        userID,
			ChangeType:    "decrease",
			PointChange:   -200,
			OperationType: "point_exchange",
			Reason:        "reason_learningExchange",
			CurrentPoints: 100,
			Remark:        "minutes:2",
			GVA_MODEL: global.GVA_MODEL{
				CreatedAt: now.Add(-1 * time.Hour),
			},
		},
		{
			UserID:        otherUserID,
			ChangeType:    "increase",
			PointChange:   50,
			OperationType: "checkin_reward",
			Reason:        "reason_signInReward",
			CurrentPoints: 50,
			Remark:        "",
			GVA_MODEL: global.GVA_MODEL{
				CreatedAt: now.Add(-30 * time.Minute),
			},
		},
	}
	if err = db.Create(&records).Error; err != nil {
		t.Fatalf("seed records failed: %v", err)
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	checkinRecords := []model.CheckinRecord{
		{UserID: userID, CheckinDate: today.AddDate(0, 0, -2), PointAward: 100},
		{UserID: userID, CheckinDate: today.AddDate(0, 0, -1), PointAward: 105},
		{UserID: otherUserID, CheckinDate: today, PointAward: 100},
	}
	if err = db.Create(&checkinRecords).Error; err != nil {
		t.Fatalf("seed checkin records failed: %v", err)
	}

	global.GVA_DB = db

	j := utils.NewJWT()
	claims := j.CreateClaims(requestModel.BaseClaims{
		UUID:        uuid.New(),
		ID:          userID,
		Username:    "checkin_api_user",
		NickName:    "checkin_api_user",
		AuthorityId: 8080,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		t.Fatalf("create token failed: %v", err)
	}

	return token
}

func performCheckinAPIRequest(t *testing.T, method, target, token string, body []byte, handler gin.HandlerFunc) *httptest.ResponseRecorder {
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

func decodeCheckinAPIResp(t *testing.T, recorder *httptest.ResponseRecorder) checkinAPIResp {
	t.Helper()

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected http 200, got %d, body=%s", recorder.Code, recorder.Body.String())
	}

	var resp checkinAPIResp
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response failed: %v body=%s", err, recorder.Body.String())
	}
	return resp
}

func TestCheckinAPI_GetPointRecordList_Success(t *testing.T) {
	token := setupCheckinAPITest(t)
	api := CheckinApi{}

	resp := performCheckinAPIRequest(t, http.MethodGet, "/englishLearning/checkin/getPointRecordList?page=1&pageSize=10", token, nil, api.GetPointRecordList)
	decoded := decodeCheckinAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var data pointRecordRespData
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

	if data.List[0].OperationType != "point_exchange" {
		t.Fatalf("expected latest record is point_exchange, got=%s", data.List[0].OperationType)
	}
	if data.List[1].OperationType != "checkin_reward" {
		t.Fatalf("expected second record is checkin_reward, got=%s", data.List[1].OperationType)
	}
}

func TestCheckinAPI_GetPointRecordList_ParamError(t *testing.T) {
	token := setupCheckinAPITest(t)
	api := CheckinApi{}

	resp := performCheckinAPIRequest(t, http.MethodGet, "/englishLearning/checkin/getPointRecordList?page=bad&pageSize=10", token, nil, api.GetPointRecordList)
	decoded := decodeCheckinAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected parameter failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "参数错误") {
		t.Fatalf("expected parameter error message, got=%s", decoded.Msg)
	}
}

func TestCheckinAPI_GetPointRecordList_AuthFailure(t *testing.T) {
	_ = setupCheckinAPITest(t)
	api := CheckinApi{}

	resp := performCheckinAPIRequest(t, http.MethodGet, "/englishLearning/checkin/getPointRecordList?page=1&pageSize=10", "", nil, api.GetPointRecordList)
	decoded := decodeCheckinAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected auth failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "获取用户信息失败") {
		t.Fatalf("expected auth failure message, got=%s", decoded.Msg)
	}
}

func TestCheckinAPI_GetPointRecordList_NormalizedPage(t *testing.T) {
	token := setupCheckinAPITest(t)
	api := CheckinApi{}

	resp := performCheckinAPIRequest(t, http.MethodGet, "/englishLearning/checkin/getPointRecordList", token, nil, api.GetPointRecordList)
	decoded := decodeCheckinAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var data pointRecordRespData
	if err := json.Unmarshal(decoded.Data, &data); err != nil {
		t.Fatalf("decode data failed: %v raw=%s", err, string(decoded.Data))
	}

	if data.Page != 1 || data.PageSize != 10 {
		t.Fatalf("expected normalized page/pageSize=1/10, got=%d/%d", data.Page, data.PageSize)
	}
}

func TestCheckinAPI_GetCheckinRecordList_Success(t *testing.T) {
	token := setupCheckinAPITest(t)
	api := CheckinApi{}

	resp := performCheckinAPIRequest(t, http.MethodGet, "/englishLearning/checkin/getCheckinRecordList?page=1&pageSize=10", token, nil, api.GetCheckinRecordList)
	decoded := decodeCheckinAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var data checkinRecordRespData
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
	if !data.List[0].CheckinDate.After(data.List[1].CheckinDate) && !data.List[0].CheckinDate.Equal(data.List[1].CheckinDate) {
		t.Fatalf("expected checkinDate desc order, first=%s second=%s", data.List[0].CheckinDate, data.List[1].CheckinDate)
	}
}

func TestCheckinAPI_GetCheckinRecordList_ParamError(t *testing.T) {
	token := setupCheckinAPITest(t)
	api := CheckinApi{}

	resp := performCheckinAPIRequest(t, http.MethodGet, "/englishLearning/checkin/getCheckinRecordList?page=bad&pageSize=10", token, nil, api.GetCheckinRecordList)
	decoded := decodeCheckinAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected parameter failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "参数错误") {
		t.Fatalf("expected parameter error message, got=%s", decoded.Msg)
	}
}

func TestCheckinAPI_GetCheckinRecordList_AuthFailure(t *testing.T) {
	_ = setupCheckinAPITest(t)
	api := CheckinApi{}

	resp := performCheckinAPIRequest(t, http.MethodGet, "/englishLearning/checkin/getCheckinRecordList?page=1&pageSize=10", "", nil, api.GetCheckinRecordList)
	decoded := decodeCheckinAPIResp(t, resp)
	if decoded.Code == 0 {
		t.Fatalf("expected auth failure, got success msg=%s", decoded.Msg)
	}
	if !strings.Contains(decoded.Msg, "获取用户信息失败") {
		t.Fatalf("expected auth failure message, got=%s", decoded.Msg)
	}
}

func TestCheckinAPI_GetCheckinRecordList_NormalizedPage(t *testing.T) {
	token := setupCheckinAPITest(t)
	api := CheckinApi{}

	resp := performCheckinAPIRequest(t, http.MethodGet, "/englishLearning/checkin/getCheckinRecordList", token, nil, api.GetCheckinRecordList)
	decoded := decodeCheckinAPIResp(t, resp)
	if decoded.Code != 0 {
		t.Fatalf("expected success, code=%d msg=%s", decoded.Code, decoded.Msg)
	}

	var data checkinRecordRespData
	if err := json.Unmarshal(decoded.Data, &data); err != nil {
		t.Fatalf("decode data failed: %v raw=%s", err, string(decoded.Data))
	}

	if data.Page != 1 || data.PageSize != 10 {
		t.Fatalf("expected normalized page/pageSize=1/10, got=%d/%d", data.Page, data.PageSize)
	}
}
