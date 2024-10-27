package client

import (
	"errors"
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
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type ClientUserApi struct {
}

var clientUserService = service.ServiceGroupApp.ClientServiceGroup.ClientUserService

var store = base64Captcha.DefaultMemStore

func (clientUserApi *ClientUserApi) GetOpenID(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("code不能为空", c)
		return
	}
	res, err := external.GetOpenID(code)
	if err != nil {
		global.GVA_LOG.Error("获取openid失败!", zap.Error(err))
		response.FailWithMessage("获取openid失败", c)
		return
	}
	response.OkWithData(res, c)
}

func (clientUserApi *ClientUserApi) GetUserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	id := strconv.Itoa(int(userID))
	if user, err := clientUserService.GetClientUser(id); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
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
			response.FailWithMessage("用户名不存在或者密码错误", c)
			return
		}
		clientUserApi.TokenNext(c, user)
		return
	}
	// 验证码次数+1
	global.BlackCache.Increment(key, 1)
	response.FailWithMessage("验证码错误", c)
}

func (clientUserApi *ClientUserApi) Register(c *gin.Context) {
	var register clientReq.CreateUser
	err := c.ShouldBindJSON(&register)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if register.Password != register.RePassword {
		response.FailWithMessage("两次输入的密码不一致", c)
		return
	}

	var clientUser client.ClientUser
	clientUser.Username = register.Username
	clientUser.Password = register.Password
	clientUser.Avatar = "https://qmplusimg.henrongyi.top/gva_header.jpg"
	if err := clientUserService.CreateClientUser(&clientUser); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
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
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
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
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
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
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
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
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
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
		response.FailWithMessage("查询失败", c)
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
		response.FailWithMessage("无法修改："+key, c)
		return
	}

	if err := clientUserService.SetClientUserInfo(key, value, userID); err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败:"+err.Error(), c)
	} else {
		response.OkWithMessage("设置成功", c)
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
		response.FailWithMessage("获取token失败", c)
		return
	}
	if !global.GVA_CONFIG.System.UseMultipoint {
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); errors.Is(err, redis.Nil) {
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	} else if err != nil {
		global.GVA_LOG.Error("设置登录状态失败!", zap.Error(err))
		response.FailWithMessage("设置登录状态失败", c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
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
