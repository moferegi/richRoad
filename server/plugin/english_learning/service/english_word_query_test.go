package service

import (
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	englishReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(127.0.0.1:9910)/gorm?charset=utf8mb4&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatalf("open dry-run db failed: %v", err)
	}

	return db
}

func TestBuildWordListPageQuery_UsesGroupAndAggregateOrder(t *testing.T) {
	db := newDryRunDB(t)
	var rows []model.EnglishWord
	stmt := buildWordListPageQuery(db, englishReq.WordListSearch{ChapterID: 123}).
		Limit(10).
		Offset(0).
		Find(&rows).Statement

	sql := strings.ToLower(stmt.SQL.String())
	requiredParts := []string{
		"group by",
		"min(cw.sort)",
		"max(cw.id)",
		"where cw.chapter_id = ?",
	}

	for _, part := range requiredParts {
		if !strings.Contains(sql, part) {
			t.Fatalf("expected sql to contain %q, got: %s", part, sql)
		}
	}

	if strings.Contains(sql, "select distinct") {
		t.Fatalf("page query should not use DISTINCT under only_full_group_by, got: %s", sql)
	}
}

func TestBuildWordListCountQuery_UsesDistinctID(t *testing.T) {
	db := newDryRunDB(t)
	var total int64
	stmt := buildWordListCountQuery(db, englishReq.WordListSearch{}).Count(&total).Statement

	sql := strings.ToLower(stmt.SQL.String())
	if !strings.Contains(sql, "distinct") {
		t.Fatalf("count query should keep DISTINCT for unique word ids, got: %s", sql)
	}
}
