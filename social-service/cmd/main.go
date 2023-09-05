package main

import (
	"demotest/social-service/handler"
	mpb "demotest/social-service/proto/message"
	vpb "demotest/social-service/proto/relation"
	"flag"
	"fmt"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"

	"google.golang.org/grpc"
)

func main() {
	IP := flag.String("ip", "127.0.0.1", "ip地址")
	Port := flag.Int("port", 8886, "端口号")
	flag.Parse()
	fmt.Print("ip: ", *IP)
	fmt.Print("  port: ", *Port)
	fmt.Println("  Service is running")
	//让连接进行复用。
	server := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second, // 服务器在收到 keep-alive ping 后的最小等待时间
			PermitWithoutStream: true,            // 允许在没有活动流的情况下发送 keep-alive ping
		}),
	)
	//注册Relation
	vpb.RegisterUserServiceServer(server, &handler.RelationServer{})
	//注册Message
	mpb.RegisterMessageServerServer(server,&handler.MessageServer{})
	lis, err := net.Listen("tcp","127.0.0.1:8886")
	if err != nil {
		panic("faild to liston " + err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		panic("faild to start grpc" + err.Error())
	}

}
