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

	tryonParsingStatusDisabled = "disabled"
	tryonParsingStatusPending  = "pending"
	tryonParsingStatusSuccess  = "success"
	tryonParsingStatusFailed   = "failed"

	tryonBeautifyStatusDisabled   = "disabled"
	tryonBeautifyStatusProcessing = "processing"
	tryonBeautifyStatusSuccess    = "success"
	tryonBeautifyStatusFailed     = "failed"

	tryonOpConsume         = "tryon_consume"
	tryonOpRefund          = "tryon_refund"
	tryonOpParsingConsume  = "tryon_parsing_consume"
	tryonOpParsingRefund   = "tryon_parsing_refund"
	tryonOpRefinerConsume  = "tryon_refiner_consume"
	tryonOpRefinerRefund   = "tryon_refiner_refund"
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
	aliyunShoeTryonURL      = "https://dashscope.aliyuncs.com/api/v1/services/aigc/virtualmodel/generation"
	aliyunParsingProcessURL = "https://dashscope.aliyuncs.com/api/v1/services/vision/image-process/process"
	aliyunTryonRefinerModel = "aitryon-refiner"
	aliyunShoeTryonModel    = "shoemodel-v1"
	aliyunRetouchSkinURL    = "https://facebody.cn-shanghai.aliyuncs.com/"
	aliyunRetouchSkinAction = "RetouchSkin"
	aliyunRetouchSkinVer    = "2019-12-30"

	tryonResultStoreFolder   = "cloth-on/middle-transfer"
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
		s.removeLocalTryonMediaFile(task.TemplateImageLower)
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
	ModelUsage          string `json:"modelUsage"`
	Provider            string `json:"provider"`
	TaskCount           int64  `json:"taskCount"`
	SuccessCount        int64  `json:"successCount"`
	FailedCount         int64  `json:"failedCount"`
	ProcessingCount     int64  `json:"processingCount"`
	RefinerEnabledCount int64  `json:"refinerEnabledCount"`
	ParsingCount        int64  `json:"parsingCount"`
	TotalCostPoints     int64  `json:"totalCostPoints"`
}

type TryonTaskStatsData struct {
	Total               int64                     `json:"total"`
	Processing          int64                     `json:"processing"`
	Success             int64                     `json:"success"`
	Failed              int64                     `json:"failed"`
	TotalCostPoints     int64                     `json:"totalCostPoints"`
	RefinerEnabledCount int64                     `json:"refinerEnabledCount"`
	ParsingEnabledCount int64                     `json:"parsingEnabledCount"`
	ModelCallTotal      int64                     `json:"modelCallTotal"`
	ModelStats          []TryonTaskModelStatsItem `json:"modelStats"`
}

type tryonModelConfig struct {
	Key                   string   `json:"key"`
	ModelUsage            string   `json:"modelUsage"`
	BeautifyModelKey      string   `json:"beautifyModelKey"`
	ParsingModelKey       string   `json:"parsingModelKey"`
	Enabled               *bool    `json:"enabled"`
	SupportsRefiner       *bool    `json:"supportsRefiner"`
	SupportsBeautify      *bool    `json:"supportsBeautify"`
	SceneType             string   `json:"sceneType"`
	Scenes                []string `json:"scenes"`
	Model                 string   `json:"model"`
	Cost                  int      `json:"cost"`
	Provider              string   `json:"provider"`
	Mode                  string   `json:"mode"`
	URL                   string   `json:"url"`
	Token                 string   `json:"token"`
	TokenBackups          []string `json:"tokenBackups"`
	BackupTokens          []string `json:"backupTokens"`
	ProviderURL           string   `json:"providerUrl"`
	ProviderToken         string   `json:"providerToken"`
	TaskQueryURL          string   `json:"taskQueryUrl"`
	FreeQuotaTotal        int      `json:"freeQuotaTotal"`
	RefinerFreeQuotaTotal int      `json:"refinerFreeQuotaTotal"`
	RefinerModel          string   `json:"refinerModel"`
	RefinerModelKey       string   `json:"refinerModelKey"`
	RefinerExtraCost      int      `json:"refinerExtraCost"`
	RefinerExtraPoints    int      `json:"refinerExtraPoints"`
	ParsingExtraCost      int      `json:"parsingExtraCost"`
	ParsingExtraPoints    int      `json:"parsingExtraPoints"`
	RefinerURL            string   `json:"refinerUrl"`
	RefinerToken          string   `json:"refinerToken"`
	RefinerTaskQueryURL   string   `json:"refinerTaskQueryUrl"`
	BeautifyModel         string   `json:"beautifyModel"`
	BeautifyExtraCost     int      `json:"beautifyExtraCost"`
	BeautifyExtraPoints   int      `json:"beautifyExtraPoints"`
	BeautifyURL           string   `json:"beautifyUrl"`
	BeautifyToken         string   `json:"beautifyToken"`
	BeautifyAccessKeyID   string   `json:"beautifyAccessKeyId"`
	BeautifyAccessKeySec  string   `json:"beautifyAccessKeySecret"`
	BeautifySecurityToken string   `json:"beautifySecurityToken"`
	BeautifyRetouch       float64  `json:"beautifyRetouchDegree"`
	BeautifyWhitening     float64  `json:"beautifyWhiteningDegree"`
	RefinerGender         string   `json:"refinerGender"`
	ApiName               string   `json:"apiName"`
	GarmentDes            string   `json:"garmentDes"`
	IsChecked             *bool    `json:"isChecked"`
	IsCheckedCrop         *bool    `json:"isCheckedCrop"`
	DenoiseSteps          int      `json:"denoiseSteps"`
	Seed                  int      `json:"seed"`
	Resolution            int      `json:"resolution"`
	RestoreFace           *bool    `json:"restoreFace"`
	ClothesType           []string `json:"clothesType"`
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
	switch usage {
	case "tryon", "beautify", "refiner", "parsing":
		return usage
	}

	if m.SupportsBeautify != nil && *m.SupportsBeautify {
		return "beautify"
	}

	modelName := resolveAliyunModelName(m, "")
	if isAliyunParsingModel(m.Provider, modelName, m.endpointURL(), "") {
		return "parsing"
	}

	if strings.Contains(strings.ToLower(strings.TrimSpace(m.Key)), "refiner") {
		return "refiner"
	}

	return "tryon"
}

func (m *tryonModelConfig) isTryonModelUsage() bool {
	return m == nil || m.usageValue() == "tryon"
}

func (m *tryonModelConfig) isBeautifyModelUsage() bool {
	return m != nil && m.usageValue() == "beautify"
}

func (m *tryonModelConfig) isRefinerModelUsage() bool {
	return m != nil && m.usageValue() == "refiner"
}

func (m *tryonModelConfig) isParsingModelUsage() bool {
	if m == nil {
		return false
	}
	if m.usageValue() == "parsing" {
		return true
	}
	modelName := resolveAliyunModelName(m, "")
	return isAliyunParsingModel(m.Provider, modelName, m.endpointURL(), "")
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
	for _, token := range m.authTokenCandidates() {
		return token
	}
	return ""
}

func appendUniqueTokenCandidate(dst []string, value string) []string {
	token := strings.TrimSpace(value)
	if token == "" {
		return dst
	}
	for _, existing := range dst {
		if existing == token {
			return dst
		}
	}
	return append(dst, token)
}

func (m *tryonModelConfig) tokenBackupsValue() []string {
	if m == nil {
		return nil
	}
	result := make([]string, 0, len(m.TokenBackups)+len(m.BackupTokens))
	for _, item := range m.TokenBackups {
		result = appendUniqueTokenCandidate(result, item)
	}
	for _, item := range m.BackupTokens {
		result = appendUniqueTokenCandidate(result, item)
	}
	return result
}

func (m *tryonModelConfig) authTokenCandidates() []string {
	if m == nil {
		return nil
	}
	result := make([]string, 0, 1+len(m.TokenBackups)+len(m.BackupTokens)+1)
	result = appendUniqueTokenCandidate(result, m.Token)
	for _, backup := range m.tokenBackupsValue() {
		result = appendUniqueTokenCandidate(result, backup)
	}
	result = appendUniqueTokenCandidate(result, m.ProviderToken)
	return result
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
	if !m.isTryonModelUsage() {
		return false
	}
	if m.SupportsRefiner != nil {
		return *m.SupportsRefiner
	}
	if strings.TrimSpace(m.RefinerModelKey) != "" {
		return true
	}

	modelLower := strings.ToLower(strings.TrimSpace(m.Model))
	if modelLower == "aitryon" || modelLower == "aitryon-plus" {
		return true
	}

	keyLower := strings.ToLower(strings.TrimSpace(m.Key))
	return strings.Contains(keyLower, "aliyun_aitryon") && !strings.Contains(keyLower, "parsing")
}

