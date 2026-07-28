package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"gorm.io/gorm"
)

type DiaryService struct{}

// DiaryCategory
func (s *DiaryService) CreateDiaryCategory(category *model.DiaryCategory) error {
	return global.GVA_DB.Create(category).Error
}

func (s *DiaryService) UpdateDiaryCategory(category model.DiaryCategory) error {
	return global.GVA_DB.Model(&model.DiaryCategory{}).Where("id = ?", category.ID).Updates(&category).Error
}

func (s *DiaryService) DeleteDiaryCategory(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var diaryIDs []uint
		if err := tx.Model(&model.Diary{}).Where("category_id = ?", id).Pluck("id", &diaryIDs).Error; err != nil {
			return err
		}
		if len(diaryIDs) > 0 {
			if err := tx.Delete(&model.DiarySentence{}, "diary_id IN ?", diaryIDs).Error; err != nil {
				return err
			}
			if err := tx.Where("diary_id IN ?", diaryIDs).Delete(&clientModel.DiaryTagRelation{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.Diary{}, "category_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.DiaryCategory{}, "id = ?", id).Error
	})
}

func (s *DiaryService) GetDiaryCategory(id uint) (model.DiaryCategory, error) {
	var category model.DiaryCategory
	err := global.GVA_DB.Where("id = ?", id).First(&category).Error
	return category, err
}

func (s *DiaryService) GetDiaryCategoryList(info request.DiaryCategorySearch) (list []model.DiaryCategory, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Model(&model.DiaryCategory{})
	if info.ShowHome != nil {
		db = db.Where("show_home = ?", *info.ShowHome)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("sort ASC, id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error
	return
}

// Diary
func (s *DiaryService) CreateDiaryWithTags(diary *model.Diary, tagIds []uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(diary).Error; err != nil {
			return err
		}
		return syncDiaryTags(tx, diary.ID, tagIds)
	})
}

func (s *DiaryService) UpdateDiaryWithTags(diary model.Diary, tagIds []uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Diary{}).Where("id = ?", diary.ID).Updates(&diary).Error; err != nil {
			return err
		}
		return syncDiaryTags(tx, diary.ID, tagIds)
	})
}

func syncDiaryTags(tx *gorm.DB, diaryID uint, tagIds []uint) error {
	if err := tx.Where("diary_id = ?", diaryID).Delete(&clientModel.DiaryTagRelation{}).Error; err != nil {
		return err
	}
	if len(tagIds) == 0 {
		return nil
	}
	tags := make([]clientModel.DiaryTagRelation, 0, len(tagIds))
	for _, tagID := range tagIds {
		if tagID > 0 {
			tags = append(tags, clientModel.DiaryTagRelation{DiaryID: diaryID, TagID: tagID})
		}
	}
	if len(tags) > 0 {
		return tx.Create(&tags).Error
	}
	return nil
}

func (s *DiaryService) GetDiaryTags(diaryID uint) ([]clientModel.DiaryTag, error) {
	var tags []clientModel.DiaryTag
	err := global.GVA_DB.Table("diary_tags").
		Joins("JOIN diary_tag_relations ON diary_tag_relations.tag_id = diary_tags.id").
		Where("diary_tag_relations.diary_id = ?", diaryID).
		Find(&tags).Error
	return tags, err
}

func (s *DiaryService) GetDiaryTagsBatch(diaryIDs []uint) (map[uint][]clientModel.DiaryTag, error) {
	if len(diaryIDs) == 0 {
		return map[uint][]clientModel.DiaryTag{}, nil
	}
	type tagRow struct {
		DiaryID uint `gorm:"column:diary_id"`
		clientModel.DiaryTag
	}
	var rows []tagRow
	err := global.GVA_DB.Table("diary_tags").
		Select("diary_tags.*, diary_tag_relations.diary_id").
		Joins("JOIN diary_tag_relations ON diary_tag_relations.tag_id = diary_tags.id").
		Where("diary_tag_relations.diary_id IN ?", diaryIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[uint][]clientModel.DiaryTag, len(rows))
	for _, row := range rows {
		result[row.DiaryID] = append(result[row.DiaryID], row.DiaryTag)
	}
	return result, nil
}

