package service

import (
	"regexp"
	"strings"

	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"gorm.io/gorm"
)

type VideoSubtitleService struct{}

// ParseAndHighlightSubtitle 解析字幕，碰对单词库，生成带高亮标签的入库句表
func (s *VideoSubtitleService) ParseAndHighlightSubtitle(episodeID uint, items []request.SubtitleItem) error {
	// 1. 全部提取当前系统的基础英语词库 (缓存供下文匹对使用)
	// 实际生产中可以放入 Redis 缓存或 Trie 树实现更快匹配
	var words []model.EnglishWord
	if err := global.GVA_DB.Select("id", "word").Find(&words).Error; err != nil {
		return err
	}

	wordMap := make(map[string]uint)
	for _, w := range words {
		// 为了忽略大小写，统一下转小写
		wordMap[strings.ToLower(w.Word)] = w.ID
	}

	// 清理该单集原有的老句子库，防止重复解析
	global.GVA_DB.Where("episode_id = ?", episodeID).Delete(&model.VideoSentence{})

	// 2. 高效正则分词器 (拆分出英文单词部分用于匹配，保留原符号)
	// 仅匹配英文字母，不包含标点
	re := regexp.MustCompile(`([a-zA-Z]+)`)

	var sentencesToInsert []model.VideoSentence
	for _, item := range items {
		// ReplaceAllStringFunc 遍历当前句子的每一个单词
		highlightedEnglish := re.ReplaceAllStringFunc(item.English, func(match string) string {
			lowerMatch := strings.ToLower(match)
			if wordID, exists := wordMap[lowerMatch]; exists {
				// 如果该词在我们的词库存在，那就打上特殊标签，比如 <w id="12">match</w>
				return fmt.Sprintf(`<w id="%d">%s</w>`, wordID, match)
			}
			return match
		})

		sentencesToInsert = append(sentencesToInsert, model.VideoSentence{
			EpisodeID: episodeID,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
			English:   highlightedEnglish,
			Translate: item.Translate,
		})
	}

	// 3. 批量高效写入数据库 (Batch Insert)
	if len(sentencesToInsert) > 0 {
		return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			return tx.CreateInBatches(sentencesToInsert, 100).Error
		})
	}

	return nil
}

func (s *VideoSubtitleService) GetSentenceList(episodeID uint) (list []model.VideoSentence, err error) {
	err = global.GVA_DB.Where("episode_id = ?", episodeID).Order("start_time ASC, id ASC").Find(&list).Error
	return
}
