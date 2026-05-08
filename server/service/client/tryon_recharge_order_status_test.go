package client

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTryonRechargeOrderStatusTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}

	if err := db.AutoMigrate(
		&clientModel.ClientUser{},
		&clientModel.PointRecord{},
		&clientModel.TryonRechargeOrder{},
	); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	return db
}

func TestSubmitTryonRechargeOrderPayment_TransitionsFromPendingToReview(t *testing.T) {
	db := setupTryonRechargeOrderStatusTestDB(t)
	prevDB := global.GVA_DB
	global.GVA_DB = db
	defer func() {
		global.GVA_DB = prevDB
	}()

	user := clientModel.ClientUser{TryonPoint: 3}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	order := clientModel.TryonRechargeOrder{
		UserID:     user.ID,
		Status:     tryonRechargeOrderStatusPending,
		Points:     7,
		Amount:     990,
		Currency:   "CNY",
		PayMethod:  "qrcode",
		OutTradeNo: "sy10001AA",
		CloseTime:  time.Now().Add(time.Hour),
	}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	svc := TryonRechargeOrderService{}
	if err := svc.SubmitTryonRechargeOrderPayment(user.ID, order.ID); err != nil {
		t.Fatalf("submit payment confirmation failed: %v", err)
	}

	var updated clientModel.TryonRechargeOrder
	if err := db.First(&updated, order.ID).Error; err != nil {
		t.Fatalf("query updated order failed: %v", err)
	}
	if updated.Status != tryonRechargeOrderStatusReview {
		t.Fatalf("expected status %s, got %s", tryonRechargeOrderStatusReview, updated.Status)
	}

	if err := svc.SubmitTryonRechargeOrderPayment(user.ID, order.ID); err != nil {
		t.Fatalf("submit payment confirmation should be idempotent for status 8, got error: %v", err)
	}
}

func TestSubmitTryonRechargeOrderPayment_RejectsNonQRCodeOrder(t *testing.T) {
	db := setupTryonRechargeOrderStatusTestDB(t)
	prevDB := global.GVA_DB
	global.GVA_DB = db
	defer func() {
		global.GVA_DB = prevDB
	}()

	user := clientModel.ClientUser{}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	order := clientModel.TryonRechargeOrder{
		UserID:     user.ID,
		Status:     tryonRechargeOrderStatusPending,
		Points:     9,
		Amount:     1990,
		Currency:   "CNY",
		PayMethod:  "contact",
		OutTradeNo: "sy10002BB",
		CloseTime:  time.Now().Add(time.Hour),
	}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	svc := TryonRechargeOrderService{}
	err := svc.SubmitTryonRechargeOrderPayment(user.ID, order.ID)
	if err == nil {
		t.Fatalf("expected error for non-qrcode order, got nil")
	}
	if !strings.Contains(err.Error(), "仅扫码支付订单可提交付款确认") {
		t.Fatalf("unexpected error: %v", err)
	}

	var updated clientModel.TryonRechargeOrder
	if err := db.First(&updated, order.ID).Error; err != nil {
		t.Fatalf("query updated order failed: %v", err)
	}
	if updated.Status != tryonRechargeOrderStatusPending {
		t.Fatalf("expected status to remain %s, got %s", tryonRechargeOrderStatusPending, updated.Status)
	}
}

func TestConfirmTryonRechargeOrderPayment_AllowsReviewAndCreditsTryonPoints(t *testing.T) {
	db := setupTryonRechargeOrderStatusTestDB(t)
	prevDB := global.GVA_DB
	global.GVA_DB = db
	defer func() {
		global.GVA_DB = prevDB
	}()

	user := clientModel.ClientUser{TryonPoint: 10}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	order := clientModel.TryonRechargeOrder{
		UserID:     user.ID,
		Status:     tryonRechargeOrderStatusReview,
		Points:     6,
		Amount:     2990,
		Currency:   "CNY",
		PayMethod:  "qrcode",
		OutTradeNo: "sy10003CC",
		CloseTime:  time.Now().Add(time.Hour),
	}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	svc := TryonRechargeOrderService{}
	if err := svc.ConfirmTryonRechargeOrderPayment(context.Background(), order.ID, "  paid by admin  "); err != nil {
		t.Fatalf("confirm payment failed: %v", err)
	}

	var updatedOrder clientModel.TryonRechargeOrder
	if err := db.First(&updatedOrder, order.ID).Error; err != nil {
		t.Fatalf("query updated order failed: %v", err)
	}
	if updatedOrder.Status != tryonRechargeOrderStatusPaid {
		t.Fatalf("expected status %s, got %s", tryonRechargeOrderStatusPaid, updatedOrder.Status)
	}
	if updatedOrder.PaidAt == nil {
		t.Fatalf("expected paid_at to be set")
	}
	if updatedOrder.Remark != "paid by admin" {
		t.Fatalf("expected trimmed remark, got %q", updatedOrder.Remark)
	}

	var updatedUser clientModel.ClientUser
	if err := db.First(&updatedUser, user.ID).Error; err != nil {
		t.Fatalf("query updated user failed: %v", err)
	}
	if updatedUser.TryonPoint != 16 {
		t.Fatalf("expected tryon points to be 16, got %d", updatedUser.TryonPoint)
	}

	var record clientModel.PointRecord
	if err := db.Where("related_order_id = ? AND operation_type = ?", int(order.ID), "tryon_recharge").First(&record).Error; err != nil {
		t.Fatalf("query point record failed: %v", err)
	}
	if record.AssetType == nil || *record.AssetType != clientModel.AssetTypeTryonPoint {
		t.Fatalf("expected asset_type %s, got %+v", clientModel.AssetTypeTryonPoint, record.AssetType)
	}
	if record.PointChange == nil || *record.PointChange != 6 {
		t.Fatalf("expected point change 6, got %+v", record.PointChange)
	}
}

func TestConfirmTryonRechargeOrderPayment_PaidOrderIsIdempotent(t *testing.T) {
	db := setupTryonRechargeOrderStatusTestDB(t)
	prevDB := global.GVA_DB
	global.GVA_DB = db
	defer func() {
		global.GVA_DB = prevDB
	}()

	now := time.Now()
	user := clientModel.ClientUser{TryonPoint: 9}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	order := clientModel.TryonRechargeOrder{
		UserID:     user.ID,
		Status:     tryonRechargeOrderStatusPaid,
		Points:     3,
		Amount:     990,
		Currency:   "CNY",
		PayMethod:  "qrcode",
		OutTradeNo: "sy10004DD",
		CloseTime:  now.Add(time.Hour),
		PaidAt:     &now,
	}
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("create order failed: %v", err)
	}

	svc := TryonRechargeOrderService{}
	if err := svc.ConfirmTryonRechargeOrderPayment(context.Background(), order.ID, "noop"); err != nil {
		t.Fatalf("confirm payment should be idempotent for paid order, got error: %v", err)
	}

	var updatedUser clientModel.ClientUser
	if err := db.First(&updatedUser, user.ID).Error; err != nil {
		t.Fatalf("query updated user failed: %v", err)
	}
	if updatedUser.TryonPoint != 9 {
		t.Fatalf("expected tryon points to remain 9, got %d", updatedUser.TryonPoint)
	}

	var count int64
	if err := db.Model(&clientModel.PointRecord{}).Count(&count).Error; err != nil {
		t.Fatalf("count point records failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no point record for idempotent paid confirmation, got %d", count)
	}
}