func (s *DiaryService) DeleteDiary(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.DiarySentence{}, "diary_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Where("diary_id = ?", id).Delete(&clientModel.DiaryTagRelation{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Diary{}, "id = ?", id).Error
	})
}

func (s *DiaryService) GetDiary(id uint) (model.Diary, error) {
	var diary model.Diary
	err := global.GVA_DB.Where("id = ?", id).First(&diary).Error
	return diary, err
}

func (s *DiaryService) GetDiaryList(info request.DiarySearch) (list []model.Diary, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Model(&model.Diary{})
	if info.CategoryID > 0 {
		db = db.Where("category_id = ?", info.CategoryID)
	}
	keyword := strings.TrimSpace(info.Keyword)
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("sort ASC, id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error
	return
}

// GetDiaryListByTag 按分类+标签筛选日记
func (s *DiaryService) GetDiaryListByTag(info request.DiaryTagFilterSearch) (list []model.Diary, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)

	diaryTable := resolveContentTableName("diary", "diaries")
	tagTable := "diary_tag_relations"

	base := global.GVA_DB.Table(diaryTable + " AS d").
		Joins("JOIN " + tagTable + " AS dtr ON dtr.diary_id = d.id")

	if info.CategoryID > 0 {
		base = base.Where("d.category_id = ?", info.CategoryID)
	}
	if len(info.TagIds) > 0 {
		base = base.Where("dtr.tag_id IN ?", info.TagIds)
	}

	countDB := base.Session(&gorm.Session{})
	err = countDB.Select("COUNT(DISTINCT d.id)").Scan(&total).Error
	if err != nil {
		return
	}

	orderClause := "d.sort ASC, d.id DESC"
	dataDB := base.Session(&gorm.Session{})
	err = dataDB.Select("DISTINCT d.*").
		Order(orderClause).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&list).Error
	return
}

// DiarySentence
func (s *DiaryService) GetDiarySentenceList(diaryID uint) (list []model.DiarySentence, err error) {
	err = global.GVA_DB.Where("diary_id = ?", diaryID).Order("start_time ASC, id ASC").Find(&list).Error
	return
}

// GetDiarySentenceListWithAuth 获取日记字幕句子列表（含试看过滤）
// hasFullAuth 为 true 时返回全部内容；为 false 时超出试看比例的句子内容清空并标记 locked
func (s *DiaryService) GetDiarySentenceListWithAuth(diaryID uint, hasFullAuth bool, trialPercent int) (list []model.DiarySentence, err error) {
	err = global.GVA_DB.Where("diary_id = ?", diaryID).Order("start_time ASC, id ASC").Find(&list).Error
	if err != nil {
		return
	}
	if hasFullAuth {
		return
	}
	// 计算最大试看时间
	var diary model.Diary
	if dbErr := global.GVA_DB.Where("id = ?", diaryID).Select("duration").First(&diary).Error; dbErr != nil {
		return list, nil // 查不到日记信息时不做过滤
	}
	duration := diary.Duration
	if duration <= 0 {
		// 数据库未设置时长时，从字幕数据中推算（取最大 endTime），与前端实际音频时长对齐
		for i := range list {
			if list[i].EndTime > duration {
				duration = list[i].EndTime
			}
		}
	}
	if duration <= 0 {
		return
	}
	safePercent := trialPercent
	if safePercent < 1 {
		safePercent = 8
	}
	maxTrialTime := duration * float64(safePercent) / 100.0
	for i := range list {
		// startTime >= maxTrialTime 的句子完全锁定（内容清空，时间区间保留）
		// startTime < maxTrialTime 的句子即使 endTime 超出试看范围也显示（跨边界情况）
		if list[i].StartTime >= maxTrialTime {
			list[i].Locked = true
			list[i].English = ""
			list[i].Translate = ""
		}
	}
	return
}

