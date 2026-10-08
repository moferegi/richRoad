package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
)

type AdService struct{}

// --- 广告位置 ---

func (s *AdService) CreateAdPosition(p client.AdPosition) error {
	return global.GVA_DB.Create(&p).Error
}

func (s *AdService) DeleteAdPosition(id uint) error {
	return global.GVA_DB.Delete(&client.AdPosition{}, id).Error
}

func (s *AdService) UpdateAdPosition(p client.AdPosition) error {
	return global.GVA_DB.Model(&client.AdPosition{}).Where("id = ?", p.ID).Updates(&p).Error
}

func (s *AdService) GetAdPositionList(info clientReq.AdPositionSearch) (list []client.AdPosition, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.AdPosition{})

	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.PositionKey != "" {
		db = db.Where("position_key = ?", info.PositionKey)
	}
	if info.IsEnabled != nil {
		db = db.Where("is_enabled = ?", *info.IsEnabled)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("sort DESC, id ASC").Limit(limit).Offset(offset).Find(&list).Error
	return
}

// --- 广告视频 ---

func (s *AdService) CreateAdVideo(v client.AdVideo) error {
	return global.GVA_DB.Create(&v).Error
}

func (s *AdService) DeleteAdVideo(id uint) error {
	return global.GVA_DB.Delete(&client.AdVideo{}, id).Error
}

func (s *AdService) UpdateAdVideo(v client.AdVideo) error {
	return global.GVA_DB.Model(&client.AdVideo{}).Where("id = ?", v.ID).Updates(&v).Error
}

func (s *AdService) GetAdVideoList(info clientReq.AdVideoSearch) (list []client.AdVideo, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.AdVideo{})

	if info.PositionID > 0 {
		db = db.Where("position_id = ?", info.PositionID)
	}
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}
	if info.MediaType != "" {
		db = db.Where("media_type = ?", info.MediaType)
	}
	if info.IsEnabled != nil {
		db = db.Where("is_enabled = ?", *info.IsEnabled)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("sort DESC, id ASC").Limit(limit).Offset(offset).Find(&list).Error
	return
}

// GetAdByPosition 根据位置标记获取启用广告列表（Uni端用）
func (s *AdService) GetAdByPosition(positionKey string) (videos []client.AdVideo, err error) {
	var position client.AdPosition
	err = global.GVA_DB.Where("position_key = ? AND is_enabled = ?", positionKey, true).First(&position).Error
	if err != nil {
		return nil, err
	}

	err = global.GVA_DB.Where("position_id = ? AND is_enabled = ?", position.ID, true).
		Order("sort DESC, id ASC").Find(&videos).Error
	return
}

// --- 观看记录 ---

func (s *AdService) CreateWatchRecord(r client.AdWatchRecord) error {
	return global.GVA_DB.Create(&r).Error
}

func (s *AdService) GetWatchRecordList(info clientReq.AdWatchRecordSearch) (list []client.AdWatchRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.AdWatchRecord{})

	if info.UserID > 0 {
		db = db.Where("user_id = ?", info.UserID)
	}
	if info.AdVideoID > 0 {
		db = db.Where("ad_video_id = ?", info.AdVideoID)
	}
	if info.PositionID > 0 {
		db = db.Where("position_id = ?", info.PositionID)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id DESC").Limit(limit).Offset(offset).Find(&list).Error
	return
}