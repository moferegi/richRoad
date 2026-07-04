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
	"os"
	"path"
	"path/filepath"
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

func (s *EnglishWordService) cacheGeneratedAudioURL(text, voice, audioURL string, timeoutMS int64) (string, error) {
	audioURL = strings.TrimSpace(audioURL)
	if audioURL == "" {
		return "", errors.New("音频地址为空")
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
	savedURL, _, err := upload.UploadFileToFolder(oss, header, "english-learn/audio")
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

// CreateWord 创建单词库实体，集成自动TTS生成以及分类映射
func (s *EnglishWordService) CreateWord(word model.EnglishWord, categoryIDs []uint, chapterIDs []uint) error {
	word.Word = strings.TrimSpace(word.Word)
	if word.Word == "" {
		return errors.New("单词不能为空")
	}
	normalizedWord := normalizeWordKey(word.Word)

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
		return nil
	})
}

func (s *EnglishWordService) generateWordAudioURL(text, voiceConfigKey, voiceEnvKey, defaultVoice string) (string, error) {
	voice := getLearningStringSetting(voiceConfigKey, voiceEnvKey, defaultVoice)
	return s.requestTTSAudioURL(text, voice)
}

func (s *EnglishWordService) requestTTSAudioURL(text, voice string) (string, error) {
	return s.requestTTSAudioURLWithOptions(text, voice, true)
}

func (s *EnglishWordService) requestTTSAudioURLWithOptions(text, voice string, persistToStore bool) (string, error) {
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
		return "", errLearningTTSProviderURLNotConfigured
	}

	timeoutMS := utils.GetInt64Setting("learning_tts_timeout_ms", "LEARNING_TTS_TIMEOUT_MS", defaultLearningTTSTimeoutMS)
	if timeoutMS < 1000 {
		timeoutMS = 1000
	}
	if timeoutMS > 60000 {
		timeoutMS = 60000
	}

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
		return "", err
	}
	defer resp.Body.Close()

	rawResp, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("TTS服务响应异常(%d)", resp.StatusCode)
	}

	audioURL := extractTTSAudioURL(rawResp)
	if audioURL == "" {
		return "", errLearningTTSMissingAudioURL
	}

	if persistToStore {
		cachedURL, cacheErr := s.cacheGeneratedAudioURL(text, voice, audioURL, timeoutMS)
		if cacheErr != nil {
			logEnglishLearningWarn("在线TTS音频落盘失败，回退原始地址", zap.String("word", text), zap.String("voice", voice), zap.String("audioURL", audioURL), zap.Error(cacheErr))
		} else if strings.TrimSpace(cachedURL) != "" {
			audioURL = cachedURL
		}
	}

	return audioURL, nil
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
	return word, err
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
			"word":        req.Word,
			"phonetic_us": req.PhoneticUS,
			"phonetic_uk": req.PhoneticUK,
			"audio_us":    req.AudioUS,
			"audio_uk":    req.AudioUK,
			"explanation": req.Explanation,
		}
		if err := tx.Model(&model.EnglishWord{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}

		if req.CategoryIDs == nil && req.ChapterIDs == nil {
			return nil
		}

		if err := tx.Where("word_id = ?", req.ID).Delete(&model.EnglishCategoryWord{}).Error; err != nil {
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

		return nil
	})
}

func (s *EnglishWordService) DeleteWord(id uint) error {
	if id == 0 {
		return errors.New("ID不能为空")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var existing model.EnglishWord
		if err := tx.Select("id").Where("id = ?", id).First(&existing).Error; err != nil {
			return err
		}

		if err := tx.Where("word_id = ?", id).Delete(&model.EnglishCategoryWord{}).Error; err != nil {
			return err
		}
		if err := tx.Where("word_id = ?", id).Delete(&model.EnglishWordSentence{}).Error; err != nil {
			return err
		}
		if err := tx.Where("word_id = ?", id).Delete(&model.EnglishWordErrorLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("target_type = ? AND target_id = ?", 1, id).Delete(&model.UserCollection{}).Error; err != nil {
			return err
		}

		return tx.Delete(&model.EnglishWord{}, "id = ?", id).Error
	})
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
	query := db.Table("english_words w").Joins("LEFT JOIN english_category_words cw ON cw.word_id = w.id")
	if info.CategoryID > 0 {
		query = query.Where("cw.category_id = ?", info.CategoryID)
	}
	if info.ChapterID > 0 {
		query = query.Where("cw.chapter_id = ?", info.ChapterID)
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
