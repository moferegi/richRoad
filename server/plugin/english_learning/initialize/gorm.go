package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		&model.EnglishCategory{},
		&model.EnglishChapter{},
		&model.EnglishWord{},
		&model.EnglishCategoryWord{},
		&model.EnglishWordSentence{},
		&model.EnglishWordErrorLog{},
		&model.VideoCategory{},
		&model.VideoSeries{},
		&model.VideoEpisode{},
		&model.VideoSubtitle{},
		&model.VideoSentence{},
		&model.UserLearningAsset{},
		&model.UserLearningEntitlement{},
		&model.CheckinRecord{},
		&model.CheckinStats{},
		&model.EnglishPointRecord{},
		&model.EnglishFreeTimeRecord{},
		&model.UserWatchHistory{},
		&model.UserWordHistory{},
		&model.UserCollection{},
	)
	if err != nil {
		fmt.Println("注册 English Learning 插件模型失败!", err)
	} else {
		ensureMySQLLongTextColumns(ctx)
		fmt.Println("注册 English Learning 插件模型成功!")
	}
}

func ensureMySQLLongTextColumns(ctx context.Context) {
	if global.GVA_CONFIG.System.DbType != "mysql" {
		return
	}

	alterSQL := []string{
		"ALTER TABLE video_episodes MODIFY COLUMN video_url LONGTEXT",
		"ALTER TABLE video_subtitles MODIFY COLUMN subtitle_url LONGTEXT",
		"ALTER TABLE video_series MODIFY COLUMN cover_url LONGTEXT",
	}

	for _, sql := range alterSQL {
		if err := global.GVA_DB.WithContext(ctx).Exec(sql).Error; err != nil {
			global.GVA_LOG.Warn("English Learning URL字段升级失败", zap.String("sql", sql), zap.Error(err))
		}
	}
}
