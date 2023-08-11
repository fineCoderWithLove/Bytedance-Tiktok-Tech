package api

import (
	"context"
	"douyin/douyin-api/proto/favorite"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"strconv"
)

var (
	favoriteSrvClient favorite.FavoriteServiceClient
)

func init() {
	conn, err := grpc.Dial("localhost:8881", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		zap.S().Errorw("[FavoriteAction]连接【点赞服务失败】", "msg", err.Error())
	}
	favoriteSrvClient = favorite.NewFavoriteServiceClient(conn)

}
func FavoriteAction(ctx *gin.Context) {
	ctx.Query("token")
	vidStr := ctx.Query("video_id")
	vid, err := strconv.ParseInt(vidStr, 64, 10)
	actionTypeStr := ctx.Query("action_type")
	actionType, err := strconv.ParseInt(actionTypeStr, 32, 10)

	favoriteActionRequest := &favorite.FavoriteActionRequest{
		UserId:     1,
		VideoId:    vid,
		ActionType: int32(actionType),
	}

	resp, err := favoriteSrvClient.FavoriteAction(context.Background(), favoriteActionRequest)
	if err != nil {
		zap.S().Errorw("调用接口失败", "msg", err.Error())
		ctx.JSON(int(resp.StatusCode), gin.H{
			"error": resp.StatusMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": 200,
		"status_msg":  "ok",
	})
	return
}

//TODO 待定
//func IsFavorite(userId int64, videoId int64) (resp *favorite.IsFavoriteResponse) {
//	isFavoriteRequest := &favorite.IsFavoriteRequest{
//		UserId:  userId,
//		VideoId: videoId,
//	}
//	response, err := favoriteSrvClient.IsFavorite(context.Background(), isFavoriteRequest)
//	if err != nil {
//		//TODO 调用失败
//	}
//	resp = response
//
//	return
//}
