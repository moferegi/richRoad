package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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

type entitlementAPIResp struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

type entitlementItem struct {
	ID           uint   `json:"ID"`
	UserID       uint   `json:"userId"`
	ResourceType string `json:"resourceType"`
	ResourceID   uint   `json:"resourceId"`
}

type entitlementPageData struct {
	List     []entitlementItem `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

func setupEntitlementAPITest(t *testing.T) (string, uint) {
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
	global.GVA_CONFIG.JWT.SigningKey = "entitlement-api-test-sign-key"
	global.GVA_CONFIG.JWT.ExpiresTime = "2h"
	global.GVA_CONFIG.JWT.BufferTime = "10m"
	global.GVA_CONFIG.JWT.Issuer = "entitlement-api-test"

	dsn := fmt.Sprintf("file:entitlement_api_test_%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	if err = db.AutoMigrate(&model.UserLearningEntitlement{}, &model.VideoSeries{}, &model.VideoEpisode{}); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	series := model.VideoSeries{Name: `{"zh":"授权测试剧集"}`}
	if err = db.Create(&series).Error; err != nil {
		t.Fatalf("seed series failed: %v", err)
	}

	global.GVA_DB = db

	j := utils.NewJWT()
	claims := j.CreateClaims(requestModel.BaseClaims{
		UUID:        uuid.New(),
		ID:          8801,
		Username:    "entitlement_admin",
		NickName:    "entitlement_admin",
		AuthorityId: 9999,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		t.Fatalf("create token failed: %v", err)
	}

	return token, series.ID
}

func performEntitlementAPIRequest(t *testing.T, method, target, token string, body []byte, handler gin.HandlerFunc) *httptest.ResponseRecorder {
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

func decodeEntitlementAPIResp(t *testing.T, recorder *httptest.ResponseRecorder) entitlementAPIResp {
	t.Helper()

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected http 200, got %d, body=%s", recorder.Code, recorder.Body.String())
	}

	var resp entitlementAPIResp
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response failed: %v body=%s", err, recorder.Body.String())
	}
	return resp
}

func TestUserLearningAssetApi_GrantRevokeEntitlement_Success(t *testing.T) {
	token, seriesID := setupEntitlementAPITest(t)
	api := UserLearningAssetApi{}

	grantPayload := map[string]interface{}{
		"userId":       9908,
		"resourceType": model.ResourceTypeVideoSeries,
		"resourceId":   seriesID,
		"remark":       "api grant",
	}
	grantBody, _ := json.Marshal(grantPayload)
	grantResp := performEntitlementAPIRequest(t, http.MethodPost, "/englishLearning/asset/grantEntitlement", token, grantBody, api.GrantEntitlement)
	grantDecoded := decodeEntitlementAPIResp(t, grantResp)
	if grantDecoded.Code != 0 {
		t.Fatalf("expected grant success, code=%d msg=%s", grantDecoded.Code, grantDecoded.Msg)
	}

	listResp := performEntitlementAPIRequest(t, http.MethodGet, "/englishLearning/asset/getEntitlementList?userId=9908&page=1&pageSize=10", token, nil, api.GetEntitlementList)
	listDecoded := decodeEntitlementAPIResp(t, listResp)
	if listDecoded.Code != 0 {
		t.Fatalf("expected list success, code=%d msg=%s", listDecoded.Code, listDecoded.Msg)
	}

	var pageData entitlementPageData
	if err := json.Unmarshal(listDecoded.Data, &pageData); err != nil {
		t.Fatalf("decode list data failed: %v raw=%s", err, string(listDecoded.Data))
	}
	if pageData.Total != 1 || len(pageData.List) != 1 {
		t.Fatalf("expected one entitlement, total=%d size=%d", pageData.Total, len(pageData.List))
	}
	if pageData.List[0].ResourceType != model.ResourceTypeVideoSeries || pageData.List[0].ResourceID != seriesID {
		t.Fatalf("unexpected entitlement row: %+v", pageData.List[0])
	}

	revokePayload := map[string]interface{}{
		"userId":       9908,
		"resourceType": model.ResourceTypeVideoSeries,
		"resourceId":   seriesID,
	}
	revokeBody, _ := json.Marshal(revokePayload)
	revokeResp := performEntitlementAPIRequest(t, http.MethodDelete, "/englishLearning/asset/revokeEntitlement", token, revokeBody, api.RevokeEntitlement)
	revokeDecoded := decodeEntitlementAPIResp(t, revokeResp)
	if revokeDecoded.Code != 0 {
		t.Fatalf("expected revoke success, code=%d msg=%s", revokeDecoded.Code, revokeDecoded.Msg)
	}

	listRespAfter := performEntitlementAPIRequest(t, http.MethodGet, "/englishLearning/asset/getEntitlementList?userId=9908&page=1&pageSize=10", token, nil, api.GetEntitlementList)
	listDecodedAfter := decodeEntitlementAPIResp(t, listRespAfter)
	if listDecodedAfter.Code != 0 {
		t.Fatalf("expected list success after revoke, code=%d msg=%s", listDecodedAfter.Code, listDecodedAfter.Msg)
	}

	var pageDataAfter entitlementPageData
	if err := json.Unmarshal(listDecodedAfter.Data, &pageDataAfter); err != nil {
		t.Fatalf("decode list data after revoke failed: %v raw=%s", err, string(listDecodedAfter.Data))
	}
	if pageDataAfter.Total != 0 || len(pageDataAfter.List) != 0 {
		t.Fatalf("expected no entitlement after revoke, total=%d size=%d", pageDataAfter.Total, len(pageDataAfter.List))
	}
}