// ParseAndSaveDiarySentences 解析日记字幕（JSON手动输入，支持关键词高亮）
func (s *DiaryService) ParseAndSaveDiarySentences(diaryID uint, items []request.DiarySubtitleItem, keywordIds []uint) error {
	wordMap, err := loadDiaryWordMap()
	if err != nil {
		return err
	}

	// 如果有选中关键词，只高亮这些
	if len(keywordIds) > 0 {
		keywordIDSet := make(map[uint]bool, len(keywordIds))
		for _, id := range keywordIds {
			keywordIDSet[id] = true
		}
		filteredWordMap := make(map[string]uint)
		for word, id := range wordMap {
			if keywordIDSet[id] {
				filteredWordMap[word] = id
			}
		}
		wordMap = filteredWordMap
	}

	sentences := make([]model.DiarySentence, 0, len(items))
	for _, item := range items {
		sentences = append(sentences, model.DiarySentence{
			DiaryID:   diaryID,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			English:   highlightDiarySentence(item.English, wordMap),
			Translate: item.Translate,
		})
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("diary_id = ?", diaryID).Delete(&model.DiarySentence{}).Error; err != nil {
			return err
		}
		if len(sentences) > 0 {
			return tx.CreateInBatches(sentences, 100).Error
		}
		return nil
	})
}

// ScanDiaryKeywords 扫描日记字幕文件提取关键词
func (s *DiaryService) ScanDiaryKeywords(req request.ScanDiaryKeywordsReq) ([]request.DiaryKeywordItem, error) {
	englishURL := strings.TrimSpace(req.EnglishSubtitleURL)
	if englishURL == "" {
		return nil, errors.New("英文字幕文件不能为空")
	}

	content, err := fetchRemoteSubtitleContent(englishURL)
	if err != nil {
		return nil, fmt.Errorf("读取英文字幕失败: %w", err)
	}
	lines := parseSubtitleLines(content)
	if len(lines) == 0 {
		return nil, errors.New("英文字幕解析失败或内容为空")
	}

	wordMap, err := loadDiaryWordMap()
	if err != nil {
		return nil, err
	}

	// 统计字幕文件中每个单词的出现次数
	wordCounter := make(map[string]int)
	for _, line := range lines {
		for _, match := range subtitleWordRegexp.FindAllString(line.Text, -1) {
			wordCounter[strings.ToLower(match)]++
		}
	}

	// 构建结果：字幕中出现的所有单词（含未入库的）
	type wordEntry struct {
		word  string
		count int
	}
	entries := make([]wordEntry, 0, len(wordCounter))
	for w, c := range wordCounter {
		entries = append(entries, wordEntry{word: w, count: c})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].count > entries[j].count })

	result := make([]request.DiaryKeywordItem, 0, len(entries))
	for _, e := range entries {
		id, matched := wordMap[e.word]
		result = append(result, request.DiaryKeywordItem{
			Word:    e.word,
			WordID:  id,
			Matched: matched,
			Count:   e.count,
		})
	}
	return result, nil
}

