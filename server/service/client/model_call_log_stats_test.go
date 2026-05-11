package client

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupModelCallLogStatsTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}

	if err := db.AutoMigrate(&clientModel.ModelCallLog{}, &clientModel.TryonStatsEvent{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	return db
}

func TestSaveModelCallLog_InferUsageAndIdentityFromHistory(t *testing.T) {
	db := setupModelCallLogStatsTestDB(t)

	prevDB := global.GVA_DB
	global.GVA_DB = db
	defer func() {
		global.GVA_DB = prevDB
	}()

	traceCtx := &modelCallTraceContext{
		UserID:            7,
		TaskID:            99,
		TaskNo:            "tryon_task_99",
		SceneType:         "clothes",
		RefinerCostPoints: 3,
	}

	saveModelCallLog(modelCallLogInput{
		TraceContext: traceCtx,
		CallStage:    "refiner_aliyun_async",
		ModelKey:     "aliyun_aitryon_refiner",
		ModelUsage:   "refiner",
		ModelName:    "aitryon-refiner",
		Provider:     "aliyun",
		Status:       tryonTaskStatusProcessing,
	})
	saveModelCallLog(modelCallLogInput{
		TraceContext: traceCtx,
		CallStage:    "refiner_aliyun_query",
		Status:       tryonTaskStatusSuccess,
	})

	var logs []clientModel.ModelCallLog
	if err := db.Order("id ASC").Find(&logs).Error; err != nil {
		t.Fatalf("query logs failed: %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 model logs, got %d", len(logs))
	}

	lastLog := logs[1]
	if lastLog.ModelUsage != "refiner" {
		t.Fatalf("expected inferred usage refiner, got %s", lastLog.ModelUsage)
	}
	if lastLog.ModelKey != "aliyun_aitryon_refiner" {
		t.Fatalf("expected inferred model key aliyun_aitryon_refiner, got %s", lastLog.ModelKey)
	}
	if lastLog.Provider != "aliyun" {
		t.Fatalf("expected inferred provider aliyun, got %s", lastLog.Provider)
	}
	if lastLog.CostPoints != 3 {
		t.Fatalf("expected inferred refiner cost 3, got %d", lastLog.CostPoints)
	}

	var event clientModel.TryonStatsEvent
	if err := db.Where("source_log_id = ?", lastLog.ID).First(&event).Error; err != nil {
		t.Fatalf("query tryon stats event failed: %v", err)
	}
	if event.ModelUsage != "refiner" {
		t.Fatalf("expected event usage refiner, got %s", event.ModelUsage)
	}
	if event.ModelKey != "aliyun_aitryon_refiner" {
		t.Fatalf("expected event model key aliyun_aitryon_refiner, got %s", event.ModelKey)
	}
	if event.Provider != "aliyun" {
		t.Fatalf("expected event provider aliyun, got %s", event.Provider)
	}
	if event.CostPoints != 3 {
		t.Fatalf("expected event cost points 3, got %d", event.CostPoints)
	}
}
