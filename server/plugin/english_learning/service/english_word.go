package service

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EnglishWordService struct{}

type WordErrorLogItem struct {
	model.EnglishWordErrorLog
	Word *model.EnglishWord `json:"word,omitempty"`
}

type TTSPreflightResult struct {
	Word       string `json:"word"`
	VoiceUS    string `json:"voiceUs,omitempty"`
	VoiceUK    string `json:"voiceUk,omitempty"`
	AudioUS    string `json:"audioUs,omitempty"`
	AudioUK    string `json:"audioUk,omitempty"`
	DurationMS int64  `json:"durationMs"`
}

type UpsertWordFromSQLResult struct {
	Action           string `json:"action"`
	WordID           uint   `json:"wordId"`
	AudioUSGenerated bool   `json:"audioUsGenerated"`
	AudioUKGenerated bool   `json:"audioUkGenerated"`
}

type BatchFillWordFromDictionaryItem struct {
	Word       string   `json:"word"`
	WordID     uint     `json:"wordId"`
	Status     string   `json:"status"`
	Message    string   `json:"message"`
	Applied    []string `json:"applied,omitempty"`
	Translated []string `json:"translated,omitempty"`
}

type BatchFillWordFromDictionaryResult struct {
	CategoryID uint                              `json:"categoryId"`
	Total      int                               `json:"total"`
	Success    int                               `json:"success"`
	Skipped    int                               `json:"skipped"`
	Failed     int                               `json:"failed"`
	Items      []BatchFillWordFromDictionaryItem `json:"items"`
}

type dictionaryEntry struct {
	Phonetics []struct {
		Text  string `json:"text"`
		Audio string `json:"audio"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Example    string `json:"example"`
		} `json:"definitions"`
	} `json:"meanings"`
}

const (
	defaultLearningTTSTimeoutMS int64 = 8000
	defaultLearningTTSVoiceUS         = "en-US-JennyNeural"
	defaultLearningTTSVoiceUK         = "en-GB-SoniaNeural"
)

var (
	errLearningTTSDisabled                 = errors.New("TTS自动生成未开启")
	errLearningTTSProviderURLNotConfigured = errors.New("TTS服务地址未配置")
	errLearningTTSMissingAudioURL          = errors.New("TTS服务未返回音频地址")
)

func getLearningStringSetting(sysConfigKey, envKey, defaultValue string) string {
	if value, ok := utils.GetSysConfigRawValue(sysConfigKey); ok {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}

	if strings.TrimSpace(envKey) != "" {
		envValue := strings.TrimSpace(os.Getenv(envKey))
		if envValue != "" {
			return envValue
		}
	}

	return defaultValue
}

func logEnglishLearningWarn(message string, fields ...zap.Field) {
	if global.GVA_LOG == nil {
		return
	}
	global.GVA_LOG.Warn(message, fields...)
}

func normalizeWordKey(word string) string {
	return strings.ToLower(strings.TrimSpace(word))
}

func detectAudioExtFromURL(audioURL string) string {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(audioURL)))
	switch ext {
	case ".mp3", ".wav", ".ogg", ".m4a", ".aac":
		return ext
	default:
		return ""
	}
}

func detectAudioExtFromContentType(contentType string) string {
	mediaType := strings.ToLower(strings.TrimSpace(contentType))
	if index := strings.Index(mediaType, ";"); index >= 0 {
		mediaType = strings.TrimSpace(mediaType[:index])
	}

	switch mediaType {
	case "audio/mpeg", "audio/mp3":
		return ".mp3"
	case "audio/wav", "audio/x-wav":
		return ".wav"
	case "audio/ogg":
		return ".ogg"
	case "audio/aac":
		return ".aac"
	case "audio/mp4", "audio/x-m4a":
		return ".m4a"
	default:
		return ""
	}
}

func buildAudioUploadFileName(text, voice, sourceURL, ext string) string {
	if ext == "" {
		ext = ".mp3"
	}
	source := normalizeWordKey(text) + "|" + strings.TrimSpace(voice) + "|" + strings.TrimSpace(sourceURL)
	hash := sha1.Sum([]byte(source))
	return fmt.Sprintf("word-%x-%d%s", hash[:8], time.Now().Unix(), ext)
}

func buildWordAudioFolder(word string) string {
	normalized := normalizeWordKey(word)
	if normalized == "" {
		return "english-learn/audio/word"
	}

	var builder strings.Builder
	builder.Grow(len(normalized))
	lastDash := false
	for _, r := range normalized {
		isDigit := r >= '0' && r <= '9'
		isLowerAlpha := r >= 'a' && r <= 'z'
		if isDigit || isLowerAlpha {
			builder.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			builder.WriteByte('-')
			lastDash = true
		}
	}

	folderKey := strings.Trim(builder.String(), "-")
	if folderKey == "" {
		folderKey = "word"
	}
	if len(folderKey) > 64 {
		folderKey = strings.Trim(folderKey[:64], "-")
		if folderKey == "" {
			folderKey = "word"
		}
	}

	return "english-learn/audio/" + folderKey
}

func buildMultipartFileHeader(fileName, contentType string, data []byte) (*multipart.FileHeader, error) {
	if len(data) == 0 {
		return nil, errors.New("音频内容为空")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, strings.ReplaceAll(fileName, `"`, "")))
	if strings.TrimSpace(contentType) != "" {
		partHeader.Set("Content-Type", contentType)
	}

	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return nil, err
	}
	if _, err = part.Write(data); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/internal-upload", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err = req.ParseMultipartForm(int64(len(data) + 1024)); err != nil {
		return nil, err
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}
	_ = file.Close()
	return header, nil
}

