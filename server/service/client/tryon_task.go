package client

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/external"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	tryonTaskStatusProcessing = "processing"
	tryonTaskStatusSuccess    = "success"
	tryonTaskStatusFailed     = "failed"
	tryonTaskStatusUnknown    = "unknown"

	tryonOpConsume        = "tryon_consume"
	tryonOpRefund         = "tryon_refund"
	tryonOpGuestInit      = "tryon_guest_init"
	tryonOpRegisterReward = "tryon_register_reward"

	aliyunTryonSynthesisURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/image2image/image-synthesis"
	aliyunParsingProcessURL = "https://dashscope.aliyuncs.com/api/v1/services/vision/image-process/process"
)

type TryonTaskService struct{}

type TryonTaskTrendItem struct {
	Date       string `json:"date"`
	Total      int64  `json:"total"`
	Processing int64  `json:"processing"`
	Success    int64  `json:"success"`
	Failed     int64  `json:"failed"`
}

type tryonModelConfig struct {
	Key           string   `json:"key"`
	Enabled       *bool    `json:"enabled"`
	SceneType     string   `json:"sceneType"`
	Scenes        []string `json:"scenes"`
	Model         string   `json:"model"`
	Cost          int      `json:"cost"`
	Provider      string   `json:"provider"`
	Mode          string   `json:"mode"`
	URL           string   `json:"url"`
	Token         string   `json:"token"`
	ProviderURL   string   `json:"providerUrl"`
	ProviderToken string   `json:"providerToken"`
	TaskQueryURL  string   `json:"taskQueryUrl"`
	Resolution    int      `json:"resolution"`
	RestoreFace   *bool    `json:"restoreFace"`
	ClothesType   []string `json:"clothesType"`
}

func (m *tryonModelConfig) isEnabled() bool {
	if m == nil {
		return false
	}
	if m.Enabled == nil {
		return true
	}
	return *m.Enabled
}

func (m *tryonModelConfig) matchScene(sceneType string) bool {
	if m == nil {
		return false
	}
	sceneType = strings.TrimSpace(sceneType)
	if sceneType == "" {
		return true
	}
	if len(m.Scenes) > 0 {
		for _, scene := range m.Scenes {
			if strings.EqualFold(strings.TrimSpace(scene), sceneType) {
				return true
			}
		}
		return false
	}
	if strings.TrimSpace(m.SceneType) == "" {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(m.SceneType), sceneType)
}

func (m *tryonModelConfig) endpointURL() string {
	if m == nil {
		return ""
	}
	if strings.TrimSpace(m.URL) != "" {
		return strings.TrimSpace(m.URL)
	}
	return strings.TrimSpace(m.ProviderURL)
}

func (m *tryonModelConfig) authToken() string {
	if m == nil {
		return ""
	}
	if strings.TrimSpace(m.Token) != "" {
		return strings.TrimSpace(m.Token)
	}
	return strings.TrimSpace(m.ProviderToken)
}

func (m *tryonModelConfig) queryTaskURL() string {
	if m == nil {
		return ""
	}
	if strings.TrimSpace(m.TaskQueryURL) != "" {
		return strings.TrimSpace(m.TaskQueryURL)
	}
	return ""
}

func (m *tryonModelConfig) restoreFaceValue() bool {
	if m == nil || m.RestoreFace == nil {
		return true
	}
	return *m.RestoreFace
}

func (m *tryonModelConfig) resolutionValue() int {
	if m == nil {
		return -1
	}
	if m.Resolution == 0 {
		return -1
	}
	return m.Resolution
}

type tryonProviderResponse struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	ResultImage string `json:"resultImage"`
	ResultURL   string `json:"resultUrl"`
	Data        struct {
		ResultImage string `json:"resultImage"`
		ResultURL   string `json:"resultUrl"`
		Image       string `json:"image"`
		URL         string `json:"url"`
	} `json:"data"`
}

type tryonInvokeResult struct {
	Status         string
	ResultImage    string
	ProviderTaskID string
	ErrorMessage   string
}

type dashscopeTryonCreateResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Output    struct {
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"`
	} `json:"output"`
}

type dashscopeTaskQueryResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Output    struct {
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"`
		ImageURL   string `json:"image_url"`
		Code       string `json:"code"`
		Message    string `json:"message"`
	} `json:"output"`
}

type dashscopeParsingResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Output    struct {
		Bbox          [][]int  `json:"bbox"`
		ParsingImgURL []string `json:"parsing_img_url"`
		CropImgURL    []string `json:"crop_img_url"`
	} `json:"output"`
}

