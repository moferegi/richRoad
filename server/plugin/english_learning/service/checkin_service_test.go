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

func setupCheckinServiceTestDB(t *testing.T) {
	t.Helper()

	originalDB := global.GVA_DB
	t.Cleanup(func() {
		global.GVA_DB = originalDB
	})

	dsn := fmt.Sprintf("file:checkin_service_test_%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	if err = db.AutoMigrate(
		&model.UserLearningAsset{},
		&model.CheckinRecord{},
		&model.CheckinStats{},
		&model.EnglishPointRecord{},
		&model.EnglishFreeTimeRecord{},
	); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	global.GVA_DB = db
}

func TestCheckinService_PerformCheckin_CreatePointRecord(t *testing.T) {
	setupCheckinServiceTestDB(t)

	svc := CheckinService{}
	userID := uint(9301)

	award, err := svc.PerformCheckin(userID, 100, 5, 10)
	if err != nil {
		t.Fatalf("perform checkin failed: %v", err)
	}
	if award != 100 {
		t.Fatalf("expected award=100, got=%d", award)
	}

	var asset model.UserLearningAsset
	if err := global.GVA_DB.Where("user_id = ?", userID).First(&asset).Error; err != nil {
		t.Fatalf("query asset failed: %v", err)
	}
	if asset.TotalPoints != 100 {
		t.Fatalf("expected total_points=100, got=%d", asset.TotalPoints)
	}

	var records []model.EnglishPointRecord
	if err := global.GVA_DB.Where("user_id = ?", userID).Order("id ASC").Find(&records).Error; err != nil {
		t.Fatalf("query point records failed: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 point record, got=%d", len(records))
	}

	record := records[0]
	if record.ChangeType != "increase" {
		t.Fatalf("expected changeType=increase, got=%s", record.ChangeType)
	}
	if record.PointChange != 100 {
		t.Fatalf("expected pointChange=100, got=%d", record.PointChange)
	}
	if record.OperationType != checkinOpReward {
		t.Fatalf("expected operationType=%s, got=%s", checkinOpReward, record.OperationType)
	}
	if record.Reason != checkinReasonSignIn {
		t.Fatalf("expected reason=%s, got=%s", checkinReasonSignIn, record.Reason)
	}
	if record.CurrentPoints != 100 {
		t.Fatalf("expected currentPoints=100, got=%d", record.CurrentPoints)
	}

	_, err = svc.PerformCheckin(userID, 100, 5, 10)
	if err == nil {
		t.Fatalf("expected duplicate checkin error, got nil")
	}

	var count int64
	if err := global.GVA_DB.Model(&model.EnglishPointRecord{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		t.Fatalf("count point records failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected point record count still 1, got=%d", count)
	}
}

func TestCheckinService_ExchangePoints_CreatePointRecord(t *testing.T) {
	setupCheckinServiceTestDB(t)

	svc := CheckinService{}
	userID := uint(9302)

	asset := model.UserLearningAsset{
		UserID:      userID,
		TotalPoints: 500,
		FreeMinutes: 0,
	}
	if err := global.GVA_DB.Create(&asset).Error; err != nil {
		t.Fatalf("seed asset failed: %v", err)
	}

	mins, err := svc.ExchangePoints(userID, 200, 100)
	if err != nil {
		t.Fatalf("exchange points failed: %v", err)
	}
	if mins != 2 {
		t.Fatalf("expected exchange minutes=2, got=%d", mins)
	}

	var updated model.UserLearningAsset
	if err := global.GVA_DB.Where("user_id = ?", userID).First(&updated).Error; err != nil {
		t.Fatalf("query updated asset failed: %v", err)
	}
	if updated.TotalPoints != 300 {
		t.Fatalf("expected total_points=300, got=%d", updated.TotalPoints)
	}
	if updated.FreeMinutes != 2 {
		t.Fatalf("expected free_minutes=2, got=%d", updated.FreeMinutes)
	}

	var records []model.EnglishPointRecord
	if err := global.GVA_DB.Where("user_id = ?", userID).Order("id DESC").Find(&records).Error; err != nil {
		t.Fatalf("query point records failed: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 point record, got=%d", len(records))
	}

	record := records[0]
	if record.ChangeType != "decrease" {
		t.Fatalf("expected changeType=decrease, got=%s", record.ChangeType)
	}
	if record.PointChange != -200 {
		t.Fatalf("expected pointChange=-200, got=%d", record.PointChange)
	}
	if record.OperationType != checkinOpPointExchange {
		t.Fatalf("expected operationType=%s, got=%s", checkinOpPointExchange, record.OperationType)
	}
	if record.Reason != checkinReasonPointDeduct {
		t.Fatalf("expected reason=%s, got=%s", checkinReasonPointDeduct, record.Reason)
	}
	if record.CurrentPoints != 300 {
		t.Fatalf("expected currentPoints=300, got=%d", record.CurrentPoints)
	}

	var freeTimeRecords []model.EnglishFreeTimeRecord
	if err := global.GVA_DB.Where("user_id = ?", userID).Order("id DESC").Find(&freeTimeRecords).Error; err != nil {
		t.Fatalf("query free time records failed: %v", err)
	}
	if len(freeTimeRecords) != 1 {
		t.Fatalf("expected 1 free time record, got=%d", len(freeTimeRecords))
	}

	freeTimeRecord := freeTimeRecords[0]
	if freeTimeRecord.ChangeType != "increase" {
		t.Fatalf("expected freeTime changeType=increase, got=%s", freeTimeRecord.ChangeType)
	}
	if freeTimeRecord.MinuteChange != 2 {
		t.Fatalf("expected freeTime minuteChange=2, got=%d", freeTimeRecord.MinuteChange)
	}
	if freeTimeRecord.OperationType != checkinOpPointExchange {
		t.Fatalf("expected freeTime operationType=%s, got=%s", checkinOpPointExchange, freeTimeRecord.OperationType)
	}
	if freeTimeRecord.Reason != checkinReasonFreeTimeAdd {
		t.Fatalf("expected freeTime reason=%s, got=%s", checkinReasonFreeTimeAdd, freeTimeRecord.Reason)
	}
	if freeTimeRecord.CurrentFreeMinutes != 2 {
		t.Fatalf("expected currentFreeMinutes=2, got=%d", freeTimeRecord.CurrentFreeMinutes)
	}
}

func TestCheckinService_GetCheckinRecordList(t *testing.T) {
	setupCheckinServiceTestDB(t)

	svc := CheckinService{}
	userID := uint(9303)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	twoDaysAgo := today.AddDate(0, 0, -2)

	records := []model.CheckinRecord{
		{UserID: userID, CheckinDate: twoDaysAgo, PointAward: 100},
		{UserID: userID, CheckinDate: yesterday, PointAward: 105},
		{UserID: 9304, CheckinDate: today, PointAward: 100},
	}
	if err := global.GVA_DB.Create(&records).Error; err != nil {
		t.Fatalf("seed checkin records failed: %v", err)
	}

	info := englishReq.CheckinRecordSearch{PageInfo: commonReq.PageInfo{Page: 1, PageSize: 10}}
	list, total, page, pageSize, err := svc.GetCheckinRecordList(userID, info)
	if err != nil {
		t.Fatalf("get checkin record list failed: %v", err)
	}

	if total != 2 {
		t.Fatalf("expected total=2, got=%d", total)
	}
	if page != 1 || pageSize != 10 {
		t.Fatalf("expected page/pageSize=1/10, got=%d/%d", page, pageSize)
	}
	if len(list) != 2 {
		t.Fatalf("expected list size=2, got=%d", len(list))
	}
	if !list[0].CheckinDate.Equal(yesterday) {
		t.Fatalf("expected first checkinDate=%s, got=%s", yesterday.Format("2006-01-02"), list[0].CheckinDate.Format("2006-01-02"))
	}
	if !list[1].CheckinDate.Equal(twoDaysAgo) {
		t.Fatalf("expected second checkinDate=%s, got=%s", twoDaysAgo.Format("2006-01-02"), list[1].CheckinDate.Format("2006-01-02"))
	}
}
