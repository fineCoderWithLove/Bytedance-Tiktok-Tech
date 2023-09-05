package api

import (
	"context"
	rpb "demotest/douyin-api/proto/relation"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

/*
	1.用户关注操作接口 Post
*/
func UserRelationAction(ctx *gin.Context) {
	relationConn, err := grpc.Dial("127.0.0.1:8886", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[UserRelationAction]连接【social-service】失败，检查网络或者端口", "msg", err.Error())
	}

	defer relationConn.Close()
	// 生成 gRPC 客户端调用接口
	baseSrvClient := rpb.NewUserServiceClient(relationConn);
	userid, _ := strconv.Atoi(ctx.Query("to_user_id"))
	actiontype, _ := strconv.Atoi(ctx.Query("action_type"))
	resp, err := baseSrvClient.UserRelationAction(context.Background(), &rpb.UserRelationActionReq{
		Token:      ctx.Query("token"),
		ToUserId:   int64(userid),
		ActionType: int32(actiontype),
	})
	fmt.Println(resp)
	if err != nil {
		zap.S().Errorw("[api]调用【UserRelationAction】接口失败", "msg", err.Error())
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
	2.查询用户的关注列表接口  Get
*/
func UserRelationFollowList(ctx *gin.Context) {
	relationConn, err := grpc.Dial("127.0.0.1:8886", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[UserRelationFollowList]连接【social-service】失败，检查网络或者端口", "msg", err.Error())
	}
	defer relationConn.Close()
	// 生成 gRPC 客户端调用接口
	baseSrvClient := rpb.NewUserServiceClient(relationConn);
	userid, _ := strconv.Atoi(ctx.Query("user_id"))
	resp, err := baseSrvClient.UserRelationFollowList(context.Background(), &rpb.UserRelationFollowListReq{
		Token:  ctx.Query("token"),
		UserId: int64(userid),
	})
	if err != nil {
		zap.S().Errorw("[api]调用【UserRelationFollowList】接口失败", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"user_list":   resp.UserList,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_list":   resp.UserList,
	})

}

/*
	3.查询用户粉丝列表接口 Get
*/
func UserRelationFollowerList(ctx *gin.Context) {
	relationConn, err := grpc.Dial("127.0.0.1:8886", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[UserRelationFollowList]连接【social-service】失败，检查网络或者端口", "msg", err.Error())
	}
	defer relationConn.Close()
	// 生成 gRPC 客户端调用接口
	baseSrvClient := rpb.NewUserServiceClient(relationConn);
	userid, _ := strconv.Atoi(ctx.Query("user_id"))
	resp, err := baseSrvClient.UserRelationFollowerList(context.Background(), &rpb.UserRelationFollowerListReq{
		Token:  ctx.Query("token"),
		UserId: int64(userid),
	})
	if err != nil {
		zap.S().Errorw("[api]调用【UserRelationFollowList】接口失败", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"user_list":   resp.UserList,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_list":   resp.UserList,
	})
}

/*
	4.查询用户好友列表接口  Get
*/
func UserRelationFriendList(ctx *gin.Context) {
	relationConn, err := grpc.Dial("127.0.0.1:8886", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[UserRelationFriendList]连接【social-service】失败，检查网络或者端口", "msg", err.Error())
	}
	defer relationConn.Close()
	// 生成 gRPC 客户端调用接口
	baseSrvClient := rpb.NewUserServiceClient(relationConn);
	userid, _ := strconv.Atoi(ctx.Query("user_id"))
	resp, err := baseSrvClient.UserRelationFriendList(context.Background(), &rpb.UserRelationFriendListReq{
		Token:  ctx.Query("token"),
		UserId: int64(userid),
	})
	if err != nil {
		zap.S().Errorw("[api]调用【UserRelationFollowList】接口失败", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"user_list":   resp.UserList,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_list":   resp.UserList,
	})
}

