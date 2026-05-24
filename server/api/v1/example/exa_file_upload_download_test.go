package example

import (
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	"github.com/glebarez/sqlite"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func setupSignURLTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&shop.Category{}, &shop.Good{}, &shop.Sku{}, &shop.Order{}, &shop.OrderDetail{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	previousDB := global.GVA_DB
	global.GVA_DB = db
	t.Cleanup(func() { global.GVA_DB = previousDB })
	return db
}

func TestPublicSignMediaPathAllowsPreviewMedia(t *testing.T) {
	db := setupSignURLTestDB(t)
	active := true
	price := 100.0
	good := shop.Good{
		Status: &active,
		Price:  &price,
		Banner: datatypes.JSON([]byte(`[{"src":"/media/public-preview.mp4"}]`)),
	}
	if err := db.Create(&good).Error; err != nil {
		t.Fatalf("create good: %v", err)
	}

	if !isPublicSignMediaPathAllowed("/media/public-preview.mp4", 0) {
		t.Fatalf("public preview media should be signable")
	}
}

func TestPublicSignMediaPathRequiresPaidOrderForProtectedSkuMedia(t *testing.T) {
	db := setupSignURLTestDB(t)
	active := true
	price := 100.0
	good := shop.Good{Status: &active, Price: &price}
	if err := db.Create(&good).Error; err != nil {
		t.Fatalf("create good: %v", err)
	}
	sku := shop.Sku{
		GoodID: good.ID,
		Name:   "episode",
		Specs:  datatypes.JSON([]byte(`[{"label":"videoUrl","value":"/media/paid-episode.mp4"}]`)),
	}
	if err := db.Create(&sku).Error; err != nil {
		t.Fatalf("create sku: %v", err)
	}

	if isPublicSignMediaPathAllowed("/media/paid-episode.mp4", 1001) {
		t.Fatalf("unpaid protected media should not be signable")
	}

	now := time.Now()
	order := shop.Order{UserID: 1001, Status: "1", PaidAt: &now}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}
	detail := shop.OrderDetail{OrderID: order.ID, GoodID: good.ID, SKUID: sku.ID, Quantity: 1, Price: 100}
	if err := db.Create(&detail).Error; err != nil {
		t.Fatalf("create order detail: %v", err)
	}

	if !isPublicSignMediaPathAllowed("/media/paid-episode.mp4", 1001) {
		t.Fatalf("paid protected media should be signable")
	}
	if isPublicSignMediaPathAllowed("/media/paid-episode.mp4", 1002) {
		t.Fatalf("another user should not be able to sign protected media")
	}
}