// CreateTryonTask 创建并处理试衣任务：先扣币，模型失败时全额退币
func (s *TryonTaskService) CreateTryonTask(ctx context.Context, userID uint, req clientReq.CreateTryonTaskReq) (task client.TryonTask, reused bool, err error) {
	if userID == 0 {
		return task, false, errors.New("用户未登录")
	}

	requestID := strings.TrimSpace(req.RequestID)
	if requestID == "" {
		requestID = generateTryonRequestID()
	}

	var existingTask client.TryonTask
	findErr := global.GVA_DB.Where("user_id = ? AND request_id = ?", userID, requestID).First(&existingTask).Error
	if findErr == nil {
		return existingTask, true, nil
	}
	if findErr != nil && !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return task, false, findErr
	}

	sysConfigService := SysConfigService{}
	costPoints := sysConfigService.GetConfigIntByKey("tryon_cost_points", 1)
	if costPoints < 1 {
		costPoints = 1
	}

	modelCfg, modelErr := s.resolveTryonModelConfig(req.SceneType, req.ModelKey)
	if modelErr != nil {
		return task, false, modelErr
	}
	if modelCfg != nil && modelCfg.Cost > 0 {
		costPoints = modelCfg.Cost
	}

	providerName := "external"
	if modelCfg != nil {
		if strings.TrimSpace(modelCfg.Key) != "" {
			providerName = strings.TrimSpace(modelCfg.Key)
		} else if strings.TrimSpace(modelCfg.Provider) != "" {
			providerName = strings.TrimSpace(modelCfg.Provider)
		}
	}

	task = client.TryonTask{
		UserID:        userID,
		RequestID:     requestID,
		TaskNo:        generateTryonTaskNo(),
		SceneType:     req.SceneType,
		Status:        tryonTaskStatusProcessing,
		SourceImage:   req.SourceImage,
		TemplateImage: req.TemplateImage,
		Provider:      providerName,
		CostPoints:    costPoints,
	}

	if err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx", tx)
		if giftErr := s.ensureGuestInitPoints(txCtx, userID); giftErr != nil {
			return giftErr
		}

		if createErr := tx.Create(&task).Error; createErr != nil {
			return createErr
		}

		pointRecordService := PointRecordService{}
		record := buildPointRecord(userID, "decrease", costPoints, tryonOpConsume, "试衣任务扣费", "tryon_task:"+task.TaskNo)
		if pointErr := pointRecordService.CreatePointRecord(txCtx, record); pointErr != nil {
			return pointErr
		}
		return nil
	}); err != nil {
		return task, false, err
	}

	invokeResult, invokeErr := s.invokeTryonProvider(req, modelCfg)
	if invokeErr != nil {
		refundErr := s.markFailedAndRefund(ctx, &task, invokeErr.Error())
		if refundErr != nil {
			global.GVA_LOG.Error("试衣失败退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
			return task, false, errors.New("试衣失败且退币异常，请联系管理员")
		}
		_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
		return task, false, invokeErr
	}

	switch invokeResult.Status {
	case tryonTaskStatusProcessing:
		updateResult := global.GVA_DB.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ? AND status = ?", task.ID, userID, tryonTaskStatusProcessing).
			Updates(map[string]interface{}{
				"provider_task_id": strings.TrimSpace(invokeResult.ProviderTaskID),
				"error_message":    "",
			})
		if updateResult.Error != nil {
			return task, false, updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return task, false, errors.New("任务状态已变更，无法更新异步任务")
		}
	case tryonTaskStatusSuccess:
		completeAt := time.Now()
		updateResult := global.GVA_DB.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ? AND status = ?", task.ID, userID, tryonTaskStatusProcessing).
			Updates(map[string]interface{}{
				"status":           tryonTaskStatusSuccess,
				"provider_task_id": strings.TrimSpace(invokeResult.ProviderTaskID),
				"result_image":     strings.TrimSpace(invokeResult.ResultImage),
				"error_message":    "",
				"completed_at":     completeAt,
			})
		if updateResult.Error != nil {
			return task, false, updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return task, false, errors.New("任务状态已变更，无法完成")
		}
	case tryonTaskStatusFailed:
		failReason := invokeResult.ErrorMessage
		if strings.TrimSpace(failReason) == "" {
			failReason = "试衣模型处理失败"
		}
		refundErr := s.markFailedAndRefund(ctx, &task, failReason)
		if refundErr != nil {
			global.GVA_LOG.Error("试衣失败退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
			return task, false, errors.New("试衣失败且退币异常，请联系管理员")
		}
		_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
		return task, false, errors.New(failReason)
	default:
		refundErr := s.markFailedAndRefund(ctx, &task, "试衣模型返回未知状态")
		if refundErr != nil {
			global.GVA_LOG.Error("试衣失败退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
			return task, false, errors.New("试衣失败且退币异常，请联系管理员")
		}
		_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
		return task, false, errors.New("试衣模型返回未知状态")
	}

	_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
	return task, false, nil
}

// GetTryonTaskByID 按ID获取当前用户的试衣任务
func (s *TryonTaskService) GetTryonTaskByID(ctx context.Context, id uint, userID uint) (task client.TryonTask, err error) {
	err = global.GVA_DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error
	if err != nil {
		return
	}
	if refreshErr := s.refreshTryonTaskStatus(ctx, &task); refreshErr != nil {
		global.GVA_LOG.Warn("刷新试衣任务状态失败", zap.Error(refreshErr), zap.Uint("taskID", task.ID))
	}
	return
}

// GetMyTryonTaskList 获取当前用户试衣任务列表
func (s *TryonTaskService) GetMyTryonTaskList(userID uint, info clientReq.TryonTaskSearch) (list []client.TryonTask, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.TryonTask{}).Where("user_id = ?", userID)

	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.SceneType != "" {
		db = db.Where("scene_type = ?", info.SceneType)
	}
	if info.StartCreatedAt != nil {
		db = db.Where("created_at >= ?", info.StartCreatedAt)
	}
	if info.EndCreatedAt != nil {
		db = db.Where("created_at <= ?", info.EndCreatedAt)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return
	}
	ctx := context.Background()
	for i := range list {
		if !strings.EqualFold(strings.TrimSpace(list[i].Status), tryonTaskStatusProcessing) {
			continue
		}
		if strings.TrimSpace(list[i].ProviderTaskID) == "" {
			continue
		}
		if refreshErr := s.refreshTryonTaskStatus(ctx, &list[i]); refreshErr != nil {
			global.GVA_LOG.Warn("列表刷新试衣任务状态失败", zap.Error(refreshErr), zap.Uint("taskID", list[i].ID))
		}
	}
	return
}

// GetTryonTaskList 获取管理端试衣任务列表
func (s *TryonTaskService) GetTryonTaskList(info clientReq.TryonTaskSearch) (list []client.TryonTask, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, true)

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return
	}
	ctx := context.Background()
	for i := range list {
		if !strings.EqualFold(strings.TrimSpace(list[i].Status), tryonTaskStatusProcessing) {
			continue
		}
		if strings.TrimSpace(list[i].ProviderTaskID) == "" {
			continue
		}
		if refreshErr := s.refreshTryonTaskStatus(ctx, &list[i]); refreshErr != nil {
			global.GVA_LOG.Warn("管理端列表刷新试衣任务状态失败", zap.Error(refreshErr), zap.Uint("taskID", list[i].ID))
		}
	}
	return
}

