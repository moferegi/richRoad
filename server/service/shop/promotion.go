package shop

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type PromotionService struct{}

// CreatePromotion 创建促销信息记录
// Author [yourname](https://github.com/yourname)
func (PromoService *PromotionService) CreatePromotion(ctx context.Context, Promo *shop.Promotion) (err error) {
	err = global.GVA_DB.Create(Promo).Error
	return err
}

// DeletePromotion 删除促销信息记录
// Author [yourname](https://github.com/yourname)
func (PromoService *PromotionService) DeletePromotion(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Promotion{}, "id = ?", ID).Error
	return err
}

// DeletePromotionByIds 批量删除促销信息记录
// Author [yourname](https://github.com/yourname)
func (PromoService *PromotionService) DeletePromotionByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Promotion{}, "id in ?", IDs).Error
	return err
}

// UpdatePromotion 更新促销信息记录
// Author [yourname](https://github.com/yourname)
func (PromoService *PromotionService) UpdatePromotion(ctx context.Context, Promo shop.Promotion) (err error) {
	err = global.GVA_DB.Model(&shop.Promotion{}).Where("id = ?", Promo.ID).Updates(&Promo).Error
	return err
}

// GetPromotion 根据ID获取促销信息记录
// Author [yourname](https://github.com/yourname)
func (PromoService *PromotionService) GetPromotion(ctx context.Context, ID string) (Promo shop.Promotion, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&Promo).Error
	return
}

// GetPromotionInfoList 分页获取促销信息记录
// Author [yourname](https://github.com/yourname)
func (PromoService *PromotionService) GetPromotionInfoList(ctx context.Context, info shopReq.PromotionSearch) (list []shop.Promotion, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.Promotion{})
	var Promos []shop.Promotion
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.Title != nil && *info.Title != "" {
		db = db.Where("title = ?", *info.Title)
	}
	if info.Category != nil && *info.Category != "" {
		db = db.Where("category LIKE ?", "%"+*info.Category+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["ID"] = true
	orderMap["CreatedAt"] = true
	orderMap["promotion_name"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&Promos).Error
	return Promos, total, err
}

func (PromoService *PromotionService) GetPromotionPublic(ctx context.Context) (promotion shop.Promotion, err error) {
	// 查询第一个启用状态(IsAction=true)的促销记录
	err = global.GVA_DB.WithContext(ctx).
		Where("is_action = ?", true). // 查询启用状态的记录
		First(&promotion).Error

	return promotion, err
}
