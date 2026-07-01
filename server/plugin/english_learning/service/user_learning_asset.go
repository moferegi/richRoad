package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	englishReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserLearningAssetService struct{}

const (
	heartbeatMaxUsingSecs      = 120
	heartbeatProgressTolerance = 8.0
	heartbeatCounterTTL        = 24 * time.Hour
	heartbeatMetaTTL           = 2 * time.Hour
	freeTimeOpHeartbeatConsume = "heartbeat_consume"
	freeTimeReasonWatchConsume = "reason_learningWatchConsume"
)

// HandleHeartbeat 处理前端定期上报的心跳 (防薅羊毛、更新观看进度和扣除免费时长)
func (s *UserLearningAssetService) HandleHeartbeat(userID uint, episodeID uint, progressSecs float64, usingSecs int) error {
	if userID == 0 || episodeID == 0 {
		return errors.New("心跳参数异常")
	}

	progressSecs = normalizeHeartbeatProgress(progressSecs)
	if usingSecs < 0 {
		usingSecs = 0
	}
	if usingSecs > heartbeatMaxUsingSecs {
		usingSecs = heartbeatMaxUsingSecs
	}

	ctx := context.Background()

	// 1. 查询并更新视频观看历史
	var history model.UserWatchHistory
	historyResult := global.GVA_DB.Where("user_id = ? AND episode_id = ?", userID, episodeID).Limit(1).Find(&history)
	if historyResult.Error != nil {
		return historyResult.Error
	}
	historyNotFound := historyResult.RowsAffected == 0

	prevProgress := 0.0
	if !historyNotFound {
		prevProgress = normalizeHeartbeatProgress(history.ProgressSecs)
	}

	trustedUsingSecs := trustHeartbeatUsingSeconds(usingSecs, prevProgress, progressSecs)
	trustedUsingSecs, guardErr := applyHeartbeatRateGuard(ctx, userID, episodeID, trustedUsingSecs)
	if guardErr != nil {
		return guardErr
	}

	storedProgress := progressSecs
	if storedProgress < prevProgress {
		// 防止恶意回退进度影响后续可信秒数判断
		storedProgress = prevProgress
	}

	if historyNotFound {
		history = model.UserWatchHistory{UserID: userID, EpisodeID: episodeID, ProgressSecs: storedProgress}
		if createErr := global.GVA_DB.Create(&history).Error; createErr != nil {
			return createErr
		}
	} else {
		if updateErr := global.GVA_DB.Model(&history).Update("progress_secs", storedProgress).Error; updateErr != nil {
			return updateErr
		}
	}

	if trustedUsingSecs <= 0 {
		return nil
	}

	// 2. 检查用户的免费时长资产
	var asset model.UserLearningAsset
	assetResult := global.GVA_DB.Where("user_id = ?", userID).Limit(1).Find(&asset)
	if assetResult.Error != nil {
		return assetResult.Error
	}
	if assetResult.RowsAffected == 0 {
		ensuredAsset, ensureErr := ServiceGroupApp.LearningAuthzService.EnsureUserAsset(userID)
		if ensureErr != nil {
			return ensureErr
		}
		asset = ensuredAsset
	}

	// 3. 时长扣除逻辑
	now := time.Now()
	// 如果仍在绝对全场免费阶段(如注册赠送24小时内) 或 者 VIP，不扣除 FreeMinutes
	hasFreeTime := asset.FreeTimeExpire != nil && asset.FreeTimeExpire.After(now)
	isVip := asset.IsVip != nil && *asset.IsVip

	if hasFreeTime || isVip {
		return nil // 无需扣除桶内余额
	}
	if asset.FreeMinutes <= 0 {
		return nil
	}

	// 利用 Redis 作为计数缓冲桶，满 60 秒才去扣去数据库的 1 freeMinute
	redisKey := fmt.Sprintf("english_learning:heartbeat_secs:%d", userID)
	if global.GVA_REDIS == nil {
		deductMins := trustedUsingSecs / 60
		deductMins = reflectDeduct(asset.FreeMinutes, deductMins)
		if deductMins <= 0 {
			return nil
		}

		currentFreeMinutes := asset.FreeMinutes - deductMins
		remark := fmt.Sprintf("episode:%d,usingSecs:%d", episodeID, trustedUsingSecs)
		return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&model.UserLearningAsset{}).
				Where("user_id = ? AND free_minutes > 0", userID).
				UpdateColumn("free_minutes", gorm.Expr("free_minutes - ?", deductMins)).Error; err != nil {
				return err
			}

			return s.createFreeTimeRecord(tx, userID, "decrease", -deductMins, freeTimeOpHeartbeatConsume, freeTimeReasonWatchConsume, currentFreeMinutes, episodeID, remark)
		})
	}

	// 在原有的累积心跳秒数上增加
	accumulatedSecs, err := global.GVA_REDIS.IncrBy(ctx, redisKey, int64(trustedUsingSecs)).Result()
	if err != nil {
		return err
	}
	_ = global.GVA_REDIS.Expire(ctx, redisKey, heartbeatCounterTTL)

	// 如果累积超过60秒，换算为扣除的分钟数
	if accumulatedSecs >= 60 {
		deductMins := int(accumulatedSecs / 60)
		remainSecs := int(accumulatedSecs % 60)
		deductMins = reflectDeduct(asset.FreeMinutes, deductMins)

		if deductMins > 0 {
			currentFreeMinutes := asset.FreeMinutes - deductMins
			remark := fmt.Sprintf("episode:%d,usingSecs:%d", episodeID, trustedUsingSecs)
			if updateErr := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
				if err := tx.Model(&model.UserLearningAsset{}).
					Where("user_id = ? AND free_minutes > 0", userID).
					UpdateColumn("free_minutes", gorm.Expr("free_minutes - ?", deductMins)).Error; err != nil {
					return err
				}

				return s.createFreeTimeRecord(tx, userID, "decrease", -deductMins, freeTimeOpHeartbeatConsume, freeTimeReasonWatchConsume, currentFreeMinutes, episodeID, remark)
			}); updateErr != nil {
				return updateErr
			}
		}

		// 重置 Redis 中的余数
		_ = global.GVA_REDIS.Set(ctx, redisKey, remainSecs, heartbeatCounterTTL)
	}

	return nil
}

