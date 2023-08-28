package test

import (
	"douyin/base-service/global"
	"fmt"
	"testing"
)

type Video struct {
	VideoID       int64
	UserID        int64
	User          User `gorm:"foreignKey:UserID"`
	PlayURL       string
	CoverURL      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
	Title         string
}

type User struct {
	ID              int64
	Name            string
	FollowCount     int64
	FollowerCount   int64
	IsFollow        bool
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  int64
	WorkCount       int64
	FavoriteCount   int64
}

func TestVideoStream(t *testing.T) {
	var videos []Video
	page := 1     // 当前页数
	pageSize := 3 // 每页记录数

	offset := (page - 1) * pageSize

	result := global.DB.
		Table("videos").
		Order("created_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&videos)

	if result.Error != nil {
		// 处理查询错误
	}
	var newvideos []Video
	// 处理查询结果
	for _, video := range videos {
		// 处理视频信息
		var user User
		result = global.DB.
			Table("user").
			Where(video.UserID).
			Find(&user)
		// ...
		fmt.Println("查询出来的user信息")
		fmt.Println(user)
		video.User = user
		fmt.Println("封装的user信息")
		fmt.Println(video)
		newvideos = append(newvideos, video)
	}
	fmt.Println(newvideos)

}
func TestSelectUser(t *testing.T) {
	var user User

	resp := global.DB.Table("user").
		First(&user, 20)
	if resp.Error != nil {
		fmt.Println("查询失败:", resp.Error)
		return
	}

	fmt.Println("User ID:", user.ID)
	fmt.Println("User Name:", user.Name)
	fmt.Println("Follow Count:", user.FollowCount)
	fmt.Println("Follower Count:", user.FollowerCount)
	fmt.Println("Avatar:", user.Avatar)
	fmt.Println("Background Image:", user.BackgroundImage)
	fmt.Println("Signature:", user.Signature)
	fmt.Println("Total Favorited:", user.TotalFavorited)
	fmt.Println("Work Count:", user.WorkCount)
	fmt.Println("Favorite Count:", user.FavoriteCount)
}

func TestVideoList(t *testing.T) {
	// 根据 user_id 查询用户信息
	var user User
	if err := global.DB.Table("user").Where("id = ?", 20).First(&user).Error; err != nil {
		t.Errorf("failed to query user: %v", err)
		return
	}
	fmt.Println(user)

	// 根据 user_id 查询该用户发布的所有视频信息
	var videos []Video
	if err := global.DB.Table("videos").
		Select("videos.*, CASE WHEN favorites.user_id IS NOT NULL THEN true ELSE false END as isfavourite").
		Joins("LEFT JOIN favorites ON videos.video_id = favorites.video_id AND favorites.user_id = ?", 20).
		Find(&videos).Error; err != nil {
		t.Errorf("failed to query videos: %v", err)
		return
	}

	// 将查询到的用户信息放入每个视频对象中
	for i := range videos {
		videos[i].User = user
	}
	fmt.Println(videos)
}
