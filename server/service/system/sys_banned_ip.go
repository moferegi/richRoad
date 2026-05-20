package system

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	sysReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
)

type BannedIPService struct{}

const (
	bannedIPRedisHash = "sys:banned_ips" // Redis hash: ip -> expiry_unix (0=永久)

	attackCounterKeyPrefix = "ip:attack:" // Redis string key: ip:attack:{attack_type}:{ip}

	attackAutoBanEnabledConfigKey = "security_attack_auto_ban_enabled"
	attackAutoBanEnabledEnvKey    = "CS_ATTACK_AUTO_BAN_ENABLED"

	attackAutoBanWindowConfigKey = "security_attack_auto_ban_window_seconds"
	attackAutoBanWindowEnvKey    = "CS_ATTACK_AUTO_BAN_WINDOW_SECONDS"

	attackAutoBanDurationConfigKey = "security_attack_auto_ban_duration_minutes"
	attackAutoBanDurationEnvKey    = "CS_ATTACK_AUTO_BAN_DURATION_MINUTES"

	attackAutoBanThresholdDefaultConfigKey = "security_attack_auto_ban_threshold_default"
	attackAutoBanThresholdDefaultEnvKey    = "CS_ATTACK_AUTO_BAN_THRESHOLD_DEFAULT"

	attackAutoBanThresholdSysErrorConfigKey = "security_attack_auto_ban_threshold_sys_error_rate_limit"
	attackAutoBanThresholdSysErrorEnvKey    = "CS_ATTACK_AUTO_BAN_THRESHOLD_SYS_ERROR_RATE_LIMIT"

	defaultAttackAutoBanEnabled                 = true
	defaultAttackAutoBanWindowSeconds     int64 = 3600
	defaultAttackAutoBanDurationMins      int64 = 60
	defaultAttackAutoBanThresholdDefault  int64 = 10
	defaultAttackAutoBanThresholdSysError int64 = 30
)

// BanIP 封禁指定IP
func (s *BannedIPService) BanIP(ip, reason, bannedBy string, durationMinutes int, isAuto bool) error {
	var expiredAt *time.Time
	if durationMinutes > 0 {
		t := time.Now().Add(time.Duration(durationMinutes) * time.Minute)
		expiredAt = &t
	}

	var existing sysModel.SysBannedIP
	result := global.GVA_DB.Where("ip = ?", ip).First(&existing)
	if result.Error == nil {
		// 已存在，更新封禁信息
		updates := map[string]interface{}{
			"reason":     reason,
			"banned_by":  bannedBy,
			"is_auto":    isAuto,
			"expired_at": expiredAt,
		}
		if err := global.GVA_DB.Model(&existing).Updates(updates).Error; err != nil {
			return err
		}
	} else {
		// 新建封禁记录
		ban := sysModel.SysBannedIP{
			IP:        ip,
			Reason:    reason,
			BannedBy:  bannedBy,
			IsAuto:    isAuto,
			ExpiredAt: expiredAt,
		}
		if err := global.GVA_DB.Create(&ban).Error; err != nil {
			return err
		}
	}

	s.syncToRedis(ip, expiredAt)
	return nil
}

// UnbanIP 解封指定IP
func (s *BannedIPService) UnbanIP(ip string) error {
	if err := global.GVA_DB.Where("ip = ?", ip).Delete(&sysModel.SysBannedIP{}).Error; err != nil {
		return err
	}
	if global.GVA_REDIS != nil {
		global.GVA_REDIS.HDel(context.Background(), bannedIPRedisHash, ip)
	}
	return nil
}

// IsBanned 检查IP是否被封禁（优先查Redis，降级查DB）
func (s *BannedIPService) IsBanned(ip string) bool {
	if global.GVA_REDIS != nil {
		val, err := global.GVA_REDIS.HGet(context.Background(), bannedIPRedisHash, ip).Result()
		if err == nil {
			if val == "0" {
				return true // 永久封禁
			}
			expiry, parseErr := strconv.ParseInt(val, 10, 64)
			if parseErr == nil {
				if time.Now().Unix() < expiry {
					return true // 临时封禁未过期
				}
				// 过期，清理
				global.GVA_REDIS.HDel(context.Background(), bannedIPRedisHash, ip)
				global.GVA_DB.Where("ip = ?", ip).Delete(&sysModel.SysBannedIP{})
				return false
			}
		}
		return false
	}

	// DB 降级查询
	var banned sysModel.SysBannedIP
	if err := global.GVA_DB.Where("ip = ?", ip).First(&banned).Error; err != nil {
		return false
	}
	if banned.ExpiredAt == nil {
		return true
	}
	if time.Now().After(*banned.ExpiredAt) {
		global.GVA_DB.Where("ip = ?", ip).Delete(&sysModel.SysBannedIP{})
		return false
	}
	return true
}

