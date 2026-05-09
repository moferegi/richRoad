package client

import (
	"bufio"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/flipped-aurora/gin-vue-admin/server/external"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	tryonTaskStatusProcessing = "processing"
	tryonTaskStatusSuccess    = "success"
	tryonTaskStatusFailed     = "failed"
	tryonTaskStatusUnknown    = "unknown"

	tryonRefinerStatusDisabled   = "disabled"
	tryonRefinerStatusPending    = "pending"
	tryonRefinerStatusProcessing = "processing"
	tryonRefinerStatusSuccess    = "success"
	tryonRefinerStatusFailed     = "failed"

	tryonBeautifyStatusDisabled   = "disabled"
	tryonBeautifyStatusProcessing = "processing"
	tryonBeautifyStatusSuccess    = "success"
	tryonBeautifyStatusFailed     = "failed"

	tryonOpConsume         = "tryon_consume"
	tryonOpRefund          = "tryon_refund"
	tryonOpBeautifyConsume = "tryon_beautify_consume"
	tryonOpBeautifyRefund  = "tryon_beautify_refund"
	tryonOpGuestInit       = "tryon_guest_init"
	tryonOpRegisterReward  = "tryon_register_reward"
	tryonOpInviteReward    = "tryon_invite_register_reward"

	tryonReasonTaskDeduct     = "reason_tryonTaskDeduct"
	tryonReasonTaskRefund     = "reason_tryonTaskRefund"
	tryonReasonBeautifyDeduct = "reason_tryonBeautifyDeduct"
	tryonReasonBeautifyRefund = "reason_tryonBeautifyRefund"
	tryonReasonGuestInit      = "reason_tryonGuestInit"
	tryonReasonRegister       = "reason_tryonRegisterReward"
	tryonReasonInviteReward   = "reason_tryonInviteReward"

	aliyunTryonSynthesisURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/image2image/image-synthesis"
	aliyunParsingProcessURL = "https://dashscope.aliyuncs.com/api/v1/services/vision/image-process/process"
	aliyunTryonRefinerModel = "aitryon-refiner"
	aliyunRetouchSkinURL    = "https://facebody.cn-shanghai.aliyuncs.com/"
	aliyunRetouchSkinAction = "RetouchSkin"
	aliyunRetouchSkinVer    = "2019-12-30"

	tryonResultStoreFolder   = "cloth-on/uni-get"
	tryonResultMaxImageBytes = 25 * 1024 * 1024
)

type TryonTaskService struct{}

// DeleteTryonTask 删除单个试衣任务（管理端）
func (s *TryonTaskService) DeleteTryonTask(id uint) error {
	if id == 0 {
		return errors.New("invalidID")
	}

	var task client.TryonTask
	if err := global.GVA_DB.Where("id = ?", id).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tryonTaskNotFound")
		}
		return err
	}

	if err := global.GVA_DB.Delete(&client.TryonTask{}, id).Error; err != nil {
		return err
	}

	s.cleanupLocalTryonMediaFiles([]client.TryonTask{task})
	return nil
}

// DeleteMyTryonTaskByTaskNo 删除当前用户的试衣任务（客户端）
func (s *TryonTaskService) DeleteMyTryonTaskByTaskNo(userID uint, taskNo string) error {
	if userID == 0 {
		return errors.New("loginRequired")
	}

	taskNo = strings.TrimSpace(taskNo)
	if taskNo == "" {
		return errors.New("invalidTaskNo")
	}

	var task client.TryonTask
	if err := global.GVA_DB.Where("user_id = ? AND task_no = ?", userID, taskNo).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tryonTaskNotFound")
		}
		return err
	}

	if err := global.GVA_DB.Delete(&client.TryonTask{}, task.ID).Error; err != nil {
		return err
	}

	s.cleanupLocalTryonMediaFiles([]client.TryonTask{task})
	return nil
}

func (s *TryonTaskService) cleanupLocalTryonMediaFiles(tasks []client.TryonTask) {
	if !strings.EqualFold(strings.TrimSpace(global.GVA_CONFIG.System.OssType), "local") {
		return
	}

	for _, task := range tasks {
		s.removeLocalTryonMediaFile(task.SourceImage)
		s.removeLocalTryonMediaFile(task.TemplateImage)
		s.removeLocalTryonMediaFile(task.ResultImage)
	}
}

func (s *TryonTaskService) removeLocalTryonMediaFile(rawURL string) {
	fileKey := extractLocalTryonFileKey(rawURL)
	if fileKey == "" {
		return
	}

	storeRoot := strings.TrimSpace(global.GVA_CONFIG.Local.StorePath)
	if storeRoot == "" {
		return
	}

	targetPath := filepath.Clean(filepath.Join(storeRoot, filepath.FromSlash(fileKey)))
	rootPath := filepath.Clean(storeRoot)
	targetPathLower := strings.ToLower(targetPath)
	rootPathLower := strings.ToLower(rootPath)
	if targetPathLower != rootPathLower && !strings.HasPrefix(targetPathLower, rootPathLower+string(os.PathSeparator)) {
		return
	}

	if _, err := os.Stat(targetPath); err != nil {
		return
	}

	if err := os.Remove(targetPath); err != nil {
		global.GVA_LOG.Warn("删除本地试衣文件失败", zap.Error(err), zap.String("path", targetPath))
	}
}

func extractLocalTryonFileKey(rawURL string) string {
	value := strings.TrimSpace(rawURL)
	if value == "" {
		return ""
	}

	pathValue := strings.ReplaceAll(value, "\\", "/")
	if parsed, err := url.Parse(value); err == nil && parsed != nil && strings.TrimSpace(parsed.Path) != "" {
		pathValue = strings.ReplaceAll(parsed.Path, "\\", "/")
	}

	lowerPath := strings.ToLower(pathValue)
	if idx := strings.Index(lowerPath, "/uploads/file/"); idx >= 0 {
		pathValue = pathValue[idx+len("/uploads/file/"):]
	} else if strings.HasPrefix(lowerPath, "uploads/file/") {
		pathValue = pathValue[len("uploads/file/"):]
	} else {
		pathValue = strings.TrimPrefix(pathValue, "/")
	}

	pathValue = strings.TrimSpace(strings.TrimPrefix(pathValue, "/"))
	pathValue = strings.TrimPrefix(pathValue, "./")
	if pathValue == "" || strings.Contains(pathValue, "..") {
		return ""
	}

	return pathValue
}

// DeleteTryonTaskByIds 批量删除试衣任务（管理端）
func (s *TryonTaskService) DeleteTryonTaskByIds(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("invalidIDs")
	}

	var tasks []client.TryonTask
	if err := global.GVA_DB.Where("id IN ?", ids).Find(&tasks).Error; err != nil {
		return err
	}
	if len(tasks) == 0 {
		return errors.New("tryonTaskNotFound")
	}

	if err := global.GVA_DB.Where("id IN ?", ids).Delete(&client.TryonTask{}).Error; err != nil {
		return err
	}

	s.cleanupLocalTryonMediaFiles(tasks)
	return nil
}

type TryonTaskTrendItem struct {
	Date       string `json:"date"`
	Total      int64  `json:"total"`
	Processing int64  `json:"processing"`
	Success    int64  `json:"success"`
	Failed     int64  `json:"failed"`
}

type TryonTaskModelStatsItem struct {
	ModelKey            string `json:"modelKey"`
	TaskCount           int64  `json:"taskCount"`
	SuccessCount        int64  `json:"successCount"`
	FailedCount         int64  `json:"failedCount"`
	ProcessingCount     int64  `json:"processingCount"`
	RefinerEnabledCount int64  `json:"refinerEnabledCount"`
	TotalCostPoints     int64  `json:"totalCostPoints"`
}

type TryonTaskStatsData struct {
	Total               int64                     `json:"total"`
	Processing          int64                     `json:"processing"`
	Success             int64                     `json:"success"`
	Failed              int64                     `json:"failed"`
	TotalCostPoints     int64                     `json:"totalCostPoints"`
	RefinerEnabledCount int64                     `json:"refinerEnabledCount"`
	ModelCallTotal      int64                     `json:"modelCallTotal"`
	ModelStats          []TryonTaskModelStatsItem `json:"modelStats"`
}

