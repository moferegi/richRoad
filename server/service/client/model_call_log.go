package client

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

const (
	modelCallStatusSuccess = "success"
	modelCallStatusError   = "error"
)

type modelCallTraceContext struct {
	UserID               uint
	TaskID               uint
	TaskNo               string
	RequestID            string
	SceneType            string
	TemplatePart         string
	SourceImage          string
	TemplateImage        string
	EntrySource          string
	Behavior             string
	BaseCostPoints       int
	RefinerCostPoints    int
	ParsingCostPoints    int
	BeautifyCostPoints   int
	BaseRefundPoints     int
	RefinerRefundPoints  int
	ParsingRefundPoints  int
	BeautifyRefundPoints int
	TotalRefundPoints    int
	TokenValue           string
	TokenSlot            string
}

type modelCallLogInput struct {
	TraceContext    *modelCallTraceContext
	CallStage       string
	ModelKey        string
	ModelUsage      string
	ModelName       string
	Provider        string
	EndpointURL     string
	QueryURL        string
	RelatedModels   map[string]string
	RequestPayload  interface{}
	ResponsePayload interface{}
	ResultImage     string
	Status          string
	ErrorMessage    string
	StartedAt       time.Time
	FinishedAt      time.Time
	TokenValue      string
	TokenSlot       string
}

// GetModelCallLogList 分页获取模型调用日志
func (s *SysConfigService) GetModelCallLogList(info clientReq.ModelCallLogSearch) (list []client.ModelCallLog, total int64, err error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 20
	}

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.ModelCallLog{})

	if strings.TrimSpace(info.StartCreatedAt) != "" {
		db = db.Where("created_at >= ?", strings.TrimSpace(info.StartCreatedAt))
	}
	if strings.TrimSpace(info.EndCreatedAt) != "" {
		db = db.Where("created_at <= ?", strings.TrimSpace(info.EndCreatedAt))
	}
	if info.UserID > 0 {
		db = db.Where("user_id = ?", info.UserID)
	}
	if strings.TrimSpace(info.SceneType) != "" {
		db = db.Where("scene_type = ?", strings.TrimSpace(info.SceneType))
	}
	if strings.TrimSpace(info.Behavior) != "" {
		db = db.Where("behavior = ?", strings.TrimSpace(info.Behavior))
	}
	if strings.TrimSpace(info.CallStage) != "" {
		db = db.Where("call_stage = ?", strings.TrimSpace(info.CallStage))
	}
	if strings.TrimSpace(info.ModelKey) != "" {
		db = db.Where("model_key LIKE ?", "%"+strings.TrimSpace(info.ModelKey)+"%")
	}
	if strings.TrimSpace(info.ModelUsage) != "" {
		db = db.Where("model_usage = ?", strings.TrimSpace(info.ModelUsage))
	}
	if strings.TrimSpace(info.Provider) != "" {
		db = db.Where("provider = ?", strings.TrimSpace(info.Provider))
	}
	if info.OnlyFailed {
		db = db.Where("status IN ?", []string{tryonTaskStatusFailed, modelCallStatusError})
	} else if strings.TrimSpace(info.Status) != "" {
		db = db.Where("status = ?", strings.TrimSpace(info.Status))
	}
	if strings.TrimSpace(info.TaskNo) != "" {
		db = db.Where("task_no LIKE ?", "%"+strings.TrimSpace(info.TaskNo)+"%")
	}
	if strings.TrimSpace(info.RequestID) != "" {
		db = db.Where("request_id LIKE ?", "%"+strings.TrimSpace(info.RequestID)+"%")
	}

	if strings.TrimSpace(info.Keyword) != "" {
		keyword := "%" + strings.TrimSpace(info.Keyword) + "%"
		db = db.Where("task_no LIKE ? OR request_id LIKE ? OR model_key LIKE ? OR model_name LIKE ? OR provider LIKE ? OR error_message LIKE ?", keyword, keyword, keyword, keyword, keyword, keyword)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&list).Error
	return
}

func inferModelCallEntrySource(sceneType string, templatePart string) string {
	scene := strings.ToLower(strings.TrimSpace(sceneType))
	part := strings.ToLower(strings.TrimSpace(templatePart))

	switch scene {
	case "shoes":
		return "uni_shoes_room"
	case "takeoff":
		return "uni_generate_page"
	case "clothes":
		if part == "upper" || part == "lower" {
			return "uni_generate_page"
		}
		return "uni_tryon_room"
	default:
		return "uni_unknown"
	}
}

