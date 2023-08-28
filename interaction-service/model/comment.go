package model

import (
	"gorm.io/gorm"
	"time"
)

// Comment 评论表 /*
type Comment struct {
	CommentId  int64          `json:"comment_id" column:"comment_id" gorm:"not null;index:comment_video"`
	VideoId    int64          `json:"video_id" column:"video_id" gorm:"not null;index:comment_video"`
	UserId     int64          `json:"user_id" column:"user_id" gorm:"not null"`
	Content    string         `json:"content" column:"content"`
	CreateDate *time.Time     `json:"create_date" gorm:"autoCreateTime"column:"create_date"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" column:"deleted_at"`
}
