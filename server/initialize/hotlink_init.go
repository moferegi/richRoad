package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitHotlinkData 注册防盗链相关API和Casbin权限
func InitHotlinkData() {
	if global.GVA_DB == nil {
		return
	}
	db := global.GVA_DB

	initHotlinkApis(db)
	initHotlinkCasbin(db)
}

func initHotlinkApis(db *gorm.DB) {
	apis := []sysModel.SysApi{
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/signURL", Description: "生成防盗链签名URL"},
		{ApiGroup: "文件上传与下载", Method: "GET", Path: "/fileUploadAndDownload/hotlinkConfig", Description: "获取防盗链配置"},
	}
	for _, api := range apis {
		var count int64
		db.Model(&sysModel.SysApi{}).Where("path = ? AND method = ?", api.Path, api.Method).Count(&count)
		if count == 0 {
			if err := db.Create(&api).Error; err != nil {
				global.GVA_LOG.Error("注册防盗链API失败", zap.String("path", api.Path), zap.Error(err))
			}
		}
	}
}

func initHotlinkCasbin(db *gorm.DB) {
	authIDs := []string{"888", "8881", "8080", "9528"}
	rules := []struct {
		Path   string
		Method string
	}{
		{"/fileUploadAndDownload/signURL", "POST"},
		{"/fileUploadAndDownload/hotlinkConfig", "GET"},
		{"/fileUploadAndDownload/listFolders", "GET"},
	}

	for _, authID := range authIDs {
		for _, rule := range rules {
			var count int64
			db.Table("casbin_rule").Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
				"p", authID, rule.Path, rule.Method).Count(&count)
			if count == 0 {
				db.Exec("INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)",
					"p", authID, rule.Path, rule.Method)
			}
		}
	}
}
