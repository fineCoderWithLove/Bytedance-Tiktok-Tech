package api

import (
	"context"
	"douyin/douyin-api/proto/user"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

/*
用户的详情
*/
func GetUserDetail(ctx *gin.Context) {
	userConn, err := grpc.Dial("127.0.0.1:8887", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserDetail]连接【base-service】失败，检查网络或者端口", "msg", err.Error())
	}
	defer userConn.Close()

	// 生成 gRPC 客户端调用接口
	baseSrvClient := user.NewUserServiceClient(userConn)

	resp, err := baseSrvClient.UserDetail(context.Background(), &user.DetailRep{
		UserId: ctx.Query("user_id"),
		Token:  ctx.Query("token"),
	})
	if err != nil {
		zap.S().Errorw("[api]调用【GetUserDetail】接口失败", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"status_msg":  "Internal Server Error",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  "success",
		"user":        resp.User,
	})
}

/*
用户的注册接口,把请求头转json输出再去格式化一下
*/
func PrintRequestHeaders(r *http.Request) {
	headers := make(map[string]string)
	for key, value := range r.Header {
		headers[key] = value[0]
	}

	jsonData, err := json.MarshalIndent(headers, "", "  ")
	if err != nil {
		fmt.Println("Failed to convert headers to JSON:", err.Error())
		return
	}

	fmt.Println(string(jsonData))
}

func UserRegister(ctx *gin.Context) {
	// 设置 Content-Type
	//ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	zap.S().Info("[api]开始调用【UserRegister】方法")
	userConn, err := grpc.Dial("127.0.0.1:8887", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserDetail]连接【base-service】失败，检查网络或者端口", "msg", err.Error())
	}
	defer userConn.Close()
	PrintRequestHeaders(ctx.Request)
	fmt.Println("从请求中获取的值")
	fmt.Println(ctx.Query("username"))
	fmt.Println(ctx.Query("password"))
	fmt.Println("----------------------------------------")
	// 生成 gRPC 客户端调用接口
	baseSrvClient := user.NewUserServiceClient(userConn)
	resp, err := baseSrvClient.UserRegister(context.Background(), &user.RegisterReq{
		Username: ctx.Query("username"),
		Password: ctx.Query("password"),
	})
	if err != nil {
		zap.S().Errorw("[api]调用【UserRegister】接口失败", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"status_msg":  "Internal Server Error",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_id":     resp.UserId,
		"token":       resp.Token,
	})
}

/*
用户的登录接口
*/
func UserLogin(ctx *gin.Context) {
	zap.S().Info("[api]开始调用【UserRegister】方法")
	userConn, err := grpc.Dial("127.0.0.1:8887", grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserDetail]连接【base-service】失败，检查网络或者端口", "msg", err.Error())
	}
	defer userConn.Close()
	fmt.Println(ctx.Query("username"))
	fmt.Println(ctx.Query("password"))
	fmt.Println("----------------------------------------")
	// 生成 gRPC 客户端调用接口
	baseSrvClient := user.NewUserServiceClient(userConn)
	resp, err := baseSrvClient.UserLogin(context.Background(), &user.LoginReq{
		Username: ctx.Query("username"),
		Password: ctx.Query("password"),
	})
	if err != nil {
		zap.S().Errorw("[api]调用【UserRegister】接口失败", "msg", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"status_msg":  "Internal Server Error",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_id":     resp.UserId,
		"token":       resp.Token,
	})
}
