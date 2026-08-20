package client

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// GameCategory 游戏大分类
type GameCategory struct {
	global.GVA_MODEL
	GameKey     string `json:"gameKey" form:"gameKey" gorm:"column:game_key;size:64;uniqueIndex;comment:游戏标识;"`   // 游戏标识: 24point / 36point / password_crack
	Name        string `json:"name" form:"name" gorm:"column:name;type:json;comment:多语言名称;"`                      // 多语言名称
	Description string `json:"description" form:"description" gorm:"column:description;type:json;comment:多语言描述;"` // 多语言描述
	Sort        int    `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序;"`                         // 排序
	Status      int    `json:"status" form:"status" gorm:"column:status;default:1;comment:1=启用 0=禁用;"`            // 状态
}

func (GameCategory) TableName() string {
	return "game_categories"
}

// GameDifficultyCategory 游戏难度分类
type GameDifficultyCategory struct {
	global.GVA_MODEL
	GameID uint   `json:"gameID" form:"gameID" gorm:"column:game_id;index;comment:关联 game_categories.id;"` // 关联游戏ID
	Name   string `json:"name" form:"name" gorm:"column:name;type:json;comment:多语言名称;"`                    // 多语言名称
	Sort   int    `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序（越高越难）;"`                 // 排序
	Status int    `json:"status" form:"status" gorm:"column:status;default:1;comment:1=启用 0=禁用;"`          // 状态
}

func (GameDifficultyCategory) TableName() string {
	return "game_difficulty_categories"
}

// GameLevel 关卡
type GameLevel struct {
	global.GVA_MODEL
	CategoryID   uint   `json:"categoryID" form:"categoryID" gorm:"column:category_id;index;comment:关联 game_difficulty_categories.id;"` // 关联难度分类ID
	LevelNumber  int    `json:"levelNumber" form:"levelNumber" gorm:"column:level_number;comment:该分类下第几关;"`                             // 关卡序号
	Numbers      string `json:"numbers" form:"numbers" gorm:"column:numbers;size:128;comment:4个数字逗号分隔;"`                                // 4个数字
	TargetResult int    `json:"targetResult" form:"targetResult" gorm:"column:target_result;default:24;comment:目标结果;"`                  // 目标结果
	Sort         int    `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序;"`                                              // 排序
}

func (GameLevel) TableName() string {
	return "game_levels"
}

// GameUserProgress 用户闯关记录
type GameUserProgress struct {
	global.GVA_MODEL
	UserID       uint       `json:"userID" form:"userID" gorm:"column:user_id;index;comment:用户ID;"`                       // 用户ID
	LevelID      uint       `json:"levelID" form:"levelID" gorm:"column:level_id;index;comment:关联 game_levels.id;"`       // 关卡ID
	Status       int        `json:"status" form:"status" gorm:"column:status;default:0;comment:1=已通过 0=未通过;"`             // 状态
	PassedAt     *time.Time `json:"passedAt" form:"passedAt" gorm:"column:passed_at;comment:通过时间;"`                       // 通过时间
	AttemptCount int        `json:"attemptCount" form:"attemptCount" gorm:"column:attempt_count;default:0;comment:尝试次数;"` // 尝试次数
}

func (GameUserProgress) TableName() string {
	return "game_user_progress"
}

// SetUserProgressReq 管理端设置用户进度请求
type SetUserProgressReq struct {
	UserID  uint `json:"userID" binding:"required"`
	LevelID uint `json:"levelID" binding:"required"`
	Status  int  `json:"status" binding:"oneof=0 1"`
}
