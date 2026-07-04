package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"gorm.io/gorm"
)

type VideoSubtitleService struct{}

type subtitleLine struct {
	StartTime float64
	EndTime   float64
	Text      string
}

var subtitleWordRegexp = regexp.MustCompile(`([a-zA-Z]+)`)
var subtitleTagRegexp = regexp.MustCompile(`<[^>]+>`)
var subtitleBlankSplitRegexp = regexp.MustCompile(`\n\s*\n+`)

// ParseAndHighlightSubtitle 解析字幕，碰对单词库，生成带高亮标签的入库句表
func (s *VideoSubtitleService) ParseAndHighlightSubtitle(episodeID uint, items []request.SubtitleItem) error {
	wordMap, err := s.loadWordMap()
	if err != nil {
		return err
	}

	var sentencesToInsert []model.VideoSentence
	for _, item := range items {
		sentencesToInsert = append(sentencesToInsert, model.VideoSentence{
			EpisodeID: episodeID,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			English:   highlightSubtitleSentence(item.English, wordMap),
			Translate: item.Translate,
		})
	}

	return s.replaceEpisodeSentences(episodeID, sentencesToInsert)
}

func (s *VideoSubtitleService) ParseAndHighlightSubtitleFromFiles(req request.ParseSubtitleFilesReq) error {
	if req.EpisodeID == 0 {
		return errors.New("episodeId不能为空")
	}

	englishURL := strings.TrimSpace(req.EnglishSubtitleURL)
	if englishURL == "" {
		return errors.New("英文字幕文件不能为空")
	}

	englishContent, err := fetchRemoteSubtitleContent(englishURL)
	if err != nil {
		return fmt.Errorf("读取英文字幕失败: %w", err)
	}
	englishLines := parseSubtitleLines(englishContent)
	if len(englishLines) == 0 {
		return errors.New("英文字幕解析失败或内容为空")
	}

	wordMap, err := s.loadWordMap()
	if err != nil {
		return err
	}

	translationByLang := make(map[string][]string)
	subtitleRows := make([]model.VideoSubtitle, 0, 1+len(req.TranslationSubtitle))
	subtitleRows = append(subtitleRows, model.VideoSubtitle{
		EpisodeID:   req.EpisodeID,
		Language:    "en",
		Format:      inferSubtitleFormat(englishURL),
		SubtitleUrl: englishURL,
	})

	for _, item := range req.TranslationSubtitle {
		lang := normalizeSubtitleLanguage(item.Language)
		subtitleURL := strings.TrimSpace(item.SubtitleURL)
		if lang == "" || subtitleURL == "" || lang == "en" {
			continue
		}

		content, readErr := fetchRemoteSubtitleContent(subtitleURL)
		if readErr != nil {
			return fmt.Errorf("读取字幕失败(%s): %w", lang, readErr)
		}
		lines := parseSubtitleLines(content)
		if len(lines) == 0 {
			return fmt.Errorf("字幕解析失败(%s)", lang)
		}

		translationByLang[lang] = alignSubtitleText(englishLines, lines)
		subtitleRows = append(subtitleRows, model.VideoSubtitle{
			EpisodeID:   req.EpisodeID,
			Language:    lang,
			Format:      inferSubtitleFormat(subtitleURL),
			SubtitleUrl: subtitleURL,
		})
	}

	sentences := make([]model.VideoSentence, 0, len(englishLines))
	for idx, line := range englishLines {
		translateMap := make(map[string]string)
		for lang, rows := range translationByLang {
			if idx >= len(rows) {
				continue
			}
			text := strings.TrimSpace(rows[idx])
			if text == "" {
				continue
			}
			translateMap[lang] = text
		}

		translateJSON := "{}"
		if len(translateMap) > 0 {
			if raw, marshalErr := json.Marshal(translateMap); marshalErr == nil {
				translateJSON = string(raw)
			}
		}

		sentences = append(sentences, model.VideoSentence{
			EpisodeID: req.EpisodeID,
			StartTime: line.StartTime,
			EndTime:   line.EndTime,
			English:   highlightSubtitleSentence(line.Text, wordMap),
			Translate: translateJSON,
		})
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Where("episode_id = ?", req.EpisodeID).Delete(&model.VideoSubtitle{}).Error; err != nil {
			return err
		}
		if len(subtitleRows) > 0 {
			if err = tx.CreateInBatches(subtitleRows, 100).Error; err != nil {
				return err
			}
		}

		if err = tx.Where("episode_id = ?", req.EpisodeID).Delete(&model.VideoSentence{}).Error; err != nil {
			return err
		}
		if len(sentences) > 0 {
			if err = tx.CreateInBatches(sentences, 100).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *VideoSubtitleService) replaceEpisodeSentences(episodeID uint, sentences []model.VideoSentence) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("episode_id = ?", episodeID).Delete(&model.VideoSentence{}).Error; err != nil {
			return err
		}
		if len(sentences) == 0 {
			return nil
		}
		return tx.CreateInBatches(sentences, 100).Error
	})
}

func (s *VideoSubtitleService) loadWordMap() (map[string]uint, error) {
	var words []model.EnglishWord
	if err := global.GVA_DB.Select("id", "word").Find(&words).Error; err != nil {
		return nil, err
	}

	wordMap := make(map[string]uint, len(words))
	for _, w := range words {
		wordMap[strings.ToLower(strings.TrimSpace(w.Word))] = w.ID
	}
	return wordMap, nil
}

func highlightSubtitleSentence(text string, wordMap map[string]uint) string {
	return subtitleWordRegexp.ReplaceAllStringFunc(text, func(match string) string {
		lowerMatch := strings.ToLower(match)
		if wordID, exists := wordMap[lowerMatch]; exists {
			return fmt.Sprintf(`<w id="%d">%s</w>`, wordID, match)
		}
		return match
	})
}

