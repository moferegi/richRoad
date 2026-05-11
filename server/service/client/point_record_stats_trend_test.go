package client

import (
	"context"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTryonPointStatsTrendTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}

	if err := db.AutoMigrate(&clientModel.TryonPointStatsEvent{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	return db
}

func TestGetTryonPointStats_AutoTrendFallsBackToLatestDataWindow(t *testing.T) {
	db := setupTryonPointStatsTrendTestDB(t)

	prevDB := global.GVA_DB
	global.GVA_DB = db
	defer func() {
		global.GVA_DB = prevDB
	}()

	latestWithData := time.Now().AddDate(0, 0, -60)
	event := clientModel.TryonPointStatsEvent{
		SourcePointRecordID: 1,
		EventAt:             latestWithData,
		UserID:              100,
		ChangeType:          "increase",
		OperationType:       "tryon_recharge",
		PointChange:         12,
		PointAmount:         12,
		Reason:              "充值",
	}
	if err := db.Create(&event).Error; err != nil {
		t.Fatalf("create stats event failed: %v", err)
	}

	svc := PointRecordService{}
	stats, err := svc.GetTryonPointStats(context.Background(), clientReq.TryonPointStatsSearch{})
	if err != nil {
		t.Fatalf("get tryon point stats failed: %v", err)
	}
	if len(stats.DailyTrend) != 14 {
		t.Fatalf("expected auto fallback trend window with 14 days, got %d", len(stats.DailyTrend))
	}

	targetDate := latestWithData.Format("2006-01-02")
	found := false
	for _, item := range stats.DailyTrend {
		if item.Date != targetDate {
			continue
		}
		found = true
		if item.Granted != 12 {
			t.Fatalf("expected granted=12 on latest data day, got %d", item.Granted)
		}
	}
	if !found {
		t.Fatalf("expected trend to contain latest data date %s", targetDate)
	}
}
