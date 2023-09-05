package api

import (
	"context"
	mpb "demotest/douyin-api/proto/message"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net/http"
	"strconv"
	"time"
)

/*
	5.用户发送消息接口 Post
*/
func UserMessageAction(ctx *gin.Context) {

	relationConn, err := grpc.Dial("127.0.0.1:8886", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[UserRelationAction]连接【social-service】失败，检查网络或者端口", "msg", err.Error())
	}

	defer relationConn.Close()
	// 生成 gRPC 客户端调用接口
	baseSrvClient := mpb.NewMessageServerClient(relationConn);
	touserid, _ := strconv.Atoi(ctx.Query("to_user_id"))
	actiontype, _ := strconv.Atoi(ctx.Query("action_type"))
	//接收参数
	resp, err := baseSrvClient.MessageAction(context.Background(), &mpb.MessageActionReq{
		Token:      ctx.Query("token"),
		ToUserId:   int64(touserid),
		ActionType: int32(actiontype),
		Content:    ctx.Query("content"),
	})
	fmt.Println(resp)
	if err != nil {
		zap.S().Errorw("[api]调用【UserMessageAction】接口失败", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
	})
}

/*
	6.查询用户聊天记录
*/
func UserMessageChat(ctx *gin.Context) {
	conn, err := grpc.Dial(
		"127.0.0.1:8886",
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
			Time:                30 * time.Second, // 客户端发送 keep-alive ping 的时间间隔
			Timeout:             20 * time.Second, // 客户端等待 keep-alive ping 的响应超时时间
			PermitWithoutStream: true,             // 允许在没有活动流的情况下发送 keep-alive ping
		}),
	)
	defer conn.Close()

	baseSrvClient := mpb.NewMessageServerClient(conn);

	touserid, _ := strconv.Atoi(ctx.Query("to_user_id"))
	PreMsgTime, _ := strconv.Atoi(ctx.Query("pre_msg_time"))
	//接收参数
	resp, err := baseSrvClient.MessageChat(context.Background(), &mpb.MessageChatReq{
		Token:      ctx.Query("token"),
		ToUserId:   int64(touserid),
		PreMsgTime: int64(PreMsgTime),
	})
	if err != nil {
		zap.S().Errorw("[api]调用【UserMessageAction】接口失败", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code":  resp.StatusCode,
			"status_msg":   resp.StatusMsg,
			"message_list": resp.MessageList,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code":  resp.StatusCode,
		"status_msg":   resp.StatusMsg,
		"message_list": resp.MessageList,
	})
}
