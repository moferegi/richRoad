package shop

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system" // 假设用户模型
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CouponOrderUserService struct{}

// 雪花算法节点，用于生成唯一ID
var (
	snowflakeNode *snowflake.Node
	nodeOnce      sync.Once
	nodeErr       error
)

// generateSnowflakeID 生成雪花算法ID
func generateSnowflakeID() (string, error) {
	nodeOnce.Do(func() {
		// 创建一个节点，节点ID为1
		snowflakeNode, nodeErr = snowflake.NewNode(1)
	})

	if nodeErr != nil {
		return "", nodeErr
	}

	// 生成ID
	id := snowflakeNode.Generate()
	return id.String(), nil
}

// CreateCouponOrderUser 创建优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) CreateCouponOrderUser(ctx context.Context, cou *shop.CouponOrderUser) (err error) {
	// 自动生成优惠券编号
	if cou.CouponNum == "" {
		cou.CouponNum = fmt.Sprintf("CPN%d%04d", time.Now().UnixMilli(), cou.UserID%10000)
	}
	// 默认未使用
	if cou.Status == nil {
		unused := false
		cou.Status = &unused
	}
	// 默认领取时间
	if cou.ClaimedAt == nil {
		now := time.Now()
		cou.ClaimedAt = &now
	}
	// 同步 ShopUserID（查询用此字段）
	if cou.ShopUserID == nil && cou.UserID > 0 {
		uid := int(cou.UserID)
		cou.ShopUserID = &uid
	}
	err = global.GVA_DB.Create(cou).Error
	return err
}

// DeleteCouponOrderUser 删除优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) DeleteCouponOrderUser(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.CouponOrderUser{}, "id = ?", ID).Error
	return err
}

// DeleteCouponOrderUserByIds 批量删除优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) DeleteCouponOrderUserByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.CouponOrderUser{}, "id in ?", IDs).Error
	return err
}

// UpdateCouponOrderUser 更新优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) UpdateCouponOrderUser(ctx context.Context, cou shop.CouponOrderUser) (err error) {
	err = global.GVA_DB.Model(&shop.CouponOrderUser{}).Where("id = ?", cou.ID).Updates(&cou).Error
	return err
}

// GetCouponOrderUser 根据ID获取优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) GetCouponOrderUser(ctx context.Context, ID string) (cou shop.CouponOrderUser, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&cou).Error
	return
}

// GetCouponOrderUserInfoList 分页获取优惠券记录
// Author [yourname](https://github.com/yourname)
func (couService *CouponOrderUserService) GetCouponOrderUserInfoList(ctx context.Context, info shopReq.CouponOrderUserSearch) (list []shop.CouponOrderUser, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&shop.CouponOrderUser{})
	var cous []shop.CouponOrderUser
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	if info.CouponID != nil {
		db = db.Where("coupon_id = ?", *info.CouponID)
	}
	if info.OrderID != nil {
		db = db.Where("order_id = ?", *info.OrderID)
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	if info.UserID != nil {
		db = db.Where("user_id = ?", *info.UserID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["id"] = true
	orderMap["created_at"] = true
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

	err = db.Find(&cous).Error
	return cous, total, err
}
func (couService *CouponOrderUserService) GetCouponOrderUserDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	couponID := make([]map[string]any, 0)

	global.GVA_DB.Table("shop_coupon").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&couponID)
	res["couponID"] = couponID
	orderID := make([]map[string]any, 0)

	global.GVA_DB.Table("shop_order").Where("deleted_at IS NULL").Select("out_trade_no as label,id as value").Scan(&orderID)
	res["orderID"] = orderID
	return
}
func (couService *CouponOrderUserService) GetCouponOrderUserPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

