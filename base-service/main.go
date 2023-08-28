package main

import (
	"douyin/base-service/handler"
	pb "douyin/base-service/proto"
<<<<<<< HEAD
=======
	vpb "douyin/base-service/videoproto"
>>>>>>> cba9c25843da297a4159b839c47e609847fe7bed
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

<<<<<<< HEAD
func main()  {
	IP := flag.String("ip","127.0.0.1","ip地址")
	Port := flag.Int("port",8887,"端口号")
	flag.Parse()
	fmt.Print("ip: ",*IP)
	fmt.Print("  port: ",*Port)
	fmt.Println("  Service is running")
	server := grpc.NewServer()
	//pb.RegisterUserServer(server,&handler.UserServe{})
	pb.RegisterUserServiceServer(server,&handler.UserServe{})
	lis,err := net.Listen("tcp",fmt.Sprintf("%s:%d",*IP,*Port))
	if err != nil {
		panic("faild to liston "+err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		panic("faild to start grpc"+err.Error())
	}


}


=======
func main() {
	IP := flag.String("ip", "127.0.0.1", "ip地址")
	Port := flag.Int("port", 8887, "端口号")
	flag.Parse()
	fmt.Print("ip: ", *IP)
	fmt.Print("  port: ", *Port)
	fmt.Println("  Service is running")
	server := grpc.NewServer()
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
>>>>>>> cba9c25843da297a4159b839c47e609847fe7bed
