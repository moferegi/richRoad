package client

import (
	"context"
	"errors"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PointRecordService struct{}

type TryonPointModelStatsItem struct {
	ModelKey            string `json:"modelKey"`
	TaskCount           int64  `json:"taskCount"`
	RefinerEnabledCount int64  `json:"refinerEnabledCount"`
	BeautifyUsedCount   int64  `json:"beautifyUsedCount"`
	TotalCostPoints     int64  `json:"totalCostPoints"`
	BeautifyCostPoints  int64  `json:"beautifyCostPoints"`
}

type TryonPointStatsData struct {
	RegisterRewardTotal int64                      `json:"registerRewardTotal"`
	InviteRewardTotal   int64                      `json:"inviteRewardTotal"`
	RechargeTotal       int64                      `json:"rechargeTotal"`
	AdminIncreaseTotal  int64                      `json:"adminIncreaseTotal"`
	AdminDecreaseTotal  int64                      `json:"adminDecreaseTotal"`
	TotalGranted        int64                      `json:"totalGranted"`
	TotalUsed           int64                      `json:"totalUsed"`
	BeautifyUsedTotal   int64                      `json:"beautifyUsedTotal"`
	ModelCallTotal      int64                      `json:"modelCallTotal"`
	ModelCostTotal      int64                      `json:"modelCostTotal"`
	BeautifyCostTotal   int64                      `json:"beautifyCostTotal"`
	ModelStats          []TryonPointModelStatsItem `json:"modelStats"`
}

func normalizeAssetType(assetType *string) (string, error) {
	if assetType == nil || strings.TrimSpace(*assetType) == "" {
		return client.AssetTypePoint, nil
	}

	value := strings.TrimSpace(*assetType)
	switch value {
	case client.AssetTypePoint, client.AssetTypeTryonPoint:
		return value, nil
	default:
		return "", errors.New("assetTypeMustPointOrTryonPoint")
	}
}

// CreatePointRecord 创建积分记录管理记录
// Author [yourname](https://github.com/yourname)
func (cprService *PointRecordService) CreatePointRecord(ctx context.Context, cpr *client.PointRecord) (err error) {
	// 参数验证
	if cpr.UserId == nil || cpr.ChangeType == nil || cpr.PointChange == nil {
		return errors.New("pointRecordRequiredFields")
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
		return errors.New("pointChangeTypeInvalid")
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
	assetType, err := normalizeAssetType(cpr.AssetType)
	if err != nil {
		return err
	}
	cpr.AssetType = &assetType

	// 1. 获取当前用户信息并锁定记录
	var user client.ClientUser
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", *cpr.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("userNotExist")
		}
		return err
	}

	// 2. 计算新的资产总数
	currentBalance := 0
	fieldName := "point"
	insufficientMessage := "pointInsufficient"
	switch assetType {
	case client.AssetTypePoint:
		currentBalance = user.Point
		fieldName = "point"
		insufficientMessage = "pointInsufficient"
	case client.AssetTypeTryonPoint:
		currentBalance = user.TryonPoint
		fieldName = "tryon_point"
		insufficientMessage = "tryonPointInsufficient"
	}

	newPoints := currentBalance + actualPointChange
	if newPoints < 0 {
		return errors.New(insufficientMessage)
	}

	// 3. 更新用户资产
	if err := tx.Model(&user).Update(fieldName, newPoints).Error; err != nil {
		return err
	}

	// 4. 设置积分记录的当前资产和实际变化值
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
	if info.AssetType != nil && strings.TrimSpace(*info.AssetType) != "" {
		assetType, normalizeErr := normalizeAssetType(info.AssetType)
		if normalizeErr != nil {
			return nil, 0, normalizeErr
		}
		db = db.Where("asset_type = ?", assetType)
	} else {
		// 兼容历史数据：旧记录 asset_type 为空，默认归为积分记录。
		db = db.Where("(asset_type = ? OR asset_type IS NULL)", client.AssetTypePoint)
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
	} else {
		db = db.Order("created_at desc")
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&cprs).Error
	return cprs, total, err
}