// GetTryonTaskStats 获取管理端试衣任务统计（按筛选条件）
func (s *TryonTaskService) GetTryonTaskStats(info clientReq.TryonTaskSearch) (total, processing, success, failed int64, err error) {
	if err = s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, false).Count(&total).Error; err != nil {
		return
	}
	if err = s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, false).Where("status = ?", tryonTaskStatusProcessing).Count(&processing).Error; err != nil {
		return
	}
	if err = s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, false).Where("status = ?", tryonTaskStatusSuccess).Count(&success).Error; err != nil {
		return
	}
	err = s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, false).Where("status = ?", tryonTaskStatusFailed).Count(&failed).Error
	return
}

// GetTryonTaskTrend 获取管理端试衣任务趋势（按天）
func (s *TryonTaskService) GetTryonTaskTrend(info clientReq.TryonTaskSearch) (list []TryonTaskTrendItem, err error) {
	trendDays := info.TrendDays
	if trendDays <= 0 {
		trendDays = 7
	}
	if trendDays > 90 {
		trendDays = 90
	}

	now := time.Now()
	var startTime, endTime time.Time
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		if info.EndCreatedAt.Before(*info.StartCreatedAt) {
			return nil, errors.New("开始时间不能晚于结束时间")
		}
		startTime = *info.StartCreatedAt
		endTime = *info.EndCreatedAt
	} else {
		dayStartNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startTime = dayStartNow.AddDate(0, 0, -trendDays+1)
		endTime = now
	}

	startDay := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	endDay := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())
	if endDay.Before(startDay) {
		return nil, errors.New("开始时间不能晚于结束时间")
	}

	queryInfo := info
	queryInfo.Status = ""
	queryInfo.StartCreatedAt = &startTime
	queryInfo.EndCreatedAt = &endTime

	var taskList []client.TryonTask
	err = s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), queryInfo, false).
		Select("created_at", "status").
		Find(&taskList).Error
	if err != nil {
		return nil, err
	}

	trendMap := make(map[string]*TryonTaskTrendItem)
	for day := startDay; !day.After(endDay); day = day.AddDate(0, 0, 1) {
		dateKey := day.Format("2006-01-02")
		trendMap[dateKey] = &TryonTaskTrendItem{Date: dateKey}
	}

	for _, task := range taskList {
		dateKey := task.CreatedAt.Format("2006-01-02")
		item, ok := trendMap[dateKey]
		if !ok {
			continue
		}
		item.Total++
		switch task.Status {
		case tryonTaskStatusProcessing:
			item.Processing++
		case tryonTaskStatusSuccess:
			item.Success++
		case tryonTaskStatusFailed:
			item.Failed++
		}
	}

	for day := startDay; !day.After(endDay); day = day.AddDate(0, 0, 1) {
		dateKey := day.Format("2006-01-02")
		if item, ok := trendMap[dateKey]; ok {
			list = append(list, *item)
		}
	}
	return
}

func (s *TryonTaskService) applyTryonTaskFilters(db *gorm.DB, info clientReq.TryonTaskSearch, includeStatus bool) *gorm.DB {
	if info.UserID > 0 {
		db = db.Where("user_id = ?", info.UserID)
	}
	if strings.TrimSpace(info.TaskNo) != "" {
		db = db.Where("task_no = ?", strings.TrimSpace(info.TaskNo))
	}
	if strings.TrimSpace(info.RequestID) != "" {
		db = db.Where("request_id = ?", strings.TrimSpace(info.RequestID))
	}
	if strings.TrimSpace(info.SceneType) != "" {
		db = db.Where("scene_type = ?", strings.TrimSpace(info.SceneType))
	}
	if includeStatus && strings.TrimSpace(info.Status) != "" {
		db = db.Where("status = ?", strings.TrimSpace(info.Status))
	}
	if info.StartCreatedAt != nil {
		db = db.Where("created_at >= ?", info.StartCreatedAt)
	}
	if info.EndCreatedAt != nil {
		db = db.Where("created_at <= ?", info.EndCreatedAt)
	}
	return db
}

func (s *TryonTaskService) markFailedAndRefund(ctx context.Context, task *client.TryonTask, failReason string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		failReason = truncateTryonError(failReason)
		updateResult := tx.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ? AND status = ?", task.ID, task.UserID, tryonTaskStatusProcessing).
			Updates(map[string]interface{}{
				"status":        tryonTaskStatusFailed,
				"error_message": failReason,
				"refund_points": task.CostPoints,
				"completed_at":  now,
				"refunded_at":   now,
			})
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return errors.New("任务状态已变更，无法执行退币")
		}

		pointRecordService := PointRecordService{}
		txCtx := context.WithValue(ctx, "tx", tx)
		record := buildPointRecord(task.UserID, "increase", task.CostPoints, tryonOpRefund, "试衣任务失败退币", "tryon_task:"+task.TaskNo)
		if err := pointRecordService.CreatePointRecord(txCtx, record); err != nil {
			return err
		}
		return nil
	})
}

