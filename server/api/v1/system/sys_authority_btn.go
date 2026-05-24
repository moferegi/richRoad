package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorityBtnApi struct{}

func isAuthorityBtnManageRole(authorityID uint) bool {
	return authorityID == 888 || authorityID == 8881
}

// GetAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   获取权限按钮
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SysAuthorityBtnReq                                      true  "菜单id, 角色id, 选中的按钮id"
// @Success   200   {object}  response.Response{data=response.SysAuthorityBtnRes,msg=string}  "返回列表成功"
// @Router    /authorityBtn/getAuthorityBtn [post]
func (a *AuthorityBtnApi) GetAuthorityBtn(c *gin.Context) {
	adminAuthorityID := utils.GetUserAuthorityId(c)
	if !isAuthorityBtnManageRole(adminAuthorityID) {
		response.FailWithMessage("无权限操作", c)
		return
	}
	var req request.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := authorityService.CheckAuthorityIDAuth(adminAuthorityID, req.AuthorityId); err != nil {
		response.FailWithMessage("无权限操作", c)
		return
	}
	res, err := authorityBtnService.GetAuthorityBtn(req)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(res, "查询成功", c)
}

// SetAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   设置权限按钮
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SysAuthorityBtnReq     true  "菜单id, 角色id, 选中的按钮id"
// @Success   200   {object}  response.Response{msg=string}  "返回列表成功"
// @Router    /authorityBtn/setAuthorityBtn [post]
func (a *AuthorityBtnApi) SetAuthorityBtn(c *gin.Context) {
	adminAuthorityID := utils.GetUserAuthorityId(c)
	if !isAuthorityBtnManageRole(adminAuthorityID) {
		response.FailWithMessage("无权限操作", c)
		return
	}
	var req request.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := authorityService.CheckAuthorityIDAuth(adminAuthorityID, req.AuthorityId); err != nil {
		response.FailWithMessage("无权限操作", c)
		return
	}
	err = authorityBtnService.SetAuthorityBtn(req)
	if err != nil {
		global.GVA_LOG.Error("分配失败!", zap.Error(err))
		response.FailWithMessage("分配失败", c)
		return
	}
	response.OkWithMessage("分配成功", c)
}

// CanRemoveAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   设置权限按钮
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}  "删除成功"
// @Router    /authorityBtn/canRemoveAuthorityBtn [post]
func (a *AuthorityBtnApi) CanRemoveAuthorityBtn(c *gin.Context) {
	if !isAuthorityBtnManageRole(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage("无权限操作", c)
		return
	}
	id := c.Query("id")
	err := authorityBtnService.CanRemoveAuthorityBtn(id)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