// ParseAndSaveDiarySentencesFromFiles 从字幕文件解析并入库日记字幕
func (s *DiaryService) ParseAndSaveDiarySentencesFromFiles(req request.ParseDiarySubtitleFilesReq) error {
	if req.DiaryID == 0 {
		return errors.New("diaryId不能为空")
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

	// 解析英式英文字幕文件（仅取时间轴，文字内容以美式为准）
	var ukLines []subtitleLine
	ukSubtitleURL := strings.TrimSpace(req.EnglishSubtitleURLUk)
	if ukSubtitleURL != "" {
		ukContent, ukErr := fetchRemoteSubtitleContent(ukSubtitleURL)
		if ukErr != nil {
			return fmt.Errorf("读取英式英文字幕失败: %w", ukErr)
		}
		ukLines = parseSubtitleLines(ukContent)
		if len(ukLines) == 0 {
			return errors.New("英式英文字幕解析失败或内容为空")
		}
	}

	// 构建关键词ID集合
	keywordIDSet := make(map[uint]bool, len(req.KeywordIDs))
	for _, id := range req.KeywordIDs {
		keywordIDSet[id] = true
	}

	wordMap, err := loadDiaryWordMap()
	if err != nil {
		return err
	}

	// 筛选后的高亮映射
	filteredWordMap := make(map[string]uint)
	for word, id := range wordMap {
		if keywordIDSet[id] {
			filteredWordMap[word] = id
		}
	}

	translationByLang := make(map[string][]string)
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
	}

	sentences := make([]model.DiarySentence, 0, len(englishLines))
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

		sentence := model.DiarySentence{
			DiaryID:   req.DiaryID,
			StartTime: line.StartTime,
			EndTime:   line.EndTime,
			English:   highlightDiarySentence(line.Text, filteredWordMap),
			Translate: translateJSON,
		}
		// 若有英式字幕，按索引对齐填充 UK 时间轴
		if len(ukLines) > 0 && idx < len(ukLines) {
			sentence.StartTimeUk = ukLines[idx].StartTime
			sentence.EndTimeUk = ukLines[idx].EndTime
		}
		sentences = append(sentences, sentence)
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 保存字幕文件 URL 到日记记录
		urlUpdates := map[string]interface{}{
			"english_subtitle_url": englishURL,
		}
		if ukSubtitleURL != "" {
			urlUpdates["english_subtitle_url_uk"] = ukSubtitleURL
		}
		// 保存翻译字幕文件 URL 映射
		if len(req.TranslationSubtitle) > 0 {
			translationURLMap := make(map[string]string)
			for _, item := range req.TranslationSubtitle {
				lang := normalizeSubtitleLanguage(item.Language)
				url := strings.TrimSpace(item.SubtitleURL)
				if lang == "" || url == "" || lang == "en" {
					continue
				}
				translationURLMap[lang] = url
			}
			if len(translationURLMap) > 0 {
				if raw, marshalErr := json.Marshal(translationURLMap); marshalErr == nil {
					urlUpdates["translation_subtitle_urls"] = string(raw)
				}
			}
		}
		if err = tx.Model(&model.Diary{}).Where("id = ?", req.DiaryID).Updates(urlUpdates).Error; err != nil {
			return err
		}

		if err = tx.Where("diary_id = ?", req.DiaryID).Delete(&model.DiarySentence{}).Error; err != nil {
			return err
		}
		if len(sentences) > 0 {
			return tx.CreateInBatches(sentences, 100).Error
		}
		return nil
	})
}

