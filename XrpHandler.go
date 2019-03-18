package main

import (
	"encoding/json"
	"fmt"
	// "go/token"
	"io/ioutil"
	"net/http"
	// "bytes"
	// "time"
	// "errors"
	// "./bch"
	"./ripple"
	"log"
	// "strconv"
	// "./trx"
	// "github.com/sasaxie/go-client-api/service"
	"./rpcclient"
)

func dealwithXrpMethod(t rpcclient.JsonRequest, response *rpcclient.Response2, retError *rpcclient.ErrorStruct) {
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
	url := "http://" + cfg.Xrp.Host + ":" + cfg.Xrp.Port
	s := ripple.NewRPCClient("xrp", url, "10s", cfg.Xrp.Username, cfg.Xrp.Password)

	switch t.Method {

	case "getbalance":
		var from string
		if len(t.Params) > 0 {
			from = fmt.Sprintf("%v", t.Params[0])
		}
		b, err := s.GetBalance(from)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("balance:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "send":
		var from string
		var to string
		var amount string
		var memo string
		if len(t.Params) >= 3 {
			from = fmt.Sprintf("%v", t.Params[0])
			to = fmt.Sprintf("%v", t.Params[1])
			amount = fmt.Sprintf("%v", t.Params[2])
			memo = fmt.Sprintf("%v", t.Params[3])
		}
		b, err := s.Send(cfg.Xrp.Private, from, to, amount, memo)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("send:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "verify":
		var account string
		var hei string
		var seq string
		if len(t.Params) >= 3 {
			account = fmt.Sprintf("%v", t.Params[0])
			hei = fmt.Sprintf("%v", t.Params[1])
			seq = fmt.Sprintf("%v", t.Params[2])
		}
		b, err := s.Verify(account, hei, seq)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("verifyIn:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "gettransaction":
		var txid string
		if len(t.Params) >= 1 {
			txid = fmt.Sprintf("%v", t.Params[0])
		}
		b, err := s.GetTransaction(txid)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("gettransaction:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	default:
		retError.Message = "Method not found"
		response.Error = retError
	}

}
func XrpHandler(w http.ResponseWriter, r *http.Request) {
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
			dealwithXrpMethod(t, &response, &retError)
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
