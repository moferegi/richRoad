package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/service"
	"github.com/gin-gonic/gin"
)

type ConfigApi struct{}

// GetCsConfig 获取客服系统配置（公开接口，UniApp 调用）
// @Tags CustomerService
// @Summary 获取客服系统配置
// @Produce application/json
// @Success 200 {object} response.Response{data=model.CsConfig} "获取成功"
// @Router /cs/config/get [get]
func (a *ConfigApi) GetCsConfig(c *gin.Context) {
	cfg, err := service.Service.ConfigService.GetConfig()
	if err != nil {
		response.FailWithMessage("获取配置失败", c)
		return
	}
	response.OkWithData(cfg, c)
}

// UpdateCsConfig 更新客服系统配置（管理员接口）
// @Tags CustomerServiceAdmin
// @Summary 更新客服系统配置（plat开关）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body map[string]bool true "{ \"platEnabled\": true }"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /cs/admin/config/update [put]
func (a *ConfigApi) UpdateCsConfig(c *gin.Context) {
	var body struct {
		PlatEnabled bool `json:"platEnabled"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.Service.ConfigService.UpdatePlatEnabled(body.PlatEnabled); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
