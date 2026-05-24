package shop

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type KefuService struct{}

// CreateKefu 创建客服记录
// Author [yourname](https://github.com/yourname)
func (kefuService *KefuService) CreateKefu(ctx context.Context, kefu *shop.Kefu) (err error) {
	err = global.GVA_DB.Create(kefu).Error
	return err
}

// DeleteKefu 删除客服记录
// Author [yourname](https://github.com/yourname)
func (kefuService *KefuService) DeleteKefu(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Kefu{}, "id = ?", ID).Error
	return err
}

// DeleteKefuByIds 批量删除客服记录
// Author [yourname](https://github.com/yourname)
func (kefuService *KefuService) DeleteKefuByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Kefu{}, "id in ?", IDs).Error
	return err
}

// UpdateKefu 更新客服记录
// Author [yourname](https://github.com/yourname)
func (kefuService *KefuService) UpdateKefu(ctx context.Context, kefu shop.Kefu) (err error) {
	err = global.GVA_DB.Model(&shop.Kefu{}).Where("id = ?", kefu.ID).Updates(&kefu).Error
	return err
}

// GetKefu 根据ID获取客服记录
// Author [yourname](https://github.com/yourname)
func (kefuService *KefuService) GetKefu(ctx context.Context, ID string) (kefu shop.Kefu, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&kefu).Error
	return
}

// GetKefuInfoList 分页获取客服记录
// Author [yourname](https://github.com/yourname)
func (kefuService *KefuService) GetKefuInfoList(ctx context.Context, info shopReq.KefuSearch) (list []shop.Kefu, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.Kefu{})
	var kefus []shop.Kefu
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["id"] = true
	orderMap["created_at"] = true
	orderMap["name"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&kefus).Error
	return kefus, total, err
}
func (kefuService *KefuService) GetKefuPublic(ctx context.Context) (list []shop.KefuPublic, err error) {
	err = global.GVA_DB.WithContext(ctx).
		Model(&shop.Kefu{}).
		Select("name", "name_i18n", "avatar", "external_avatar", "qr_code", "contact_id", "status", "link").
		Where("status IN ?", []string{"在线", "忙碌"}).
		Find(&list).Error
	return
}
