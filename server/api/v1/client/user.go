package client

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/external"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ClientUserApi struct {
}

var clientUserService = service.ServiceGroupApp.ClientServiceGroup.ClientUserService

var store = base64Captcha.DefaultMemStore

func (clientUserApi *ClientUserApi) GetOpenID(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage(i18n.T(c, "codeEmpty"), c)
		return
	}
	res, err := external.GetOpenID(code)
	if err != nil {
		global.GVA_LOG.Error("获取openid失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "openidFail"), c)
		return
	}
	response.OkWithData(res, c)
}

func (clientUserApi *ClientUserApi) GetUserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	id := strconv.Itoa(int(userID))
	if user, err := clientUserService.GetClientUser(id); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "queryFail"), c)
	} else {
		response.OkWithData(user, c)
	}
}

func (clientUserApi *ClientUserApi) Login(c *gin.Context) {

	var l clientReq.Login
	err := c.ShouldBindJSON(&l)
	key := c.ClientIP()

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.LoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 判断验证码是否开启
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool = openCaptcha == 0 || openCaptcha < interfaceToInt(v)

	if !oc || (l.CaptchaId != "" && l.Captcha != "" && store.Verify(l.CaptchaId, l.Captcha, true)) {
		user, err := clientUserService.Login(&l)
		if err != nil {
			global.GVA_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
			// 验证码次数+1
			global.BlackCache.Increment(key, 1)
			response.FailWithMessage(i18n.T(c, "loginFail"), c)
			return
		}
		clientUserApi.TokenNext(c, user)
		return
	}
	// 验证码次数+1
	global.BlackCache.Increment(key, 1)
	response.FailWithMessage(i18n.T(c, "captchaError"), c)
}

func (clientUserApi *ClientUserApi) Register(c *gin.Context) {
	var register clientReq.CreateUser
	err := c.ShouldBindJSON(&register)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if register.Password != register.RePassword {
		response.FailWithMessage(i18n.T(c, "passwordMismatch"), c)
		return
	}

	var clientUser client.ClientUser
	clientUser.Username = register.Username
	clientUser.Password = register.Password
	clientUser.Avatar = "https://qmplusimg.henrongyi.top/gva_header.jpg"

	// 处理邀请码
	if register.InviteCode != "" {
		var inviter client.ClientUser
		if err := global.GVA_DB.Where("invite_code = ?", register.InviteCode).First(&inviter).Error; err == nil {
			clientUser.InvitedBy = inviter.ID
		}
	}

	if err := clientUserService.CreateClientUser(&clientUser); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "createFail"), c)
		return
	}

	// 邀请奖励：给邀请人发放 sub_register 营销奖励（积分+优惠券）
	if clientUser.InvitedBy > 0 {
		// 检查 sub_register 配置是否要求下级首单
		reward, rewardErr := marketingRewardService.GetMarketingRewardByType("sub_register")
		if rewardErr == nil && (reward.SubRequireOrder == nil || !*reward.SubRequireOrder) {
			// 不要求首单，直接奖励
			if err := marketingRewardService.TriggerReward(clientUser.InvitedBy, "sub_register", "invite_reward", "邀请用户 "+clientUser.Username+" 注册奖励", 0); err != nil {
				global.GVA_LOG.Error("邀请奖励发放失败", zap.Error(err))
			}
		}
		// 如果 SubRequireOrder=true，等下级首单付款时再发，由 order 服务处理
	}

	// 注册奖励：给新用户发放 register 营销奖励（积分+优惠券）
	if err := marketingRewardService.TriggerReward(clientUser.ID, "register", "register_reward", "新用户注册奖励", 0); err != nil {
		global.GVA_LOG.Error("注册奖励发放失败", zap.Error(err))
	}

	response.OkWithMessage(i18n.T(c, "createSuccess"), c)
}