func inferModelCallBehavior(sceneType string, templatePart string) string {
	scene := strings.ToLower(strings.TrimSpace(sceneType))
	part := strings.ToLower(strings.TrimSpace(templatePart))

	switch scene {
	case "shoes":
		return "shoes_generate"
	case "takeoff":
		return "takeoff_generate"
	case "clothes":
		if part == "upper" {
			return "clothes_generate_upper"
		}
		if part == "lower" {
			return "clothes_generate_lower"
		}
		return "clothes_generate"
	default:
		return "unknown_generate"
	}
}

func createModelCallTraceContext(userID uint, task client.TryonTask, req clientReq.CreateTryonTaskReq) *modelCallTraceContext {
	ctx := &modelCallTraceContext{
		UserID:               userID,
		TaskID:               task.ID,
		TaskNo:               strings.TrimSpace(task.TaskNo),
		RequestID:            strings.TrimSpace(req.RequestID),
		SceneType:            strings.TrimSpace(req.SceneType),
		TemplatePart:         strings.TrimSpace(req.TemplatePart),
		SourceImage:          strings.TrimSpace(req.SourceImage),
		TemplateImage:        strings.TrimSpace(req.TemplateImage),
		BaseCostPoints:       task.BaseCostPoints,
		RefinerCostPoints:    task.RefinerCostPoints,
		ParsingCostPoints:    task.ParsingCostPoints,
		BeautifyCostPoints:   task.BeautifyCost,
		BaseRefundPoints:     0,
		RefinerRefundPoints:  0,
		ParsingRefundPoints:  0,
		BeautifyRefundPoints: 0,
		TotalRefundPoints:    0,
	}
	ctx.EntrySource = inferModelCallEntrySource(ctx.SceneType, ctx.TemplatePart)
	ctx.Behavior = inferModelCallBehavior(ctx.SceneType, ctx.TemplatePart)
	return ctx
}

func createBeautifyTraceContext(userID uint, task client.TryonTask) *modelCallTraceContext {
	ctx := &modelCallTraceContext{
		UserID:               userID,
		TaskID:               task.ID,
		TaskNo:               strings.TrimSpace(task.TaskNo),
		RequestID:            strings.TrimSpace(task.RequestID),
		SceneType:            strings.TrimSpace(task.SceneType),
		SourceImage:          strings.TrimSpace(task.SourceImage),
		TemplateImage:        strings.TrimSpace(task.TemplateImage),
		BaseCostPoints:       task.BaseCostPoints,
		RefinerCostPoints:    task.RefinerCostPoints,
		ParsingCostPoints:    task.ParsingCostPoints,
		BeautifyCostPoints:   task.BeautifyCost,
		RefinerRefundPoints:  task.RefinerRefund,
		ParsingRefundPoints:  task.ParsingRefund,
		BeautifyRefundPoints: task.BeautifyRefund,
		TotalRefundPoints:    task.RefundPoints,
	}
	ctx.EntrySource = inferModelCallEntrySource(ctx.SceneType, "")
	ctx.Behavior = "beautify_generate"
	return ctx
}

func createTraceContextFromTask(task *client.TryonTask) *modelCallTraceContext {
	if task == nil {
		return nil
	}
	ctx := &modelCallTraceContext{
		UserID:               task.UserID,
		TaskID:               task.ID,
		TaskNo:               strings.TrimSpace(task.TaskNo),
		RequestID:            strings.TrimSpace(task.RequestID),
		SceneType:            strings.TrimSpace(task.SceneType),
		SourceImage:          strings.TrimSpace(task.SourceImage),
		TemplateImage:        strings.TrimSpace(task.TemplateImage),
		BaseCostPoints:       task.BaseCostPoints,
		RefinerCostPoints:    task.RefinerCostPoints,
		ParsingCostPoints:    task.ParsingCostPoints,
		BeautifyCostPoints:   task.BeautifyCost,
		RefinerRefundPoints:  task.RefinerRefund,
		ParsingRefundPoints:  task.ParsingRefund,
		BeautifyRefundPoints: task.BeautifyRefund,
		TotalRefundPoints:    task.RefundPoints,
	}
	baseRefund := task.RefundPoints - task.RefinerRefund - task.ParsingRefund - task.BeautifyRefund
	if baseRefund > 0 {
		ctx.BaseRefundPoints = baseRefund
	}
	ctx.EntrySource = inferModelCallEntrySource(ctx.SceneType, "")
	ctx.Behavior = inferModelCallBehavior(ctx.SceneType, "")
	return ctx
}

