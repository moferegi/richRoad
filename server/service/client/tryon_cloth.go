package client

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
)

const (
	tryonClothCategoryUpper    = "upper"
	tryonClothCategoryLower    = "lower"
	tryonClothCategoryOnepiece = "onepiece"
	tryonClothCategoryShoes    = "shoes"
)

var tryonClothCategoryLabel = map[string]string{
	tryonClothCategoryUpper:    "上装",
	tryonClothCategoryLower:    "下装",
	tryonClothCategoryOnepiece: "连体装",
	tryonClothCategoryShoes:    "鞋",
}

var tryonClothCategorySet = map[string]struct{}{
	tryonClothCategoryUpper:    {},
	tryonClothCategoryLower:    {},
	tryonClothCategoryOnepiece: {},
	tryonClothCategoryShoes:    {},
}

type TryonClothService struct{}

func normalizeTryonClothCategory(raw string) (string, bool) {
	category := strings.TrimSpace(strings.ToLower(raw))
	_, ok := tryonClothCategorySet[category]
	return category, ok
}

func normalizeTryonClothCategories(raw []string) []string {
	if len(raw) == 0 {
		return nil
	}
	result := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, item := range raw {
		category, ok := normalizeTryonClothCategory(item)
		if !ok {
			continue
		}
		if _, exists := seen[category]; exists {
			continue
		}
		seen[category] = struct{}{}
		result = append(result, category)
	}
	return result
}

func buildTryonClothDefaultName(category string) string {
	label := tryonClothCategoryLabel[category]
	if label == "" {
		label = "衣物"
	}
	return fmt.Sprintf("%s%d", label, time.Now().Unix())
}

// CreateTryonCloth 创建我的衣橱
func (s *TryonClothService) CreateTryonCloth(userID uint, req clientReq.CreateTryonClothReq) (cloth client.TryonCloth, err error) {
	if userID == 0 {
		return cloth, errors.New("loginRequired")
	}

	category, ok := normalizeTryonClothCategory(req.Category)
	if !ok {
		return cloth, errors.New("tryonClothCategoryInvalid")
	}

	image := strings.TrimSpace(req.Image)
	if image == "" {
		return cloth, errors.New("clothesImageRequired")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = buildTryonClothDefaultName(category)
	}

	cloth = client.TryonCloth{
		UserID:   userID,
		Name:     name,
		Category: category,
		Image:    image,
	}

	err = global.GVA_DB.Create(&cloth).Error
	return
}

// UpdateTryonCloth 更新我的衣橱
func (s *TryonClothService) UpdateTryonCloth(userID uint, authorityID uint, req clientReq.UpdateTryonClothReq) error {
	if req.ID == 0 {
		return errors.New("invalidID")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errors.New("tryonClothNameRequired")
	}

	updateData := map[string]interface{}{
		"name": name,
	}
	if strings.TrimSpace(req.Category) != "" {
		category, ok := normalizeTryonClothCategory(req.Category)
		if !ok {
			return errors.New("tryonClothCategoryInvalid")
		}
		updateData["category"] = category
	}

	db := global.GVA_DB.Model(&client.TryonCloth{}).Where("id = ?", req.ID)
	if authorityID != 888 {
		db = db.Where("user_id = ?", userID)
	}

	result := db.Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("tryonClothNotFoundOrNoPermission")
	}
	return nil
}

// DeleteTryonCloth 删除单个我的衣橱
func (s *TryonClothService) DeleteTryonCloth(userID uint, authorityID uint, id uint) error {
	if id == 0 {
		return errors.New("invalidID")
	}

	db := global.GVA_DB.Where("id = ?", id)
	if authorityID != 888 {
		db = db.Where("user_id = ?", userID)
	}

	result := db.Delete(&client.TryonCloth{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("tryonClothNotFoundOrNoPermission")
	}
	return nil
}

// DeleteTryonClothByIds 批量删除我的衣橱
func (s *TryonClothService) DeleteTryonClothByIds(userID uint, authorityID uint, ids []uint) error {
	if len(ids) == 0 {
		return errors.New("invalidIDs")
	}

	db := global.GVA_DB.Where("id IN ?", ids)
	if authorityID != 888 {
		db = db.Where("user_id = ?", userID)
	}

	result := db.Delete(&client.TryonCloth{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("tryonClothNotFoundOrNoPermission")
	}
	return nil
}

// GetMyTryonClothList 获取我的衣橱列表
func (s *TryonClothService) GetMyTryonClothList(userID uint, categories []string) (list []client.TryonCloth, err error) {
	if userID == 0 {
		return list, errors.New("loginRequired")
	}

	db := global.GVA_DB.Where("user_id = ?", userID)
	normalized := normalizeTryonClothCategories(categories)
	if len(normalized) > 0 {
		db = db.Where("category IN ?", normalized)
	}

	err = db.Order("id desc").Find(&list).Error
	return
}

// GetTryonClothList 获取我的衣橱列表（管理端）
func (s *TryonClothService) GetTryonClothList(info clientReq.TryonClothSearch) (list []client.TryonCloth, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.TryonCloth{})

	if info.UserID != 0 {
		db = db.Where("user_id = ?", info.UserID)
	}
	if strings.TrimSpace(info.Name) != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(info.Name)+"%")
	}
	if strings.TrimSpace(info.Category) != "" {
		category, ok := normalizeTryonClothCategory(info.Category)
		if !ok {
			return list, total, errors.New("tryonClothCategoryInvalid")
		}
		db = db.Where("category = ?", category)
	}
	if info.StartCreatedAt != nil {
		db = db.Where("created_at >= ?", info.StartCreatedAt)
	}
	if info.EndCreatedAt != nil {
		db = db.Where("created_at <= ?", info.EndCreatedAt)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&list).Error
	return
}
