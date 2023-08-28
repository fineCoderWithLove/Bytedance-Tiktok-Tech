package handler

import (
	"context"
	"douyin/base-service/global"
	"douyin/base-service/util"
	vpb "douyin/base-service/videoproto"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
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
// 辅助函数：将 []*vpb.Video 转换为 []Video
func ConvertProtoToVideos(protoVideos []*vpb.Video) []Video {
	videos := make([]Video, len(protoVideos))
	for i, protoVideo := range protoVideos {
		videos[i] = Video{
			Author: User{ // 假设 User 结构体中的字段与 vpb.User 一致
				ID:              protoVideo.Author.Id,
				Name:            protoVideo.Author.Name,
				FollowCount:     protoVideo.Author.FollowCount,
				FollowerCount:   protoVideo.Author.FollowerCount,
				IsFollow:        protoVideo.Author.IsFollow,
				Avatar:          protoVideo.Author.Avatar,
				BackgroundImage: protoVideo.Author.BackgroundImage,
				Signature:       protoVideo.Author.Signature,
				TotalFavorited:  protoVideo.Author.TotalFavorited,
				WorkCount:       protoVideo.Author.WorkCount,
				FavoriteCount:   protoVideo.Author.FavoriteCount,
			},
			VideoID:       protoVideo.Id,
			CoverURL:      protoVideo.CoverUrl,
			FavoriteCount: protoVideo.FavoriteCount,
			IsFavorite:    protoVideo.IsFavorite,
			Title:         protoVideo.Title,
		}
	}
	return videos
}

func (s *VideoServe) VideoList(ctx context.Context, req *vpb.VideoListReq) (*vpb.VideoListResp, error) {
	fmt.Println(req.UserId)
	if req.Token == "" {
		resp := &vpb.VideoListResp{
			StatusCode: 400,
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
	fmt.Println("---------------------")
	fmt.Println(user)
	// 根据 user_id 查询该用户发布的所有视频信息
	var videos []Video
	if err := global.DB.Table("videos").
		Select("videos.*, CASE WHEN favorites.user_id IS NOT NULL THEN true ELSE false END AS isfavourite").
		Joins("LEFT JOIN favorites ON videos.video_id = favorites.video_id").
		Where("videos.user_id = ?", req.UserId).
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
			FavoriteCount: video.FavoriteCount,
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
	fmt.Println(req.Token)
	fmt.Println(req.Title)
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
	fmt.Println("the filename is " + uploadFileName)
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

func (s *VideoServe) VideoStream(ctx context.Context, req *vpb.VideoStreamReq) (*vpb.VideoStreamResp, error) {
	/*
	   1.如果没有传递时间的参数，将时间倒序排列查询每次查询只查询3条记录，demo
	   2.如果有时间的参数那就按照时间来倒序查找然后将每次最后一条视频的创建时间当作下一次请求的时间
	   3.查询videolist的时候还需要查询作者的信息
	   question:如何知道是拿一个用户的第几次请求？？？
	*/
	page := 1     // 当前页数
	pageSize := 3 // 每页记录数
	offset := (page - 1) * pageSize
	var videos []Video
	var result *gorm.DB
	var nextTime int64
	zap.S().Info("准备调用视频流服务")
	// 又可能超出了传输字节的限制
	fmt.Println(req.LatestTime)
	if req.LatestTime == 0 {
		// 没有提供时间，按当前时间倒序查询
		fmt.Println("没有提供时间按照当前值查询")
		result = global.DB.
			Table("videos").
			Order("created_time DESC").
			Offset(offset).
			Limit(pageSize).
			Find(&videos)

	} else {
		fmt.Println("有参数的情况")
		latestTime := req.LatestTime
		//format := "2006-01-02 15:04:05"

		// 将时间戳转换为时间对象

		// 使用提供的时间进行倒序查询
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
	//定义返回的protobuf
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

	fmt.Println(newvideos)

	var resp vpb.VideoStreamResp
	var videoList []*vpb.Video

	for _, video := range newvideos {
		/*
		调用方法实现
		 */
		FavoriteCount, commentCount, _ := util.GetVideoFavoriteAndCommentCount(video.VideoID)
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
	fmt.Println(resp)
	return &resp, nil
}
