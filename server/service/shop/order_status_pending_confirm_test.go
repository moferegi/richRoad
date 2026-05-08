package shop

import (
	"strconv"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	shopModel "github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func setupShopOrderPendingConfirmTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}

	if err := db.AutoMigrate(
		&clientModel.ClientUser{},
		&clientModel.PointRecord{},
		&shopModel.Order{},
		&shopModel.OrderDetail{},
	); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	return db
}

func TestUpdateOrderStatusForUser_ToPendingConfirmWithQRCode(t *testing.T) {
	db := setupShopOrderPendingConfirmTestDB(t)
	prevDB := global.GVA_DB
	prevLog := global.GVA_LOG
	global.GVA_DB = db
	global.GVA_LOG = zap.NewNop()
	defer func() {
		global.GVA_DB = prevDB
		global.GVA_LOG = prevLog
	}()

	user := clientModel.ClientUser{}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	order := shopModel.Order{UserID: user.ID, Status: "0", PayMethod: "qrcode", TotalPrice: 0}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	svc := OrderService{}
	if err := svc.UpdateOrderStatusForUser(user.ID, strconv.Itoa(int(order.ID)), shopOrderStatusPendingConfirm); err != nil {
		t.Fatalf("update order status failed: %v", err)
	}

	var updated shopModel.Order
	if err := db.First(&updated, order.ID).Error; err != nil {
		t.Fatalf("query updated order failed: %v", err)
	}
	if updated.Status != shopOrderStatusPendingConfirm {
		t.Fatalf("expected status %s, got %s", shopOrderStatusPendingConfirm, updated.Status)
	}
}

func TestUpdateOrderStatusForUser_ToPendingConfirmRejectsNonQRCode(t *testing.T) {
	db := setupShopOrderPendingConfirmTestDB(t)
	prevDB := global.GVA_DB
	prevLog := global.GVA_LOG
	global.GVA_DB = db
	global.GVA_LOG = zap.NewNop()
	defer func() {
		global.GVA_DB = prevDB
		global.GVA_LOG = prevLog
	}()

	user := clientModel.ClientUser{}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	order := shopModel.Order{UserID: user.ID, Status: "0", PayMethod: "contact", TotalPrice: 0}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	svc := OrderService{}
	err := svc.UpdateOrderStatusForUser(user.ID, strconv.Itoa(int(order.ID)), shopOrderStatusPendingConfirm)
	if err == nil {
		t.Fatalf("expected error for non-qrcode payment method, got nil")
	}
	if !strings.Contains(err.Error(), "仅扫码支付订单可提交付款确认") {
		t.Fatalf("unexpected error: %v", err)
	}

	var updated shopModel.Order
	if err := db.First(&updated, order.ID).Error; err != nil {
		t.Fatalf("query updated order failed: %v", err)
	}
	if updated.Status != "0" {
		t.Fatalf("expected status to remain 0, got %s", updated.Status)
	}
}

func TestConfirmPayment_AllowsPendingConfirmStatus(t *testing.T) {
	db := setupShopOrderPendingConfirmTestDB(t)
	prevDB := global.GVA_DB
	prevLog := global.GVA_LOG
	global.GVA_DB = db
	global.GVA_LOG = zap.NewNop()
	defer func() {
		global.GVA_DB = prevDB
		global.GVA_LOG = prevLog
	}()

	user := clientModel.ClientUser{}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	order := shopModel.Order{UserID: user.ID, Status: shopOrderStatusPendingConfirm, PayMethod: "qrcode", TotalPrice: 0}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	svc := OrderService{}
	if err := svc.ConfirmPayment(strconv.Itoa(int(order.ID))); err != nil {
		t.Fatalf("confirm payment failed: %v", err)
	}

	var updated shopModel.Order
	if err := db.First(&updated, order.ID).Error; err != nil {
		t.Fatalf("query updated order failed: %v", err)
	}
	if updated.Status != "1" {
		t.Fatalf("expected status 1 after confirm payment, got %s", updated.Status)
	}
	if updated.PaidAt == nil {
		t.Fatalf("expected paid_at to be set after confirm payment")
	}
}

func TestConfirmPayment_RejectsNonPendingStatuses(t *testing.T) {
	db := setupShopOrderPendingConfirmTestDB(t)
	prevDB := global.GVA_DB
	prevLog := global.GVA_LOG
	global.GVA_DB = db
	global.GVA_LOG = zap.NewNop()
	defer func() {
		global.GVA_DB = prevDB
		global.GVA_LOG = prevLog
	}()

	user := clientModel.ClientUser{}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	order := shopModel.Order{UserID: user.ID, Status: "2", PayMethod: "qrcode", TotalPrice: 0}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	svc := OrderService{}
	err := svc.ConfirmPayment(strconv.Itoa(int(order.ID)))
	if err == nil {
		t.Fatalf("expected error for non pending status, got nil")
	}
	if !strings.Contains(err.Error(), "只能对待付款/待后台确认订单确认收款") {
		t.Fatalf("unexpected error: %v", err)
	}
}