func cloneTraceContext(ctx *modelCallTraceContext) *modelCallTraceContext {
	if ctx == nil {
		return nil
	}
	copyCtx := *ctx
	return &copyCtx
}

func attachTokenTraceContext(traceCtx *modelCallTraceContext, token string, tokenSlot string) *modelCallTraceContext {
	ctx := cloneTraceContext(traceCtx)
	if ctx == nil {
		ctx = &modelCallTraceContext{}
	}
	ctx.TokenValue = normalizeTokenValue(token)
	ctx.TokenSlot = strings.TrimSpace(tokenSlot)
	return ctx
}

func errorToString(err error) string {
	if err == nil {
		return ""
	}
	return strings.TrimSpace(err.Error())
}

func traceModelUsage(traceCtx *modelCallTraceContext) string {
	if traceCtx == nil {
		return ""
	}
	behavior := strings.ToLower(strings.TrimSpace(traceCtx.Behavior))
	switch {
	case strings.Contains(behavior, "beautify"):
		return "beautify"
	case strings.Contains(behavior, "refiner"):
		return "refiner"
	case strings.Contains(behavior, "parsing"):
		return "parsing"
	default:
		return ""
	}
}

func inferModelUsageByCallStage(callStage string) string {
	stage := strings.ToLower(strings.TrimSpace(callStage))
	switch {
	case strings.Contains(stage, "refiner"):
		return "refiner"
	case strings.Contains(stage, "parsing"):
		return "parsing"
	case strings.Contains(stage, "beautify"):
		return "beautify"
	case strings.Contains(stage, "tryon"):
		return "tryon"
	default:
		return ""
	}
}

func normalizeModelUsageForLog(inputUsage string, callStage string, traceCtx *modelCallTraceContext) string {
	usage := strings.TrimSpace(inputUsage)
	if usage == "" {
		usage = traceModelUsage(traceCtx)
	}
	if strings.TrimSpace(usage) == "" {
		usage = inferModelUsageByCallStage(callStage)
	}
	return normalizeTryonModelUsageValue(usage)
}

func resolveModelIdentityFromRecentLog(taskID uint, usage string, modelKey string, modelName string, provider string) (string, string, string) {
	trimmedUsage := strings.ToLower(strings.TrimSpace(usage))
	resolvedModelKey := strings.TrimSpace(modelKey)
	resolvedModelName := strings.TrimSpace(modelName)
	resolvedProvider := strings.TrimSpace(provider)

	if taskID == 0 || trimmedUsage == "" {
		if resolvedProvider == "" {
			resolvedProvider = inferProviderFromModelKey(resolvedModelKey)
		}
		return resolvedModelKey, resolvedModelName, resolvedProvider
	}

	type latestModelIdentity struct {
		ModelKey  string
		ModelName string
		Provider  string
	}
	var latest latestModelIdentity
	err := global.GVA_DB.Model(&client.ModelCallLog{}).
		Select("model_key, model_name, provider").
		Where("task_id = ? AND LOWER(TRIM(model_usage)) = ? AND (COALESCE(NULLIF(TRIM(model_key), ''), '') <> '' OR COALESCE(NULLIF(TRIM(provider), ''), '') <> '' OR COALESCE(NULLIF(TRIM(model_name), ''), '') <> '')", taskID, trimmedUsage).
		Order("id DESC").
		Take(&latest).Error
	if err == nil {
		if resolvedModelKey == "" {
			resolvedModelKey = strings.TrimSpace(latest.ModelKey)
		}
		if resolvedModelName == "" {
			resolvedModelName = strings.TrimSpace(latest.ModelName)
		}
		if resolvedProvider == "" {
			resolvedProvider = strings.TrimSpace(latest.Provider)
		}
	}

	if resolvedProvider == "" {
		resolvedProvider = inferProviderFromModelKey(resolvedModelKey)
	}

	return resolvedModelKey, resolvedModelName, resolvedProvider
}