func (s *UserLearningAssetService) GetFreeTimeRecordList(userID uint, info englishReq.FreeTimeRecordSearch) ([]model.EnglishFreeTimeRecord, int64, int, int, error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	limit := pageSize
	offset := pageSize * (page - 1)

	db := global.GVA_DB.Model(&model.EnglishFreeTimeRecord{}).Where("user_id = ?", userID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, page, pageSize, err
	}

	list := make([]model.EnglishFreeTimeRecord, 0, limit)
	if err := db.Order("created_at DESC").Order("id DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, page, pageSize, err
	}

	return list, total, page, pageSize, nil
}

func (s *UserLearningAssetService) createFreeTimeRecord(tx *gorm.DB, userID uint, changeType string, minuteChange int, operationType string, reason string, currentFreeMinutes int, relatedEpisodeID uint, remark string) error {
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

func normalizeHeartbeatProgress(progressSecs float64) float64 {
	if math.IsNaN(progressSecs) || math.IsInf(progressSecs, 0) || progressSecs < 0 {
		return 0
	}
	return progressSecs
}

func trustHeartbeatUsingSeconds(rawUsing int, previousProgress float64, currentProgress float64) int {
	if rawUsing <= 0 {
		return 0
	}

	forwardDelta := currentProgress - previousProgress
	if forwardDelta < 0 {
		forwardDelta = 0
	}

	maxByProgress := int(math.Ceil(forwardDelta + heartbeatProgressTolerance))
	trusted := rawUsing

	if maxByProgress == 0 && trusted > 5 {
		return 0
	}
	if maxByProgress > 0 && trusted > maxByProgress {
		trusted = maxByProgress
	}
	if trusted < 0 {
		return 0
	}
	if trusted > heartbeatMaxUsingSecs {
		trusted = heartbeatMaxUsingSecs
	}
	return trusted
}

func applyHeartbeatRateGuard(ctx context.Context, userID uint, episodeID uint, usingSecs int) (int, error) {
	if usingSecs <= 0 || global.GVA_REDIS == nil {
		return usingSecs, nil
	}

	lastTsKey := fmt.Sprintf("english_learning:heartbeat_last_ts:%d:%d", userID, episodeID)
	now := time.Now().Unix()

	lastTsRaw, err := global.GVA_REDIS.Get(ctx, lastTsKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return usingSecs, err
	}

	if strings.TrimSpace(lastTsRaw) != "" {
		lastTs, parseErr := strconv.ParseInt(lastTsRaw, 10, 64)
		if parseErr == nil {
			interval := now - lastTs
			if interval <= 0 {
				if usingSecs > 2 {
					usingSecs = 2
				}
			} else {
				limit := int(interval) + 2
				if limit < 0 {
					limit = 0
				}
				if usingSecs > limit {
					usingSecs = limit
				}
			}
		}
	}

	if setErr := global.GVA_REDIS.Set(ctx, lastTsKey, now, heartbeatMetaTTL).Err(); setErr != nil {
		global.GVA_LOG.Warn("写入心跳频率保护时间戳失败", zap.Error(setErr), zap.Uint("userID", userID), zap.Uint("episodeID", episodeID))
	}

	if usingSecs < 0 {
		usingSecs = 0
	}
	return usingSecs, nil
}

// 辅助方法：确保扣除后不会变成负数
func reflectDeduct(current, deduct int) int {
	if deduct < 0 {
		return 0
	}
	if current < deduct {
		return current
	}
	return deduct
}