// CreateClientUser 创建客户端用户
// @Tags ClientUser
// @Summary 创建客户端用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.ClientUser true "创建客户端用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /clientUser/createClientUser [post]
func (clientUserApi *ClientUserApi) CreateClientUser(c *gin.Context) {
	var clientUser client.ClientUser
	err := c.ShouldBindJSON(&clientUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	clientUser.CreatedBy = utils.GetUserID(c)

	if err := clientUserService.CreateClientUser(&clientUser); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "createFail"), c)
	} else {
		response.OkWithMessage(i18n.T(c, "createSuccess"), c)
	}
}

// DeleteClientUser 删除客户端用户
// @Tags ClientUser
// @Summary 删除客户端用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.ClientUser true "删除客户端用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /clientUser/deleteClientUser [delete]
func (clientUserApi *ClientUserApi) DeleteClientUser(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	if err := clientUserService.DeleteClientUser(ID, userID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "deleteFail"), c)
	} else {
		response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
	}
}

// DeleteClientUserByIds 批量删除客户端用户
// @Tags ClientUser
// @Summary 批量删除客户端用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /clientUser/deleteClientUserByIds [delete]
func (clientUserApi *ClientUserApi) DeleteClientUserByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	if err := clientUserService.DeleteClientUserByIds(IDs, userID); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "batchDeleteFail"), c)
	} else {
		response.OkWithMessage(i18n.T(c, "batchDeleteSuccess"), c)
	}
}

// UpdateClientUser 更新客户端用户
// @Tags ClientUser
// @Summary 更新客户端用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.ClientUser true "更新客户端用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /clientUser/updateClientUser [put]
func (clientUserApi *ClientUserApi) UpdateClientUser(c *gin.Context) {
	var clientUser client.ClientUser
	err := c.ShouldBindJSON(&clientUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	clientUser.UpdatedBy = utils.GetUserID(c)

	if err := clientUserService.UpdateClientUser(clientUser); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "updateFail"), c)
	} else {
		response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
	}
}

// FindClientUser 用id查询客户端用户
// @Tags ClientUser
// @Summary 用id查询客户端用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query client.ClientUser true "用id查询客户端用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /clientUser/findClientUser [get]
func (clientUserApi *ClientUserApi) FindClientUser(c *gin.Context) {
	ID := c.Query("ID")
	if reclientUser, err := clientUserService.GetClientUser(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "queryFail"), c)
	} else {
		response.OkWithData(gin.H{"reclientUser": reclientUser}, c)
	}
}

// GetClientUserList 分页获取客户端用户列表
// @Tags ClientUser
// @Summary 分页获取客户端用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.ClientUserSearch true "分页获取客户端用户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /clientUser/getClientUserList [get]
func (clientUserApi *ClientUserApi) GetClientUserList(c *gin.Context) {
	var pageInfo clientReq.ClientUserSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := clientUserService.GetClientUserInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, i18n.T(c, "getSuccess"), c)
	}
}

// GetClientUserPublic 不需要鉴权的客户端用户接口
// @Tags ClientUser
// @Summary 不需要鉴权的客户端用户接口
// @accept application/json
// @Produce application/json
// @Param data query clientReq.ClientUserSearch true "分页获取客户端用户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /clientUser/setClientUserInfo [post]
func (clientUserApi *ClientUserApi) SetClientUserInfo(c *gin.Context) {

	var req clientReq.UpdateKV
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	key := req.Key
	value := req.Value
	userID := utils.GetUserID(c)

	keyWhiteList := []string{"avatar", "nickname", "gender", "phone", "email"}
	inWhiteList := false
	for i := range keyWhiteList {
		if keyWhiteList[i] == key {
			inWhiteList = true
			break
		}
	}
	if !inWhiteList {
		response.FailWithMessage(i18n.T(c, "cannotModify")+":"+key, c)
		return
	}

	if err := clientUserService.SetClientUserInfo(key, value, userID); err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "setFail")+":"+err.Error(), c)
	} else {
		response.OkWithMessage(i18n.T(c, "setSuccess"), c)
	}
}

