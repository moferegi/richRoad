package client

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	"gorm.io/gorm"
)

type CollectService struct {
}

func scopeCollectOwner(db *gorm.DB, userID uint, allowAll bool) (*gorm.DB, error) {
	if allowAll {
		return db, nil
	}
	if userID == 0 {
		return db, errors.New("noPermission")
	}
	return db.Where("user_id = ?", userID), nil
}

// CreateCollect 创建收藏记录
// Author [piexlmax](https://github.com/piexlmax)
func (collectService *CollectService) CreateCollect(collect *client.Collect) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var coll client.Collect
		var s shop.Good
		sErr := tx.
			Where("id = ? AND status = ?", collect.GoodID, true).
			Where("category_id IS NULL OR category_id NOT IN (SELECT id FROM shop_category WHERE show_in_uni = ?)", false).
			First(&s).Error
		if sErr != nil {
			return errors.New("seriesNotFound")
		}
		ferr := tx.First(&coll, "user_id = ? and good_id = ?", collect.UserID, collect.GoodID).Error
		if ferr != nil {
			e := tx.Create(collect).Error
			if e != nil {
				return e
			}
			e = tx.Model(&s).Update("collect_num", gorm.Expr("collect_num + ?", 1)).Error
			return e
		}
		e := tx.Unscoped().Delete(&coll, "id = ?", coll.ID).Error
		if e != nil {
			return e
		}
		e = tx.Model(&s).Update("collect_num", gorm.Expr("collect_num - ?", 1)).Error
		return e
	})
}

// DeleteCollect 删除收藏记录
// Author [piexlmax](https://github.com/piexlmax)
func (collectService *CollectService) DeleteCollect(ID string, userID uint, allowAll bool) (err error) {
	db, err := scopeCollectOwner(global.GVA_DB.Where("id = ?", ID), userID, allowAll)
	if err != nil {
		return err
	}
	err = db.Delete(&client.Collect{}).Error
	return err
}

// DeleteCollectByIds 批量删除收藏记录
// Author [piexlmax](https://github.com/piexlmax)
func (collectService *CollectService) DeleteCollectByIds(IDs []string, userID uint, allowAll bool) (err error) {
	db, err := scopeCollectOwner(global.GVA_DB.Where("id in ?", IDs), userID, allowAll)
	if err != nil {
		return err
	}
	err = db.Delete(&[]client.Collect{}).Error
	return err
}

// UpdateCollect 更新收藏记录
// Author [piexlmax](https://github.com/piexlmax)
func (collectService *CollectService) UpdateCollect(collect client.Collect, userID uint, allowAll bool) (err error) {
	db, err := scopeCollectOwner(global.GVA_DB.Model(&client.Collect{}).Where("id = ?", collect.ID), userID, allowAll)
	if err != nil {
		return err
	}
	if !allowAll {
		collect.UserID = userID
	}
	err = db.Updates(&collect).Error
	return err
}

// GetCollect 根据ID获取收藏记录
// Author [piexlmax](https://github.com/piexlmax)
func (collectService *CollectService) GetCollect(userID uint, ID string) (ok bool, err error) {
	var collect client.Collect
	ok = true
	err = global.GVA_DB.Where("good_id = ? and user_id = ?", ID, userID).First(&collect).Error
	if err != nil {
		ok = false
	}
	return ok, nil
}

// GetCollectInfoList 分页获取收藏记录
// Author [piexlmax](https://github.com/piexlmax)
func (collectService *CollectService) GetCollectInfoList(info clientReq.CollectSearch) (list []shop.Good, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&shop.Good{}).
		Joins("JOIN client_collect ON client_collect.good_id = shop_good.id").
		Where("client_collect.user_id = ?", info.UserID).
		Where("client_collect.deleted_at IS NULL").
		Where("shop_good.status = ?", true).
		Where("shop_good.category_id IS NULL OR shop_good.category_id NOT IN (SELECT id FROM shop_category WHERE show_in_uni = ?)", false)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("client_collect.created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Order("client_collect.updated_at DESC").Find(&list).Error
	return list, total, err
}
