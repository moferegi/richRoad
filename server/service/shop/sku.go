package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type SkuService struct {
}

// CreateSku 创建sku记录
// Author [piexlmax](https://github.com/piexlmax)
func (skuService *SkuService) CreateSku(sku *shop.Sku) (err error) {
	err = global.GVA_DB.Create(sku).Error
	return err
}

// DeleteSku 删除sku记录
// Author [piexlmax](https://github.com/piexlmax)
func (skuService *SkuService) DeleteSku(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Sku{}, "id = ?", ID).Error
	return err
}

// DeleteSkuByIds 批量删除sku记录
// Author [piexlmax](https://github.com/piexlmax)
func (skuService *SkuService) DeleteSkuByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Sku{}, "id in ?", IDs).Error
	return err
}

// UpdateSku 更新sku记录
// Author [piexlmax](https://github.com/piexlmax)
func (skuService *SkuService) UpdateSku(sku shop.Sku) (err error) {
	err = global.GVA_DB.Model(&shop.Sku{}).Where("id = ?", sku.ID).Updates(&sku).Error
	return err
}

// GetSku 根据ID获取sku记录
// Author [piexlmax](https://github.com/piexlmax)
func (skuService *SkuService) GetSku(ID string) (sku shop.Sku, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sku).Error
	return
}

// GetSkuInfoList 分页获取sku记录
// Author [piexlmax](https://github.com/piexlmax)
func (skuService *SkuService) GetSkuInfoList(info shopReq.SkuSearch) (list []shop.Sku, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.Sku{})
	var skus []shop.Sku
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.GoodID != nil {
		db = db.Where("good_id = ?", *info.GoodID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&skus).Error
	return skus, total, err
}
