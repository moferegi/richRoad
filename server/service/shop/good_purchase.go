package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type GoodPurchaseService struct{}

// CreateGoodPurchase 创建进货记录
func (s *GoodPurchaseService) CreateGoodPurchase(purchase *shop.GoodPurchase) (err error) {
	purchase.TotalCost = purchase.Quantity * purchase.UnitCost
	err = global.GVA_DB.Create(purchase).Error
	return err
}

// DeleteGoodPurchase 删除进货记录
func (s *GoodPurchaseService) DeleteGoodPurchase(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.GoodPurchase{}, "id = ?", ID).Error
	return err
}

// UpdateGoodPurchase 更新进货记录
func (s *GoodPurchaseService) UpdateGoodPurchase(purchase shop.GoodPurchase) (err error) {
	purchase.TotalCost = purchase.Quantity * purchase.UnitCost
	err = global.GVA_DB.Model(&shop.GoodPurchase{}).Where("id = ?", purchase.ID).Updates(&purchase).Error
	return err
}

// GetGoodPurchase 根据ID获取进货记录
func (s *GoodPurchaseService) GetGoodPurchase(ID string) (purchase shop.GoodPurchase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).Preload("Good").Preload("Sku").First(&purchase).Error
	return
}

// GetGoodPurchaseList 分页获取进货记录列表
func (s *GoodPurchaseService) GetGoodPurchaseList(info shopReq.GoodPurchaseSearch) (list []shop.GoodPurchase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&shop.GoodPurchase{})

	if info.GoodID != nil {
		db = db.Where("good_id = ?", *info.GoodID)
	}
	if info.SkuID != nil {
		db = db.Where("sku_id = ?", *info.SkuID)
	}
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Order("created_at DESC").Preload("Good").Preload("Sku").Find(&list).Error
	return
}

// GetGoodPurchaseSummary 获取商品进货汇总
func (s *GoodPurchaseService) GetGoodPurchaseSummary(goodID uint) (totalQty int64, totalCost int64, err error) {
	type Result struct {
		TotalQty  int64
		TotalCost int64
	}
	var result Result
	err = global.GVA_DB.Model(&shop.GoodPurchase{}).
		Where("good_id = ?", goodID).
		Select("COALESCE(SUM(quantity),0) as total_qty, COALESCE(SUM(total_cost),0) as total_cost").
		Scan(&result).Error
	return result.TotalQty, result.TotalCost, err
}