func (s *EnglishWordService) cacheGeneratedAudioURL(text, voice, audioURL string, timeoutMS int64, saveFolder string) (string, error) {
	audioURL = strings.TrimSpace(audioURL)
	if audioURL == "" {
		return "", errors.New("音频地址为空")
	}
	saveFolder = strings.TrimSpace(saveFolder)
	if saveFolder == "" {
		saveFolder = "english-learn/audio"
	}

	if strings.Contains(strings.ToLower(audioURL), "/english-learn/audio/") {
		return audioURL, nil
	}

	if timeoutMS < 1000 {
		timeoutMS = 1000
	}
	if timeoutMS > 60000 {
		timeoutMS = 60000
	}

	client := &http.Client{Timeout: time.Duration(timeoutMS) * time.Millisecond}
	resp, err := client.Get(audioURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("下载TTS音频失败(%d)", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return "", err
	}
	if len(body) == 0 {
		return "", errors.New("下载到的TTS音频为空")
	}

	ext := detectAudioExtFromURL(audioURL)
	if ext == "" {
		ext = detectAudioExtFromContentType(resp.Header.Get("Content-Type"))
	}
	fileName := buildAudioUploadFileName(text, voice, path.Base(audioURL), ext)
	header, err := buildMultipartFileHeader(fileName, resp.Header.Get("Content-Type"), body)
	if err != nil {
		return "", err
	}

	oss := upload.NewOss()
	savedURL, _, err := upload.UploadFileToFolder(oss, header, saveFolder)
	if err != nil {
		return "", err
	}

	return savedURL, nil
}

func (s *EnglishWordService) buildCategoryWordAssociations(tx *gorm.DB, wordID uint, categoryIDs []uint, chapterIDs []uint) ([]model.EnglishCategoryWord, error) {
	categoryIDs = uniqueUintIDs(categoryIDs)
	chapterIDs = uniqueUintIDs(chapterIDs)

	chapterCategoryMap := make(map[uint]uint)
	if len(chapterIDs) > 0 {
		var chapters []model.EnglishChapter
		if err := tx.Select("id", "category_id").Where("id IN ?", chapterIDs).Find(&chapters).Error; err != nil {
			return nil, err
		}
		if len(chapters) != len(chapterIDs) {
			return nil, errors.New("存在无效章节ID")
		}
		for _, chapter := range chapters {
			chapterCategoryMap[chapter.ID] = chapter.CategoryID
			if chapter.CategoryID > 0 {
				categoryIDs = append(categoryIDs, chapter.CategoryID)
			}
		}
	}

	categoryIDs = uniqueUintIDs(categoryIDs)
	associations := make([]model.EnglishCategoryWord, 0, len(categoryIDs)+len(chapterIDs))
	dupGuard := make(map[string]struct{})

	for _, categoryID := range categoryIDs {
		if categoryID == 0 {
			continue
		}
		key := fmt.Sprintf("%d-0", categoryID)
		if _, exists := dupGuard[key]; exists {
			continue
		}
		dupGuard[key] = struct{}{}
		associations = append(associations, model.EnglishCategoryWord{
			CategoryID: categoryID,
			WordID:     wordID,
		})
	}

	for _, chapterID := range chapterIDs {
		if chapterID == 0 {
			continue
		}
		categoryID := chapterCategoryMap[chapterID]
		key := fmt.Sprintf("%d-%d", categoryID, chapterID)
		if _, exists := dupGuard[key]; exists {
			continue
		}
		dupGuard[key] = struct{}{}
		associations = append(associations, model.EnglishCategoryWord{
			CategoryID: categoryID,
			ChapterID:  chapterID,
			WordID:     wordID,
		})
	}

	return associations, nil
}

func (s *EnglishWordService) buildWordSentenceEntities(word string, sentenceReqs []request.EnglishWordSentenceReq) []model.EnglishWordSentence {
	if len(sentenceReqs) == 0 {
		return []model.EnglishWordSentence{}
	}

	entities := make([]model.EnglishWordSentence, 0, len(sentenceReqs))
	for index, item := range sentenceReqs {
		source := strings.TrimSpace(item.Source)
		if source == "" {
			continue
		}

		entity := model.EnglishWordSentence{
			Source:    source,
			Translate: strings.TrimSpace(item.Translate),
			AudioUS:   strings.TrimSpace(item.AudioUS),
			AudioUK:   strings.TrimSpace(item.AudioUK),
			VideoID:   item.VideoID,
			Sort:      item.Sort,
		}
		if entity.Sort == 0 {
			entity.Sort = index + 1
		}

		if entity.AudioUS == "" {
			audioUS, err := s.generateWordAudioURLWithFolder(source, "learning_tts_voice_us", "LEARNING_TTS_VOICE_US", defaultLearningTTSVoiceUS, buildWordAudioFolder(word))
			if err == nil {
				entity.AudioUS = audioUS
			} else if !errors.Is(err, errLearningTTSDisabled) && !errors.Is(err, errLearningTTSProviderURLNotConfigured) {
				logEnglishLearningWarn("自动生成造句美式发音失败", zap.String("word", word), zap.String("sentence", source), zap.Error(err))
			}
		}

		if entity.AudioUK == "" {
			audioUK, err := s.generateWordAudioURLWithFolder(source, "learning_tts_voice_uk", "LEARNING_TTS_VOICE_UK", defaultLearningTTSVoiceUK, buildWordAudioFolder(word))
			if err == nil {
				entity.AudioUK = audioUK
			} else if !errors.Is(err, errLearningTTSDisabled) && !errors.Is(err, errLearningTTSProviderURLNotConfigured) {
				logEnglishLearningWarn("自动生成造句英式发音失败", zap.String("word", word), zap.String("sentence", source), zap.Error(err))
			}
		}

		entities = append(entities, entity)
	}

	return entities
}

func (s *EnglishWordService) replaceWordSentences(tx *gorm.DB, wordID uint, sentences []model.EnglishWordSentence) error {
	if err := tx.Unscoped().Where("word_id = ?", wordID).Delete(&model.EnglishWordSentence{}).Error; err != nil {
		return err
	}
	if len(sentences) == 0 {
		return nil
	}
	for i := range sentences {
		sentences[i].WordID = wordID
	}
	return tx.CreateInBatches(sentences, 100).Error
}

// CreateWord 创建单词库实体，集成自动TTS生成以及分类映射
func (s *EnglishWordService) CreateWord(word model.EnglishWord, categoryIDs []uint, chapterIDs []uint, sentenceReqs []request.EnglishWordSentenceReq) error {
	word.Word = strings.TrimSpace(word.Word)
	if word.Word == "" {
		return errors.New("单词不能为空")
	}
	normalizedWord := normalizeWordKey(word.Word)
	sentenceEntities := s.buildWordSentenceEntities(word.Word, sentenceReqs)

	// 2. 智能 TTS 旁路服务：未手工上传音频时，尝试通过外部TTS服务自动生成。
	if word.AudioUS == "" {
		audioUS, err := s.generateWordAudioURL(word.Word, "learning_tts_voice_us", "LEARNING_TTS_VOICE_US", defaultLearningTTSVoiceUS)
		if err == nil {
			word.AudioUS = audioUS
		} else if !errors.Is(err, errLearningTTSDisabled) && !errors.Is(err, errLearningTTSProviderURLNotConfigured) {
			global.GVA_LOG.Warn("自动生成美式发音失败", zap.String("word", word.Word), zap.Error(err))
		}
	}
	if word.AudioUK == "" {
		audioUK, err := s.generateWordAudioURL(word.Word, "learning_tts_voice_uk", "LEARNING_TTS_VOICE_UK", defaultLearningTTSVoiceUK)
		if err == nil {
			word.AudioUK = audioUK
		} else if !errors.Is(err, errLearningTTSDisabled) && !errors.Is(err, errLearningTTSProviderURLNotConfigured) {
			global.GVA_LOG.Warn("自动生成英式发音失败", zap.String("word", word.Word), zap.Error(err))
		}
	}

	// 3. 强一致性事务控制：创建主词条并建立多对多关系
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.EnglishWord{}).Where("LOWER(TRIM(word)) = ?", normalizedWord).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("单词已存在不可重复添加")
		}

		if err := tx.Create(&word).Error; err != nil {
			return err
		}

		associations, err := s.buildCategoryWordAssociations(tx, word.ID, categoryIDs, chapterIDs)
		if err != nil {
			return err
		}
		for _, assoc := range associations {
			if err = tx.Create(&assoc).Error; err != nil {
				return err
			}
		}

		if err = s.replaceWordSentences(tx, word.ID, sentenceEntities); err != nil {
			return err
		}
		return nil
	})
}

func (s *EnglishWordService) generateWordAudioURL(text, voiceConfigKey, voiceEnvKey, defaultVoice string) (string, error) {
	return s.generateWordAudioURLWithFolder(text, voiceConfigKey, voiceEnvKey, defaultVoice, buildWordAudioFolder(text))
}

