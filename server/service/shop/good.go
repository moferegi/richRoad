package shop

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"gorm.io/gorm"
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
func (goodService *GoodService) GetGood(ID string, userID uint, authority uint) (good shop.Good, err error) {
	err = global.GVA_DB.Where("id = ?", ID).Preload("SKUS").First(&good).Error

	if err != nil {
		return good, errors.New("商品不存在")
	}
	// 浏览量+1
	err = global.GVA_DB.Model(&good).Where("id = ?", ID).UpdateColumn("view_num", gorm.Expr("view_num + ?", 1)).Error
	if err != nil {
		return good, err
	}
	if userID != 0 && authority != 888 {
		var history shop.History
		// 先查一下history表第一条是不是当前访问的这个 如果不是 则创建一条 并且清理掉之前的那条 一个用户最多保留30条历史记录
		err = global.GVA_DB.Where("good_id = ? AND user_id = ?", ID, userID).First(&history).Error
		id, _ := strconv.Atoi(ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				history = shop.History{
					GoodID: id,
					UserID: userID,
				}
				err = global.GVA_DB.Create(&history).Error
			} else {
				return good, err
			}
		} else {
			// 如果存在 则更新一下
			history.GoodID = id
			err = global.GVA_DB.Save(&history).Error
		}
	}
	return
}

func (goodService *GoodService) GetGoodHistory(userID uint) (goods []shop.Good, err error) {
	var histories []shop.History
	err = global.GVA_DB.Where("user_id = ?", userID).Order("updated_at DESC").Limit(30).Find(&histories).Error
	if err != nil {
		return nil, err
	}
	for _, history := range histories {
		var good shop.Good
		err = global.GVA_DB.Where("id = ?", history.GoodID).Preload("SKUS").First(&good).Error
		if err == nil {
			goods = append(goods, good)
		}
	}
	return goods, nil
}

// GetGoodInfoList 分页获取商品记录
// Author [piexlmax](https://github.com/piexlmax)
// GetGoodInfoList 分页获取商品记录
// Author `https://github.com/piexlmax`
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

	// 统一搜索逻辑：优先使用keyword，如果没有则使用title
	searchTerm := ""
	if info.Keyword != "" {
		searchTerm = info.Keyword
	} else if info.Title != "" {
		searchTerm = info.Title
	}

	if searchTerm != "" {
		// 在商品标题、描述中搜索，对于JSON字段tags使用JSON_UNQUOTE和JSON_SEARCH
		db = db.Where("title LIKE ? OR description LIKE ? OR JSON_SEARCH(tags, 'one', ?) IS NOT NULL",
			"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")
	}

	if info.CategoryID != 0 {
		db = db.Where("category_id = ?", info.CategoryID)
	} else if info.ExcludeHiddenCategories {
		// 排除showInUni=false的分类下的商品
		db = db.Where("category_id NOT IN (SELECT id FROM shop_category WHERE show_in_uni = ?)", false)
	}

	if info.Recommend != nil {
		db = db.Where("recommend = ?", info.Recommend)
	}

	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}

	if info.IsPresale != nil {
		db = db.Where("is_presale = ?", info.IsPresale)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 排序支持
	orderClause := "id desc"
	if info.OrderBy != "" {
		allowedCols := map[string]bool{"price": true, "sale_num": true, "view_num": true, "collect_num": true, "created_at": true, "sort": true, "presale_sort": true}
		if allowedCols[info.OrderBy] {
			dir := "asc"
			if info.OrderDir == "desc" {
				dir = "desc"
			}
			orderClause = info.OrderBy + " " + dir
		}
	}
	db = db.Order(orderClause)

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Preload("SKUS", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, good_id, inventory")
	}).Find(&goods).Error
	if err != nil {
		return goods, total, err
	}

	// 用最新的 tag 数据刷新商品 tags JSON（确保 color/nameI18n 同步）
	enrichGoodsTags(goods)

	return goods, total, err
}

// enrichGoodsTags 用数据库中最新的 Tag 信息刷新商品 tags JSON
func enrichGoodsTags(goods []shop.Good) {
	// 1. 收集所有 tag ID
	tagIDSet := map[uint]struct{}{}
	type tagEntry struct {
		ID       uint   `json:"ID"`
		Name     string `json:"name"`
		NameI18n string `json:"nameI18n"`
		Color    string `json:"color"`
	}
	goodTagsMap := make([][]tagEntry, len(goods))
	for i, g := range goods {
		if len(g.Tags) == 0 {
			continue
		}
		var entries []tagEntry
		if err := json.Unmarshal(g.Tags, &entries); err != nil {
			continue
		}
		goodTagsMap[i] = entries
		for _, e := range entries {
			if e.ID > 0 {
				tagIDSet[e.ID] = struct{}{}
			}
		}
	}
	if len(tagIDSet) == 0 {
		return
	}

	// 2. 批量查询最新 tag 数据
	ids := make([]uint, 0, len(tagIDSet))
	for id := range tagIDSet {
		ids = append(ids, id)
	}
	var dbTags []shop.Tag
	if err := global.GVA_DB.Where("id IN ?", ids).Find(&dbTags).Error; err != nil {
		return
	}
	tagMap := map[uint]shop.Tag{}
	for _, t := range dbTags {
		tagMap[t.ID] = t
	}

	// 3. 用最新数据替换并写回 JSON
	for i, entries := range goodTagsMap {
		if entries == nil {
			continue
		}
		changed := false
		for j, e := range entries {
			if t, ok := tagMap[e.ID]; ok {
				newEntry := tagEntry{ID: t.ID}
				if t.Name != nil {
					newEntry.Name = *t.Name
				}
				newEntry.NameI18n = t.NameI18n
				if t.Color != nil {
					newEntry.Color = *t.Color
				}
				if newEntry.Name != e.Name || newEntry.NameI18n != e.NameI18n || newEntry.Color != e.Color {
					changed = true
				}
				entries[j] = newEntry
			}
		}
		if changed {
			if data, err := json.Marshal(entries); err == nil {
				goods[i].Tags = data
			}
		}
	}
}
