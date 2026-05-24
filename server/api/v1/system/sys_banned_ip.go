package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	sysReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BannedIPApi struct{}

var bannedIPService = service.ServiceGroupApp.SystemServiceGroup.BannedIPService

func isBannedIPManageRole(authorityID uint) bool {
	return authorityID == 888 || authorityID == 8881
}

// BanIP 手动封禁IP
// @Tags     安全管理
// @Summary  封禁指定IP
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body sysReq.BanIPRequest true "IP封禁请求"
// @Success  200  {object} response.Response{msg=string} "封禁成功"
// @Router   /sysBannedIP/banIP [post]
func (b *BannedIPApi) BanIP(c *gin.Context) {
	if !isBannedIPManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage("无权限操作", c)
		return
	}
	var req sysReq.BanIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userInfo := utils.GetUserInfo(c)
	if userInfo == nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	if err := bannedIPService.BanIP(req.IP, req.Reason, userInfo.Username, req.Duration, false); err != nil {
		global.GVA_LOG.Error("封禁IP失败", zap.Error(err))
		response.FailWithMessage("封禁失败", c)
		return
	}
	response.OkWithMessage("封禁成功", c)
}

// UnbanIP 解封IP
// @Tags     安全管理
// @Summary  解封指定IP
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body sysReq.UnbanIPRequest true "解封IP请求"
// @Success  200  {object} response.Response{msg=string} "解封成功"
// @Router   /sysBannedIP/unbanIP [post]
func (b *BannedIPApi) UnbanIP(c *gin.Context) {
	if !isBannedIPManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage("无权限操作", c)
		return
	}
	var req sysReq.UnbanIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := bannedIPService.UnbanIP(req.IP); err != nil {
		global.GVA_LOG.Error("解封IP失败", zap.Error(err))
		response.FailWithMessage("解封失败", c)
		return
	}
	response.OkWithMessage("解封成功", c)
}

// GetBannedIPList 获取封禁IP列表
// @Tags     安全管理
// @Summary  分页获取封禁IP列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query sysReq.BannedIPSearch true "查询参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /sysBannedIP/getBannedIPList [get]
func (b *BannedIPApi) GetBannedIPList(c *gin.Context) {
	if !isBannedIPManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage("无权限操作", c)
		return
	}
	var req sysReq.BannedIPSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	list, total, err := bannedIPService.GetBannedIPList(req)
	if err != nil {
		global.GVA_LOG.Error("获取封禁IP列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// GetAttackStats 获取IP攻击统计
// @Tags     安全管理
// @Summary  获取IP攻击统计数据
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    hours query int false "统计最近N小时，默认24"
// @Success  200  {object} response.Response{data=[]system.AttackStatItem,msg=string} "获取成功"
// @Router   /sysBannedIP/getAttackStats [get]
func (b *BannedIPApi) GetAttackStats(c *gin.Context) {
	if !isBannedIPManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage("无权限操作", c)
		return
	}
	var req sysReq.AttackStatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	items, err := bannedIPService.GetAttackStats(req.Hours)
	if err != nil {
		global.GVA_LOG.Error("获取攻击统计失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(items, "获取成功", c)
}