type tryonModelConfig struct {
	Key                              string   `json:"key"`
	ModelUsage                       string   `json:"modelUsage"`
	BeautifyModelKey                 string   `json:"beautifyModelKey"`
	Enabled                          *bool    `json:"enabled"`
	SupportsRefiner                  *bool    `json:"supportsRefiner"`
	SupportsBeautify                 *bool    `json:"supportsBeautify"`
	SceneType                        string   `json:"sceneType"`
	Scenes                           []string `json:"scenes"`
	Model                            string   `json:"model"`
	Cost                             int      `json:"cost"`
	Provider                         string   `json:"provider"`
	Mode                             string   `json:"mode"`
	URL                              string   `json:"url"`
	Token                            string   `json:"token"`
	ProviderURL                      string   `json:"providerUrl"`
	ProviderToken                    string   `json:"providerToken"`
	TaskQueryURL                     string   `json:"taskQueryUrl"`
	FreeQuotaTotal                   int      `json:"freeQuotaTotal"`
	RefinerFreeQuotaTotal            int      `json:"refinerFreeQuotaTotal"`
	RefinerModel                     string   `json:"refinerModel"`
	RefinerExtraCost                 int      `json:"refinerExtraCost"`
	RefinerExtraPoints               int      `json:"refinerExtraPoints"`
	RefinerURL                       string   `json:"refinerUrl"`
	RefinerToken                     string   `json:"refinerToken"`
	RefinerTaskQueryURL              string   `json:"refinerTaskQueryUrl"`
	BeautifyModel                    string   `json:"beautifyModel"`
	BeautifyExtraCost                int      `json:"beautifyExtraCost"`
	BeautifyExtraPoints              int      `json:"beautifyExtraPoints"`
	BeautifyURL                      string   `json:"beautifyUrl"`
	BeautifyToken                    string   `json:"beautifyToken"`
	BeautifyAccessKeyID              string   `json:"beautifyAccessKeyId"`
	BeautifyAccessKeySec             string   `json:"beautifyAccessKeySecret"`
	BeautifySecurityToken            string   `json:"beautifySecurityToken"`
	BeautifyRetouch                  float64  `json:"beautifyRetouchDegree"`
	BeautifyWhitening                float64  `json:"beautifyWhiteningDegree"`
	RefinerGender                    string   `json:"refinerGender"`
	ApiName                          string   `json:"apiName"`
	GarmentDes                       string   `json:"garmentDes"`
	IsChecked                        *bool    `json:"isChecked"`
	IsCheckedCrop                    *bool    `json:"isCheckedCrop"`
	DenoiseSteps                     int      `json:"denoiseSteps"`
	Seed                             int      `json:"seed"`
	Resolution                       int      `json:"resolution"`
	RestoreFace                      *bool    `json:"restoreFace"`
	ClothesType                      []string `json:"clothesType"`
	AutoEnableAliyunParsingUpperOnly *bool    `json:"autoEnableAliyunParsingUpperOnly"`
	AutoEnableAliyunParsingLowerOnly *bool    `json:"autoEnableAliyunParsingLowerOnly"`
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

func (m *tryonModelConfig) usageValue() string {
	if m == nil {
		return "tryon"
	}
	usage := strings.ToLower(strings.TrimSpace(m.ModelUsage))
	if usage == "beautify" {
		return "beautify"
	}
	return "tryon"
}

func (m *tryonModelConfig) isTryonModelUsage() bool {
	return m == nil || m.usageValue() == "tryon"
}

func (m *tryonModelConfig) isBeautifyModelUsage() bool {
	return m != nil && m.usageValue() == "beautify"
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

func (m *tryonModelConfig) supportsRefinerValue() bool {
	if m == nil {
		return false
	}
	if m.isBeautifyModelUsage() {
		return false
	}
	if m.SupportsRefiner != nil {
		return *m.SupportsRefiner
	}

	modelLower := strings.ToLower(strings.TrimSpace(m.Model))
	if modelLower == "aitryon" || modelLower == "aitryon-plus" {
		return true
	}

	keyLower := strings.ToLower(strings.TrimSpace(m.Key))
	return strings.Contains(keyLower, "aliyun_aitryon") && !strings.Contains(keyLower, "parsing")
}

func (m *tryonModelConfig) refinerModelValue() string {
	if m == nil || strings.TrimSpace(m.RefinerModel) == "" {
		return aliyunTryonRefinerModel
	}
	return strings.TrimSpace(m.RefinerModel)
}

func (m *tryonModelConfig) refinerExtraCostValue() int {
	if m == nil {
		return 0
	}
	if m.RefinerExtraCost > 0 {
		return m.RefinerExtraCost
	}
	if m.RefinerExtraPoints > 0 {
		return m.RefinerExtraPoints
	}
	return 0
}

func (m *tryonModelConfig) supportsBeautifyValue() bool {
	if m == nil {
		return false
	}
	if m.isBeautifyModelUsage() {
		return true
	}
	if m.SupportsBeautify != nil {
		return *m.SupportsBeautify
	}
	if strings.TrimSpace(m.BeautifyModel) != "" || strings.TrimSpace(m.BeautifyURL) != "" || strings.TrimSpace(m.BeautifyToken) != "" {
		return true
	}
	return false
}

func (m *tryonModelConfig) beautifyModelValue() string {
	if m == nil {
		return aliyunRetouchSkinAction
	}
	if strings.TrimSpace(m.BeautifyModel) != "" {
		return strings.TrimSpace(m.BeautifyModel)
	}
	if m.isBeautifyModelUsage() && strings.TrimSpace(m.Model) != "" {
		return strings.TrimSpace(m.Model)
	}
	return aliyunRetouchSkinAction
}

func (m *tryonModelConfig) beautifyExtraCostValue() int {
	if m == nil {
		return 0
	}
	if m.BeautifyExtraCost > 0 {
		return m.BeautifyExtraCost
	}
	if m.BeautifyExtraPoints > 0 {
		return m.BeautifyExtraPoints
	}
	if m.isBeautifyModelUsage() && m.Cost > 0 {
		return m.Cost
	}
	return 0
}

func (m *tryonModelConfig) beautifyEndpointURL() string {
	if m != nil && strings.TrimSpace(m.BeautifyURL) != "" {
		return strings.TrimSpace(m.BeautifyURL)
	}
	if m != nil && m.isBeautifyModelUsage() {
		if strings.TrimSpace(m.URL) != "" {
			return strings.TrimSpace(m.URL)
		}
		if strings.TrimSpace(m.ProviderURL) != "" {
			return strings.TrimSpace(m.ProviderURL)
		}
	}
	return aliyunRetouchSkinURL
}

func (m *tryonModelConfig) beautifyAccessCredentials(fallbackToken string) (accessKeyID string, accessKeySecret string, securityToken string) {
	if m != nil {
		accessKeyID = strings.TrimSpace(m.BeautifyAccessKeyID)
		accessKeySecret = strings.TrimSpace(m.BeautifyAccessKeySec)
		securityToken = strings.TrimSpace(m.BeautifySecurityToken)
		if accessKeyID != "" && accessKeySecret != "" {
			return accessKeyID, accessKeySecret, securityToken
		}
		if id, sec, token, ok := parseAKSKFromToken(m.BeautifyToken); ok {
			if id != "" && sec != "" {
				return id, sec, firstNonEmptyString(token, securityToken)
			}
		}
		if m.isBeautifyModelUsage() {
			if id, sec, token, ok := parseAKSKFromToken(m.Token); ok {
				if id != "" && sec != "" {
					return id, sec, firstNonEmptyString(token, securityToken)
				}
			}
			if id, sec, token, ok := parseAKSKFromToken(m.ProviderToken); ok {
				if id != "" && sec != "" {
					return id, sec, firstNonEmptyString(token, securityToken)
				}
			}
		}
	}

	if id, sec, token, ok := parseAKSKFromToken(fallbackToken); ok {
		return id, sec, token
	}

	return accessKeyID, accessKeySecret, securityToken
}

func (m *tryonModelConfig) beautifyRetouchDegreeValue(override float64) float64 {
	if override > 0 {
		return clampBeautifyDegree(override)
	}
	if m != nil && m.BeautifyRetouch > 0 {
		return clampBeautifyDegree(m.BeautifyRetouch)
	}
	return 70
}

func (m *tryonModelConfig) beautifyWhiteningDegreeValue(override float64) float64 {
	if override > 0 {
		return clampBeautifyDegree(override)
	}
	if m != nil && m.BeautifyWhitening > 0 {
		return clampBeautifyDegree(m.BeautifyWhitening)
	}
	return 30
}

func (m *tryonModelConfig) refinerEndpointURL(fallback string) string {
	if m != nil && strings.TrimSpace(m.RefinerURL) != "" {
		return strings.TrimSpace(m.RefinerURL)
	}
	fallback = strings.TrimSpace(fallback)
	if fallback == "" {
		return aliyunTryonSynthesisURL
	}
	return fallback
}

func (m *tryonModelConfig) refinerAuthToken(fallback string) string {
	if m != nil && strings.TrimSpace(m.RefinerToken) != "" {
		return strings.TrimSpace(m.RefinerToken)
	}
	if m != nil && strings.TrimSpace(m.Token) != "" {
		return strings.TrimSpace(m.Token)
	}
	return strings.TrimSpace(fallback)
}

func (m *tryonModelConfig) refinerTaskQueryURL() string {
	if m != nil && strings.TrimSpace(m.RefinerTaskQueryURL) != "" {
		return strings.TrimSpace(m.RefinerTaskQueryURL)
	}
	return m.queryTaskURL()
}

func (m *tryonModelConfig) refinerGenderValue() string {
	if m == nil || strings.TrimSpace(m.RefinerGender) == "" {
		return "woman"
	}
	gender := strings.ToLower(strings.TrimSpace(m.RefinerGender))
	if gender != "woman" && gender != "man" {
		return "woman"
	}
	return gender
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

func (m *tryonModelConfig) gradioApiName() string {
	if m == nil || strings.TrimSpace(m.ApiName) == "" {
		return "/tryon"
	}
	apiName := strings.TrimSpace(m.ApiName)
	if !strings.HasPrefix(apiName, "/") {
		apiName = "/" + apiName
	}
	return apiName
}

func (m *tryonModelConfig) garmentDescriptionValue() string {
	if m == nil || strings.TrimSpace(m.GarmentDes) == "" {
		return "clothing item"
	}
	return strings.TrimSpace(m.GarmentDes)
}

func (m *tryonModelConfig) gradioIsCheckedValue() bool {
	if m == nil || m.IsChecked == nil {
		return true
	}
	return *m.IsChecked
}

func (m *tryonModelConfig) gradioIsCheckedCropValue() bool {
	if m == nil || m.IsCheckedCrop == nil {
		return false
	}
	return *m.IsCheckedCrop
}

func (m *tryonModelConfig) gradioDenoiseStepsValue() int {
	if m == nil || m.DenoiseSteps <= 0 {
		return 30
	}
	return m.DenoiseSteps
}

func (m *tryonModelConfig) gradioSeedValue() int {
	if m == nil || m.Seed == 0 {
		return 42
	}
	return m.Seed
}

func (m *tryonModelConfig) autoEnableAliyunParsingUpperOnlyValue() bool {
	if m == nil || m.AutoEnableAliyunParsingUpperOnly == nil {
		return true
	}
	return *m.AutoEnableAliyunParsingUpperOnly
}

func (m *tryonModelConfig) autoEnableAliyunParsingLowerOnlyValue() bool {
	if m == nil || m.AutoEnableAliyunParsingLowerOnly == nil {
		return true
	}
	return *m.AutoEnableAliyunParsingLowerOnly
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
		TaskID        string          `json:"task_id"`
		TaskStatus    string          `json:"task_status"`
		ImageURL      json.RawMessage `json:"image_url"`
		ParsingImgURL []string        `json:"parsing_img_url"`
		CropImgURL    []string        `json:"crop_img_url"`
		Code          string          `json:"code"`
		Message       string          `json:"message"`
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

type aliyunRetouchSkinResponse struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
	Data      struct {
		ImageURL string `json:"ImageURL"`
	} `json:"Data"`
}

// CreateTryonTask 创建并处理试衣任务：先扣币，模型失败时全额退币
func (s *TryonTaskService) CreateTryonTask(ctx context.Context, userID uint, req clientReq.CreateTryonTaskReq) (task client.TryonTask, reused bool, err error) {
	if userID == 0 {
		return task, false, errors.New("loginRequired")
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
	if !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return task, false, findErr
	}

	sysConfigService := SysConfigService{}
	baseCostPoints := sysConfigService.GetConfigIntByKey("tryon_cost_points", 1)
	if baseCostPoints < 1 {
		baseCostPoints = 1
	}

	modelCfg, modelErr := s.resolveTryonModelConfig(req.SceneType, req.ModelKey)
	if modelErr != nil {
		return task, false, modelErr
	}
	if modelCfg != nil && modelCfg.Cost > 0 {
		baseCostPoints = modelCfg.Cost
	}

	providerName := "external"
	if modelCfg != nil {
		if strings.TrimSpace(modelCfg.Key) != "" {
			providerName = strings.TrimSpace(modelCfg.Key)
		} else if strings.TrimSpace(modelCfg.Provider) != "" {
			providerName = strings.TrimSpace(modelCfg.Provider)
		}
	}

	enableRefiner := req.EnableRefiner && canEnableAliyunRefiner(req.SceneType, modelCfg)
	extraCostPoints := 0
	if enableRefiner {
		extraCostPoints = modelCfg.refinerExtraCostValue()
	}

	costPoints := baseCostPoints + extraCostPoints
	if costPoints < 1 {
		costPoints = 1
	}

	refinerStatus := tryonRefinerStatusDisabled
	if enableRefiner {
		refinerStatus = tryonRefinerStatusPending
	}

	task = client.TryonTask{
		UserID:         userID,
		RequestID:      requestID,
		TaskNo:         generateTryonTaskNo(),
		SceneType:      req.SceneType,
		Status:         tryonTaskStatusProcessing,
		SourceImage:    req.SourceImage,
		TemplateImage:  req.TemplateImage,
		Provider:       providerName,
		EnableRefiner:  enableRefiner,
		RefinerStatus:  refinerStatus,
		BeautifyStatus: tryonBeautifyStatusDisabled,
		CostPoints:     costPoints,
	}

	if err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx", tx)
		if createErr := tx.Create(&task).Error; createErr != nil {
			return createErr
		}

		pointRecordService := PointRecordService{}
		record := buildPointRecord(userID, client.AssetTypeTryonPoint, "decrease", costPoints, tryonOpConsume, tryonReasonTaskDeduct, "tryon_task:"+task.TaskNo)
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
			return task, false, errors.New("tryonRefundFailedContactAdmin")
		}
		_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
		return task, false, invokeErr
	}

	switch invokeResult.Status {
	case tryonTaskStatusProcessing:
		updates := map[string]interface{}{
			"provider_task_id": strings.TrimSpace(invokeResult.ProviderTaskID),
			"error_message":    "",
		}
		if task.EnableRefiner {
			updates["refiner_status"] = tryonRefinerStatusPending
			updates["refiner_task_id"] = ""
		} else {
			updates["refiner_status"] = tryonRefinerStatusDisabled
			updates["refiner_task_id"] = ""
		}

		updateResult := global.GVA_DB.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ? AND status = ?", task.ID, userID, tryonTaskStatusProcessing).
			Updates(updates)
		if updateResult.Error != nil {
			return task, false, updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return task, false, errors.New("tryonTaskStatusConflictUpdateAsync")
		}
	case tryonTaskStatusSuccess:
		if task.EnableRefiner {
			refinerResult, refinerErr := s.invokeAliyunRefinerAsync(req, strings.TrimSpace(invokeResult.ResultImage), modelCfg)
			if refinerErr != nil {
				refundErr := s.markFailedAndRefund(ctx, &task, refinerErr.Error())
				if refundErr != nil {
					global.GVA_LOG.Error("试衣精修失败退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
					return task, false, errors.New("tryonRefinerRefundFailedContactAdmin")
				}
				_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
				return task, false, refinerErr
			}

			switch refinerResult.Status {
			case tryonTaskStatusProcessing:
				updateResult := global.GVA_DB.Model(&client.TryonTask{}).
					Where("id = ? AND user_id = ? AND status = ?", task.ID, userID, tryonTaskStatusProcessing).
					Updates(map[string]interface{}{
						"provider_task_id": strings.TrimSpace(refinerResult.ProviderTaskID),
						"refiner_task_id":  strings.TrimSpace(refinerResult.ProviderTaskID),
						"refiner_status":   tryonRefinerStatusProcessing,
						"error_message":    "",
					})
				if updateResult.Error != nil {
					return task, false, updateResult.Error
				}
				if updateResult.RowsAffected == 0 {
					return task, false, errors.New("tryonTaskStatusConflictUpdateRefiner")
				}
			case tryonTaskStatusSuccess:
				storedResultImage := s.persistTryonResultImage(strings.TrimSpace(refinerResult.ResultImage), task.TaskNo)
				completeAt := time.Now()
				updateResult := global.GVA_DB.Model(&client.TryonTask{}).
					Where("id = ? AND user_id = ? AND status = ?", task.ID, userID, tryonTaskStatusProcessing).
					Updates(map[string]interface{}{
						"status":           tryonTaskStatusSuccess,
						"provider_task_id": strings.TrimSpace(refinerResult.ProviderTaskID),
						"refiner_task_id":  strings.TrimSpace(refinerResult.ProviderTaskID),
						"refiner_status":   tryonRefinerStatusSuccess,
						"result_image":     storedResultImage,
						"error_message":    "",
						"completed_at":     completeAt,
					})
				if updateResult.Error != nil {
					return task, false, updateResult.Error
				}
				if updateResult.RowsAffected == 0 {
					return task, false, errors.New("tryonTaskStatusConflictComplete")
				}
			case tryonTaskStatusFailed:
				failReason := refinerResult.ErrorMessage
				if strings.TrimSpace(failReason) == "" {
					failReason = "tryonRefinerProcessFail"
				}
				refundErr := s.markFailedAndRefund(ctx, &task, failReason)
				if refundErr != nil {
					global.GVA_LOG.Error("试衣精修失败退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
					return task, false, errors.New("tryonRefinerRefundFailedContactAdmin")
				}
				_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
				return task, false, errors.New(failReason)
			default:
				refundErr := s.markFailedAndRefund(ctx, &task, "tryonRefinerStatusUnknown")
				if refundErr != nil {
					global.GVA_LOG.Error("试衣精修失败退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
					return task, false, errors.New("tryonRefinerRefundFailedContactAdmin")
				}
				_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
				return task, false, errors.New("tryonRefinerStatusUnknown")
			}
		} else {
			storedResultImage := s.persistTryonResultImage(strings.TrimSpace(invokeResult.ResultImage), task.TaskNo)
			completeAt := time.Now()
			updateResult := global.GVA_DB.Model(&client.TryonTask{}).
				Where("id = ? AND user_id = ? AND status = ?", task.ID, userID, tryonTaskStatusProcessing).
				Updates(map[string]interface{}{
					"status":           tryonTaskStatusSuccess,
					"provider_task_id": strings.TrimSpace(invokeResult.ProviderTaskID),
					"refiner_status":   tryonRefinerStatusDisabled,
					"refiner_task_id":  "",
					"result_image":     storedResultImage,
					"error_message":    "",
					"completed_at":     completeAt,
				})
			if updateResult.Error != nil {
				return task, false, updateResult.Error
			}
			if updateResult.RowsAffected == 0 {
				return task, false, errors.New("tryonTaskStatusConflictComplete")
			}
		}
	case tryonTaskStatusFailed:
		failReason := invokeResult.ErrorMessage
		if strings.TrimSpace(failReason) == "" {
			failReason = "tryonModelProcessFail"
		}
		refundErr := s.markFailedAndRefund(ctx, &task, failReason)
		if refundErr != nil {
			global.GVA_LOG.Error("试衣失败退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
			return task, false, errors.New("tryonRefundFailedContactAdmin")
		}
		_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
		return task, false, errors.New(failReason)
	default:
		refundErr := s.markFailedAndRefund(ctx, &task, "tryonModelStatusUnknown")
		if refundErr != nil {
			global.GVA_LOG.Error("试衣失败退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
			return task, false, errors.New("tryonRefundFailedContactAdmin")
		}
		_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
		return task, false, errors.New("tryonModelStatusUnknown")
	}

	_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
	return task, false, nil
}

// ApplyTryonBeautify 对试衣结果执行智能美肤（每个试衣任务仅可执行一次）
func (s *TryonTaskService) ApplyTryonBeautify(ctx context.Context, userID uint, req clientReq.ApplyTryonBeautifyReq) (task client.TryonTask, err error) {
	if userID == 0 {
		return task, errors.New("loginRequired")
	}
	if req.TaskID == 0 {
		return task, errors.New("invalidID")
	}

	var modelCfg *tryonModelConfig
	beautifyCost := 0

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if findErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", req.TaskID, userID).
			First(&task).Error; findErr != nil {
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return errors.New("tryonTaskNotFound")
			}
			return findErr
		}

		if !strings.EqualFold(strings.TrimSpace(task.Status), tryonTaskStatusSuccess) || strings.TrimSpace(task.ResultImage) == "" {
			return errors.New("tryonTaskNotReadyForBeautify")
		}
		if hasTryonBeautifyUsed(task) {
			return errors.New("tryonBeautifyAlreadyUsed")
		}

		resolvedModelCfg, modelErr := s.resolveBeautifyModelConfig(task.SceneType, task.Provider, req.BeautifyModelKey)
		if modelErr != nil {
			return modelErr
		}
		if resolvedModelCfg == nil || !resolvedModelCfg.supportsBeautifyValue() {
			return errors.New("tryonBeautifyUnsupported")
		}

		cfgCopy := *resolvedModelCfg
		modelCfg = &cfgCopy

		beautifyCost = modelCfg.beautifyExtraCostValue()
		if beautifyCost < 0 {
			beautifyCost = 0
		}

		beautifyTaskNo := strings.TrimSpace(task.BeautifyTaskNo)
		if beautifyTaskNo == "" {
			beautifyTaskNo = generateTryonBeautifyTaskNo()
		}

		updates := map[string]interface{}{
			"beautify_status":  tryonBeautifyStatusProcessing,
			"beautify_task_no": beautifyTaskNo,
			"beautify_cost":    beautifyCost,
			"beautify_refund":  0,
			"beautify_error":   "",
			"beautify_result":  "",
			"beautify_at":      nil,
		}
		if updateErr := tx.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ?", task.ID, task.UserID).
			Updates(updates).Error; updateErr != nil {
			return updateErr
		}

		if beautifyCost > 0 {
			pointRecordService := PointRecordService{}
			txCtx := context.WithValue(ctx, "tx", tx)
			record := buildPointRecord(task.UserID, client.AssetTypeTryonPoint, "decrease", beautifyCost, tryonOpBeautifyConsume, tryonReasonBeautifyDeduct, "tryon_beautify:"+beautifyTaskNo)
			if pointErr := pointRecordService.CreatePointRecord(txCtx, record); pointErr != nil {
				return pointErr
			}
		}

		task.BeautifyStatus = tryonBeautifyStatusProcessing
		task.BeautifyTaskNo = beautifyTaskNo
		task.BeautifyCost = beautifyCost
		task.BeautifyRefund = 0
		task.BeautifyError = ""
		task.BeautifyResult = ""
		task.BeautifyAt = nil
		return nil
	})
	if err != nil {
		return task, err
	}

	if modelCfg == nil {
		return task, errors.New("tryonBeautifyUnsupported")
	}

	resultImage, invokeErr := s.invokeTryonBeautifyProvider(task.SceneType, task.ResultImage, modelCfg, req)
	if invokeErr != nil {
		refundErr := s.markTryonBeautifyFailedAndRefund(ctx, &task, beautifyCost, invokeErr.Error())
		if refundErr != nil {
			global.GVA_LOG.Error("智能美肤失败且退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
			return task, errors.New("tryonBeautifyRefundFailedContactAdmin")
		}
		_ = global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, userID).First(&task).Error
		return task, invokeErr
	}

	storedResultImage := s.persistTryonResultImage(strings.TrimSpace(resultImage), task.TaskNo+"_beautify")
	now := time.Now()
	updateResult := global.GVA_DB.Model(&client.TryonTask{}).
		Where("id = ? AND user_id = ? AND beautify_status = ?", task.ID, userID, tryonBeautifyStatusProcessing).
		Updates(map[string]interface{}{
			"beautify_status": tryonBeautifyStatusSuccess,
			"beautify_result": storedResultImage,
			"beautify_error":  "",
			"beautify_at":     now,
		})
	if updateResult.Error != nil {
		return task, updateResult.Error
	}
	if updateResult.RowsAffected == 0 {
		if getErr := global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, userID).First(&task).Error; getErr != nil {
			return task, getErr
		}
		if strings.EqualFold(strings.TrimSpace(task.BeautifyStatus), tryonBeautifyStatusSuccess) {
			return task, nil
		}
		return task, errors.New("tryonBeautifyStatusConflict")
	}

	if err = global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, userID).First(&task).Error; err != nil {
		return task, err
	}
	return task, nil
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
func (s *TryonTaskService) GetTryonTaskStats(info clientReq.TryonTaskSearch) (stats TryonTaskStatsData, err error) {
	baseQuery := s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, false)

	if err = baseQuery.Count(&stats.Total).Error; err != nil {
		return stats, err
	}
	if err = s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, false).Where("status = ?", tryonTaskStatusProcessing).Count(&stats.Processing).Error; err != nil {
		return stats, err
	}
	if err = s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, false).Where("status = ?", tryonTaskStatusSuccess).Count(&stats.Success).Error; err != nil {
		return stats, err
	}
	if err = s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, false).Where("status = ?", tryonTaskStatusFailed).Count(&stats.Failed).Error; err != nil {
		return stats, err
	}

	type tryonTaskStatsSummaryRow struct {
		TotalCostPoints     int64 `json:"totalCostPoints"`
		RefinerEnabledCount int64 `json:"refinerEnabledCount"`
	}
	var summary tryonTaskStatsSummaryRow
	if err = s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, false).
		Select("COALESCE(SUM(cost_points), 0) as total_cost_points, COALESCE(SUM(CASE WHEN enable_refiner THEN 1 ELSE 0 END), 0) as refiner_enabled_count").
		Scan(&summary).Error; err != nil {
		return stats, err
	}
	stats.TotalCostPoints = summary.TotalCostPoints
	stats.RefinerEnabledCount = summary.RefinerEnabledCount

	type tryonTaskModelAggRow struct {
		Provider            string `json:"provider"`
		TaskCount           int64  `json:"taskCount"`
		SuccessCount        int64  `json:"successCount"`
		FailedCount         int64  `json:"failedCount"`
		ProcessingCount     int64  `json:"processingCount"`
		RefinerEnabledCount int64  `json:"refinerEnabledCount"`
		TotalCostPoints     int64  `json:"totalCostPoints"`
	}
	var rows []tryonTaskModelAggRow
	err = s.applyTryonTaskFilters(global.GVA_DB.Model(&client.TryonTask{}), info, false).
		Select("provider, COUNT(*) as task_count, COALESCE(SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END), 0) as success_count, COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0) as failed_count, COALESCE(SUM(CASE WHEN status = 'processing' THEN 1 ELSE 0 END), 0) as processing_count, COALESCE(SUM(CASE WHEN enable_refiner THEN 1 ELSE 0 END), 0) as refiner_enabled_count, COALESCE(SUM(cost_points), 0) as total_cost_points").
		Group("provider").
		Order("task_count DESC").
		Scan(&rows).Error
	if err != nil {
		return stats, err
	}

	stats.ModelStats = make([]TryonTaskModelStatsItem, 0, len(rows))
	for _, row := range rows {
		modelKey := strings.TrimSpace(row.Provider)
		if modelKey == "" {
			modelKey = "default"
		}
		stats.ModelStats = append(stats.ModelStats, TryonTaskModelStatsItem{
			ModelKey:            modelKey,
			TaskCount:           row.TaskCount,
			SuccessCount:        row.SuccessCount,
			FailedCount:         row.FailedCount,
			ProcessingCount:     row.ProcessingCount,
			RefinerEnabledCount: row.RefinerEnabledCount,
			TotalCostPoints:     row.TotalCostPoints,
		})
		stats.ModelCallTotal += row.TaskCount
	}

	return stats, nil
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
			return nil, errors.New("tryonTimeRangeInvalid")
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
		return nil, errors.New("tryonTimeRangeInvalid")
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
		updates := map[string]interface{}{
			"status":        tryonTaskStatusFailed,
			"error_message": failReason,
			"refund_points": task.CostPoints,
			"completed_at":  now,
			"refunded_at":   now,
		}
		if task.EnableRefiner {
			updates["refiner_status"] = tryonRefinerStatusFailed
		}
		updateResult := tx.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ? AND status = ?", task.ID, task.UserID, tryonTaskStatusProcessing).
			Updates(updates)
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return errors.New("tryonTaskStatusConflictRefund")
		}

		pointRecordService := PointRecordService{}
		txCtx := context.WithValue(ctx, "tx", tx)
		record := buildPointRecord(task.UserID, client.AssetTypeTryonPoint, "increase", task.CostPoints, tryonOpRefund, tryonReasonTaskRefund, "tryon_task:"+task.TaskNo)
		if err := pointRecordService.CreatePointRecord(txCtx, record); err != nil {
			return err
		}
		return nil
	})
}

