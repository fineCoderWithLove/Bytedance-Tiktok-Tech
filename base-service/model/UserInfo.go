package model

import "time"

type Video struct {
	 VideoId int32
	 PlayUrl string
	 CoverUrl string
	 FavoriteCount int32
	 CommentCount int32
	 Title string
	 CreatedTime time.Time
}