// GrantRegisterRewardPoints 发放注册奖励试衣币（幂等）
func (s *TryonTaskService) GrantRegisterRewardPoints(ctx context.Context, userID uint) error {
	if userID == 0 {
		return errors.New("用户ID无效")
	}

	rewardPoints := s.getConfigIntAllowZero("tryon_register_reward_points", 8)
	if rewardPoints <= 0 {
		return nil
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx", tx)
		return s.grantPointsIfNotExists(txCtx, tx, userID, tryonOpRegisterReward, rewardPoints, "注册奖励试衣币", "register")
	})
}

func (s *TryonTaskService) loadTryonModelConfigList() ([]tryonModelConfig, error) {
	sysConfigService := SysConfigService{}
	raw, err := sysConfigService.GetConfigByKey("tryon_models")
	if err != nil || strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	list := make([]tryonModelConfig, 0)
	if unmarshalErr := json.Unmarshal([]byte(raw), &list); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return list, nil
}

func (s *TryonTaskService) resolveTryonModelConfig(sceneType string, modelKey string) (*tryonModelConfig, error) {
	list, err := s.loadTryonModelConfigList()
	if err != nil {
		return nil, errors.New("试衣模型配置格式错误")
	}
	if len(list) == 0 {
		return nil, nil
	}

	modelKey = strings.TrimSpace(modelKey)
	var fallback *tryonModelConfig

	for i := range list {
		item := &list[i]
		if !item.isEnabled() {
			continue
		}
		if !item.matchScene(sceneType) {
			continue
		}
		if fallback == nil {
			fallback = item
		}
		if modelKey != "" && strings.EqualFold(strings.TrimSpace(item.Key), modelKey) {
			return item, nil
		}
	}

	if modelKey != "" {
		return nil, errors.New("所选模型不可用或未启用")
	}

	return fallback, nil
}

func (s *TryonTaskService) invokeTryonProvider(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig) (tryonInvokeResult, error) {
	sysConfigService := SysConfigService{}
	providerMode, _ := sysConfigService.GetConfigByKey("tryon_provider_mode")
	providerURL, _ := sysConfigService.GetConfigByKey("tryon_provider_url")
	providerToken, _ := sysConfigService.GetConfigByKey("tryon_provider_token")

	if modelCfg != nil {
		if strings.TrimSpace(modelCfg.Mode) != "" {
			providerMode = strings.TrimSpace(modelCfg.Mode)
		}
		if modelURL := modelCfg.endpointURL(); modelURL != "" {
			providerURL = modelURL
		}
		if modelToken := modelCfg.authToken(); modelToken != "" {
			providerToken = modelToken
		}
	}

	if strings.EqualFold(providerMode, "mock_success") {
		return tryonInvokeResult{
			Status:      tryonTaskStatusSuccess,
			ResultImage: req.SourceImage,
		}, nil
	}

	modelName := resolveAliyunModelName(modelCfg, req.SceneType)
	providerName := ""
	if modelCfg != nil {
		providerName = modelCfg.Provider
	}

	resolvedReq := req
	if isAliyunParsingModel(providerName, modelName, providerURL, req.SceneType) || isAliyunTryonModel(providerName, modelName, providerURL, req.SceneType) {
		normalizedSource, normalizeErr := s.normalizeAliyunMediaURL(req.SourceImage)
		if normalizeErr != nil {
			return tryonInvokeResult{}, normalizeErr
		}
		resolvedReq.SourceImage = normalizedSource

		if !strings.EqualFold(strings.TrimSpace(req.SceneType), "takeoff") {
			normalizedTemplate, normalizeErr := s.normalizeAliyunMediaURL(req.TemplateImage)
			if normalizeErr != nil {
				return tryonInvokeResult{}, normalizeErr
			}
			resolvedReq.TemplateImage = normalizedTemplate
		}
	}

	if isAliyunParsingModel(providerName, modelName, providerURL, req.SceneType) {
		return s.invokeAliyunParsing(resolvedReq, modelCfg, providerToken, providerURL, modelName)
	}

	if isAliyunTryonModel(providerName, modelName, providerURL, req.SceneType) {
		return s.invokeAliyunTryonAsync(resolvedReq, modelCfg, providerToken, providerURL, modelName)
	}

	if strings.TrimSpace(providerURL) == "" {
		return tryonInvokeResult{}, errors.New("试衣模型服务未配置")
	}

	headers := map[string]string{}
	if strings.TrimSpace(providerToken) != "" {
		headers["Authorization"] = "Bearer " + strings.TrimSpace(providerToken)
	}

	body := map[string]interface{}{
		"sceneType":     resolvedReq.SceneType,
		"sourceImage":   resolvedReq.SourceImage,
		"templateImage": resolvedReq.TemplateImage,
	}
	if strings.TrimSpace(resolvedReq.ModelKey) != "" {
		body["modelKey"] = strings.TrimSpace(resolvedReq.ModelKey)
	}

	resp, err := external.HttpRequest[tryonProviderResponse](external.RequestParams{
		URL:     strings.TrimSpace(providerURL),
		Method:  "POST",
		Headers: headers,
		Body:    body,
	})
	if err != nil {
		return tryonInvokeResult{}, err
	}
	if resp == nil {
		return tryonInvokeResult{}, errors.New("试衣模型返回为空")
	}

	resultImage := firstNonEmptyString(
		resp.ResultImage,
		resp.ResultURL,
		resp.Data.ResultImage,
		resp.Data.ResultURL,
		resp.Data.Image,
		resp.Data.URL,
	)
	if resultImage == "" {
		if strings.TrimSpace(resp.Msg) != "" {
			return tryonInvokeResult{}, errors.New(strings.TrimSpace(resp.Msg))
		}
		return tryonInvokeResult{}, errors.New("试衣模型返回结果为空")
	}

	return tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage}, nil
}

