package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardApi struct{}

var dashboardService = service.ServiceGroupApp.ShopServiceGroup.DashboardService

// GetDashboardOverview 获取数据看板概览
// @Tags Dashboard
// @Summary 获取数据看板概览
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param days query int false "统计天数(默认30)"
// @Success 200 {object} response.Response{data=shop.DashboardOverview,msg=string} "获取成功"
// @Router /dashboard/getOverview [get]
func (api *DashboardApi) GetDashboardOverview(c *gin.Context) {
	if data, err := dashboardService.GetDashboardOverview(); err != nil {
		global.GVA_LOG.Error("获取看板数据失败!", zap.Error(err))
		response.FailWithMessage("获取看板数据失败", c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}
