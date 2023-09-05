package util

import (
	"github.com/gin-gonic/gin"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"time"
)
//利用连接池对象生成一个连接
func InitFactory(address string, MaxConn int, MaxRelaxedConn int, ctx *gin.Context) *grpcpool.ClientConn {
	// 创建 gRPC 连接的工厂函数
	factory := func() (*grpc.ClientConn, error) {
		// 创建 gRPC 连接
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	// 创建连接池
	pool, err := grpcpool.New(factory, MaxConn, MaxRelaxedConn, time.Minute)
	if err != nil {
		// 处理错误
	}
	defer pool.Close()
	conn, err := pool.Get(ctx)
	if err != nil {
		// 处理错误
	}
	return conn
}