func (m *tryonModelConfig) refinerModelValue() string {
	if m == nil {
		return aliyunTryonRefinerModel
	}
	if strings.TrimSpace(m.RefinerModel) != "" {
		return strings.TrimSpace(m.RefinerModel)
	}
	if m.isRefinerModelUsage() && strings.TrimSpace(m.Model) != "" {
		return strings.TrimSpace(m.Model)
	}
	if m.isTryonModelUsage() && strings.TrimSpace(m.Model) == "aitryon-refiner" {
		return strings.TrimSpace(m.Model)
	}
	if m.isRefinerModelUsage() && strings.TrimSpace(m.Model) == "" {
		return aliyunTryonRefinerModel
	}
	if strings.TrimSpace(m.RefinerModel) == "" {
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

func (m *tryonModelConfig) parsingCostValue() int {
	if m == nil {
		return 0
	}
	if m.ParsingExtraCost > 0 {
		return m.ParsingExtraCost
	}
	if m.ParsingExtraPoints > 0 {
		return m.ParsingExtraPoints
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
	return 0
}

func (m *tryonModelConfig) beautifyWhiteningDegreeValue(override float64) float64 {
	if override > 0 {
		return clampBeautifyDegree(override)
	}
	if m != nil && m.BeautifyWhitening > 0 {
		return clampBeautifyDegree(m.BeautifyWhitening)
	}
	return 0
}

func (m *tryonModelConfig) refinerEndpointURL(fallback string) string {
	if m != nil && m.isRefinerModelUsage() {
		if strings.TrimSpace(m.URL) != "" {
			return strings.TrimSpace(m.URL)
		}
		if strings.TrimSpace(m.ProviderURL) != "" {
			return strings.TrimSpace(m.ProviderURL)
		}
	}
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
	for _, token := range m.refinerAuthTokenCandidates(fallback) {
		return token
	}
	return ""
}

func (m *tryonModelConfig) refinerAuthTokenCandidates(fallback string) []string {
	result := make([]string, 0)
	if m != nil && m.isRefinerModelUsage() {
		result = appendUniqueTokenCandidate(result, m.Token)
		for _, backup := range m.tokenBackupsValue() {
			result = appendUniqueTokenCandidate(result, backup)
		}
		result = appendUniqueTokenCandidate(result, m.ProviderToken)
	}
	if m != nil {
		result = appendUniqueTokenCandidate(result, m.RefinerToken)
		if !m.isRefinerModelUsage() {
			result = appendUniqueTokenCandidate(result, m.Token)
			for _, backup := range m.tokenBackupsValue() {
				result = appendUniqueTokenCandidate(result, backup)
			}
			result = appendUniqueTokenCandidate(result, m.ProviderToken)
		}
	}
	result = appendUniqueTokenCandidate(result, fallback)
	return result
}

func (m *tryonModelConfig) refinerTaskQueryURL() string {
	if m != nil && m.isRefinerModelUsage() {
		if strings.TrimSpace(m.TaskQueryURL) != "" {
			return strings.TrimSpace(m.TaskQueryURL)
		}
	}
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

	ParsingAttempted bool
	ParsingSucceeded bool
	ParsingModelKey  string
	ParsingError     string
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

type dashscopeTaskResultItem struct {
	URL       string          `json:"url"`
	ImageURL  json.RawMessage `json:"image_url"`
	ResultURL string          `json:"result_url"`
}

type dashscopeTaskQueryResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Output    struct {
		TaskID        string                    `json:"task_id"`
		TaskStatus    string                    `json:"task_status"`
		ImageURL      json.RawMessage           `json:"image_url"`
		ImageURLs     []string                  `json:"image_urls"`
		ResultURL     string                    `json:"result_url"`
		ResultURLs    []string                  `json:"result_urls"`
		Results       []dashscopeTaskResultItem `json:"results"`
		ParsingImgURL []string                  `json:"parsing_img_url"`
		CropImgURL    []string                  `json:"crop_img_url"`
		Code          string                    `json:"code"`
		Message       string                    `json:"message"`
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
	enableParsing := req.EnableParsing && canEnableAliyunParsing(req.SceneType, modelCfg, req.TemplatePart)
	refinerCostPoints := 0
	parsingCostPoints := 0
	resolvedRefinerModelKey := ""
	resolvedParsingModelKey := ""
	if enableRefiner {
		refinerModelCfg, resolveErr := s.resolveRefinerModelConfig(req.SceneType, modelCfg)
		if resolveErr != nil {
			return task, false, resolveErr
		}
		if refinerModelCfg == nil {
			return task, false, errors.New("tryonRefinerModelUnavailable")
		}
		refinerCostPoints = refinerModelCfg.refinerExtraCostValue()
		resolvedRefinerModelKey = strings.TrimSpace(refinerModelCfg.Key)
		if resolvedRefinerModelKey == "" {
			resolvedRefinerModelKey = strings.TrimSpace(modelCfg.RefinerModelKey)
		}
	}
	if enableParsing {
		preferredParsingModelKey := ""
		if modelCfg != nil {
			preferredParsingModelKey = strings.TrimSpace(modelCfg.ParsingModelKey)
		}
		parsingModelCfg, resolveErr := s.resolveAliyunParsingAssistModelConfig(req.SceneType, preferredParsingModelKey)
		if resolveErr != nil {
			global.GVA_LOG.Warn("解析分割增强模型配置失败，按0扣费继续", zap.Error(resolveErr), zap.String("sceneType", req.SceneType), zap.String("preferredParsingModelKey", preferredParsingModelKey))
		} else if parsingModelCfg != nil {
			parsingCostPoints = parsingModelCfg.parsingCostValue()
			resolvedParsingModelKey = strings.TrimSpace(parsingModelCfg.Key)
		}
		if resolvedParsingModelKey == "" && modelCfg != nil {
			resolvedParsingModelKey = strings.TrimSpace(modelCfg.ParsingModelKey)
		}
	}

	costPoints := baseCostPoints + refinerCostPoints + parsingCostPoints
	if costPoints < 1 {
		costPoints = 1
	}

	refinerStatus := tryonRefinerStatusDisabled
	if enableRefiner {
		refinerStatus = tryonRefinerStatusPending
	}

	parsingStatus := tryonParsingStatusDisabled
	if enableParsing {
		parsingStatus = tryonParsingStatusPending
	}

	task = client.TryonTask{
		UserID:             userID,
		RequestID:          requestID,
		TaskNo:             generateTryonTaskNo(),
		SceneType:          req.SceneType,
		Status:             tryonTaskStatusProcessing,
		SourceImage:        req.SourceImage,
		TemplateImage:      req.TemplateImage,
		TemplateImageLower: req.TemplateImageLower,
		Provider:           providerName,
		EnableRefiner:      enableRefiner,
		EnableParsing:      enableParsing,
		ParsingModelKey:    resolvedParsingModelKey,
		ParsingStatus:      parsingStatus,
		ParsingError:       "",
		RefinerModelKey:    resolvedRefinerModelKey,
		RefinerStatus:      refinerStatus,
		BeautifyStatus:     tryonBeautifyStatusDisabled,
		BaseCostPoints:     baseCostPoints,
		ParsingCostPoints:  parsingCostPoints,
		RefinerCostPoints:  refinerCostPoints,
		CostPoints:         costPoints,
	}

	if err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx", tx)
		if createErr := tx.Create(&task).Error; createErr != nil {
			return createErr
		}

		pointRecordService := PointRecordService{}
		if baseCostPoints > 0 {
			record := buildPointRecord(userID, client.AssetTypeTryonPoint, "decrease", baseCostPoints, tryonOpConsume, tryonReasonTaskDeduct, "tryon_task:"+task.TaskNo)
			if pointErr := pointRecordService.CreatePointRecord(txCtx, record); pointErr != nil {
				return pointErr
			}
		}
		if refinerCostPoints > 0 {
			remark := "tryon_refiner:" + task.TaskNo
			if strings.TrimSpace(resolvedRefinerModelKey) != "" {
				remark = remark + ":" + strings.TrimSpace(resolvedRefinerModelKey)
			}
			record := buildPointRecord(userID, client.AssetTypeTryonPoint, "decrease", refinerCostPoints, tryonOpRefinerConsume, tryonReasonTaskDeduct, remark)
			if pointErr := pointRecordService.CreatePointRecord(txCtx, record); pointErr != nil {
				return pointErr
			}
		}
		if parsingCostPoints > 0 {
			remark := buildParsingPointRemark(task.TaskNo, resolvedParsingModelKey)
			record := buildPointRecord(userID, client.AssetTypeTryonPoint, "decrease", parsingCostPoints, tryonOpParsingConsume, tryonReasonTaskDeduct, remark)
			if pointErr := pointRecordService.CreatePointRecord(txCtx, record); pointErr != nil {
				return pointErr
			}
		}
		return nil
	}); err != nil {
		return task, false, err
	}

	resolvedReq := req
	resolvedReq.EnableRefiner = enableRefiner
	resolvedReq.EnableParsing = enableParsing
	traceCtx := createModelCallTraceContext(userID, task, resolvedReq)
	invokeResult, invokeErr := s.invokeTryonProvider(resolvedReq, modelCfg, traceCtx)
	if invokeErr != nil {
		refundErr := s.markFailedAndRefund(ctx, &task, invokeErr.Error())
		if refundErr != nil {
			global.GVA_LOG.Error("试衣失败退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
			return task, false, errors.New("tryonRefundFailedContactAdmin")
		}
		_ = global.GVA_DB.Where("id = ?", task.ID).First(&task).Error
		return task, false, invokeErr
	}

	parsingStageUpdates := buildParsingStageUpdates(task, invokeResult)

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
		mergeTryonTaskUpdates(updates, parsingStageUpdates)

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
			storedMainResultImage := s.persistTryonResultImage(strings.TrimSpace(invokeResult.ResultImage), task.TaskNo)
			refinerResult, refinerErr := s.invokeAliyunRefinerAsync(resolvedReq, strings.TrimSpace(invokeResult.ResultImage), modelCfg, traceCtx)
			if refinerErr != nil {
				fallbackErr := s.markTryonRefinerFailedAndKeepSuccess(ctx, &task, storedMainResultImage, refinerErr.Error(), "", parsingStageUpdates)
				if fallbackErr != nil {
					global.GVA_LOG.Error("试衣精修失败回退主图且退币异常", zap.Error(fallbackErr), zap.Uint("taskID", task.ID))
					return task, false, errors.New("tryonRefinerRefundFailedContactAdmin")
				}
				break
			}

			switch refinerResult.Status {
			case tryonTaskStatusProcessing:
				updates := map[string]interface{}{
					"provider_task_id": strings.TrimSpace(refinerResult.ProviderTaskID),
					"refiner_task_id":  strings.TrimSpace(refinerResult.ProviderTaskID),
					"refiner_status":   tryonRefinerStatusProcessing,
					"result_image":     storedMainResultImage,
					"error_message":    "",
				}
				mergeTryonTaskUpdates(updates, parsingStageUpdates)
				updateResult := global.GVA_DB.Model(&client.TryonTask{}).
					Where("id = ? AND user_id = ? AND status = ?", task.ID, userID, tryonTaskStatusProcessing).
					Updates(updates)
				if updateResult.Error != nil {
					return task, false, updateResult.Error
				}
				if updateResult.RowsAffected == 0 {
					return task, false, errors.New("tryonTaskStatusConflictUpdateRefiner")
				}
			case tryonTaskStatusSuccess:
				storedResultImage := s.persistTryonResultImage(strings.TrimSpace(refinerResult.ResultImage), task.TaskNo)
				completeAt := time.Now()
				updates := map[string]interface{}{
					"status":           tryonTaskStatusSuccess,
					"provider_task_id": strings.TrimSpace(refinerResult.ProviderTaskID),
					"refiner_task_id":  strings.TrimSpace(refinerResult.ProviderTaskID),
					"refiner_status":   tryonRefinerStatusSuccess,
					"result_image":     storedResultImage,
					"error_message":    "",
					"completed_at":     completeAt,
				}
				mergeTryonTaskUpdates(updates, parsingStageUpdates)
				updateResult := global.GVA_DB.Model(&client.TryonTask{}).
					Where("id = ? AND user_id = ? AND status = ?", task.ID, userID, tryonTaskStatusProcessing).
					Updates(updates)
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
				fallbackErr := s.markTryonRefinerFailedAndKeepSuccess(ctx, &task, storedMainResultImage, failReason, strings.TrimSpace(refinerResult.ProviderTaskID), parsingStageUpdates)
				if fallbackErr != nil {
					global.GVA_LOG.Error("试衣精修失败回退主图且退币异常", zap.Error(fallbackErr), zap.Uint("taskID", task.ID))
					return task, false, errors.New("tryonRefinerRefundFailedContactAdmin")
				}
			default:
				fallbackErr := s.markTryonRefinerFailedAndKeepSuccess(ctx, &task, storedMainResultImage, "tryonRefinerStatusUnknown", strings.TrimSpace(refinerResult.ProviderTaskID), parsingStageUpdates)
				if fallbackErr != nil {
					global.GVA_LOG.Error("试衣精修未知状态回退主图且退币异常", zap.Error(fallbackErr), zap.Uint("taskID", task.ID))
					return task, false, errors.New("tryonRefinerRefundFailedContactAdmin")
				}
			}
		} else {
			storedResultImage := s.persistTryonResultImage(strings.TrimSpace(invokeResult.ResultImage), task.TaskNo)
			completeAt := time.Now()
			updates := map[string]interface{}{
				"status":           tryonTaskStatusSuccess,
				"provider_task_id": strings.TrimSpace(invokeResult.ProviderTaskID),
				"refiner_status":   tryonRefinerStatusDisabled,
				"refiner_task_id":  "",
				"result_image":     storedResultImage,
				"error_message":    "",
				"completed_at":     completeAt,
			}
			mergeTryonTaskUpdates(updates, parsingStageUpdates)
			updateResult := global.GVA_DB.Model(&client.TryonTask{}).
				Where("id = ? AND user_id = ? AND status = ?", task.ID, userID, tryonTaskStatusProcessing).
				Updates(updates)
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
	if refundErr := s.refundParsingIfNeeded(ctx, &task); refundErr != nil {
		global.GVA_LOG.Error("试衣分割增强失败退币异常", zap.Error(refundErr), zap.Uint("taskID", task.ID))
		return task, false, errors.New("tryonParsingRefundFailedContactAdmin")
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

	beautifyTraceCtx := createBeautifyTraceContext(userID, task)
	resultImage, invokeErr := s.invokeTryonBeautifyProvider(task.SceneType, task.ResultImage, modelCfg, req, beautifyTraceCtx)
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

// RefreshProcessingTryonTasks 批量刷新processing试衣任务状态（定时补偿）
func (s *TryonTaskService) RefreshProcessingTryonTasks(ctx context.Context, db *gorm.DB, limit int) (scanned int, changed int, err error) {
	if db == nil {
		db = global.GVA_DB
	}
	if db == nil {
		return 0, 0, errors.New("dbNotInitialized")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if limit <= 0 {
		limit = 60
	}
	if limit > 500 {
		limit = 500
	}

	var list []client.TryonTask
	err = db.Model(&client.TryonTask{}).
		Where("status = ? AND provider_task_id <> ''", tryonTaskStatusProcessing).
		Order("id asc").
		Limit(limit).
		Find(&list).Error
	if err != nil {
		return 0, 0, err
	}

	scanned = len(list)
	for i := range list {
		beforeStatus := strings.ToLower(strings.TrimSpace(list[i].Status))
		beforeRefinerStatus := strings.ToLower(strings.TrimSpace(list[i].RefinerStatus))

		if refreshErr := s.refreshTryonTaskStatus(ctx, &list[i]); refreshErr != nil {
			global.GVA_LOG.Warn("定时补偿刷新试衣任务状态失败", zap.Error(refreshErr), zap.Uint("taskID", list[i].ID))
			continue
		}

		afterStatus := strings.ToLower(strings.TrimSpace(list[i].Status))
		afterRefinerStatus := strings.ToLower(strings.TrimSpace(list[i].RefinerStatus))
		if beforeStatus != afterStatus || beforeRefinerStatus != afterRefinerStatus {
			changed++
		}
	}

	return
}

// GetTryonTaskStats 获取管理端试衣任务统计（按筛选条件）
func (s *TryonTaskService) GetTryonTaskStats(info clientReq.TryonTaskSearch) (stats TryonTaskStatsData, err error) {
	queryInfo := info
	queryInfo.Status = ""

	latestIDSubQuery := s.buildLatestTryonStatsEventIDSubQuery(queryInfo)
	buildLatestQuery := func() *gorm.DB {
		return s.applyTryonStatsEventFilters(
			global.GVA_DB.Table("client_tryon_stats_event as tse").Where("tse.id IN (?)", latestIDSubQuery),
			queryInfo,
			"tse",
			false,
		)
	}

	type tryonTaskStatsSummaryRow struct {
		Total             int64 `json:"total"`
		Processing        int64 `json:"processing"`
		Success           int64 `json:"success"`
		Failed            int64 `json:"failed"`
		TotalCostPoints   int64 `json:"totalCostPoints"`
		RefinerUsageCount int64 `json:"refinerUsageCount"`
		ParsingUsageCount int64 `json:"parsingUsageCount"`
	}
	var summary tryonTaskStatsSummaryRow
	if err = buildLatestQuery().
		Select("COALESCE(SUM(CASE WHEN COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon') = 'tryon' THEN 1 ELSE 0 END), 0) as total, COALESCE(SUM(CASE WHEN COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon') = 'tryon' AND tse.status = 'processing' THEN 1 ELSE 0 END), 0) as processing, COALESCE(SUM(CASE WHEN COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon') = 'tryon' AND tse.status = 'success' THEN 1 ELSE 0 END), 0) as success, COALESCE(SUM(CASE WHEN COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon') = 'tryon' AND (tse.status = 'failed' OR tse.status = 'error') THEN 1 ELSE 0 END), 0) as failed, COALESCE(SUM(CASE WHEN tse.status = 'success' THEN GREATEST(COALESCE(tse.cost_points, 0) - COALESCE(tse.refund_points, 0), 0) ELSE 0 END), 0) as total_cost_points, COALESCE(SUM(CASE WHEN COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon') = 'refiner' THEN 1 ELSE 0 END), 0) as refiner_usage_count, COALESCE(SUM(CASE WHEN COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon') = 'parsing' THEN 1 ELSE 0 END), 0) as parsing_usage_count").
		Scan(&summary).Error; err != nil {
		return stats, err
	}

	stats.Total = summary.Total
	stats.Processing = summary.Processing
	stats.Success = summary.Success
	stats.Failed = summary.Failed
	stats.TotalCostPoints = summary.TotalCostPoints
	stats.RefinerEnabledCount = summary.RefinerUsageCount
	stats.ParsingEnabledCount = summary.ParsingUsageCount

	type tryonTaskModelAggRow struct {
		ModelKey            string `json:"modelKey"`
		ModelUsage          string `json:"modelUsage"`
		Provider            string `json:"provider"`
		TaskCount           int64  `json:"taskCount"`
		SuccessCount        int64  `json:"successCount"`
		FailedCount         int64  `json:"failedCount"`
		ProcessingCount     int64  `json:"processingCount"`
		RefinerEnabledCount int64  `json:"refinerEnabledCount"`
		ParsingCount        int64  `json:"parsingCount"`
		TotalCostPoints     int64  `json:"totalCostPoints"`
	}

	var rows []tryonTaskModelAggRow
	if err = buildLatestQuery().
		Select("COALESCE(NULLIF(TRIM(tse.model_key), ''), 'default') as model_key, COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon') as model_usage, COALESCE(NULLIF(TRIM(tse.provider), ''), 'default') as provider, COALESCE(SUM(CASE WHEN tse.status = 'success' OR tse.status = 'failed' OR tse.status = 'error' THEN 1 ELSE 0 END), 0) as task_count, COALESCE(SUM(CASE WHEN tse.status = 'success' THEN 1 ELSE 0 END), 0) as success_count, COALESCE(SUM(CASE WHEN tse.status = 'failed' OR tse.status = 'error' THEN 1 ELSE 0 END), 0) as failed_count, COALESCE(SUM(CASE WHEN tse.status = 'processing' THEN 1 ELSE 0 END), 0) as processing_count, COALESCE(SUM(CASE WHEN COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon') = 'refiner' THEN 1 ELSE 0 END), 0) as refiner_enabled_count, COALESCE(SUM(CASE WHEN COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon') = 'parsing' THEN 1 ELSE 0 END), 0) as parsing_count, COALESCE(SUM(CASE WHEN tse.status = 'success' THEN GREATEST(COALESCE(tse.cost_points, 0) - COALESCE(tse.refund_points, 0), 0) ELSE 0 END), 0) as total_cost_points").
		Group("COALESCE(NULLIF(TRIM(tse.model_key), ''), 'default'), COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon'), COALESCE(NULLIF(TRIM(tse.provider), ''), 'default')").
		Order("task_count DESC, success_count DESC").
		Scan(&rows).Error; err != nil {
		return stats, err
	}

	modelProviderMap := loadTryonModelProviderMap()
	modelUsageMap := loadTryonModelUsageMap()
	defaultRefinerByProvider := loadDefaultModelKeyByUsageAndProvider("refiner")
	type tryonTaskStatsModelAggKey struct {
		ModelUsage string
		ModelKey   string
		Provider   string
	}
	modelStatsMap := make(map[tryonTaskStatsModelAggKey]*TryonTaskModelStatsItem, len(rows))
	for _, row := range rows {
		modelKey := strings.TrimSpace(row.ModelKey)
		if modelKey == "" {
			modelKey = "default"
		}

		modelUsage := normalizeTryonModelUsageValue(row.ModelUsage)
		modelProvider := strings.TrimSpace(row.Provider)
		if modelProvider == "" || strings.EqualFold(modelProvider, "default") {
			modelProvider = resolveTryonModelProvider(modelKey, modelProviderMap)
		}
		if modelProvider == "" {
			modelProvider = "default"
		}

		if modelUsage == "refiner" {
			modelKey = normalizeRefinerStatsModelKey(modelKey, modelProvider, modelUsageMap, defaultRefinerByProvider)
			resolvedProvider := resolveTryonModelProvider(modelKey, modelProviderMap)
			if strings.TrimSpace(resolvedProvider) != "" {
				modelProvider = resolvedProvider
			}
		}

		aggKey := tryonTaskStatsModelAggKey{
			ModelUsage: strings.ToLower(strings.TrimSpace(modelUsage)),
			ModelKey:   strings.ToLower(strings.TrimSpace(modelKey)),
			Provider:   strings.ToLower(strings.TrimSpace(modelProvider)),
		}
		if item, exists := modelStatsMap[aggKey]; exists {
			item.TaskCount += row.TaskCount
			item.SuccessCount += row.SuccessCount
			item.FailedCount += row.FailedCount
			item.ProcessingCount += row.ProcessingCount
			item.RefinerEnabledCount += row.RefinerEnabledCount
			item.ParsingCount += row.ParsingCount
			item.TotalCostPoints += row.TotalCostPoints
			continue
		}

		item := TryonTaskModelStatsItem{
			ModelKey:            modelKey,
			ModelUsage:          modelUsage,
			Provider:            modelProvider,
			TaskCount:           row.TaskCount,
			SuccessCount:        row.SuccessCount,
			FailedCount:         row.FailedCount,
			ProcessingCount:     row.ProcessingCount,
			RefinerEnabledCount: row.RefinerEnabledCount,
			ParsingCount:        row.ParsingCount,
			TotalCostPoints:     row.TotalCostPoints,
		}
		modelStatsMap[aggKey] = &item
	}

	stats.ModelStats = make([]TryonTaskModelStatsItem, 0, len(modelStatsMap))
	stats.ModelCallTotal = 0
	for _, item := range modelStatsMap {
		stats.ModelStats = append(stats.ModelStats, *item)
		stats.ModelCallTotal += item.TaskCount
	}

	sort.SliceStable(stats.ModelStats, func(i, j int) bool {
		if stats.ModelStats[i].TaskCount == stats.ModelStats[j].TaskCount {
			leftKey := strings.TrimSpace(stats.ModelStats[i].ModelKey)
			rightKey := strings.TrimSpace(stats.ModelStats[j].ModelKey)
			if strings.EqualFold(leftKey, rightKey) {
				return strings.TrimSpace(stats.ModelStats[i].ModelUsage) < strings.TrimSpace(stats.ModelStats[j].ModelUsage)
			}
			return strings.ToLower(leftKey) < strings.ToLower(rightKey)
		}
		return stats.ModelStats[i].TaskCount > stats.ModelStats[j].TaskCount
	})

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

	latestIDSubQuery := s.buildLatestTryonStatsEventIDSubQuery(queryInfo)
	type tryonTaskTrendLogRow struct {
		CreatedAt time.Time `json:"createdAt"`
		Status    string    `json:"status"`
	}
	var logList []tryonTaskTrendLogRow
	err = s.applyTryonStatsEventFilters(
		global.GVA_DB.Table("client_tryon_stats_event as tse").Where("tse.id IN (?)", latestIDSubQuery),
		queryInfo,
		"tse",
		false,
	).
		Where("COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon') = ?", "tryon").
		Select("tse.event_at as created_at", "tse.status").
		Order("tse.event_at asc").
		Scan(&logList).Error
	if err != nil {
		return nil, err
	}

	trendMap := make(map[string]*TryonTaskTrendItem)
	for day := startDay; !day.After(endDay); day = day.AddDate(0, 0, 1) {
		dateKey := day.Format("2006-01-02")
		trendMap[dateKey] = &TryonTaskTrendItem{Date: dateKey}
	}

	for _, logRow := range logList {
		dateKey := logRow.CreatedAt.Format("2006-01-02")
		item, ok := trendMap[dateKey]
		if !ok {
			continue
		}
		item.Total++
		switch strings.ToLower(strings.TrimSpace(logRow.Status)) {
		case tryonTaskStatusProcessing:
			item.Processing++
		case tryonTaskStatusSuccess:
			item.Success++
		case tryonTaskStatusFailed, modelCallStatusError:
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

func (s *TryonTaskService) buildLatestTryonStatsEventIDSubQuery(info clientReq.TryonTaskSearch) *gorm.DB {
	return s.applyTryonStatsEventFilters(
		global.GVA_DB.Table("client_tryon_stats_event as tse").Where("tse.task_id > 0"),
		info,
		"tse",
		false,
	).
		Select("MAX(tse.id) as id").
		Group("tse.task_id, COALESCE(NULLIF(TRIM(tse.model_usage), ''), 'tryon')")
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

func (s *TryonTaskService) applyTryonStatsEventFilters(db *gorm.DB, info clientReq.TryonTaskSearch, alias string, includeStatus bool) *gorm.DB {
	prefix := strings.TrimSpace(alias)
	column := func(name string) string {
		if prefix == "" {
			return name
		}
		return prefix + "." + name
	}

	if info.UserID > 0 {
		db = db.Where(column("user_id")+" = ?", info.UserID)
	}
	if strings.TrimSpace(info.TaskNo) != "" {
		db = db.Where(column("task_no")+" = ?", strings.TrimSpace(info.TaskNo))
	}
	if strings.TrimSpace(info.RequestID) != "" {
		db = db.Where(column("request_id")+" = ?", strings.TrimSpace(info.RequestID))
	}
	if strings.TrimSpace(info.SceneType) != "" {
		db = db.Where(column("scene_type")+" = ?", strings.TrimSpace(info.SceneType))
	}
	if includeStatus && strings.TrimSpace(info.Status) != "" {
		status := strings.ToLower(strings.TrimSpace(info.Status))
		switch status {
		case tryonTaskStatusFailed, modelCallStatusError:
			db = db.Where(column("status")+" IN ?", []string{tryonTaskStatusFailed, modelCallStatusError})
		default:
			db = db.Where(column("status")+" = ?", status)
		}
	}
	if info.StartCreatedAt != nil {
		db = db.Where(column("event_at")+" >= ?", *info.StartCreatedAt)
	}
	if info.EndCreatedAt != nil {
		db = db.Where(column("event_at")+" <= ?", *info.EndCreatedAt)
	}

	usageFilterExpr := fmt.Sprintf("COALESCE(NULLIF(TRIM(%s), ''), 'tryon') IN ?", column("model_usage"))
	db = db.Where(usageFilterExpr, []string{"tryon", "refiner", "parsing", "beautify"})
	return db
}

func mergeTryonTaskUpdates(base map[string]interface{}, patch map[string]interface{}) {
	if len(base) == 0 || len(patch) == 0 {
		return
	}
	for key, value := range patch {
		base[key] = value
	}
}

func buildParsingStageUpdates(task client.TryonTask, invokeResult tryonInvokeResult) map[string]interface{} {
	updates := map[string]interface{}{}
	if !task.EnableParsing {
		updates["parsing_status"] = tryonParsingStatusDisabled
		updates["parsing_model_key"] = ""
		updates["parsing_error"] = ""
		return updates
	}

	modelKey := strings.TrimSpace(task.ParsingModelKey)
	if strings.TrimSpace(invokeResult.ParsingModelKey) != "" {
		modelKey = strings.TrimSpace(invokeResult.ParsingModelKey)
	}
	updates["parsing_model_key"] = modelKey

	if !invokeResult.ParsingAttempted {
		updates["parsing_status"] = tryonParsingStatusDisabled
		updates["parsing_error"] = ""
		return updates
	}

	if invokeResult.ParsingSucceeded {
		updates["parsing_status"] = tryonParsingStatusSuccess
		updates["parsing_error"] = ""
		return updates
	}

	updates["parsing_status"] = tryonParsingStatusFailed
	updates["parsing_error"] = truncateTryonError(firstNonEmptyString(invokeResult.ParsingError, "aliyunParsingNoUsableGarmentArea"))
	return updates
}

func buildParsingPointRemark(taskNo string, parsingModelKey string) string {
	remark := "tryon_parsing:" + strings.TrimSpace(taskNo)
	if strings.TrimSpace(parsingModelKey) != "" {
		remark += ":" + strings.TrimSpace(parsingModelKey)
	}
	return remark
}

func resolveParsingStatus(task *client.TryonTask, extraUpdates map[string]interface{}) string {
	if extraUpdates != nil {
		if raw, ok := extraUpdates["parsing_status"]; ok {
			text := strings.ToLower(strings.TrimSpace(fmt.Sprint(raw)))
			if text != "" {
				return text
			}
		}
	}
	if task == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(task.ParsingStatus))
}

func (s *TryonTaskService) refundParsingIfNeeded(ctx context.Context, task *client.TryonTask) error {
	if task == nil || task.ID == 0 || task.UserID == 0 {
		return nil
	}
	if !task.EnableParsing {
		return nil
	}
	if resolveParsingStatus(task, nil) != tryonParsingStatusFailed {
		return nil
	}
	if task.ParsingCostPoints <= 0 || task.ParsingRefund >= task.ParsingCostPoints {
		return nil
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var latest client.TryonTask
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", task.ID, task.UserID).
			First(&latest).Error; err != nil {
			return err
		}

		if !latest.EnableParsing || resolveParsingStatus(&latest, nil) != tryonParsingStatusFailed {
			*task = latest
			return nil
		}
		if latest.ParsingCostPoints <= 0 || latest.ParsingRefund >= latest.ParsingCostPoints {
			*task = latest
			return nil
		}

		refundPoints := latest.ParsingCostPoints - latest.ParsingRefund
		if refundPoints <= 0 {
			*task = latest
			return nil
		}

		now := time.Now()
		updateData := map[string]interface{}{
			"parsing_refund": latest.ParsingRefund + refundPoints,
			"refund_points":  latest.RefundPoints + refundPoints,
			"refunded_at":    now,
		}
		if err := tx.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ?", latest.ID, latest.UserID).
			Updates(updateData).Error; err != nil {
			return err
		}

		pointRecordService := PointRecordService{}
		txCtx := context.WithValue(ctx, "tx", tx)
		record := buildPointRecord(latest.UserID, client.AssetTypeTryonPoint, "increase", refundPoints, tryonOpParsingRefund, tryonReasonTaskRefund, buildParsingPointRemark(latest.TaskNo, latest.ParsingModelKey))
		if err := pointRecordService.CreatePointRecord(txCtx, record); err != nil {
			return err
		}

		latest.ParsingRefund += refundPoints
		latest.RefundPoints += refundPoints
		latest.RefundedAt = &now
		*task = latest
		return nil
	})
}

func (s *TryonTaskService) markFailedAndRefund(ctx context.Context, task *client.TryonTask, failReason string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		failReason = truncateTryonError(failReason)

		baseRefundPoints := task.BaseCostPoints
		parsingRefundPoints := task.ParsingCostPoints
		refinerRefundPoints := task.RefinerCostPoints
		if baseRefundPoints < 0 {
			baseRefundPoints = 0
		}
		if parsingRefundPoints < 0 {
			parsingRefundPoints = 0
		}
		if refinerRefundPoints < 0 {
			refinerRefundPoints = 0
		}
		totalRefundPoints := baseRefundPoints + parsingRefundPoints + refinerRefundPoints
		if totalRefundPoints <= 0 {
			totalRefundPoints = task.CostPoints
			baseRefundPoints = totalRefundPoints
			parsingRefundPoints = 0
			refinerRefundPoints = 0
		}

		updates := map[string]interface{}{
			"status":        tryonTaskStatusFailed,
			"error_message": failReason,
			"refund_points": totalRefundPoints,
			"completed_at":  now,
			"refunded_at":   now,
		}
		if task.EnableRefiner {
			updates["refiner_status"] = tryonRefinerStatusFailed
		}
		if parsingRefundPoints > 0 {
			updates["parsing_refund"] = parsingRefundPoints
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
		if baseRefundPoints > 0 {
			record := buildPointRecord(task.UserID, client.AssetTypeTryonPoint, "increase", baseRefundPoints, tryonOpRefund, tryonReasonTaskRefund, "tryon_task:"+task.TaskNo)
			if err := pointRecordService.CreatePointRecord(txCtx, record); err != nil {
				return err
			}
		}
		if refinerRefundPoints > 0 {
			remark := "tryon_refiner:" + task.TaskNo
			if strings.TrimSpace(task.RefinerModelKey) != "" {
				remark = remark + ":" + strings.TrimSpace(task.RefinerModelKey)
			}
			record := buildPointRecord(task.UserID, client.AssetTypeTryonPoint, "increase", refinerRefundPoints, tryonOpRefinerRefund, tryonReasonTaskRefund, remark)
			if err := pointRecordService.CreatePointRecord(txCtx, record); err != nil {
				return err
			}
		}
		if parsingRefundPoints > 0 {
			record := buildPointRecord(task.UserID, client.AssetTypeTryonPoint, "increase", parsingRefundPoints, tryonOpParsingRefund, tryonReasonTaskRefund, buildParsingPointRemark(task.TaskNo, task.ParsingModelKey))
			if err := pointRecordService.CreatePointRecord(txCtx, record); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TryonTaskService) markTryonRefinerFailedAndKeepSuccess(ctx context.Context, task *client.TryonTask, mainResultImage string, failReason string, refinerTaskID string, extraUpdates map[string]interface{}) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		failReason = truncateTryonError(failReason)

		storedResultImage := strings.TrimSpace(mainResultImage)
		if storedResultImage == "" {
			storedResultImage = strings.TrimSpace(task.ResultImage)
		}
		if storedResultImage == "" {
			return errors.New("tryonResultImageRequired")
		}

		refinerRefundPoints := task.RefinerCostPoints
		if refinerRefundPoints < 0 {
			refinerRefundPoints = 0
		}
		parsingRefundPoints := 0
		if task.EnableParsing && resolveParsingStatus(task, extraUpdates) == tryonParsingStatusFailed {
			parsingRefundPoints = task.ParsingCostPoints - task.ParsingRefund
			if parsingRefundPoints < 0 {
				parsingRefundPoints = 0
			}
		}

		updates := map[string]interface{}{
			"status":          tryonTaskStatusSuccess,
			"result_image":    storedResultImage,
			"refiner_status":  tryonRefinerStatusFailed,
			"refiner_task_id": strings.TrimSpace(refinerTaskID),
			"refiner_refund":  refinerRefundPoints,
			"error_message":   failReason,
			"completed_at":    now,
		}
		totalExtraRefund := refinerRefundPoints + parsingRefundPoints
		if totalExtraRefund > 0 {
			updates["refund_points"] = task.RefundPoints + totalExtraRefund
			updates["refunded_at"] = now
		}
		if parsingRefundPoints > 0 {
			updates["parsing_refund"] = task.ParsingRefund + parsingRefundPoints
		}
		if strings.TrimSpace(refinerTaskID) != "" {
			updates["provider_task_id"] = strings.TrimSpace(refinerTaskID)
		}
		mergeTryonTaskUpdates(updates, extraUpdates)

		updateResult := tx.Model(&client.TryonTask{}).
			Where("id = ? AND user_id = ? AND status = ?", task.ID, task.UserID, tryonTaskStatusProcessing).
			Updates(updates)
		if updateResult.Error != nil {
			return updateResult.Error
		}
		if updateResult.RowsAffected == 0 {
			return errors.New("tryonTaskStatusConflictRefund")
		}

		if refinerRefundPoints > 0 {
			pointRecordService := PointRecordService{}
			txCtx := context.WithValue(ctx, "tx", tx)
			remark := "tryon_refiner:" + task.TaskNo
			if strings.TrimSpace(task.RefinerModelKey) != "" {
				remark = remark + ":" + strings.TrimSpace(task.RefinerModelKey)
			}
			record := buildPointRecord(task.UserID, client.AssetTypeTryonPoint, "increase", refinerRefundPoints, tryonOpRefinerRefund, tryonReasonTaskRefund, remark)
			if err := pointRecordService.CreatePointRecord(txCtx, record); err != nil {
				return err
			}
		}
		if parsingRefundPoints > 0 {
			pointRecordService := PointRecordService{}
			txCtx := context.WithValue(ctx, "tx", tx)
			record := buildPointRecord(task.UserID, client.AssetTypeTryonPoint, "increase", parsingRefundPoints, tryonOpParsingRefund, tryonReasonTaskRefund, buildParsingPointRemark(task.TaskNo, task.ParsingModelKey))
			if err := pointRecordService.CreatePointRecord(txCtx, record); err != nil {
				return err
			}
		}

		task.Status = tryonTaskStatusSuccess
		task.ResultImage = storedResultImage
		task.RefinerStatus = tryonRefinerStatusFailed
		task.RefinerTaskID = strings.TrimSpace(refinerTaskID)
		task.RefinerRefund = refinerRefundPoints
		task.ErrorMessage = failReason
		task.CompletedAt = &now
		if totalExtraRefund > 0 {
			task.RefundPoints += totalExtraRefund
			task.RefundedAt = &now
		}
		if parsingRefundPoints > 0 {
			task.ParsingRefund += parsingRefundPoints
		}
		if strings.TrimSpace(refinerTaskID) != "" {
			task.ProviderTaskID = strings.TrimSpace(refinerTaskID)
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

func (s *TryonTaskService) loadTryonModelConfigList(_ string) ([]tryonModelConfig, error) {
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

func findEnabledRefinerModelByKey(list []tryonModelConfig, sceneType string, modelKey string) *tryonModelConfig {
	modelKey = strings.TrimSpace(modelKey)
	if modelKey == "" {
		return nil
	}

	for i := range list {
		item := &list[i]
		if !item.isEnabled() || !item.isRefinerModelUsage() {
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
		if !item.isEnabled() || !item.isRefinerModelUsage() {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(item.Key), modelKey) {
			return item
		}
	}

	return nil
}

func (s *TryonTaskService) resolveRefinerModelConfig(sceneType string, baseModelCfg *tryonModelConfig) (*tryonModelConfig, error) {
	if baseModelCfg == nil {
		return nil, nil
	}

	refinerModelKey := strings.TrimSpace(baseModelCfg.RefinerModelKey)
	if refinerModelKey == "" {
		// Backward compatibility for legacy embedded refiner configs.
		hasLegacyEmbeddedRefiner := strings.TrimSpace(baseModelCfg.RefinerModel) != "" ||
			strings.TrimSpace(baseModelCfg.RefinerURL) != "" ||
			strings.TrimSpace(baseModelCfg.RefinerToken) != "" ||
			strings.TrimSpace(baseModelCfg.RefinerTaskQueryURL) != "" ||
			baseModelCfg.RefinerExtraCost > 0 ||
			baseModelCfg.RefinerExtraPoints > 0
		if hasLegacyEmbeddedRefiner && baseModelCfg.supportsRefinerValue() {
			return baseModelCfg, nil
		}
		return nil, errors.New("tryonRefinerModelUnavailable")
	}

	list, err := s.loadTryonModelConfigList(sceneType)
	if err != nil {
		return nil, errors.New("tryonModelConfigInvalid")
	}
	if len(list) == 0 {
		return nil, errors.New("tryonRefinerModelUnavailable")
	}

	refinerModelCfg := findEnabledRefinerModelByKey(list, sceneType, refinerModelKey)
	if refinerModelCfg == nil {
		return nil, errors.New("tryonRefinerModelUnavailable")
	}

	return refinerModelCfg, nil
}

func findEnabledTryonModelByKey(list []tryonModelConfig, sceneType string, modelKey string) *tryonModelConfig {
	modelKey = strings.TrimSpace(modelKey)
	if modelKey == "" {
		return nil
	}

	for i := range list {
		item := &list[i]
		if !item.isEnabled() || !item.isTryonModelUsage() {
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
		if !item.isEnabled() || !item.isTryonModelUsage() {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(item.Key), modelKey) {
			return item
		}
	}

	return nil
}

func (s *TryonTaskService) resolveBeautifyModelConfig(sceneType string, tryonModelKey string, preferredBeautifyModelKey string) (*tryonModelConfig, error) {
	list, err := s.loadTryonModelConfigList(sceneType)
	if err != nil {
		return nil, errors.New("tryonModelConfigInvalid")
	}
	if len(list) == 0 {
		return nil, nil
	}

	_ = strings.TrimSpace(tryonModelKey)

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

func (s *TryonTaskService) invokeTryonProvider(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, traceCtx *modelCallTraceContext) (tryonInvokeResult, error) {
	if modelCfg == nil {
		return tryonInvokeResult{}, errors.New("tryonModelUnavailable")
	}

	providerMode := strings.TrimSpace(modelCfg.Mode)
	if providerMode == "" {
		providerMode = "prod"
	}
	if strings.EqualFold(providerMode, "mock_success") {
		return s.invokeTryonProviderWithToken(req, modelCfg, traceCtx, "")
	}

	providerURL := strings.TrimSpace(modelCfg.endpointURL())
	modelName := resolveAliyunModelName(modelCfg, req.SceneType)
	providerName := strings.TrimSpace(modelCfg.Provider)
	isAliyunModel := isAliyunParsingModel(providerName, modelName, providerURL, req.SceneType) ||
		isAliyunTryonModel(providerName, modelName, providerURL, req.SceneType)
	if !isAliyunModel {
		return s.invokeTryonProviderWithToken(req, modelCfg, traceCtx, "")
	}

	tokenCandidates := modelCfg.authTokenCandidates()
	if len(tokenCandidates) == 0 {
		return s.invokeTryonProviderWithToken(req, modelCfg, traceCtx, "")
	}

	for index, token := range tokenCandidates {
		result, err := s.invokeTryonProviderWithToken(req, modelCfg, traceCtx, token)
		if shouldSwitchTryonModelToken(err, result) && index < len(tokenCandidates)-1 {
			global.GVA_LOG.Warn("试衣模型鉴权失败，自动切换备用token重试",
				zap.String("modelKey", strings.TrimSpace(modelCfg.Key)),
				zap.String("sceneType", strings.TrimSpace(req.SceneType)),
				zap.Int("attempt", index+1),
				zap.Int("total", len(tokenCandidates)),
				zap.String("reason", tryonTokenSwitchReason(err, result)),
			)
			continue
		}
		return result, err
	}

	return s.invokeTryonProviderWithToken(req, modelCfg, traceCtx, tokenCandidates[len(tokenCandidates)-1])
}

func (s *TryonTaskService) invokeTryonProviderWithToken(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, traceCtx *modelCallTraceContext, overrideToken string) (tryonInvokeResult, error) {
	if modelCfg == nil {
		return tryonInvokeResult{}, errors.New("tryonModelUnavailable")
	}

	providerMode := strings.TrimSpace(modelCfg.Mode)
	if providerMode == "" {
		providerMode = "prod"
	}
	providerURL := strings.TrimSpace(modelCfg.endpointURL())
	providerToken := strings.TrimSpace(overrideToken)
	if providerToken == "" {
		providerToken = strings.TrimSpace(modelCfg.authToken())
	}

	if strings.EqualFold(providerMode, "mock_success") {
		return tryonInvokeResult{
			Status:      tryonTaskStatusSuccess,
			ResultImage: req.SourceImage,
		}, nil
	}

	modelName := resolveAliyunModelName(modelCfg, req.SceneType)
	providerName := strings.TrimSpace(modelCfg.Provider)

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

			if strings.TrimSpace(req.TemplateImageLower) != "" {
				normalizedTemplateLower, normalizeErr := s.normalizeAliyunMediaURL(req.TemplateImageLower)
				if normalizeErr != nil {
					return tryonInvokeResult{}, normalizeErr
				}
				resolvedReq.TemplateImageLower = normalizedTemplateLower
			}
		}
	}

	if isAliyunParsing {
		return s.invokeAliyunParsing(resolvedReq, modelCfg, providerToken, providerURL, modelName, false, traceCtx)
	}

	parsedCompanionGarmentImage := ""
	parsingAttempted := false
	parsingSucceeded := false
	parsingErrorMessage := ""
	parsingModelKey := strings.TrimSpace(modelCfg.ParsingModelKey)
	if isAliyunTryon && shouldEnableAliyunParsingForSinglePart(resolvedReq.SceneType, modelName, resolvedReq.TemplatePart, req.EnableParsing) {
		parsingAttempted = true
		parsedCompanionGarment, resolvedParsingModelKey, parseErr := s.tryBuildAliyunParsedCompanionGarment(resolvedReq, providerToken, modelCfg, traceCtx)
		if strings.TrimSpace(resolvedParsingModelKey) != "" {
			parsingModelKey = strings.TrimSpace(resolvedParsingModelKey)
		}
		if parseErr != nil {
			parsingErrorMessage = parseErr.Error()
			global.GVA_LOG.Warn("自动阿里取衣分割失败，回退仅使用原始单件模板图继续试衣",
				zap.String("modelKey", strings.TrimSpace(modelCfg.Key)),
				zap.String("modelName", strings.TrimSpace(modelName)),
				zap.String("templatePart", strings.TrimSpace(resolvedReq.TemplatePart)),
				zap.Error(parseErr),
			)
		} else if strings.TrimSpace(parsedCompanionGarment) != "" {
			parsedCompanionGarmentImage = strings.TrimSpace(parsedCompanionGarment)
			parsingSucceeded = true
		}
	}

	if isAliyunTryon {
		result, invokeErr := s.invokeAliyunTryonAsync(resolvedReq, modelCfg, providerToken, providerURL, modelName, parsedCompanionGarmentImage, traceCtx)
		if invokeErr != nil {
			return tryonInvokeResult{}, invokeErr
		}
		result.ParsingAttempted = parsingAttempted
		result.ParsingSucceeded = parsingSucceeded
		result.ParsingModelKey = strings.TrimSpace(parsingModelKey)
		result.ParsingError = strings.TrimSpace(parsingErrorMessage)
		return result, nil
	}

	if isGradioTryon {
		return s.invokeGradioTryon(resolvedReq, modelCfg, providerToken, providerURL, traceCtx)
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

	startedAt := time.Now()
	resp, err := external.HttpRequest[tryonProviderResponse](external.RequestParams{
		URL:     strings.TrimSpace(providerURL),
		Method:  "POST",
		Headers: headers,
		Body:    body,
	})
	if err != nil {
		saveModelCallLog(modelCallLogInput{
			TraceContext:   traceCtx,
			CallStage:      "tryon_custom_provider",
			ModelKey:       strings.TrimSpace(modelCfg.Key),
			ModelUsage:     modelCfg.usageValue(),
			ModelName:      strings.TrimSpace(modelName),
			Provider:       strings.TrimSpace(providerName),
			EndpointURL:    strings.TrimSpace(providerURL),
			RelatedModels:  map[string]string{"refinerModelKey": strings.TrimSpace(modelCfg.RefinerModelKey), "parsingModelKey": strings.TrimSpace(modelCfg.ParsingModelKey)},
			RequestPayload: body,
			Status:         modelCallStatusError,
			ErrorMessage:   err.Error(),
			StartedAt:      startedAt,
			FinishedAt:     time.Now(),
		})
		return tryonInvokeResult{}, err
	}
	if resp == nil {
		saveModelCallLog(modelCallLogInput{
			TraceContext:   traceCtx,
			CallStage:      "tryon_custom_provider",
			ModelKey:       strings.TrimSpace(modelCfg.Key),
			ModelUsage:     modelCfg.usageValue(),
			ModelName:      strings.TrimSpace(modelName),
			Provider:       strings.TrimSpace(providerName),
			EndpointURL:    strings.TrimSpace(providerURL),
			RelatedModels:  map[string]string{"refinerModelKey": strings.TrimSpace(modelCfg.RefinerModelKey), "parsingModelKey": strings.TrimSpace(modelCfg.ParsingModelKey)},
			RequestPayload: body,
			Status:         modelCallStatusError,
			ErrorMessage:   "tryonModelResponseEmpty",
			StartedAt:      startedAt,
			FinishedAt:     time.Now(),
		})
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
			saveModelCallLog(modelCallLogInput{
				TraceContext:    traceCtx,
				CallStage:       "tryon_custom_provider",
				ModelKey:        strings.TrimSpace(modelCfg.Key),
				ModelUsage:      modelCfg.usageValue(),
				ModelName:       strings.TrimSpace(modelName),
				Provider:        strings.TrimSpace(providerName),
				EndpointURL:     strings.TrimSpace(providerURL),
				RelatedModels:   map[string]string{"refinerModelKey": strings.TrimSpace(modelCfg.RefinerModelKey), "parsingModelKey": strings.TrimSpace(modelCfg.ParsingModelKey)},
				RequestPayload:  body,
				ResponsePayload: resp,
				Status:          modelCallStatusError,
				ErrorMessage:    strings.TrimSpace(resp.Msg),
				StartedAt:       startedAt,
				FinishedAt:      time.Now(),
			})
			return tryonInvokeResult{}, errors.New(strings.TrimSpace(resp.Msg))
		}
		saveModelCallLog(modelCallLogInput{
			TraceContext:    traceCtx,
			CallStage:       "tryon_custom_provider",
			ModelKey:        strings.TrimSpace(modelCfg.Key),
			ModelUsage:      modelCfg.usageValue(),
			ModelName:       strings.TrimSpace(modelName),
			Provider:        strings.TrimSpace(providerName),
			EndpointURL:     strings.TrimSpace(providerURL),
			RelatedModels:   map[string]string{"refinerModelKey": strings.TrimSpace(modelCfg.RefinerModelKey), "parsingModelKey": strings.TrimSpace(modelCfg.ParsingModelKey)},
			RequestPayload:  body,
			ResponsePayload: resp,
			Status:          modelCallStatusError,
			ErrorMessage:    "tryonModelResultEmpty",
			StartedAt:       startedAt,
			FinishedAt:      time.Now(),
		})
		return tryonInvokeResult{}, errors.New("tryonModelResultEmpty")
	}

	saveModelCallLog(modelCallLogInput{
		TraceContext:    traceCtx,
		CallStage:       "tryon_custom_provider",
		ModelKey:        strings.TrimSpace(modelCfg.Key),
		ModelUsage:      modelCfg.usageValue(),
		ModelName:       strings.TrimSpace(modelName),
		Provider:        strings.TrimSpace(providerName),
		EndpointURL:     strings.TrimSpace(providerURL),
		RelatedModels:   map[string]string{"refinerModelKey": strings.TrimSpace(modelCfg.RefinerModelKey), "parsingModelKey": strings.TrimSpace(modelCfg.ParsingModelKey)},
		RequestPayload:  body,
		ResponsePayload: resp,
		ResultImage:     strings.TrimSpace(resultImage),
		Status:          tryonTaskStatusSuccess,
		StartedAt:       startedAt,
		FinishedAt:      time.Now(),
	})

	return tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage}, nil
}

func (s *TryonTaskService) invokeGradioTryon(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, providerToken string, providerURL string, traceCtx *modelCallTraceContext) (result tryonInvokeResult, err error) {
	startedAt := time.Now()
	var requestPayload interface{}
	var responsePayload interface{}
	defer func() {
		status := strings.TrimSpace(result.Status)
		if status == "" {
			if err != nil {
				status = modelCallStatusError
			} else {
				status = tryonTaskStatusUnknown
			}
		}

		saveModelCallLog(modelCallLogInput{
			TraceContext:    traceCtx,
			CallStage:       "tryon_gradio",
			ModelKey:        strings.TrimSpace(modelCfg.Key),
			ModelUsage:      modelCfg.usageValue(),
			ModelName:       strings.TrimSpace(modelCfg.Model),
			Provider:        strings.TrimSpace(modelCfg.Provider),
			EndpointURL:     strings.TrimSpace(providerURL),
			RelatedModels:   map[string]string{"refinerModelKey": strings.TrimSpace(modelCfg.RefinerModelKey), "parsingModelKey": strings.TrimSpace(modelCfg.ParsingModelKey)},
			RequestPayload:  requestPayload,
			ResponsePayload: responsePayload,
			ResultImage:     strings.TrimSpace(result.ResultImage),
			Status:          status,
			ErrorMessage:    errorToString(err),
			StartedAt:       startedAt,
			FinishedAt:      time.Now(),
		})
	}()

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
	requestPayload = body

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
	responsePayload = map[string]interface{}{
		"statusCode": response.StatusCode,
		"body":       string(responseBody),
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
			result = tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage}
			return result, nil
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

	result = tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage, ProviderTaskID: eventID}
	return result, nil
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
	if part == "upper" || part == "lower" || part == "upper_lower" || part == "dress" {
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
	if part == "upper_lower" {
		return []string{"upper", "lower"}
	}
	if part == "dress" {
		return []string{"upper", "lower"}
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

func shouldEnableAliyunParsingForSinglePart(sceneType string, modelName string, templatePart string, enableParsing bool) bool {
	if !enableParsing || !strings.EqualFold(strings.TrimSpace(sceneType), "clothes") {
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
	return part == "upper" || part == "lower" || part == "upper_lower" || part == "dress"
}

func isAliyunParsingAssistModel(item *tryonModelConfig) bool {
	if item == nil {
		return false
	}
	if item.isParsingModelUsage() {
		return true
	}
	// Backward compatibility: old configs may still mark parsing entries as tryon.
	itemModelName := resolveAliyunModelName(item, "")
	return isAliyunParsingModel(item.Provider, itemModelName, item.endpointURL(), "")
}

func findEnabledAliyunParsingModelByKey(list []tryonModelConfig, sceneType string, modelKey string) *tryonModelConfig {
	modelKey = strings.TrimSpace(modelKey)
	if modelKey == "" {
		return nil
	}

	for i := range list {
		item := &list[i]
		if !isAliyunParsingAssistModel(item) {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(item.Key), modelKey) {
			continue
		}
		if item.isEnabled() && item.matchScene(sceneType) {
			return item
		}
	}

	for i := range list {
		item := &list[i]
		if !isAliyunParsingAssistModel(item) {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(item.Key), modelKey) {
			return item
		}
	}

	return nil
}

func (s *TryonTaskService) resolveAliyunParsingAssistModelConfig(sceneType string, preferredParsingModelKey string) (*tryonModelConfig, error) {
	list, err := s.loadTryonModelConfigList(sceneType)
	if err != nil {
		return nil, errors.New("tryonModelConfigInvalid")
	}
	if len(list) == 0 {
		return nil, nil
	}

	preferredParsingModelKey = strings.TrimSpace(preferredParsingModelKey)
	if preferredParsingModelKey != "" {
		if preferred := findEnabledAliyunParsingModelByKey(list, sceneType, preferredParsingModelKey); preferred != nil {
			return preferred, nil
		}
		return nil, errors.New("aliyunParsingModelUnavailable")
	}

	var enabledSceneMatched *tryonModelConfig
	var enabledFallback *tryonModelConfig
	var configuredSceneMatched *tryonModelConfig
	var configuredFallback *tryonModelConfig

	for i := range list {
		item := &list[i]
		if !isAliyunParsingAssistModel(item) {
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

func (s *TryonTaskService) tryBuildAliyunParsedCompanionGarment(req clientReq.CreateTryonTaskReq, fallbackToken string, currentModelCfg *tryonModelConfig, traceCtx *modelCallTraceContext) (string, string, error) {
	part := normalizeTryonTemplatePart(req.TemplatePart)
	if part == "" {
		return "", "", nil
	}

	preferredParsingModelKey := ""
	if currentModelCfg != nil {
		preferredParsingModelKey = strings.TrimSpace(currentModelCfg.ParsingModelKey)
	}

	parsingCfg, err := s.resolveAliyunParsingAssistModelConfig(req.SceneType, preferredParsingModelKey)
	if err != nil {
		return "", "", err
	}
	if parsingCfg == nil {
		return "", "", errors.New("aliyunParsingModelUnavailable")
	}

	cfgCopy := *parsingCfg
	parseClothesType := resolveAliyunParsingClothesTypeForTemplatePart(part)
	if len(parseClothesType) == 0 {
		parseClothesType = []string{part}
	}

	attemptClothesTypes := [][]string{parseClothesType}

	parseReq := clientReq.CreateTryonTaskReq{
		SceneType: "takeoff",
		// Parsing assist should segment from the model/source image, then reuse crop_img_url as companion garment image.
		SourceImage: strings.TrimSpace(req.SourceImage),
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

		parseTraceCtx := cloneTraceContext(traceCtx)
		if parseTraceCtx != nil {
			parseTraceCtx.SceneType = "takeoff"
			parseTraceCtx.TemplatePart = ""
			parseTraceCtx.SourceImage = strings.TrimSpace(parseReq.SourceImage)
			parseTraceCtx.TemplateImage = ""
			parseTraceCtx.EntrySource = inferModelCallEntrySource(parseTraceCtx.SceneType, parseTraceCtx.TemplatePart)
			parseTraceCtx.Behavior = inferModelCallBehavior(parseTraceCtx.SceneType, parseTraceCtx.TemplatePart)
		}

		// Keep preferParsingImage=false so crop_img_url is selected first for downstream companion garment image.
		parseResult, parseErr := s.invokeAliyunParsing(parseReq, &cfgCopy, parseProviderToken, parseProviderURL, parseModelName, false, parseTraceCtx)
		if parseErr != nil {
			return "", strings.TrimSpace(cfgCopy.Key), parseErr
		}

		if strings.EqualFold(strings.TrimSpace(parseResult.Status), tryonTaskStatusSuccess) && strings.TrimSpace(parseResult.ResultImage) != "" {
			normalizedCompanion, normalizeErr := s.normalizeAliyunMediaURL(strings.TrimSpace(parseResult.ResultImage))
			if normalizeErr != nil {
				return "", strings.TrimSpace(cfgCopy.Key), normalizeErr
			}
			return strings.TrimSpace(normalizedCompanion), strings.TrimSpace(cfgCopy.Key), nil
		}

		errMsg := strings.TrimSpace(parseResult.ErrorMessage)
		if errMsg == "" {
			errMsg = "aliyunParsingNoUsableGarmentArea"
		}
		lastErr = errors.New(errMsg)

		return "", strings.TrimSpace(cfgCopy.Key), lastErr
	}

	return "", strings.TrimSpace(cfgCopy.Key), lastErr
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

func buildAliyunTryonInput(personImageURL string, templateImageURL string, templateImageLowerURL string, templatePart string, companionGarmentImageURL string) map[string]interface{} {
	input := map[string]interface{}{
		"person_image_url": strings.TrimSpace(personImageURL),
	}

	templateImageURL = strings.TrimSpace(templateImageURL)
	templateImageLowerURL = strings.TrimSpace(templateImageLowerURL)
	companionGarmentImageURL = strings.TrimSpace(companionGarmentImageURL)

	part := normalizeTryonTemplatePart(templatePart)
	if part == "lower" {
		input["bottom_garment_url"] = templateImageURL
		if companionGarmentImageURL != "" {
			input["top_garment_url"] = companionGarmentImageURL
		}
		return input
	}
	if part == "" || part == "upper_lower" {
		if templateImageURL != "" {
			input["top_garment_url"] = templateImageURL
		}
		if templateImageLowerURL != "" {
			input["bottom_garment_url"] = templateImageLowerURL
		}
		if companionGarmentImageURL != "" {
			if _, ok := input["top_garment_url"]; !ok {
				input["top_garment_url"] = companionGarmentImageURL
			}
			if _, ok := input["bottom_garment_url"]; !ok {
				input["bottom_garment_url"] = companionGarmentImageURL
			}
		}
		return input
	}

	input["top_garment_url"] = templateImageURL
	if part == "upper" {
		if templateImageLowerURL != "" {
			input["bottom_garment_url"] = templateImageLowerURL
		} else if companionGarmentImageURL != "" {
			input["bottom_garment_url"] = companionGarmentImageURL
		}
		return input
	}

	if part == "dress" && companionGarmentImageURL != "" {
		input["bottom_garment_url"] = companionGarmentImageURL
	}
	return input
}

func isAliyunShoeTryonModel(sceneType string, modelName string, providerURL string) bool {
	if !strings.EqualFold(strings.TrimSpace(sceneType), "shoes") {
		return false
	}

	if strings.EqualFold(strings.TrimSpace(modelName), aliyunShoeTryonModel) {
		return true
	}

	urlLower := strings.ToLower(strings.TrimSpace(providerURL))
	if strings.Contains(urlLower, "/virtualmodel/generation") {
		return true
	}

	return false
}

func buildAliyunShoeTryonInput(templateImageURL string, shoeImageURL string, shoeImageLowerURL string) map[string]interface{} {
	input := map[string]interface{}{
		"template_image_url": strings.TrimSpace(templateImageURL),
	}

	shoeURLSet := make(map[string]struct{})
	shoeURLs := make([]string, 0, 2)
	for _, candidate := range []string{strings.TrimSpace(shoeImageURL), strings.TrimSpace(shoeImageLowerURL)} {
		if candidate == "" {
			continue
		}
		if _, exists := shoeURLSet[candidate]; exists {
			continue
		}
		shoeURLSet[candidate] = struct{}{}
		shoeURLs = append(shoeURLs, candidate)
	}
	if len(shoeURLs) > 0 {
		input["shoe_image_url"] = shoeURLs
	}

	return input
}

func (s *TryonTaskService) invokeAliyunTryonAsync(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, providerToken string, providerURL string, modelName string, companionGarmentImageURL string, traceCtx *modelCallTraceContext) (result tryonInvokeResult, err error) {
	startedAt := time.Now()
	var requestPayload interface{}
	var responsePayload interface{}
	defer func() {
		status := strings.TrimSpace(result.Status)
		if status == "" {
			if err != nil {
				status = modelCallStatusError
			} else {
				status = tryonTaskStatusUnknown
			}
		}

		saveModelCallLog(modelCallLogInput{
			TraceContext:    traceCtx,
			CallStage:       "tryon_aliyun_async",
			ModelKey:        strings.TrimSpace(modelCfg.Key),
			ModelUsage:      modelCfg.usageValue(),
			ModelName:       strings.TrimSpace(modelName),
			Provider:        strings.TrimSpace(modelCfg.Provider),
			EndpointURL:     strings.TrimSpace(providerURL),
			QueryURL:        strings.TrimSpace(modelCfg.queryTaskURL()),
			RelatedModels:   map[string]string{"refinerModelKey": strings.TrimSpace(modelCfg.RefinerModelKey), "parsingModelKey": strings.TrimSpace(modelCfg.ParsingModelKey)},
			RequestPayload:  requestPayload,
			ResponsePayload: responsePayload,
			ResultImage:     strings.TrimSpace(result.ResultImage),
			Status:          status,
			ErrorMessage:    errorToString(err),
			StartedAt:       startedAt,
			FinishedAt:      time.Now(),
		})
	}()

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

	inputPayload := buildAliyunTryonInput(req.SourceImage, req.TemplateImage, req.TemplateImageLower, req.TemplatePart, companionGarmentImageURL)
	parametersPayload := map[string]interface{}{
		"resolution":   modelCfg.resolutionValue(),
		"restore_face": modelCfg.restoreFaceValue(),
	}
	if isAliyunShoeTryonModel(req.SceneType, resolvedModelName, createURL) {
		inputPayload = buildAliyunShoeTryonInput(req.SourceImage, req.TemplateImage, req.TemplateImageLower)
		parametersPayload = map[string]interface{}{
			"n": 1,
		}
	}

	body := map[string]interface{}{
		"model":      resolvedModelName,
		"input":      inputPayload,
		"parameters": parametersPayload,
	}
	requestPayload = body

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
	responsePayload = resp
	if resp == nil {
		return tryonInvokeResult{}, errors.New("aliyunTryonResponseEmpty")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		result = tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code)}
		return result, nil
	}

	taskID := strings.TrimSpace(resp.Output.TaskID)
	if taskID == "" {
		return tryonInvokeResult{}, errors.New("aliyunTryonTaskIDMissing")
	}

	status := normalizeProviderTaskStatus(resp.Output.TaskStatus)
	if status == tryonTaskStatusFailed {
		result = tryonInvokeResult{Status: tryonTaskStatusFailed, ProviderTaskID: taskID, ErrorMessage: formatAliyunModelError(resp.Message, resp.Code)}
		return result, nil
	}
	if status == tryonTaskStatusSuccess {
		queryTokenCandidates := append([]string{strings.TrimSpace(providerToken)}, modelCfg.authTokenCandidates()...)
		queryResult, queryErr := s.queryAliyunTryonTaskWithTokenFallback(taskID, queryTokenCandidates, createURL, modelCfg.queryTaskURL(), traceCtx, "tryon_aliyun_query", strings.TrimSpace(modelCfg.Key))
		if queryErr != nil {
			return tryonInvokeResult{}, queryErr
		}
		queryResult.ProviderTaskID = taskID
		result = queryResult
		return result, nil
	}

	result = tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}
	return result, nil
}

func (s *TryonTaskService) invokeAliyunRefinerAsync(req clientReq.CreateTryonTaskReq, coarseImageURL string, modelCfg *tryonModelConfig, traceCtx *modelCallTraceContext) (result tryonInvokeResult, err error) {
	startedAt := time.Now()
	var requestPayload interface{}
	var responsePayload interface{}
	var resolvedRefinerModelCfg *tryonModelConfig
	defer func() {
		modelKey := ""
		modelUsage := "refiner"
		modelName := ""
		provider := ""
		endpointURL := ""
		queryURL := ""
		relatedModels := map[string]string{}

		if modelCfg != nil {
			relatedModels["baseModelKey"] = strings.TrimSpace(modelCfg.Key)
			relatedModels["baseModelUsage"] = modelCfg.usageValue()
		}
		if resolvedRefinerModelCfg != nil {
			modelKey = strings.TrimSpace(resolvedRefinerModelCfg.Key)
			modelName = strings.TrimSpace(resolvedRefinerModelCfg.refinerModelValue())
			provider = strings.TrimSpace(resolvedRefinerModelCfg.Provider)
			endpointURL = strings.TrimSpace(resolvedRefinerModelCfg.refinerEndpointURL(strings.TrimSpace(resolvedRefinerModelCfg.endpointURL())))
			queryURL = strings.TrimSpace(resolvedRefinerModelCfg.refinerTaskQueryURL())
			relatedModels["refinerModelKey"] = strings.TrimSpace(resolvedRefinerModelCfg.Key)
		}

		status := strings.TrimSpace(result.Status)
		if status == "" {
			if err != nil {
				status = modelCallStatusError
			} else {
				status = tryonTaskStatusUnknown
			}
		}

		saveModelCallLog(modelCallLogInput{
			TraceContext:    traceCtx,
			CallStage:       "refiner_aliyun_async",
			ModelKey:        modelKey,
			ModelUsage:      modelUsage,
			ModelName:       modelName,
			Provider:        provider,
			EndpointURL:     endpointURL,
			QueryURL:        queryURL,
			RelatedModels:   relatedModels,
			RequestPayload:  requestPayload,
			ResponsePayload: responsePayload,
			ResultImage:     strings.TrimSpace(result.ResultImage),
			Status:          status,
			ErrorMessage:    errorToString(err),
			StartedAt:       startedAt,
			FinishedAt:      time.Now(),
		})
	}()

	if modelCfg == nil {
		return tryonInvokeResult{}, errors.New("tryonModelUnavailable")
	}

	refinerModelCfg, resolveErr := s.resolveRefinerModelConfig(req.SceneType, modelCfg)
	if resolveErr != nil {
		return tryonInvokeResult{}, resolveErr
	}
	if refinerModelCfg == nil {
		return tryonInvokeResult{}, errors.New("tryonRefinerModelUnavailable")
	}
	resolvedRefinerModelCfg = refinerModelCfg

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
	normalizedTemplateLower := ""
	if strings.TrimSpace(req.TemplateImageLower) != "" {
		normalizedTemplateLower, normalizeErr = s.normalizeAliyunMediaURL(req.TemplateImageLower)
		if normalizeErr != nil {
			return tryonInvokeResult{}, normalizeErr
		}
	}
	normalizedCoarse, normalizeErr := s.normalizeAliyunMediaURL(coarseImageURL)
	if normalizeErr != nil {
		return tryonInvokeResult{}, normalizeErr
	}

	providerURL := strings.TrimSpace(refinerModelCfg.endpointURL())
	refinerURL := refinerModelCfg.refinerEndpointURL(providerURL)
	refinerTokenCandidates := refinerModelCfg.refinerAuthTokenCandidates(strings.TrimSpace(refinerModelCfg.authToken()))
	if len(refinerTokenCandidates) == 0 {
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

	input := buildAliyunTryonInput(normalizedSource, normalizedTemplate, normalizedTemplateLower, req.TemplatePart, "")
	input["coarse_image_url"] = strings.TrimSpace(normalizedCoarse)

	body := map[string]interface{}{
		"model": refinerModelCfg.refinerModelValue(),
		"input": input,
		"parameters": map[string]interface{}{
			"gender": refinerModelCfg.refinerGenderValue(),
		},
	}
	requestPayload = body

	for index, refinerToken := range refinerTokenCandidates {
		resp, callErr := external.HttpRequest[dashscopeTryonCreateResponse](external.RequestParams{
			URL:    refinerURL,
			Method: "POST",
			Headers: map[string]string{
				"Authorization":     "Bearer " + strings.TrimSpace(refinerToken),
				"X-DashScope-Async": "enable",
			},
			Body: body,
		})
		if callErr != nil {
			global.GVA_LOG.Warn("阿里精修HTTP请求失败",
				zap.String("refinerURLHost", refinerURLHost),
				zap.Error(callErr),
			)
			if shouldSwitchTryonModelToken(callErr, tryonInvokeResult{}) && index < len(refinerTokenCandidates)-1 {
				global.GVA_LOG.Warn("阿里精修鉴权失败，自动切换备用token重试",
					zap.String("modelKey", strings.TrimSpace(refinerModelCfg.Key)),
					zap.Int("attempt", index+1),
					zap.Int("total", len(refinerTokenCandidates)),
					zap.String("reason", tryonTokenSwitchReason(callErr, tryonInvokeResult{})),
				)
				continue
			}
			return tryonInvokeResult{}, callErr
		}

		responsePayload = resp
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
			failedResult := tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: mappedErr}
			if shouldSwitchTryonModelToken(nil, failedResult) && index < len(refinerTokenCandidates)-1 {
				global.GVA_LOG.Warn("阿里精修鉴权失败，自动切换备用token重试",
					zap.String("modelKey", strings.TrimSpace(refinerModelCfg.Key)),
					zap.Int("attempt", index+1),
					zap.Int("total", len(refinerTokenCandidates)),
					zap.String("reason", tryonTokenSwitchReason(nil, failedResult)),
				)
				continue
			}
			result = failedResult
			return result, nil
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
			failedResult := tryonInvokeResult{Status: tryonTaskStatusFailed, ProviderTaskID: taskID, ErrorMessage: mappedErr}
			if shouldSwitchTryonModelToken(nil, failedResult) && index < len(refinerTokenCandidates)-1 {
				global.GVA_LOG.Warn("阿里精修鉴权失败，自动切换备用token重试",
					zap.String("modelKey", strings.TrimSpace(refinerModelCfg.Key)),
					zap.Int("attempt", index+1),
					zap.Int("total", len(refinerTokenCandidates)),
					zap.String("reason", tryonTokenSwitchReason(nil, failedResult)),
				)
				continue
			}
			result = failedResult
			return result, nil
		}
		if status == tryonTaskStatusSuccess {
			queryTokenCandidates := append([]string{strings.TrimSpace(refinerToken)}, refinerTokenCandidates...)
			queryResult, queryErr := s.queryAliyunTryonTaskWithTokenFallback(taskID, queryTokenCandidates, refinerURL, refinerModelCfg.refinerTaskQueryURL(), traceCtx, "refiner_aliyun_query", strings.TrimSpace(refinerModelCfg.Key))
			if queryErr != nil {
				return tryonInvokeResult{}, queryErr
			}
			queryResult.ProviderTaskID = taskID
			result = queryResult
			return result, nil
		}

		result = tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}
		return result, nil
	}

	return tryonInvokeResult{}, errors.New("aliyunRefinerApiKeyMissing")
}

func (s *TryonTaskService) invokeAliyunParsing(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, providerToken string, providerURL string, modelName string, preferParsingImage bool, traceCtx *modelCallTraceContext) (result tryonInvokeResult, err error) {
	startedAt := time.Now()
	var requestPayload interface{}
	var responsePayload interface{}
	defer func() {
		status := strings.TrimSpace(result.Status)
		if status == "" {
			if err != nil {
				status = modelCallStatusError
			} else {
				status = tryonTaskStatusUnknown
			}
		}

		modelKey := ""
		modelUsage := "parsing"
		provider := ""
		if modelCfg != nil {
			modelKey = strings.TrimSpace(modelCfg.Key)
			provider = strings.TrimSpace(modelCfg.Provider)
		}

		queryURL := ""
		if modelCfg != nil {
			queryURL = strings.TrimSpace(modelCfg.queryTaskURL())
		}

		saveModelCallLog(modelCallLogInput{
			TraceContext:    traceCtx,
			CallStage:       "parsing_aliyun_sync",
			ModelKey:        modelKey,
			ModelUsage:      modelUsage,
			ModelName:       strings.TrimSpace(modelName),
			Provider:        provider,
			EndpointURL:     strings.TrimSpace(providerURL),
			QueryURL:        queryURL,
			RequestPayload:  requestPayload,
			ResponsePayload: responsePayload,
			ResultImage:     strings.TrimSpace(result.ResultImage),
			Status:          status,
			ErrorMessage:    errorToString(err),
			StartedAt:       startedAt,
			FinishedAt:      time.Now(),
		})
	}()

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
	requestPayload = body

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
			return s.invokeAliyunParsingAsync(req, modelCfg, providerToken, processURL, resolvedModelName, clothesType, preferParsingImage, traceCtx)
		}
		return tryonInvokeResult{}, err
	}
	responsePayload = resp
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
			return s.invokeAliyunParsingAsync(req, modelCfg, providerToken, processURL, resolvedModelName, clothesType, preferParsingImage, traceCtx)
		}
		result = tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: mappedErr}
		return result, nil
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
		result = tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: "aliyunParsingNoUsableGarmentArea"}
		return result, nil
	}

	result = tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage}
	return result, nil
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

func (s *TryonTaskService) invokeAliyunParsingAsync(req clientReq.CreateTryonTaskReq, modelCfg *tryonModelConfig, providerToken string, processURL string, modelName string, clothesType []string, preferParsingImage bool, traceCtx *modelCallTraceContext) (result tryonInvokeResult, err error) {
	startedAt := time.Now()
	var requestPayload interface{}
	var responsePayload interface{}
	defer func() {
		status := strings.TrimSpace(result.Status)
		if status == "" {
			if err != nil {
				status = modelCallStatusError
			} else {
				status = tryonTaskStatusUnknown
			}
		}

		modelKey := ""
		modelUsage := "parsing"
		provider := ""
		queryURL := ""
		if modelCfg != nil {
			modelKey = strings.TrimSpace(modelCfg.Key)
			provider = strings.TrimSpace(modelCfg.Provider)
			queryURL = strings.TrimSpace(modelCfg.queryTaskURL())
		}

		saveModelCallLog(modelCallLogInput{
			TraceContext:    traceCtx,
			CallStage:       "parsing_aliyun_async",
			ModelKey:        modelKey,
			ModelUsage:      modelUsage,
			ModelName:       strings.TrimSpace(modelName),
			Provider:        provider,
			EndpointURL:     strings.TrimSpace(processURL),
			QueryURL:        queryURL,
			RequestPayload:  requestPayload,
			ResponsePayload: responsePayload,
			ResultImage:     strings.TrimSpace(result.ResultImage),
			Status:          status,
			ErrorMessage:    errorToString(err),
			StartedAt:       startedAt,
			FinishedAt:      time.Now(),
		})
	}()

	body := map[string]interface{}{
		"model": modelName,
		"input": map[string]interface{}{
			"image_url": strings.TrimSpace(req.SourceImage),
		},
		"parameters": map[string]interface{}{
			"clothes_type": clothesType,
		},
	}
	requestPayload = body

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
	responsePayload = resp
	if resp == nil {
		return tryonInvokeResult{}, errors.New("aliyunParsingResponseEmpty")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		result = tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code)}
		return result, nil
	}

	taskID := strings.TrimSpace(resp.Output.TaskID)
	if taskID == "" {
		return tryonInvokeResult{}, errors.New("aliyunTryonTaskIDMissing")
	}

	status := normalizeProviderTaskStatus(resp.Output.TaskStatus)
	if status == tryonTaskStatusFailed {
		result = tryonInvokeResult{Status: tryonTaskStatusFailed, ProviderTaskID: taskID, ErrorMessage: formatAliyunModelError(resp.Message, resp.Code)}
		return result, nil
	}

	taskQueryURL := ""
	if modelCfg != nil {
		taskQueryURL = modelCfg.queryTaskURL()
	}

	const maxAttempts = 20
	const interval = 1500 * time.Millisecond

	for attempt := 0; attempt < maxAttempts; attempt++ {
		queryResult, queryErr := s.queryAliyunParsingTask(taskID, providerToken, processURL, taskQueryURL, preferParsingImage, traceCtx)
		if queryErr != nil {
			return tryonInvokeResult{}, queryErr
		}
		if queryResult.Status == tryonTaskStatusProcessing || queryResult.Status == tryonTaskStatusUnknown {
			if attempt == maxAttempts-1 {
				break
			}
			time.Sleep(interval)
			continue
		}
		result = queryResult
		return result, nil
	}

	result = tryonInvokeResult{Status: tryonTaskStatusFailed, ProviderTaskID: taskID, ErrorMessage: "tryonTaskFailed"}
	return result, nil
}

func (s *TryonTaskService) queryAliyunParsingTask(taskID string, providerToken string, createURL string, taskQueryURL string, preferParsingImage bool, traceCtx *modelCallTraceContext) (result tryonInvokeResult, err error) {
	startedAt := time.Now()
	var responsePayload interface{}
	defer func() {
		status := strings.TrimSpace(result.Status)
		if status == "" {
			if err != nil {
				status = modelCallStatusError
			} else {
				status = tryonTaskStatusUnknown
			}
		}
		saveModelCallLog(modelCallLogInput{
			TraceContext:    traceCtx,
			CallStage:       "parsing_aliyun_query",
			ModelUsage:      "parsing",
			EndpointURL:     strings.TrimSpace(createURL),
			QueryURL:        strings.TrimSpace(taskQueryURL),
			ResponsePayload: responsePayload,
			ResultImage:     strings.TrimSpace(result.ResultImage),
			Status:          status,
			ErrorMessage:    errorToString(err),
			StartedAt:       startedAt,
			FinishedAt:      time.Now(),
		})
	}()

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
	responsePayload = resp
	if resp == nil {
		return tryonInvokeResult{}, errors.New("aliyunTaskQueryResponseEmpty")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		result = tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code), ProviderTaskID: taskID}
		return result, nil
	}

	status := normalizeProviderTaskStatus(resp.Output.TaskStatus)
	switch status {
	case tryonTaskStatusSuccess:
		resultImage := pickAliyunParsingResultImage(preferParsingImage, resp.Output.ParsingImgURL, resp.Output.CropImgURL, resp.Output.ImageURL)
		if strings.TrimSpace(resultImage) == "" {
			result = tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: "aliyunTaskSuccessNoImage", ProviderTaskID: taskID}
			return result, nil
		}
		result = tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: strings.TrimSpace(resultImage), ProviderTaskID: taskID}
		return result, nil
	case tryonTaskStatusFailed:
		result = tryonInvokeResult{
			Status:         tryonTaskStatusFailed,
			ProviderTaskID: taskID,
			ErrorMessage: formatAliyunModelError(
				resp.Output.Message,
				resp.Output.Code,
				resp.Message,
				resp.Code,
				"tryonTaskFailed",
			),
		}
		return result, nil
	case tryonTaskStatusProcessing, tryonTaskStatusUnknown:
		result = tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}
		return result, nil
	default:
		result = tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}
		return result, nil
	}
}

