package client

// VideoEpisodeTag 视频单集与标签关联表
type VideoEpisodeTag struct {
	ID        uint `json:"ID" gorm:"primarykey"`
	EpisodeID uint `json:"episodeId" gorm:"column:episode_id;index;comment:视频单集ID"`
	TagID     uint `json:"tagId" gorm:"column:tag_id;index;comment:视频标签ID"`
}

// TableName 自定义表名 video_episode_tags
func (VideoEpisodeTag) TableName() string {
	return "video_episode_tags"
}
