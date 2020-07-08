package main

import (
	"rpc-go/myFirstGrpcPackage"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

// server继承自动生成的服务类
type server struct {
	myFirstGrpcPackage.UnimplementedGymServer
}

// 实现相应接口
// 普通rpc方法
func (s *server) DemoOne(ctx context.Context, in *myFirstGrpcPackage.RequestDemoOne) (*myFirstGrpcPackage.ResponseDemoOne, error) {
	fmt.Println("there is a heaven above you! "+in.Name)
	return &myFirstGrpcPackage.ResponseDemoOne{Code: 200, Msg: "ok",}, nil
}

// 服务端流rpc方法
func (s *server) DemoTwo(rect *myFirstGrpcPackage.RequestDemoTwo, stream myFirstGrpcPackage.Gym_DemoTwoServer) error {
	return nil
}

// 客户端流rpc方法
func (s *server) DemoThree(stream myFirstGrpcPackage.Gym_DemoThreeServer) error {
	return  nil
}

// 双向流rpc方法
func (s *server) DemoFour(stream myFirstGrpcPackage.Gym_DemoFourServer) error {
	return nil
}

func main() {
	fmt.Print("rpc服务端启动")
	//指定要监听客户端请求的端口
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//创建一个gRPC server的实例
	s := grpc.NewServer()
	//使用gRPC server注册服务实现
	myFirstGrpcPackage.RegisterGymServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}