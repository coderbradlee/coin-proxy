package trx

import (
	"fmt"
	"github.com/sasaxie/go-client-api/service"
	"testing"
)

func Test_Client(t *testing.T) {
	// if (*num < 0) || (strings.EqualFold("", *grpcAddress) && len(*grpcAddress) == 0) {
	// 	log.Fatalln("./get-block-by-num -grpcAddress localhost" +
	// 		":50051 -number <block number>")
	// }

	client := service.NewGrpcClient("192.168.1.152:50051")
	client.Start()
	defer client.Conn.Close()

	block := client.GetBlockByNum(1157236)

	fmt.Printf("block: %v\n", block)
}
