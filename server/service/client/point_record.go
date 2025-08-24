package client

import (
	"context"
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"gorm.io/gorm"
)

type PointRecordService struct{}

// CreatePointRecord 创建积分记录管理记录
// Author [yourname](https://github.com/yourname)
func (cprService *PointRecordService) CreatePointRecord(ctx context.Context, cpr *client.PointRecord) (err error) {
	// 参数验证
	if cpr.UserId == nil || cpr.ChangeType == nil || cpr.PointChange == nil {
		return errors.New("用户ID、增减类型和积分变化不能为空")
	}

	// 根据增减类型调整积分变化值
	var actualPointChange int
	if *cpr.ChangeType == "increase" {
		actualPointChange = *cpr.PointChange
		if actualPointChange < 0 {
			actualPointChange = -actualPointChange // 确保增加时为正数
		}
	} else if *cpr.ChangeType == "decrease" {
		actualPointChange = -*cpr.PointChange
		if actualPointChange > 0 {
			actualPointChange = -actualPointChange // 确保减少时为负数
		}
	} else {
		return errors.New("增减类型必须是 increase 或 decrease")
	}

	// 从context中获取事务，如果没有则使用全局DB创建新事务
	var db *gorm.DB
	if txFromCtx := ctx.Value("tx"); txFromCtx != nil {
		if tx, ok := txFromCtx.(*gorm.DB); ok {
			db = tx
		} else {
			db = global.GVA_DB
		}
	} else {
		db = global.GVA_DB
	}

	// 如果从context获取到事务，直接使用；否则创建新事务
	if ctx.Value("tx") != nil {
		// 直接在现有事务中执行
		tx := db
		return cprService.executePointRecordLogic(tx, cpr, actualPointChange)
	} else {
		// 创建新事务
		return db.Transaction(func(tx *gorm.DB) error {
			return cprService.executePointRecordLogic(tx, cpr, actualPointChange)
		})
	}
}

// executePointRecordLogic 执行积分记录的核心逻辑
func (cprService *PointRecordService) executePointRecordLogic(tx *gorm.DB, cpr *client.PointRecord, actualPointChange int) error {
	// 1. 获取当前用户信息并锁定记录
	var user client.ClientUser
	if err := tx.Where("id = ?", *cpr.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 2. 计算新的积分总数
	newPoints := user.Point + actualPointChange
	if newPoints < 0 {
		return errors.New("积分不足，无法扣除")
	}

	// 3. 更新用户积分
	if err := tx.Model(&user).Update("point", newPoints).Error; err != nil {
		return err
	}

	// 4. 设置积分记录的当前积分和实际变化值
	cpr.CurrentPoints = &newPoints
	cpr.PointChange = &actualPointChange

	// 5. 创建积分记录
	if err := tx.Create(cpr).Error; err != nil {
		return err
	}

	return nil
}

// DeletePointRecord 删除积分记录管理记录
// Author [yourname](https://github.com/yourname)
func (cprService *PointRecordService) DeletePointRecord(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&client.PointRecord{}, "id = ?", ID).Error
	return err
}

// DeletePointRecordByIds 批量删除积分记录管理记录
// Author [yourname](https://github.com/yourname)
func (cprService *PointRecordService) DeletePointRecordByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]client.PointRecord{}, "id in ?", IDs).Error
	return err
}

// UpdatePointRecord 更新积分记录管理记录
// Author [yourname](https://github.com/yourname)
func (cprService *PointRecordService) UpdatePointRecord(ctx context.Context, cpr client.PointRecord) (err error) {
	err = global.GVA_DB.Model(&client.PointRecord{}).Where("id = ?", cpr.ID).Updates(&cpr).Error
	return err
}

// GetPointRecord 根据ID获取积分记录管理记录
// Author [yourname](https://github.com/yourname)
func (cprService *PointRecordService) GetPointRecord(ctx context.Context, ID string) (cpr client.PointRecord, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&cpr).Error
	return
}

// GetPointRecordInfoList 分页获取积分记录管理记录
// Author [yourname](https://github.com/yourname)
func (cprService *PointRecordService) GetPointRecordInfoList(ctx context.Context, info clientReq.PointRecordSearch) (list []client.PointRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&client.PointRecord{})
	var cprs []client.PointRecord
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.UserId != nil {
		db = db.Where("user_id = ?", *info.UserId)
	}
	if info.ChangeType != nil && *info.ChangeType != "" {
		db = db.Where("change_type = ?", *info.ChangeType)
	}
	if info.PointChange != nil {
		db = db.Where("point_change = ?", *info.PointChange)
	}
	if info.OperationType != nil && *info.OperationType != "" {
		db = db.Where("operation_type = ?", *info.OperationType)
	}
	if info.Reason != nil && *info.Reason != "" {
		db = db.Where("reason LIKE ?", "%"+*info.Reason+"%")
	}
	if info.CurrentPoints != nil {
		db = db.Where("current_points = ?", *info.CurrentPoints)
	}
	if info.RelatedOrderId != nil {
		db = db.Where("related_order_id = ?", *info.RelatedOrderId)
	}
	if info.Remark != nil && *info.Remark != "" {
		db = db.Where("remark LIKE ?", "%"+*info.Remark+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["id"] = true
	orderMap["created_at"] = true
	orderMap["point_change"] = true
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

	err = db.Find(&cprs).Error
	return cprs, total, err
}
func (cprService *PointRecordService) GetPointRecordPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
