package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type GoodService struct {
}

// CreateGood 创建商品记录
// Author [piexlmax](https://github.com/piexlmax)
func (goodService *GoodService) CreateGood(good *shop.Good) (err error) {
	err = global.GVA_DB.Create(good).Error
	return err
}

// DeleteGood 删除商品记录
// Author [piexlmax](https://github.com/piexlmax)
func (goodService *GoodService) DeleteGood(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Good{}, "id = ?", ID).Error
	return err
}

// DeleteGoodByIds 批量删除商品记录
// Author [piexlmax](https://github.com/piexlmax)
func (goodService *GoodService) DeleteGoodByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Good{}, "id in ?", IDs).Error
	return err
}

// UpdateGood 更新商品记录
// Author [piexlmax](https://github.com/piexlmax)
func (goodService *GoodService) UpdateGood(good shop.Good) (err error) {
	err = global.GVA_DB.Model(&shop.Good{}).Where("id = ?", good.ID).Updates(&good).Error
	return err
}

// GetGood 根据ID获取商品记录
// Author [piexlmax](https://github.com/piexlmax)
func (goodService *GoodService) GetGood(ID string) (good shop.Good, err error) {
	err = global.GVA_DB.Where("id = ?", ID).Preload("SKUS").First(&good).Error
	return
}

// GetGoodInfoList 分页获取商品记录
// Author [piexlmax](https://github.com/piexlmax)
func (goodService *GoodService) GetGoodInfoList(info shopReq.GoodSearch) (list []shop.Good, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.Good{})
	var goods []shop.Good
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}
	if info.CategoryID != 0 {
		db = db.Where("category_id = ?", info.CategoryID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&goods).Error
	return goods, total, err
}