func fetchRemoteSubtitleContent(rawURL string) (string, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(strings.TrimSpace(rawURL))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("HTTP状态异常(%d)", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseSubtitleLines(content string) []subtitleLine {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	content = strings.TrimPrefix(content, "\ufeff")
	content = strings.TrimSpace(content)
	if content == "" {
		return []subtitleLine{}
	}

	blocks := subtitleBlankSplitRegexp.Split(content, -1)
	lines := make([]subtitleLine, 0, len(blocks))
	for _, block := range blocks {
		item, ok := parseSubtitleBlock(block)
		if !ok {
			continue
		}
		lines = append(lines, item)
	}

	return lines
}

func parseSubtitleBlock(block string) (subtitleLine, bool) {
	rawLines := strings.Split(strings.TrimSpace(block), "\n")
	if len(rawLines) < 2 {
		return subtitleLine{}, false
	}

	timeLineIndex := -1
	for idx, line := range rawLines {
		if strings.Contains(line, "-->") {
			timeLineIndex = idx
			break
		}
	}
	if timeLineIndex < 0 || timeLineIndex >= len(rawLines)-1 {
		return subtitleLine{}, false
	}

	startTime, endTime, ok := parseSubtitleTimeRange(rawLines[timeLineIndex])
	if !ok {
		return subtitleLine{}, false
	}

	text := cleanSubtitleText(strings.Join(rawLines[timeLineIndex+1:], " "))
	if text == "" {
		return subtitleLine{}, false
	}

	return subtitleLine{StartTime: startTime, EndTime: endTime, Text: text}, true
}

func parseSubtitleTimeRange(raw string) (float64, float64, bool) {
	parts := strings.Split(raw, "-->")
	if len(parts) != 2 {
		return 0, 0, false
	}

	startRaw := strings.TrimSpace(parts[0])
	endRaw := strings.TrimSpace(parts[1])
	if fields := strings.Fields(endRaw); len(fields) > 0 {
		endRaw = fields[0]
	}

	startTime, ok := parseSubtitleTimestamp(startRaw)
	if !ok {
		return 0, 0, false
	}
	endTime, ok := parseSubtitleTimestamp(endRaw)
	if !ok {
		return 0, 0, false
	}
	if endTime < startTime {
		endTime = startTime
	}
	return startTime, endTime, true
}

func parseSubtitleTimestamp(raw string) (float64, bool) {
	text := strings.TrimSpace(strings.ReplaceAll(raw, ",", "."))
	if text == "" {
		return 0, false
	}

	parts := strings.Split(text, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, false
	}

	parseIntPart := func(value string) (float64, bool) {
		n, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return 0, false
		}
		return float64(n), true
	}

	parseFloatPart := func(value string) (float64, bool) {
		n, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return 0, false
		}
		return n, true
	}

	if len(parts) == 2 {
		minutes, ok := parseIntPart(parts[0])
		if !ok {
			return 0, false
		}
		seconds, ok := parseFloatPart(parts[1])
		if !ok {
			return 0, false
		}
		return minutes*60 + seconds, true
	}

	hours, ok := parseIntPart(parts[0])
	if !ok {
		return 0, false
	}
	minutes, ok := parseIntPart(parts[1])
	if !ok {
		return 0, false
	}
	seconds, ok := parseFloatPart(parts[2])
	if !ok {
		return 0, false
	}

	return hours*3600 + minutes*60 + seconds, true
}

func cleanSubtitleText(raw string) string {
	text := strings.TrimSpace(raw)
	if text == "" {
		return ""
	}

	text = subtitleTagRegexp.ReplaceAllString(text, " ")
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "{\\an8}", "")
	text = strings.ReplaceAll(text, "{\\an7}", "")
	text = strings.TrimSpace(strings.Join(strings.Fields(text), " "))
	return text
}

func normalizeSubtitleLanguage(raw string) string {
	code := strings.TrimSpace(raw)
	if code == "" {
		return ""
	}

	parts := strings.Split(code, "-")
	for idx, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if idx == 0 {
			parts[idx] = strings.ToLower(part)
			continue
		}
		if len(part) <= 3 {
			parts[idx] = strings.ToUpper(part)
		} else {
			parts[idx] = strings.ToLower(part)
		}
	}

	return strings.Join(parts, "-")
}

func inferSubtitleFormat(rawURL string) string {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(strings.TrimSpace(rawURL))), ".")
	if ext == "" {
		return "srt"
	}
	return ext
}

func alignSubtitleText(base []subtitleLine, target []subtitleLine) []string {
	result := make([]string, len(base))
	if len(base) == 0 || len(target) == 0 {
		return result
	}

	for idx, baseLine := range base {
		if idx < len(target) {
			if math.Abs(target[idx].StartTime-baseLine.StartTime) <= 1.5 {
				result[idx] = target[idx].Text
				continue
			}
		}

		bestIndex := -1
		bestDelta := math.MaxFloat64
		for targetIndex, targetLine := range target {
			delta := math.Abs(targetLine.StartTime - baseLine.StartTime)
			if delta < bestDelta {
				bestDelta = delta
				bestIndex = targetIndex
			}
		}

		if bestIndex >= 0 && bestDelta <= 1.5 {
			result[idx] = target[bestIndex].Text
			continue
		}

		if idx < len(target) {
			result[idx] = target[idx].Text
		}
	}

	return result
}

func (s *VideoSubtitleService) GetSentenceList(episodeID uint) (list []model.VideoSentence, err error) {
	err = global.GVA_DB.Where("episode_id = ?", episodeID).Order("start_time ASC, id ASC").Find(&list).Error
	return
}
