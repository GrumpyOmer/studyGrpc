package main

import (
	"io"
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
	fmt.Println("普通rpc请求")
	//获取客户端发送的消息
	fmt.Println("rpc: "+in.Name)
	return &myFirstGrpcPackage.ResponseDemoOne{Code: 200, Msg: in.Name,}, nil
}

// 服务端流rpc方法
func (s *server) DemoTwo(rect *myFirstGrpcPackage.RequestDemoTwo, stream myFirstGrpcPackage.Gym_DemoTwoServer) error {
	fmt.Println("服务端rpc流请求")
	//获取客户端发送的消息
	fmt.Println( "server flow: "+rect.Name)
	result := []byte(rect.Name)
	for _,v:= range result{
		//服务端流式响应
		if err:= stream.Send(&myFirstGrpcPackage.ResponseDemoTwo{Code:200, Msg:string(v)});err!= nil{
			return err
		}
	}
	return nil
}

// 客户端流rpc方法
func (s *server) DemoThree(stream myFirstGrpcPackage.Gym_DemoThreeServer) error {
	fmt.Println("客户端rpc流请求")
	//客户端流结果集
	var result string
	for {
		//获取客户端流式请求
		something,err := stream.Recv()
		//获取客户端发送的消息
		fmt.Println("client flow: "+something.GetName())
		//接受完毕全部客户端请求 直接响应结果
		if err == io.EOF {
			return stream.SendAndClose(&myFirstGrpcPackage.ResponseDemoThree{Code:200, Msg: result})
		}
		//其他错误
		if err != nil {
			return err
		}
		result += something.GetName()
	}
}

// 双向流rpc方法
func (s *server) DemoFour(stream myFirstGrpcPackage.Gym_DemoFourServer) error {
	fmt.Println("双向rpc流请求")
	for {
		//获取客户端流请求
		something,err := stream.Recv()
		//获取客户端流消息
		fmt.Println("two way flow: "+something.GetName())
		//接受完毕全部客户端流请求
		if err == io.EOF {
			return nil
		}
		//其他错误
		if err != nil {
			return err
		}
		//响应服务端流

		//方式1 边接收边响应
		if err := stream.Send(&myFirstGrpcPackage.ResponseDemoFour{Code:200,Msg:something.GetName()});err != nil{
			return err
		}
		//方式2 接受所有结果后一次性响应 参考客户端流rpc方法
	}
}

func main() {
	fmt.Println("rpc服务端启动")
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