func hasTryonBeautifyUsed(task client.TryonTask) bool {
	status := strings.ToLower(strings.TrimSpace(task.BeautifyStatus))
	if status == "" || status == tryonBeautifyStatusDisabled {
		return strings.TrimSpace(task.BeautifyTaskNo) != "" || strings.TrimSpace(task.BeautifyResult) != ""
	}
	return status == tryonBeautifyStatusProcessing || status == tryonBeautifyStatusSuccess || status == tryonBeautifyStatusFailed
}

func (s *TryonTaskService) markTryonBeautifyFailedAndRefund(ctx context.Context, task *client.TryonTask, costPoints int, failReason string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		failReason = truncateTryonError(failReason)

		refundPoints := costPoints
		if refundPoints < 0 {
			refundPoints = 0
		}

		updateData := map[string]interface{}{
			"beautify_status": tryonBeautifyStatusFailed,
			"beautify_error":  failReason,
			"beautify_refund": refundPoints,
			"beautify_at":     now,
		}

		updateResult := tx.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ? AND beautify_status = ?", task.ID, task.UserID, tryonBeautifyStatusProcessing).
			Updates(updateData)
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return errors.New("tryonBeautifyStatusConflict")
		}

		if refundPoints > 0 {
			pointRecordService := PointRecordService{}
			txCtx := context.WithValue(ctx, "tx", tx)
			remark := "tryon_beautify:" + strings.TrimSpace(task.BeautifyTaskNo)
			record := buildPointRecord(task.UserID, client.AssetTypeTryonPoint, "increase", refundPoints, tryonOpBeautifyRefund, tryonReasonBeautifyRefund, remark)
			if err := pointRecordService.CreatePointRecord(txCtx, record); err != nil {
				return err
			}
		}

		task.BeautifyStatus = tryonBeautifyStatusFailed
		task.BeautifyError = failReason
		task.BeautifyRefund = refundPoints
		task.BeautifyAt = &now
		return nil
	})
}

