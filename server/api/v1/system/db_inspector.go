package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DBInspectorApi struct{}

func isDBInspectorAdmin(authorityID uint) bool {
	return authorityID == 888 || authorityID == 8881
}

// GetOverview 获取数据库巡检总览
// @Tags DBInspector
// @Summary 获取数据库巡检总览（数据库/Redis状态 + 表体积/用途）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query systemReq.DBInspectorOverviewSearch false "分页与关键词"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /dbInspector/getOverview [get]
func (dbInspectorApi *DBInspectorApi) GetOverview(c *gin.Context) {
	if !isDBInspectorAdmin(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage("无权限，仅超级管理员可查看", c)
		return
	}

	var req systemReq.DBInspectorOverviewSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := dbInspectorService.GetOverview(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("获取数据库巡检总览失败", zap.Error(err))
		response.FailWithMessage("获取数据库巡检总览失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

// AutoFix 自动修复数据库或Redis连接
// @Tags DBInspector
// @Summary 自动修复数据库或Redis连接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.DBInspectorAutoFixReq true "修复目标"
// @Success 200 {object} response.Response{data=object,msg=string} "执行完成"
// @Router /dbInspector/autoFix [post]
func (dbInspectorApi *DBInspectorApi) AutoFix(c *gin.Context) {
	if !isDBInspectorAdmin(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage("无权限，仅超级管理员可操作", c)
		return
	}

	var req systemReq.DBInspectorAutoFixReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := dbInspectorService.AutoFix(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("数据库巡检自动修复失败", zap.Error(err))
		response.FailWithMessage("自动修复失败: "+err.Error(), c)
		return
	}
	if result.Success {
		response.OkWithDetailed(result, result.Message, c)
		return
	}
	response.FailWithDetailed(result, result.Message, c)
}

// DeleteRecordsByRange 按日期范围真删除记录并联动清理文件
// @Tags DBInspector
// @Summary 按日期范围真删除记录并联动清理文件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.DBInspectorDeleteByRangeReq true "删除参数"
// @Success 200 {object} response.Response{data=object,msg=string} "删除完成"
// @Router /dbInspector/deleteRecordsByRange [post]
func (dbInspectorApi *DBInspectorApi) DeleteRecordsByRange(c *gin.Context) {
	if !isDBInspectorAdmin(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage("无权限，仅超级管理员可操作", c)
		return
	}

	var req systemReq.DBInspectorDeleteByRangeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := dbInspectorService.DeleteRecordsByRange(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("数据库巡检真删除失败", zap.Error(err))
		response.FailWithMessage("真删除失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "真删除完成", c)
}
