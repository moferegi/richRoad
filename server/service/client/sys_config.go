package client

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
)

type SysConfigService struct{}

// GetSysConfigList 分页获取系统参数列表
func (s *SysConfigService) GetSysConfigList(info clientReq.SysConfigSearch) (list []client.SysConfig, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.SysConfig{})

	if info.ConfigGroup != "" {
		db = db.Where("config_group = ?", info.ConfigGroup)
	}
	if info.ConfigKey != "" {
		db = db.Where("config_key LIKE ?", "%"+info.ConfigKey+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("config_group, id").Limit(limit).Offset(offset).Find(&list).Error
	return
}

// UpdateSysConfig 更新参数值（只允许改value和remark）
func (s *SysConfigService) UpdateSysConfig(id uint, configValue string, remark string) error {
	return global.GVA_DB.Model(&client.SysConfig{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"config_value": configValue,
			"remark":       remark,
		}).Error
}

// GetConfigByKey 根据key获取配置值
func (s *SysConfigService) GetConfigByKey(key string) (string, error) {
	var config client.SysConfig
	err := global.GVA_DB.Where("config_key = ?", key).First(&config).Error
	if err != nil {
		return "", err
	}
	return config.ConfigValue, nil
}

// GetConfigIntByKey 获取int类型参数
func (s *SysConfigService) GetConfigIntByKey(key string, defaultVal int) int {
	val, err := s.GetConfigByKey(key)
	if err != nil || val == "" {
		return defaultVal
	}
	var result int
	_, _ = fmt.Sscanf(val, "%d", &result)
	if result <= 0 {
		return defaultVal
	}
	return result
}

// GetConfigByGroup 按分组获取所有配置
func (s *SysConfigService) GetConfigByGroup(group string) (list []client.SysConfig, err error) {
	err = global.GVA_DB.Where("config_group = ?", group).Order("id").Find(&list).Error
	return
}

// GetAnnouncementConfig 获取公告配置（便捷方法）
func (s *SysConfigService) GetAnnouncementConfig() (map[string]string, error) {
	list, err := s.GetConfigByGroup("announcement")
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, cfg := range list {
		result[cfg.ConfigKey] = cfg.ConfigValue
	}
	return result, nil
}