func (s *EnglishWordService) generateWordAudioURLWithFolder(text, voiceConfigKey, voiceEnvKey, defaultVoice, saveFolder string) (string, error) {
	voice := getLearningStringSetting(voiceConfigKey, voiceEnvKey, defaultVoice)
	return s.requestTTSAudioURLWithSaveFolder(text, voice, true, saveFolder)
}

func (s *EnglishWordService) requestTTSAudioURL(text, voice string) (string, error) {
	return s.requestTTSAudioURLWithSaveFolder(text, voice, true, buildWordAudioFolder(text))
}

func (s *EnglishWordService) requestTTSAudioURLWithOptions(text, voice string, persistToStore bool) (string, error) {
	return s.requestTTSAudioURLWithSaveFolder(text, voice, persistToStore, buildWordAudioFolder(text))
}

func signLearningAudioURL(rawAudioURL string) string {
	rawAudioURL = strings.TrimSpace(rawAudioURL)
	if rawAudioURL == "" {
		return ""
	}

	cdnDomain := strings.TrimSpace(global.GVA_CONFIG.Hotlink.CdnDomain)
	if cdnDomain == "" {
		return rawAudioURL
	}

	if strings.HasPrefix(rawAudioURL, "http://") || strings.HasPrefix(rawAudioURL, "https://") {
		audioParsed, audioErr := url.Parse(rawAudioURL)
		cdnParsed, cdnErr := url.Parse(cdnDomain)
		if audioErr != nil || cdnErr != nil {
			return rawAudioURL
		}
		if !strings.EqualFold(strings.TrimSpace(audioParsed.Host), strings.TrimSpace(cdnParsed.Host)) {
			return rawAudioURL
		}
	}

	return SignLearningVideoURL(rawAudioURL)
}

func applyWordAudioSignatures(word *model.EnglishWord) {
	if word == nil {
		return
	}
	word.AudioUS = signLearningAudioURL(word.AudioUS)
	word.AudioUK = signLearningAudioURL(word.AudioUK)
}

func applySentenceAudioSignatures(sentences []model.EnglishWordSentence) {
	for i := range sentences {
		sentences[i].AudioUS = signLearningAudioURL(sentences[i].AudioUS)
		sentences[i].AudioUK = signLearningAudioURL(sentences[i].AudioUK)
	}
}

func (s *EnglishWordService) requestTTSAudioURLWithSaveFolder(text, voice string, persistToStore bool, saveFolder string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", errors.New("单词不能为空")
	}

	enabled := utils.GetBoolSetting("learning_tts_enabled", "LEARNING_TTS_ENABLED", false)
	if !enabled {
		return "", errLearningTTSDisabled
	}

	providerURL := getLearningStringSetting("learning_tts_provider_url", "LEARNING_TTS_PROVIDER_URL", "")
	if providerURL == "" {
		fallback := buildFreeFallbackTTSAudioURL(text, voice)
		if fallback == "" {
			return "", errLearningTTSProviderURLNotConfigured
		}
		return s.persistAudioURLIfNeeded(text, voice, fallback, timeoutMSForTTS(), persistToStore, saveFolder), nil
	}

	timeoutMS := timeoutMSForTTS()

	bodyBytes, err := json.Marshal(map[string]string{
		"text":   text,
		"voice":  strings.TrimSpace(voice),
		"format": "mp3",
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, providerURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey := getLearningStringSetting("learning_tts_api_key", "LEARNING_TTS_API_KEY", ""); apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	httpClient := &http.Client{Timeout: time.Duration(timeoutMS) * time.Millisecond}
	resp, err := httpClient.Do(req)
	if err != nil {
		fallback := ""
		if isLocalTTSProvider(providerURL) {
			fallback = buildFreeFallbackTTSAudioURL(text, voice)
		}
		if fallback != "" {
			logEnglishLearningWarn("主TTS服务不可用，回退免费TTS", zap.String("provider", providerURL), zap.Error(err))
			return s.persistAudioURLIfNeeded(text, voice, fallback, timeoutMS, persistToStore, saveFolder), nil
		}
		return "", fmt.Errorf("TTS服务不可用(provider=%s): %w", providerURL, err)
	}
	defer resp.Body.Close()

	rawResp, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		errBody := strings.TrimSpace(string(rawResp))
		if len(errBody) > 160 {
			errBody = errBody[:160] + "..."
		}
		fallback := ""
		if isLocalTTSProvider(providerURL) {
			fallback = buildFreeFallbackTTSAudioURL(text, voice)
		}
		if fallback != "" {
			logEnglishLearningWarn("主TTS服务状态异常，回退免费TTS", zap.String("provider", providerURL), zap.Int("status", resp.StatusCode))
			return s.persistAudioURLIfNeeded(text, voice, fallback, timeoutMS, persistToStore, saveFolder), nil
		}
		if errBody != "" {
			return "", fmt.Errorf("TTS服务响应异常(%d, provider=%s): %s", resp.StatusCode, providerURL, errBody)
		}
		return "", fmt.Errorf("TTS服务响应异常(%d, provider=%s)", resp.StatusCode, providerURL)
	}

	audioURL := extractTTSAudioURL(rawResp)
	if audioURL == "" {
		fallback := ""
		if isLocalTTSProvider(providerURL) {
			fallback = buildFreeFallbackTTSAudioURL(text, voice)
		}
		if fallback != "" {
			logEnglishLearningWarn("主TTS服务未返回音频地址，回退免费TTS", zap.String("provider", providerURL))
			return s.persistAudioURLIfNeeded(text, voice, fallback, timeoutMS, persistToStore, saveFolder), nil
		}
		return "", errLearningTTSMissingAudioURL
	}

	return s.persistAudioURLIfNeeded(text, voice, audioURL, timeoutMS, persistToStore, saveFolder), nil
}

func timeoutMSForTTS() int64 {
	timeoutMS := utils.GetInt64Setting("learning_tts_timeout_ms", "LEARNING_TTS_TIMEOUT_MS", defaultLearningTTSTimeoutMS)
	if timeoutMS < 1000 {
		timeoutMS = 1000
	}
	if timeoutMS > 60000 {
		timeoutMS = 60000
	}
	return timeoutMS
}

func (s *EnglishWordService) persistAudioURLIfNeeded(text, voice, audioURL string, timeoutMS int64, persistToStore bool, saveFolder string) string {
	audioURL = strings.TrimSpace(audioURL)
	if !persistToStore || audioURL == "" {
		return audioURL
	}

	cachedURL, cacheErr := s.cacheGeneratedAudioURL(text, voice, audioURL, timeoutMS, saveFolder)
	if cacheErr != nil {
		logEnglishLearningWarn("在线TTS音频落盘失败，回退原始地址", zap.String("word", text), zap.String("voice", voice), zap.String("audioURL", audioURL), zap.Error(cacheErr))
		return audioURL
	}
	if strings.TrimSpace(cachedURL) != "" {
		return cachedURL
	}
	return audioURL
}

