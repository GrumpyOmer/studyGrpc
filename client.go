package main

import (
	"rpc-go/lightweight"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	fmt.Print("rpc客户端启动")

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := lightweight.NewGymClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.BodyBuilding(ctx, &lightweight.Person{
		Name: "chenqionghe",
		Actions: []string{"深蹲", "卧推", "硬拉"},
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("code: %d, msg: %s", r.Code, r.Msg)

	t, err := c.DemoBuilding(ctx, &lightweight.RequestDemo{
		Name: "guns n'r roses",
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("code: %d, msg: %s", t.Code, t.Msg)
}