// GrantRegisterRewardPoints 发放注册奖励试衣币（幂等）
func (s *TryonTaskService) GrantRegisterRewardPoints(ctx context.Context, userID uint) error {
	if userID == 0 {
		return errors.New("invalidUserID")
	}

	rewardPoints := s.getConfigIntAllowZero("tryon_register_reward_points", 8)
	if rewardPoints <= 0 {
		return nil
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx", tx)
		return s.grantPointsIfNotExists(txCtx, tx, userID, tryonOpRegisterReward, rewardPoints, tryonReasonRegister, "register")
	})
}

// GrantInviteRegisterRewardPoints 发放邀请下级注册奖励试衣币（按被邀请用户幂等）
func (s *TryonTaskService) GrantInviteRegisterRewardPoints(ctx context.Context, inviterID uint, invitedUserID uint, invitedUsername string) error {
	if inviterID == 0 || invitedUserID == 0 || inviterID == invitedUserID {
		return nil
	}

	rewardPoints := s.getConfigIntAllowZero("tryon_invite_register_reward_points", 0)
	if rewardPoints <= 0 {
		return nil
	}

	remark := fmt.Sprintf("invite_register:%d", invitedUserID)

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx", tx)
		return s.grantPointsIfNotExists(txCtx, tx, inviterID, tryonOpInviteReward, rewardPoints, tryonReasonInviteReward, remark)
	})
}

func (s *TryonTaskService) loadTryonModelConfigList(sceneType string) ([]tryonModelConfig, error) {
	sysConfigService := SysConfigService{}
	scene := strings.ToLower(strings.TrimSpace(sceneType))
	keys := []string{"tryon_models"}
	if scene == "shoes" {
		keys = []string{"shoe_models", "tryon_models"}
	}

	for _, key := range keys {
		raw, err := sysConfigService.GetConfigByKey(key)
		if err != nil || strings.TrimSpace(raw) == "" {
			continue
		}

		list := make([]tryonModelConfig, 0)
		if unmarshalErr := json.Unmarshal([]byte(raw), &list); unmarshalErr != nil {
			return nil, unmarshalErr
		}
		return list, nil
	}

	return nil, nil
}

func (s *TryonTaskService) resolveTryonModelConfig(sceneType string, modelKey string) (*tryonModelConfig, error) {
	list, err := s.loadTryonModelConfigList(sceneType)
	if err != nil {
		return nil, errors.New("tryonModelConfigInvalid")
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
		if !item.isTryonModelUsage() {
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
		return nil, errors.New("tryonModelUnavailable")
	}

	return fallback, nil
}

func findEnabledBeautifyModelByKey(list []tryonModelConfig, sceneType string, modelKey string) *tryonModelConfig {
	modelKey = strings.TrimSpace(modelKey)
	if modelKey == "" {
		return nil
	}

	for i := range list {
		item := &list[i]
		if !item.isEnabled() || !item.isBeautifyModelUsage() || !item.supportsBeautifyValue() {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(item.Key), modelKey) {
			continue
		}
		if item.matchScene(sceneType) {
			return item
		}
	}

	for i := range list {
		item := &list[i]
		if !item.isEnabled() || !item.isBeautifyModelUsage() || !item.supportsBeautifyValue() {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(item.Key), modelKey) {
			return item
		}
	}

	return nil
}

func (s *TryonTaskService) resolveBeautifyModelConfig(sceneType string, _ string, preferredBeautifyModelKey string) (*tryonModelConfig, error) {
	list, err := s.loadTryonModelConfigList(sceneType)
	if err != nil {
		return nil, errors.New("tryonModelConfigInvalid")
	}
	if len(list) == 0 {
		return nil, nil
	}

	var beautifyFallback *tryonModelConfig
	var beautifyAnyScene *tryonModelConfig

	for i := range list {
		item := &list[i]
		if !item.isEnabled() || !item.isBeautifyModelUsage() || !item.supportsBeautifyValue() {
			continue
		}
		if item.matchScene(sceneType) {
			if beautifyFallback == nil {
				beautifyFallback = item
			}
			continue
		}
		if beautifyAnyScene == nil {
			beautifyAnyScene = item
		}
	}

	if preferred := findEnabledBeautifyModelByKey(list, sceneType, preferredBeautifyModelKey); preferred != nil {
		return preferred, nil
	}

	if beautifyFallback != nil {
		return beautifyFallback, nil
	}
	if beautifyAnyScene != nil {
		return beautifyAnyScene, nil
	}

	return nil, nil
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

	isAliyunParsing := isAliyunParsingModel(providerName, modelName, providerURL, req.SceneType)
	isAliyunTryon := isAliyunTryonModel(providerName, modelName, providerURL, req.SceneType)
	isGradioTryon := isGradioTryonModel(providerName, providerURL, modelName)

	resolvedReq := req
	resolvedReq.TemplatePart = normalizeTryonTemplatePart(req.TemplatePart)
	if isAliyunParsing || isAliyunTryon || isGradioTryon {
		normalizedSource, normalizeErr := s.normalizeAliyunMediaURL(req.SourceImage)
		if normalizeErr != nil {
			return tryonInvokeResult{}, normalizeErr
		}
		resolvedReq.SourceImage = normalizedSource

		if !isAliyunParsing {
			normalizedTemplate, normalizeErr := s.normalizeAliyunMediaURL(req.TemplateImage)
			if normalizeErr != nil {
				return tryonInvokeResult{}, normalizeErr
			}
			resolvedReq.TemplateImage = normalizedTemplate
		}
	}

	if isAliyunParsing {
		return s.invokeAliyunParsing(resolvedReq, modelCfg, providerToken, providerURL, modelName, false)
	}

	if isAliyunTryon && shouldAutoAliyunParsingForSinglePart(resolvedReq.SceneType, modelCfg, modelName, resolvedReq.TemplatePart) {
		parsedTemplate, parseErr := s.tryBuildAliyunParsedTemplate(resolvedReq, providerToken)
		if parseErr != nil {
			global.GVA_LOG.Warn("自动阿里取衣分割失败，回退原始模板图继续试衣",
				zap.String("modelKey", strings.TrimSpace(modelCfg.Key)),
				zap.String("modelName", strings.TrimSpace(modelName)),
				zap.String("templatePart", strings.TrimSpace(resolvedReq.TemplatePart)),
				zap.Error(parseErr),
			)
		} else if strings.TrimSpace(parsedTemplate) != "" {
			resolvedReq.TemplateImage = strings.TrimSpace(parsedTemplate)
		}
	}

	if isAliyunTryon {
		return s.invokeAliyunTryonAsync(resolvedReq, modelCfg, providerToken, providerURL, modelName)
	}

	if isGradioTryon {
		return s.invokeGradioTryon(resolvedReq, modelCfg, providerToken, providerURL)
	}

	if strings.TrimSpace(providerURL) == "" {
		return tryonInvokeResult{}, errors.New("tryonModelServiceNotConfigured")
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
	if strings.TrimSpace(resolvedReq.TemplatePart) != "" {
		body["templatePart"] = strings.TrimSpace(resolvedReq.TemplatePart)
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
		return tryonInvokeResult{}, errors.New("tryonModelResponseEmpty")
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
		return tryonInvokeResult{}, errors.New("tryonModelResultEmpty")
	}

	return tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage}, nil
}

func (s *TryonTaskService) invokeGradioTryon(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, providerToken string, providerURL string) (tryonInvokeResult, error) {
	if strings.TrimSpace(req.SourceImage) == "" {
		return tryonInvokeResult{}, errors.New("modelImageRequired")
	}
	if strings.TrimSpace(req.TemplateImage) == "" {
		return tryonInvokeResult{}, errors.New("clothesImageRequired")
	}

	callURL, rootURL, err := buildGradioCallURL(providerURL, modelCfg.gradioApiName())
	if err != nil {
		return tryonInvokeResult{}, err
	}

	garmentDescription := resolveGradioGarmentDescription(modelCfg.garmentDescriptionValue(), req.TemplatePart)

	body := map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"background": gradioFileData(req.SourceImage),
				"layers":     []interface{}{},
				"composite":  nil,
			},
			gradioFileData(req.TemplateImage),
			garmentDescription,
			modelCfg.gradioIsCheckedValue(),
			modelCfg.gradioIsCheckedCropValue(),
			modelCfg.gradioDenoiseStepsValue(),
			modelCfg.gradioSeedValue(),
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return tryonInvokeResult{}, err
	}

	client := &http.Client{Timeout: 10 * time.Minute}
	request, err := http.NewRequest(http.MethodPost, callURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return tryonInvokeResult{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(providerToken) != "" {
		request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(providerToken))
	}

	response, err := client.Do(request)
	if err != nil {
		return tryonInvokeResult{}, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return tryonInvokeResult{}, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return tryonInvokeResult{}, errors.New("gradioTryonRequestFailed")
	}

	var createResp struct {
		EventID string          `json:"event_id"`
		Data    json.RawMessage `json:"data"`
		Error   string          `json:"error"`
	}
	if err := json.Unmarshal(responseBody, &createResp); err != nil {
		return tryonInvokeResult{}, errors.New("gradioTryonCreateRespParseFail")
	}
	if strings.TrimSpace(createResp.Error) != "" {
		return tryonInvokeResult{}, errors.New(strings.TrimSpace(createResp.Error))
	}
	if len(createResp.Data) > 0 && string(createResp.Data) != "null" {
		resultImage := extractGradioOutputImage(createResp.Data, rootURL)
		if resultImage != "" {
			return tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage}, nil
		}
	}

	eventID := strings.TrimSpace(createResp.EventID)
	if eventID == "" {
		return tryonInvokeResult{}, errors.New("gradioEventIDMissing")
	}

	resultImage, err := s.waitGradioResult(client, callURL, rootURL, eventID, providerToken)
	if err != nil {
		return tryonInvokeResult{}, err
	}
	if resultImage == "" {
		return tryonInvokeResult{}, errors.New("gradioResultEmpty")
	}

	return tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage, ProviderTaskID: eventID}, nil
}

func (s *TryonTaskService) waitGradioResult(client *http.Client, callURL string, rootURL string, eventID string, providerToken string) (string, error) {
	eventURL := strings.TrimRight(callURL, "/") + "/" + url.PathEscape(eventID)
	request, err := http.NewRequest(http.MethodGet, eventURL, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Accept", "text/event-stream")
	if strings.TrimSpace(providerToken) != "" {
		request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(providerToken))
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", errors.New("gradioTryonQueryFailed")
	}

	scanner := bufio.NewScanner(response.Body)
	scanner.Buffer(make([]byte, 1024), 10*1024*1024)
	currentEvent := ""
	lastError := ""
	sawErrorEvent := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "event:") {
			currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if strings.EqualFold(currentEvent, "error") {
			sawErrorEvent = true
			if payload == "" || payload == "null" || payload == "{}" {
				if strings.TrimSpace(lastError) == "" {
					lastError = "gradioTryonErrorWithoutDetails"
				}
				continue
			}
			lastError = parseGradioErrorMessage([]byte(payload))
			continue
		}

		if payload == "" || payload == "null" {
			continue
		}

		resultImage := extractGradioOutputImage([]byte(payload), rootURL)
		if resultImage != "" {
			return resultImage, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if lastError != "" {
		return "", errors.New(lastError)
	}
	if sawErrorEvent {
		return "", errors.New("gradioTryonErrorWithoutDetails")
	}
	return "", nil
}

func isGradioTryonModel(providerName string, providerURL string, modelName string) bool {
	provider := strings.ToLower(strings.TrimSpace(providerName))
	if provider == "gradio" || provider == "huggingface" || provider == "hf_space" || provider == "hf-space" {
		return true
	}
	joined := strings.ToLower(strings.TrimSpace(providerURL) + " " + strings.TrimSpace(modelName))
	return strings.Contains(joined, "hf.space") || strings.Contains(joined, "gradio") || strings.Contains(joined, "idm-vton")
}

func buildGradioCallURL(rawURL string, apiName string) (string, string, error) {
	baseURL := normalizeGradioBaseURL(rawURL)
	if baseURL == "" {
		return "", "", errors.New("gradioModelURLNotConfigured")
	}

	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", "", errors.New("gradioModelURLInvalid")
	}

	apiPath := strings.Trim(strings.TrimSpace(apiName), "/")
	if apiPath == "" {
		apiPath = "tryon"
	}

	rootURL := parsed.Scheme + "://" + parsed.Host
	pathValue := strings.TrimRight(parsed.EscapedPath(), "/")
	if strings.Contains(pathValue, "/call/") {
		return strings.TrimRight(baseURL, "/"), rootURL, nil
	}
	if strings.Contains(pathValue, "/gradio_api/call/") {
		return strings.TrimRight(baseURL, "/"), rootURL, nil
	}
	if pathValue == "/call" {
		return strings.TrimRight(baseURL, "/") + "/" + apiPath, rootURL, nil
	}
	if strings.HasSuffix(pathValue, "/gradio_api") {
		return strings.TrimRight(baseURL, "/") + "/call/" + apiPath, rootURL, nil
	}
	return strings.TrimRight(rootURL+pathValue, "/") + "/call/" + apiPath, rootURL, nil
}

