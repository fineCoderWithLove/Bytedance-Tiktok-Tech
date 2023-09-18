package handler

import (
	"context"
	"demotest/base-service/global"
	"demotest/base-service/util"
	vpb "demotest/base-service/videoproto"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

/*
VideoPublish(context.Context, *PublishReq) (*PublishResp, error)
VideoList(context.Context, *VideoListReq) (*VideoListResp, error)
VideoStream(context.Context, *VideoStreamReq) (*VideoStreamResp, error)
*/
type VideoServe struct {
	vpb.UnimplementedVideoServiceServer
}
type Video struct {
	VideoID       int64
	UserID        int64
	Author        User `gorm:"foreignKey:UserID"`
	PlayURL       string
	CoverURL      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
	Title         string
	CreatedTime   int64
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

/*
	1.用户发布视频列表
	类似于视频流接口
*/

func (s *VideoServe) VideoList(ctx context.Context, req *vpb.VideoListReq) (*vpb.VideoListResp, error) {
	fmt.Println(req.UserId)
	if req.Token == "" {
		resp := &vpb.VideoListResp{
			StatusCode: 403,
			StatusMsg:  "请先登录",
			VideoList:  nil,
		}
		return resp, nil
	}
	// 根据 user_id 查询用户信息
	var user User
	if err := global.DB.Table("user").Where("id = ?", req.UserId).First(&user).Error; err != nil {
		resp := &vpb.VideoListResp{
			StatusCode: 500,
			StatusMsg:  "error",
			VideoList:  nil,
		}
		return resp, nil
	}
	// 根据 user_id 查询该用户发布的所有视频信息
	var videos []Video
	if err := global.DB.Table("videos").
		Select("videos.*, MAX(CASE WHEN favorites.user_id IS NOT NULL THEN true ELSE false END) AS isfavourite").
		Joins("LEFT JOIN favorites ON videos.video_id = favorites.video_id").
		Where("videos.user_id = ?", req.UserId).
		Group("videos.video_id").
		Find(&videos).Error; err != nil {
		resp := &vpb.VideoListResp{
			StatusCode: 500,
			StatusMsg:  err.Error(),
			VideoList:  nil,
		}
		return resp, nil
	}
	// 将查询到的用户信息放入每个视频对象中
	for i := range videos {
		videos[i].Author = user
	}

	var videoList []*vpb.Video

	for _, video := range videos {
		FavoriteCount, _, _ := util.GetVideoFavoriteAndCommentCount(video.VideoID)
		videoProto := &vpb.Video{
			Author: &vpb.User{
				Id:              video.Author.ID,
				Name:            video.Author.Name,
				FollowCount:     video.Author.FollowCount,
				FollowerCount:   video.Author.FollowerCount,
				IsFollow:        false,
				Avatar:          video.Author.Avatar,
				BackgroundImage: video.Author.BackgroundImage,
				Signature:       video.Author.Signature,
				TotalFavorited:  video.Author.TotalFavorited,
				WorkCount:       video.Author.WorkCount,
				FavoriteCount:   video.Author.FavoriteCount,
			},
			Id:            video.VideoID,
			PlayUrl:       video.CoverURL,
			CoverUrl:      video.CoverURL,
			FavoriteCount: FavoriteCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		}

		videoList = append(videoList, videoProto)
	}
	resp := vpb.VideoListResp{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
	}
	return &resp, nil
}

/*
2.用户投稿接口PublishResp
*/
func (s *VideoServe) VideoPublish(ctx context.Context, req *vpb.PublishReq) (*vpb.PublishResp, error) {

	// 配置七牛云的访问密钥和存储空间
	accessKey := "VWgOopLKiABRZUwFG57_AcJ-dpJm9S31S0IQqqDg"
	secretKey := "kFfl96WyzVhY7SKCibfEIL8HktjV2fV5I3TryGeV"
	bucket := "tokendouyin"
	// 创建七牛云的认证对象
	mac := qbox.NewMac(accessKey, secretKey)

	// 创建七牛云的上传配置
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	// 构建上传到七牛云的文件名
	uploadFileName := time.Now().Format("20060102150405") + ".mp4" // 设置上传到七牛云的文件名

	// 将字节数组写入临时文件
	tmpFile, err := ioutil.TempFile("", "temp_video.*")
	if err != nil {
		// 处理错误
		fmt.Println("wrong 111")
	}
	defer os.Remove(tmpFile.Name()) // 删除临时文件

	_, err = tmpFile.Write(req.Data)
	if err != nil {
		// 处理错误
		fmt.Println("wrong 222")
	}

	// 打开临时文件进行上传
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		// 处理错误
		fmt.Println("wrong 333")
	}

	defer file.Close()

	// 创建上传到七牛云的上传参数
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	upToken := putPolicy.UploadToken(mac)
	// 上传视频文件到七牛云
	err = formUploader.PutFile(ctx, &ret, upToken, uploadFileName, uploadFileName, &putExtra)
	if err != nil {
		// 处理错误
		fmt.Println(err.Error())
	}

	response := &vpb.PublishResp{
		StatusCode:    "200",
		StatusMessage: "success",
	}

	return response, nil
}

/*
	3.视频流接口
*/

type Attention struct {
	UserID   int64
	ToUserId int64
}

func (s *VideoServe) VideoStream(ctx context.Context, req *vpb.VideoStreamReq) (*vpb.VideoStreamResp, error) {
	uid := global.RS.Get(req.Token).Val()
	Intuid, _ := strconv.Atoi(uid)
	resp := vpb.VideoStreamResp{}
	page := 1     // 当前页数
	pageSize := 3 // 每页记录数
	offset := (page - 1) * pageSize
	var videos []Video
	var result *gorm.DB
	var nextTime int64
	zap.S().Info("准备调用视频流服务")
	// 又可能超出了传输字节的限制
	if req.LatestTime == 0 {
		// 没有提供时间，按当前时间倒序查询
		result = global.DB.
			Table("videos").
			Order("created_time DESC").
			Offset(offset).
			Limit(pageSize).
			Find(&videos)

	} else {
		latestTime := req.LatestTime
		result = global.DB.
			Table("videos").
			Where("created_time < ?", latestTime).
			Order("created_time DESC").
			Offset(offset).
			Limit(pageSize).
			Find(&videos)
	}
	//不一定会查询出来三个视频，所以需要每次查询出来遍历·
	for _, videos := range videos {
		nextTime = videos.CreatedTime
	}
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
		video.Author = user
		newvideos = append(newvideos, video)
	}
	var videoList []*vpb.Video
	for _, video := range newvideos {
		/*
			调用方法实现
		*/
		var Isfollow bool
		FavoriteCount, commentCount, _ := util.GetVideoFavoriteAndCommentCount(video.VideoID)
		var attentions []Attention
		result := global.DB.Table("attention").Where("user_id = ? AND to_user_id = ?", Intuid, video.Author.ID).Find(&attentions)
		if result.RowsAffected != 0 {
			Isfollow = true
		}
		//查询视频流的时候要从redis中查询出用户的所有信息
		TotalFavorited, userFavoriteCount, FollowerCount, FollowCount, _ := util.GetUserAllUserData(video.UserID)
		//返回获赞数量，粉丝数量，关注数量
		videoProto := &vpb.Video{
			Author: &vpb.User{
				Id:            video.Author.ID,
				Name:          video.Author.Name,
				FollowCount:   FollowCount,
				FollowerCount:  FollowerCount,
				//要查询是否关注

				IsFollow:        Isfollow,
				Avatar:          video.Author.Avatar,
				BackgroundImage: video.Author.BackgroundImage,
				Signature:       video.Author.Signature,
				TotalFavorited:  TotalFavorited,
				WorkCount:       video.Author.WorkCount,
				FavoriteCount:   userFavoriteCount,
			},
			Id:            video.VideoID,
			PlayUrl:       video.PlayURL,
			CoverUrl:      video.CoverURL,
			FavoriteCount: FavoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		}

		videoList = append(videoList, videoProto)
	}
	resp.VideoList = videoList
	//TODO  修改下一次请求的时间,nextTime如何进行返回

	resp = vpb.VideoStreamResp{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
		NextTime:   nextTime,
	}
	return &resp, nil
}
