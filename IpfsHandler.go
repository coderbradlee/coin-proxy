package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "bytes"
	// "time"
	// "errors"
	// "./bch"
	"./rpcclient"
	"log"
	// "strconv"
	"./ipfs"
	// "github.com/sasaxie/go-client-api/service"
	// "encoding/hex"
)

func dealwithIpfsMethod(t rpcclient.JsonRequest, response *rpcclient.Response2, retError *rpcclient.ErrorStruct) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		// fmt.Println("c")
		if err := recover(); err != nil {
			// fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			retError.Message = fmt.Sprintf("%s", err)
			response.Error = retError
			return
		}
		// fmt.Println("d")
	}()
	s := ipfs.NewRPCClient(cfg.Ipfs.Url, cfg.Ipfs.Peerid)

	switch t.Method {

	case "add":
		var key string
		var value string
		if len(t.Params) > 1 {
			key = fmt.Sprintf("%v", t.Params[0])
			value = fmt.Sprintf("%v", t.Params[1])
		} else {
			retError.Message = "wrong params"
			response.Error = retError
			return
		}
		b, err := s.Add(key, value, ipfs_dir)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("add:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "get":
		var pass string
		if len(t.Params) > 0 {
			pass = fmt.Sprintf("%v", t.Params[0])
		} else {
			retError.Message = "wrong params"
			response.Error = retError
			return
		}
		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.Get(pass, ipfs_dir)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("Get:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	default:
		retError.Message = "Method not found"
		response.Error = retError
	}

}
func IpfsHandler(w http.ResponseWriter, r *http.Request) {
	////////////////////////////////
	addr := r.Header.Get("X-Real-IP")
	if addr == "" {
		addr = r.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = r.RemoteAddr
		}
	}
	log.Println("Method:", r.Method)
	/////////////////////////////////////////////////////////////////
	if r.Method != "POST" {
		fmt.Fprint(w, "this interface should be post!")
	} else {
		// var ret string
		body, _ := ioutil.ReadAll(r.Body)

		var t rpcclient.JsonRequest
		err_decode := json.Unmarshal(body, &t)
		log.Println("request:", t)
		defer r.Body.Close()
		var response rpcclient.Response2
		response.ID = t.ID
		var retError rpcclient.ErrorStruct
		retError.Code = -1
		// {"result":null,"error":{"code":-32601,"message":"Method not found"},"id":"curltext"}
		if err_decode != nil {
			retError.Message = "json unmarshal error"
			response.Error = &retError
		} else {
			dealwithIpfsMethod(t, &response, &retError)
		}

		var content string
		if b, err := json.Marshal(response); err == nil {
			content = string(b)
		}
		fmt.Fprint(w, content)
		log_str := fmt.Sprintf("Started %s %s for %s:%s response:%s", r.Method, r.URL.Path, addr, body, content)
		log.Println(log_str)
	}

}
