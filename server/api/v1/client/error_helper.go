package client

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
)

func failClientWithErr(c *gin.Context, err error) {
	response.FailWithMessage(resolveClientErrMessage(c, err, "fail"), c)
}

func failClientWithKey(c *gin.Context, key string) {
	key = strings.TrimSpace(key)
	if key == "" || !i18n.HasKey(key) {
		response.FailWithMessage(i18n.T(c, "fail"), c)
		return
	}
	response.FailWithMessage(i18n.T(c, key), c)
}

func isClientAdminAuthority(authorityID uint) bool {
	return authorityID == 888 || authorityID == 8881
}

func resolveClientErrMessage(c *gin.Context, err error, fallbackKey string) string {
	fallbackKey = strings.TrimSpace(fallbackKey)
	if fallbackKey == "" {
		fallbackKey = "fail"
	}
	if err == nil {
		return i18n.T(c, fallbackKey)
	}
	key := strings.TrimSpace(err.Error())
	if key != "" && i18n.HasKey(key) {
		return i18n.T(c, key)
	}
	return i18n.T(c, fallbackKey)
}
