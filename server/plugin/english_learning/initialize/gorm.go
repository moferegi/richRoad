package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
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
		fmt.Println("注册 English Learning 插件模型成功!")
	}
}
