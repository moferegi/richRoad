package shop

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type SkuSpecService struct{}

// CreateSkuSpec 创建SKU规格字典
func (s *SkuSpecService) CreateSkuSpec(ctx context.Context, spec *shop.SkuSpec) (err error) {
	err = global.GVA_DB.Create(spec).Error
	return err
}

// DeleteSkuSpec 删除SKU规格字典
func (s *SkuSpecService) DeleteSkuSpec(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.SkuSpec{}, "id = ?", ID).Error
	return err
}

// DeleteSkuSpecByIds 批量删除SKU规格字典
func (s *SkuSpecService) DeleteSkuSpecByIds(ctx context.Context, IDs []int) (err error) {
	err = global.GVA_DB.Delete(&[]shop.SkuSpec{}, "id in ?", IDs).Error
	return err
}

// UpdateSkuSpec 更新SKU规格字典
func (s *SkuSpecService) UpdateSkuSpec(ctx context.Context, spec shop.SkuSpec) (err error) {
	err = global.GVA_DB.Model(&shop.SkuSpec{}).Where("id = ?", spec.ID).Updates(&spec).Error
	return err
}

// GetSkuSpec 根据ID获取SKU规格字典
func (s *SkuSpecService) GetSkuSpec(ctx context.Context, ID string) (spec shop.SkuSpec, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&spec).Error
	return
}

// GetSkuSpecInfoList 分页获取SKU规格字典
func (s *SkuSpecService) GetSkuSpecInfoList(ctx context.Context, info shopReq.SkuSpecSearch) (list []shop.SkuSpec, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&shop.SkuSpec{})
	if info.Type != "" {
		db = db.Where("type = ?", info.Type)
	}
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["id"] = true
	orderMap["created_at"] = true
	orderMap["sort"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	} else {
		db = db.Order("sort asc, id desc")
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&list).Error
	return list, total, err
}

// GetAllSkuSpecs 获取所有SKU规格字典（不分页，用于下拉选择）
func (s *SkuSpecService) GetAllSkuSpecs(ctx context.Context, specType string) (list []shop.SkuSpec, err error) {
	db := global.GVA_DB.Model(&shop.SkuSpec{})
	if specType != "" {
		db = db.Where("type = ?", specType)
	}
	err = db.Order("sort asc, id desc").Find(&list).Error
	return
}