// TokenNext 登录以后签发jwt
func (clientUserApi *ClientUserApi) TokenNext(c *gin.Context, user client.ClientUser) {
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.Nickname,
		Username:    user.Username,
		AuthorityId: 8080,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "tokenFail"), c)
		return
	}
	if !global.GVA_CONFIG.System.UseMultipoint {
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, i18n.T(c, "loginSuccess"), c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); errors.Is(err, redis.Nil) {
		if err := utils.SetRedisJWT(token, user.Username); err != nil {
			global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
			response.FailWithMessage(i18n.T(c, "loginStatusFail"), c)
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, i18n.T(c, "loginSuccess"), c)
	} else if err != nil {
		global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "loginStatusFail"), c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage(i18n.T(c, "jwtBlacklistFail"), c)
			return
		}
		if err := utils.SetRedisJWT(token, user.Username); err != nil {
			response.FailWithMessage(i18n.T(c, "loginStatusFail"), c)
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, i18n.T(c, "loginSuccess"), c)
	}
}

// 类型转换
func interfaceToInt(v interface{}) (i int) {
	switch v := v.(type) {
	case int:
		i = v
	default:
		i = 0
	}
	return
}

// GetSubordinates 获取下级用户列表
// @Tags ClientUser
// @Summary 获取用户的下级列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param userID query int true "用户ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /clientUser/getSubordinates [get]
func (clientUserApi *ClientUserApi) GetSubordinates(c *gin.Context) {
	userIDStr := c.Query("userID")
	if userIDStr == "" {
		response.FailWithMessage("用户ID不能为空", c)
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.FailWithMessage("用户ID格式错误", c)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	list, total, err := clientUserService.GetSubordinates(uint(userID), page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("获取下级失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}

// GetMyInviteInfo 获取当前用户的邀请信息（uni-app端）
// @Tags ClientUser
// @Summary 获取我的邀请信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /clientUser/getMyInviteInfo [get]
func (clientUserApi *ClientUserApi) GetMyInviteInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	user, err := clientUserService.GetClientUser(strconv.Itoa(int(userID)))
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	subordinateCount, _ := clientUserService.GetSubordinateCount(userID)
	response.OkWithDetailed(gin.H{
		"inviteCode":       user.InviteCode,
		"subordinateCount": subordinateCount,
		"point":            user.Point,
	}, "获取成功", c)
}

// GetMySubordinates 获取当前用户的下级列表（uni-app端）
// @Tags ClientUser
// @Summary 获取我的下级列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /clientUser/getMySubordinates [get]
func (clientUserApi *ClientUserApi) GetMySubordinates(c *gin.Context) {
	userID := utils.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	list, total, err := clientUserService.GetSubordinates(userID, page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("获取下级失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}

// PhoneLogin 手机号+密码登录
// @Tags ClientUser
// @Summary 手机号+密码登录
// @accept application/json
// @Produce application/json
// @Param data body clientReq.PhoneLoginRequest true "手机号登录参数"
// @Success 200 {object} response.Response{data=systemRes.LoginResponse,msg=string} "登录成功"
// @Router /clientUser/phoneLogin [post]
func (clientUserApi *ClientUserApi) PhoneLogin(c *gin.Context) {
	var req clientReq.PhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证码校验
	if !store.Verify(req.CaptchaId, req.Captcha, true) {
		response.FailWithMessage(i18n.T(c, "captchaError"), c)
		return
	}

	// 检查登录失败限制
	var securityService = service.ServiceGroupApp.ClientServiceGroup.SecurityService
	allowed, waitSec, _ := securityService.CheckLoginFail(req.Phone)
	if !allowed {
		response.FailWithMessage(i18n.T(c, "loginLocked")+" "+strconv.Itoa(waitSec)+"s", c)
		return
	}

	user, err := clientUserService.LoginByPhone(req.AreaCode, req.Phone, req.Password)
	if err != nil {
		global.GVA_LOG.Error("手机号登录失败!", zap.Error(err))
		securityService.RecordLoginFail(req.Phone)
		response.FailWithMessage(i18n.T(c, "loginFail"), c)
		return
	}

	securityService.ClearLoginFail(req.Phone)
	clientUserApi.TokenNext(c, user)
}

// PhoneRegister 手机号注册
// @Tags ClientUser
// @Summary 手机号注册
// @accept application/json
// @Produce application/json
// @Param data body clientReq.PhoneRegisterRequest true "手机号注册参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"注册成功"}"
// @Router /clientUser/phoneRegister [post]
func (clientUserApi *ClientUserApi) PhoneRegister(c *gin.Context) {
	var req clientReq.PhoneRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证码校验
	if !store.Verify(req.CaptchaId, req.Captcha, true) {
		response.FailWithMessage(i18n.T(c, "captchaError"), c)
		return
	}

	// IP注册限制
	var securityService = service.ServiceGroupApp.ClientServiceGroup.SecurityService
	allowed, _ := securityService.CheckRegisterIPLimit(c.ClientIP())
	if !allowed {
		response.FailWithMessage(i18n.T(c, "registerIPLimit"), c)
		return
	}

	// 验证手机号格式（通过PhoneAreaCode的正则）
	var areaCode client.PhoneAreaCode
	if err := global.GVA_DB.Where("area_code = ? AND is_enabled = ?", req.AreaCode, true).First(&areaCode).Error; err != nil {
		response.FailWithMessage(i18n.T(c, "areaCodeInvalid"), c)
		return
	}
	if areaCode.PhoneRegex != "" {
		matched, _ := regexp.MatchString(areaCode.PhoneRegex, req.Phone)
		if !matched {
			response.FailWithMessage(i18n.T(c, "phoneFormatError"), c)
			return
		}
	}

	// 注册
	clientUser, err := clientUserService.RegisterByPhone(req.AreaCode, req.Phone, req.Password, req.InviteCode)
	if err != nil {
		global.GVA_LOG.Error("手机号注册失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 记录IP注册次数
	securityService.IncrementRegisterIP(c.ClientIP())

	// 邀请奖励：给邀请人发放 sub_register 营销奖励（积分+优惠券）
	if clientUser.InvitedBy > 0 {
		reward, rewardErr := marketingRewardService.GetMarketingRewardByType("sub_register")
		if rewardErr == nil && (reward.SubRequireOrder == nil || !*reward.SubRequireOrder) {
			if err := marketingRewardService.TriggerReward(clientUser.InvitedBy, "sub_register", "invite_reward", "邀请用户 "+clientUser.Username+" 注册奖励", 0); err != nil {
				global.GVA_LOG.Error("邀请奖励发放失败", zap.Error(err))
			}
		}
	}

	// 注册奖励：给新用户发放 register 营销奖励（积分+优惠券）
	if err := marketingRewardService.TriggerReward(clientUser.ID, "register", "register_reward", "新用户注册奖励", 0); err != nil {
		global.GVA_LOG.Error("注册奖励发放失败", zap.Error(err))
	}

	response.OkWithMessage(i18n.T(c, "createSuccess"), c)
}

// ChangePassword 修改密码
// @Tags ClientUser
// @Summary 修改密码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body clientReq.ChangePasswordRequest true "修改密码参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /clientUser/changePassword [post]
func (clientUserApi *ClientUserApi) ChangePassword(c *gin.Context) {
	var req clientReq.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证码校验
	if !store.Verify(req.CaptchaId, req.Captcha, true) {
		response.FailWithMessage(i18n.T(c, "captchaError"), c)
		return
	}

	userID := utils.GetUserID(c)
	if err := clientUserService.ChangePassword(userID, req.OldPassword, req.NewPassword, req.Method); err != nil {
		global.GVA_LOG.Error("修改密码失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "changeSuccess"), c)
}
