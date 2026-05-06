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

type TryonModelService struct{}

// CreateTryonModel 创建我的模特
func (s *TryonModelService) CreateTryonModel(userID uint, req clientReq.CreateTryonModelReq) (model client.TryonModel, err error) {
	if userID == 0 {
		return model, errors.New("用户未登录")
	}

	image := strings.TrimSpace(req.Image)
	if image == "" {
		return model, errors.New("模特图片不能为空")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = fmt.Sprintf("我的模特%d", time.Now().Unix())
	}

	model = client.TryonModel{
		UserID: userID,
		Name:   name,
		Image:  image,
	}

	err = global.GVA_DB.Create(&model).Error
	return
}

// UpdateTryonModel 重命名我的模特
func (s *TryonModelService) UpdateTryonModel(userID uint, authorityID uint, req clientReq.UpdateTryonModelReq) error {
	if req.ID == 0 {
		return errors.New("ID参数错误")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errors.New("模特名称不能为空")
	}

	db := global.GVA_DB.Model(&client.TryonModel{}).Where("id = ?", req.ID)
	if authorityID != 888 {
		db = db.Where("user_id = ?", userID)
	}

	result := db.Updates(map[string]interface{}{
		"name": name,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("模特不存在或无权限")
	}

	return nil
}

// DeleteTryonModel 删除单个我的模特
func (s *TryonModelService) DeleteTryonModel(userID uint, authorityID uint, id uint) error {
	if id == 0 {
		return errors.New("ID参数错误")
	}

	db := global.GVA_DB.Where("id = ?", id)
	if authorityID != 888 {
		db = db.Where("user_id = ?", userID)
	}

	result := db.Delete(&client.TryonModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("模特不存在或无权限")
	}

	return nil
}

// DeleteTryonModelByIds 批量删除我的模特
func (s *TryonModelService) DeleteTryonModelByIds(userID uint, authorityID uint, ids []uint) error {
	if len(ids) == 0 {
		return errors.New("IDs参数错误")
	}

	db := global.GVA_DB.Where("id IN ?", ids)
	if authorityID != 888 {
		db = db.Where("user_id = ?", userID)
	}

	result := db.Delete(&client.TryonModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("模特不存在或无权限")
	}

	return nil
}

// GetMyTryonModelList 获取我的模特列表
func (s *TryonModelService) GetMyTryonModelList(userID uint) (list []client.TryonModel, err error) {
	if userID == 0 {
		return list, errors.New("用户未登录")
	}

	err = global.GVA_DB.Where("user_id = ?", userID).Order("id desc").Find(&list).Error
	return
}

// GetTryonModelList 获取模特列表（管理端）
func (s *TryonModelService) GetTryonModelList(info clientReq.TryonModelSearch) (list []client.TryonModel, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&client.TryonModel{})
	if info.UserID != 0 {
		db = db.Where("user_id = ?", info.UserID)
	}
	if strings.TrimSpace(info.Name) != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(info.Name)+"%")
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