func (s *TryonTaskService) invokeAliyunTryonAsync(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, providerToken string, providerURL string, modelName string) (tryonInvokeResult, error) {
	if strings.TrimSpace(providerToken) == "" {
		return tryonInvokeResult{}, errors.New("阿里试衣API Key未配置")
	}

	createURL := strings.TrimSpace(providerURL)
	if createURL == "" {
		createURL = aliyunTryonSynthesisURL
	}

	resolvedModelName := strings.TrimSpace(modelName)
	if resolvedModelName == "" {
		resolvedModelName = "aitryon"
	}

	if strings.TrimSpace(req.SourceImage) == "" {
		return tryonInvokeResult{}, errors.New("模特图不能为空")
	}
	if strings.TrimSpace(req.TemplateImage) == "" {
		return tryonInvokeResult{}, errors.New("服饰图不能为空")
	}

	body := map[string]interface{}{
		"model": resolvedModelName,
		"input": map[string]interface{}{
			"person_image_url": strings.TrimSpace(req.SourceImage),
			"top_garment_url":  strings.TrimSpace(req.TemplateImage),
		},
		"parameters": map[string]interface{}{
			"resolution":   modelCfg.resolutionValue(),
			"restore_face": modelCfg.restoreFaceValue(),
		},
	}

	resp, err := external.HttpRequest[dashscopeTryonCreateResponse](external.RequestParams{
		URL:    createURL,
		Method: "POST",
		Headers: map[string]string{
			"Authorization":     "Bearer " + strings.TrimSpace(providerToken),
			"X-DashScope-Async": "enable",
		},
		Body: body,
	})
	if err != nil {
		return tryonInvokeResult{}, err
	}
	if resp == nil {
		return tryonInvokeResult{}, errors.New("阿里试衣返回为空")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code)}, nil
	}

	taskID := strings.TrimSpace(resp.Output.TaskID)
	if taskID == "" {
		return tryonInvokeResult{}, errors.New("阿里试衣未返回task_id")
	}

	status := normalizeProviderTaskStatus(resp.Output.TaskStatus)
	if status == tryonTaskStatusFailed {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ProviderTaskID: taskID, ErrorMessage: formatAliyunModelError(resp.Message, resp.Code)}, nil
	}
	if status == tryonTaskStatusSuccess {
		queryResult, queryErr := s.queryAliyunTryonTask(taskID, strings.TrimSpace(providerToken), createURL, modelCfg.queryTaskURL())
		if queryErr != nil {
			return tryonInvokeResult{}, queryErr
		}
		queryResult.ProviderTaskID = taskID
		return queryResult, nil
	}

	return tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}, nil
}

func (s *TryonTaskService) invokeAliyunParsing(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, providerToken string, providerURL string, modelName string) (tryonInvokeResult, error) {
	if strings.TrimSpace(providerToken) == "" {
		return tryonInvokeResult{}, errors.New("阿里图片分割API Key未配置")
	}

	processURL := strings.TrimSpace(providerURL)
	if processURL == "" {
		processURL = aliyunParsingProcessURL
	}

	resolvedModelName := strings.TrimSpace(modelName)
	if resolvedModelName == "" {
		resolvedModelName = "aitryon-parsing-v1"
	}

	clothesType := []string{"upper"}
	if modelCfg != nil && len(modelCfg.ClothesType) > 0 {
		clothesType = modelCfg.ClothesType
	}

	body := map[string]interface{}{
		"model": resolvedModelName,
		"input": map[string]interface{}{
			"image_url": strings.TrimSpace(req.SourceImage),
		},
		"parameters": map[string]interface{}{
			"clothes_type": clothesType,
		},
	}

	resp, err := external.HttpRequest[dashscopeParsingResponse](external.RequestParams{
		URL:    processURL,
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + strings.TrimSpace(providerToken),
		},
		Body: body,
	})
	if err != nil {
		return tryonInvokeResult{}, err
	}
	if resp == nil {
		return tryonInvokeResult{}, errors.New("阿里图片分割返回为空")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code)}, nil
	}

	resultImage := ""
	for _, item := range resp.Output.CropImgURL {
		resultImage = strings.TrimSpace(item)
		if resultImage != "" {
			break
		}
	}
	if resultImage == "" {
		for _, item := range resp.Output.ParsingImgURL {
			resultImage = strings.TrimSpace(item)
			if resultImage != "" {
				break
			}
		}
	}

	if resultImage == "" {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: "图片分割未识别到可用服饰区域"}, nil
	}

	return tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage}, nil
}

