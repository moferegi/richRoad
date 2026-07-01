package service

import (
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	englishModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	englishReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupEnglishWordReportErrorTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:english_word_report_error_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db failed: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(&englishModel.EnglishWord{}, &englishModel.EnglishWordErrorLog{}); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	return db
}

func TestReportWordError_AccumulatesWrongCount(t *testing.T) {
	db := setupEnglishWordReportErrorTestDB(t)
	originalDB := global.GVA_DB
	global.GVA_DB = db
	t.Cleanup(func() {
		global.GVA_DB = originalDB
	})

	word := englishModel.EnglishWord{Word: "report_error_regression_word", Explanation: "{}"}
	if err := db.Create(&word).Error; err != nil {
		t.Fatalf("create word failed: %v", err)
	}

	svc := EnglishWordService{}
	firstReq := englishReq.ReportWordErrorReq{
		WordID:       word.ID,
		CategoryID:   10,
		ChapterID:    100,
		WrongIndex:   1,
		ExpectedChar: "ab",
		InputChar:    "xy",
	}
	if err := svc.ReportWordError(9527, firstReq); err != nil {
		t.Fatalf("first report failed: %v", err)
	}

	secondReq := englishReq.ReportWordErrorReq{
		WordID:       word.ID,
		CategoryID:   11,
		ChapterID:    101,
		WrongIndex:   3,
		ExpectedChar: "c",
		InputChar:    "z",
	}
	if err := svc.ReportWordError(9527, secondReq); err != nil {
		t.Fatalf("second report failed: %v", err)
	}

	var row englishModel.EnglishWordErrorLog
	if err := db.Where("user_id = ? AND word_id = ?", 9527, word.ID).First(&row).Error; err != nil {
		t.Fatalf("query error log failed: %v", err)
	}

	if row.WrongCount != 2 {
		t.Fatalf("expected wrongCount=2, got=%d", row.WrongCount)
	}
	if row.CategoryID != secondReq.CategoryID || row.ChapterID != secondReq.ChapterID {
		t.Fatalf("expected latest category/chapter to be updated, got category=%d chapter=%d", row.CategoryID, row.ChapterID)
	}
	if row.LastWrongIndex != secondReq.WrongIndex {
		t.Fatalf("expected lastWrongIndex=%d, got=%d", secondReq.WrongIndex, row.LastWrongIndex)
	}
	if row.LastExpectedChar != "c" || row.LastInputChar != "z" {
		t.Fatalf("expected last chars to be c/z, got=%s/%s", row.LastExpectedChar, row.LastInputChar)
	}
}

func TestReportWordError_ReturnsErrorWhenWordMissing(t *testing.T) {
	db := setupEnglishWordReportErrorTestDB(t)
	originalDB := global.GVA_DB
	global.GVA_DB = db
	t.Cleanup(func() {
		global.GVA_DB = originalDB
	})

	svc := EnglishWordService{}
	err := svc.ReportWordError(9527, englishReq.ReportWordErrorReq{WordID: 999999, WrongIndex: 0})
	if err == nil {
		t.Fatal("expected missing-word error, got nil")
	}
	if !strings.Contains(err.Error(), "单词不存在") {
		t.Fatalf("expected missing-word error message, got: %v", err)
	}
}
