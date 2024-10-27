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

// CreateCollect 创建收藏记录
// Author [piexlmax](https://github.com/piexlmax)
func (collectService *CollectService) CreateCollect(collect *client.Collect) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var coll client.Collect
		var s shop.Good
		sErr := tx.First(&s, "id = ?", collect.GoodID).Error
		if sErr != nil {
			return errors.New("商品不存在")
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
func (collectService *CollectService) DeleteCollect(ID string) (err error) {
	err = global.GVA_DB.Delete(&client.Collect{}, "id = ?", ID).Error
	return err
}

// DeleteCollectByIds 批量删除收藏记录
// Author [piexlmax](https://github.com/piexlmax)
func (collectService *CollectService) DeleteCollectByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]client.Collect{}, "id in ?", IDs).Error
	return err
}

// UpdateCollect 更新收藏记录
// Author [piexlmax](https://github.com/piexlmax)
func (collectService *CollectService) UpdateCollect(collect client.Collect) (err error) {
	err = global.GVA_DB.Model(&client.Collect{}).Where("id = ?", collect.ID).Updates(&collect).Error
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
	// 创建db
	db := global.GVA_DB.Model(&client.Collect{})
	var collects []client.Collect
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	db = db.Where("user_id = ?", info.UserID)

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&collects).Error

	var goodIDs []uint
	for _, collect := range collects {
		goodIDs = append(goodIDs, collect.GoodID)
	}

	var goods []shop.Good
	err = global.GVA_DB.Where("id in ?", goodIDs).Find(&goods).Error

	return goods, total, err
}
