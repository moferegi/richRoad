package initialize

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	wxpayModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const migrateClientTryonPointKey = "migration_client_user_tryon_point_v1"
const migrateTryonPointStatsEventKey = "migration_tryon_point_stats_event_v1"
const migrateTryonStatsEventKey = "migration_tryon_stats_event_v1"

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(client.ClientUser{}, shop.Banner{}, shop.Category{}, shop.Good{}, shop.Sku{}, shop.Cart{}, shop.CreateOrder{}, shop.OrderDetail{}, client.Address{}, client.Collect{}, shop.Comment{}, shop.Tag{}, shop.Coupon{}, shop.CouponOrderUser{}, shop.Promotion{}, shop.History{}, client.PointRecord{}, client.TryonPointStatsEvent{}, client.TryonRechargeOrder{}, client.TryonTask{}, client.TryonModel{}, client.TryonCloth{}, client.ModelCallLog{}, client.TryonStatsEvent{}, wxpayModel.Order{}, shop.Kefu{}, client.VisitorLog{}, client.VisitorSummary{}, client.SysConfig{}, shop.GoodPurchase{}, shop.QrcodePayment{}, shop.Popup{}, shop.MarketingReward{}, client.SysLanguage{}, client.PhoneAreaCode{}, client.SignIn{}, client.ExternalLinkDomain{}, client.LoginFailRecord{}, client.RegisterIPRecord{}, shop.SkuSpec{}, client.VideoTag{}, client.VideoEpisodeTag{}, client.DiaryTag{}, client.DiaryTagRelation{}, client.GameCategory{}, client.GameDifficultyCategory{}, client.GameLevel{}, client.GameUserProgress{}, client.PwdGameLevel{})
	if err != nil {
		return err
	}
	if err = migrateClientTryonPoint(db); err != nil {
		return err
	}
	if err = migrateTryonPointStatsEvent(db); err != nil {
		return err
	}
	if err = migrateTryonStatsEvent(db); err != nil {
		return err
	}
	return nil
}
func migrateClientTryonPoint(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var marker client.SysConfig
		err := tx.Where("config_key = ?", migrateClientTryonPointKey).First(&marker).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err = tx.Model(&client.ClientUser{}).Where("point > 0 AND (tryon_point = 0 OR tryon_point IS NULL)").Update("tryon_point", gorm.Expr("point")).Error; err != nil {
			return err
		}
		flag := client.SysConfig{ConfigKey: migrateClientTryonPointKey, ConfigValue: "1", ConfigName: "用户资产拆分迁移标记", ConfigGroup: "migration", Remark: "v1: 初始化 tryon_point=point，避免历史资产丢失"}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&flag).Error
	})
}
func migrateTryonPointStatsEvent(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var marker client.SysConfig
		err := tx.Where("config_key = ?", migrateTryonPointStatsEventKey).First(&marker).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		const batchSize = 500
		var lastID uint
		for {
			pointRecords := make([]client.PointRecord, 0, batchSize)
			query := tx.Where("asset_type = ?", client.AssetTypeTryonPoint)
			if lastID > 0 {
				query = query.Where("id > ?", lastID)
			}
			if err = query.Order("id ASC").Limit(batchSize).Find(&pointRecords).Error; err != nil {
				return err
			}
			if len(pointRecords) == 0 {
				break
			}
			events := make([]client.TryonPointStatsEvent, 0, len(pointRecords))
			for _, record := range pointRecords {
				eventAt := record.CreatedAt
				if eventAt.IsZero() {
					eventAt = time.Now()
				}
				userID := uint(0)
				if record.UserId != nil && *record.UserId > 0 {
					userID = uint(*record.UserId)
				}
				pointChange := 0
				if record.PointChange != nil {
					pointChange = *record.PointChange
				}
				pointAmount := pointChange
				if pointAmount < 0 {
					pointAmount = -pointAmount
				}
				changeType := ""
				if record.ChangeType != nil {
					changeType = strings.TrimSpace(*record.ChangeType)
				}
				operationType := ""
				if record.OperationType != nil {
					operationType = strings.TrimSpace(*record.OperationType)
				}
				reason := ""
				if record.Reason != nil {
					reason = strings.TrimSpace(*record.Reason)
				}
				events = append(events, client.TryonPointStatsEvent{SourcePointRecordID: record.ID, EventAt: eventAt, UserID: userID, ChangeType: changeType, OperationType: operationType, PointChange: pointChange, PointAmount: pointAmount, Reason: reason})
			}
			if len(events) > 0 {
				if err = tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "source_point_record_id"}}, DoNothing: true}).Create(&events).Error; err != nil {
					return err
				}
			}
			lastID = pointRecords[len(pointRecords)-1].ID
		}
		flag := client.SysConfig{ConfigKey: migrateTryonPointStatsEventKey, ConfigValue: "1", ConfigName: "试衣币统计事实表迁移标记", ConfigGroup: "migration", Remark: "v1: 从 client_point_records 回填 tryon_point 事实事件"}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&flag).Error
	})
}
func inferTryonModelUsageByCallStageForMigration(callStage string) string {
	stage := strings.ToLower(strings.TrimSpace(callStage))
	switch {
	case strings.Contains(stage, "refiner"):
		return "refiner"
	case strings.Contains(stage, "parsing"):
		return "parsing"
	case strings.Contains(stage, "beautify"):
		return "beautify"
	default:
		return "tryon"
	}
}
func normalizeTryonModelUsageForMigration(rawUsage string, callStage string) string {
	usage := strings.ToLower(strings.TrimSpace(rawUsage))
	switch usage {
	case "tryon", "refiner", "parsing", "beautify":
		return usage
	default:
		return inferTryonModelUsageByCallStageForMigration(callStage)
	}
}
func inferTryonProviderFromModelKeyForMigration(modelKey string) string {
	key := strings.ToLower(strings.TrimSpace(modelKey))
	if key == "" || key == "default" {
		return ""
	}
	if strings.Contains(key, "aliyun") || strings.Contains(key, "dashscope") || strings.Contains(key, "aitryon") {
		return "aliyun"
	}
	if strings.Contains(key, "gradio") || strings.Contains(key, "huggingface") || strings.Contains(key, "hf") {
		return "gradio"
	}
	return ""
}
func buildTryonStatsIdentityKeyForMigration(taskID uint, usage string) string {
	return strconv.FormatUint(uint64(taskID), 10) + "|" + strings.ToLower(strings.TrimSpace(usage))
}
func migrateTryonStatsEvent(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var marker client.SysConfig
		err := tx.Where("config_key = ?", migrateTryonStatsEventKey).First(&marker).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		const batchSize = 500
		var lastID uint
		type modelIdentity struct {
			ModelKey  string
			ModelName string
			Provider  string
		}
		identityCache := make(map[string]modelIdentity)
		for {
			logs := make([]client.ModelCallLog, 0, batchSize)
			query := tx.Order("id ASC")
			if lastID > 0 {
				query = query.Where("id > ?", lastID)
			}
			if err = query.Limit(batchSize).Find(&logs).Error; err != nil {
				return err
			}
			if len(logs) == 0 {
				break
			}
			events := make([]client.TryonStatsEvent, 0, len(logs))
			for _, logRow := range logs {
				usage := normalizeTryonModelUsageForMigration(logRow.ModelUsage, logRow.CallStage)
				identityKey := buildTryonStatsIdentityKeyForMigration(logRow.TaskID, usage)
				modelKey := strings.TrimSpace(logRow.ModelKey)
				modelName := strings.TrimSpace(logRow.ModelName)
				provider := strings.TrimSpace(logRow.Provider)
				if cached, exists := identityCache[identityKey]; exists {
					if modelKey == "" {
						modelKey = cached.ModelKey
					}
					if modelName == "" {
						modelName = cached.ModelName
					}
					if provider == "" {
						provider = cached.Provider
					}
				}
				if provider == "" {
					provider = inferTryonProviderFromModelKeyForMigration(modelKey)
				}
				if strings.TrimSpace(modelKey) != "" || strings.TrimSpace(modelName) != "" || strings.TrimSpace(provider) != "" {
					identityCache[identityKey] = modelIdentity{ModelKey: modelKey, ModelName: modelName, Provider: provider}
				}
				eventAt := logRow.CreatedAt
				if eventAt.IsZero() {
					eventAt = time.Now()
				}
				events = append(events, client.TryonStatsEvent{SourceLogID: logRow.ID, EventAt: eventAt, UserID: logRow.UserID, TaskID: logRow.TaskID, TaskNo: strings.TrimSpace(logRow.TaskNo), RequestID: strings.TrimSpace(logRow.RequestID), SceneType: strings.TrimSpace(logRow.SceneType), ModelUsage: usage, ModelKey: modelKey, Provider: provider, Status: strings.TrimSpace(logRow.Status), CostPoints: logRow.CostPoints, RefundPoints: logRow.RefundPoints})
			}
			if len(events) > 0 {
				if err = tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "source_log_id"}}, DoUpdates: clause.AssignmentColumns([]string{"event_at", "user_id", "task_id", "task_no", "request_id", "scene_type", "model_usage", "model_key", "provider", "status", "cost_points", "refund_points", "updated_at"})}).Create(&events).Error; err != nil {
					return err
				}
			}
			lastID = logs[len(logs)-1].ID
		}
		flag := client.SysConfig{ConfigKey: migrateTryonStatsEventKey, ConfigValue: "1", ConfigName: "试衣模型统计事实表迁移标记", ConfigGroup: "migration", Remark: "v1: 从 client_model_call_log 回填并修复模型统计事实事件"}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&flag).Error
	})
}