func modelCallPointsByUsage(traceCtx *modelCallTraceContext, usage string) (costPoints int, refundPoints int) {
	if traceCtx == nil {
		return 0, 0
	}

	usage = strings.ToLower(strings.TrimSpace(usage))
	switch usage {
	case "refiner":
		costPoints = traceCtx.RefinerCostPoints
		refundPoints = traceCtx.RefinerRefundPoints
	case "parsing":
		costPoints = traceCtx.ParsingCostPoints
		refundPoints = traceCtx.ParsingRefundPoints
	case "beautify":
		costPoints = traceCtx.BeautifyCostPoints
		refundPoints = traceCtx.BeautifyRefundPoints
	default:
		costPoints = traceCtx.BaseCostPoints
		refundPoints = traceCtx.BaseRefundPoints
		if refundPoints <= 0 && traceCtx.TotalRefundPoints > 0 {
			baseRefund := traceCtx.TotalRefundPoints - traceCtx.RefinerRefundPoints - traceCtx.ParsingRefundPoints - traceCtx.BeautifyRefundPoints
			if baseRefund > 0 {
				refundPoints = baseRefund
			}
		}
	}

	if costPoints < 0 {
		costPoints = 0
	}
	if refundPoints < 0 {
		refundPoints = 0
	}
	return costPoints, refundPoints
}

func sanitizePayloadForLog(payload interface{}) interface{} {
	if payload == nil {
		return nil
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return payload
	}

	var decoded interface{}
	if err = json.Unmarshal(raw, &decoded); err != nil {
		return payload
	}

	return sanitizePayloadNode(decoded, "")
}

func sanitizePayloadNode(value interface{}, parentKey string) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(typed))
		for k, v := range typed {
			if isSensitiveLogKey(k) {
				result[k] = "***"
				continue
			}
			result[k] = sanitizePayloadNode(v, k)
		}
		return result
	case []interface{}:
		result := make([]interface{}, 0, len(typed))
		for _, item := range typed {
			result = append(result, sanitizePayloadNode(item, parentKey))
		}
		return result
	case string:
		if isSensitiveLogKey(parentKey) {
			return "***"
		}
		return truncateLogText(typed, 4096)
	default:
		return typed
	}
}

func isSensitiveLogKey(key string) bool {
	key = strings.ToLower(strings.TrimSpace(key))
	if key == "" {
		return false
	}
	sensitiveParts := []string{
		"token",
		"authorization",
		"accesskey",
		"access_key",
		"secret",
		"signature",
		"password",
		"credential",
	}
	for _, part := range sensitiveParts {
		if strings.Contains(key, part) {
			return true
		}
	}
	return false
}

func truncateLogText(value string, maxLen int) string {
	value = strings.TrimSpace(value)
	if maxLen <= 0 || len(value) <= maxLen {
		return value
	}
	if maxLen <= 32 {
		return value[:maxLen]
	}
	return value[:maxLen-16] + "...[truncated]"
}

func marshalLogPayload(payload interface{}) string {
	safePayload := sanitizePayloadForLog(payload)
	encoded, err := json.Marshal(safePayload)
	if err != nil {
		return truncateLogText("\"<payload_marshal_failed>\"", 64000)
	}
	return truncateLogText(string(encoded), 64000)
}

func buildTryonStatsEventFromModelCallLog(logRow *client.ModelCallLog) client.TryonStatsEvent {
	eventAt := time.Now()
	if logRow != nil && !logRow.CreatedAt.IsZero() {
		eventAt = logRow.CreatedAt
	}

	if logRow == nil {
		return client.TryonStatsEvent{EventAt: eventAt}
	}

	modelUsage := normalizeModelUsageForLog(logRow.ModelUsage, logRow.CallStage, nil)
	modelKey := strings.TrimSpace(logRow.ModelKey)
	provider := strings.TrimSpace(logRow.Provider)
	if provider == "" {
		provider = inferProviderFromModelKey(modelKey)
	}

	return client.TryonStatsEvent{
		SourceLogID:  logRow.ID,
		EventAt:      eventAt,
		UserID:       logRow.UserID,
		TaskID:       logRow.TaskID,
		TaskNo:       logRow.TaskNo,
		RequestID:    logRow.RequestID,
		SceneType:    logRow.SceneType,
		ModelUsage:   modelUsage,
		ModelKey:     modelKey,
		Provider:     provider,
		Status:       logRow.Status,
		CostPoints:   logRow.CostPoints,
		RefundPoints: logRow.RefundPoints,
	}
}

