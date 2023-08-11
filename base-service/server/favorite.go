package main

import (
	"douyin/base-service/dao"
	"douyin/base-service/global"
	"douyin/base-service/handler"
	"douyin/base-service/proto/favorite"
	"google.golang.org/grpc"
	"net"
)

func main() {

	server := grpc.NewServer()
	favorite.RegisterFavoriteServiceServer(server, &handler.FavoriteService{})

	listen, err := net.Listen("tcp", "127.0.0.1:8881")
	if err != nil {
		panic(err)
	}

	dao.SetDefault(global.DB)

	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
