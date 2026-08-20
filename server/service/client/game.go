package client

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"gorm.io/gorm"
)

type GameService struct{}

// ==================== 游戏大分类 ====================

func (s *GameService) CreateGameCategory(cat *client.GameCategory) error {
	return global.GVA_DB.Create(cat).Error
}

func (s *GameService) UpdateGameCategory(cat *client.GameCategory) error {
	return global.GVA_DB.Model(&client.GameCategory{}).Where("id = ?", cat.ID).Updates(map[string]interface{}{
		"game_key":    cat.GameKey,
		"name":        cat.Name,
		"description": cat.Description,
		"sort":        cat.Sort,
		"status":      cat.Status,
	}).Error
}

func (s *GameService) DeleteGameCategory(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除关联的难度分类
		if err := tx.Where("game_id = ?", id).Delete(&client.GameDifficultyCategory{}).Error; err != nil {
			return err
		}
		return tx.Delete(&client.GameCategory{}, id).Error
	})
}

func (s *GameService) GetGameCategoryList() (list []client.GameCategory, err error) {
	err = global.GVA_DB.Where("status = 1").Order("sort ASC").Find(&list).Error
	return
}

func (s *GameService) GetGameCategoryByID(id uint) (cat client.GameCategory, err error) {
	err = global.GVA_DB.First(&cat, id).Error
	return
}

func (s *GameService) GetGameCategoryListAdmin(page, pageSize int) (list []client.GameCategory, total int64, err error) {
	db := global.GVA_DB.Model(&client.GameCategory{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if pageSize > 0 {
		db = db.Limit(pageSize).Offset(pageSize * (page - 1))
	}
	err = db.Order("sort ASC").Find(&list).Error
	return
}

// ==================== 游戏难度分类 ====================

func (s *GameService) CreateDifficultyCategory(cat *client.GameDifficultyCategory) error {
	return global.GVA_DB.Create(cat).Error
}

func (s *GameService) UpdateDifficultyCategory(cat *client.GameDifficultyCategory) error {
	return global.GVA_DB.Model(&client.GameDifficultyCategory{}).Where("id = ?", cat.ID).Updates(map[string]interface{}{
		"game_id": cat.GameID,
		"name":    cat.Name,
		"sort":    cat.Sort,
		"status":  cat.Status,
	}).Error
}

func (s *GameService) DeleteDifficultyCategory(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除关联的关卡
		if err := tx.Where("category_id = ?", id).Delete(&client.GameLevel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&client.GameDifficultyCategory{}, id).Error
	})
}

func (s *GameService) GetDifficultyCategoryList(gameID uint) (list []client.GameDifficultyCategory, err error) {
	err = global.GVA_DB.Where("game_id = ? AND status = 1", gameID).Order("sort ASC").Find(&list).Error
	return
}

func (s *GameService) GetDifficultyCategoryByID(id uint) (cat client.GameDifficultyCategory, err error) {
	err = global.GVA_DB.First(&cat, id).Error
	return
}

func (s *GameService) GetDifficultyCategoryListAdmin(gameID uint, page, pageSize int) (list []client.GameDifficultyCategory, total int64, err error) {
	db := global.GVA_DB.Model(&client.GameDifficultyCategory{}).Where("game_id = ?", gameID)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if pageSize > 0 {
		db = db.Limit(pageSize).Offset(pageSize * (page - 1))
	}
	err = db.Order("sort ASC").Find(&list).Error
	return
}

// ==================== 关卡 ====================

func (s *GameService) CreateLevel(level *client.GameLevel) error {
	return global.GVA_DB.Create(level).Error
}

func (s *GameService) UpdateLevel(level *client.GameLevel) error {
	return global.GVA_DB.Model(&client.GameLevel{}).Where("id = ?", level.ID).Updates(map[string]interface{}{
		"category_id":   level.CategoryID,
		"level_number":  level.LevelNumber,
		"numbers":       level.Numbers,
		"target_result": level.TargetResult,
		"sort":          level.Sort,
	}).Error
}

func (s *GameService) DeleteLevel(id uint) error {
	return global.GVA_DB.Delete(&client.GameLevel{}, id).Error
}

func (s *GameService) GetLevelList(categoryID uint) (list []client.GameLevel, err error) {
	err = global.GVA_DB.Where("category_id = ?", categoryID).Order("sort ASC, level_number ASC").Find(&list).Error
	return
}

func (s *GameService) GetLevelListAdmin(categoryID uint, page, pageSize int) (list []client.GameLevel, total int64, err error) {
	db := global.GVA_DB.Model(&client.GameLevel{}).Where("category_id = ?", categoryID)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if pageSize > 0 {
		db = db.Limit(pageSize).Offset(pageSize * (page - 1))
	}
	err = db.Order("sort ASC, level_number ASC").Find(&list).Error
	return
}

func (s *GameService) GetLevelByID(id uint) (level client.GameLevel, err error) {
	err = global.GVA_DB.First(&level, id).Error
	return
}

