package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	englishReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	clientService "github.com/flipped-aurora/gin-vue-admin/server/service/client"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LearningAuthzService struct{}

type VideoEpisodeAccessDecision struct {
	HasGlobalFullAuth      bool   `json:"hasGlobalFullAuth"`
	HasResourceEntitlement bool   `json:"hasResourceEntitlement"`
	HasFullAuth            bool   `json:"hasFullAuth"`
	TrialOnly              bool   `json:"trialOnly"`
	Reason                 string `json:"reason"`
}

const (
	defaultLearningNewUserFreeHours = 24
	authzSourceManualGrant          = "manual_grant"
)

func normalizeResourceType(resourceType string) string {
	return strings.ToLower(strings.TrimSpace(resourceType))
}

func isSupportedResourceType(resourceType string) bool {
	switch normalizeResourceType(resourceType) {
	case model.ResourceTypeEnglishCategory, model.ResourceTypeVideoSeries, model.ResourceTypeVideoEpisode:
		return true
	default:
		return false
	}
}

func getLearningConfigIntAllowZero(key string, defaultValue int) int {
	svc := clientService.SysConfigService{}
	raw, err := svc.GetConfigByKey(key)
	if err != nil {
		return defaultValue
	}
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value < 0 {
		return defaultValue
	}
	return value
}

func (s *LearningAuthzService) buildNewUserFreeExpire(now time.Time) *time.Time {
	hours := getLearningConfigIntAllowZero("learning_new_user_free_hours", defaultLearningNewUserFreeHours)
	if hours <= 0 {
		return nil
	}
	expireAt := now.Add(time.Duration(hours) * time.Hour)
	return &expireAt
}

func (s *LearningAuthzService) EnsureUserAsset(userID uint) (model.UserLearningAsset, error) {
	return s.EnsureUserAssetWithTx(global.GVA_DB, userID)
}

func (s *LearningAuthzService) EnsureUserAssetWithTx(tx *gorm.DB, userID uint) (model.UserLearningAsset, error) {
	if userID == 0 {
		return model.UserLearningAsset{}, errors.New("用户ID不能为空")
	}

	var asset model.UserLearningAsset
	result := tx.Where("user_id = ?", userID).Limit(1).Find(&asset)
	if result.Error != nil {
		return model.UserLearningAsset{}, result.Error
	}
	if result.RowsAffected > 0 {
		return asset, nil
	}

	asset = model.UserLearningAsset{
		UserID:         userID,
		FreeTimeExpire: s.buildNewUserFreeExpire(time.Now()),
	}
	if err := tx.Create(&asset).Error; err != nil {
		return model.UserLearningAsset{}, err
	}
	return asset, nil
}

func (s *LearningAuthzService) HasGlobalFullAuth(asset model.UserLearningAsset, now time.Time) bool {
	if asset.FreeTimeExpire != nil && asset.FreeTimeExpire.After(now) {
		return true
	}
	if asset.FreeMinutes > 0 {
		return true
	}
	if asset.IsVip != nil && *asset.IsVip {
		return true
	}
	return false
}

func (s *LearningAuthzService) ResolveGlobalFullAuth(userID uint) (bool, model.UserLearningAsset, error) {
	asset, err := s.EnsureUserAsset(userID)
	if err != nil {
		return false, model.UserLearningAsset{}, err
	}
	return s.HasGlobalFullAuth(asset, time.Now()), asset, nil
}

func (s *LearningAuthzService) hasDirectActiveEntitlement(userID uint, resourceType string, resourceID uint, now time.Time) (bool, error) {
	if userID == 0 || resourceID == 0 {
		return false, nil
	}

	var count int64
	err := global.GVA_DB.Model(&model.UserLearningEntitlement{}).
		Where("user_id = ? AND resource_type = ? AND resource_id = ?", userID, resourceType, resourceID).
		Where("expire_at IS NULL OR expire_at > ?", now).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *LearningAuthzService) HasActiveEntitlement(userID uint, resourceType string, resourceID uint, now time.Time) (bool, error) {
	resourceType = normalizeResourceType(resourceType)
	if !isSupportedResourceType(resourceType) || resourceID == 0 || userID == 0 {
		return false, nil
	}

	direct, err := s.hasDirectActiveEntitlement(userID, resourceType, resourceID, now)
	if err != nil || direct {
		return direct, err
	}

	if resourceType != model.ResourceTypeVideoEpisode {
		return false, nil
	}

	var episode model.VideoEpisode
	result := global.GVA_DB.Select("id, series_id").Where("id = ?", resourceID).Limit(1).Find(&episode)
	if result.Error != nil || result.RowsAffected == 0 || episode.SeriesID == 0 {
		return false, result.Error
	}

	return s.hasDirectActiveEntitlement(userID, model.ResourceTypeVideoSeries, episode.SeriesID, now)
}

func (s *LearningAuthzService) isCategoryRestricted(category model.EnglishCategory) bool {
	needVip := category.NeedVip != nil && *category.NeedVip
	return category.Price > 0 || needVip
}

func (s *LearningAuthzService) isSeriesRestricted(series model.VideoSeries) bool {
	needVip := series.NeedVip != nil && *series.NeedVip
	return series.Price > 0 || needVip
}

