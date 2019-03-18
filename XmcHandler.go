package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "bytes"
	// "time"
	// "errors"
	"./rpcclient"
	"./xmr"
	"log"
	// "strconv"
	// "io/ioutil"
)

func dealwithXmcMethod(t rpcclient.JsonRequest, response *rpcclient.Response2, retError *rpcclient.ErrorStruct) {

	s := xmr.NewRPCClient("xmc", cfg.Xmc.Url, cfg.Xmc.Username, cfg.Xmc.Password, "10s")
	log.Println(cfg.Xmc)
	switch t.Method {
	case "getbalance":
		var account string
		account = "0"
		if len(t.Params) > 0 {
			account = fmt.Sprintf("%v", t.Params[0])
		}

		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.GetBalance(account)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("balance:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	// case "getbalanceofaddress":
	// 	var addr string
	// 	if len(t.Params) > 0 {
	// 		addr = fmt.Sprintf("%v", t.Params[0])
	// 	}

	// 	b, err := s.GetBalanceOfAddr(addr)

	// 	if err != nil {
	// 		retError.Message = err.Error()
	// 		response.Error = retError
	// 	} else {
	// 		log.Println("balance:", b)
	// 		bytes, _ := json.Marshal(fmt.Sprintf("%.8f", b))
	// 		response.Result = bytes
	// 		response.Error = nil
	// 	}
	case "getnewaddress":
		var index string
		index = "0"
		if len(t.Params) > 0 {
			index = fmt.Sprintf("%v", t.Params[0])
		}

		b, err := s.GetAddress(index)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("newaddress:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	// // case "listtransaction":
	// // 	var addr string
	// // 	if len(t.Params) > 0 {
	// // 		addr = fmt.Sprintf("%v", t.Params[0])
	// // 	}

	// // 	// change, points, balance, err := s.GetBalance(addr)
	// // 	b, err := s.ListTransaction(addr)

	// // 	if err != nil {
	// // 		retError.Message = err.Error()
	// // 		response.Error = retError
	// // 	} else {
	// // 		log.Println("transaction list:", b)
	// // 		bytes, _ := json.Marshal(b)
	// // 		response.Result = bytes
	// // 		response.Error = nil
	// // 	}
	case "gettransaction":
		var account_index string
		var min_height string
		var max_height string
		if len(t.Params) > 2 {
			account_index = fmt.Sprintf("%v", t.Params[0])
			min_height = fmt.Sprintf("%v", t.Params[1])
			max_height = fmt.Sprintf("%v", t.Params[2])
		} else {
			retError.Message = "params is less"
			response.Error = retError
			return
		}

		b, err := s.GetTransaction(account_index, min_height, max_height)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("transaction list:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "send":
		var fromIndex string
		fromIndex = "0"
		var to string
		var amount string
		if len(t.Params) >= 3 {
			fromIndex = fmt.Sprintf("%v", t.Params[0])
			to = fmt.Sprintf("%v", t.Params[1])
			amount = fmt.Sprintf("%v", t.Params[2])
		}
		hash, f, err := s.Transfer(fromIndex, to, amount)
		h := struct {
			Hash string `json:"hash"`
			Fee  string `json:"fee"`
		}{
			Hash: hash,
			Fee:  f,
		}
		log.Println("hash:", h)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			bytes, _ := json.Marshal(h)
			response.Result = bytes
			response.Error = nil
		}

	default:
		retError.Message = "Method not found"
		response.Error = retError
	}

}
func XmcHandler(w http.ResponseWriter, r *http.Request) {
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
			dealwithXmcMethod(t, &response, &retError)
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
