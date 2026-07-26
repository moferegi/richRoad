package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"gorm.io/gorm"
)

type AdminService struct{}

// ---- 管理端列表条目类型（含用户信息） ----

// AdminCheckinRecordItem 管理端签到记录条目
type AdminCheckinRecordItem struct {
	model.CheckinRecord
	UserName string `json:"userName" gorm:"column:user_name"`
	NickName string `json:"nickName" gorm:"column:nick_name"`
}

// AdminPointRecordItem 管理端积分记录条目
type AdminPointRecordItem struct {
	model.EnglishPointRecord
	UserName string `json:"userName" gorm:"column:user_name"`
	NickName string `json:"nickName" gorm:"column:nick_name"`
}

// AdminFreeTimeRecordItem 管理端时长记录条目
type AdminFreeTimeRecordItem struct {
	model.EnglishFreeTimeRecord
	UserName string `json:"userName" gorm:"column:user_name"`
	NickName string `json:"nickName" gorm:"column:nick_name"`
}

// AdminCollectionItem 管理端收藏记录条目
type AdminCollectionItem struct {
	model.UserCollection
	UserName string `json:"userName" gorm:"column:user_name"`
	NickName string `json:"nickName" gorm:"column:nick_name"`
}

// AdminWordErrorLogItem 管理端错词本条目
type AdminWordErrorLogItem struct {
	model.EnglishWordErrorLog
	UserName     string `json:"userName" gorm:"column:user_name"`
	NickName     string `json:"nickName" gorm:"column:nick_name"`
	WordName     string `json:"wordName" gorm:"column:word_name"`
	CategoryName string `json:"categoryName" gorm:"column:category_name"`
	ChapterName  string `json:"chapterName" gorm:"column:chapter_name"`
}