func (s *TryonTaskService) queryAliyunTryonTask(taskID string, providerToken string, createURL string, taskQueryURL string) (tryonInvokeResult, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return tryonInvokeResult{}, errors.New("task_id不能为空")
	}
	if strings.TrimSpace(providerToken) == "" {
		return tryonInvokeResult{}, errors.New("阿里试衣API Key未配置")
	}

	queryURL := buildAliyunTaskQueryURL(createURL, taskID)
	if strings.TrimSpace(taskQueryURL) != "" {
		queryURL = formatTaskQueryURL(taskQueryURL, taskID)
	}

	resp, err := external.HttpRequest[dashscopeTaskQueryResponse](external.RequestParams{
		URL:    queryURL,
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + strings.TrimSpace(providerToken),
		},
		Body: map[string]interface{}{},
	})
	if err != nil {
		return tryonInvokeResult{}, err
	}
	if resp == nil {
		return tryonInvokeResult{}, errors.New("查询阿里试衣任务返回为空")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code), ProviderTaskID: taskID}, nil
	}

	status := normalizeProviderTaskStatus(resp.Output.TaskStatus)
	switch status {
	case tryonTaskStatusSuccess:
		resultImage := strings.TrimSpace(resp.Output.ImageURL)
		if resultImage == "" {
			return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: "试衣任务成功但未返回图片", ProviderTaskID: taskID}, nil
		}
		return tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage, ProviderTaskID: taskID}, nil
	case tryonTaskStatusFailed:
		return tryonInvokeResult{
			Status:         tryonTaskStatusFailed,
			ProviderTaskID: taskID,
			ErrorMessage: formatAliyunModelError(
				resp.Output.Message,
				resp.Output.Code,
				resp.Message,
				resp.Code,
				"试衣任务失败",
			),
		}, nil
	case tryonTaskStatusProcessing, tryonTaskStatusUnknown:
		return tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}, nil
	default:
		return tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}, nil
	}
}

func formatAliyunModelError(values ...string) string {
	msg := firstNonEmptyString(values...)
	if msg == "" {
		return "试衣任务失败"
	}

	lowerMsg := strings.ToLower(msg)
	if strings.Contains(lowerMsg, "unable to download the media resource") || strings.Contains(lowerMsg, "invalid image url") {
		return "模型服务无法下载图片资源，请确认 sourceImage/templateImage 为公网可访问 URL（不要使用 localhost 或内网地址）"
	}
	if strings.Contains(lowerMsg, "download ") && strings.Contains(lowerMsg, "refused") {
		return "模型服务无法下载图片资源（URL被拒绝），请更换可直连的图片域名后重试"
	}
	if strings.Contains(lowerMsg, "invalidurl") {
		return "图片地址无效或模型服务不可达，请更换可直连的图片域名后重试"
	}
	if strings.Contains(lowerMsg, "data inspection") {
		return "图片风控或资源检测失败，请更换图片后重试"
	}
	return msg
}

func (s *TryonTaskService) normalizeAliyunMediaURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", errors.New("图片地址不能为空")
	}
	if strings.HasPrefix(strings.ToLower(rawURL), "data:") {
		return "", errors.New("图片地址不能是 data URI，请先上传并使用公网 URL")
	}

	publicBaseURL := s.getTryonMediaPublicBaseURL()

	if !isAbsoluteHTTPURL(rawURL) {
		if publicBaseURL == "" {
			return "", errors.New("图片地址不是公网 URL，请在系统参数配置 tryon_media_public_base_url 后重试")
		}
		return joinBaseURLAndResourcePath(publicBaseURL, rawURL)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() == "" {
		return "", errors.New("图片地址格式错误")
	}

	if isLoopbackOrPrivateHost(parsed.Hostname()) {
		if publicBaseURL == "" {
			return "", errors.New("当前图片地址为 localhost/内网地址，模型服务无法访问；请配置 tryon_media_public_base_url 为公网域名")
		}
		resourcePath := strings.TrimSpace(parsed.Path)
		if resourcePath == "" {
			resourcePath = "/"
		}
		if parsed.RawQuery != "" {
			resourcePath = resourcePath + "?" + parsed.RawQuery
		}
		return joinBaseURLAndResourcePath(publicBaseURL, resourcePath)
	}

	// If a public media base is configured, prefer it for known internal media hosts.
	if publicBaseURL != "" {
		targetHost := urlHostname(publicBaseURL)
		currentHost := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
		if targetHost != "" && currentHost != "" && !strings.EqualFold(currentHost, targetHost) {
			rewriteHosts := buildTryonMediaRewriteHosts()
			if _, shouldRewrite := rewriteHosts[currentHost]; shouldRewrite {
				resourcePath := strings.TrimSpace(parsed.Path)
				if resourcePath == "" {
					resourcePath = "/"
				}
				if parsed.RawQuery != "" {
					resourcePath = resourcePath + "?" + parsed.RawQuery
				}
				return joinBaseURLAndResourcePath(publicBaseURL, resourcePath)
			}
		}
	}

	return rawURL, nil
}

func (s *TryonTaskService) getTryonMediaPublicBaseURL() string {
	sysConfigService := SysConfigService{}
	publicBaseURL, _ := sysConfigService.GetConfigByKey("tryon_media_public_base_url")
	publicBaseURL = strings.TrimSpace(publicBaseURL)
	if publicBaseURL != "" {
		return strings.TrimRight(publicBaseURL, "/")
	}

	cdnDomain := strings.TrimSpace(global.GVA_CONFIG.Hotlink.CdnDomain)
	if strings.HasPrefix(strings.ToLower(cdnDomain), "http://") || strings.HasPrefix(strings.ToLower(cdnDomain), "https://") {
		return strings.TrimRight(cdnDomain, "/")
	}
	return ""
}

func isAbsoluteHTTPURL(rawURL string) bool {
	lowerURL := strings.ToLower(strings.TrimSpace(rawURL))
	return strings.HasPrefix(lowerURL, "http://") || strings.HasPrefix(lowerURL, "https://")
}

func urlHostname(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}
	if !isAbsoluteHTTPURL(rawURL) {
		return ""
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(parsed.Hostname()))
}

