package main

import (
	"../xlog"
	// "context"
	// "crypto/ecdsa"
	_ "encoding/hex"
	"fmt"
	"log"
	// "math/big"
	// "time"
	// "./martini"
	"net/http"
	// "./eosgo-client/common"
	// "./eosgo-client/rpc"
	// yoyow "./yoyow"
	// "encoding/json"
	// eos "github.com/eoscanada/eos-go"
	// yoyow "github.com/scorum/bitshares-go"
	// "github.com/scorum/bitshares-go/types"
	"io/ioutil"
	"os"
	"strconv"
)

type Conf struct {
	Host   string `json:"host,omitempty"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Filter string `json:"filter,omitempty"`
}

var cfg Conf //proxy.Config

func logPanics(function func(http.ResponseWriter,
	*http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Println(fmt.Sprintf("[%v] caught panic: %v", request.RemoteAddr, x))
				fmt.Println(fmt.Sprintf("[%v] caught panic: %v", request.RemoteAddr, x))
			}
		}()
		function(writer, request)
	}
}

func start() {
	s := NewRPCClient("etc", cfg.Host, "10s")
	start, err := strconv.ParseInt(cfg.From, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	end, err := strconv.ParseInt(cfg.To, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := os.OpenFile("./logs/out", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Println("os OpenFile error: ", err)
		return
	}
	defer f.Close()

	// f.WriteString("another content")

	for i := start; i <= end; i++ {
		block, err := s.GetBlockByHeight(i)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, v := range block.Transactions {
			log.Println(v.Hash)
			tx, err := s.GetTransaction(v.Hash)
			if err != nil {
				log.Println(err)
				continue
			}
			if tx.From == cfg.Filter || tx.To == cfg.Filter {
				out := tx.From + "," + tx.To + "," + tx.Value + "," + tx.Input + "\n"
				// log.Println(out)
				// writeFile(out)
				f.WriteString(out)
			}
		}
	}
}
func writeFile(content string) (err error) {
	con := []byte(content)

	err = ioutil.WriteFile("./out", con, 0644)

	if err != nil {
		log.Println("write file err:", err)
		return err
	}
	return nil
}
func main() {
	xlog.XX()
	if !LoadConfig("config.toml", &cfg) {
		return
	}
	log.Println(cfg)
	start()
	// startMartiniForEos()
	quit := make(chan bool)
	<-quit
}
