package middleware

import (
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	systemMiddleware "github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

const (
	readRateLimitWindowSeconds = 60
	readRateLimitPerWindow     = 120
)

var antiScrapePathPrefixes = []string{
	"/api/englishLearning/content/getCategoryList",
	"/api/englishLearning/content/getChapterList",
	"/api/englishLearning/content/getVideoCategoryList",
	"/api/englishLearning/content/getVideoSeriesList",
	"/api/englishLearning/content/findVideoSeries",
	"/api/englishLearning/content/findVideoEpisode",
	"/api/englishLearning/content/getVideoEpisodeList",
	"/api/englishLearning/word/getWordList",
	"/api/englishLearning/word/findWord",
	"/api/englishLearning/word/getErrorLogList",
	"/api/englishLearning/video/getSentenceList",
	"/api/englishLearning/userData/getWatchHistoryList",
	"/api/englishLearning/userData/getCollectionList",
	"/api/englishLearning/userData/getCollectionDetailList",
	"/api/englishLearning/checkin/getCheckinRecordList",
	"/api/englishLearning/checkin/getPointRecordList",
	"/api/englishLearning/asset/getFreeTimeRecordList",
	"/api/popup/getActivePopups",
}

func LearningReadRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.GVA_REDIS == nil {
			c.Next()
			return
		}
		if c.Request.Method != "GET" {
			c.Next()
			return
		}
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		if !isAntiScrapeTargetPath(path) {
			c.Next()
			return
		}

		uid := "0"
		if value, ok := c.Get(CtxEnglishUserID); ok {
			if userID, castOK := value.(uint); castOK && userID > 0 {
				uid = fmt.Sprintf("%d", userID)
			}
		}
		key := fmt.Sprintf("EL_READ_LIMIT:%s:%s:%s", uid, c.ClientIP(), path)
		if err := systemMiddleware.SetLimitWithTime(key, readRateLimitPerWindow, readRateLimitWindowSeconds); err != nil {
			c.JSON(200, gin.H{"code": 7, "data": nil, "msg": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

func isAntiScrapeTargetPath(path string) bool {
	if path == "" {
		return false
	}
	for _, prefix := range antiScrapePathPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