func buildFreeFallbackTTSAudioURL(text, voice string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	baseURL := strings.TrimSpace(getLearningStringSetting("learning_tts_fallback_url", "LEARNING_TTS_FALLBACK_URL", "https://dict.youdao.com/dictvoice"))
	if baseURL == "" {
		return ""
	}

	voiceType := "1"
	voiceLower := strings.ToLower(strings.TrimSpace(voice))
	if strings.Contains(voiceLower, "uk") || strings.Contains(voiceLower, "gb") {
		voiceType = "2"
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	query := parsed.Query()
	query.Set("audio", text)
	query.Set("type", voiceType)
	query.Set("le", "en")
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func isLocalTTSProvider(providerURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(providerURL))
	if err != nil {
		return false
	}
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if host != "127.0.0.1" && host != "localhost" {
		return false
	}
	port := strings.TrimSpace(parsed.Port())
	if port != "19090" {
		return false
	}
	return strings.Contains(strings.ToLower(strings.TrimSpace(parsed.Path)), "tts")
}

func extractTTSAudioURL(raw []byte) string {
	rawText := strings.TrimSpace(string(raw))
	if strings.HasPrefix(rawText, "http://") || strings.HasPrefix(rawText, "https://") {
		return rawText
	}

	payload := map[string]interface{}{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}

	paths := [][]string{
		{"audioUrl"},
		{"audioURL"},
		{"url"},
		{"data", "audioUrl"},
		{"data", "audioURL"},
		{"data", "url"},
	}

	for _, p := range paths {
		if value := lookupNestedString(payload, p...); value != "" {
			return value
		}
	}

	return ""
}

func lookupNestedString(payload map[string]interface{}, path ...string) string {
	var cursor interface{} = payload
	for _, key := range path {
		asMap, ok := cursor.(map[string]interface{})
		if !ok {
			return ""
		}
		cursor = asMap[key]
	}

	value, ok := cursor.(string)
	if !ok {
		return ""
	}

	return strings.TrimSpace(value)
}

func (s *EnglishWordService) GetWord(id uint) (model.EnglishWord, error) {
	var word model.EnglishWord
	err := global.GVA_DB.Where("id = ?", id).First(&word).Error
	if err == nil {
		applyWordAudioSignatures(&word)
	}
	return word, err
}

func (s *EnglishWordService) GetWordBindingIDs(wordID uint) ([]uint, []uint, error) {
	if wordID == 0 {
		return []uint{}, []uint{}, nil
	}

	var relations []model.EnglishCategoryWord
	if err := global.GVA_DB.
		Select("category_id", "chapter_id").
		Where("word_id = ?", wordID).
		Find(&relations).Error; err != nil {
		return nil, nil, err
	}

	categorySet := map[uint]struct{}{}
	chapterSet := map[uint]struct{}{}
	categoryIDs := make([]uint, 0, len(relations))
	chapterIDs := make([]uint, 0, len(relations))

	for _, relation := range relations {
		if relation.CategoryID > 0 {
			if _, exists := categorySet[relation.CategoryID]; !exists {
				categorySet[relation.CategoryID] = struct{}{}
				categoryIDs = append(categoryIDs, relation.CategoryID)
			}
		}
		if relation.ChapterID > 0 {
			if _, exists := chapterSet[relation.ChapterID]; !exists {
				chapterSet[relation.ChapterID] = struct{}{}
				chapterIDs = append(chapterIDs, relation.ChapterID)
			}
		}
	}

	return categoryIDs, chapterIDs, nil
}

func (s *EnglishWordService) UpdateWord(req request.UpdateEnglishWordReq) error {
	if req.ID == 0 {
		return errors.New("ID不能为空")
	}

	req.Word = strings.TrimSpace(req.Word)
	if req.Word == "" {
		return errors.New("单词不能为空")
	}
	normalizedWord := normalizeWordKey(req.Word)

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var existing model.EnglishWord
		if err := tx.Select("id").Where("id = ?", req.ID).First(&existing).Error; err != nil {
			return err
		}

		var duplicateCount int64
		if err := tx.Model(&model.EnglishWord{}).
			Where("LOWER(TRIM(word)) = ?", normalizedWord).
			Where("id <> ?", req.ID).
			Count(&duplicateCount).Error; err != nil {
			return err
		}
		if duplicateCount > 0 {
			return errors.New("单词已存在不可重复添加")
		}

		updates := map[string]interface{}{
			"word":           req.Word,
			"phonetic_us":    req.PhoneticUS,
			"phonetic_uk":    req.PhoneticUK,
			"part_of_speech": strings.TrimSpace(req.PartOfSpeech),
			"audio_us":       req.AudioUS,
			"audio_uk":       req.AudioUK,
			"explanation":    req.Explanation,
		}
		if err := tx.Model(&model.EnglishWord{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}

		if req.CategoryIDs == nil && req.ChapterIDs == nil && req.Sentences == nil {
			return nil
		}

		if req.CategoryIDs != nil || req.ChapterIDs != nil {
			if err := tx.Unscoped().Where("word_id = ?", req.ID).Delete(&model.EnglishCategoryWord{}).Error; err != nil {
				return err
			}

			associations, err := s.buildCategoryWordAssociations(tx, req.ID, req.CategoryIDs, req.ChapterIDs)
			if err != nil {
				return err
			}

			for _, assoc := range associations {
				if err = tx.Create(&assoc).Error; err != nil {
					return err
				}
			}
		}

		if req.Sentences != nil {
			sentenceEntities := s.buildWordSentenceEntities(req.Word, req.Sentences)
			if err := s.replaceWordSentences(tx, req.ID, sentenceEntities); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *EnglishWordService) DeleteWord(id uint) error {
	if id == 0 {
		return errors.New("ID不能为空")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var existing model.EnglishWord
		if err := tx.Unscoped().Select("id").Where("id = ?", id).First(&existing).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Where("word_id = ?", id).Delete(&model.EnglishCategoryWord{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("word_id = ?", id).Delete(&model.EnglishWordSentence{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("word_id = ?", id).Delete(&model.EnglishWordErrorLog{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("target_type = ? AND target_id = ?", 1, id).Delete(&model.UserCollection{}).Error; err != nil {
			return err
		}

		return tx.Unscoped().Delete(&model.EnglishWord{}, "id = ?", id).Error
	})
}

func (s *EnglishWordService) GetWordSentences(wordID uint) ([]model.EnglishWordSentence, error) {
	if wordID == 0 {
		return []model.EnglishWordSentence{}, nil
	}
	list := make([]model.EnglishWordSentence, 0)
	err := global.GVA_DB.Where("word_id = ?", wordID).Order("sort ASC, id ASC").Find(&list).Error
	if err == nil && len(list) > 0 {
		applySentenceAudioSignatures(list)
	}
	return list, err
}

func resolveRegenerateFlag(flag *bool) bool {
	if flag == nil {
		return true
	}
	return *flag
}

func (s *EnglishWordService) RegenerateWordAudio(req request.RegenerateWordAudioReq) error {
	if req.ID == 0 {
		return errors.New("ID不能为空")
	}

	regenerateUS := resolveRegenerateFlag(req.RegenerateUS)
	regenerateUK := resolveRegenerateFlag(req.RegenerateUK)
	if !regenerateUS && !regenerateUK {
		return errors.New("请至少选择一种发音进行重生成")
	}

	var word model.EnglishWord
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&word).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if regenerateUS {
		audioUS, err := s.generateWordAudioURL(word.Word, "learning_tts_voice_us", "LEARNING_TTS_VOICE_US", defaultLearningTTSVoiceUS)
		if err != nil {
			return err
		}
		updates["audio_us"] = audioUS
	}

	if regenerateUK {
		audioUK, err := s.generateWordAudioURL(word.Word, "learning_tts_voice_uk", "LEARNING_TTS_VOICE_UK", defaultLearningTTSVoiceUK)
		if err != nil {
			return err
		}
		updates["audio_uk"] = audioUK
	}

	if len(updates) == 0 {
		return nil
	}

	return global.GVA_DB.Model(&model.EnglishWord{}).Where("id = ?", req.ID).Updates(updates).Error
}

func (s *EnglishWordService) PreflightTTS(req request.TTSPreflightReq) (TTSPreflightResult, error) {
	result := TTSPreflightResult{}

	word := strings.TrimSpace(req.Word)
	if word == "" {
		word = "destiny"
	}
	result.Word = word

	checkUS := resolveRegenerateFlag(req.CheckUS)
	checkUK := resolveRegenerateFlag(req.CheckUK)
	if !checkUS && !checkUK {
		return result, errors.New("请至少选择一种发音进行预检")
	}

	start := time.Now()
	if checkUS {
		result.VoiceUS = getLearningStringSetting("learning_tts_voice_us", "LEARNING_TTS_VOICE_US", defaultLearningTTSVoiceUS)
		audioUS, err := s.requestTTSAudioURLWithOptions(word, result.VoiceUS, false)
		if err != nil {
			return result, fmt.Errorf("美式发音预检失败: %w", err)
		}
		result.AudioUS = audioUS
	}

	if checkUK {
		result.VoiceUK = getLearningStringSetting("learning_tts_voice_uk", "LEARNING_TTS_VOICE_UK", defaultLearningTTSVoiceUK)
		audioUK, err := s.requestTTSAudioURLWithOptions(word, result.VoiceUK, false)
		if err != nil {
			return result, fmt.Errorf("英式发音预检失败: %w", err)
		}
		result.AudioUK = audioUK
	}

	result.DurationMS = time.Since(start).Milliseconds()
	return result, nil
}

func buildWordListBaseQuery(db *gorm.DB, info request.WordListSearch) *gorm.DB {
	query := db.Table("english_words w").
		Joins("LEFT JOIN english_category_words cw ON cw.word_id = w.id AND cw.deleted_at IS NULL")
	keyword := strings.ToLower(strings.TrimSpace(info.Keyword))
	if info.CategoryID > 0 {
		query = query.Where("cw.category_id = ?", info.CategoryID)
	}
	if info.ChapterID > 0 {
		query = query.Where("cw.chapter_id = ?", info.ChapterID)
	}
	if keyword != "" {
		// FuzzySearch: nil 或 true → 模糊 LIKE; 显式 false → 精准 =
		if info.FuzzySearch != nil && !*info.FuzzySearch {
			query = query.Where("LOWER(TRIM(w.word)) = ?", keyword)
		} else {
			query = query.Where("LOWER(TRIM(w.word)) LIKE ?", "%"+keyword+"%")
		}
	}
	return query
}

func buildWordListCountQuery(db *gorm.DB, info request.WordListSearch) *gorm.DB {
	return buildWordListBaseQuery(db, info).Session(&gorm.Session{}).Distinct("w.id")
}

func buildWordListPageQuery(db *gorm.DB, info request.WordListSearch) *gorm.DB {
	return buildWordListBaseQuery(db, info).Session(&gorm.Session{}).
		Select("w.*").
		Group("w.id").
		Order("COALESCE(MIN(cw.sort), 0) ASC, MAX(cw.id) DESC")
}

func (s *EnglishWordService) GetWordListByChapter(info request.WordListSearch) (list []model.EnglishWord, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	err = buildWordListCountQuery(global.GVA_DB, info).Count(&total).Error
	if err != nil {
		return
	}
	err = buildWordListPageQuery(global.GVA_DB, info).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list).Error
	if err != nil || len(list) == 0 {
		return
	}

	wordIDs := make([]uint, 0, len(list))
	for _, item := range list {
		if item.ID > 0 {
			wordIDs = append(wordIDs, item.ID)
		}
	}
	wordIDs = uniqueUintIDs(wordIDs)
	if len(wordIDs) == 0 {
		return
	}

	var sentences []model.EnglishWordSentence
	err = global.GVA_DB.Where("word_id IN ?", wordIDs).Order("sort ASC, id ASC").Find(&sentences).Error
	if err != nil {
		return
	}

	sentenceMap := make(map[uint][]model.EnglishWordSentence, len(wordIDs))
	for _, sentence := range sentences {
		sentenceMap[sentence.WordID] = append(sentenceMap[sentence.WordID], sentence)
	}
	for i := range list {
		list[i].Sentences = sentenceMap[list[i].ID]
		applyWordAudioSignatures(&list[i])
		if len(list[i].Sentences) > 0 {
			applySentenceAudioSignatures(list[i].Sentences)
		}
	}
	return
}

func sanitizeCharToken(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	runes := []rune(raw)
	if len(runes) > 1 {
		return string(runes[:1])
	}
	return raw
}

func uniqueUintIDs(ids []uint) []uint {
	set := make(map[uint]struct{}, len(ids))
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, exists := set[id]; exists {
			continue
		}
		set[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func (s *EnglishWordService) ReportWordError(userID uint, req request.ReportWordErrorReq) error {
	if userID == 0 || req.WordID == 0 {
		return errors.New("参数错误")
	}

	var word model.EnglishWord
	if err := global.GVA_DB.Select("id").Where("id = ?", req.WordID).First(&word).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("单词不存在")
		}
		return err
	}

	logRow := model.EnglishWordErrorLog{
		UserID:           userID,
		WordID:           req.WordID,
		CategoryID:       req.CategoryID,
		ChapterID:        req.ChapterID,
		WrongCount:       1,
		LastWrongIndex:   req.WrongIndex,
		LastExpectedChar: sanitizeCharToken(req.ExpectedChar),
		LastInputChar:    sanitizeCharToken(req.InputChar),
	}

	updates := map[string]interface{}{
		"wrong_count":        gorm.Expr("wrong_count + ?", 1),
		"category_id":        req.CategoryID,
		"chapter_id":         req.ChapterID,
		"last_wrong_index":   req.WrongIndex,
		"last_expected_char": sanitizeCharToken(req.ExpectedChar),
		"last_input_char":    sanitizeCharToken(req.InputChar),
		"updated_at":         time.Now(),
	}

	return global.GVA_DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "word_id"}},
		DoUpdates: clause.Assignments(updates),
	}).Create(&logRow).Error
}

func (s *EnglishWordService) GetWordErrorLogList(userID uint, info request.WordErrorLogSearch) (list []WordErrorLogItem, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Model(&model.EnglishWordErrorLog{}).Where("user_id = ?", userID)
	if info.CategoryID > 0 {
		db = db.Where("category_id = ?", info.CategoryID)
	}
	if info.ChapterID > 0 {
		db = db.Where("chapter_id = ?", info.ChapterID)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	var logs []model.EnglishWordErrorLog
	err = db.Order("wrong_count DESC, updated_at DESC, id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&logs).Error
	if err != nil || len(logs) == 0 {
		return
	}

	wordIDs := make([]uint, 0, len(logs))
	for _, row := range logs {
		wordIDs = append(wordIDs, row.WordID)
	}
	wordIDs = uniqueUintIDs(wordIDs)

	wordMap := make(map[uint]model.EnglishWord)
	if len(wordIDs) > 0 {
		var words []model.EnglishWord
		err = global.GVA_DB.Where("id IN ?", wordIDs).Find(&words).Error
		if err != nil {
			return
		}
		for _, word := range words {
			wordMap[word.ID] = word
		}
	}

	list = make([]WordErrorLogItem, 0, len(logs))
	for _, row := range logs {
		item := WordErrorLogItem{EnglishWordErrorLog: row}
		if word, ok := wordMap[row.WordID]; ok {
			wordCopy := word
			item.Word = &wordCopy
		}
		list = append(list, item)
	}

	return
}

func (s *EnglishWordService) DeleteWordErrorLog(userID uint, wordID uint) error {
	if userID == 0 || wordID == 0 {
		return errors.New("参数错误")
	}

	return global.GVA_DB.Where("user_id = ? AND word_id = ?", userID, wordID).Delete(&model.EnglishWordErrorLog{}).Error
}

func normalizeI18nJSONWithZh(text string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return `{"zh":""}`
	}
	buf, err := json.Marshal(map[string]string{"zh": trimmed})
	if err != nil {
		return `{"zh":""}`
	}
	return string(buf)
}

func mergeI18nZhText(existingRaw, zhText string) string {
	zhText = strings.TrimSpace(zhText)
	if strings.TrimSpace(existingRaw) == "" {
		return normalizeI18nJSONWithZh(zhText)
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(existingRaw), &obj); err != nil || obj == nil {
		return normalizeI18nJSONWithZh(zhText)
	}
	obj["zh"] = zhText
	buf, err := json.Marshal(obj)
	if err != nil {
		return normalizeI18nJSONWithZh(zhText)
	}
	return string(buf)
}

func ensureWordAssociation(tx *gorm.DB, categoryID, chapterID, wordID uint) error {
	if categoryID == 0 || wordID == 0 {
		return errors.New("分类或单词参数错误")
	}

	var relation model.EnglishCategoryWord
	err := tx.Where("category_id = ? AND chapter_id = ? AND word_id = ?", categoryID, chapterID, wordID).First(&relation).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return tx.Create(&model.EnglishCategoryWord{
		CategoryID: categoryID,
		ChapterID:  chapterID,
		WordID:     wordID,
	}).Error
}

func (s *EnglishWordService) UpsertWordFromSQL(req request.UpsertWordFromSQLReq) (UpsertWordFromSQLResult, error) {
	result := UpsertWordFromSQLResult{}
	req.Word = strings.TrimSpace(req.Word)
	if req.Word == "" {
		return result, errors.New("单词不能为空")
	}
	if req.CategoryID == 0 {
		return result, errors.New("分类不能为空")
	}
	normalizedWord := normalizeWordKey(req.Word)
	translate := strings.TrimSpace(req.TranslateZh)

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var category model.EnglishCategory
		if err := tx.Select("id").Where("id = ?", req.CategoryID).First(&category).Error; err != nil {
			return errors.New("分类不存在")
		}

		if req.ChapterID > 0 {
			var chapter model.EnglishChapter
			if err := tx.Select("id", "category_id").Where("id = ?", req.ChapterID).First(&chapter).Error; err != nil {
				return errors.New("章节不存在")
			}
			if chapter.CategoryID != req.CategoryID {
				return errors.New("章节不属于所选分类")
			}
		}

		var word model.EnglishWord
		err := tx.Where("LOWER(TRIM(word)) = ?", normalizedWord).First(&word).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			word = model.EnglishWord{
				Word:        req.Word,
				Explanation: normalizeI18nJSONWithZh(translate),
			}
			if req.GenerateAudio {
				if audioUS, audioErr := s.generateWordAudioURL(req.Word, "learning_tts_voice_us", "LEARNING_TTS_VOICE_US", defaultLearningTTSVoiceUS); audioErr == nil {
					word.AudioUS = audioUS
					result.AudioUSGenerated = true
				}
				if audioUK, audioErr := s.generateWordAudioURL(req.Word, "learning_tts_voice_uk", "LEARNING_TTS_VOICE_UK", defaultLearningTTSVoiceUK); audioErr == nil {
					word.AudioUK = audioUK
					result.AudioUKGenerated = true
				}
			}
			if err = tx.Create(&word).Error; err != nil {
				return err
			}
			result.Action = "create"
		} else {
			updates := map[string]interface{}{
				"explanation": mergeI18nZhText(word.Explanation, translate),
			}
			if req.GenerateAudio {
				if strings.TrimSpace(word.AudioUS) == "" {
					if audioUS, audioErr := s.generateWordAudioURL(word.Word, "learning_tts_voice_us", "LEARNING_TTS_VOICE_US", defaultLearningTTSVoiceUS); audioErr == nil {
						updates["audio_us"] = audioUS
						result.AudioUSGenerated = true
					}
				}
				if strings.TrimSpace(word.AudioUK) == "" {
					if audioUK, audioErr := s.generateWordAudioURL(word.Word, "learning_tts_voice_uk", "LEARNING_TTS_VOICE_UK", defaultLearningTTSVoiceUK); audioErr == nil {
						updates["audio_uk"] = audioUK
						result.AudioUKGenerated = true
					}
				}
			}
			if err = tx.Model(&model.EnglishWord{}).Where("id = ?", word.ID).Updates(updates).Error; err != nil {
				return err
			}
			result.Action = "update"
		}

		if err = ensureWordAssociation(tx, req.CategoryID, 0, word.ID); err != nil {
			return err
		}
		if req.ChapterID > 0 {
			if err = ensureWordAssociation(tx, req.CategoryID, req.ChapterID, word.ID); err != nil {
				return err
			}
		}

		result.WordID = word.ID
		return nil
	})

	return result, err
}

