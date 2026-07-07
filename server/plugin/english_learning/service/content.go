package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"gorm.io/gorm"
)

type ContentService struct{}

var storageKeyNormalizeRegexp = regexp.MustCompile(`[^a-z0-9_-]+`)

func normalizeStorageKey(raw string) string {
	key := strings.ToLower(strings.TrimSpace(raw))
	if key == "" {
		return ""
	}
	key = storageKeyNormalizeRegexp.ReplaceAllString(key, "-")
	key = strings.Trim(key, "-")
	return key
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return page, pageSize
}

func resolveContentTableName(baseName string, fallback string) string {
	if global.GVA_DB != nil && global.GVA_DB.NamingStrategy != nil {
		resolved := strings.TrimSpace(global.GVA_DB.NamingStrategy.TableName(baseName))
		if resolved != "" {
			return resolved
		}
	}

	fallback = strings.TrimSpace(fallback)
	if fallback != "" {
		return fallback
	}

	return baseName
}

func resolveVideoSeriesOrderExpr(seriesAlias string) string {
	alias := strings.TrimSpace(seriesAlias)
	if alias == "" {
		alias = "video_series"
	}

	if global.GVA_DB != nil && global.GVA_DB.Migrator().HasColumn(&model.VideoSeries{}, "sort") {
		return fmt.Sprintf("%s.sort ASC, %s.id DESC", alias, alias)
	}

	return fmt.Sprintf("%s.id DESC", alias)
}

// EnglishCategory
func (s *ContentService) CreateCategory(category *model.EnglishCategory) error {
	return global.GVA_DB.Create(category).Error
}

func (s *ContentService) UpdateCategory(category model.EnglishCategory) error {
	return global.GVA_DB.Model(&model.EnglishCategory{}).Where("id = ?", category.ID).Updates(&category).Error
}

func (s *ContentService) DeleteCategory(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.EnglishCategory{}, "id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.EnglishChapter{}, "category_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.EnglishCategoryWord{}, "category_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *ContentService) GetCategory(id uint) (model.EnglishCategory, error) {
	var category model.EnglishCategory
	err := global.GVA_DB.Where("id = ?", id).First(&category).Error
	return category, err
}