func (s *TryonTaskService) invokeTryonBeautifyProvider(sceneType string, resultImage string, modelCfg *tryonModelConfig, req clientReq.ApplyTryonBeautifyReq, traceCtx *modelCallTraceContext) (outputImage string, err error) {
	startedAt := time.Now()
	var requestPayload interface{}
	var responsePayload interface{}
	defer func() {
		modelKey := ""
		modelUsage := "beautify"
		modelName := ""
		provider := ""
		endpointURL := ""
		if modelCfg != nil {
			modelKey = strings.TrimSpace(modelCfg.Key)
			modelName = strings.TrimSpace(modelCfg.beautifyModelValue())
			provider = strings.TrimSpace(modelCfg.Provider)
			endpointURL = strings.TrimSpace(firstNonEmptyString(modelCfg.endpointURL(), modelCfg.beautifyEndpointURL()))
		}

		status := modelCallStatusSuccess
		if err != nil {
			status = modelCallStatusError
		}

		saveModelCallLog(modelCallLogInput{
			TraceContext:    traceCtx,
			CallStage:       "beautify_provider",
			ModelKey:        modelKey,
			ModelUsage:      modelUsage,
			ModelName:       modelName,
			Provider:        provider,
			EndpointURL:     endpointURL,
			RequestPayload:  requestPayload,
			ResponsePayload: responsePayload,
			ResultImage:     strings.TrimSpace(outputImage),
			Status:          status,
			ErrorMessage:    errorToString(err),
			StartedAt:       startedAt,
			FinishedAt:      time.Now(),
		})
	}()

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
		outputImage = strings.TrimSpace(resultImage)
		return outputImage, nil
	}

	modelName := modelCfg.beautifyModelValue()
	providerName := strings.TrimSpace(modelCfg.Provider)
	modelKey := strings.TrimSpace(modelCfg.Key)
	accessKeyID, accessKeySecret, _ := modelCfg.beautifyAccessCredentials(providerToken)
	hasAliyunCredentials := strings.TrimSpace(accessKeyID) != "" && strings.TrimSpace(accessKeySecret) != ""

	if isAliyunBeautifyModel(providerName, modelName, providerURL) || hasAliyunCredentials {
		return s.invokeAliyunRetouchSkin(resultImage, modelCfg, providerToken, req, traceCtx)
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
		"operationType": "beautify",
		"sceneType":     strings.TrimSpace(sceneType),
		"imageUrl":      strings.TrimSpace(normalizedImage),
		"sourceImage":   strings.TrimSpace(normalizedImage),
		"resultImage":   strings.TrimSpace(normalizedImage),
	}
	if retouchDegree > 0 {
		body["retouchDegree"] = retouchDegree
	}
	if whiteningDegree > 0 {
		body["whiteningDegree"] = whiteningDegree
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
	requestPayload = body

	resp, err := external.HttpRequest[tryonProviderResponse](external.RequestParams{
		URL:     strings.TrimSpace(targetURL),
		Method:  "POST",
		Headers: headers,
		Body:    body,
	})
	if err != nil {
		return "", err
	}
	responsePayload = resp
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

	outputImage = strings.TrimSpace(resultURL)
	return outputImage, nil
}

