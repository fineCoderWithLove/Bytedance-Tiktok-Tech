package api

import (
	"context"
	"douyin/douyin-api/global"
	vpb "douyin/douyin-api/proto/video"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var accessKey = "VWgOopLKiABRZUwFG57_AcJ-dpJm9S31S0IQqqDg"
var secretKey = "kFfl96WyzVhY7SKCibfEIL8HktjV2fV5I3TryGeV"
var bucket = "tokendouyin"

type Video struct {
	//VideoID       int64 `gorm:"primary_key;auto_increment"`
	UserID        int64
	User          User `gorm:"foreignKey:UserID"`
	PlayURL       string
	CoverURL      string
	FavoriteCount int64
	CommentCount  int64
	//IsFavorite    bool
	Title       string
	CreatedTime int64
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

func extractVideoFrame(videoPath, imagePath string) error {
	cmd := exec.Command("./ffmpeg-6.0-essentials_build/bin/ffmpeg.exe", "-i", videoPath, "-ss", "00:00:00.001", "-vframes", "1", imagePath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
func qiniuUpload(ctx context.Context, filename string) string {
	localFile := filename
	key := filename // 上传到七牛云的文件名

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = true

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(ctx, &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return "err"
	}
	fmt.Println(ret.Key, ret.Hash)
	//说明上传图片已经完成了
	//http://ryo02ixq2.hb-bkt.clouddn.com/20230823122640_de42a5197c56350888fbca9c7f1b8eb9_raw.mp4
	return "http://ryo02ixq2.hb-bkt.clouddn.com/" + ret.Key
}

/*
视频发布的接口
*/
func VideoPublish(c *gin.Context) {
	// 获取表单参数值
	token := c.Request.FormValue("token")
	title := c.Request.FormValue("title")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您尚未登录"})
	}
	// 打印参数值
	fmt.Println("--------")
	fmt.Println(c.Request.FormValue("data"))
	fmt.Println("--------")
	fmt.Println("token是多少:", token)
	fmt.Println("title是多少:", title)
	zap.S().Info("开始调用发布视频的接口")
	// 解析上传的文件
	err := c.Request.ParseMultipartForm(32 << 24) // 设置最大文件大小
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取上传的文件
	file, handler, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// 生成唯一的文件名
	filename := time.Now().Format("20060102150405") + "_" + handler.Filename
	fmt.Println(filename)

	// 创建一个新文件来保存上传的内容
	out, err := os.Create(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer out.Close()

	// 将上传的文件内容复制到新文件中
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 调用七牛云的上传接口
	qiniuUpload(c.Request.Context(), filename)
	// 生成图片文件名
	imageFilename := time.Now().Format("20060102150405") + ".jpg"

	// 生成图片保存路径
	imagePath := "http://ryo02ixq2.hb-bkt.clouddn.com/" + imageFilename

	// 提取视频第一帧作为图片
	err = extractVideoFrame(filename, imageFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("取帧率错误")
		fmt.Println(err.Error())
	}

	// 调用七牛云的上传接口，上传图片
	imagePath = qiniuUpload(c.Request.Context(), imageFilename)
	fmt.Println(imagePath)

	// TODO 实现逻辑将token查询出userid，并且添加到mysql中
	key := token

	// 查询键的值
	value, err := global.RS.Get(string(key)).Result()
	fmt.Println("------------------------")
	fmt.Println(value)
	if value == "0" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "您还没有登录！"})
	} else {
		// 先添加把视频添加到videos库中
		// 设置时区为上海
		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			fmt.Println("Failed to load location:", err)
			return
		}

		// 获取当前时间
		now := time.Now().In(location)

		// 将时间转换为时间戳
		timestamp := now.Unix()
		userid, _ := strconv.Atoi(value)
		video := Video{
			UserID:        int64(userid),
			PlayURL:       "http://ryo02ixq2.hb-bkt.clouddn.com/" + filename,
			CoverURL:      imagePath,
			FavoriteCount: 0,
			CommentCount:  0,
			Title:         title,
			CreatedTime:   timestamp,
		}
		result := global.DB.Create(&video) // 通过数据的指针来创建
		if result != nil {
			//创建成功继续
			c.JSON(http.StatusOK, gin.H{"message": "视频上传成功！"})
			// 上传成功
		}

	}
	// 删除本地文件
	err = os.Remove(filename)
	if err != nil {
		// 处理删除文件的错误，例如打印日志或返回错误信息给客户端
	}
}

/*
2.视频流接口
*/
func VideoStream(ctx *gin.Context) {
	//因为客户端传递过来的是string类型的时间所以需要进行处理,latestTime是一个int64类型的，考虑如何进行转换
	fmt.Println("string的格式是")
	fmt.Println(ctx.Query("latest_time"))
	timestamp := int64(0)
	if ctx.Query("latest_time") != "0" {
		//latestTimeStr := ctx.Query("latest_time")
		//layout := "2006-01-02 15:04:05"
		//latestTime, err := time.Parse(layout, latestTimeStr)
		//if err != nil {
		//	// 处理解析错误情况
		//	// 例如返回错误信息给客户端或使用默认值
		//	zap.S().Errorw("无法解析 latest_time", "error", err)
		//	// 设置默认值
		//	latestTime = time.Time{}
		//}
		//timestamp = latestTime.Unix()
		seconds, _ := strconv.Atoi(ctx.Query("latest_time"))
		timestamp = int64(seconds)
	} else {
		fmt.Println("参数为空")
		timestamp = 0
	}

	zap.S().Info("[api]开始调用【VideoStream】方法")
	videoconn, err := grpc.Dial("127.0.0.1:8887", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserDetail]连接【base-service】失败，检查网络或者端口", "msg", err.Error())
	}
	defer videoconn.Close()
	// 生成 gRPC 客户端调用接口
	baseSrvClient := vpb.NewVideoServiceClient(videoconn)
	resp, err := baseSrvClient.VideoStream(context.Background(), &vpb.VideoStreamReq{
		Token:      ctx.Query("token"),
		LatestTime: timestamp,
	})

	fmt.Println("resp的数据是")
	fmt.Println(resp)

	if err != nil {
		// todo 返回失败信息
	}
	reqtime := int64(resp.NextTime)
	tm := time.Unix(reqtime, 0)
	strTimestamp := tm.Format("2006-01-02 15:04:05")
	fmt.Println(strTimestamp)

	ctx.JSON(http.StatusOK, gin.H{
		"next_time":   resp.NextTime,
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"video_list":  resp.VideoList,
	})

}

/*
3.用户发布的视频列表
1)类似于视频流接口
*/
func VideoList(ctx *gin.Context) {
	zap.S().Info("[api]开始调用【VideoList】方法")
	videoconn, err := grpc.Dial("127.0.0.1:8887", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[VideoList]连接【base-service】失败，检查网络或者端口", "msg", err.Error())
	}
	defer videoconn.Close()
	// 生成 gRPC 客户端调用接口
	baseSrvClient := vpb.NewVideoServiceClient(videoconn)
	resp, err := baseSrvClient.VideoList(context.Background(), &vpb.VideoListReq{
		Token:  ctx.Query("token"),
		UserId: ctx.Query("user_id"),
	})
	fmt.Println(resp)

	if err != nil {
		// todo 返回失败信息
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"video_list":  resp.VideoList,
	})
}
