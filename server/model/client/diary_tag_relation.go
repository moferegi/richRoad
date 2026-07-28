package client

// DiaryTagRelation 日记与标签关联表
type DiaryTagRelation struct {
	ID      uint `json:"ID" gorm:"primarykey"`
	DiaryID uint `json:"diaryId" gorm:"column:diary_id;index;comment:日记ID"`
	TagID   uint `json:"tagId" gorm:"column:tag_id;index;comment:日记标签ID"`
}

// TableName 自定义表名 diary_tag_relations
func (DiaryTagRelation) TableName() string {
	return "diary_tag_relations"
}