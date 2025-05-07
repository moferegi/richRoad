
package shop

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
    shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type TagService struct {}
// CreateTag 创建标签记录
// Author [yourname](https://github.com/yourname)
func (tagService *TagService) CreateTag(ctx context.Context, tag *shop.Tag) (err error) {
	err = global.GVA_DB.Create(tag).Error
	return err
}

// DeleteTag 删除标签记录
// Author [yourname](https://github.com/yourname)
func (tagService *TagService)DeleteTag(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Tag{},"id = ?",ID).Error
	return err
}

// DeleteTagByIds 批量删除标签记录
// Author [yourname](https://github.com/yourname)
func (tagService *TagService)DeleteTagByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Tag{},"id in ?",IDs).Error
	return err
}

// UpdateTag 更新标签记录
// Author [yourname](https://github.com/yourname)
func (tagService *TagService)UpdateTag(ctx context.Context, tag shop.Tag) (err error) {
	err = global.GVA_DB.Model(&shop.Tag{}).Where("id = ?",tag.ID).Updates(&tag).Error
	return err
}

// GetTag 根据ID获取标签记录
// Author [yourname](https://github.com/yourname)
func (tagService *TagService)GetTag(ctx context.Context, ID string) (tag shop.Tag, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&tag).Error
	return
}
// GetTagInfoList 分页获取标签记录
// Author [yourname](https://github.com/yourname)
func (tagService *TagService)GetTagInfoList(ctx context.Context, info shopReq.TagSearch) (list []shop.Tag, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&shop.Tag{})
    var tags []shop.Tag
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
           orderMap["ID"] = true
           orderMap["CreatedAt"] = true
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

	err = db.Find(&tags).Error
	return  tags, total, err
}
func (tagService *TagService)GetTagPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
