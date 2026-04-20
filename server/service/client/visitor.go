package client

import (
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type VisitorService struct {
}

// 今日统计缓存（避免高频查库）
var (
	statsCache     map[string]interface{}
	statsCacheTime time.Time
	statsCacheMu   sync.RWMutex
	statsCacheTTL  = 30 * time.Second
)

// Heartbeat 记录访客心跳
func (visitorService *VisitorService) Heartbeat(log *client.VisitorLog) error {
	// 冗余日期字段，避免查询时使用 DATE() 函数破坏索引
	log.CreatedDate = time.Now().Format("2006-01-02")
	// 解析 IP 归属地
	log.Location = utils.GetIPLocation(log.IP)
	return global.GVA_DB.Create(log).Error
}

// GetVisitorLogList 分页获取访客日志列表
func (visitorService *VisitorService) GetVisitorLogList(info clientReq.VisitorLogSearch) (list []client.VisitorLog, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.VisitorLog{})

	if info.StartDate != "" {
		db = db.Where("created_at >= ?", info.StartDate)
	}
	if info.EndDate != "" {
		db = db.Where("created_at <= ?", info.EndDate+" 23:59:59")
	}
	if info.VisitorID != "" {
		db = db.Where("visitor_id = ?", info.VisitorID)
	}
	if info.IP != "" {
		db = db.Where("ip = ?", info.IP)
	}
	if info.Platform != "" {
		db = db.Where("platform = ?", info.Platform)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	// 排序支持
	orderClause := "id desc"
	if info.OrderBy != "" {
		allowedCols := map[string]bool{"created_at": true}
		if allowedCols[info.OrderBy] {
			dir := "asc"
			if info.OrderDir == "desc" {
				dir = "desc"
			}
			orderClause = info.OrderBy + " " + dir
		}
	}
	err = db.Order(orderClause).Limit(limit).Offset(offset).Find(&list).Error
	return
}

// GetVisitorSummaryList 分页获取访客汇总列表
func (visitorService *VisitorService) GetVisitorSummaryList(info clientReq.VisitorSummarySearch) (list []client.VisitorSummary, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.VisitorSummary{})

	if info.StartDate != "" {
		db = db.Where("date >= ?", info.StartDate)
	}
	if info.EndDate != "" {
		db = db.Where("date <= ?", info.EndDate)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("date desc").Limit(limit).Offset(offset).Find(&list).Error
	return
}

// GetTodayStats 获取今日实时统计（带30秒缓存）
func (visitorService *VisitorService) GetTodayStats() (map[string]interface{}, error) {
	// 读缓存
	statsCacheMu.RLock()
	if statsCache != nil && time.Since(statsCacheTime) < statsCacheTTL {
		result := statsCache
		statsCacheMu.RUnlock()
		return result, nil
	}
	statsCacheMu.RUnlock()

	today := time.Now().Format("2006-01-02")
	db := global.GVA_DB

	// 使用冗余日期字段查询，命中索引
	var pv int64
	if err := db.Model(&client.VisitorLog{}).Where("created_date = ?", today).Count(&pv).Error; err != nil {
		return nil, err
	}

	var uv int64
	if err := db.Model(&client.VisitorLog{}).Where("created_date = ?", today).Distinct("visitor_id").Count(&uv).Error; err != nil {
		return nil, err
	}

	var registeredUser int64
	if err := db.Model(&client.VisitorLog{}).Where("created_date = ?", today).Where("user_id > 0").Distinct("user_id").Count(&registeredUser).Error; err != nil {
		return nil, err
	}

	var newVisitor int64
	db.Model(&client.VisitorLog{}).
		Where("created_date = ? AND visitor_id NOT IN (?)",
			today,
			db.Model(&client.VisitorLog{}).Select("visitor_id").Where("created_date < ?", today),
		).Distinct("visitor_id").Count(&newVisitor)

	result := map[string]interface{}{
		"date":           today,
		"pv":             pv,
		"uv":             uv,
		"newVisitor":     newVisitor,
		"registeredUser": registeredUser,
	}

	// 写缓存
	statsCacheMu.Lock()
	statsCache = result
	statsCacheTime = time.Now()
	statsCacheMu.Unlock()

	return result, nil
}

// AggregateDailySummary 聚合指定日期的汇总数据
func (visitorService *VisitorService) AggregateDailySummary(date string) error {
	db := global.GVA_DB

	// 兼容旧数据：优先用 created_date 字段，若无数据则回退到 DATE(created_at)
	dateCondition := "created_date = ?"
	beforeCondition := "created_date < ?"
	var checkCount int64
	db.Model(&client.VisitorLog{}).Where("created_date = ?", date).Count(&checkCount)
	if checkCount == 0 {
		// 旧数据可能没有 created_date，回退到 DATE(created_at)
		dateCondition = "DATE(created_at) = ?"
		beforeCondition = "DATE(created_at) < ?"
	}

	var pv int64
	db.Model(&client.VisitorLog{}).Where(dateCondition, date).Count(&pv)

	var uv int64
	db.Model(&client.VisitorLog{}).Where(dateCondition, date).Distinct("visitor_id").Count(&uv)

	var registeredUser int64
	db.Model(&client.VisitorLog{}).Where(dateCondition, date).Where("user_id > 0").Distinct("user_id").Count(&registeredUser)

	// 统计新访客：当天第一次出现的visitor_id
	var newVisitor int64
	db.Model(&client.VisitorLog{}).
		Where(dateCondition+" AND visitor_id NOT IN (?)",
			date,
			db.Model(&client.VisitorLog{}).Select("visitor_id").Where(beforeCondition, date),
		).Distinct("visitor_id").Count(&newVisitor)

	summary := client.VisitorSummary{
		Date:           date,
		PV:             pv,
		UV:             uv,
		NewVisitor:     newVisitor,
		RegisteredUser: registeredUser,
	}

	return db.Where("date = ?", date).Assign(summary).FirstOrCreate(&summary).Error
}

// CleanOldLogs 清理超过指定天数的访客日志
func (visitorService *VisitorService) CleanOldLogs(retainDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retainDays).Format("2006-01-02")
	return global.GVA_DB.Where("created_date < ?", cutoff).Delete(&client.VisitorLog{}).Error
}
