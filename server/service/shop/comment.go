package shop

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentService struct{}

// CreateComment 创建用户评论记录
// Author [piexlmax](https://github.com/piexlmax)
func (commentService *CommentService) CreateComment(comment *shop.Comment) (err error) {

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order shop.Order
		ferr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&order, "id = ? AND user_id = ? AND status = ?", comment.OrderID, comment.UserID, 3).
			Error
		if ferr != nil {
			return errors.New("未找到订单")
		}

		var orderDetail shop.OrderDetail
		dErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&orderDetail, "order_id = ? AND good_id = ? AND sku_id = ?", comment.OrderID, comment.GoodID, comment.SKUID).Error
		if dErr != nil {
			return errors.New("未找到订单详情")
		}
		if orderDetail.IsComment {
			return errors.New("订单已评价")
		}
		result := tx.Model(&shop.OrderDetail{}).Where("id = ? AND is_comment = ?", orderDetail.ID, false).Update("is_comment", true)
		if result.Error != nil {
			return errors.New("更新订单详情失败")
		}
		if result.RowsAffected == 0 {
			return errors.New("订单已评价")
		}

		var remaining int64
		if err = tx.Model(&shop.OrderDetail{}).Where("order_id = ? AND is_comment = ?", comment.OrderID, false).Count(&remaining).Error; err != nil {
			return errors.New("查询订单评价状态失败")
		}
		if remaining == 0 {
			err = tx.Model(&order).Update("status", 7).Error
			if err != nil {
				return errors.New("更新订单状态失败")
			}
		}

		err = tx.Create(comment).Error
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

// GetComment 根据商品ID获取公开评论记录
// Author [piexlmax](https://github.com/piexlmax)
func (commentService *CommentService) GetComment(ID string) (comments []shop.CommentPublic, err error) {
	var count int64
	err = global.GVA_DB.Model(&shop.Good{}).
		Where("id = ? AND status = ?", ID, true).
		Where("category_id IS NULL OR category_id NOT IN (SELECT id FROM shop_category WHERE show_in_uni = ?)", false).
		Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return []shop.CommentPublic{}, nil
	}

	var records []shop.Comment
	err = global.GVA_DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "avatar", "nickname")
	}).Where("good_id = ?", ID).Find(&records).Error
	if err != nil {
		return nil, err
	}

	comments = make([]shop.CommentPublic, 0, len(records))
	for _, record := range records {
		item := shop.CommentPublic{
			ID:          record.ID,
			CreatedAt:   record.CreatedAt,
			GoodID:      record.GoodID,
			SKUID:       record.SKUID,
			Pics:        record.Pics,
			Rating:      record.Rating,
			Content:     record.Content,
			ShopReply:   record.ShopReply,
			ShopReplyAt: record.ShopReplyAt,
		}
		if record.User != nil {
			item.User = shop.CommentPublicUser{
				ID:       record.User.ID,
				Avatar:   record.User.Avatar,
				Nickname: record.User.Nickname,
			}
		}
		comments = append(comments, item)
	}
	return comments, nil
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