func (s *TryonTaskService) invokeAliyunRetouchSkin(resultImage string, modelCfg *tryonModelConfig, fallbackToken string, req clientReq.ApplyTryonBeautifyReq, traceCtx *modelCallTraceContext) (outputImage string, err error) {
	startedAt := time.Now()
	var requestPayload interface{}
	var responsePayload interface{}
	defer func() {
		modelKey := ""
		modelUsage := "beautify"
		modelName := ""
		provider := ""
		endpointURL := ""
		if modelCfg != nil {
			modelKey = strings.TrimSpace(modelCfg.Key)
			modelName = strings.TrimSpace(modelCfg.beautifyModelValue())
			provider = strings.TrimSpace(modelCfg.Provider)
			endpointURL = strings.TrimSpace(modelCfg.beautifyEndpointURL())
		}

		status := modelCallStatusSuccess
		if err != nil {
			status = modelCallStatusError
		}

		saveModelCallLog(modelCallLogInput{
			TraceContext:    traceCtx,
			CallStage:       "beautify_aliyun_rpc",
			ModelKey:        modelKey,
			ModelUsage:      modelUsage,
			ModelName:       modelName,
			Provider:        provider,
			EndpointURL:     endpointURL,
			RequestPayload:  requestPayload,
			ResponsePayload: responsePayload,
			ResultImage:     strings.TrimSpace(outputImage),
			Status:          status,
			ErrorMessage:    errorToString(err),
			StartedAt:       startedAt,
			FinishedAt:      time.Now(),
		})
	}()

	normalizedImage, normalizeErr := s.normalizeAliyunRetouchImageURL(resultImage)
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
	if action == "" || !strings.EqualFold(action, aliyunRetouchSkinAction) {
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
	}
	if degree, ok := normalizeAliyunBeautifyDegree(retouchDegree); ok {
		params["RetouchDegree"] = formatAliyunBeautifyDegree(degree)
	}
	if degree, ok := normalizeAliyunBeautifyDegree(whiteningDegree); ok {
		params["WhiteningDegree"] = formatAliyunBeautifyDegree(degree)
	}
	if strings.TrimSpace(securityToken) != "" {
		params["SecurityToken"] = strings.TrimSpace(securityToken)
	}

	signature := signAliyunRPCParams(params, strings.TrimSpace(accessKeySecret))
	params["Signature"] = signature
	requestPayload = params

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
	responsePayload = map[string]interface{}{"statusCode": resp.StatusCode, "body": string(responseBody)}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", errors.New(formatAliyunModelError(string(responseBody), "tryonBeautifyRequestFailed"))
	}

	var apiResp aliyunRetouchSkinResponse
	if err := json.Unmarshal(responseBody, &apiResp); err != nil {
		return "", errors.New("tryonBeautifyResponseInvalid")
	}
	responsePayload = apiResp

	if strings.TrimSpace(apiResp.Code) != "" {
		return "", errors.New(formatAliyunModelError(apiResp.Message, apiResp.Code))
	}

	resultURL := strings.TrimSpace(apiResp.Data.ImageURL)
	if resultURL == "" {
		return "", errors.New("tryonBeautifyResultEmpty")
	}
	outputImage = resultURL
	return outputImage, nil
}

