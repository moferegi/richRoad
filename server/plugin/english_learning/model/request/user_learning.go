package request

import commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type HeartbeatReq struct {
	EpisodeID     uint    `json:"episodeId"`
	ProgressSecs  float64 `json:"progressSecs"`
	UsingTimeSecs int     `json:"usingTimeSecs"` // 距离上次心跳间隔的实际观看秒数(如30)
}

type FreeTimeRecordSearch struct {
	commonReq.PageInfo
}
