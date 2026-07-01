package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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

// CreateWord 创建单词库实体，集成自动TTS生成以及分类映射
func (s *EnglishWordService) CreateWord(word model.EnglishWord, chapterIDs []uint) error {
	// 1. 唯一性校验
	var count int64
	global.GVA_DB.Model(&model.EnglishWord{}).Where("word = ?", word.Word).Count(&count)
	if count > 0 {
		return errors.New("单词已存在不可重复添加")
	}

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
		if err := tx.Create(&word).Error; err != nil {
			return err
		}

		for _, cid := range chapterIDs {
			assoc := model.EnglishCategoryWord{
				ChapterID: cid,
				WordID:    word.ID,
			}
			// 自动查找归属大类
			var chapter model.EnglishChapter
			if err := tx.Where("id = ?", cid).First(&chapter).Error; err == nil {
				assoc.CategoryID = chapter.CategoryID
			}

			if err := tx.Create(&assoc).Error; err != nil {
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

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var existing model.EnglishWord
		if err := tx.Select("id").Where("id = ?", req.ID).First(&existing).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"word":        req.Word,
			"phonetic_us": req.PhoneticUS,
			"phonetic_uk": req.PhoneticUK,
			"explanation": req.Explanation,
		}
		if err := tx.Model(&model.EnglishWord{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}

		if req.ChapterIDs == nil {
			return nil
		}

		if err := tx.Where("word_id = ?", req.ID).Delete(&model.EnglishCategoryWord{}).Error; err != nil {
			return err
		}

		chapterIDs := uniqueUintIDs(req.ChapterIDs)
		for _, chapterID := range chapterIDs {
			assoc := model.EnglishCategoryWord{
				ChapterID: chapterID,
				WordID:    req.ID,
			}

			var chapter model.EnglishChapter
			if err := tx.Where("id = ?", chapterID).First(&chapter).Error; err == nil {
				assoc.CategoryID = chapter.CategoryID
			}

			if err := tx.Create(&assoc).Error; err != nil {
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
		audioUS, err := s.requestTTSAudioURL(word, result.VoiceUS)
		if err != nil {
			return result, fmt.Errorf("美式发音预检失败: %w", err)
		}
		result.AudioUS = audioUS
	}

	if checkUK {
		result.VoiceUK = getLearningStringSetting("learning_tts_voice_uk", "LEARNING_TTS_VOICE_UK", defaultLearningTTSVoiceUK)
		audioUK, err := s.requestTTSAudioURL(word, result.VoiceUK)
		if err != nil {
			return result, fmt.Errorf("英式发音预检失败: %w", err)
		}
		result.AudioUK = audioUK
	}

	result.DurationMS = time.Since(start).Milliseconds()
	return result, nil
}

func buildWordListBaseQuery(db *gorm.DB, info request.WordListSearch) *gorm.DB {
	query := db.Table("english_words w").Joins("JOIN english_category_words cw ON cw.word_id = w.id")
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
		Order("MIN(cw.sort) ASC, MAX(cw.id) DESC")
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