func (s *LearningAuthzService) CanAccessCategory(userID uint, categoryID uint) (bool, error) {
	if userID == 0 || categoryID == 0 {
		return false, nil
	}

	var category model.EnglishCategory
	result := global.GVA_DB.Where("id = ?", categoryID).Limit(1).Find(&category)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, gorm.ErrRecordNotFound
	}

	if !s.isCategoryRestricted(category) {
		return true, nil
	}

	globalFull, _, err := s.ResolveGlobalFullAuth(userID)
	if err != nil {
		return false, err
	}
	if globalFull {
		return true, nil
	}

	return s.HasActiveEntitlement(userID, model.ResourceTypeEnglishCategory, categoryID, time.Now())
}

func (s *LearningAuthzService) EvaluateVideoEpisodeAccess(userID uint, episode model.VideoEpisode) (VideoEpisodeAccessDecision, error) {
	decision := VideoEpisodeAccessDecision{}
	if userID == 0 {
		decision.TrialOnly = true
		decision.Reason = "anonymous"
		return decision, nil
	}

	globalFull, _, err := s.ResolveGlobalFullAuth(userID)
	if err != nil {
		return decision, err
	}
	decision.HasGlobalFullAuth = globalFull
	if globalFull {
		decision.HasFullAuth = true
		decision.Reason = "global_full_auth"
		return decision, nil
	}

	var series model.VideoSeries
	if episode.SeriesID > 0 {
		if err = global.GVA_DB.Select("id, price, need_vip").Where("id = ?", episode.SeriesID).First(&series).Error; err != nil {
			return decision, err
		}
	}

	if !s.isSeriesRestricted(series) {
		decision.HasFullAuth = true
		decision.Reason = "resource_free"
		return decision, nil
	}

	entitled, err := s.HasActiveEntitlement(userID, model.ResourceTypeVideoEpisode, episode.ID, time.Now())
	if err != nil {
		return decision, err
	}
	decision.HasResourceEntitlement = entitled
	if entitled {
		decision.HasFullAuth = true
		decision.Reason = "resource_entitlement"
		return decision, nil
	}

	decision.TrialOnly = true
	decision.Reason = "trial_only"
	return decision, nil
}

func (s *LearningAuthzService) resourceExists(resourceType string, resourceID uint) (bool, error) {
	if resourceID == 0 {
		return false, nil
	}

	resourceType = normalizeResourceType(resourceType)
	var count int64
	switch resourceType {
	case model.ResourceTypeEnglishCategory:
		err := global.GVA_DB.Model(&model.EnglishCategory{}).Where("id = ?", resourceID).Count(&count).Error
		return count > 0, err
	case model.ResourceTypeVideoSeries:
		err := global.GVA_DB.Model(&model.VideoSeries{}).Where("id = ?", resourceID).Count(&count).Error
		return count > 0, err
	case model.ResourceTypeVideoEpisode:
		err := global.GVA_DB.Model(&model.VideoEpisode{}).Where("id = ?", resourceID).Count(&count).Error
		return count > 0, err
	default:
		return false, errors.New("不支持的资源类型")
	}
}

func (s *LearningAuthzService) GrantEntitlement(req englishReq.GrantEntitlementReq, grantedBy uint) error {
	if req.UserID == 0 || req.ResourceID == 0 {
		return errors.New("用户ID和资源ID不能为空")
	}

	resourceType := normalizeResourceType(req.ResourceType)
	if !isSupportedResourceType(resourceType) {
		return errors.New("资源类型不支持")
	}

	if req.ExpireAt != nil && !req.ExpireAt.After(time.Now()) {
		return errors.New("过期时间必须晚于当前时间")
	}

	exists, err := s.resourceExists(resourceType, req.ResourceID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("目标资源不存在")
	}

	entitlement := model.UserLearningEntitlement{
		UserID:       req.UserID,
		ResourceType: resourceType,
		ResourceID:   req.ResourceID,
		Source:       authzSourceManualGrant,
		GrantedBy:    grantedBy,
		ExpireAt:     req.ExpireAt,
		Remark:       strings.TrimSpace(req.Remark),
	}

	return global.GVA_DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "resource_type"}, {Name: "resource_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"source":     authzSourceManualGrant,
			"granted_by": grantedBy,
			"expire_at":  req.ExpireAt,
			"remark":     strings.TrimSpace(req.Remark),
			"updated_at": time.Now(),
		}),
	}).Create(&entitlement).Error
}

func (s *LearningAuthzService) RevokeEntitlement(req englishReq.RevokeEntitlementReq) error {
	if req.UserID == 0 || req.ResourceID == 0 {
		return errors.New("用户ID和资源ID不能为空")
	}

	resourceType := normalizeResourceType(req.ResourceType)
	if !isSupportedResourceType(resourceType) {
		return errors.New("资源类型不支持")
	}

	return global.GVA_DB.Delete(&model.UserLearningEntitlement{}, "user_id = ? AND resource_type = ? AND resource_id = ?", req.UserID, resourceType, req.ResourceID).Error
}

func (s *LearningAuthzService) GetEntitlementList(info englishReq.EntitlementSearch) (list []model.UserLearningEntitlement, total int64, page int, pageSize int, err error) {
	page, pageSize = normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Model(&model.UserLearningEntitlement{})

	if info.UserID > 0 {
		db = db.Where("user_id = ?", info.UserID)
	}
	if strings.TrimSpace(info.ResourceType) != "" {
		resourceType := normalizeResourceType(info.ResourceType)
		if !isSupportedResourceType(resourceType) {
			err = errors.New("资源类型不支持")
			return
		}
		db = db.Where("resource_type = ?", resourceType)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error
	return
}
