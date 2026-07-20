package client

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
)

type VideoTagService struct{}

// CreateVideoTag 创建视频标签记录
func (s *VideoTagService) CreateVideoTag(ctx context.Context, tag *client.VideoTag) (err error) {
	err = global.GVA_DB.Create(tag).Error
	return err
}

// DeleteVideoTag 删除视频标签记录
func (s *VideoTagService) DeleteVideoTag(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&client.VideoTag{}, "id = ?", ID).Error
	return err
}

// DeleteVideoTagByIds 批量删除视频标签记录
func (s *VideoTagService) DeleteVideoTagByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]client.VideoTag{}, "id in ?", IDs).Error
	return err
}

// UpdateVideoTag 更新视频标签记录
func (s *VideoTagService) UpdateVideoTag(ctx context.Context, tag client.VideoTag) (err error) {
	err = global.GVA_DB.Model(&client.VideoTag{}).Where("id = ?", tag.ID).Updates(&tag).Error
	return err
}

// GetVideoTag 根据ID获取视频标签记录
func (s *VideoTagService) GetVideoTag(ctx context.Context, ID string) (tag client.VideoTag, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&tag).Error
	return
}

// GetVideoTagInfoList 分页获取视频标签记录
func (s *VideoTagService) GetVideoTagInfoList(ctx context.Context, info clientReq.VideoTagSearch) (list []client.VideoTag, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.VideoTag{})
	var tags []client.VideoTag
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

// GetVideoTagPublic 获取公开视频标签列表
func (s *VideoTagService) GetVideoTagPublic(ctx context.Context) (list []client.VideoTag, err error) {
	err = global.GVA_DB.Model(&client.VideoTag{}).Order("created_at DESC").Find(&list).Error
	return
}
