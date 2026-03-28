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
