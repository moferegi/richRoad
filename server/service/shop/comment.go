package shop

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"gorm.io/gorm"
)

type CommentService struct{}

// CreateComment 创建用户评论记录
// Author [piexlmax](https://github.com/piexlmax)
func (commentService *CommentService) CreateComment(comment *shop.Comment) (err error) {

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var orderDetail shop.OrderDetail
		dErr := global.GVA_DB.First(&orderDetail, "order_id = ? && good_id = ? && sku_id = ?", comment.OrderID, comment.GoodID, comment.SKUID).Error
		if dErr != nil {
			return errors.New("未找到订单详情")
		}
		e := global.GVA_DB.Model(&orderDetail).Update("is_comment", true).Error
		if e != nil {
			return errors.New("更新订单详情失败")
		}

		var order shop.Order
		ferr := global.GVA_DB.
			Preload("Detail").
			First(&order, "id = ? && user_id = ? && status = ?", comment.OrderID, comment.UserID, 3).
			Error
		if ferr != nil {
			return errors.New("未找到订单")
		}
		needUpdateStatus := true
		for _, detail := range order.Detail {
			if !detail.IsComment {
				needUpdateStatus = false
				break
			}
		}
		if needUpdateStatus {
			err = global.GVA_DB.Model(&order).Update("status", 7).Error
			if err != nil {
				return errors.New("更新订单状态失败")
			}
		}

		err = global.GVA_DB.Create(comment).Error
		return err
	})
	return err
}

// DeleteComment 删除用户评论记录
// Author [piexlmax](https://github.com/piexlmax)
func (commentService *CommentService) DeleteComment(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Comment{}, "id = ?", ID).Error
	return err
}

// DeleteCommentByIds 批量删除用户评论记录
// Author [piexlmax](https://github.com/piexlmax)
func (commentService *CommentService) DeleteCommentByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Comment{}, "id in ?", IDs).Error
	return err
}

// UpdateComment 更新用户评论记录
// Author [piexlmax](https://github.com/piexlmax)
func (commentService *CommentService) UpdateComment(comment shop.Comment) (err error) {
	err = global.GVA_DB.Model(&shop.Comment{}).Where("id = ?", comment.ID).Updates(&comment).Error
	return err
}

// GetComment 根据ID获取用户评论记录
// Author [piexlmax](https://github.com/piexlmax)
func (commentService *CommentService) GetComment(ID string) (comment []shop.Comment, err error) {
	err = global.GVA_DB.Preload("User").Where("good_id = ?", ID).Find(&comment).Error
	return
}

func (commentService *CommentService) GetCommentBk(ID string) (comment shop.Comment, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&comment).Error
	return
}

// GetCommentInfoList 分页获取用户评论记录
// Author [piexlmax](https://github.com/piexlmax)
func (commentService *CommentService) GetCommentInfoList(info shopReq.CommentSearch) (list []shop.Comment, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.Comment{})
	var comments []shop.Comment
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.OrderID != nil {
		db = db.Where("order_id = ?", info.OrderID)
	}
	if info.GoodID != nil {
		db = db.Where("good_id = ?", info.GoodID)
	}
	if info.SKUID != nil {
		db = db.Where("SKUID = ?", info.SKUID)
	}
	if info.StartShopReplyAt != nil && info.EndShopReplyAt != nil {
		db = db.Where("shop_reply_at BETWEEN ? AND ? ", info.StartShopReplyAt, info.EndShopReplyAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&comments).Error
	return comments, total, err
}
