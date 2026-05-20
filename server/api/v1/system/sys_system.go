package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemApi struct{}

func isSystemConfigRole(authorityID uint) bool {
	return authorityID == 888 || authorityID == 8881 || authorityID == 9528
}

func isSystemSensitiveRole(authorityID uint) bool {
	return authorityID == 888
}

// GetSystemConfig
// @Tags      System
// @Summary   获取配置文件内容
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysConfigResponse,msg=string}  "获取配置文件内容,返回包括系统配置"
// @Router    /system/getSystemConfig [post]
func (s *SystemApi) GetSystemConfig(c *gin.Context) {
	if !isSystemConfigRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	config, err := systemConfigService.GetSystemConfig()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(systemRes.SysConfigResponse{Config: config}, "获取成功", c)
}

// SetSystemConfig
// @Tags      System
// @Summary   设置配置文件内容
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      system.System                   true  "设置配置文件内容"
// @Success   200   {object}  response.Response{data=string}  "设置配置文件内容"
// @Router    /system/setSystemConfig [post]
func (s *SystemApi) SetSystemConfig(c *gin.Context) {
	if !isSystemConfigRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	var sys system.System
	err := c.ShouldBindJSON(&sys)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = systemConfigService.SetSystemConfig(sys)
	if err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败", c)
		return
	}
	response.OkWithMessage("设置成功", c)
}

// ReloadSystem
// @Tags      System
// @Summary   重载系统
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}  "重载系统"
// @Router    /system/reloadSystem [post]
func (s *SystemApi) ReloadSystem(c *gin.Context) {
	if !isSystemSensitiveRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	// 触发系统重载事件
	err := utils.GlobalSystemEvents.TriggerReload()
	if err != nil {
		global.GVA_LOG.Error("重载系统失败!", zap.Error(err))
		response.FailWithMessage("重载系统失败", c)
		return
	}
	response.OkWithMessage("重载系统成功", c)
}

// GetServerInfo
// @Tags      System
// @Summary   获取服务器信息
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "获取服务器信息"
// @Router    /system/getServerInfo [post]
func (s *SystemApi) GetServerInfo(c *gin.Context) {
	if !isSystemSensitiveRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	server, err := systemConfigService.GetServerInfo()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"server": server}, "获取成功", c)
}
