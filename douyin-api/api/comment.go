package api

import (
	"demotest/douyin-api/global"
	"demotest/douyin-api/globalinit/constant"
	"demotest/douyin-api/proto/comment"
	"fmt"
	"strconv"
			"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	commentServiceClient comment.CommentServiceClient
)

func init() {
	conn, err := grpc.Dial(constant.CommentServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw(constant.CommentServiceName+"连接失败", "msg", err.Error())
		panic(err)
	}
	commentServiceClient = comment.NewCommentServiceClient(conn)

}
func CommentAction(ctx *gin.Context) {
	req := &comment.CommentActionRequest{}

	if token := ctx.Query("token"); token == "" {
		zap.S().Errorw("token不存在")
		ctx.JSON(constant.CommentActionErrCode, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	} else {
		req.Token = token
		uid := global.RS.Get(token).Val()
		if uid == "" {
			zap.S().Errorw("uid错误：", uid)
			ctx.JSON(constant.CommentActionErrCode, gin.H{
				"error": constant.ErrorMsg,
			})
			return
		} else {
			id, _ := strconv.ParseInt(uid, 10, 64)
			req.UserId = id
		}
	}

	vidStr := ctx.Query("video_id")
	actionTypeStr := ctx.Query("action_type")
	if vidStr == "" || actionTypeStr == "" {
		zap.S().Errorw("参数错误")
		ctx.JSON(constant.CommentActionErrCode, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	} else {

		if vid, err := strconv.ParseInt(vidStr, 10, 64); err != nil {
			fmt.Println(vidStr)
			zap.S().Errorw("vid:%s 转换异常", vid)
			ctx.JSON(constant.CommentActionErrCode, gin.H{
				"error": constant.ErrorMsg,
			})
			return
		} else {
			req.VideoId = vid
		}

		if actionType, err := strconv.ParseInt(actionTypeStr, 10, 64); err != nil {
			zap.S().Errorw("actionType:%s 转换异常", actionType)
			ctx.JSON(constant.CommentActionErrCode, gin.H{
				"error": constant.ErrorMsg,
			})
			return
		} else {
			req.ActionType = actionType
		}
	}

	if req.ActionType == 1 {
		if commentText := ctx.Query("comment_text"); commentText == "" {
			zap.S().Errorw("commentText不能为空")
			ctx.JSON(constant.CommentActionErrCode, gin.H{
				"error": constant.ErrorMsg,
			})
			return
		} else {
			req.CommentText = &commentText
		}
	} else if req.ActionType == 2 {
		commentIdStr := ctx.Query("comment_id")
		if commentIdStr == "" {
			zap.S().Errorw("comment_id不能为空")
			ctx.JSON(constant.CommentActionErrCode, gin.H{
				"error": constant.ErrorMsg,
			})
			return
		} else {
			req.CommentId = &commentIdStr
		}
	} else {
		zap.S().Errorw("请求参数错误")
		ctx.JSON(constant.CommentActionErrCode, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	}
	zap.S().Info("调用CommentAction，参数：token:", req.Token,
		"uid:", req.UserId,
		"vid:", req.VideoId,
		"actionType:", req.ActionType,
		"commentId:", req.CommentId,
		"commentText", req.CommentText)
	resp, err := commentServiceClient.CommentAction(ctx, req)
	if err != nil {
		zap.S().Errorw("CommentAction调用失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	}

	ctx.JSON(0, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"comment":     resp.Comment,
	})
	return
}
func CommentList(ctx *gin.Context) {

	request := &comment.CommentListRequest{}
	if token := ctx.Query("token"); token == "" {
		zap.S().Errorw("token不存在")
	} else {
		request.Token = token
	}

	vidStr := ctx.Query("video_id")
	if vid, err := strconv.ParseInt(vidStr, 10, 64); err != nil {
		zap.S().Errorw("video_id转换失败")
		ctx.JSON(constant.CommentActionErrCode, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	} else {
		request.VideoId = vid
	}

	zap.S().Info("调用CommentList，参数：token:", request.Token,
		"vid:", request.VideoId)
	resp, err := commentServiceClient.CommentList(ctx, request)
	if err != nil {
		zap.S().Errorw("调用CommentList失败")
		ctx.JSON(constant.CommentActionErrCode, gin.H{
			"error": constant.ErrorMsg,
		})
		return
	}

	ctx.JSON(int(resp.StatusCode), gin.H{
		"status_code":  resp.StatusCode,
		"status_msg":   resp.StatusMsg,
		"comment_list": resp.CommentList,
	})

}
