package middleware

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/golang-jwt/jwt/v5"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

const (
	wsAllowQueryTokenConfigKey = "security_ws_allow_query_token"
	wsAllowQueryTokenEnvKey    = "CS_WS_ALLOW_QUERY_TOKEN"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := getTokenFromRequest(c)
		if token == "" {
			response.NoAuth(i18n.T(c, "notLogin"), c)
			c.Abort()
			return
		}
		if isBlacklist(token) {
			response.NoAuth(i18n.T(c, "tokenInvalid"), c)
			utils.ClearToken(c)
			c.Abort()
			return
		}
		j := utils.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				response.NoAuth(i18n.T(c, "tokenExpired"), c)
				utils.ClearToken(c)
				c.Abort()
				return
			}
			response.NoAuth(i18n.T(c, "tokenInvalid"), c)
			utils.ClearToken(c)
			c.Abort()
			return
		}

		// 已登录用户被管理员禁用 需要使该用户的jwt失效 此处比较消耗性能 如果需要 请自行打开
		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开

		//if user, err := userService.FindUserByUuid(claims.UUID.String()); err != nil || user.Enable == 2 {
		//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		//	response.FailWithDetailed(gin.H{"reload": true}, i18n.T(c, "userDisabled"), c)
		//	c.Abort()
		//}
		// 客户端用户封禁检查(AuthorityId=8080)
		if claims.AuthorityId == 8080 {
			var u client.ClientUser
			if err := global.GVA_DB.Select("banned").Where("id = ?", claims.BaseClaims.ID).First(&u).Error; err == nil {
				if u.Banned != nil && *u.Banned {
					response.FailWithDetailed(gin.H{"banned": true, "reload": true}, i18n.T(c, "accountBanned"), c)
					c.Abort()
					return
				}
			}
		}
		c.Set("claims", claims)
		if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime {
			dr, _ := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(dr))
			newToken, _ := j.CreateTokenByOldToken(token, *claims)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
			utils.SetToken(c, newToken, int(dr.Seconds()/60))
			if global.GVA_CONFIG.System.UseMultipoint {
				// 记录新的活跃jwt（管理员单设备模式）
				_ = utils.SetRedisJWT(newToken, newClaims.Username)
			}
			if global.GVA_CONFIG.System.MaxLoginDevices > 0 && claims.AuthorityId == 8080 {
				// 客户端用户：续签时替换设备列表中的旧 token
				utils.RenewDeviceToken(newClaims.Username, token, newToken, dr)
			}
		}
		c.Next()

		if newToken, exists := c.Get("new-token"); exists {
			c.Header("new-token", newToken.(string))
		}
		if newExpiresAt, exists := c.Get("new-expires-at"); exists {
			c.Header("new-expires-at", newExpiresAt.(string))
		}
	}
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool

func isBlacklist(jwt string) bool {
	_, ok := global.BlackCache.Get(jwt)
	return ok
}

func getTokenFromRequest(c *gin.Context) string {
	token := strings.TrimSpace(c.GetHeader("x-token"))
	if token != "" {
		return token
	}

	token = tokenFromAuthorization(c.GetHeader("Authorization"))
	if token != "" {
		return token
	}

	token = tokenFromSecWebSocketProtocol(c.GetHeader("Sec-WebSocket-Protocol"))
	if token != "" {
		return token
	}

	queryToken := strings.TrimSpace(c.Query("token"))
	if queryToken != "" && isWSQueryTokenAllowed(c.Request.URL.Path) {
		return queryToken
	}

	// cookie 兜底（保留现有兼容行为）
	return utils.GetToken(c)
}

func tokenFromAuthorization(auth string) string {
	auth = strings.TrimSpace(auth)
	if auth == "" {
		return ""
	}
	const bearer = "Bearer "
	if len(auth) > len(bearer) && strings.EqualFold(auth[:len(bearer)], bearer) {
		return strings.TrimSpace(auth[len(bearer):])
	}
	if strings.Count(auth, ".") >= 2 {
		return auth
	}
	return ""
}

func tokenFromSecWebSocketProtocol(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	parts := strings.Split(v, ",")
	tokens := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			tokens = append(tokens, t)
		}
	}
	if len(tokens) == 0 {
		return ""
	}
	if len(tokens) >= 2 && (strings.EqualFold(tokens[0], "bearer") || strings.EqualFold(tokens[0], "token")) {
		return tokens[1]
	}
	for _, t := range tokens {
		if strings.Count(t, ".") >= 2 {
			return t
		}
	}
	return ""
}

func isWSQueryTokenAllowed(path string) bool {
	path = strings.TrimSpace(path)
	if !(strings.HasSuffix(path, "/cs/ws") || strings.HasSuffix(path, "/cs/wsAgent")) {
		return false
	}
	return utils.GetBoolSetting(wsAllowQueryTokenConfigKey, wsAllowQueryTokenEnvKey, false)
}
