package main

import (
	"fmt"
	"github.com/sasaxie/go-client-api/service"
)

func test() {
	client := service.NewGrpcClient("192.168.1.152:50051")
	client.Start()
	defer client.Conn.Close()

	block := client.GetBlockByNum(1157236)

	fmt.Println("block: ", block)
}
