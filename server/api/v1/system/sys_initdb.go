package system

import (
	"net"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type DBApi struct{}

const (
	initAPIEnabledConfigKey = "security_init_api_enabled"
	initAPIEnabledEnvKey    = "CS_INIT_API_ENABLED"

	initAPIPrivateOnlyConfigKey = "security_init_api_private_network_only"
	initAPIPrivateOnlyEnvKey    = "CS_INIT_API_PRIVATE_NETWORK_ONLY"
)

func isPrivateOrLoopbackIP(ip string) bool {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return false
	}
	if parsed.IsLoopback() || parsed.IsPrivate() || parsed.IsLinkLocalUnicast() {
		return true
	}
	return false
}

func (i *DBApi) checkInitAPIAccess(c *gin.Context) bool {
	defaultEnabled := global.GVA_DB == nil
	if !utils.GetBoolSetting(initAPIEnabledConfigKey, initAPIEnabledEnvKey, defaultEnabled) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return false
	}

	if !utils.GetBoolSetting(initAPIPrivateOnlyConfigKey, initAPIPrivateOnlyEnvKey, true) {
		return true
	}

	if !isPrivateOrLoopbackIP(c.ClientIP()) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return false
	}

	return true
}

// InitDB
// @Tags     InitDB
// @Summary  初始化用户数据库
// @Produce  application/json
// @Param    data  body      request.InitDB                  true  "初始化数据库参数"
// @Success  200   {object}  response.Response{data=string}  "初始化用户数据库"
// @Router   /init/initdb [post]
func (i *DBApi) InitDB(c *gin.Context) {
	if !i.checkInitAPIAccess(c) {
		return
	}

	if global.GVA_DB != nil {
		global.GVA_LOG.Error("已存在数据库配置!")
		response.FailWithMessage("已存在数据库配置", c)
		return
	}
	var dbInfo request.InitDB
	if err := c.ShouldBindJSON(&dbInfo); err != nil {
		global.GVA_LOG.Error("参数校验不通过!", zap.Error(err))
		response.FailWithMessage("参数校验不通过", c)
		return
	}
	if err := initDBService.InitDB(dbInfo); err != nil {
		global.GVA_LOG.Error("自动创建数据库失败!", zap.Error(err))
		response.FailWithMessage("自动创建数据库失败，请查看后台日志，检查后在进行初始化", c)
		return
	}
	response.OkWithMessage("自动创建数据库成功", c)
}

// CheckDB
// @Tags     CheckDB
// @Summary  初始化用户数据库
// @Produce  application/json
// @Success  200  {object}  response.Response{data=map[string]interface{},msg=string}  "初始化用户数据库"
// @Router   /init/checkdb [post]
func (i *DBApi) CheckDB(c *gin.Context) {
	if !i.checkInitAPIAccess(c) {
		return
	}

	var (
		message  = "前往初始化数据库"
		needInit = true
	)

	if global.GVA_DB != nil {
		message = "数据库无需初始化"
		needInit = false
	}
	global.GVA_LOG.Info(message)
	response.OkWithDetailed(gin.H{"needInit": needInit}, message, c)
}