func (couService *CouponOrderUserService) GetAllClaimCoupon(ctx context.Context, userID uint, goodIds []int) (coupons []map[string]interface{}, err error) {
	// 初始化返回结果切片
	coupons = make([]map[string]interface{}, 0)
	shopUserID := int(userID)

	// 1. 查询所有启用的优惠券（包括已过期的，后面按是否已领分别处理）
	var allCoupons []shop.Coupon
	if err = global.GVA_DB.Where("status = ?", true).Find(&allCoupons).Error; err != nil {
		return nil, errors.New("查询优惠券失败: " + err.Error())
	}

	// 2. 一次性查询用户所有已领取的券记录
	var allUserCoupons []shop.CouponOrderUser
	global.GVA_DB.Where("shop_user_id = ?", shopUserID).Find(&allUserCoupons)
	// 按 coupon_id 分组
	userCouponMap := make(map[uint][]shop.CouponOrderUser)
	for _, uc := range allUserCoupons {
		if uc.CouponID != nil {
			userCouponMap[uint(*uc.CouponID)] = append(userCouponMap[uint(*uc.CouponID)], uc)
		}
	}

	now := time.Now()
	for _, coupon := range allCoupons {
		inPeriod := coupon.StartTime != nil && coupon.EndTime != nil &&
			!now.Before(*coupon.StartTime) && !now.After(*coupon.EndTime)

		userCoupons := userCouponMap[coupon.ID]

		// 如果券不在有效期内且用户没有已领取的未使用券，跳过
		hasUnused := false
		for _, uc := range userCoupons {
			if uc.OrderID == nil {
				hasUnused = true
				break
			}
		}
		if !inPeriod && !hasUnused {
			continue
		}

		// 构建基础券数据
		couponData := map[string]interface{}{
			"couponNum":       0,
			"couponID":        coupon.ID,
			"name":            "",
			"minSpend":        coupon.MinSpend,
			"discount":        coupon.Discount,
			"productID":       0,
			"startTime":       "",
			"endTime":         "",
			"status":          0,
			"canUse":          0,
			"nameI18n":        coupon.NameI18n,
			"descriptionI18n": coupon.DescriptionI18n,
			"backgroundImage": coupon.BackgroundImage,
			"externalBgPath":  coupon.ExternalBgPath,
		}
		if coupon.Name != nil {
			couponData["name"] = *coupon.Name
		}
		if coupon.StartTime != nil {
			couponData["startTime"] = coupon.StartTime.Format("2006-01-02")
		}
		if coupon.EndTime != nil {
			couponData["endTime"] = coupon.EndTime.Format("2006-01-02")
		}
		if coupon.ProductID != nil {
			couponData["productID"] = *coupon.ProductID
		}
		if coupon.Description != nil {
			couponData["description"] = *coupon.Description
		}

		// 判断该券对当前商品是否可用
		if len(goodIds) > 0 {
			allowedIDs := couponAllowedProductIDs(coupon)
			if len(allowedIDs) == 0 {
				couponData["canUse"] = 1
			} else {
				for _, goodID := range goodIds {
					if allowedIDs[goodID] {
						couponData["canUse"] = 1
						break
					}
				}
			}
		} else {
			couponData["canUse"] = 1
		}

		if len(userCoupons) == 0 {
			// 未领取过，仅在有效期内才显示可领取
			if inPeriod {
				coupons = append(coupons, couponData)
			}
		} else {
			// 为每张已领取的券生成一条记录
			for _, uc := range userCoupons {
				item := make(map[string]interface{})
				for k, v := range couponData {
					item[k] = v
				}
				item["couponNum"] = uc.CouponNum
				if uc.OrderID != nil {
					item["status"] = 1
					item["used"] = true
				} else if !inPeriod {
					// 已领取但优惠券已过期，标记为过期不可用
					item["status"] = 2
					item["used"] = false
					item["expired"] = true
					item["canUse"] = 0
				} else {
					item["status"] = 0
					item["used"] = false
				}
				coupons = append(coupons, item)
			}
		}
	}

	return coupons, nil
}

// ClaimCouponByUser 用户领取优惠券
func (couService *CouponOrderUserService) ClaimCouponByUser(ctx context.Context, userID uint, couponID int) (couponNum string, err error) {
	shopUserID := int(userID)
	now := time.Now()

	err = global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var coupon shop.Coupon
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND status = ? AND start_time <= ? AND end_time >= ?", couponID, true, now, now).
			First(&coupon).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("优惠券无效或已领完")
			}
			return errors.New("查询优惠券失败: " + err.Error())
		}

		claimed := 0
		if coupon.Claimed != nil {
			claimed = *coupon.Claimed
		}
		if coupon.Quantity != nil && claimed >= *coupon.Quantity {
			return errors.New("优惠券无效或已领完")
		}

		var existingClaim shop.CouponOrderUser
		err := tx.Where("coupon_id = ? AND shop_user_id = ?", couponID, shopUserID).First(&existingClaim).Error
		if err == nil {
			return errors.New("您已领取过该优惠券")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("检查领取状态失败: " + err.Error())
		}

		snowflakeID, err := generateSnowflakeID()
		if err != nil {
			return errors.New("生成券码失败: " + err.Error())
		}
		couponNum = snowflakeID

		status := false
		couponOrderUser := shop.CouponOrderUser{
			CouponNum:  snowflakeID,
			CouponID:   &couponID,
			UserID:     userID,
			ShopUserID: &shopUserID,
			Status:     &status,
			ClaimedAt:  &now,
		}
		if err = tx.Create(&couponOrderUser).Error; err != nil {
			return errors.New("领取优惠券失败: " + err.Error())
		}

		result := tx.Model(&shop.Coupon{}).Where("id = ?", couponID).UpdateColumn("claimed", gorm.Expr("COALESCE(claimed, 0) + ?", 1))
		if result.Error != nil {
			return errors.New("更新优惠券数量失败: " + result.Error.Error())
		}
		if result.RowsAffected == 0 {
			return errors.New("优惠券已领完或更新失败")
		}

		return nil
	})
	return couponNum, err
}