// GetAdminCheckinRecordList 管理端：签到记录列表（支持userId/日期筛选）
func (s *AdminService) GetAdminCheckinRecordList(userId uint, startDate, endDate string, page, pageSize int) (list []AdminCheckinRecordItem, total int64, err error) {
	page, pageSize = normalizePage(page, pageSize)
	db := global.GVA_DB.Table("checkin_records AS cr").
		Joins("LEFT JOIN sys_users AS u ON u.id = cr.user_id").
		Select("cr.*, u.username AS user_name, u.nick_name")
	if userId > 0 {
		db = db.Where("cr.user_id = ?", userId)
	}
	if startDate != "" {
		db = db.Where("cr.checkin_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("cr.checkin_date <= ?", endDate)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("cr.checkin_date DESC, cr.id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Scan(&list).Error
	return
}

// GetAdminPointRecordList 管理端：积分记录列表（支持userId/changeType/日期筛选）
func (s *AdminService) GetAdminPointRecordList(userId uint, changeType string, startDate, endDate string, page, pageSize int) (list []AdminPointRecordItem, total int64, err error) {
	page, pageSize = normalizePage(page, pageSize)
	db := global.GVA_DB.Table("english_point_records AS pr").
		Joins("LEFT JOIN sys_users AS u ON u.id = pr.user_id").
		Select("pr.*, u.username AS user_name, u.nick_name")
	if userId > 0 {
		db = db.Where("pr.user_id = ?", userId)
	}
	if changeType != "" {
		db = db.Where("pr.change_type = ?", changeType)
	}
	if startDate != "" {
		db = db.Where("pr.created_at >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("pr.created_at <= ?", endDate+" 23:59:59")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("pr.created_at DESC, pr.id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Scan(&list).Error
	return
}

// GetAdminFreeTimeRecordList 管理端：时长明细列表（支持userId/changeType/日期筛选）
func (s *AdminService) GetAdminFreeTimeRecordList(userId uint, changeType string, startDate, endDate string, page, pageSize int) (list []AdminFreeTimeRecordItem, total int64, err error) {
	page, pageSize = normalizePage(page, pageSize)
	db := global.GVA_DB.Table("english_free_time_records AS fr").
		Joins("LEFT JOIN sys_users AS u ON u.id = fr.user_id").
		Select("fr.*, u.username AS user_name, u.nick_name")
	if userId > 0 {
		db = db.Where("fr.user_id = ?", userId)
	}
	if changeType != "" {
		db = db.Where("fr.change_type = ?", changeType)
	}
	if startDate != "" {
		db = db.Where("fr.created_at >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("fr.created_at <= ?", endDate+" 23:59:59")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("fr.created_at DESC, fr.id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Scan(&list).Error
	return
}

// GetAdminWatchHistoryList 管理端：观看历史列表（支持userId/日期筛选，关联单集名称和用户名）
func (s *AdminService) GetAdminWatchHistoryList(userId uint, startDate, endDate string, page, pageSize int) (list []AdminWatchHistoryItem, total int64, err error) {
	page, pageSize = normalizePage(page, pageSize)
	whTable := resolveContentTableName("user_watch_history", "user_watch_histories")
	epTable := resolveContentTableName("video_episode", "video_episodes")

	db := global.GVA_DB.Table(whTable+" AS wh").
		Joins("LEFT JOIN "+epTable+" AS ve ON ve.id = wh.episode_id").
		Joins("LEFT JOIN sys_users AS u ON u.id = wh.user_id").
		Select("wh.*, ve.name AS episode_name, u.username AS user_name, u.nick_name")

	if userId > 0 {
		db = db.Where("wh.user_id = ?", userId)
	}
	if startDate != "" {
		db = db.Where("wh.updated_at >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("wh.updated_at <= ?", endDate+" 23:59:59")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("wh.updated_at DESC, wh.id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Scan(&list).Error
	return
}

// AdminWatchHistoryItem 管理端观看历史条目（含单集名称+用户名）
type AdminWatchHistoryItem struct {
	model.UserWatchHistory
	EpisodeName string `json:"episodeName" gorm:"column:episode_name"`
	UserName   string `json:"userName" gorm:"column:user_name"`
	NickName   string `json:"nickName" gorm:"column:nick_name"`
}

// GetAdminCollectionList 管理端：收藏列表（支持userId/targetType筛选）
func (s *AdminService) GetAdminCollectionList(userId uint, targetType int, page, pageSize int) (list []AdminCollectionItem, total int64, err error) {
	page, pageSize = normalizePage(page, pageSize)
	db := global.GVA_DB.Table("user_collections AS uc").
		Joins("LEFT JOIN sys_users AS u ON u.id = uc.user_id").
		Select("uc.*, u.username AS user_name, u.nick_name")
	if userId > 0 {
		db = db.Where("uc.user_id = ?", userId)
	}
	if targetType > 0 {
		db = db.Where("uc.target_type = ?", targetType)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("uc.id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Scan(&list).Error
	return
}

// GetAdminWordErrorLogList 管理端：错词本列表（支持userId/日期筛选，关联单词/分类/章节名称）
func (s *AdminService) GetAdminWordErrorLogList(userId uint, startDate, endDate string, page, pageSize int) (list []AdminWordErrorLogItem, total int64, err error) {
	page, pageSize = normalizePage(page, pageSize)
	db := global.GVA_DB.Table("english_word_error_logs AS el").
		Joins("LEFT JOIN sys_users AS u ON u.id = el.user_id").
		Joins("LEFT JOIN english_words AS ew ON ew.id = el.word_id").
		Joins("LEFT JOIN english_categories AS ec ON ec.id = el.category_id").
		Joins("LEFT JOIN english_chapters AS ech ON ech.id = el.chapter_id").
		Select("el.*, u.username AS user_name, u.nick_name, ew.word AS word_name, ec.name AS category_name, ech.name AS chapter_name")
	if userId > 0 {
		db = db.Where("el.user_id = ?", userId)
	}
	if startDate != "" {
		db = db.Where("el.created_at >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("el.created_at <= ?", endDate+" 23:59:59")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("el.created_at DESC, el.id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Scan(&list).Error
	return
}

// GetAdminUserList 管理端：获取用户列表（用于下拉筛选，返回userId+username/phone等标识）
func (s *AdminService) GetAdminUserList(keyword string) (list []AdminUserItem, err error) {
	db := global.GVA_DB.Table("sys_users").
		Select("id, username, nick_name, phone").
		Where("deleted_at IS NULL")

	if keyword != "" {
		db = db.Where("username LIKE ? OR nick_name LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	err = db.Limit(50).Scan(&list).Error
	return
}

// AdminUserItem 管理端用户搜索项
type AdminUserItem struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
}

// 确保 gorm 不被误删
var _ = gorm.Expr
