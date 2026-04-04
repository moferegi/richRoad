package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type BannerService struct {
}

// CreateBanner 创建轮播图记录
// Author [piexlmax](https://github.com/piexlmax)
func (bannerService *BannerService) CreateBanner(banner *shop.Banner) (err error) {
	err = global.GVA_DB.Create(banner).Error
	return err
}

// DeleteBanner 删除轮播图记录
// Author [piexlmax](https://github.com/piexlmax)
func (bannerService *BannerService) DeleteBanner(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Banner{}, "id = ?", ID).Error
	return err
}

// DeleteBannerByIds 批量删除轮播图记录
// Author [piexlmax](https://github.com/piexlmax)
func (bannerService *BannerService) DeleteBannerByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Banner{}, "id in ?", IDs).Error
	return err
}

// UpdateBanner 更新轮播图记录
// Author [piexlmax](https://github.com/piexlmax)
func (bannerService *BannerService) UpdateBanner(banner shop.Banner) (err error) {
	err = global.GVA_DB.Model(&shop.Banner{}).Where("id = ?", banner.ID).Updates(&banner).Error
	return err
}

// GetBanner 根据ID获取轮播图记录
// Author [piexlmax](https://github.com/piexlmax)
func (bannerService *BannerService) GetBanner(ID string) (banner shop.Banner, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&banner).Error
	return
}

// GetBannerInfoList 分页获取轮播图记录
// Author [piexlmax](https://github.com/piexlmax)
func (bannerService *BannerService) GetBannerInfoList(info shopReq.BannerSearch) (list []shop.Banner, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.Banner{})
	var banners []shop.Banner
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.IsEnabled != nil {
		db = db.Where("is_enabled = ?", *info.IsEnabled)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&banners).Error
	return banners, total, err
}