// BatchCheckWords 批量检查单词是否存在于词库（不创建，仅查询）
func (s *EnglishWordService) BatchCheckWords(req request.BatchCheckWordsReq) ([]request.BatchCheckWordsItem, error) {
	if len(req.Words) == 0 {
		return nil, errors.New("单词列表不能为空")
	}

	// 去重并规范化
	wordSet := make(map[string]bool, len(req.Words))
	normalizedWords := make([]string, 0, len(req.Words))
	for _, rawWord := range req.Words {
		word := strings.TrimSpace(strings.ToLower(rawWord))
		if word == "" || wordSet[word] {
			continue
		}
		wordSet[word] = true
		normalizedWords = append(normalizedWords, word)
	}

	if len(normalizedWords) == 0 {
		return nil, errors.New("无可检查的有效单词")
	}

	// 批量查询存在的单词
	var existingWords []model.EnglishWord
	err := global.GVA_DB.Where("LOWER(TRIM(word)) IN ?", normalizedWords).Find(&existingWords).Error
	if err != nil {
		return nil, err
	}

	// 构建已存在单词的 map
	existingMap := make(map[string]uint, len(existingWords))
	for _, ew := range existingWords {
		existingMap[strings.ToLower(strings.TrimSpace(ew.Word))] = ew.ID
	}

	// 构建结果
	result := make([]request.BatchCheckWordsItem, 0, len(normalizedWords))
	for _, word := range normalizedWords {
		wordID, exists := existingMap[word]
		result = append(result, request.BatchCheckWordsItem{
			Word:   word,
			WordID: wordID,
			Exists: exists,
		})
	}
	return result, nil
}

