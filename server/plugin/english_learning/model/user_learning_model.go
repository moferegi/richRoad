package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// UserLearningAsset 用户外挂资产与学习权限表
type UserLearningAsset struct {
	global.GVA_MODEL
	UserID         uint       `json:"userId" gorm:"uniqueIndex;comment:用户ID"`
	FreeMinutes    int        `json:"freeMinutes" gorm:"default:0;comment:剩余可用免费观看时长(分钟)"`
	TotalPoints    int        `json:"totalPoints" gorm:"default:0;comment:总累积积分(仅做冗余查询，具体以流水为准)"`
	FreeTimeExpire *time.Time `json:"freeTimeExpire" gorm:"comment:新注册赠送全部免费权益过期时间(为空代表无权益/已过期)"`
	IsVip          *bool      `json:"isVip" gorm:"default:false;comment:是否开通全站终身或包月VIP"`
}

// CheckinRecord 签到流水表
type CheckinRecord struct {
	global.GVA_MODEL
	UserID      uint      `json:"userId" gorm:"index;comment:用户ID"`
	CheckinDate time.Time `json:"checkinDate" gorm:"type:date;index;comment:打卡归属日期"`
	PointAward  int       `json:"pointAward" gorm:"comment:本次打卡获得积分数"`
}

// CheckinStats 签到与打卡数据统计(每个用户仅有一条汇总)
type CheckinStats struct {
	global.GVA_MODEL
	UserID         uint `json:"userId" gorm:"uniqueIndex;comment:用户ID"`
	ContinuousDays int  `json:"continuousDays" gorm:"default:0;comment:连续打卡天数"`
	TotalDays      int  `json:"totalDays" gorm:"default:0;comment:累计总打卡天数"`
	WordsToday     int  `json:"wordsToday" gorm:"default:0;comment:今日已背词数量(0点重置)"`
}

// EnglishPointRecord 英语学习积分流水（用于学习模块积分明细展示）
type EnglishPointRecord struct {
	global.GVA_MODEL
	UserID        uint   `json:"userId" gorm:"index;comment:用户ID"`
	ChangeType    string `json:"changeType" gorm:"size:20;comment:增减类型(increase/decrease)"`
	PointChange   int    `json:"pointChange" gorm:"comment:积分变化数量，正数为增加，负数为减少"`
	OperationType string `json:"operationType" gorm:"size:50;comment:操作类型"`
	Reason        string `json:"reason" gorm:"size:200;comment:积分变化原因i18n key"`
	CurrentPoints int    `json:"currentPoints" gorm:"comment:操作后的当前积分余额"`
	Remark        string `json:"remark" gorm:"size:500;comment:备注"`
}

// EnglishFreeTimeRecord 英语学习免费时长流水（用于学习模块时长明细展示）
type EnglishFreeTimeRecord struct {
	global.GVA_MODEL
	UserID             uint   `json:"userId" gorm:"index;comment:用户ID"`
	ChangeType         string `json:"changeType" gorm:"size:20;comment:增减类型(increase/decrease)"`
	MinuteChange       int    `json:"minuteChange" gorm:"comment:时长变化分钟数，正数为增加，负数为减少"`
	OperationType      string `json:"operationType" gorm:"size:50;comment:操作类型"`
	Reason             string `json:"reason" gorm:"size:200;comment:时长变化原因i18n key"`
	CurrentFreeMinutes int    `json:"currentFreeMinutes" gorm:"comment:操作后的当前免费分钟数"`
	RelatedEpisodeID   uint   `json:"relatedEpisodeId" gorm:"default:0;comment:关联视频单集ID"`
	Remark             string `json:"remark" gorm:"size:500;comment:备注"`
}

// UserWatchHistory 视频观看历史 (用于进度恢复)
type UserWatchHistory struct {
	global.GVA_MODEL
	UserID       uint    `json:"userId" gorm:"index;comment:用户ID"`
	EpisodeID    uint    `json:"episodeId" gorm:"index;comment:视频单集ID"`
	ProgressSecs float64 `json:"progressSecs" gorm:"comment:当前观看秒数"`
}

// UserWordHistory 单词背诵进度 (恢复上次看到哪里)
type UserWordHistory struct {
	global.GVA_MODEL
	UserID     uint `json:"userId" gorm:"uniqueIndex;comment:用户ID"`
	CategoryID uint `json:"categoryId" gorm:"comment:最后访问的分类ID"`
	ChapterID  uint `json:"chapterId" gorm:"comment:最后访问的章节ID"`
	WordIndex  int  `json:"wordIndex" gorm:"comment:最后访问的该章节内的词汇序列"`
}

// UserCollection 统一收藏表 (关联单词/句子/视频/日记/日记句子)
type UserCollection struct {
	global.GVA_MODEL
	UserID     uint `json:"userId" gorm:"index;comment:用户ID"`
	TargetType int  `json:"targetType" gorm:"comment:收藏类型 1单词 2视频句子 3视频 4日记 5日记句子"`
	TargetID   uint `json:"targetId" gorm:"comment:关联对应的模型主键ID"`
}