// GetBannedIPList 分页查询封禁IP列表
func (s *BannedIPService) GetBannedIPList(req sysReq.BannedIPSearch) (list []sysModel.SysBannedIP, total int64, err error) {
	db := global.GVA_DB.Model(&sysModel.SysBannedIP{})
	if req.IP != "" {
		db = db.Where("ip LIKE ?", "%"+req.IP+"%")
	}
	if req.IsAuto != nil {
		db = db.Where("is_auto = ?", *req.IsAuto)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&list).Error
	return
}

// AttackStatItem 攻击统计条目（全来源合并）
type AttackStatItem struct {
	IP            string    `json:"ip"`
	TotalCount    int64     `json:"totalCount"`    // 总攻击次数
	LoginFail     int64     `json:"loginFail"`     // 密码错误（客户端）
	CaptchaFail   int64     `json:"captchaFail"`   // 验证码错误
	RegisterLimit int64     `json:"registerLimit"` // 注册IP超限
	SysErrorRate  int64     `json:"sysErrorRate"`  // 错误上报频率超限
	AdminFail     int64     `json:"adminFail"`     // 管理员登录失败
	LastTime      time.Time `json:"lastTime"`
	IsBanned      bool      `json:"isBanned"`
}

// GetAttackStats 全面统计各来源攻击行为，合并 sys_attack_logs + sys_login_logs
func (s *BannedIPService) GetAttackStats(hours int) ([]AttackStatItem, error) {
	if hours <= 0 {
		hours = 24
	}
	since := time.Now().Add(-time.Duration(hours) * time.Hour)

	statsMap := map[string]*AttackStatItem{}
	getOrCreate := func(ip string) *AttackStatItem {
		if item, ok := statsMap[ip]; ok {
			return item
		}
		statsMap[ip] = &AttackStatItem{IP: ip}
		return statsMap[ip]
	}

	// 1. 查 sys_attack_logs（客户端各类攻击，按 ip+attack_type 聚合）
	type aggRow struct {
		Ip         string    `gorm:"column:ip"`
		AttackType string    `gorm:"column:attack_type"`
		Count      int64     `gorm:"column:count"`
		LastTime   time.Time `gorm:"column:last_time"`
	}
	var attackRows []aggRow
	if err := global.GVA_DB.Table("sys_attack_logs").
		Select("ip, attack_type, COUNT(*) AS count, MAX(created_at) AS last_time").
		Where("created_at >= ?", since).
		Group("ip, attack_type").
		Scan(&attackRows).Error; err != nil {
		return nil, err
	}
	for _, r := range attackRows {
		item := getOrCreate(r.Ip)
		switch r.AttackType {
		case "login_fail":
			item.LoginFail += r.Count
		case "captcha_fail":
			item.CaptchaFail += r.Count
		case "register_limit":
			item.RegisterLimit += r.Count
		case "sys_error_rate_limit":
			item.SysErrorRate += r.Count
		}
		item.TotalCount += r.Count
		if r.LastTime.After(item.LastTime) {
			item.LastTime = r.LastTime
		}
	}

	// 2. 查 sys_login_logs（管理员后台登录失败）
	type logRow struct {
		Ip       string    `gorm:"column:ip"`
		Count    int64     `gorm:"column:count"`
		LastTime time.Time `gorm:"column:last_time"`
	}
	var adminRows []logRow
	if err := global.GVA_DB.Table("sys_login_logs").
		Select("ip, COUNT(*) AS count, MAX(created_at) AS last_time").
		Where("status = ? AND created_at >= ?", false, since).
		Group("ip").
		Scan(&adminRows).Error; err != nil {
		global.GVA_LOG.Warn("查询管理员登录日志失败", zap.Error(err))
	}
	for _, r := range adminRows {
		item := getOrCreate(r.Ip)
		item.AdminFail += r.Count
		item.TotalCount += r.Count
		if r.LastTime.After(item.LastTime) {
			item.LastTime = r.LastTime
		}
	}

	// 3. 构建结果，按总次数倒排，最多返回 50 条
	items := make([]AttackStatItem, 0, len(statsMap))
	for _, item := range statsMap {
		item.IsBanned = s.IsBanned(item.IP)
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].TotalCount > items[j].TotalCount
	})
	if len(items) > 50 {
		items = items[:50]
	}
	return items, nil
}