func normalizeLangCodes(input []string) []string {
	if len(input) == 0 {
		return []string{}
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(input))
	for _, item := range input {
		lang := strings.ToLower(strings.TrimSpace(item))
		if lang == "" {
			continue
		}
		if _, exists := seen[lang]; exists {
			continue
		}
		seen[lang] = struct{}{}
		out = append(out, lang)
	}
	sort.Strings(out)
	return out
}

func parseI18nMap(raw string) map[string]string {
	out := map[string]string{}
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return out
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &obj); err != nil || obj == nil {
		out["zh"] = trimmed
		return out
	}

	for key, value := range obj {
		lang := strings.ToLower(strings.TrimSpace(key))
		if lang == "" {
			continue
		}
		text := strings.TrimSpace(fmt.Sprintf("%v", value))
		if text != "" {
			out[lang] = text
		}
	}
	return out
}

func stringifyI18nMap(data map[string]string) string {
	cleaned := map[string]string{}
	for key, value := range data {
		lang := strings.ToLower(strings.TrimSpace(key))
		text := strings.TrimSpace(value)
		if lang == "" || text == "" {
			continue
		}
		cleaned[lang] = text
	}
	if len(cleaned) == 0 {
		return `{"zh":""}`
	}
	buf, err := json.Marshal(cleaned)
	if err != nil {
		return `{"zh":""}`
	}
	return string(buf)
}

