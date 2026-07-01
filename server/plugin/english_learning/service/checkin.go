package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	englishReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"gorm.io/gorm"
)

type CheckinService struct{}

const (
	checkinOpReward          = "checkin_reward"
	checkinOpPointExchange   = "point_exchange"
	checkinReasonSignIn      = "reason_signInReward"
	checkinReasonPointDeduct = "reason_learningExchange"
	checkinReasonFreeTimeAdd = "reason_learningExchange"
)

// PerformCheckin 执行打卡逻辑
// 参数: userID, 基础积分basePoint(100), 每日递增数值increment(5), 轮回天数cycleDays(10)
func (s *CheckinService) PerformCheckin(userID uint, basePoint int, increment int, cycleDays int) (int, error) {
	var currentPointsAward int

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		// 截断到当天的零点
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		yesterday := today.AddDate(0, 0, -1)

		// 1. 判断今天是否已经打卡
		var count int64
		tx.Model(&model.CheckinRecord{}).Where("user_id = ? AND checkin_date = ?", userID, today).Count(&count)
		if count > 0 {
			return errors.New("今日已打卡，请勿重复操作")
		}

		// 2. 获取打卡统计
		var stats model.CheckinStats
		statsResult := tx.Where("user_id = ?", userID).Limit(1).Find(&stats)
		if statsResult.Error != nil {
			return statsResult.Error
		}
		if statsResult.RowsAffected == 0 {
			stats.UserID = userID
		}

		// 3. 判断昨天是否打卡，计算连续打卡天数
		var yestCount int64
		tx.Model(&model.CheckinRecord{}).Where("user_id = ? AND checkin_date = ?", userID, yesterday).Count(&yestCount)

		if yestCount > 0 {
			stats.ContinuousDays += 1
		} else {
			// 中断了，重新开始
			stats.ContinuousDays = 1
		}
		stats.TotalDays += 1

		// 4. 计算本次应发积分 (按照轮回周期计算)
		// 例: 第1天0加成(100)，第2天(105)...第10天(145)。第11天取模为0，落回(100)
		bonusDays := (stats.ContinuousDays - 1) % cycleDays
		currentPointsAward = basePoint + (bonusDays * increment)

		// 5. 写入打卡流水
		record := model.CheckinRecord{
			UserID:      userID,
			CheckinDate: today,
			PointAward:  currentPointsAward,
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		// 6. 更新打卡统计
		if stats.ID == 0 {
			if err := tx.Create(&stats).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Save(&stats).Error; err != nil {
				return err
			}
		}

		// 7. 发放积分到用户学习资产表 (同时兼容其他营销发放逻辑)
		asset, ensureErr := ServiceGroupApp.LearningAuthzService.EnsureUserAssetWithTx(tx, userID)
		if ensureErr != nil {
			return ensureErr
		}

		currentPoints := asset.TotalPoints + currentPointsAward
		if err := tx.Model(&asset).UpdateColumn("total_points", gorm.Expr("total_points + ?", currentPointsAward)).Error; err != nil {
			return err
		}

		return s.createPointRecord(tx, userID, "increase", currentPointsAward, checkinOpReward, checkinReasonSignIn, currentPoints, "")
	})

	return currentPointsAward, err
}