// IssueCouponToAllUsers 管理员向所有用户发放优惠券
// 注意：此函数可能需要根据实际用户表结构进行调整，这里假设用户表为 system.SysUser
func (couService *CouponOrderUserService) IssueCouponToAllUsers(ctx context.Context, couponID int) (err error) {
	// 1. 获取所有用户ID (假设用户表为 system_users, 字段为 id)
	// 您需要根据您的实际用户表结构修改此部分
	var userIDs []uint
	// 假设 shop_user 表存在并且有 id 字段
	if err = global.GVA_DB.WithContext(ctx).Table("shop_user").Pluck("id", &userIDs).Error; err != nil {
		// 如果 shop_user 不存在，尝试 system.SysUser
		// 注意：gin-vue-admin 的用户表通常是 sys_users, 模型是 system.SysUser
		// 如果您的用户体系不同，请修改这里的查询
		if err = global.GVA_DB.WithContext(ctx).Model(&system.SysUser{}).Pluck("id", &userIDs).Error; err != nil {
			return errors.New("获取用户列表失败: " + err.Error())
		}
	}

	if len(userIDs) == 0 {
		return errors.New("没有找到任何用户")
	}

	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var coupon shop.Coupon
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND status = ? AND start_time <= NOW() AND end_time >= NOW()", couponID, true).
			First(&coupon).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("优惠券无效")
			}
			return errors.New("查询优惠券失败: " + err.Error())
		}

		claimed := 0
		if coupon.Claimed != nil {
			claimed = *coupon.Claimed
		}
		remaining := -1
		if coupon.Quantity != nil {
			remaining = *coupon.Quantity - claimed
			if remaining <= 0 {
				return errors.New("优惠券无效或已领完")
			}
		}

		issuedCount := 0
		status := false // false 表示未使用
		now := time.Now()

		for _, userID := range userIDs {
			if remaining >= 0 && issuedCount >= remaining {
				global.GVA_LOG.Warn("优惠券库存不足，停止发放", zap.Int("couponID", couponID))
				break
			}

			var existingClaim shop.CouponOrderUser
			shopUserID := int(userID)
			findErr := tx.Where("coupon_id = ? AND shop_user_id = ?", couponID, shopUserID).First(&existingClaim).Error
			if findErr == nil {
				continue
			}
			if !errors.Is(findErr, gorm.ErrRecordNotFound) {
				return errors.New("检查用户领取状态失败: " + findErr.Error())
			}

			snowflakeID, err := generateSnowflakeID()
			if err != nil {
				return errors.New("生成券码失败: " + err.Error())
			}

			couponOrderUser := shop.CouponOrderUser{
				CouponNum:  snowflakeID,
				CouponID:   &couponID,
				UserID:     userID,
				ShopUserID: &shopUserID,
				Status:     &status,
				ClaimedAt:  &now,
			}
			if err = tx.Create(&couponOrderUser).Error; err != nil {
				return errors.New("为用户发放优惠券失败: " + err.Error())
			}
			issuedCount++
		}

		if issuedCount == 0 {
			return nil
		}

		result := tx.Model(&shop.Coupon{}).
			Where("id = ? AND (quantity IS NULL OR COALESCE(claimed, 0) + ? <= quantity)", couponID, issuedCount).
			UpdateColumn("claimed", gorm.Expr("COALESCE(claimed, 0) + ?", issuedCount))
		if result.Error != nil {
			return errors.New("批量更新优惠券数量失败: " + result.Error.Error())
		}
		if result.RowsAffected == 0 {
			return errors.New("优惠券库存不足，发放失败")
		}

		return nil
	})
}

// couponAllowedProductIDs 解析优惠券允许的商品ID集合
// 优先使用 ProductIDs（逗号分隔），兼容旧字段 ProductID（单个）
func couponAllowedProductIDs(coupon shop.Coupon) map[int]bool {
	ids := make(map[int]bool)
	if coupon.ProductIDs != "" {
		parts := strings.Split(coupon.ProductIDs, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if id, err := strconv.Atoi(p); err == nil && id > 0 {
				ids[id] = true
			}
		}
	}
	if len(ids) == 0 && coupon.ProductID != nil && *coupon.ProductID > 0 {
		ids[*coupon.ProductID] = true
	}
	return ids
}
