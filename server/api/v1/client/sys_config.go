package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysConfigApi struct{}

// GetSysConfigList 分页获取系统参数列表
// @Tags SysConfig
// @Summary 分页获取系统参数列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.SysConfigSearch true "分页获取系统参数列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysConfig/getSysConfigList [get]
func (s *SysConfigApi) GetSysConfigList(c *gin.Context) {
	var pageInfo request.SysConfigSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysConfigService.GetSysConfigList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// UpdateSysConfig 更新系统参数
// @Tags SysConfig
// @Summary 更新系统参数
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body updateSysConfigReq true "更新系统参数"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sysConfig/updateSysConfig [put]
func (s *SysConfigApi) UpdateSysConfig(c *gin.Context) {
	var req struct {
		ID          uint   `json:"id" binding:"required"`
		ConfigValue string `json:"configValue"`
		Remark      string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysConfigService.UpdateSysConfig(req.ID, req.ConfigValue, req.Remark); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetSysConfigByGroup 按分组获取系统参数
// @Tags SysConfig
// @Summary 按分组获取系统参数
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param configGroup query string true "配置分组"
// @Success 200 {object} response.Response{data=[]client.SysConfig,msg=string} "获取成功"
// @Router /sysConfig/getSysConfigByGroup [get]
func (s *SysConfigApi) GetSysConfigByGroup(c *gin.Context) {
	group := c.Query("configGroup")
	if group == "" {
		response.FailWithMessage("configGroup不能为空", c)
		return
	}
	list, err := sysConfigService.GetConfigByGroup(group)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetSysConfigByKey 按key获取单条系统参数
// @Tags SysConfig
// @Summary 按key获取单条系统参数
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param configKey query string true "配置键名"
// @Success 200 {object} response.Response{data=string,msg=string} "获取成功"
// @Router /sysConfig/getSysConfigByKey [get]
func (s *SysConfigApi) GetSysConfigByKey(c *gin.Context) {
	key := c.Query("configKey")
	if key == "" {
		response.FailWithMessage("configKey不能为空", c)
		return
	}
	val, err := sysConfigService.GetConfigByKey(key)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(val, "获取成功", c)
}

// GetLoginConfig 获取登录相关配置（公开接口）
// @Tags SysConfig
// @Summary 获取登录相关配置
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]string,msg=string} "获取成功"
// @Router /sysConfig/getLoginConfig [get]
func (s *SysConfigApi) GetLoginConfig(c *gin.Context) {
	phoneEnabled, _ := sysConfigService.GetConfigByKey("phone_login_enabled")
	defaultMethod, _ := sysConfigService.GetConfigByKey("default_login_method")
	response.OkWithDetailed(map[string]string{
		"phone_login_enabled":  phoneEnabled,
		"default_login_method": defaultMethod,
	}, "获取成功", c)
}

// GetAnnouncementConfig 获取公告配置
// @Tags SysConfig
// @Summary 获取公告配置
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]string,msg=string} "获取成功"
// @Router /sysConfig/getAnnouncementConfig [get]
func (s *SysConfigApi) GetAnnouncementConfig(c *gin.Context) {
	result, err := sysConfigService.GetAnnouncementConfig()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}
