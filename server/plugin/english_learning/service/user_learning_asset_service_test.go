package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	englishReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setupUserLearningAssetServiceTestDB(t *testing.T) {
	t.Helper()

	originalDB := global.GVA_DB
	originalRedis := global.GVA_REDIS
	t.Cleanup(func() {
		global.GVA_DB = originalDB
		global.GVA_REDIS = originalRedis
	})

	dsn := fmt.Sprintf("file:user_learning_asset_service_test_%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	if err = db.AutoMigrate(
		&model.UserLearningAsset{},
		&model.UserWatchHistory{},
		&model.EnglishFreeTimeRecord{},
	); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	global.GVA_DB = db
	global.GVA_REDIS = nil
}

func TestUserLearningAssetService_HandleHeartbeat_CreateFreeTimeRecord(t *testing.T) {
	setupUserLearningAssetServiceTestDB(t)

	svc := UserLearningAssetService{}
	userID := uint(9401)
	episodeID := uint(701)

	asset := model.UserLearningAsset{
		UserID:      userID,
		FreeMinutes: 3,
	}
	if err := global.GVA_DB.Create(&asset).Error; err != nil {
		t.Fatalf("seed asset failed: %v", err)
	}

	if err := svc.HandleHeartbeat(userID, episodeID, 70, 70); err != nil {
		t.Fatalf("handle heartbeat failed: %v", err)
	}

	var updatedAsset model.UserLearningAsset
	if err := global.GVA_DB.Where("user_id = ?", userID).First(&updatedAsset).Error; err != nil {
		t.Fatalf("query asset failed: %v", err)
	}
	if updatedAsset.FreeMinutes != 2 {
		t.Fatalf("expected free_minutes=2, got=%d", updatedAsset.FreeMinutes)
	}

	var history model.UserWatchHistory
	if err := global.GVA_DB.Where("user_id = ? AND episode_id = ?", userID, episodeID).First(&history).Error; err != nil {
		t.Fatalf("query watch history failed: %v", err)
	}
	if history.ProgressSecs != 70 {
		t.Fatalf("expected progress_secs=70, got=%f", history.ProgressSecs)
	}

	var records []model.EnglishFreeTimeRecord
	if err := global.GVA_DB.Where("user_id = ?", userID).Order("id DESC").Find(&records).Error; err != nil {
		t.Fatalf("query free time records failed: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 free time record, got=%d", len(records))
	}

	record := records[0]
	if record.ChangeType != "decrease" {
		t.Fatalf("expected changeType=decrease, got=%s", record.ChangeType)
	}
	if record.MinuteChange != -1 {
		t.Fatalf("expected minuteChange=-1, got=%d", record.MinuteChange)
	}
	if record.OperationType != freeTimeOpHeartbeatConsume {
		t.Fatalf("expected operationType=%s, got=%s", freeTimeOpHeartbeatConsume, record.OperationType)
	}
	if record.Reason != freeTimeReasonWatchConsume {
		t.Fatalf("expected reason=%s, got=%s", freeTimeReasonWatchConsume, record.Reason)
	}
	if record.CurrentFreeMinutes != 2 {
		t.Fatalf("expected currentFreeMinutes=2, got=%d", record.CurrentFreeMinutes)
	}
	if record.RelatedEpisodeID != episodeID {
		t.Fatalf("expected relatedEpisodeID=%d, got=%d", episodeID, record.RelatedEpisodeID)
	}
}

func TestUserLearningAssetService_GetFreeTimeRecordList(t *testing.T) {
	setupUserLearningAssetServiceTestDB(t)

	svc := UserLearningAssetService{}
	userID := uint(9402)
	now := time.Now()

	records := []model.EnglishFreeTimeRecord{
		{
			UserID:             userID,
			ChangeType:         "increase",
			MinuteChange:       2,
			OperationType:      "point_exchange",
			Reason:             "reason_learningExchange",
			CurrentFreeMinutes: 2,
			RelatedEpisodeID:   0,
			Remark:             "points:200",
			GVA_MODEL:          global.GVA_MODEL{CreatedAt: now.Add(-2 * time.Hour)},
		},
		{
			UserID:             userID,
			ChangeType:         "decrease",
			MinuteChange:       -1,
			OperationType:      freeTimeOpHeartbeatConsume,
			Reason:             freeTimeReasonWatchConsume,
			CurrentFreeMinutes: 1,
			RelatedEpisodeID:   801,
			Remark:             "episode:801,usingSecs:65",
			GVA_MODEL:          global.GVA_MODEL{CreatedAt: now.Add(-1 * time.Hour)},
		},
		{
			UserID:             9403,
			ChangeType:         "increase",
			MinuteChange:       3,
			OperationType:      "point_exchange",
			Reason:             "reason_learningExchange",
			CurrentFreeMinutes: 3,
			RelatedEpisodeID:   0,
			Remark:             "points:300",
			GVA_MODEL:          global.GVA_MODEL{CreatedAt: now.Add(-30 * time.Minute)},
		},
	}
	if err := global.GVA_DB.Create(&records).Error; err != nil {
		t.Fatalf("seed records failed: %v", err)
	}

	info := englishReq.FreeTimeRecordSearch{PageInfo: commonReq.PageInfo{Page: 1, PageSize: 2}}
	list, total, page, pageSize, err := svc.GetFreeTimeRecordList(userID, info)
	if err != nil {
		t.Fatalf("get free time record list failed: %v", err)
	}

	if total != 2 {
		t.Fatalf("expected total=2, got=%d", total)
	}
	if page != 1 || pageSize != 2 {
		t.Fatalf("expected page/pageSize=1/2, got=%d/%d", page, pageSize)
	}
	if len(list) != 2 {
		t.Fatalf("expected list size=2, got=%d", len(list))
	}
	if list[0].OperationType != freeTimeOpHeartbeatConsume {
		t.Fatalf("expected first record operationType=%s, got=%s", freeTimeOpHeartbeatConsume, list[0].OperationType)
	}
	if list[1].OperationType != "point_exchange" {
		t.Fatalf("expected second record operationType=point_exchange, got=%s", list[1].OperationType)
	}
}
