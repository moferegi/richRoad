package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	englishReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setupLearningAuthzServiceTestDB(t *testing.T) {
	t.Helper()

	originalDB := global.GVA_DB
	t.Cleanup(func() {
		global.GVA_DB = originalDB
	})

	dsn := fmt.Sprintf("file:learning_authz_service_test_%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	if err = db.AutoMigrate(
		&clientModel.SysConfig{},
		&model.UserLearningAsset{},
		&model.UserLearningEntitlement{},
		&model.VideoSeries{},
		&model.VideoEpisode{},
		&model.EnglishCategory{},
	); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	global.GVA_DB = db
}

func TestLearningAuthzService_EnsureUserAsset_AssignsDefaultFreeExpire(t *testing.T) {
	setupLearningAuthzServiceTestDB(t)

	svc := LearningAuthzService{}
	userID := uint(9901)
	start := time.Now()

	asset, err := svc.EnsureUserAsset(userID)
	if err != nil {
		t.Fatalf("ensure user asset failed: %v", err)
	}
	if asset.UserID != userID {
		t.Fatalf("expected userId=%d, got=%d", userID, asset.UserID)
	}
	if asset.FreeTimeExpire == nil {
		t.Fatalf("expected freeTimeExpire not nil")
	}

	minExpected := start.Add(23 * time.Hour)
	maxExpected := start.Add(25 * time.Hour)
	if asset.FreeTimeExpire.Before(minExpected) || asset.FreeTimeExpire.After(maxExpected) {
		t.Fatalf("expected freeTimeExpire within [23h,25h], got=%s", asset.FreeTimeExpire.Format(time.RFC3339))
	}

	asset2, err := svc.EnsureUserAsset(userID)
	if err != nil {
		t.Fatalf("ensure user asset second time failed: %v", err)
	}
	if asset2.ID != asset.ID {
		t.Fatalf("expected same asset record, got first=%d second=%d", asset.ID, asset2.ID)
	}

	var count int64
	if err = global.GVA_DB.Model(&model.UserLearningAsset{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		t.Fatalf("count asset failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected only one asset row, got=%d", count)
	}
}

func TestLearningAuthzService_EvaluateVideoEpisodeAccess_BySeriesEntitlement(t *testing.T) {
	setupLearningAuthzServiceTestDB(t)

	svc := LearningAuthzService{}
	userID := uint(9902)

	isVip := false
	asset := model.UserLearningAsset{UserID: userID, FreeMinutes: 0, IsVip: &isVip}
	if err := global.GVA_DB.Create(&asset).Error; err != nil {
		t.Fatalf("seed asset failed: %v", err)
	}

	series := model.VideoSeries{
		Name:    `{"zh":"测试剧集"}`,
		Price:   12.5,
		NeedVip: &isVip,
	}
	if err := global.GVA_DB.Create(&series).Error; err != nil {
		t.Fatalf("seed video series failed: %v", err)
	}

	episode := model.VideoEpisode{SeriesID: series.ID, Name: `{"zh":"第1集"}`, TrialPercent: 8}
	if err := global.GVA_DB.Create(&episode).Error; err != nil {
		t.Fatalf("seed video episode failed: %v", err)
	}

	decisionBefore, err := svc.EvaluateVideoEpisodeAccess(userID, episode)
	if err != nil {
		t.Fatalf("evaluate episode access before entitlement failed: %v", err)
	}
	if decisionBefore.HasFullAuth {
		t.Fatalf("expected trial before entitlement")
	}
	if !decisionBefore.TrialOnly {
		t.Fatalf("expected trialOnly=true before entitlement")
	}

	grantReq := englishReq.GrantEntitlementReq{
		UserID:       userID,
		ResourceType: model.ResourceTypeVideoSeries,
		ResourceID:   series.ID,
		Remark:       "test grant",
	}
	if err = svc.GrantEntitlement(grantReq, 1); err != nil {
		t.Fatalf("grant entitlement failed: %v", err)
	}

	decisionAfter, err := svc.EvaluateVideoEpisodeAccess(userID, episode)
	if err != nil {
		t.Fatalf("evaluate episode access after entitlement failed: %v", err)
	}
	if !decisionAfter.HasFullAuth {
		t.Fatalf("expected full auth after entitlement")
	}
	if !decisionAfter.HasResourceEntitlement {
		t.Fatalf("expected hasResourceEntitlement=true after entitlement")
	}
	if decisionAfter.TrialOnly {
		t.Fatalf("expected trialOnly=false after entitlement")
	}
}
