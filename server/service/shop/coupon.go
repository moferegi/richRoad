package shop

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type CouponService struct{}

// CreateCoupon 创建优惠券记录
// Author [yourname](https://github.com/yourname)
func (CouService *CouponService) CreateCoupon(ctx context.Context, Cou *shop.Coupon) (err error) {
	err = global.GVA_DB.Create(Cou).Error
	return err
}

// DeleteCoupon 删除优惠券记录
// Author [yourname](https://github.com/yourname)
func (CouService *CouponService) DeleteCoupon(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Coupon{}, "id = ?", ID).Error
	return err
}

// DeleteCouponByIds 批量删除优惠券记录
// Author [yourname](https://github.com/yourname)
func (CouService *CouponService) DeleteCouponByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Coupon{}, "id in ?", IDs).Error
	return err
}

// UpdateCoupon 更新优惠券记录
// Author [yourname](https://github.com/yourname)
func (CouService *CouponService) UpdateCoupon(ctx context.Context, Cou shop.Coupon) (err error) {
	err = global.GVA_DB.Model(&shop.Coupon{}).Where("id = ?", Cou.ID).Updates(&Cou).Error
	return err
}

// GetCoupon 根据ID获取优惠券记录
// Author [yourname](https://github.com/yourname)
func (CouService *CouponService) GetCoupon(ctx context.Context, ID string) (Cou shop.Coupon, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&Cou).Error
	return
}

// GetCouponInfoList 分页获取优惠券记录
// Author [yourname](https://github.com/yourname)
func (CouService *CouponService) GetCouponInfoList(ctx context.Context, info shopReq.CouponSearch) (list []shop.Coupon, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.Coupon{})
	var Cous []shop.Coupon
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	if info.Type != "" {
		db = db.Where("type = ?", info.Type)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["ID"] = true
	orderMap["CreatedAt"] = true
	orderMap["name"] = true
	orderMap["start_time"] = true
	orderMap["end_time"] = true
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

	err = db.Find(&Cous).Error
	return Cous, total, err
}
func (CouService *CouponService) GetCouponDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	productID := make([]map[string]any, 0)

	global.GVA_DB.Table("shop_good").Where("deleted_at IS NULL").Select("title as label,id as value").Scan(&productID)
	res["productID"] = productID
	return
}
func (CouService *CouponService) GetCouponPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