func (s *TryonTaskService) queryAliyunTryonTask(taskID string, providerToken string, createURL string, taskQueryURL string, traceCtx *modelCallTraceContext, callStage string) (result tryonInvokeResult, err error) {
	startedAt := time.Now()
	var responsePayload interface{}
	defer func() {
		stage := strings.TrimSpace(callStage)
		if stage == "" {
			stage = "tryon_aliyun_query"
		}
		status := strings.TrimSpace(result.Status)
		if status == "" {
			if err != nil {
				status = modelCallStatusError
			} else {
				status = tryonTaskStatusUnknown
			}
		}
		saveModelCallLog(modelCallLogInput{
			TraceContext:    traceCtx,
			CallStage:       stage,
			ModelUsage:      traceModelUsage(traceCtx),
			EndpointURL:     strings.TrimSpace(createURL),
			QueryURL:        strings.TrimSpace(taskQueryURL),
			ResponsePayload: responsePayload,
			ResultImage:     strings.TrimSpace(result.ResultImage),
			Status:          status,
			ErrorMessage:    errorToString(err),
			StartedAt:       startedAt,
			FinishedAt:      time.Now(),
		})
	}()

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
	responsePayload = resp
	if resp == nil {
		return tryonInvokeResult{}, errors.New("aliyunTaskQueryResponseEmpty")
	}

	if code := strings.TrimSpace(resp.Code); code != "" {
		result = tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: formatAliyunModelError(resp.Message, code), ProviderTaskID: taskID}
		return result, nil
	}

	status := normalizeProviderTaskStatus(resp.Output.TaskStatus)
	switch status {
	case tryonTaskStatusSuccess:
		resultImage := extractDashscopeImageURL(resp.Output.ImageURL)
		if resultImage == "" {
			for _, item := range resp.Output.ImageURLs {
				if strings.TrimSpace(item) != "" {
					resultImage = strings.TrimSpace(item)
					break
				}
			}
		}
		if resultImage == "" {
			resultImage = strings.TrimSpace(resp.Output.ResultURL)
		}
		if resultImage == "" {
			for _, item := range resp.Output.ResultURLs {
				if strings.TrimSpace(item) != "" {
					resultImage = strings.TrimSpace(item)
					break
				}
			}
		}
		if resultImage == "" {
			for _, item := range resp.Output.Results {
				resultImage = firstNonEmptyString(
					strings.TrimSpace(item.URL),
					extractDashscopeImageURL(item.ImageURL),
					strings.TrimSpace(item.ResultURL),
				)
				if resultImage != "" {
					break
				}
			}
		}
		if resultImage == "" {
			result = tryonInvokeResult{Status: tryonTaskStatusFailed, ErrorMessage: "aliyunTaskSuccessNoImage", ProviderTaskID: taskID}
			return result, nil
		}
		result = tryonInvokeResult{Status: tryonTaskStatusSuccess, ResultImage: resultImage, ProviderTaskID: taskID}
		return result, nil
	case tryonTaskStatusFailed:
		result = tryonInvokeResult{
			Status:         tryonTaskStatusFailed,
			ProviderTaskID: taskID,
			ErrorMessage: formatAliyunModelError(
				resp.Output.Message,
				resp.Output.Code,
				resp.Message,
				resp.Code,
				"tryonTaskFailed",
			),
		}
		return result, nil
	case tryonTaskStatusProcessing, tryonTaskStatusUnknown:
		result = tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}
		return result, nil
	default:
		result = tryonInvokeResult{Status: tryonTaskStatusProcessing, ProviderTaskID: taskID}
		return result, nil
	}
}