func normalizeGradioBaseURL(rawURL string) string {
	value := strings.TrimSpace(rawURL)
	if value == "" {
		return ""
	}
	if !strings.Contains(value, "://") {
		if strings.Contains(value, "/") && !strings.Contains(value, ".") {
			value = strings.ToLower(strings.ReplaceAll(value, "/", "-")) + ".hf.space"
		}
		value = "https://" + value
	}
	return strings.TrimRight(value, "/")
}

func gradioFileData(rawURL string) map[string]interface{} {
	value := strings.TrimSpace(rawURL)
	origName := "image.png"
	if parsed, err := url.Parse(value); err == nil && parsed != nil {
		baseName := strings.TrimSpace(path.Base(parsed.Path))
		if baseName != "" && baseName != "." && baseName != "/" {
			origName = baseName
		}
	}

	return map[string]interface{}{
		"path":      value,
		"url":       value,
		"orig_name": origName,
		"meta": map[string]interface{}{
			"_type": "gradio.FileData",
		},
	}
}

func parseGradioErrorMessage(payload []byte) string {
	var errObj map[string]interface{}
	if json.Unmarshal(payload, &errObj) == nil {
		for _, key := range []string{"error", "message", "detail"} {
			if msg := strings.TrimSpace(fmt.Sprint(errObj[key])); msg != "" && msg != "<nil>" {
				return msg
			}
		}
	}
	msg := strings.TrimSpace(string(payload))
	if msg == "" {
		return "gradioTryonTaskFailed"
	}
	return msg
}

func extractGradioOutputImage(payload []byte, rootURL string) string {
	var value interface{}
	if err := json.Unmarshal(payload, &value); err != nil {
		return ""
	}
	return extractGradioImageValue(value, rootURL, true)
}

func extractGradioImageValue(value interface{}, rootURL string, preferFirst bool) string {
	switch typed := value.(type) {
	case []interface{}:
		if len(typed) == 0 {
			return ""
		}
		if preferFirst {
			if result := extractGradioImageValue(typed[0], rootURL, false); result != "" {
				return result
			}
		}
		for _, item := range typed {
			if result := extractGradioImageValue(item, rootURL, false); result != "" {
				return result
			}
		}
	case map[string]interface{}:
		for _, key := range []string{"data", "output", "result", "value", "image", "images", "file"} {
			if nested, ok := typed[key]; ok {
				if result := extractGradioImageValue(nested, rootURL, true); result != "" {
					return result
				}
			}
		}
		for _, key := range []string{"url", "path"} {
			if result := formatGradioFileURL(strings.TrimSpace(fmt.Sprint(typed[key])), rootURL); result != "" {
				return result
			}
		}
	case string:
		return formatGradioFileURL(typed, rootURL)
	}
	return ""
}

func formatGradioFileURL(value string, rootURL string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "<nil>" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(value), "http://") || strings.HasPrefix(strings.ToLower(value), "https://") {
		return value
	}
	rootURL = strings.TrimRight(strings.TrimSpace(rootURL), "/")
	if rootURL == "" {
		return value
	}
	if strings.HasPrefix(value, "/file=") {
		return rootURL + value
	}
	if strings.HasPrefix(value, "/tmp/") {
		return rootURL + "/file=" + value
	}
	if strings.HasPrefix(value, "/") {
		return rootURL + value
	}
	return rootURL + "/file=" + value
}

func normalizeTryonTemplatePart(value string) string {
	part := strings.ToLower(strings.TrimSpace(value))
	if part == "upper" || part == "lower" {
		return part
	}
	return ""
}

func resolveAliyunParsingClothesTypeForTemplatePart(templatePart string) []string {
	part := normalizeTryonTemplatePart(templatePart)
	if part == "upper" {
		// For single upper-garment try-on assist, use lower mask extraction to keep lower body unchanged.
		return []string{"lower"}
	}
	if part == "lower" {
		// For single lower-garment try-on assist, use upper mask extraction to keep upper body unchanged.
		return []string{"upper"}
	}
	return nil
}

func isAliyunParsingNoUsableGarmentAreaError(raw string) bool {
	text := strings.ToLower(strings.TrimSpace(raw))
	if text == "" {
		return false
	}
	if strings.Contains(text, "aliyunparsingnousablegarmentarea") {
		return true
	}
	if strings.Contains(text, "no usable garment area") {
		return true
	}
	if strings.Contains(text, "no usable") && strings.Contains(text, "garment") {
		return true
	}
	return false
}

func sameAliyunClothesType(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !strings.EqualFold(strings.TrimSpace(a[i]), strings.TrimSpace(b[i])) {
			return false
		}
	}
	return true
}

func shouldAutoAliyunParsingForSinglePart(sceneType string, modelCfg *tryonModelConfig, modelName string, templatePart string) bool {
	if modelCfg == nil || !strings.EqualFold(strings.TrimSpace(sceneType), "clothes") {
		return false
	}

	part := normalizeTryonTemplatePart(templatePart)
	if part == "" {
		return false
	}

	modelLower := strings.ToLower(strings.TrimSpace(modelName))
	if modelLower != "aitryon" && modelLower != "aitryon-plus" {
		return false
	}

	if part == "upper" {
		return modelCfg.autoEnableAliyunParsingUpperOnlyValue()
	}
	if part == "lower" {
		return modelCfg.autoEnableAliyunParsingLowerOnlyValue()
	}

	return false
}

func (s *TryonTaskService) resolveAliyunParsingAssistModelConfig(sceneType string) (*tryonModelConfig, error) {
	list, err := s.loadTryonModelConfigList(sceneType)
	if err != nil {
		return nil, errors.New("tryonModelConfigInvalid")
	}
	if len(list) == 0 {
		return nil, nil
	}

	var enabledSceneMatched *tryonModelConfig
	var enabledFallback *tryonModelConfig
	var configuredSceneMatched *tryonModelConfig
	var configuredFallback *tryonModelConfig

	for i := range list {
		item := &list[i]
		if !item.isTryonModelUsage() {
			continue
		}

		// Use strict parsing model detection here. Passing sceneType=takeoff would
		// make isAliyunParsingModel return true for all models, which can wrongly
		// select normal try-on models (e.g. aitryon) for parsing assist.
		itemModelName := resolveAliyunModelName(item, "")
		if !isAliyunParsingModel(item.Provider, itemModelName, item.endpointURL(), "") {
			continue
		}

		if item.isEnabled() {
			if item.matchScene(sceneType) && enabledSceneMatched == nil {
				enabledSceneMatched = item
			}
			if enabledFallback == nil {
				enabledFallback = item
			}
		}

		if item.matchScene(sceneType) && configuredSceneMatched == nil {
			configuredSceneMatched = item
		}
		if configuredFallback == nil {
			configuredFallback = item
		}
	}

	if enabledSceneMatched != nil {
		return enabledSceneMatched, nil
	}
	if enabledFallback != nil {
		return enabledFallback, nil
	}
	if configuredSceneMatched != nil {
		return configuredSceneMatched, nil
	}
	return configuredFallback, nil
}

func (s *TryonTaskService) tryBuildAliyunParsedTemplate(req clientReq.CreateTryonTaskReq, fallbackToken string) (string, error) {
	part := normalizeTryonTemplatePart(req.TemplatePart)
	if part == "" {
		return "", nil
	}

	parsingCfg, err := s.resolveAliyunParsingAssistModelConfig(req.SceneType)
	if err != nil {
		return "", err
	}
	if parsingCfg == nil {
		return "", errors.New("aliyunParsingModelUnavailable")
	}

	cfgCopy := *parsingCfg
	parseClothesType := resolveAliyunParsingClothesTypeForTemplatePart(part)
	if len(parseClothesType) == 0 {
		parseClothesType = []string{part}
	}

	attemptClothesTypes := [][]string{parseClothesType}
	fallbackClothesType := []string{part}
	if !sameAliyunClothesType(parseClothesType, fallbackClothesType) {
		attemptClothesTypes = append(attemptClothesTypes, fallbackClothesType)
	}

	parseReq := clientReq.CreateTryonTaskReq{
		SceneType:   "takeoff",
		SourceImage: strings.TrimSpace(req.TemplateImage),
	}

	parseProviderURL := strings.TrimSpace(cfgCopy.endpointURL())
	if parseProviderURL == "" {
		parseProviderURL = aliyunParsingProcessURL
	}

	parseProviderToken := strings.TrimSpace(cfgCopy.authToken())
	if parseProviderToken == "" {
		parseProviderToken = strings.TrimSpace(fallbackToken)
	}

	parseModelName := resolveAliyunModelName(&cfgCopy, "takeoff")
	lastErr := errors.New("aliyunParsingNoUsableGarmentArea")

	for i, clothesTypeAttempt := range attemptClothesTypes {
		cfgCopy.ClothesType = clothesTypeAttempt
		global.GVA_LOG.Info("自动阿里取衣分割参数",
			zap.String("templatePart", part),
			zap.Strings("clothesType", cfgCopy.ClothesType),
			zap.String("modelKey", strings.TrimSpace(cfgCopy.Key)),
			zap.Int("attempt", i+1),
			zap.Int("attemptTotal", len(attemptClothesTypes)),
		)

		parseResult, parseErr := s.invokeAliyunParsing(parseReq, &cfgCopy, parseProviderToken, parseProviderURL, parseModelName, true)
		if parseErr != nil {
			return "", parseErr
		}

		if strings.EqualFold(strings.TrimSpace(parseResult.Status), tryonTaskStatusSuccess) && strings.TrimSpace(parseResult.ResultImage) != "" {
			normalizedTemplate, normalizeErr := s.normalizeAliyunMediaURL(strings.TrimSpace(parseResult.ResultImage))
			if normalizeErr != nil {
				return "", normalizeErr
			}
			return strings.TrimSpace(normalizedTemplate), nil
		}

		errMsg := strings.TrimSpace(parseResult.ErrorMessage)
		if errMsg == "" {
			errMsg = "aliyunParsingNoUsableGarmentArea"
		}
		lastErr = errors.New(errMsg)

		if i < len(attemptClothesTypes)-1 && isAliyunParsingNoUsableGarmentAreaError(errMsg) {
			global.GVA_LOG.Warn("自动阿里取衣分割无可用区域，尝试备用clothesType",
				zap.String("templatePart", part),
				zap.String("modelKey", strings.TrimSpace(cfgCopy.Key)),
				zap.Strings("failedClothesType", clothesTypeAttempt),
				zap.String("error", errMsg),
			)
			continue
		}

		return "", lastErr
	}

	return "", lastErr
}

func resolveGradioGarmentDescription(configured string, templatePart string) string {
	base := strings.TrimSpace(configured)
	if base == "" {
		base = "clothing item"
	}

	part := normalizeTryonTemplatePart(templatePart)
	if part == "" {
		return base
	}

	lowerBase := strings.ToLower(base)
	if part == "lower" {
		if lowerBase == "clothing item" || lowerBase == "upper garment" {
			return "lower garment"
		}
		return base
	}

	if part == "upper" {
		if lowerBase == "clothing item" || lowerBase == "lower garment" {
			return "upper garment"
		}
	}

	return base
}

func buildAliyunTryonInput(personImageURL string, templateImageURL string, templatePart string) map[string]interface{} {
	input := map[string]interface{}{
		"person_image_url": strings.TrimSpace(personImageURL),
	}

	part := normalizeTryonTemplatePart(templatePart)
	if part == "lower" {
		input["bottom_garment_url"] = strings.TrimSpace(templateImageURL)
		return input
	}

	input["top_garment_url"] = strings.TrimSpace(templateImageURL)
	return input
}