// loadDiaryWordMap 加载单词库映射
func loadDiaryWordMap() (map[string]uint, error) {
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

// highlightDiarySentence 对英文字幕句子进行关键词高亮
func highlightDiarySentence(text string, wordMap map[string]uint) string {
	return subtitleWordRegexp.ReplaceAllStringFunc(text, func(match string) string {
		lowerMatch := strings.ToLower(match)
		if wordID, exists := wordMap[lowerMatch]; exists {
			return fmt.Sprintf(`<w id="%d">%s</w>`, wordID, match)
		}
		return match
	})
}

// GetDiaryKeywords 获取日记字幕中出现的所有单词列表
// 遍历字幕句子的所有单词，匹配词库的标记 matched=true 并附 wordId，
// 已高亮（有 <w id> 标签）的单词 count>0，未高亮但字幕中出现的 count=0
func (s *DiaryService) GetDiaryKeywords(diaryID uint) ([]request.DiaryKeywordItem, error) {
	var sentences []model.DiarySentence
	if err := global.GVA_DB.Where("diary_id = ?", diaryID).Select("english").Find(&sentences).Error; err != nil {
		return nil, err
	}

	// 加载词库映射
	wordMap, err := loadDiaryWordMap()
	if err != nil {
		return nil, err
	}

	// 统计字幕中所有单词（去标签后的纯文本）
	wordCounter := make(map[string]int)     // word -> 出现次数
	highlightedIDs := make(map[string]uint) // 已高亮的 word -> wordID
	wTagRegexp := regexp.MustCompile(`<w id="(\d+)">([^<]+)</w>`)
	cleanTagRegexp := regexp.MustCompile(`<[^>]+>`)

	for _, sentence := range sentences {
		// 先提取已高亮的单词（有 <w id> 标签）
		matches := wTagRegexp.FindAllStringSubmatch(sentence.English, -1)
		for _, m := range matches {
			if len(m) >= 3 {
				word := strings.ToLower(strings.TrimSpace(m[2]))
				if word == "" {
					continue
				}
				if parsedID, err := strconv.ParseUint(m[1], 10, 64); err == nil {
					highlightedIDs[word] = uint(parsedID)
				}
			}
		}
		// 去除所有标签后统计所有单词出现次数
		cleanText := cleanTagRegexp.ReplaceAllString(sentence.English, " ")
		for _, match := range subtitleWordRegexp.FindAllString(cleanText, -1) {
			wordCounter[strings.ToLower(match)]++
		}
	}

	// 构建结果：字幕中出现的所有单词
	result := make([]request.DiaryKeywordItem, 0, len(wordCounter))
	for word, count := range wordCounter {
		wordID, inLibrary := wordMap[word]
		matched := inLibrary
		// 优先使用已高亮的 wordID（即使词库已删，仍保留高亮标签中的 ID）
		if hlID, hl := highlightedIDs[word]; hl {
			wordID = hlID
		}
		result = append(result, request.DiaryKeywordItem{
			Word:    word,
			WordID:  wordID,
			Matched: matched,
			Count:   count,
		})
	}

	// 已匹配的排前面，然后按出现次数降序
	sort.Slice(result, func(i, j int) bool {
		if result[i].Matched != result[j].Matched {
			return result[i].Matched
		}
		return result[i].Count > result[j].Count
	})

	return result, nil
}

// RehighlightDiarySentences 对已有日记字幕句子去标签后重新高亮
func (s *DiaryService) RehighlightDiarySentences(req request.RehighlightDiarySentencesReq) error {
	if req.DiaryID == 0 {
		return errors.New("diaryId不能为空")
	}

	keywordIDSet := make(map[uint]bool, len(req.KeywordIDs))
	for _, id := range req.KeywordIDs {
		keywordIDSet[id] = true
	}

	wordMap, err := loadDiaryWordMap()
	if err != nil {
		return err
	}

	filteredWordMap := make(map[string]uint)
	for word, id := range wordMap {
		if keywordIDSet[id] {
			filteredWordMap[word] = id
		}
	}

	var sentences []model.DiarySentence
	if err := global.GVA_DB.Where("diary_id = ?", req.DiaryID).Find(&sentences).Error; err != nil {
		return err
	}

	stripWTagRegexp := regexp.MustCompile(`<w[^>]*>|</w>`)
	for i := range sentences {
		cleanText := stripWTagRegexp.ReplaceAllString(sentences[i].English, "")
		sentences[i].English = highlightDiarySentence(cleanText, filteredWordMap)
	}

	if len(sentences) == 0 {
		return nil
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, s := range sentences {
			if err := tx.Model(&model.DiarySentence{}).Where("id = ?", s.ID).Update("english", s.English).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateDiarySentenceList 批量更新日记字幕句子
func (s *DiaryService) UpdateDiarySentenceList(req request.UpdateDiarySentenceListReq) error {
	if req.DiaryID == 0 {
		return errors.New("diaryId不能为空")
	}
	if len(req.Sentences) == 0 {
		return errors.New("sentences不能为空")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Sentences {
			if item.ID == 0 {
				continue
			}

			updates := map[string]interface{}{
				"english":    strings.TrimSpace(item.English),
				"translate":  strings.TrimSpace(item.Translate),
				"start_time": item.StartTime,
				"end_time":   item.EndTime,
			}

			res := tx.Model(&model.DiarySentence{}).
				Where("id = ? AND diary_id = ?", item.ID, req.DiaryID).
				Updates(updates)
			if res.Error != nil {
				return res.Error
			}
		}
		return nil
	})
}