func (s *ContentService) GetCategoryList(info request.EnglishCategorySearch) (list []model.EnglishCategory, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Model(&model.EnglishCategory{})
	if info.NeedVip != nil {
		db = db.Where("need_vip = ?", *info.NeedVip)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("sort ASC, id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error
	return
}

// EnglishChapter
func (s *ContentService) CreateChapter(chapter *model.EnglishChapter) error {
	return global.GVA_DB.Create(chapter).Error
}

func (s *ContentService) UpdateChapter(chapter model.EnglishChapter) error {
	return global.GVA_DB.Model(&model.EnglishChapter{}).Where("id = ?", chapter.ID).Updates(&chapter).Error
}

func (s *ContentService) DeleteChapter(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.EnglishChapter{}, "id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.EnglishCategoryWord{}, "chapter_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *ContentService) GetChapter(id uint) (model.EnglishChapter, error) {
	var chapter model.EnglishChapter
	err := global.GVA_DB.Where("id = ?", id).First(&chapter).Error
	return chapter, err
}

func (s *ContentService) GetChapterList(info request.EnglishChapterSearch) (list []model.EnglishChapter, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Model(&model.EnglishChapter{})
	if info.CategoryID > 0 {
		db = db.Where("category_id = ?", info.CategoryID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("sort ASC, id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error
	return
}

// VideoCategory
func (s *ContentService) CreateVideoCategory(category *model.VideoCategory) error {
	category.StorageKey = normalizeStorageKey(category.StorageKey)
	return global.GVA_DB.Create(category).Error
}

func (s *ContentService) UpdateVideoCategory(category model.VideoCategory) error {
	category.StorageKey = normalizeStorageKey(category.StorageKey)
	return global.GVA_DB.Model(&model.VideoCategory{}).Where("id = ?", category.ID).Updates(&category).Error
}

func (s *ContentService) DeleteVideoCategory(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var seriesIDs []uint
		if err := tx.Model(&model.VideoSeries{}).Where("category_id = ?", id).Pluck("id", &seriesIDs).Error; err != nil {
			return err
		}
		if len(seriesIDs) > 0 {
			var episodeIDs []uint
			if err := tx.Model(&model.VideoEpisode{}).Where("series_id IN ?", seriesIDs).Pluck("id", &episodeIDs).Error; err != nil {
				return err
			}
			if len(episodeIDs) > 0 {
				if err := tx.Delete(&model.VideoSubtitle{}, "episode_id IN ?", episodeIDs).Error; err != nil {
					return err
				}
				if err := tx.Delete(&model.VideoSentence{}, "episode_id IN ?", episodeIDs).Error; err != nil {
					return err
				}
			}
			if err := tx.Delete(&model.VideoEpisode{}, "series_id IN ?", seriesIDs).Error; err != nil {
				return err
			}
			if err := tx.Delete(&model.VideoSeries{}, "id IN ?", seriesIDs).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&model.VideoCategory{}, "id = ?", id).Error
	})
}

func (s *ContentService) GetVideoCategory(id uint) (model.VideoCategory, error) {
	var category model.VideoCategory
	err := global.GVA_DB.Where("id = ?", id).First(&category).Error
	return category, err
}

func (s *ContentService) GetVideoCategoryList(info request.VideoCategorySearch) (list []model.VideoCategory, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Model(&model.VideoCategory{})
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

// VideoSeries
func (s *ContentService) CreateVideoSeries(series *model.VideoSeries) error {
	return global.GVA_DB.Create(series).Error
}

func (s *ContentService) UpdateVideoSeries(series model.VideoSeries) error {
	return global.GVA_DB.Model(&model.VideoSeries{}).Where("id = ?", series.ID).Updates(&series).Error
}

func (s *ContentService) DeleteVideoSeries(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var episodeIDs []uint
		if err := tx.Model(&model.VideoEpisode{}).Where("series_id = ?", id).Pluck("id", &episodeIDs).Error; err != nil {
			return err
		}
		if len(episodeIDs) > 0 {
			if err := tx.Delete(&model.VideoSubtitle{}, "episode_id IN ?", episodeIDs).Error; err != nil {
				return err
			}
			if err := tx.Delete(&model.VideoSentence{}, "episode_id IN ?", episodeIDs).Error; err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.VideoEpisode{}, "series_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.VideoSeries{}, "id = ?", id).Error
	})
}

func (s *ContentService) GetVideoSeries(id uint) (model.VideoSeries, error) {
	var series model.VideoSeries
	seriesTable := resolveContentTableName("video_series", "video_series")
	seriesAlias := "vs"
	videoEpisodeTable := resolveContentTableName("video_episode", "video_episodes")
	userWatchHistoryTable := resolveContentTableName("user_watch_history", "user_watch_histories")
	viewCountExpr := fmt.Sprintf("COALESCE(COUNT(DISTINCT %s.id), 0) as view_count", userWatchHistoryTable)
	userCountExpr := fmt.Sprintf("COALESCE(COUNT(DISTINCT %s.user_id), 0) as user_count", userWatchHistoryTable)
	episodeJoin := fmt.Sprintf("LEFT JOIN %s ON %s.series_id = %s.id", videoEpisodeTable, videoEpisodeTable, seriesAlias)
	historyJoin := fmt.Sprintf("LEFT JOIN %s ON %s.episode_id = %s.id", userWatchHistoryTable, userWatchHistoryTable, videoEpisodeTable)

	err := global.GVA_DB.
		Model(&model.VideoSeries{}).
		Table(seriesTable+" AS "+seriesAlias).
		Select(seriesAlias+".*, "+viewCountExpr+", "+userCountExpr).
		Joins(episodeJoin).
		Joins(historyJoin).
		Where(seriesAlias+".id = ?", id).
		Group(seriesAlias + ".id").
		First(&series).Error
	return series, err
}

func (s *ContentService) GetVideoSeriesList(info request.VideoSeriesSearch) (list []model.VideoSeries, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	seriesTable := resolveContentTableName("video_series", "video_series")
	seriesAlias := "vs"
	videoEpisodeTable := resolveContentTableName("video_episode", "video_episodes")
	userWatchHistoryTable := resolveContentTableName("user_watch_history", "user_watch_histories")
	viewCountExpr := fmt.Sprintf("COALESCE(COUNT(DISTINCT %s.id), 0) as view_count", userWatchHistoryTable)
	userCountExpr := fmt.Sprintf("COALESCE(COUNT(DISTINCT %s.user_id), 0) as user_count", userWatchHistoryTable)
	episodeJoin := fmt.Sprintf("LEFT JOIN %s ON %s.series_id = %s.id", videoEpisodeTable, videoEpisodeTable, seriesAlias)
	historyJoin := fmt.Sprintf("LEFT JOIN %s ON %s.episode_id = %s.id", userWatchHistoryTable, userWatchHistoryTable, videoEpisodeTable)
	orderExpr := resolveVideoSeriesOrderExpr(seriesAlias)

	// 构建基础查询
	baseQuery := global.GVA_DB.Model(&model.VideoSeries{}).Table(seriesTable + " AS " + seriesAlias)
	if info.CategoryID > 0 {
		baseQuery = baseQuery.Where(seriesAlias+".category_id = ?", info.CategoryID)
	}
	if info.ShowHome != nil {
		baseQuery = baseQuery.Where(seriesAlias+".show_home = ?", *info.ShowHome)
	}
	keyword := strings.TrimSpace(info.Keyword)
	if keyword != "" {
		baseQuery = baseQuery.Where(seriesAlias+".name LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	err = baseQuery.Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询，并关联统计数据
	err = baseQuery.
		Select(seriesAlias + ".*, " + viewCountExpr + ", " + userCountExpr).
		Joins(episodeJoin).
		Joins(historyJoin).
		Group(seriesAlias + ".id").
		Order(orderExpr).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&list).Error

	return
}

// VideoEpisode
func (s *ContentService) CreateVideoEpisode(episode *model.VideoEpisode) error {
	return global.GVA_DB.Create(episode).Error
}

func (s *ContentService) UpdateVideoEpisode(episode model.VideoEpisode) error {
	return global.GVA_DB.Model(&model.VideoEpisode{}).Where("id = ?", episode.ID).Updates(&episode).Error
}

func (s *ContentService) DeleteVideoEpisode(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.VideoSubtitle{}, "episode_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.VideoSentence{}, "episode_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.VideoEpisode{}, "id = ?", id).Error
	})
}

func (s *ContentService) GetVideoEpisode(id uint) (model.VideoEpisode, error) {
	var episode model.VideoEpisode
	err := global.GVA_DB.Where("id = ?", id).First(&episode).Error
	return episode, err
}

func (s *ContentService) GetVideoEpisodeList(info request.VideoEpisodeSearch) (list []model.VideoEpisode, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Model(&model.VideoEpisode{})
	if info.SeriesID > 0 {
		db = db.Where("series_id = ?", info.SeriesID)
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