func (s *TryonTaskService) invokeAliyunTryonAsync(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, providerToken string, providerURL string, modelName string) (tryonInvokeResult, error) {
	if strings.TrimSpace(providerToken) == "" {
		return tryonInvokeResult{}, errors.New("aliyunTryonApiKeyMissing")
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
		return tryonInvokeResult{}, errors.New("modelImageRequired")
	}
	if strings.TrimSpace(req.TemplateImage) == "" {
		return tryonInvokeResult{}, errors.New("clothesImageRequired")
	}

	body := map[string]interface{}{
		"model": resolvedModelName,
		"input": buildAliyunTryonInput(req.SourceImage, req.TemplateImage, req.TemplatePart),
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
		return tryonInvokeResult{}, errors.New("aliyunTryonResponseEmpty")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code)}, nil
	}

	taskID := strings.TrimSpace(resp.Output.TaskID)
	if taskID == "" {
		return tryonInvokeResult{}, errors.New("aliyunTryonTaskIDMissing")
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

func (s *TryonTaskService) invokeAliyunRefinerAsync(req clientReq.CreateTryonTaskReq, coarseImageURL string, modelCfg *tryonModelConfig) (tryonInvokeResult, error) {
	if strings.TrimSpace(coarseImageURL) == "" {
		return tryonInvokeResult{}, errors.New("refinerInputImageRequired")
	}
	if strings.TrimSpace(req.SourceImage) == "" {
		return tryonInvokeResult{}, errors.New("modelImageRequired")
	}
	if strings.TrimSpace(req.TemplateImage) == "" {
		return tryonInvokeResult{}, errors.New("clothesImageRequired")
	}

	normalizedSource, normalizeErr := s.normalizeAliyunMediaURL(req.SourceImage)
	if normalizeErr != nil {
		return tryonInvokeResult{}, normalizeErr
	}
	normalizedTemplate, normalizeErr := s.normalizeAliyunMediaURL(req.TemplateImage)
	if normalizeErr != nil {
		return tryonInvokeResult{}, normalizeErr
	}
	normalizedCoarse, normalizeErr := s.normalizeAliyunMediaURL(coarseImageURL)
	if normalizeErr != nil {
		return tryonInvokeResult{}, normalizeErr
	}

	sysConfigService := SysConfigService{}
	providerURL, _ := sysConfigService.GetConfigByKey("tryon_provider_url")
	providerToken, _ := sysConfigService.GetConfigByKey("tryon_provider_token")

	refinerURL := modelCfg.refinerEndpointURL(providerURL)
	refinerToken := modelCfg.refinerAuthToken(providerToken)
	if strings.TrimSpace(refinerToken) == "" {
		return tryonInvokeResult{}, errors.New("aliyunRefinerApiKeyMissing")
	}

	refinerURLHost := urlHostname(refinerURL)
	sourceHost := urlHostname(normalizedSource)
	templateHost := urlHostname(normalizedTemplate)
	coarseHost := urlHostname(normalizedCoarse)
	global.GVA_LOG.Info("阿里精修调用诊断",
		zap.String("sceneType", strings.TrimSpace(req.SceneType)),
		zap.String("refinerURLHost", refinerURLHost),
		zap.String("sourceHost", sourceHost),
		zap.String("templateHost", templateHost),
		zap.String("coarseHost", coarseHost),
	)

	input := buildAliyunTryonInput(normalizedSource, normalizedTemplate, req.TemplatePart)
	input["coarse_image_url"] = strings.TrimSpace(normalizedCoarse)

	body := map[string]interface{}{
		"model": modelCfg.refinerModelValue(),
		"input": input,
		"parameters": map[string]interface{}{
			"gender": modelCfg.refinerGenderValue(),
		},
	}

	resp, err := external.HttpRequest[dashscopeTryonCreateResponse](external.RequestParams{
		URL:    refinerURL,
		Method: "POST",
		Headers: map[string]string{
			"Authorization":     "Bearer " + strings.TrimSpace(refinerToken),
			"X-DashScope-Async": "enable",
		},
		Body: body,
	})
	if err != nil {
		global.GVA_LOG.Warn("阿里精修HTTP请求失败",
			zap.String("refinerURLHost", refinerURLHost),
			zap.Error(err),
		)
		return tryonInvokeResult{}, err
	}
	if resp == nil {
		global.GVA_LOG.Warn("阿里精修响应为空",
			zap.String("refinerURLHost", refinerURLHost),
		)
		return tryonInvokeResult{}, errors.New("aliyunRefinerResponseEmpty")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		mappedErr := formatAliyunModelError(resp.Message, code)
		global.GVA_LOG.Warn("阿里精修返回业务失败",
			zap.String("refinerURLHost", refinerURLHost),
			zap.String("aliyunCode", code),
			zap.String("aliyunTaskStatus", strings.TrimSpace(resp.Output.TaskStatus)),
			zap.String("mappedError", mappedErr),
		)
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: mappedErr}, nil
	}

	taskID := strings.TrimSpace(resp.Output.TaskID)
	if taskID == "" {
		return tryonInvokeResult{}, errors.New("aliyunRefinerTaskIDMissing")
	}

	status := normalizeProviderTaskStatus(resp.Output.TaskStatus)
	if status == tryonTaskStatusFailed {
		mappedErr := formatAliyunModelError(resp.Message, resp.Code)
		global.GVA_LOG.Warn("阿里精修任务状态失败",
			zap.String("refinerURLHost", refinerURLHost),
			zap.String("aliyunTaskStatus", strings.TrimSpace(resp.Output.TaskStatus)),
			zap.String("mappedError", mappedErr),
		)
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ProviderTaskID: taskID, ErrorMessage: mappedErr}, nil
	}
	if status == tryonTaskStatusSuccess {
		queryResult, queryErr := s.queryAliyunTryonTask(taskID, strings.TrimSpace(refinerToken), refinerURL, modelCfg.refinerTaskQueryURL())
		if queryErr != nil {
			return tryonInvokeResult{}, queryErr
		}
		queryResult.ProviderTaskID = taskID
		return queryResult, nil
	}

	return tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}, nil
}

func (s *TryonTaskService) invokeAliyunParsing(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, providerToken string, providerURL string, modelName string, preferParsingImage bool) (tryonInvokeResult, error) {
	if strings.TrimSpace(providerToken) == "" {
		return tryonInvokeResult{}, errors.New("aliyunParsingApiKeyMissing")
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
		if shouldFallbackToAliyunParsingAsync(err.Error()) {
			global.GVA_LOG.Info("阿里取衣分割同步请求异常，自动切换异步",
				zap.String("model", resolvedModelName),
				zap.String("processURLHost", urlHostname(processURL)),
				zap.Error(err),
			)
			return s.invokeAliyunParsingAsync(req, modelCfg, providerToken, processURL, resolvedModelName, clothesType, preferParsingImage)
		}
		return tryonInvokeResult{}, err
	}
	if resp == nil {
		return tryonInvokeResult{}, errors.New("aliyunParsingResponseEmpty")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		mappedErr := formatAliyunModelError(resp.Message, code)
		if shouldFallbackToAliyunParsingAsync(code, resp.Message, mappedErr) {
			global.GVA_LOG.Info("阿里取衣分割同步调用受限，自动切换异步",
				zap.String("model", resolvedModelName),
				zap.String("processURLHost", urlHostname(processURL)),
				zap.String("aliyunCode", code),
				zap.String("aliyunMessage", strings.TrimSpace(resp.Message)),
			)
			return s.invokeAliyunParsingAsync(req, modelCfg, providerToken, processURL, resolvedModelName, clothesType, preferParsingImage)
		}
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: mappedErr}, nil
	}

	resultImage := ""
	candidateGroups := [][]string{resp.Output.CropImgURL, resp.Output.ParsingImgURL}
	if preferParsingImage {
		candidateGroups = [][]string{resp.Output.ParsingImgURL, resp.Output.CropImgURL}
	}

	for _, group := range candidateGroups {
		for _, item := range group {
			resultImage = strings.TrimSpace(item)
			if resultImage != "" {
				break
			}
		}
		if resultImage != "" {
			break
		}
	}

	if resultImage == "" {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: "aliyunParsingNoUsableGarmentArea"}, nil
	}

	return tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage}, nil
}

func shouldFallbackToAliyunParsingAsync(values ...string) bool {
	for _, value := range values {
		text := strings.ToLower(strings.TrimSpace(value))
		if text == "" {
			continue
		}
		if strings.Contains(text, "does not support synchronous calls") {
			return true
		}
		if strings.Contains(text, "not support synchronous") {
			return true
		}
		if strings.Contains(text, "synchronous calls are not supported") {
			return true
		}
		if strings.Contains(text, "不支持同步") || strings.Contains(text, "仅支持异步") {
			return true
		}
	}
	return false
}

func pickAliyunParsingResultImage(preferParsingImage bool, parsingURLs []string, cropURLs []string, fallbackImageURL json.RawMessage) string {
	candidateGroups := [][]string{cropURLs, parsingURLs}
	if preferParsingImage {
		candidateGroups = [][]string{parsingURLs, cropURLs}
	}

	for _, group := range candidateGroups {
		for _, item := range group {
			value := strings.TrimSpace(item)
			if value != "" {
				return value
			}
		}
	}

	return extractDashscopeImageURL(fallbackImageURL)
}

func (s *TryonTaskService) invokeAliyunParsingAsync(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, providerToken string, processURL string, modelName string, clothesType []string, preferParsingImage bool) (tryonInvokeResult, error) {
	body := map[string]interface{}{
		"model": modelName,
		"input": map[string]interface{}{
			"image_url": strings.TrimSpace(req.SourceImage),
		},
		"parameters": map[string]interface{}{
			"clothes_type": clothesType,
		},
	}

	resp, err := external.HttpRequest[dashscopeTryonCreateResponse](external.RequestParams{
		URL:    processURL,
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
		return tryonInvokeResult{}, errors.New("aliyunParsingResponseEmpty")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code)}, nil
	}

	taskID := strings.TrimSpace(resp.Output.TaskID)
	if taskID == "" {
		return tryonInvokeResult{}, errors.New("aliyunTryonTaskIDMissing")
	}

	status := normalizeProviderTaskStatus(resp.Output.TaskStatus)
	if status == tryonTaskStatusFailed {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ProviderTaskID: taskID, ErrorMessage: formatAliyunModelError(resp.Message, resp.Code)}, nil
	}

	taskQueryURL := ""
	if modelCfg != nil {
		taskQueryURL = modelCfg.queryTaskURL()
	}

	const maxAttempts = 20
	const interval = 1500 * time.Millisecond

	for attempt := 0; attempt < maxAttempts; attempt++ {
		result, queryErr := s.queryAliyunParsingTask(taskID, providerToken, processURL, taskQueryURL, preferParsingImage)
		if queryErr != nil {
			return tryonInvokeResult{}, queryErr
		}
		if result.Status == tryonTaskStatusProcessing || result.Status == tryonTaskStatusUnknown {
			if attempt == maxAttempts-1 {
				break
			}
			time.Sleep(interval)
			continue
		}
		return result, nil
	}

	return tryonInvokeResult{Status: tryonTaskStatusFailed, ProviderTaskID: taskID, ErrorMessage: "tryonTaskFailed"}, nil
}

func (s *TryonTaskService) queryAliyunParsingTask(taskID string, providerToken string, createURL string, taskQueryURL string, preferParsingImage bool) (tryonInvokeResult, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return tryonInvokeResult{}, errors.New("aliyunTaskIDRequired")
	}
	if strings.TrimSpace(providerToken) == "" {
		return tryonInvokeResult{}, errors.New("aliyunParsingApiKeyMissing")
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
		return tryonInvokeResult{}, errors.New("aliyunTaskQueryResponseEmpty")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code), ProviderTaskID: taskID}, nil
	}

	status := normalizeProviderTaskStatus(resp.Output.TaskStatus)
	switch status {
	case tryonTaskStatusSuccess:
		resultImage := pickAliyunParsingResultImage(preferParsingImage, resp.Output.ParsingImgURL, resp.Output.CropImgURL, resp.Output.ImageURL)
		if strings.TrimSpace(resultImage) == "" {
			return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: "aliyunTaskSuccessNoImage", ProviderTaskID: taskID}, nil
		}
		return tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: strings.TrimSpace(resultImage), ProviderTaskID: taskID}, nil
	case tryonTaskStatusFailed:
		return tryonInvokeResult{
			Status:         tryonTaskStatusFailed,
			ProviderTaskID: taskID,
			ErrorMessage: formatAliyunModelError(
				resp.Output.Message,
				resp.Output.Code,
				resp.Message,
				resp.Code,
				"tryonTaskFailed",
			),
		}, nil
	case tryonTaskStatusProcessing, tryonTaskStatusUnknown:
		return tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}, nil
	default:
		return tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}, nil
	}
}

