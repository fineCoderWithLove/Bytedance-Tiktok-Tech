package model

import "time"

type Video struct {
	VideoId       int32
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int32
	CommentCount  int32
	Title         string
	CreatedTime   time.Time
}
type User struct {
	Id              int64
	Name            string
	FollowCount     int64
	FollowerCount   int64
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  int64
	WorkCount       int64
	FavoriteCount   int64
	Password        string
}

// 注册时候返回的数据
type UserResp struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}
