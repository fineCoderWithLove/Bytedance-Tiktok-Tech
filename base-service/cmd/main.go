package main

import (
	"demotest/base-service/handler"
	pb "demotest/base-service/proto"
	vpb "demotest/base-service/videoproto"
	"flag"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	IP := flag.String("ip", "127.0.0.1", "ip地址")
	Port := flag.Int("port", 8887, "端口号")
	flag.Parse()
	fmt.Print("ip: ", *IP)
	fmt.Print("  port: ", *Port)
	fmt.Println("  Service is running")
	server := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second, // 服务器在收到 keep-alive ping 后的最小等待时间
			PermitWithoutStream: true,            // 允许在没有活动流的情况下发送 keep-alive ping
		}),
	)
	//pb.RegisterUserServer(server,&handler.UserServe{})
	pb.RegisterUserServiceServer(server, &handler.UserServe{})
	vpb.RegisterVideoServiceServer(server, &handler.VideoServe{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("faild to liston " + err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		panic("faild to start grpc" + err.Error())
	}

}