func buildTryonMediaRewriteHosts() map[string]struct{} {
	hosts := map[string]struct{}{}
	for _, raw := range []string{global.GVA_CONFIG.Hotlink.CdnDomain, global.GVA_CONFIG.AwsS3.BaseURL} {
		host := urlHostname(raw)
		if host == "" {
			continue
		}
		hosts[host] = struct{}{}
	}
	return hosts
}

func joinBaseURLAndResourcePath(baseURL string, resourcePath string) (string, error) {
	baseURL = strings.TrimSpace(baseURL)
	resourcePath = strings.TrimSpace(resourcePath)
	if baseURL == "" {
		return "", errors.New("公网地址前缀不能为空")
	}

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil || parsedBaseURL.Scheme == "" || parsedBaseURL.Hostname() == "" {
		return "", errors.New("tryon_media_public_base_url 配置格式错误，请填写如 https://back.example.com")
	}

	if strings.HasPrefix(strings.ToLower(resourcePath), "http://") || strings.HasPrefix(strings.ToLower(resourcePath), "https://") {
		resourceURL, parseErr := url.Parse(resourcePath)
		if parseErr == nil {
			resourcePath = strings.TrimSpace(resourceURL.Path)
			if resourcePath == "" {
				resourcePath = "/"
			}
			if resourceURL.RawQuery != "" {
				resourcePath = resourcePath + "?" + resourceURL.RawQuery
			}
		}
	}

	if !strings.HasPrefix(resourcePath, "/") {
		resourcePath = "/" + resourcePath
	}
	return strings.TrimRight(baseURL, "/") + resourcePath, nil
}

func isLoopbackOrPrivateHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	host = strings.Trim(host, "[]")
	if host == "" {
		return true
	}

	if host == "localhost" || host == "0.0.0.0" || host == "::1" || host == "host.docker.internal" {
		return true
	}
	if strings.HasSuffix(host, ".local") {
		return true
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return true
	}
	return false
}

func (s *TryonTaskService) refreshTryonTaskStatus(ctx context.Context, task *client.TryonTask) error {
	if task == nil {
		return nil
	}
	if !strings.EqualFold(strings.TrimSpace(task.Status), tryonTaskStatusProcessing) {
		return nil
	}
	if strings.TrimSpace(task.ProviderTaskID) == "" {
		return nil
	}

	modelCfg, _ := s.resolveTryonModelConfig(task.SceneType, task.Provider)

	sysConfigService := SysConfigService{}
	providerURL, _ := sysConfigService.GetConfigByKey("tryon_provider_url")
	providerToken, _ := sysConfigService.GetConfigByKey("tryon_provider_token")
	providerName := ""
	modelName := ""
	taskQueryURL := ""
	if modelCfg != nil {
		if modelURL := modelCfg.endpointURL(); modelURL != "" {
			providerURL = modelURL
		}
		if modelToken := modelCfg.authToken(); modelToken != "" {
			providerToken = modelToken
		}
		providerName = modelCfg.Provider
		modelName = resolveAliyunModelName(modelCfg, task.SceneType)
		taskQueryURL = modelCfg.queryTaskURL()
	}

	if !isAliyunTryonModel(providerName, modelName, providerURL, task.SceneType) {
		return nil
	}

	result, err := s.queryAliyunTryonTask(task.ProviderTaskID, providerToken, providerURL, taskQueryURL)
	if err != nil {
		return err
	}

	switch result.Status {
	case tryonTaskStatusProcessing:
		return nil
	case tryonTaskStatusSuccess:
		now := time.Now()
		updateResult := global.GVA_DB.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ? AND status = ?", task.ID, task.UserID, tryonTaskStatusProcessing).
			Updates(map[string]interface{}{
				"status":        tryonTaskStatusSuccess,
				"result_image":  strings.TrimSpace(result.ResultImage),
				"error_message": "",
				"completed_at":  now,
			})
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
		}
		task.Status = tryonTaskStatusSuccess
		task.ResultImage = strings.TrimSpace(result.ResultImage)
		task.ErrorMessage = ""
		task.CompletedAt = &now
		return nil
	case tryonTaskStatusFailed:
		if err = s.markFailedAndRefund(ctx, task, result.ErrorMessage); err != nil {
			return err
		}
		return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
	default:
		return nil
	}
}

func resolveAliyunModelName(modelCfg *tryonModelConfig, sceneType string) string {
	if modelCfg != nil {
		if strings.TrimSpace(modelCfg.Model) != "" {
			return strings.TrimSpace(modelCfg.Model)
		}
		keyLower := strings.ToLower(strings.TrimSpace(modelCfg.Key))
		providerLower := strings.ToLower(strings.TrimSpace(modelCfg.Provider))
		switch {
		case strings.Contains(keyLower, "parsing") || strings.Contains(providerLower, "parsing"):
			return "aitryon-parsing-v1"
		case strings.Contains(keyLower, "plus") || strings.Contains(providerLower, "plus"):
			return "aitryon-plus"
		case strings.Contains(keyLower, "outfit") || strings.Contains(keyLower, "aitryon") || strings.Contains(providerLower, "outfit") || strings.Contains(providerLower, "aitryon"):
			return "aitryon"
		}
	}

	if strings.EqualFold(strings.TrimSpace(sceneType), "takeoff") {
		return "aitryon-parsing-v1"
	}

	return ""
}

func isAliyunTryonModel(provider string, modelName string, providerURL string, sceneType string) bool {
	if strings.EqualFold(strings.TrimSpace(sceneType), "takeoff") {
		return false
	}

	modelLower := strings.ToLower(strings.TrimSpace(modelName))
	if modelLower == "aitryon" || modelLower == "aitryon-plus" {
		return true
	}

	providerLower := strings.ToLower(strings.TrimSpace(provider))
	if strings.Contains(providerLower, "aliyun") || strings.Contains(providerLower, "dashscope") {
		if !strings.Contains(providerLower, "parsing") {
			return true
		}
	}

	urlLower := strings.ToLower(strings.TrimSpace(providerURL))
	if strings.Contains(urlLower, "dashscope.aliyuncs.com") && strings.Contains(urlLower, "image-synthesis") {
		return true
	}

	return false
}

