package main

import (
	"douyin/douyin-api/globalinit/constant"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"interaction-service/dao"
	"interaction-service/global"
	"interaction-service/handler"
	"interaction-service/proto/favorite"
	"net"
)

func main() {
	server := grpc.NewServer()
	favorite.RegisterFavoriteServiceServer(server, &handler.FavoriteService{})

	listen, err := net.Listen("tcp", constant.FavoriteServiceAddr)
	if err != nil {
		panic(err)
	}

	dao.SetDefault(global.DB)

	global.ConsoleLogger.Info(constant.FavoriteServiceName,
		zap.String("Addr: ", constant.FavoriteServiceAddr),
	)
	global.InfoLogger.Info(constant.FavoriteServiceName,
		zap.String("Addr: ", constant.FavoriteServiceAddr),
	)

	err = server.Serve(listen)

	if err != nil {
		panic(err)
	}
}
