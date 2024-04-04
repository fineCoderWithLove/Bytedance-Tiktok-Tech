package api

import (
	"demotest/douyin-api/global"
	"demotest/douyin-api/globalinit/constant"
	"demotest/douyin-api/proto/favorite"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
			"net/http"
)

var (
	favoriteSrvClient favorite.FavoriteServiceClient
)

func init() {
	conn, err := grpc.Dial(constant.FavoriteServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw(constant.FavoriteServiceName+"连接失败", "msg", err.Error())
		panic(err)
	}

	favoriteSrvClient = favorite.NewFavoriteServiceClient(conn)
}

func FavoriteList(ctx *gin.Context) {
	request := &favorite.FavoriteListRequest{}
	token := ctx.Query("token")
	request.Token = token
	var userIdStr string
	userIdStr = ctx.Query("user_id")
	id, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		zap.S().Error("uid错误：", userIdStr)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	} else {
		request.UserId = id
	}

	zap.S().Info("调用FavoriteList，参数：", "uid:", request.UserId, "token", request.Token)
	resp, err := favoriteSrvClient.FavoriteList(ctx, request)
	if err != nil {
		zap.S().Errorw("调用FavoriteList接口失败", "msg", err.Error())
		ctx.JSON(constant.FavoriteListErrCode, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	}

	ctx.JSON(int(resp.StatusCode), gin.H{
		"video_list":  resp.VideoList,
		"status_code": resp.StatusCode,
		"status_msg":  resp.GetStatusMsg(),
	})
	return
}
func FavoriteAction(ctx *gin.Context) {
	request := &favorite.FavoriteActionRequest{}
	uid := global.RS.Get(ctx.Query("token")).Val()

	if uid == "" {
		zap.S().Errorw("uid错误：", uid)
		ctx.JSON(constant.FavoriteActionLikeErrCode, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	} else {
		id, _ := strconv.ParseInt(uid, 10, 64)
		request.UserId = id
	}
	vidStr := ctx.Query("video_id")
	if vid, err := strconv.ParseInt(vidStr, 10, 64); err != nil {
		zap.S().Errorw("vid:%s 转换异常", vidStr)
		ctx.JSON(constant.FavoriteActionLikeErrCode, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	} else {
		request.VideoId = vid
	}
	actionTypeStr := ctx.Query("action_type")
	if actionType, err := strconv.ParseInt(actionTypeStr, 10, 32); err != nil {
		zap.S().Errorw("actionType:%s 转换异常", actionTypeStr)
		ctx.JSON(constant.FavoriteActionLikeErrCode, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	} else {
		request.ActionType = int32(actionType)
	}

	zap.S().Info("调用CommentAction，参数：", "uid:", request.UserId,
		"vid:", request.VideoId,
		"actionType:", request.ActionType)
	resp, err := favoriteSrvClient.FavoriteAction(ctx, request)
	if err != nil {
		zap.S().Errorw("调用CommentAction接口失败", "msg", err.Error())
		ctx.JSON(constant.FavoriteActionLikeErrCode, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	}

	ctx.JSON(int(resp.StatusCode), gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.GetStatusMsg(),
	})
	return
}