func (s *TryonTaskService) invokeTryonBeautifyProvider(sceneType string, resultImage string, modelCfg *tryonModelConfig, req clientReq.ApplyTryonBeautifyReq) (string, error) {
	if modelCfg == nil || !modelCfg.isBeautifyModelUsage() {
		return "", errors.New("tryonBeautifyUnsupported")
	}

	providerMode := strings.TrimSpace(modelCfg.Mode)
	if providerMode == "" {
		providerMode = "prod"
	}
	providerURL := strings.TrimSpace(modelCfg.endpointURL())
	providerToken := strings.TrimSpace(modelCfg.authToken())

	if strings.EqualFold(strings.TrimSpace(providerMode), "mock_success") {
		return strings.TrimSpace(resultImage), nil
	}

	modelName := ""
	providerName := ""
	modelKey := ""
	modelName = modelCfg.beautifyModelValue()
	providerName = strings.TrimSpace(modelCfg.Provider)
	modelKey = strings.TrimSpace(modelCfg.Key)

	if isAliyunBeautifyModel(providerName, modelName, providerURL) {
		return s.invokeAliyunRetouchSkin(resultImage, modelCfg, providerToken, req)
	}

	normalizedImage, normalizeErr := s.normalizeAliyunMediaURL(resultImage)
	if normalizeErr != nil {
		return "", normalizeErr
	}

	targetURL := strings.TrimSpace(providerURL)
	if targetURL == "" {
		targetURL = strings.TrimSpace(modelCfg.beautifyEndpointURL())
	}
	if strings.TrimSpace(targetURL) == "" {
		return "", errors.New("tryonModelServiceNotConfigured")
	}

	headers := map[string]string{}
	if strings.TrimSpace(providerToken) != "" {
		headers["Authorization"] = "Bearer " + strings.TrimSpace(providerToken)
	}

	retouchDegree := modelCfg.beautifyRetouchDegreeValue(req.RetouchDegree)
	whiteningDegree := modelCfg.beautifyWhiteningDegreeValue(req.WhiteningDegree)

	body := map[string]interface{}{
		"operationType":   "beautify",
		"sceneType":       strings.TrimSpace(sceneType),
		"imageUrl":        strings.TrimSpace(normalizedImage),
		"sourceImage":     strings.TrimSpace(normalizedImage),
		"resultImage":     strings.TrimSpace(normalizedImage),
		"retouchDegree":   retouchDegree,
		"whiteningDegree": whiteningDegree,
	}
	if strings.TrimSpace(modelName) != "" {
		body["model"] = strings.TrimSpace(modelName)
	}
	if strings.TrimSpace(modelKey) != "" {
		body["modelKey"] = strings.TrimSpace(modelKey)
	}
	if strings.TrimSpace(providerName) != "" {
		body["provider"] = strings.TrimSpace(providerName)
	}

	resp, err := external.HttpRequest[tryonProviderResponse](external.RequestParams{
		URL:     strings.TrimSpace(targetURL),
		Method:  "POST",
		Headers: headers,
		Body:    body,
	})
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("tryonModelResponseEmpty")
	}

	resultURL := firstNonEmptyString(
		resp.ResultImage,
		resp.ResultURL,
		resp.Data.ResultImage,
		resp.Data.ResultURL,
		resp.Data.Image,
		resp.Data.URL,
	)
	if strings.TrimSpace(resultURL) == "" {
		if strings.TrimSpace(resp.Msg) != "" {
			return "", errors.New(strings.TrimSpace(resp.Msg))
		}
		return "", errors.New("tryonBeautifyResultEmpty")
	}

	return strings.TrimSpace(resultURL), nil
}

func (s *TryonTaskService) invokeAliyunRetouchSkin(resultImage string, modelCfg *tryonModelConfig, fallbackToken string, req clientReq.ApplyTryonBeautifyReq) (string, error) {
	normalizedImage, normalizeErr := s.normalizeAliyunMediaURL(resultImage)
	if normalizeErr != nil {
		return "", normalizeErr
	}

	endpointURL, endpointErr := normalizeAliyunOpenAPIEndpoint(modelCfg.beautifyEndpointURL())
	if endpointErr != nil {
		return "", endpointErr
	}

	accessKeyID, accessKeySecret, securityToken := modelCfg.beautifyAccessCredentials(fallbackToken)
	if strings.TrimSpace(accessKeyID) == "" || strings.TrimSpace(accessKeySecret) == "" {
		return "", errors.New("aliyunBeautifyAccessKeyMissing")
	}

	action := strings.TrimSpace(modelCfg.beautifyModelValue())
	if action == "" {
		action = aliyunRetouchSkinAction
	}

	retouchDegree := modelCfg.beautifyRetouchDegreeValue(req.RetouchDegree)
	whiteningDegree := modelCfg.beautifyWhiteningDegreeValue(req.WhiteningDegree)

	params := map[string]string{
		"Action":           action,
		"Version":          aliyunRetouchSkinVer,
		"Format":           "JSON",
		"AccessKeyId":      strings.TrimSpace(accessKeyID),
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"SignatureVersion": "1.0",
		"SignatureNonce":   generateAliyunSignatureNonce(),
		"ImageURL":         strings.TrimSpace(normalizedImage),
		"RetouchDegree":    formatBeautifyDegree(retouchDegree),
		"WhiteningDegree":  formatBeautifyDegree(whiteningDegree),
	}
	if strings.TrimSpace(securityToken) != "" {
		params["SecurityToken"] = strings.TrimSpace(securityToken)
	}

	signature := signAliyunRPCParams(params, strings.TrimSpace(accessKeySecret))
	params["Signature"] = signature

	body := encodeAliyunRPCParams(params)
	httpReq, err := http.NewRequest(http.MethodPost, endpointURL, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 45 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", errors.New(formatAliyunModelError(string(responseBody), "tryonBeautifyRequestFailed"))
	}

	var apiResp aliyunRetouchSkinResponse
	if err := json.Unmarshal(responseBody, &apiResp); err != nil {
		return "", errors.New("tryonBeautifyResponseInvalid")
	}

	if strings.TrimSpace(apiResp.Code) != "" {
		return "", errors.New(formatAliyunModelError(apiResp.Message, apiResp.Code))
	}

	resultURL := strings.TrimSpace(apiResp.Data.ImageURL)
	if resultURL == "" {
		return "", errors.New("tryonBeautifyResultEmpty")
	}
	return resultURL, nil
}

func (s *TryonTaskService) queryAliyunTryonTask(taskID string, providerToken string, createURL string, taskQueryURL string) (tryonInvokeResult, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return tryonInvokeResult{}, errors.New("aliyunTaskIDRequired")
	}
	if strings.TrimSpace(providerToken) == "" {
		return tryonInvokeResult{}, errors.New("aliyunTryonApiKeyMissing")
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
		return tryonInvokeResult{}, errors.New("aliyunTaskQueryResponseEmpty")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code), ProviderTaskID: taskID}, nil
	}

	status := normalizeProviderTaskStatus(resp.Output.TaskStatus)
	switch status {
	case tryonTaskStatusSuccess:
		resultImage := extractDashscopeImageURL(resp.Output.ImageURL)
		if resultImage == "" {
			return tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: "aliyunTaskSuccessNoImage", ProviderTaskID: taskID}, nil
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
				"tryonTaskFailed",
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
		return "tryonTaskFailed"
	}

	lowerMsg := strings.ToLower(msg)
	if strings.Contains(lowerMsg, "unable to download the media resource") || strings.Contains(lowerMsg, "invalid image url") {
		return "modelCannotDownloadResource"
	}
	if strings.Contains(lowerMsg, "download ") && strings.Contains(lowerMsg, "refused") {
		return "modelDownloadRefused"
	}
	if strings.Contains(lowerMsg, "cf-mitigated") || strings.Contains(lowerMsg, "cloudflare") || strings.Contains(lowerMsg, "challenge") {
		return "modelDownloadRefused"
	}
	if strings.Contains(lowerMsg, "invalidurl") {
		return "modelImageURLInvalidOrUnreachable"
	}
	if strings.Contains(lowerMsg, "data inspection") || strings.Contains(lowerMsg, "datainspection") {
		if strings.Contains(lowerMsg, "download") || strings.Contains(lowerMsg, "resource") || strings.Contains(lowerMsg, "url") || strings.Contains(lowerMsg, "refused") || strings.Contains(lowerMsg, "403") {
			return "modelDownloadRefused"
		}
		return "modelImageRiskCheckFailed"
	}
	return msg
}

func (s *TryonTaskService) normalizeAliyunMediaURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", errors.New("mediaURLRequired")
	}
	if strings.HasPrefix(strings.ToLower(rawURL), "data:") {
		return "", errors.New("mediaURLDataURINotAllowed")
	}

	publicBaseURL := s.getTryonMediaPublicBaseURL()

	if !isAbsoluteHTTPURL(rawURL) {
		if publicBaseURL == "" {
			return "", errors.New("mediaURLPublicURLRequired")
		}
		return joinBaseURLAndResourcePath(publicBaseURL, rawURL)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() == "" {
		return "", errors.New("mediaURLInvalid")
	}

	if isLoopbackOrPrivateHost(parsed.Hostname()) {
		if publicBaseURL == "" {
			return "", errors.New("mediaURLPrivateHostNotAccessible")
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
		return "", errors.New("publicBaseURLRequired")
	}

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil || parsedBaseURL.Scheme == "" || parsedBaseURL.Hostname() == "" {
		return "", errors.New("publicBaseURLInvalid")
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
	refinerURL := strings.TrimSpace(providerURL)
	refinerToken := strings.TrimSpace(providerToken)
	refinerTaskQueryURL := ""
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
		refinerURL = modelCfg.refinerEndpointURL(providerURL)
		refinerToken = modelCfg.refinerAuthToken(providerToken)
		refinerTaskQueryURL = modelCfg.refinerTaskQueryURL()
	}
	if strings.TrimSpace(refinerTaskQueryURL) == "" {
		refinerTaskQueryURL = taskQueryURL
	}

	if !isAliyunTryonModel(providerName, modelName, providerURL, task.SceneType) {
		return nil
	}

	if task.EnableRefiner && strings.EqualFold(strings.TrimSpace(task.RefinerStatus), tryonRefinerStatusProcessing) {
		refinerTaskID := strings.TrimSpace(task.RefinerTaskID)
		if refinerTaskID == "" {
			refinerTaskID = strings.TrimSpace(task.ProviderTaskID)
		}
		if refinerTaskID == "" {
			return nil
		}

		result, err := s.queryAliyunTryonTask(refinerTaskID, refinerToken, refinerURL, refinerTaskQueryURL)
		if err != nil {
			return err
		}

		switch result.Status {
		case tryonTaskStatusProcessing:
			return nil
		case tryonTaskStatusSuccess:
			storedResultImage := s.persistTryonResultImage(strings.TrimSpace(result.ResultImage), task.TaskNo)
			now := time.Now()
			updateResult := global.GVA_DB.Model(&client.TryonTask{}).
				Where("id = ? AND user_id = ? AND status = ?", task.ID, task.UserID, tryonTaskStatusProcessing).
				Updates(map[string]interface{}{
					"status":           tryonTaskStatusSuccess,
					"provider_task_id": strings.TrimSpace(refinerTaskID),
					"refiner_task_id":  strings.TrimSpace(refinerTaskID),
					"refiner_status":   tryonRefinerStatusSuccess,
					"result_image":     storedResultImage,
					"error_message":    "",
					"completed_at":     now,
				})
			if updateResult.Error != nil {
				return updateResult.Error
			}
			if updateResult.RowsAffected == 0 {
				return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
			}
			task.Status = tryonTaskStatusSuccess
			task.ProviderTaskID = strings.TrimSpace(refinerTaskID)
			task.RefinerTaskID = strings.TrimSpace(refinerTaskID)
			task.RefinerStatus = tryonRefinerStatusSuccess
			task.ResultImage = storedResultImage
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

	result, err := s.queryAliyunTryonTask(task.ProviderTaskID, providerToken, providerURL, taskQueryURL)
	if err != nil {
		return err
	}

	switch result.Status {
	case tryonTaskStatusProcessing:
		return nil
	case tryonTaskStatusSuccess:
		if task.EnableRefiner && strings.EqualFold(strings.TrimSpace(task.RefinerStatus), tryonRefinerStatusPending) {
			refinerReq := clientReq.CreateTryonTaskReq{
				SceneType:     task.SceneType,
				SourceImage:   task.SourceImage,
				TemplateImage: task.TemplateImage,
				ModelKey:      task.Provider,
			}
			refinerResult, refinerErr := s.invokeAliyunRefinerAsync(refinerReq, strings.TrimSpace(result.ResultImage), modelCfg)
			if refinerErr != nil {
				if err = s.markFailedAndRefund(ctx, task, refinerErr.Error()); err != nil {
					return err
				}
				return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
			}

			switch refinerResult.Status {
			case tryonTaskStatusProcessing:
				refinerTaskID := strings.TrimSpace(refinerResult.ProviderTaskID)
				updateResult := global.GVA_DB.Model(&client.TryonTask{}).
					Where("id = ? AND user_id = ? AND status = ?", task.ID, task.UserID, tryonTaskStatusProcessing).
					Updates(map[string]interface{}{
						"provider_task_id": refinerTaskID,
						"refiner_task_id":  refinerTaskID,
						"refiner_status":   tryonRefinerStatusProcessing,
						"error_message":    "",
					})
				if updateResult.Error != nil {
					return updateResult.Error
				}
				if updateResult.RowsAffected == 0 {
					return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
				}
				task.ProviderTaskID = refinerTaskID
				task.RefinerTaskID = refinerTaskID
				task.RefinerStatus = tryonRefinerStatusProcessing
				task.ErrorMessage = ""
				return nil
			case tryonTaskStatusSuccess:
				storedResultImage := s.persistTryonResultImage(strings.TrimSpace(refinerResult.ResultImage), task.TaskNo)
				now := time.Now()
				refinerTaskID := strings.TrimSpace(refinerResult.ProviderTaskID)
				updateResult := global.GVA_DB.Model(&client.TryonTask{}).
					Where("id = ? AND user_id = ? AND status = ?", task.ID, task.UserID, tryonTaskStatusProcessing).
					Updates(map[string]interface{}{
						"status":           tryonTaskStatusSuccess,
						"provider_task_id": refinerTaskID,
						"refiner_task_id":  refinerTaskID,
						"refiner_status":   tryonRefinerStatusSuccess,
						"result_image":     storedResultImage,
						"error_message":    "",
						"completed_at":     now,
					})
				if updateResult.Error != nil {
					return updateResult.Error
				}
				if updateResult.RowsAffected == 0 {
					return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
				}
				task.Status = tryonTaskStatusSuccess
				task.ProviderTaskID = refinerTaskID
				task.RefinerTaskID = refinerTaskID
				task.RefinerStatus = tryonRefinerStatusSuccess
				task.ResultImage = storedResultImage
				task.ErrorMessage = ""
				task.CompletedAt = &now
				return nil
			case tryonTaskStatusFailed:
				if err = s.markFailedAndRefund(ctx, task, refinerResult.ErrorMessage); err != nil {
					return err
				}
				return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
			default:
				if err = s.markFailedAndRefund(ctx, task, "tryonRefinerStatusUnknown"); err != nil {
					return err
				}
				return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
			}
		}

		storedResultImage := s.persistTryonResultImage(strings.TrimSpace(result.ResultImage), task.TaskNo)
		now := time.Now()
		updateResult := global.GVA_DB.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ? AND status = ?", task.ID, task.UserID, tryonTaskStatusProcessing).
			Updates(map[string]interface{}{
				"status":          tryonTaskStatusSuccess,
				"result_image":    storedResultImage,
				"refiner_status":  tryonRefinerStatusDisabled,
				"refiner_task_id": "",
				"error_message":   "",
				"completed_at":    now,
			})
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
		}
		task.Status = tryonTaskStatusSuccess
		task.ResultImage = storedResultImage
		task.RefinerStatus = tryonRefinerStatusDisabled
		task.RefinerTaskID = ""
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

func (s *TryonTaskService) persistTryonResultImage(rawURL string, taskNo string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	storedURL, err := s.uploadTryonResultToUniGet(rawURL, taskNo)
	if err != nil {
		global.GVA_LOG.Warn("试衣结果图回存失败，使用原始地址", zap.Error(err), zap.String("taskNo", strings.TrimSpace(taskNo)), zap.String("url", rawURL))
		return rawURL
	}
	if strings.TrimSpace(storedURL) == "" {
		return rawURL
	}
	return storedURL
}

func (s *TryonTaskService) uploadTryonResultToUniGet(rawURL string, taskNo string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", errors.New("tryonResultURLRequired")
	}

	if strings.Contains(strings.ToLower(rawURL), "/"+tryonResultStoreFolder+"/") {
		return rawURL, nil
	}

	if !isAbsoluteHTTPURL(rawURL) {
		return rawURL, nil
	}

	if !strings.EqualFold(strings.TrimSpace(global.GVA_CONFIG.System.OssType), "aws-s3") {
		return rawURL, nil
	}

	baseURL := strings.TrimRight(strings.TrimSpace(global.GVA_CONFIG.AwsS3.BaseURL), "/")
	if baseURL == "" {
		return rawURL, nil
	}

	client := &http.Client{Timeout: 25 * time.Second}
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "gin-vue-admin/tryon-result-fetch")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", errors.New("tryonResultDownloadFailed")
	}

	limitedReader := io.LimitReader(resp.Body, tryonResultMaxImageBytes+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", errors.New("tryonResultContentEmpty")
	}
	if len(data) > tryonResultMaxImageBytes {
		return "", errors.New("tryonResultTooLarge")
	}

	contentType := strings.TrimSpace(resp.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	ext := detectTryonResultImageExt(contentType, rawURL)
	objectKey := buildTryonResultObjectKey(taskNo, ext)

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(global.GVA_CONFIG.AwsS3.Region),
		Endpoint:         aws.String(global.GVA_CONFIG.AwsS3.Endpoint),
		S3ForcePathStyle: aws.Bool(global.GVA_CONFIG.AwsS3.S3ForcePathStyle),
		DisableSSL:       aws.Bool(global.GVA_CONFIG.AwsS3.DisableSSL),
		Credentials: credentials.NewStaticCredentials(
			global.GVA_CONFIG.AwsS3.SecretID,
			global.GVA_CONFIG.AwsS3.SecretKey,
			"",
		),
	})
	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(global.GVA_CONFIG.AwsS3.Bucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(strings.Split(contentType, ";")[0]),
	})
	if err != nil {
		return "", err
	}

	return baseURL + "/" + objectKey, nil
}

