package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"rpc-go/myFirstGrpcPackage"
	"time"
)

const (
	address = "localhost:50051"
)

//rpc实例
type goRpc struct {
	rpc myFirstGrpcPackage.GymClient

}

//通过连接初始化Gym服务客户端
func(g *goRpc) newGymClient(conn *grpc.ClientConn) {
	g.rpc = myFirstGrpcPackage.NewGymClient(conn)
}

//普通rpc请求
func(g *goRpc) commonRequest(ctx context.Context, request *myFirstGrpcPackage.RequestDemoOne) {
	r, err := g.rpc.DemoOne(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("code: %d, msg: %s", r.Code, r.Msg)
}

// 服务端流rpc请求

// 客户端流rpc请求

// 双向流rpc请求


func main() {
	fmt.Print("rpc客户端启动")
	// 建立客户端与服务端之间的grpc连接
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	gymClient:= goRpc{}
	//通过连接初始化客户端 也可以封装
	gymClient.newGymClient(conn)
	// 设置请求超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//普通rpc请求调用
	//准备请求数据
	commonRequest:= myFirstGrpcPackage.RequestDemoOne{Name:"Grumpy Omer"}
	gymClient.commonRequest(ctx, &commonRequest)
}

