package main

import (
	"douyin/douyin-api/globalinit/constant"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"interaction-service/dao"
	"interaction-service/global"

	"interaction-service/handler"
	"interaction-service/proto/comment"
	"net"
)

func main() {
	server := grpc.NewServer()
	comment.RegisterCommentServiceServer(server, &handler.CommentService{})

	listen, err := net.Listen("tcp", constant.CommentServiceAddr)
	if err != nil {
		panic(err)
	}

	dao.SetDefault(global.DB)

	global.ConsoleLogger.Info(constant.CommentServiceClientName,
		zap.String("Addr: ", constant.CommentServiceAddr),
	)
	global.InfoLogger.Info(constant.UserServiceClientName,
		zap.String("Addr: ", constant.CommentServiceAddr),
	)

	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