func buildTryonResultObjectKey(taskNo string, ext string) string {
	safeTaskNo := strings.ToLower(strings.TrimSpace(taskNo))
	if safeTaskNo == "" {
		safeTaskNo = "task"
	}

	builder := strings.Builder{}
	for _, r := range safeTaskNo {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			builder.WriteRune(r)
			continue
		}
		builder.WriteRune('-')
	}
	safeTaskNo = strings.Trim(builder.String(), "-")
	if safeTaskNo == "" {
		safeTaskNo = "task"
	}
	if len(safeTaskNo) > 32 {
		safeTaskNo = safeTaskNo[:32]
	}

	ext = strings.ToLower(strings.TrimSpace(ext))
	if ext == "" || !strings.HasPrefix(ext, ".") {
		ext = ".jpg"
	}

	fileName := fmt.Sprintf("%s-%d%s", safeTaskNo, time.Now().UnixNano(), ext)
	pathParts := make([]string, 0, 3)
	if prefix := strings.Trim(strings.TrimSpace(global.GVA_CONFIG.AwsS3.PathPrefix), "/"); prefix != "" {
		pathParts = append(pathParts, prefix)
	}
	pathParts = append(pathParts, tryonResultStoreFolder, fileName)
	return strings.Join(pathParts, "/")
}

func detectTryonResultImageExt(contentType string, rawURL string) string {
	ct := strings.ToLower(strings.TrimSpace(strings.Split(strings.TrimSpace(contentType), ";")[0]))
	switch ct {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	case "image/bmp":
		return ".bmp"
	case "image/avif":
		return ".avif"
	}

	ext := ""
	if parsed, err := url.Parse(strings.TrimSpace(rawURL)); err == nil {
		ext = strings.ToLower(strings.TrimSpace(path.Ext(parsed.Path)))
	}

	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif", ".bmp", ".avif":
		if ext == ".jpeg" {
			return ".jpg"
		}
		return ext
	default:
		return ".jpg"
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
		case strings.Contains(keyLower, "refiner") || strings.Contains(providerLower, "refiner"):
			return aliyunTryonRefinerModel
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

func isAliyunBeautifyModel(provider string, modelName string, providerURL string) bool {
	modelLower := strings.ToLower(strings.TrimSpace(modelName))
	if modelLower == strings.ToLower(aliyunRetouchSkinAction) {
		return true
	}

	providerLower := strings.ToLower(strings.TrimSpace(provider))
	if strings.Contains(providerLower, "aliyun") || strings.Contains(providerLower, "dashscope") || strings.Contains(providerLower, "facebody") {
		return true
	}

	urlLower := strings.ToLower(strings.TrimSpace(providerURL))
	if strings.Contains(urlLower, "facebody.aliyuncs.com") {
		return true
	}
	if strings.Contains(urlLower, "aliyuncs.com") && strings.Contains(modelLower, "retouch") {
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

func canEnableAliyunRefiner(sceneType string, modelCfg *tryonModelConfig) bool {
	if modelCfg == nil || !modelCfg.supportsRefinerValue() {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(sceneType), "clothes") {
		return false
	}

	modelName := resolveAliyunModelName(modelCfg, sceneType)
	modelLower := strings.ToLower(strings.TrimSpace(modelName))
	return modelLower == "aitryon" || modelLower == "aitryon-plus"
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

func extractDashscopeImageURL(raw json.RawMessage) string {
	text := strings.TrimSpace(string(raw))
	if text == "" || text == "null" {
		return ""
	}

	var single string
	if err := json.Unmarshal(raw, &single); err == nil {
		return strings.TrimSpace(single)
	}

	var list []string
	if err := json.Unmarshal(raw, &list); err == nil {
		for _, item := range list {
			if strings.TrimSpace(item) != "" {
				return strings.TrimSpace(item)
			}
		}
	}

	var mixedList []interface{}
	if err := json.Unmarshal(raw, &mixedList); err == nil {
		for _, item := range mixedList {
			if urlText := strings.TrimSpace(fmt.Sprint(item)); urlText != "" && urlText != "<nil>" {
				return urlText
			}
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

func normalizeAliyunOpenAPIEndpoint(rawURL string) (string, error) {
	endpoint := strings.TrimSpace(rawURL)
	if endpoint == "" {
		endpoint = aliyunRetouchSkinURL
	}
	if !strings.Contains(endpoint, "://") {
		endpoint = "https://" + endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil || strings.TrimSpace(u.Scheme) == "" || strings.TrimSpace(u.Host) == "" {
		return "", errors.New("tryonBeautifyEndpointInvalid")
	}
	u.RawQuery = ""
	u.Fragment = ""
	if strings.TrimSpace(u.Path) == "" {
		u.Path = "/"
	}
	return u.String(), nil
}

func clampBeautifyDegree(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func formatBeautifyDegree(value float64) string {
	value = clampBeautifyDegree(value)
	return strconv.FormatFloat(value, 'f', 1, 64)
}

func parseAKSKFromToken(raw string) (accessKeyID string, accessKeySecret string, securityToken string, ok bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", "", false
	}

	if strings.HasPrefix(raw, "{") {
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &payload); err == nil {
			accessKeyID = firstNonEmptyString(
				strings.TrimSpace(fmt.Sprint(payload["accessKeyId"])),
				strings.TrimSpace(fmt.Sprint(payload["access_key_id"])),
				strings.TrimSpace(fmt.Sprint(payload["ak"])),
			)
			accessKeySecret = firstNonEmptyString(
				strings.TrimSpace(fmt.Sprint(payload["accessKeySecret"])),
				strings.TrimSpace(fmt.Sprint(payload["access_key_secret"])),
				strings.TrimSpace(fmt.Sprint(payload["sk"])),
			)
			securityToken = firstNonEmptyString(
				strings.TrimSpace(fmt.Sprint(payload["securityToken"])),
				strings.TrimSpace(fmt.Sprint(payload["security_token"])),
				strings.TrimSpace(fmt.Sprint(payload["token"])),
			)
			if accessKeyID != "" && accessKeySecret != "" {
				return accessKeyID, accessKeySecret, securityToken, true
			}
		}
	}

	for _, sep := range []string{"|", ":", ","} {
		parts := strings.Split(raw, sep)
		if len(parts) < 2 {
			continue
		}
		accessKeyID = strings.TrimSpace(parts[0])
		accessKeySecret = strings.TrimSpace(parts[1])
		if len(parts) >= 3 {
			securityToken = strings.TrimSpace(parts[2])
		}
		if accessKeyID != "" && accessKeySecret != "" {
			return accessKeyID, accessKeySecret, securityToken, true
		}
	}

	return "", "", "", false
}

func generateAliyunSignatureNonce() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func encodeAliyunRPCParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, aliyunPercentEncode(key)+"="+aliyunPercentEncode(params[key]))
	}
	return strings.Join(parts, "&")
}

func signAliyunRPCParams(params map[string]string, accessKeySecret string) string {
	canonicalized := encodeAliyunRPCParams(params)
	stringToSign := "POST&%2F&" + aliyunPercentEncode(canonicalized)

	mac := hmac.New(sha1.New, []byte(accessKeySecret+"&"))
	_, _ = mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func aliyunPercentEncode(value string) string {
	encoded := url.QueryEscape(value)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}

func generateTryonTaskNo() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "TRYON_" + time.Now().Format("20060102150405") + "_" + hex.EncodeToString(b)
}

func generateTryonBeautifyTaskNo() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "BEAUTIFY_" + time.Now().Format("20060102150405") + "_" + hex.EncodeToString(b)
}

func generateTryonRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "REQ_" + hex.EncodeToString(b)
}

func buildPointRecord(userID uint, assetType string, changeType string, points int, operationType string, reason string, remark string) *client.PointRecord {
	uid := int(userID)
	at := assetType
	ct := changeType
	p := points
	op := operationType
	rs := reason
	rm := remark
	if strings.TrimSpace(at) == "" {
		at = client.AssetTypePoint
	}

	return &client.PointRecord{
		AssetType:     &at,
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
		return errors.New("tryonGuestInitTxMissing")
	}

	guestPoints := s.getConfigIntAllowZero("tryon_guest_init_points", 0)
	if guestPoints <= 0 {
		return nil
	}

	return s.grantPointsIfNotExists(ctx, tx, userID, tryonOpGuestInit, guestPoints, tryonReasonGuestInit, "guest_first_tryon")
}

func (s *TryonTaskService) grantPointsIfNotExists(ctx context.Context, tx *gorm.DB, userID uint, operationType string, points int, reason string, remark string) error {
	uid := int(userID)
	var user client.ClientUser
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("userNotExist")
		}
		return err
	}

	var count int64
	query := tx.Model(&client.PointRecord{}).
		Where("user_id = ? AND operation_type = ? AND (asset_type = ? OR asset_type IS NULL)", uid, operationType, client.AssetTypeTryonPoint)
	if strings.TrimSpace(remark) != "" {
		query = query.Where("remark = ?", strings.TrimSpace(remark))
	}
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	pointRecordService := PointRecordService{}
	record := buildPointRecord(userID, client.AssetTypeTryonPoint, "increase", points, operationType, reason, remark)
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
