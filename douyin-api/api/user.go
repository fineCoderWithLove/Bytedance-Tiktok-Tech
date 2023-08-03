package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	pb "douyin/douyin-api/proto" // 导入生成的 Protobuf 代码
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func GetUserVideo(ctx *gin.Context) {
	userConn, err := grpc.Dial("127.0.0.1:8887", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserVideo]连接【用户服务失败】", "msg", err.Error())
	}
	defer userConn.Close()

	// 生成 gRPC 客户端调用接口
	userSrvClient := pb.NewUserClient(userConn)
	resp, err := userSrvClient.GetUserVideo(context.Background(), &pb.UserPrimary{
		UserId: 1,
	})

	if err != nil {
		zap.S().Errorw("[GetUserVideo]调用【GetUserVideo】接口失败", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	jsonData, err := protojson.MarshalOptions{Indent: "    "}.Marshal(resp)
	if err != nil {
		zap.S().Errorw("JSON marshaling failed", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	fmt.Println(jsonData)
	ctx.Data(http.StatusOK, "application/json", jsonData)
}