// ExchangePoints 积分兑换免费观看时长 (分钟)
func (s *CheckinService) ExchangePoints(userID uint, points int, rate int) (int, error) {
	if points%rate != 0 {
		return 0, errors.New("兑换积分必须是比率的整数倍")
	}

	exchangeMins := points / rate
	if exchangeMins <= 0 {
		return 0, errors.New("兑换时长必须大于0分钟")
	}

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var asset model.UserLearningAsset
		if err := tx.Where("user_id = ?", userID).First(&asset).Error; err != nil {
			return errors.New("用户资产异常")
		}

		if asset.TotalPoints < points {
			return errors.New("积分余额不足")
		}
		currentPoints := asset.TotalPoints - points
		currentFreeMinutes := asset.FreeMinutes + exchangeMins

		// 扣减积分，增加免费时长
		if err := tx.Model(&asset).Updates(map[string]interface{}{
			"total_points": gorm.Expr("total_points - ?", points),
			"free_minutes": gorm.Expr("free_minutes + ?", exchangeMins),
		}).Error; err != nil {
			return err
		}

		pointRemark := fmt.Sprintf("minutes:%d", exchangeMins)
		if err := s.createPointRecord(tx, userID, "decrease", -points, checkinOpPointExchange, checkinReasonPointDeduct, currentPoints, pointRemark); err != nil {
			return err
		}

		freeTimeRemark := fmt.Sprintf("points:%d", points)
		return s.createFreeTimeRecord(tx, userID, "increase", exchangeMins, checkinOpPointExchange, checkinReasonFreeTimeAdd, currentFreeMinutes, 0, freeTimeRemark)
	})

	return exchangeMins, err
}

func (s *CheckinService) GetStats(userID uint) (model.CheckinStats, bool, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var stats model.CheckinStats
	statsResult := global.GVA_DB.Where("user_id = ?", userID).Limit(1).Find(&stats)
	if statsResult.Error != nil {
		return model.CheckinStats{}, false, statsResult.Error
	}
	if statsResult.RowsAffected == 0 {
		stats = model.CheckinStats{UserID: userID}
	}

	var todayCount int64
	if err := global.GVA_DB.Model(&model.CheckinRecord{}).Where("user_id = ? AND checkin_date = ?", userID, today).Count(&todayCount).Error; err != nil {
		return stats, false, err
	}

	return stats, todayCount > 0, nil
}

func (s *CheckinService) GetPointRecordList(userID uint, info englishReq.PointRecordSearch) ([]model.EnglishPointRecord, int64, int, int, error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	limit := pageSize
	offset := pageSize * (page - 1)

	db := global.GVA_DB.Model(&model.EnglishPointRecord{}).Where("user_id = ?", userID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, page, pageSize, err
	}

	list := make([]model.EnglishPointRecord, 0, limit)
	if err := db.Order("created_at DESC").Order("id DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, page, pageSize, err
	}

	return list, total, page, pageSize, nil
}

func (s *CheckinService) GetCheckinRecordList(userID uint, info englishReq.CheckinRecordSearch) ([]model.CheckinRecord, int64, int, int, error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	limit := pageSize
	offset := pageSize * (page - 1)

	db := global.GVA_DB.Model(&model.CheckinRecord{}).Where("user_id = ?", userID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, page, pageSize, err
	}

	list := make([]model.CheckinRecord, 0, limit)
	if err := db.Order("checkin_date DESC").Order("id DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, page, pageSize, err
	}

	return list, total, page, pageSize, nil
}

func (s *CheckinService) createPointRecord(tx *gorm.DB, userID uint, changeType string, pointChange int, operationType string, reason string, currentPoints int, remark string) error {
	record := model.EnglishPointRecord{
		UserID:        userID,
		ChangeType:    changeType,
		PointChange:   pointChange,
		OperationType: operationType,
		Reason:        reason,
		CurrentPoints: currentPoints,
		Remark:        remark,
	}

	return tx.Create(&record).Error
}

func (s *CheckinService) createFreeTimeRecord(tx *gorm.DB, userID uint, changeType string, minuteChange int, operationType string, reason string, currentFreeMinutes int, relatedEpisodeID uint, remark string) error {
	record := model.EnglishFreeTimeRecord{
		UserID:             userID,
		ChangeType:         changeType,
		MinuteChange:       minuteChange,
		OperationType:      operationType,
		Reason:             reason,
		CurrentFreeMinutes: currentFreeMinutes,
		RelatedEpisodeID:   relatedEpisodeID,
		Remark:             remark,
	}

	return tx.Create(&record).Error
}