// ==================== 用户闯关进度 ====================

// GetUserProgress 获取用户在某个难度分类下的闯关进度
func (s *GameService) GetUserProgress(userID uint, categoryID uint) (list []client.GameUserProgress, err error) {
	err = global.GVA_DB.Where("user_id = ? AND level_id IN (SELECT id FROM game_levels WHERE category_id = ?)", userID, categoryID).Find(&list).Error
	return
}

// SubmitLevelResult 提交闯关结果
func (s *GameService) SubmitLevelResult(userID uint, levelID uint) error {
	// 检查关卡是否存在
	var level client.GameLevel
	if err := global.GVA_DB.First(&level, levelID).Error; err != nil {
		return err
	}

	// 检查前置关卡是否已通过
	var prevLevels []client.GameLevel
	if err := global.GVA_DB.Where("category_id = ? AND level_number < ?", level.CategoryID, level.LevelNumber).Find(&prevLevels).Error; err != nil {
		return err
	}
	for _, pl := range prevLevels {
		var count int64
		global.GVA_DB.Model(&client.GameUserProgress{}).Where("user_id = ? AND level_id = ? AND status = 1", userID, pl.ID).Count(&count)
		if count == 0 {
			return gorm.ErrRecordNotFound // 前置关卡未通过
		}
	}

	// 检查是否已通过
	var existing client.GameUserProgress
	err := global.GVA_DB.Where("user_id = ? AND level_id = ?", userID, levelID).First(&existing).Error
	if err == nil {
		if existing.Status == 1 {
			return nil // 已通过，幂等
		}
		// 更新为已通过
		now := time.Now()
		return global.GVA_DB.Model(&existing).Updates(map[string]interface{}{
			"status":        1,
			"passed_at":     &now,
			"attempt_count": existing.AttemptCount + 1,
		}).Error
	}

	// 新建记录
	now := time.Now()
	progress := client.GameUserProgress{
		UserID:       userID,
		LevelID:      levelID,
		Status:       1,
		PassedAt:     &now,
		AttemptCount: 1,
	}
	return global.GVA_DB.Create(&progress).Error
}

// ==================== 排行榜 ====================

type LeaderboardEntry struct {
	UserID   uint   `json:"userID"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Total    int64  `json:"total"`
}

// GetAllTimeLeaderboard 累计闯关榜
func (s *GameService) GetAllTimeLeaderboard(limit int) (list []LeaderboardEntry, err error) {
	if limit <= 0 {
		limit = 50
	}
	err = global.GVA_DB.Table("game_user_progress AS gup").
		Select("gup.user_id, cu.username, cu.nickname, COUNT(*) AS total").
		Joins("LEFT JOIN client_user AS cu ON cu.id = gup.user_id").
		Where("gup.status = 1").
		Group("gup.user_id, cu.username, cu.nickname").
		Order("total DESC").
		Limit(limit).
		Find(&list).Error
	return
}

// GetDailyLeaderboard 每日闯关榜
func (s *GameService) GetDailyLeaderboard(limit int) (list []LeaderboardEntry, err error) {
	if limit <= 0 {
		limit = 50
	}
	today := time.Now().Format("2006-01-02")
	err = global.GVA_DB.Table("game_user_progress AS gup").
		Select("gup.user_id, cu.username, cu.nickname, COUNT(*) AS total").
		Joins("LEFT JOIN client_user AS cu ON cu.id = gup.user_id").
		Where("gup.status = 1 AND DATE(gup.passed_at) = ?", today).
		Group("gup.user_id, cu.username, cu.nickname").
		Order("total DESC").
		Limit(limit).
		Find(&list).Error
	return
}

// GetUserPassedTotal 获取用户通过总关数（用于首页展示）
func (s *GameService) GetUserPassedTotal(userID uint) (total int64, err error) {
	err = global.GVA_DB.Model(&client.GameUserProgress{}).Where("user_id = ? AND status = 1", userID).Count(&total).Error
	return
}

// SetUserProgress 管理端手动设置用户进度（通关/未通关）
func (s *GameService) SetUserProgress(userID uint, levelID uint, status int) error {
	if status == 1 {
		// 设置为通关：upsert
		now := time.Now()
		return global.GVA_DB.Where("user_id = ? AND level_id = ?", userID, levelID).
			Assign(map[string]interface{}{
				"status":        1,
				"passed_at":     &now,
				"attempt_count": gorm.Expr("attempt_count + 1"),
			}).
			FirstOrCreate(&client.GameUserProgress{
				UserID:       userID,
				LevelID:      levelID,
				Status:       1,
				PassedAt:     &now,
				AttemptCount: 1,
			}).Error
	}
	// 设置为未通关：删除记录
	return global.GVA_DB.Where("user_id = ? AND level_id = ?", userID, levelID).Delete(&client.GameUserProgress{}).Error
}
