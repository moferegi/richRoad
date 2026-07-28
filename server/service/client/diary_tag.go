package client

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
)

type DiaryTagService struct{}

// CreateDiaryTag 创建日记标签
func (s *DiaryTagService) CreateDiaryTag(ctx context.Context, tag *client.DiaryTag) (err error) {
	err = global.GVA_DB.Create(tag).Error
	return err
}

// DeleteDiaryTag 删除日记标签
func (s *DiaryTagService) DeleteDiaryTag(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&client.DiaryTag{}, "id = ?", ID).Error
	return err
}

// DeleteDiaryTagByIds 批量删除日记标签
func (s *DiaryTagService) DeleteDiaryTagByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]client.DiaryTag{}, "id in ?", IDs).Error
	return err
}

// UpdateDiaryTag 更新日记标签
func (s *DiaryTagService) UpdateDiaryTag(ctx context.Context, tag client.DiaryTag) (err error) {
	err = global.GVA_DB.Model(&client.DiaryTag{}).Where("id = ?", tag.ID).Updates(&tag).Error
	return err
}

// GetDiaryTag 根据ID获取日记标签
func (s *DiaryTagService) GetDiaryTag(ctx context.Context, ID string) (tag client.DiaryTag, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&tag).Error
	return
}

// GetDiaryTagInfoList 分页获取日记标签
func (s *DiaryTagService) GetDiaryTagInfoList(ctx context.Context, info clientReq.DiaryTagSearch) (list []client.DiaryTag, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.DiaryTag{})
	var tags []client.DiaryTag
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	orderStr := "created_at DESC"
	if info.Sort != "" {
		orderStr = info.Sort
		if info.Order == "asc" {
			orderStr += " ASC"
		} else {
			orderStr += " DESC"
		}
	}
	err = db.Order(orderStr).Find(&tags).Error
	return tags, total, err
}

// GetDiaryTagPublic 获取公开日记标签列表
func (s *DiaryTagService) GetDiaryTagPublic(ctx context.Context) (list []client.DiaryTag, err error) {
	err = global.GVA_DB.Model(&client.DiaryTag{}).Order("created_at DESC").Find(&list).Error
	return
}