package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AddressApi struct {
}

var addressService = service.ServiceGroupApp.ClientServiceGroup.AddressService

// CreateAddress 创建用户地址
// @Tags Address
// @Summary 创建用户地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.Address true "创建用户地址"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /address/createAddress [post]
func (addressApi *AddressApi) CreateAddress(c *gin.Context) {
	var address client.Address
	err := c.ShouldBindJSON(&address)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		address.UserID = utils.GetUserID(c)
	}
	if err := addressService.CreateAddress(&address); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteAddress 删除用户地址
// @Tags Address
// @Summary 删除用户地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.Address true "删除用户地址"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /address/deleteAddress [delete]
func (addressApi *AddressApi) DeleteAddress(c *gin.Context) {
	ID := c.Query("ID")
	var UserID uint
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		UserID = utils.GetUserID(c)
	}
	if err := addressService.DeleteAddress(ID, UserID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteAddressByIds 批量删除用户地址
// @Tags Address
// @Summary 批量删除用户地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /address/deleteAddressByIds [delete]
func (addressApi *AddressApi) DeleteAddressByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	var UserID uint
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		UserID = utils.GetUserID(c)
	}
	if err := addressService.DeleteAddressByIds(IDs, UserID); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateAddress 更新用户地址
// @Tags Address
// @Summary 更新用户地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.Address true "更新用户地址"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /address/updateAddress [put]
func (addressApi *AddressApi) UpdateAddress(c *gin.Context) {
	var address client.Address
	err := c.ShouldBindJSON(&address)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		address.UserID = utils.GetUserID(c)
	}

	if err := addressService.UpdateAddress(address); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindAddress 用id查询用户地址
// @Tags Address
// @Summary 用id查询用户地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query client.Address true "用id查询用户地址"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /address/findAddress [get]
func (addressApi *AddressApi) FindAddress(c *gin.Context) {
	ID := c.Query("ID")
	var UserID uint
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		UserID = utils.GetUserID(c)
	}
	if readdress, err := addressService.GetAddress(ID, UserID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"readdress": readdress}, c)
	}
}

// getDefaultAddress 获取用户默认地址
// @Tags Address
// @Summary 获取用户默认地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /address/getDefaultAddress [get]
func (addressApi *AddressApi) GetDefaultAddress(c *gin.Context) {
	var UserID uint
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		UserID = utils.GetUserID(c)
	}

	address, err := addressService.GetDefaultAddress(UserID) // 接收两个返回值
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(address, c) // 使用返回的地址数据
	}
}

// GetAddressList 分页获取用户地址列表
// @Tags Address
// @Summary 分页获取用户地址列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.AddressSearch true "分页获取用户地址列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /address/getAddressList [get]
func (addressApi *AddressApi) GetAddressList(c *gin.Context) {
	var pageInfo clientReq.AddressSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		pageInfo.UserID = utils.GetUserID(c)
	}
	if list, total, err := addressService.GetAddressInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetAddressDataSource 获取Address的数据源
// @Tags Address
// @Summary 获取Address的数据源
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /address/getAddressDataSource [get]
func (addressApi *AddressApi) GetAddressDataSource(c *gin.Context) {
	// 此接口为获取数据源定义的数据
	if dataSource, err := addressService.GetAddressDataSource(); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(dataSource, c)
	}
}

// GetAddressPublic 不需要鉴权的用户地址接口
// @Tags Address
// @Summary 不需要鉴权的用户地址接口
// @accept application/json
// @Produce application/json
// @Param data query clientReq.AddressSearch true "分页获取用户地址列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /address/getAddressPublic [get]
func (addressApi *AddressApi) GetAddressPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的用户地址接口信息",
	}, "获取成功", c)
}
