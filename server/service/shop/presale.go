package shop

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"gorm.io/gorm"
)

type PresaleService struct{}

// GetPresaleGoodList 获取预售商品列表
func (s *PresaleService) GetPresaleGoodList(info shopReq.PresaleListRequest) (list []shop.Good, total int64, err error) {
	db := global.GVA_DB.Model(&shop.Good{}).Where("is_presale = ? AND presale_enabled = ? AND status = ?", true, true, true)

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if info.Limit > 0 {
		db = db.Limit(info.Limit)
	} else if info.PageSize > 0 {
		db = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1))
	}

	err = db.Order("presale_sort ASC, created_at DESC").Preload("SKUS").Find(&list).Error
	return
}

// CheckPresaleAvailable 检查预售商品是否可购买
func (s *PresaleService) CheckPresaleAvailable(goodID uint) (available bool, message string, err error) {
	var good shop.Good
	err = global.GVA_DB.Where("id = ?", goodID).First(&good).Error
	if err != nil {
		return false, "商品不存在", err
	}

	if good.IsPresale == nil || !*good.IsPresale {
		return true, "", nil // 非预售商品直接可购
	}

	if good.PresaleEnabled == nil || !*good.PresaleEnabled {
		return false, "此商品预售已结束", nil
	}

	// 检查时间
	now := time.Now()
	if good.PresaleStart != nil && now.Before(*good.PresaleStart) {
		return false, "预售尚未开始", nil
	}
	if good.PresaleEnd != nil && now.After(*good.PresaleEnd) {
		return false, "预售时间已结束", nil
	}

	// 检查数量
	if good.PresaleQty != nil && good.PresaleSold != nil && *good.PresaleSold >= *good.PresaleQty {
		return false, "预售数量已卖完", nil
	}

	return true, "", nil
}

// IncrementPresaleSold 原子增加预售已售数量(防超卖)
func (s *PresaleService) IncrementPresaleSold(tx *gorm.DB, goodID uint, qty int) error {
	if tx == nil {
		tx = global.GVA_DB
	}
	result := tx.Model(&shop.Good{}).
		Where("id = ? AND is_presale = ? AND presale_enabled = ? AND presale_sold + ? <= presale_qty", goodID, true, true, qty).
		Update("presale_sold", gorm.Expr("presale_sold + ?", qty))
	if result.RowsAffected == 0 {
		return errors.New("预售数量不足")
	}
	return result.Error
}

// GetPresaleParticipants 获取预售商品参与用户(通过订单)
func (s *PresaleService) GetPresaleParticipants(goodID uint, page, pageSize int) (list []shop.OrderRes, total int64, err error) {
	db := global.GVA_DB.Model(&shop.OrderRes{}).
		Where("is_presale = ?", true).
		Joins("JOIN shop_order_detail ON shop_order_detail.order_id = shop_order.id AND shop_order_detail.good_id = ?", goodID).
		Preload("Detail").Preload("Detail.Good").Preload("Detail.SKU")

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if pageSize > 0 {
		db = db.Limit(pageSize).Offset(pageSize * (page - 1))
	}
	err = db.Order("shop_order.created_at DESC").Find(&list).Error
	return
}
