package service

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
)

type UserDataService struct{}

type SeriesWatchProgressItem struct {
	EpisodeID    uint      `json:"episodeId"`
	ProgressSecs float64   `json:"progressSecs"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type WatchHistoryDetailItem struct {
	EpisodeID    uint      `json:"episodeId"`
	EpisodeName  string    `json:"episodeName"`
	SeriesID     uint      `json:"seriesId"`
	SeriesName   string    `json:"seriesName"`
	CoverURL     string    `json:"coverUrl"`
	ProgressSecs float64   `json:"progressSecs"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CollectionVideoDetail struct {
	EpisodeID   uint   `json:"episodeId"`
	EpisodeName string `json:"episodeName"`
	SeriesID    uint   `json:"seriesId"`
	SeriesName  string `json:"seriesName"`
	CoverURL    string `json:"coverUrl"`
	VideoURL    string `json:"videoUrl"`
}

type CollectionDetailItem struct {
	model.UserCollection
	Word     *model.EnglishWord     `json:"word,omitempty"`
	Sentence *model.VideoSentence   `json:"sentence,omitempty"`
	Video    *CollectionVideoDetail `json:"video,omitempty"`
}

func (s *UserDataService) SaveWordProgress(userID uint, req request.SaveWordProgressReq) error {
	var history model.UserWordHistory
	return global.GVA_DB.Where("user_id = ?", userID).
		Assign(model.UserWordHistory{
			CategoryID: req.CategoryID,
			ChapterID:  req.ChapterID,
			WordIndex:  req.WordIndex,
		}).FirstOrCreate(&history).Error
}

func (s *UserDataService) GetWordProgress(userID uint) (model.UserWordHistory, error) {
	var history model.UserWordHistory
	result := global.GVA_DB.Where("user_id = ?", userID).Limit(1).Find(&history)
	if result.Error != nil {
		return model.UserWordHistory{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.UserWordHistory{}, nil
	}
	return history, nil
}

func (s *UserDataService) GetWatchProgress(userID uint, episodeID uint) (model.UserWatchHistory, error) {
	var history model.UserWatchHistory
	result := global.GVA_DB.Where("user_id = ? AND episode_id = ?", userID, episodeID).Limit(1).Find(&history)
	if result.Error != nil {
		return model.UserWatchHistory{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.UserWatchHistory{}, nil
	}
	return history, nil
}

func (s *UserDataService) GetSeriesWatchProgressList(userID uint, seriesID uint) (list []SeriesWatchProgressItem, err error) {
	err = global.GVA_DB.Table("user_watch_histories AS h").
		Select("h.episode_id, h.progress_secs, h.updated_at").
		Joins("JOIN video_episodes e ON e.id = h.episode_id").
		Where("h.user_id = ? AND e.series_id = ?", userID, seriesID).
		Order("h.updated_at DESC, h.id DESC").
		Scan(&list).Error
	return
}

func (s *UserDataService) GetWatchHistoryList(userID uint, info request.WatchHistorySearch) (list []WatchHistoryDetailItem, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Table("user_watch_histories AS h").
		Joins("JOIN video_episodes e ON e.id = h.episode_id").
		Joins("LEFT JOIN video_series s ON s.id = e.series_id").
		Where("h.user_id = ?", userID)

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Select("h.episode_id, e.name AS episode_name, e.series_id, s.name AS series_name, s.cover_url, h.progress_secs, h.updated_at").
		Order("h.updated_at DESC, h.id DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list).Error
	return
}

func (s *UserDataService) Collect(userID uint, req request.CollectionReq) error {
	var existing model.UserCollection
	result := global.GVA_DB.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, req.TargetType, req.TargetID).Limit(1).Find(&existing)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return global.GVA_DB.Create(&model.UserCollection{UserID: userID, TargetType: req.TargetType, TargetID: req.TargetID}).Error
	}
	return nil
}

func (s *UserDataService) Uncollect(userID uint, req request.CollectionReq) error {
	return global.GVA_DB.Delete(&model.UserCollection{}, "user_id = ? AND target_type = ? AND target_id = ?", userID, req.TargetType, req.TargetID).Error
}

func (s *UserDataService) GetCollectionList(userID uint, info request.CollectionSearch) (list []model.UserCollection, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Model(&model.UserCollection{}).Where("user_id = ?", userID)
	if info.TargetType > 0 {
		db = db.Where("target_type = ?", info.TargetType)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error
	return
}

func uniqueUint(ids []uint) []uint {
	set := make(map[uint]struct{}, len(ids))
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := set[id]; ok {
			continue
		}
		set[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func (s *UserDataService) GetCollectionDetailList(userID uint, info request.CollectionSearch) (list []CollectionDetailItem, total int64, err error) {
	page, pageSize := normalizePage(info.Page, info.PageSize)
	db := global.GVA_DB.Model(&model.UserCollection{}).Where("user_id = ?", userID)
	if info.TargetType > 0 {
		db = db.Where("target_type = ?", info.TargetType)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	var rows []model.UserCollection
	err = db.Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&rows).Error
	if err != nil || len(rows) == 0 {
		return
	}

	wordIDs := make([]uint, 0)
	sentenceIDs := make([]uint, 0)
	videoIDs := make([]uint, 0)
	for _, row := range rows {
		switch row.TargetType {
		case 1:
			wordIDs = append(wordIDs, row.TargetID)
		case 2:
			sentenceIDs = append(sentenceIDs, row.TargetID)
		case 3:
			videoIDs = append(videoIDs, row.TargetID)
		}
	}

	wordIDs = uniqueUint(wordIDs)
	sentenceIDs = uniqueUint(sentenceIDs)
	videoIDs = uniqueUint(videoIDs)

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

	sentenceMap := make(map[uint]model.VideoSentence)
	if len(sentenceIDs) > 0 {
		var sentences []model.VideoSentence
		err = global.GVA_DB.Where("id IN ?", sentenceIDs).Find(&sentences).Error
		if err != nil {
			return
		}
		for _, sentence := range sentences {
			sentenceMap[sentence.ID] = sentence
		}
	}

	episodeMap := make(map[uint]model.VideoEpisode)
	seriesMap := make(map[uint]model.VideoSeries)
	if len(videoIDs) > 0 {
		var episodes []model.VideoEpisode
		err = global.GVA_DB.Where("id IN ?", videoIDs).Find(&episodes).Error
		if err != nil {
			return
		}
		seriesIDs := make([]uint, 0)
		for _, episode := range episodes {
			episodeMap[episode.ID] = episode
			seriesIDs = append(seriesIDs, episode.SeriesID)
		}

		missingSeriesIDs := make([]uint, 0)
		for _, id := range videoIDs {
			if _, ok := episodeMap[id]; !ok {
				missingSeriesIDs = append(missingSeriesIDs, id)
			}
		}

		seriesIDs = append(seriesIDs, missingSeriesIDs...)
		seriesIDs = uniqueUint(seriesIDs)
		if len(seriesIDs) > 0 {
			var seriesList []model.VideoSeries
			err = global.GVA_DB.Where("id IN ?", seriesIDs).Find(&seriesList).Error
			if err != nil {
				return
			}
			for _, series := range seriesList {
				seriesMap[series.ID] = series
			}
		}
	}

	list = make([]CollectionDetailItem, 0, len(rows))
	for _, row := range rows {
		item := CollectionDetailItem{UserCollection: row}
		switch row.TargetType {
		case 1:
			if word, ok := wordMap[row.TargetID]; ok {
				wordCopy := word
				item.Word = &wordCopy
			}
		case 2:
			if sentence, ok := sentenceMap[row.TargetID]; ok {
				sentenceCopy := sentence
				item.Sentence = &sentenceCopy
			}
		case 3:
			if episode, ok := episodeMap[row.TargetID]; ok {
				detail := &CollectionVideoDetail{
					EpisodeID:   episode.ID,
					EpisodeName: episode.Name,
					SeriesID:    episode.SeriesID,
					VideoURL:    SignLearningVideoURL(episode.VideoUrl),
				}
				if series, exists := seriesMap[episode.SeriesID]; exists {
					detail.SeriesName = series.Name
					detail.CoverURL = series.CoverUrl
				}
				item.Video = detail
			} else if series, ok := seriesMap[row.TargetID]; ok {
				seriesDetail := &CollectionVideoDetail{
					SeriesID:   series.ID,
					SeriesName: series.Name,
					CoverURL:   series.CoverUrl,
				}
				item.Video = seriesDetail
			}
		}
		list = append(list, item)
	}

	return
}

func (s *UserDataService) GetUserAsset(userID uint) (model.UserLearningAsset, error) {
	return ServiceGroupApp.LearningAuthzService.EnsureUserAsset(userID)
}