func pickPhonetic(entries []dictionaryEntry, preferUS bool) string {
	if len(entries) == 0 {
		return ""
	}
	candidates := make([]string, 0, 4)
	for _, entry := range entries {
		for _, item := range entry.Phonetics {
			text := strings.TrimSpace(item.Text)
			if text == "" {
				continue
			}
			audio := strings.ToLower(strings.TrimSpace(item.Audio))
			if preferUS {
				if strings.Contains(audio, "-us") || strings.Contains(audio, "/us") || strings.Contains(audio, "en-us") {
					return text
				}
			} else {
				if strings.Contains(audio, "-uk") || strings.Contains(audio, "/uk") || strings.Contains(audio, "en-gb") {
					return text
				}
			}
			candidates = append(candidates, text)
		}
	}
	if len(candidates) == 0 {
		return ""
	}
	return candidates[0]
}

func extractDictionaryData(entries []dictionaryEntry) (string, string, string, string, []string) {
	phoneticUS := pickPhonetic(entries, true)
	phoneticUK := pickPhonetic(entries, false)

	posSet := map[string]struct{}{}
	posList := make([]string, 0, 4)
	definition := ""
	exampleSet := map[string]struct{}{}
	examples := make([]string, 0, 4)

	for _, entry := range entries {
		for _, meaning := range entry.Meanings {
			pos := strings.TrimSpace(meaning.PartOfSpeech)
			if pos != "" {
				if _, exists := posSet[pos]; !exists {
					posSet[pos] = struct{}{}
					posList = append(posList, pos)
				}
			}
			for _, def := range meaning.Definitions {
				if definition == "" {
					definition = strings.TrimSpace(def.Definition)
				}
				example := strings.TrimSpace(def.Example)
				if example == "" {
					continue
				}
				if _, exists := exampleSet[example]; exists {
					continue
				}
				exampleSet[example] = struct{}{}
				examples = append(examples, example)
			}
		}
	}

	partOfSpeech := strings.Join(posList, ", ")
	return phoneticUS, phoneticUK, partOfSpeech, definition, examples
}