func (cprService *PointRecordService) GetTryonPointStats(ctx context.Context, info clientReq.TryonPointStatsSearch) (stats TryonPointStatsData, err error) {
	_ = ctx

	const (
		tryonOpRecharge  = "tryon_recharge"
		tryonOpBeautify  = "tryon_beautify_consume"
		adminAdjustOpKey = "admin_adjust_tryon_point"
	)

	sumByChangeType := func(operationType string, changeType string) (int64, error) {
		var total int64
		err := cprService.applyTryonPointStatsFilters(global.GVA_DB.Model(&client.PointRecord{}), info).
			Where("operation_type = ? AND change_type = ?", strings.TrimSpace(operationType), strings.TrimSpace(changeType)).
			Select("COALESCE(SUM(ABS(point_change)), 0)").
			Scan(&total).Error
		return total, err
	}

	if stats.RegisterRewardTotal, err = sumByChangeType(tryonOpRegisterReward, "increase"); err != nil {
		return stats, err
	}
	if stats.InviteRewardTotal, err = sumByChangeType(tryonOpInviteReward, "increase"); err != nil {
		return stats, err
	}
	if stats.RechargeTotal, err = sumByChangeType(tryonOpRecharge, "increase"); err != nil {
		return stats, err
	}
	if stats.AdminIncreaseTotal, err = sumByChangeType(adminAdjustOpKey, "increase"); err != nil {
		return stats, err
	}
	if stats.AdminDecreaseTotal, err = sumByChangeType(adminAdjustOpKey, "decrease"); err != nil {
		return stats, err
	}
	if stats.TotalUsed, err = sumByChangeType(tryonOpConsume, "decrease"); err != nil {
		return stats, err
	}
	beautifyUsedCost, sumErr := sumByChangeType(tryonOpBeautify, "decrease")
	if sumErr != nil {
		return stats, sumErr
	}
	stats.TotalUsed += beautifyUsedCost

	if err = cprService.applyTryonPointStatsFilters(global.GVA_DB.Model(&client.PointRecord{}), info).
		Where("point_change > 0").
		Select("COALESCE(SUM(point_change), 0)").
		Scan(&stats.TotalGranted).Error; err != nil {
		return stats, err
	}

	taskQuery := global.GVA_DB.Model(&client.TryonTask{})
	if len(info.CreatedAtRange) == 2 {
		taskQuery = taskQuery.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
	if info.UserId != nil {
		taskQuery = taskQuery.Where("user_id = ?", *info.UserId)
	}

	type modelAggRow struct {
		Provider            string `json:"provider"`
		TaskCount           int64  `json:"taskCount"`
		RefinerEnabledCount int64  `json:"refinerEnabledCount"`
		BeautifyUsedCount   int64  `json:"beautifyUsedCount"`
		TotalCostPoints     int64  `json:"totalCostPoints"`
		BeautifyCostPoints  int64  `json:"beautifyCostPoints"`
	}
	var rows []modelAggRow
	err = taskQuery.
		Select("provider, COUNT(*) as task_count, COALESCE(SUM(CASE WHEN enable_refiner THEN 1 ELSE 0 END), 0) as refiner_enabled_count, COALESCE(SUM(CASE WHEN beautify_status <> '' AND beautify_status <> 'disabled' THEN 1 ELSE 0 END), 0) as beautify_used_count, COALESCE(SUM(cost_points), 0) as total_cost_points, COALESCE(SUM(beautify_cost), 0) as beautify_cost_points").
		Group("provider").
		Order("task_count DESC").
		Scan(&rows).Error
	if err != nil {
		return stats, err
	}

	stats.ModelStats = make([]TryonPointModelStatsItem, 0, len(rows))
	for _, row := range rows {
		modelKey := strings.TrimSpace(row.Provider)
		if modelKey == "" {
			modelKey = "default"
		}
		stats.ModelStats = append(stats.ModelStats, TryonPointModelStatsItem{
			ModelKey:            modelKey,
			TaskCount:           row.TaskCount,
			RefinerEnabledCount: row.RefinerEnabledCount,
			BeautifyUsedCount:   row.BeautifyUsedCount,
			TotalCostPoints:     row.TotalCostPoints,
			BeautifyCostPoints:  row.BeautifyCostPoints,
		})
		stats.ModelCallTotal += row.TaskCount
		stats.ModelCostTotal += row.TotalCostPoints
		stats.BeautifyUsedTotal += row.BeautifyUsedCount
		stats.BeautifyCostTotal += row.BeautifyCostPoints
	}

	return stats, nil
}

func (cprService *PointRecordService) applyTryonPointStatsFilters(db *gorm.DB, info clientReq.TryonPointStatsSearch) *gorm.DB {
	db = db.Where("asset_type = ?", client.AssetTypeTryonPoint)
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
	if info.UserId != nil {
		db = db.Where("user_id = ?", *info.UserId)
	}
	return db
}

func (cprService *PointRecordService) GetPointRecordPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