func saveTryonStatsEvent(logRow *client.ModelCallLog) {
	if logRow == nil {
		return
	}

	eventRow := buildTryonStatsEventFromModelCallLog(logRow)
	if err := global.GVA_DB.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "source_log_id"}}, DoUpdates: clause.Assignments(map[string]interface{}{"event_at": eventRow.EventAt, "user_id": eventRow.UserID, "task_id": eventRow.TaskID, "task_no": eventRow.TaskNo, "request_id": eventRow.RequestID, "scene_type": eventRow.SceneType, "model_usage": eventRow.ModelUsage, "model_key": eventRow.ModelKey, "provider": eventRow.Provider, "status": eventRow.Status, "cost_points": eventRow.CostPoints, "refund_points": eventRow.RefundPoints, "updated_at": time.Now()})}).Create(&eventRow).Error; err != nil {
		global.GVA_LOG.Warn("写入试衣统计事实事件失败", zap.Error(err), zap.Uint("sourceLogID", logRow.ID), zap.String("taskNo", strings.TrimSpace(logRow.TaskNo)))
	}
}

func saveModelCallLog(input modelCallLogInput) {
	if strings.TrimSpace(input.CallStage) == "" {
		return
	}

	trace := input.TraceContext
	if trace == nil {
		trace = &modelCallTraceContext{}
	}

	startedAt := input.StartedAt
	if startedAt.IsZero() {
		startedAt = time.Now()
	}
	finishedAt := input.FinishedAt
	if finishedAt.IsZero() {
		finishedAt = time.Now()
	}

	durationMs := finishedAt.Sub(startedAt).Milliseconds()
	if durationMs < 0 {
		durationMs = 0
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		if strings.TrimSpace(input.ErrorMessage) != "" {
			status = modelCallStatusError
		} else {
			status = tryonTaskStatusUnknown
		}
	}

	usage := normalizeModelUsageForLog(input.ModelUsage, input.CallStage, trace)
	costPoints, refundPoints := modelCallPointsByUsage(trace, usage)
	modelKey, modelName, provider := resolveModelIdentityFromRecentLog(trace.TaskID, usage, input.ModelKey, input.ModelName, input.Provider)
	tokenValue := normalizeTokenValue(input.TokenValue)
	if tokenValue == "" {
		tokenValue = normalizeTokenValue(trace.TokenValue)
	}
	tokenSlot := strings.TrimSpace(input.TokenSlot)
	if tokenSlot == "" {
		tokenSlot = strings.TrimSpace(trace.TokenSlot)
	}
	tokenFingerprint := buildTokenFingerprint(tokenValue)

	relatedModels := ""
	if len(input.RelatedModels) > 0 {
		relatedModels = marshalLogPayload(input.RelatedModels)
	}

	logRow := client.ModelCallLog{
		UserID:           trace.UserID,
		TaskID:           trace.TaskID,
		TaskNo:           truncateLogText(trace.TaskNo, 40),
		RequestID:        truncateLogText(trace.RequestID, 80),
		SceneType:        truncateLogText(trace.SceneType, 20),
		TemplatePart:     truncateLogText(trace.TemplatePart, 20),
		EntrySource:      truncateLogText(trace.EntrySource, 80),
		Behavior:         truncateLogText(trace.Behavior, 80),
		CallStage:        truncateLogText(input.CallStage, 80),
		ModelKey:         truncateLogText(modelKey, 120),
		ModelUsage:       truncateLogText(usage, 20),
		ModelName:        truncateLogText(modelName, 120),
		Provider:         truncateLogText(provider, 80),
		EndpointURL:      truncateLogText(input.EndpointURL, 500),
		QueryURL:         truncateLogText(input.QueryURL, 500),
		RelatedModels:    relatedModels,
		SourceImage:      truncateLogText(trace.SourceImage, 1024),
		TemplateImage:    truncateLogText(trace.TemplateImage, 1024),
		ResultImage:      truncateLogText(input.ResultImage, 1024),
		RequestPayload:   marshalLogPayload(input.RequestPayload),
		ResponsePayload:  marshalLogPayload(input.ResponsePayload),
		Status:           truncateLogText(status, 32),
		ErrorMessage:     truncateLogText(strings.TrimSpace(input.ErrorMessage), 1000),
		CostPoints:       costPoints,
		RefundPoints:     refundPoints,
		DurationMs:       durationMs,
		TokenFingerprint: truncateLogText(tokenFingerprint, 64),
		TokenSlot:        truncateLogText(tokenSlot, 32),
	}

	if err := global.GVA_DB.Create(&logRow).Error; err != nil {
		global.GVA_LOG.Warn("写入模型调用日志失败", zap.Error(err), zap.String("callStage", input.CallStage), zap.String("taskNo", trace.TaskNo))
		return
	}

	saveTryonStatsEvent(&logRow)
}
