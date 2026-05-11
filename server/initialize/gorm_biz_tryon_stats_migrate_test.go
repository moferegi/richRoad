package initialize

import (
	"testing"

	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTryonStatsMigrateTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}

	if err := db.AutoMigrate(&clientModel.SysConfig{}, &clientModel.ModelCallLog{}, &clientModel.TryonStatsEvent{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	return db
}

func TestMigrateTryonStatsEvent_BackfillsMissingUsageAndModelIdentity(t *testing.T) {
	db := setupTryonStatsMigrateTestDB(t)

	logAsync := clientModel.ModelCallLog{
		TaskID:       1,
		TaskNo:       "task-1",
		CallStage:    "refiner_aliyun_async",
		ModelUsage:   "refiner",
		ModelKey:     "aliyun_aitryon_refiner",
		ModelName:    "aitryon-refiner",
		Provider:     "aliyun",
		Status:       "processing",
		CostPoints:   2,
		RefundPoints: 0,
	}
	if err := db.Create(&logAsync).Error; err != nil {
		t.Fatalf("create async log failed: %v", err)
	}

	logQuery := clientModel.ModelCallLog{
		TaskID:       1,
		TaskNo:       "task-1",
		CallStage:    "refiner_aliyun_query",
		ModelUsage:   "",
		ModelKey:     "",
		ModelName:    "",
		Provider:     "",
		Status:       "success",
		CostPoints:   2,
		RefundPoints: 0,
	}
	if err := db.Create(&logQuery).Error; err != nil {
		t.Fatalf("create query log failed: %v", err)
	}

	if err := migrateTryonStatsEvent(db); err != nil {
		t.Fatalf("migrate tryon stats events failed: %v", err)
	}

	var queryEvent clientModel.TryonStatsEvent
	if err := db.Where("source_log_id = ?", logQuery.ID).First(&queryEvent).Error; err != nil {
		t.Fatalf("query migrated event failed: %v", err)
	}
	if queryEvent.ModelUsage != "refiner" {
		t.Fatalf("expected migrated usage refiner, got %s", queryEvent.ModelUsage)
	}
	if queryEvent.ModelKey != "aliyun_aitryon_refiner" {
		t.Fatalf("expected migrated model key aliyun_aitryon_refiner, got %s", queryEvent.ModelKey)
	}
	if queryEvent.Provider != "aliyun" {
		t.Fatalf("expected migrated provider aliyun, got %s", queryEvent.Provider)
	}
	if queryEvent.Status != "success" {
		t.Fatalf("expected migrated status success, got %s", queryEvent.Status)
	}

	var marker clientModel.SysConfig
	if err := db.Where("config_key = ?", migrateTryonStatsEventKey).First(&marker).Error; err != nil {
		t.Fatalf("query migrate marker failed: %v", err)
	}
}
