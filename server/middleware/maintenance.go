package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

// Maintenance 系统维护模式中间件
// 当 client_sys_config 中 maintenance_enabled = "true" 时，
// 拦截所有非管理员的客户端请求，返回维护信息及维护页面配置
func Maintenance() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 放行管理后台相关路径（system/base/user/authority 等），仅拦截客户端业务路径
		path := c.Request.URL.Path
		// 放行：后台管理API、swagger、health、静态资源、维护状态查询接口本身
		if isAdminPath(path) {
			c.Next()
			return
		}

		// 查询维护模式开关
		var config client.SysConfig
		err := global.GVA_DB.Where("config_key = ?", "maintenance_enabled").First(&config).Error
		if err != nil {
			// 查不到配置项，说明未启用维护模式，放行
			c.Next()
			return
		}

		if config.ConfigValue != "true" {
			c.Next()
			return
		}

		// 检查请求头中是否有管理员token（AuthorityId 888 放行）
		// 如果是已鉴权的管理员请求（通过PrivateGroup的jwt中间件设置），放行
		claims, exists := c.Get("claims")
		if exists && claims != nil {
			// 管理员放行
			c.Next()
			return
		}

		// 查询所有维护模式配置，一并返回给前端
		maintenanceData := gin.H{
			"maintenance": true,
		}
		var configs []client.SysConfig
		if err := global.GVA_DB.Where("config_group = ?", "maintenance").Find(&configs).Error; err == nil {
			for _, cfg := range configs {
				maintenanceData[cfg.ConfigKey] = cfg.ConfigValue
			}
		}

		// 从配置中读取维护提示语，支持多语言JSON格式
		maintenanceMsg := "系统维护中，请稍后再试"
		if msgVal, ok := maintenanceData["maintenance_message"]; ok {
			if s, ok2 := msgVal.(string); ok2 && s != "" {
				maintenanceMsg = resolveI18nMsg(s, c.GetHeader("Accept-Language"))
			}
		}

		// 维护模式生效，返回维护信息
		c.JSON(http.StatusServiceUnavailable, response.Response{
			Code: 503,
			Msg:  maintenanceMsg,
			Data: maintenanceData,
		})
		c.Abort()
	}
}

// isAdminPath 判断是否为后台管理路径（不拦截）
func isAdminPath(path string) bool {
	adminPrefixes := []string{
		"/user/",
		"/base/",
		"/authority/",
		"/menu/",
		"/casbin/",
		"/api/",
		"/jwt/",
		"/system/",
		"/autoCode/",
		"/sysDictionary/",
		"/sysOperationRecord/",
		"/swagger/",
		"/health",
		"/sysConfig/",
		"/language/",
		"/maintenance/",
	}
	for _, prefix := range adminPrefixes {
		if strings.HasPrefix(path, prefix) || strings.Contains(path, prefix) {
			return true
		}
	}
	return false
}

// resolveI18nMsg 解析多语言JSON消息，根据Accept-Language返回对应语言
// 支持格式: {"zh":"中文","en":"English","mn":"Монгол"} 或纯文本
func resolveI18nMsg(raw string, acceptLang string) string {
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, "{") {
		return raw // 纯文本直接返回
	}
	var langMap map[string]string
	if err := json.Unmarshal([]byte(raw), &langMap); err != nil {
		return raw
	}
	// 解析Accept-Language获取首选语言
	lang := "zh"
	if acceptLang != "" {
		parts := strings.Split(acceptLang, ",")
		if len(parts) > 0 {
			lang = strings.TrimSpace(strings.Split(parts[0], ";")[0])
		}
	}
	if v, ok := langMap[lang]; ok && v != "" {
		return v
	}
	// fallback: zh > en > 第一个
	if v, ok := langMap["zh"]; ok && v != "" {
		return v
	}
	if v, ok := langMap["en"]; ok && v != "" {
		return v
	}
	for _, v := range langMap {
		if v != "" {
			return v
		}
	}
	return raw
}
