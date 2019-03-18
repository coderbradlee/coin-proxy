package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "bytes"
	// "time"
	// "errors"
	"./neo"
	"./rpcclient"
	"log"
	// "strconv"
	// "io/ioutil"
)

func dealwithNeoMethod(t rpcclient.JsonRequest, response *rpcclient.Response2, retError *rpcclient.ErrorStruct) {
	switch t.Method {
	case "getbalance":
		var account string
		// var pub string
		if len(t.Params) >= 1 {
			account = fmt.Sprintf("%v", t.Params[0])
			// pub = fmt.Sprintf("%v", t.Params[1])
		} else {
			retError.Message = "wrong params"
			response.Error = retError
			return
		}

		// change, points, balance, err := s.GetBalance(addr)
		b, err := neo.GetBalance(cfg.Neo.Url, account, cfg.Neo.NeoAssetID)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("balance:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "getnewaddress":
		var pass string
		if len(t.Params) > 0 {
			pass = fmt.Sprintf("%v", t.Params[0])
		}
		b, err := neo.GetNewAddress(cfg.Neo.Url, pass, neo_dir, cfg.Neo.WalletPass)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("newaddress:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "getheight":
		b, err := neo.GetHeight(cfg.Neo.Url)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("newaddress:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "gettransaction":
		var transactionId string
		if len(t.Params) > 0 {
			transactionId = fmt.Sprintf("%v", t.Params[0])
		}

		// change, points, balance, err := s.GetBalance(addr)
		b, err := neo.GetTransaction(cfg.Neo.Url, transactionId)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("transaction:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "getblock":
		var height string
		if len(t.Params) > 0 {
			height = fmt.Sprintf("%v", t.Params[0])
		}

		// change, points, balance, err := s.GetBalance(addr)
		b, err := neo.GetBlock(cfg.Neo.Url, height)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("getblock transaction list:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "send":
		var from string
		var to string
		var amount string
		// var assetid string
		var userPass string
		if len(t.Params) >= 4 {
			from = fmt.Sprintf("%v", t.Params[0])
			to = fmt.Sprintf("%v", t.Params[1])
			amount = fmt.Sprintf("%v", t.Params[2])
			userPass = fmt.Sprintf("%v", t.Params[3])
		}
		// coin := "NEO"
		// if assetid == cfg.Neo.GasAssetID {
		// 	coin = "GAS"
		// }
		h, err := neo.SendFrom(cfg.Neo.Url, cfg.Neo.WalletPass, from, to, amount, cfg.Neo.NeoAssetID, userPass, neo_dir)
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
func NeoHandler(w http.ResponseWriter, r *http.Request) {
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
			dealwithNeoMethod(t, &response, &retError)
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