func isAliyunParsingModel(provider string, modelName string, providerURL string, sceneType string) bool {
	if strings.EqualFold(strings.TrimSpace(sceneType), "takeoff") {
		return true
	}

	modelLower := strings.ToLower(strings.TrimSpace(modelName))
	if modelLower == "aitryon-parsing-v1" {
		return true
	}

	providerLower := strings.ToLower(strings.TrimSpace(provider))
	if strings.Contains(providerLower, "parsing") {
		return true
	}

	urlLower := strings.ToLower(strings.TrimSpace(providerURL))
	if strings.Contains(urlLower, "dashscope.aliyuncs.com") && strings.Contains(urlLower, "image-process/process") {
		return true
	}

	return false
}

func normalizeProviderTaskStatus(raw string) string {
	status := strings.ToUpper(strings.TrimSpace(raw))
	switch status {
	case "SUCCEEDED":
		return tryonTaskStatusSuccess
	case "FAILED", "CANCELED":
		return tryonTaskStatusFailed
	case "PENDING", "PRE-PROCESSING", "RUNNING", "POST-PROCESSING":
		return tryonTaskStatusProcessing
	case "UNKNOWN":
		return tryonTaskStatusUnknown
	default:
		if status == "" {
			return tryonTaskStatusUnknown
		}
		return tryonTaskStatusProcessing
	}
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func buildAliyunTaskQueryURL(createURL string, taskID string) string {
	if strings.TrimSpace(taskID) == "" {
		return ""
	}

	base := strings.TrimSpace(createURL)
	if base == "" {
		base = aliyunTryonSynthesisURL
	}

	u, err := url.Parse(base)
	if err != nil || u.Host == "" {
		return fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/tasks/%s", taskID)
	}
	u.Path = "/api/v1/tasks/" + taskID
	u.RawQuery = ""
	return u.String()
}

func formatTaskQueryURL(taskQueryURL string, taskID string) string {
	taskQueryURL = strings.TrimSpace(taskQueryURL)
	taskID = strings.TrimSpace(taskID)
	if taskQueryURL == "" {
		return ""
	}

	taskQueryURL = strings.ReplaceAll(taskQueryURL, "{task_id}", taskID)
	taskQueryURL = strings.ReplaceAll(taskQueryURL, "{taskId}", taskID)

	u, err := url.Parse(taskQueryURL)
	if err != nil {
		if strings.Contains(taskQueryURL, taskID) {
			return taskQueryURL
		}
		return strings.TrimRight(taskQueryURL, "/") + "/" + taskID
	}

	path := strings.TrimRight(u.Path, "/")
	if !strings.HasSuffix(path, "/"+taskID) {
		if strings.HasSuffix(path, "/tasks") {
			path = path + "/" + taskID
		} else if !strings.Contains(path, taskID) {
			path = path + "/" + taskID
		}
	}
	u.Path = path
	u.RawQuery = ""
	return u.String()
}

func generateTryonTaskNo() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "TRYON_" + time.Now().Format("20060102150405") + "_" + hex.EncodeToString(b)
}

func generateTryonRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "REQ_" + hex.EncodeToString(b)
}

func buildPointRecord(userID uint, changeType string, points int, operationType string, reason string, remark string) *client.PointRecord {
	uid := int(userID)
	ct := changeType
	p := points
	op := operationType
	rs := reason
	rm := remark

	return &client.PointRecord{
		UserId:        &uid,
		ChangeType:    &ct,
		PointChange:   &p,
		OperationType: &op,
		Reason:        &rs,
		Remark:        &rm,
	}
}

func (s *TryonTaskService) ensureGuestInitPoints(ctx context.Context, userID uint) error {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok || tx == nil {
		return errors.New("试衣赠币事务上下文缺失")
	}

	guestPoints := s.getConfigIntAllowZero("tryon_guest_init_points", 3)
	if guestPoints <= 0 {
		return nil
	}

	return s.grantPointsIfNotExists(ctx, tx, userID, tryonOpGuestInit, guestPoints, "游客首次试衣赠币", "guest_first_tryon")
}

func (s *TryonTaskService) grantPointsIfNotExists(ctx context.Context, tx *gorm.DB, userID uint, operationType string, points int, reason string, remark string) error {
	uid := int(userID)
	var count int64
	if err := tx.Model(&client.PointRecord{}).Where("user_id = ? AND operation_type = ?", uid, operationType).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	pointRecordService := PointRecordService{}
	record := buildPointRecord(userID, "increase", points, operationType, reason, remark)
	return pointRecordService.CreatePointRecord(ctx, record)
}

func (s *TryonTaskService) getConfigIntAllowZero(key string, defaultVal int) int {
	sysConfigService := SysConfigService{}
	val, err := sysConfigService.GetConfigByKey(key)
	if err != nil {
		return defaultVal
	}

	val = strings.TrimSpace(val)
	if val == "" {
		return defaultVal
	}

	n, parseErr := strconv.Atoi(val)
	if parseErr != nil || n < 0 {
		return defaultVal
	}
	return n
}

func truncateTryonError(msg string) string {
	msg = strings.TrimSpace(msg)
	if len(msg) <= 480 {
		return msg
	}
	return msg[:480]
}