func formatAliyunModelError(values ...string) string {
	msg := firstNonEmptyString(values...)
	if msg == "" {
		return "tryonTaskFailed"
	}

	lowerMsg := strings.ToLower(msg)
	if (strings.Contains(lowerMsg, "invalidimage.url") && (strings.Contains(lowerMsg, "oss") || strings.Contains(lowerMsg, "cdn") || strings.Contains(lowerMsg, "custom"))) ||
		strings.Contains(lowerMsg, "标准的oss域名") ||
		strings.Contains(lowerMsg, "暂不支持绑定cdn域名") ||
		strings.Contains(lowerMsg, "自定义域名") {
		return "modelAliyunRetouchNeedStandardOSS"
	}
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

func containsAnySubstring(text string, values ...string) bool {
	for _, value := range values {
		if value == "" {
			continue
		}
		if strings.Contains(text, value) {
			return true
		}
	}
	return false
}

func isTryonTokenRelatedErrorText(raw string) bool {
	text := strings.ToLower(strings.TrimSpace(raw))
	if text == "" {
		return false
	}

	if containsAnySubstring(text,
		"aliyuntryonapikeymissing",
		"aliyunparsingapikeymissing",
		"aliyunrefinerapikeymissing",
		"tokeninvalid",
		"tokenexpired",
		"tokenfail",
		"invalidapikey",
		"invalid access key",
		"invalidaccesskey",
		"unauthorized",
		"forbidden",
		"access denied",
		"accessdenied",
		"security token",
		"securitytoken",
		"signaturedoesnotmatch",
		"signature not match",
		"signature mismatch",
	) {
		return true
	}

	if containsAnySubstring(text, "api key", "apikey") {
		return true
	}

	if containsAnySubstring(text, "token") && containsAnySubstring(text,
		"invalid",
		"expired",
		"missing",
		"unauthor",
		"forbidden",
		"denied",
		"fail",
		"error",
		"失效",
		"过期",
		"无效",
		"缺失",
	) {
		return true
	}

	if containsAnySubstring(text,
		"quota",
		"out of quota",
		"quota exceeded",
		"insufficient balance",
		"余额不足",
		"额度不足",
		"欠费",
	) {
		return true
	}

	return false
}

func shouldSwitchTryonModelToken(err error, result tryonInvokeResult) bool {
	if err != nil && isTryonTokenRelatedErrorText(err.Error()) {
		return true
	}
	if isTryonTokenRelatedErrorText(result.ErrorMessage) {
		return true
	}
	return false
}

func tryonTokenSwitchReason(err error, result tryonInvokeResult) string {
	if err != nil && strings.TrimSpace(err.Error()) != "" {
		return strings.TrimSpace(err.Error())
	}
	if strings.TrimSpace(result.ErrorMessage) != "" {
		return strings.TrimSpace(result.ErrorMessage)
	}
	return "tokenInvalid"
}

func (s *TryonTaskService) queryAliyunTryonTaskWithTokenFallback(taskID string, tokenCandidates []string, createURL string, taskQueryURL string, traceCtx *modelCallTraceContext, callStage string, modelKey string) (tryonInvokeResult, error) {
	normalizedCandidates := make([]string, 0, len(tokenCandidates))
	for _, token := range tokenCandidates {
		normalizedCandidates = appendUniqueTokenCandidate(normalizedCandidates, token)
	}

	if len(normalizedCandidates) == 0 {
		return s.queryAliyunTryonTask(taskID, "", createURL, taskQueryURL, traceCtx, callStage)
	}

	for index, token := range normalizedCandidates {
		result, err := s.queryAliyunTryonTask(taskID, token, createURL, taskQueryURL, traceCtx, callStage)
		if shouldSwitchTryonModelToken(err, result) && index < len(normalizedCandidates)-1 {
			global.GVA_LOG.Warn("查询阿里任务鉴权失败，自动切换备用token重试",
				zap.String("modelKey", strings.TrimSpace(modelKey)),
				zap.String("providerTaskID", strings.TrimSpace(taskID)),
				zap.String("callStage", strings.TrimSpace(callStage)),
				zap.Int("attempt", index+1),
				zap.Int("total", len(normalizedCandidates)),
				zap.String("reason", tryonTokenSwitchReason(err, result)),
			)
			continue
		}
		return result, err
	}

	return s.queryAliyunTryonTask(taskID, normalizedCandidates[len(normalizedCandidates)-1], createURL, taskQueryURL, traceCtx, callStage)
}

func (s *TryonTaskService) normalizeAliyunRetouchImageURL(rawURL string) (string, error) {
	normalizedURL, err := s.normalizeAliyunMediaURL(rawURL)
	if err != nil {
		return "", err
	}

	parsed, err := url.Parse(normalizedURL)
	if err != nil || strings.TrimSpace(parsed.Hostname()) == "" {
		return normalizedURL, nil
	}

	hostname := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if isAliyunStandardOSSHost(hostname) {
		return normalizedURL, nil
	}

	standardBaseURL := getAliyunOSSStandardBaseURL()
	if standardBaseURL == "" {
		return normalizedURL, nil
	}

	if !s.shouldRewriteAliyunRetouchHost(hostname) {
		return normalizedURL, nil
	}

	resourcePath := strings.TrimSpace(parsed.Path)
	if resourcePath == "" {
		resourcePath = "/"
	}
	return joinBaseURLAndResourcePath(standardBaseURL, resourcePath)
}

func (s *TryonTaskService) shouldRewriteAliyunRetouchHost(currentHost string) bool {
	currentHost = strings.ToLower(strings.TrimSpace(currentHost))
	if currentHost == "" || isAliyunStandardOSSHost(currentHost) {
		return false
	}

	rewriteHosts := map[string]struct{}{}
	for _, raw := range []string{
		global.GVA_CONFIG.AliyunOSS.BucketUrl,
		s.getTryonMediaPublicBaseURL(),
	} {
		host := urlHostnameOrHost(raw)
		if host == "" {
			continue
		}
		rewriteHosts[host] = struct{}{}
	}

	if strings.EqualFold(strings.TrimSpace(global.GVA_CONFIG.System.OssType), "aliyun-oss") {
		hotlinkHost := urlHostnameOrHost(global.GVA_CONFIG.Hotlink.CdnDomain)
		if hotlinkHost != "" {
			rewriteHosts[hotlinkHost] = struct{}{}
		}
	}

	_, ok := rewriteHosts[currentHost]
	return ok
}

func getAliyunOSSStandardBaseURL() string {
	bucketName := strings.TrimSpace(global.GVA_CONFIG.AliyunOSS.BucketName)
	if bucketName == "" {
		return ""
	}

	endpointHost := urlHostnameOrHost(global.GVA_CONFIG.AliyunOSS.Endpoint)
	if endpointHost == "" {
		return ""
	}

	bucketPrefix := strings.ToLower(strings.TrimSpace(bucketName)) + "."
	endpointHost = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(endpointHost)), bucketPrefix)
	endpointHost = strings.Replace(endpointHost, "-internal", "", 1)
	if !strings.Contains(endpointHost, "aliyuncs.com") {
		return ""
	}

	return "https://" + strings.TrimSpace(bucketName) + "." + endpointHost
}

func isAliyunStandardOSSHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" {
		return false
	}
	if !strings.HasSuffix(host, ".aliyuncs.com") {
		return false
	}
	return strings.Contains(host, ".oss-") || strings.Contains(host, "-oss-")
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
	publicBaseURL := ""
	if global.GVA_DB != nil {
		sysConfigService := SysConfigService{}
		publicBaseURL, _ = sysConfigService.GetConfigByKey("tryon_media_public_base_url")
	}
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

func urlHostnameOrHost(rawValue string) string {
	rawValue = strings.TrimSpace(rawValue)
	if rawValue == "" {
		return ""
	}

	host := urlHostname(rawValue)
	if host != "" {
		return host
	}

	normalized := strings.TrimPrefix(rawValue, "//")
	normalized = strings.TrimPrefix(normalized, "/")
	if normalized == "" {
		return ""
	}
	if parsed, err := url.Parse("https://" + normalized); err == nil {
		return strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	}
	return ""
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
	if modelCfg == nil {
		return nil
	}
	traceCtx := createTraceContextFromTask(task)

	providerURL := strings.TrimSpace(modelCfg.endpointURL())
	providerTokenCandidates := modelCfg.authTokenCandidates()
	providerName := ""
	modelName := ""
	taskQueryURL := ""
	refinerModelCfg := modelCfg
	if resolvedRefinerModelCfg, resolveErr := s.resolveRefinerModelConfig(task.SceneType, modelCfg); resolveErr == nil && resolvedRefinerModelCfg != nil {
		refinerModelCfg = resolvedRefinerModelCfg
	} else if resolveErr != nil {
		global.GVA_LOG.Warn("解析精修模型失败，回退当前试衣模型配置",
			zap.String("taskNo", strings.TrimSpace(task.TaskNo)),
			zap.String("provider", strings.TrimSpace(task.Provider)),
			zap.Error(resolveErr),
		)
	}
	refinerURL := strings.TrimSpace(providerURL)
	refinerTokenCandidates := refinerModelCfg.refinerAuthTokenCandidates(strings.TrimSpace(refinerModelCfg.authToken()))
	refinerTaskQueryURL := ""
	providerName = modelCfg.Provider
	modelName = resolveAliyunModelName(modelCfg, task.SceneType)
	taskQueryURL = modelCfg.queryTaskURL()
	refinerURL = refinerModelCfg.refinerEndpointURL(strings.TrimSpace(refinerModelCfg.endpointURL()))
	refinerTaskQueryURL = refinerModelCfg.refinerTaskQueryURL()
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

		result, err := s.queryAliyunTryonTaskWithTokenFallback(refinerTaskID, refinerTokenCandidates, refinerURL, refinerTaskQueryURL, traceCtx, "refiner_aliyun_query", strings.TrimSpace(refinerModelCfg.Key))
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
			if err = s.markTryonRefinerFailedAndKeepSuccess(ctx, task, strings.TrimSpace(task.ResultImage), result.ErrorMessage, refinerTaskID, nil); err != nil {
				return err
			}
			return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
		default:
			return nil
		}
	}

	result, err := s.queryAliyunTryonTaskWithTokenFallback(task.ProviderTaskID, providerTokenCandidates, providerURL, taskQueryURL, traceCtx, "tryon_aliyun_query", strings.TrimSpace(modelCfg.Key))
	if err != nil {
		return err
	}

	switch result.Status {
	case tryonTaskStatusProcessing:
		return nil
	case tryonTaskStatusSuccess:
		if task.EnableRefiner && strings.EqualFold(strings.TrimSpace(task.RefinerStatus), tryonRefinerStatusPending) {
			storedMainResultImage := s.persistTryonResultImage(strings.TrimSpace(result.ResultImage), task.TaskNo)
			refinerModelCfg := modelCfg
			if modelCfg != nil && strings.TrimSpace(task.RefinerModelKey) != "" {
				copied := *modelCfg
				copied.RefinerModelKey = strings.TrimSpace(task.RefinerModelKey)
				refinerModelCfg = &copied
			}
			refinerReq := clientReq.CreateTryonTaskReq{
				SceneType:     task.SceneType,
				SourceImage:   task.SourceImage,
				TemplateImage: task.TemplateImage,
				ModelKey:      task.Provider,
			}
			refinerTraceCtx := cloneTraceContext(traceCtx)
			if refinerTraceCtx != nil {
				refinerTraceCtx.Behavior = "refiner_generate"
			}
			refinerResult, refinerErr := s.invokeAliyunRefinerAsync(refinerReq, strings.TrimSpace(result.ResultImage), refinerModelCfg, refinerTraceCtx)
			if refinerErr != nil {
				if err = s.markTryonRefinerFailedAndKeepSuccess(ctx, task, storedMainResultImage, refinerErr.Error(), "", nil); err != nil {
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
						"result_image":     storedMainResultImage,
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
				task.ResultImage = storedMainResultImage
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
				if err = s.markTryonRefinerFailedAndKeepSuccess(ctx, task, storedMainResultImage, refinerResult.ErrorMessage, strings.TrimSpace(refinerResult.ProviderTaskID), nil); err != nil {
					return err
				}
				return global.GVA_DB.Where("id = ? AND user_id = ?", task.ID, task.UserID).First(task).Error
			default:
				if err = s.markTryonRefinerFailedAndKeepSuccess(ctx, task, storedMainResultImage, "tryonRefinerStatusUnknown", strings.TrimSpace(refinerResult.ProviderTaskID), nil); err != nil {
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
	if strings.EqualFold(strings.TrimSpace(sceneType), "shoes") {
		return aliyunShoeTryonModel
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
	urlLower := strings.ToLower(strings.TrimSpace(providerURL))
	if modelLower == aliyunShoeTryonModel {
		if strings.Contains(providerLower, "aliyun") || strings.Contains(providerLower, "dashscope") {
			return true
		}
		if strings.Contains(urlLower, "dashscope.aliyuncs.com") || strings.Contains(urlLower, "virtualmodel/generation") {
			return true
		}
	}

	if strings.Contains(providerLower, "aliyun") || strings.Contains(providerLower, "dashscope") {
		if !strings.Contains(providerLower, "parsing") {
			return true
		}
	}

	if strings.Contains(urlLower, "dashscope.aliyuncs.com") && strings.Contains(urlLower, "image-synthesis") {
		return true
	}
	if strings.Contains(urlLower, "dashscope.aliyuncs.com") && strings.Contains(urlLower, "virtualmodel/generation") {
		return true
	}
	if strings.EqualFold(strings.TrimSpace(sceneType), "shoes") && strings.EqualFold(strings.TrimSpace(providerURL), aliyunShoeTryonURL) {
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

func canEnableAliyunParsing(sceneType string, modelCfg *tryonModelConfig, templatePart string) bool {
	if modelCfg == nil || !modelCfg.isTryonModelUsage() {
		return false
	}
	modelName := resolveAliyunModelName(modelCfg, sceneType)
	return shouldEnableAliyunParsingForSinglePart(sceneType, modelName, templatePart, true)
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

func normalizeAliyunBeautifyDegree(value float64) (float64, bool) {
	if value != value {
		return 0, false
	}
	if value <= 0 {
		return 0, false
	}
	if value > 1.5 {
		return 0, false
	}
	return value, true
}

func formatAliyunBeautifyDegree(value float64) string {
	return strconv.FormatFloat(value, 'f', 3, 64)
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