func fetchDictionaryEntries(word string) ([]dictionaryEntry, error) {
	endpoint := "https://api.dictionaryapi.dev/api/v2/entries/en/" + url.PathEscape(strings.ToLower(strings.TrimSpace(word)))
	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []dictionaryEntry{}, nil
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("dictionaryapi状态异常(%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var data []dictionaryEntry
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func requestTranslateText(endpoint, serviceName, sourceLang, targetLang, text string) (string, error) {
	payload, err := json.Marshal(map[string]string{
		"text":        text,
		"source_lang": sourceLang,
		"target_lang": targetLang,
		"service":     serviceName,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("翻译接口状态异常(%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var out struct {
		TranslatedText string `json:"translated_text"`
	}
	if err = json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	translated := strings.TrimSpace(out.TranslatedText)
	if translated == "" {
		return "", errors.New("翻译接口未返回translated_text")
	}
	return translated, nil
}

func (s *EnglishWordService) BatchFillWordFromDictionary(req request.BatchFillWordFromDictionaryReq) (BatchFillWordFromDictionaryResult, error) {
	result := BatchFillWordFromDictionaryResult{CategoryID: req.CategoryID, Items: []BatchFillWordFromDictionaryItem{}}
	if req.CategoryID == 0 {
		return result, errors.New("分类不能为空")
	}
	if !req.FillPhonetic && !req.FillPartOfSpeech && !req.FillExplanation && !req.FillSentences && !req.TranslateExplanation && !req.TranslateSentences {
		return result, errors.New("请至少选择一个补全或翻译项")
	}

	targetLangs := normalizeLangCodes(req.TargetLangs)
	if (req.TranslateExplanation || req.TranslateSentences) && len(targetLangs) == 0 {
		return result, errors.New("启用翻译时必须选择目标语言")
	}

	serviceName := strings.ToLower(strings.TrimSpace(req.TranslateService))
	if serviceName == "" {
		serviceName = "gtx"
	}
	if serviceName != "gtx" && serviceName != "edge" {
		return result, errors.New("translateService仅支持 gtx 或 edge")
	}

	translateURL := strings.TrimSpace(req.TranslateURL)
	if translateURL == "" {
		translateURL = getLearningStringSetting("learning_translate_provider_url", "LEARNING_TRANSLATE_PROVIDER_URL", "http://localhost:3000/api/free-translate")
	}

	query := global.GVA_DB.Table("english_words w").
		Select("w.*").
		Joins("JOIN english_category_words cw ON cw.word_id = w.id AND cw.deleted_at IS NULL").
		Where("cw.category_id = ?", req.CategoryID).
		Group("w.id").
		Order("COALESCE(MIN(cw.sort), 0) ASC, w.id ASC")
	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}

	words := make([]model.EnglishWord, 0)
	if err := query.Scan(&words).Error; err != nil {
		return result, err
	}

	result.Total = len(words)
	if len(words) == 0 {
		return result, nil
	}

	translateCache := map[string]string{}
	translateText := func(text, target string) (string, error) {
		cacheKey := strings.ToLower(strings.TrimSpace(text)) + "|" + target
		if value, ok := translateCache[cacheKey]; ok {
			return value, nil
		}
		translated, err := requestTranslateText(translateURL, serviceName, "en", target, text)
		if err != nil {
			return "", err
		}
		translateCache[cacheKey] = translated
		return translated, nil
	}

	for _, row := range words {
		item := BatchFillWordFromDictionaryItem{Word: row.Word, WordID: row.ID, Applied: []string{}, Translated: []string{}}
		entries, err := fetchDictionaryEntries(row.Word)
		if err != nil {
			item.Status = "failed"
			item.Message = err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
		}
		if len(entries) == 0 {
			item.Status = "skipped"
			item.Message = "dictionaryapi未找到该单词"
			result.Skipped++
			result.Items = append(result.Items, item)
			continue
		}

		phoneticUS, phoneticUK, partOfSpeech, englishDefinition, examples := extractDictionaryData(entries)
		warnings := make([]string, 0)

		err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			var word model.EnglishWord
			if err := tx.Where("id = ?", row.ID).First(&word).Error; err != nil {
				return err
			}

			updates := map[string]interface{}{}

			if req.FillPhonetic {
				if phoneticUS != "" && (req.OverwriteExisting || strings.TrimSpace(word.PhoneticUS) == "") {
					updates["phonetic_us"] = phoneticUS
					item.Applied = append(item.Applied, "phoneticUs")
				}
				if phoneticUK != "" && (req.OverwriteExisting || strings.TrimSpace(word.PhoneticUK) == "") {
					updates["phonetic_uk"] = phoneticUK
					item.Applied = append(item.Applied, "phoneticUk")
				}
			}

			if req.FillPartOfSpeech && partOfSpeech != "" && (req.OverwriteExisting || strings.TrimSpace(word.PartOfSpeech) == "") {
				updates["part_of_speech"] = partOfSpeech
				item.Applied = append(item.Applied, "partOfSpeech")
			}

			explanationMap := parseI18nMap(word.Explanation)
			explanationChanged := false
			if req.FillExplanation && englishDefinition != "" {
				if req.OverwriteExisting || strings.TrimSpace(explanationMap["en"]) == "" {
					explanationMap["en"] = englishDefinition
					explanationChanged = true
					item.Applied = append(item.Applied, "explanation")
				}
			}

			if req.TranslateExplanation {
				sourceText := strings.TrimSpace(explanationMap["en"])
				if sourceText != "" {
					for _, lang := range targetLangs {
						if lang == "en" {
							continue
						}
						if !req.OverwriteExisting && strings.TrimSpace(explanationMap[lang]) != "" {
							continue
						}
						translated, transErr := translateText(sourceText, lang)
						if transErr != nil {
							warnings = append(warnings, fmt.Sprintf("释义->%s翻译失败", lang))
							continue
						}
						explanationMap[lang] = translated
						explanationChanged = true
						item.Translated = append(item.Translated, "explanation:"+lang)
					}
				}
			}

			if explanationChanged {
				updates["explanation"] = stringifyI18nMap(explanationMap)
			}

			if len(updates) > 0 {
				if err := tx.Model(&model.EnglishWord{}).Where("id = ?", word.ID).Updates(updates).Error; err != nil {
					return err
				}
			}

			needLoadSentences := req.FillSentences || req.TranslateSentences
			if !needLoadSentences {
				return nil
			}

			sentences := make([]model.EnglishWordSentence, 0)
			if err := tx.Where("word_id = ?", word.ID).Order("sort ASC, id ASC").Find(&sentences).Error; err != nil {
				return err
			}

			if req.FillSentences {
				exampleList := make([]string, 0, len(examples))
				seen := map[string]struct{}{}
				for _, sentence := range examples {
					trimmed := strings.TrimSpace(sentence)
					if trimmed == "" {
						continue
					}
					key := strings.ToLower(trimmed)
					if _, exists := seen[key]; exists {
						continue
					}
					seen[key] = struct{}{}
					exampleList = append(exampleList, trimmed)
				}

				if req.OverwriteExisting {
					if err := tx.Unscoped().Where("word_id = ?", word.ID).Delete(&model.EnglishWordSentence{}).Error; err != nil {
						return err
					}
					sentences = []model.EnglishWordSentence{}
				}

				existingSource := map[string]struct{}{}
				maxSort := 0
				for _, sentence := range sentences {
					existingSource[strings.ToLower(strings.TrimSpace(sentence.Source))] = struct{}{}
					if sentence.Sort > maxSort {
						maxSort = sentence.Sort
					}
				}

				newRows := make([]model.EnglishWordSentence, 0)
				for _, source := range exampleList {
					if _, exists := existingSource[strings.ToLower(source)]; exists {
						continue
					}
					maxSort++
					newRows = append(newRows, model.EnglishWordSentence{
						WordID: word.ID,
						Source: source,
						Sort:   maxSort,
					})
				}

				if len(newRows) > 0 {
					if err := tx.Create(&newRows).Error; err != nil {
						return err
					}
					item.Applied = append(item.Applied, fmt.Sprintf("sentences(+%d)", len(newRows)))
					sentences = append(sentences, newRows...)
				}
			}

			if req.TranslateSentences {
				for _, sentence := range sentences {
					sourceText := strings.TrimSpace(sentence.Source)
					if sourceText == "" {
						continue
					}
					translateMap := parseI18nMap(sentence.Translate)
					changed := false
					for _, lang := range targetLangs {
						if lang == "en" {
							continue
						}
						if !req.OverwriteExisting && strings.TrimSpace(translateMap[lang]) != "" {
							continue
						}
						translated, transErr := translateText(sourceText, lang)
						if transErr != nil {
							warnings = append(warnings, fmt.Sprintf("例句[%d]->%s翻译失败", sentence.ID, lang))
							continue
						}
						translateMap[lang] = translated
						changed = true
						item.Translated = append(item.Translated, fmt.Sprintf("sentence:%d:%s", sentence.ID, lang))
					}
					if changed {
						if err := tx.Model(&model.EnglishWordSentence{}).Where("id = ?", sentence.ID).Update("translate", stringifyI18nMap(translateMap)).Error; err != nil {
							return err
						}
					}
				}
			}

			return nil
		})

		if err != nil {
			item.Status = "failed"
			item.Message = err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
		}

		if len(item.Applied) == 0 && len(item.Translated) == 0 {
			item.Status = "skipped"
			item.Message = "无可更新字段"
			result.Skipped++
		} else {
			item.Status = "success"
			if len(warnings) > 0 {
				item.Message = strings.Join(warnings, "; ")
			}
			result.Success++
		}
		result.Items = append(result.Items, item)
	}

	return result, nil
}

// BatchCreateWords 批量创建单词（仅入库单词本体，不生成TTS/关联等，用于字幕全部单词一键入库）
// 已存在的单词跳过，返回所有单词的 ID 列表
func (s *EnglishWordService) BatchCreateWords(req request.BatchCreateWordsReq) ([]request.BatchCreateWordsItem, error) {
	if len(req.Words) == 0 {
		return nil, errors.New("单词列表不能为空")
	}

	result := make([]request.BatchCreateWordsItem, 0, len(req.Words))

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, rawWord := range req.Words {
			word := strings.TrimSpace(strings.ToLower(rawWord))
			if word == "" {
				continue
			}

			var existing model.EnglishWord
			findErr := tx.Where("LOWER(TRIM(word)) = ?", word).First(&existing).Error
			if findErr == nil {
				// 已存在
				result = append(result, request.BatchCreateWordsItem{
					Word:   word,
					WordID: existing.ID,
					Newly:  false,
				})
				continue
			}
			if !errors.Is(findErr, gorm.ErrRecordNotFound) {
				return findErr
			}

			// 不存在，创建
			newWord := model.EnglishWord{
				Word:        word,
				Explanation: "{}",
			}
			if createErr := tx.Create(&newWord).Error; createErr != nil {
				return createErr
			}
			result = append(result, request.BatchCreateWordsItem{
				Word:   word,
				WordID: newWord.ID,
				Newly:  true,
			})
		}
		return nil
	})

	return result, err
}