// RecordAttack 记录攻击行为并触发自动封禁
// attackType: "login_fail" | "captcha_fail" | "register_limit" | "sys_error_rate_limit"
// 日志写入异步执行，不阻塞请求响应
func (s *BannedIPService) RecordAttack(ip, attackType, username, detail string) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		ip = "unknown"
	}
	attackType = normalizeAttackType(attackType)

	go func() {
		if global.GVA_DB != nil {
			global.GVA_DB.Create(&sysModel.SysAttackLog{
				IP:         ip,
				AttackType: attackType,
				Username:   username,
				Detail:     detail,
			})
		}
	}()

	if global.GVA_REDIS == nil {
		return
	}
	if !utils.GetBoolSetting(attackAutoBanEnabledConfigKey, attackAutoBanEnabledEnvKey, defaultAttackAutoBanEnabled) {
		return
	}

	threshold := s.getAutoBanThreshold(attackType)
	if threshold <= 0 {
		return
	}

	windowSeconds := utils.GetInt64Setting(attackAutoBanWindowConfigKey, attackAutoBanWindowEnvKey, defaultAttackAutoBanWindowSeconds)
	if windowSeconds < 1 {
		windowSeconds = defaultAttackAutoBanWindowSeconds
	}

	banDurationMins := utils.GetInt64Setting(attackAutoBanDurationConfigKey, attackAutoBanDurationEnvKey, defaultAttackAutoBanDurationMins)
	if banDurationMins < 1 {
		banDurationMins = defaultAttackAutoBanDurationMins
	}

	key := fmt.Sprintf("%s%s:%s", attackCounterKeyPrefix, attackType, ip)
	count, _ := global.GVA_REDIS.Incr(context.Background(), key).Result()
	global.GVA_REDIS.Expire(context.Background(), key, time.Duration(windowSeconds)*time.Second)

	if count >= threshold {
		reason := fmt.Sprintf("自动封禁：%s 在 %d 秒内触发 %d 次", attackType, windowSeconds, count)
		_ = s.BanIP(ip, reason, "系统自动", int(banDurationMins), true)
		global.GVA_REDIS.Del(context.Background(), key)
	}
}

func normalizeAttackType(attackType string) string {
	attackType = strings.ToLower(strings.TrimSpace(attackType))
	if attackType == "" {
		return "unknown"
	}
	return attackType
}

func (s *BannedIPService) getAutoBanThreshold(attackType string) int64 {
	if attackType == "sys_error_rate_limit" {
		threshold := utils.GetInt64Setting(
			attackAutoBanThresholdSysErrorConfigKey,
			attackAutoBanThresholdSysErrorEnvKey,
			defaultAttackAutoBanThresholdSysError,
		)
		if threshold < 1 {
			return 0
		}
		return threshold
	}

	threshold := utils.GetInt64Setting(
		attackAutoBanThresholdDefaultConfigKey,
		attackAutoBanThresholdDefaultEnvKey,
		defaultAttackAutoBanThresholdDefault,
	)
	if threshold < 1 {
		return 0
	}
	return threshold
}

// RecordIPLoginFail 向后兼容保留，内部转发 RecordAttack
func (s *BannedIPService) RecordIPLoginFail(ip string) {
	s.RecordAttack(ip, "login_fail", "", "登录失败")
}

// LoadBannedIPCache 启动时从DB加载封禁IP到Redis缓存
func (s *BannedIPService) LoadBannedIPCache() {
	if global.GVA_REDIS == nil || global.GVA_DB == nil {
		return
	}
	var bannedIPs []sysModel.SysBannedIP
	if err := global.GVA_DB.Find(&bannedIPs).Error; err != nil {
		global.GVA_LOG.Error("加载封禁IP缓存失败", zap.Error(err))
		return
	}
	for _, b := range bannedIPs {
		// 跳过已过期的临时封禁
		if b.ExpiredAt != nil && time.Now().After(*b.ExpiredAt) {
			continue
		}
		s.syncToRedis(b.IP, b.ExpiredAt)
	}
	global.GVA_LOG.Info(fmt.Sprintf("已加载 %d 个封禁IP到Redis缓存", len(bannedIPs)))
}

// syncToRedis 同步单个IP封禁状态到Redis
func (s *BannedIPService) syncToRedis(ip string, expiredAt *time.Time) {
	if global.GVA_REDIS == nil {
		return
	}
	if expiredAt == nil {
		global.GVA_REDIS.HSet(context.Background(), bannedIPRedisHash, ip, "0")
	} else {
		global.GVA_REDIS.HSet(context.Background(), bannedIPRedisHash, ip, fmt.Sprintf("%d", expiredAt.Unix()))
	}
}
