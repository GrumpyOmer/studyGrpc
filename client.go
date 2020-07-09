package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"rpc-go/myFirstGrpcPackage"
	"sync"
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
func (g *goRpc) newGymClient(conn *grpc.ClientConn) error {
	//通过连接初始化客户端
	g.rpc = myFirstGrpcPackage.NewGymClient(conn)
	return nil
}

//普通rpc请求
func (g *goRpc) commonRequest(request *myFirstGrpcPackage.RequestDemoOne) {
	// 设置请求超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := g.rpc.DemoOne(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("普通rpc响应：")
	fmt.Println(r)
}

// 服务端流rpc请求
func (g *goRpc) serverRequest(request *myFirstGrpcPackage.RequestDemoTwo) {
	// 设置请求超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := g.rpc.DemoTwo(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("服务端rpc流响应：")
	for {
		result,err:= r.Recv()
		//服务端发送消息完成，退出当前循环
		if err == io.EOF {
			break
		}
		//其他类型错误
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//正常流结果
		fmt.Println(result)
	}
}

// 客户端流rpc请求
func (g *goRpc) clientRequest(data *string) {
	// 设置请求超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := g.rpc.DemoThree(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for _,v:= range []byte(*data){
		if err:= r.Send(&myFirstGrpcPackage.RequestDemoThree{Name:string(v)});err!= nil {
			log.Fatalf("RequestDemoThree stream error:%v",err)
		}
	}
	res,err := r.CloseAndRecv()
	if err!= nil {
		fmt.Println("服务端err "+err.Error())
	}
	fmt.Println("客户端rpc流响应")
	fmt.Println(res)

}

// 双向流rpc请求
func (g *goRpc) twoWayRequest(data *string) {
	// 设置请求超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := g.rpc.DemoFour(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//开启协程接受处理流响应
	//设置等待组等待结果处理完毕
	var wG sync.WaitGroup
	wG.Add(1)
	go func() {
		for {
			something,err:= r.Recv()
			//服务端发送消息完成，退出当前循环
			if err == io.EOF {
				fmt.Println("two way rpc over")
				wG.Done()
				break
			}
			//其他类型错误一样退出
			if err != nil {
				fmt.Println("two way flow err: "+err.Error())
				break
			}
			fmt.Println("双向rpc流响应")
			fmt.Println(something)
		}
	}()
	for _,v:= range []byte(*data){
		if err:= r.Send(&myFirstGrpcPackage.RequestDemoFour{Name:string(v)});err!= nil {
			log.Fatalf("RequestDemoFour stream error:%v",err)
		}
	}
	//发送完毕 关闭请求
	r.CloseSend()
	wG.Wait()
}

func main() {
	fmt.Println("rpc客户端启动")
	//初始化rpc服务
	gymClient := goRpc{}
	// 建立gym客户端与gym服务端之间的grpc连接
	gymConn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer gymConn.Close()
	if err:= gymClient.newGymClient(gymConn);err!= nil {
		fmt.Print("rpc客户端初始化失败,err is "+err.Error())
		return
	}
	//普通rpc请求调用
	//准备请求数据
	commonRequest := myFirstGrpcPackage.RequestDemoOne{Name: "Grumpy Omer"}
	gymClient.commonRequest(&commonRequest)

	//服务端流rpc请求
	//准备请求数据
	serverRequest := myFirstGrpcPackage.RequestDemoTwo{Name: "Grumpy Omer"}
	gymClient.serverRequest(&serverRequest)

	//客户端流rpc请求
	//准备请求数据
	data := "Grumpy Omer"
	gymClient.clientRequest(&data)

	//双向流rpc请求
	//准备请求数据
	data = "Grumpy Omer"
	gymClient.twoWayRequest(&data)
}
