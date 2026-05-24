package shop

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
)

func TestCreateOrderRequiresOrderAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/order/createOrder", strings.NewReader(`{"userID":1,"status":"1","totalPrice":1}`))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("claims", &systemReq.CustomClaims{BaseClaims: systemReq.BaseClaims{ID: 1, AuthorityId: 8080}})

	(&OrderApi{}).CreateOrder(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected response status 200, got %d", recorder.Code)
	}
	if strings.Contains(recorder.Body.String(), "createSuccess") {
		t.Fatalf("non-admin order creation unexpectedly succeeded: %s", recorder.Body.String())
	}
}
