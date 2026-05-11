package client

import (
	"context"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupPointRecordTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}

	if err := db.AutoMigrate(&clientModel.ClientUser{}, &clientModel.PointRecord{}, &clientModel.TryonPointStatsEvent{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	return db
}

func TestCreatePointRecord_RespectsOuterTransactionRollback(t *testing.T) {
	db := setupPointRecordTestDB(t)
	prevDB := global.GVA_DB
	global.GVA_DB = db
	defer func() {
		global.GVA_DB = prevDB
	}()

	user := clientModel.ClientUser{Point: 10, TryonPoint: 5}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction failed: %v", tx.Error)
	}

	svc := PointRecordService{}
	assetType := clientModel.AssetTypeTryonPoint
	userID := int(user.ID)
	changeType := "decrease"
	pointChange := 2
	operationType := "tryon_consume"
	reason := "试衣任务扣费"

	record := &clientModel.PointRecord{
		AssetType:     &assetType,
		UserId:        &userID,
		ChangeType:    &changeType,
		PointChange:   &pointChange,
		OperationType: &operationType,
		Reason:        &reason,
	}

	ctxWithTx := context.WithValue(context.Background(), "tx", tx)
	if err := svc.CreatePointRecord(ctxWithTx, record); err != nil {
		t.Fatalf("create point record failed: %v", err)
	}

	var inTxUser clientModel.ClientUser
	if err := tx.First(&inTxUser, user.ID).Error; err != nil {
		t.Fatalf("query user in tx failed: %v", err)
	}
	if inTxUser.Point != 10 {
		t.Fatalf("expected point unchanged in tx, got %d", inTxUser.Point)
	}
	if inTxUser.TryonPoint != 3 {
		t.Fatalf("expected tryon point reduced to 3 in tx, got %d", inTxUser.TryonPoint)
	}

	if err := tx.Rollback().Error; err != nil {
		t.Fatalf("rollback transaction failed: %v", err)
	}

	var persistedUser clientModel.ClientUser
	if err := db.First(&persistedUser, user.ID).Error; err != nil {
		t.Fatalf("query user after rollback failed: %v", err)
	}
	if persistedUser.Point != 10 {
		t.Fatalf("expected point unchanged after rollback, got %d", persistedUser.Point)
	}
	if persistedUser.TryonPoint != 5 {
		t.Fatalf("expected tryon point rolled back to 5, got %d", persistedUser.TryonPoint)
	}

	var count int64
	if err := db.Model(&clientModel.PointRecord{}).Count(&count).Error; err != nil {
		t.Fatalf("count point records failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no point record persisted after rollback, got %d", count)
	}

	if err := db.Model(&clientModel.TryonPointStatsEvent{}).Count(&count).Error; err != nil {
		t.Fatalf("count tryon point stats events failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no tryon point stats event persisted after rollback, got %d", count)
	}
}
