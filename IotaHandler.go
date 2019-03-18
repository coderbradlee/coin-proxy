package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "bytes"
	// "time"
	// "errors"
	"./iota"
	"./rpcclient"
	"log"
	// "strconv"
)

func dealwithIotaMethod(t rpcclient.JsonRequest, response *rpcclient.Response2, retError *rpcclient.ErrorStruct) {
	s, err := iota.New(cfg.Iota.Host, cfg.Iota.Port, cfg.Iota.Method, cfg.Iota.Seed)
	if err != nil {
		log.Println("new client:", err)
		retError.Message = err.Error()
		response.Error = retError
		return
	}
	switch t.Method {
	case "getbalance":
		var account string
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
	case "getaddress":
		var index string
		var password string
		if len(t.Params) > 1 {
			index = fmt.Sprintf("%v", t.Params[0])
			password = fmt.Sprintf("%v", t.Params[1])
		}
		b, err := s.GetNewAddress(index, password, iota_dir, cfg.Iota.Seed)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("newaddress:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "listtransaction":
		var addr string
		if len(t.Params) > 0 {
			addr = fmt.Sprintf("%v", t.Params[0])
		}

		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.ListTransaction(addr)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("listtransaction:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "listbundle":
		var addr string
		if len(t.Params) > 0 {
			addr = fmt.Sprintf("%v", t.Params[0])
		}

		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.ListBundle(addr)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("ListBundle:", b)
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
		b, err := s.GetTransaction(transactionId)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("gettransaction:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "replay":
		var transactionId string
		if len(t.Params) > 0 {
			transactionId = fmt.Sprintf("%v", t.Params[0])
		}

		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.ReplayBundle(transactionId)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("replay:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "isconfirmed":
		var transactionId string
		if len(t.Params) > 0 {
			transactionId = fmt.Sprintf("%v", t.Params[0])
		}

		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.IsConfirmed(transactionId)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("IsConfirmed:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "send":
		var from string
		var password string
		var to string
		var fromindex string
		fromindex = "0"
		var amount string
		var changeAddress string
		if len(t.Params) >= 6 {
			fromindex = fmt.Sprintf("%v", t.Params[0])
			from = fmt.Sprintf("%v", t.Params[1])
			password = fmt.Sprintf("%v", t.Params[2])
			to = fmt.Sprintf("%v", t.Params[3])

			amount = fmt.Sprintf("%v", t.Params[4])
			changeAddress = fmt.Sprintf("%v", t.Params[5])
		} else {
			retError.Message = "params error"
			response.Error = retError
			return
		}
		// h, err := s.Send(from, to, amount)
		// seeds, from, recipientAddress string, fromkeyindex, toamount
		h, err := s.Send(cfg.Iota.Seed, from, to, fromindex, amount, changeAddress, password, iota_dir)
		log.Println("hash:", h)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			bytes, _ := json.Marshal(h)
			response.Result = bytes
			response.Error = nil
		}
	// case "sendIota":
	// 	var to string
	// 	var amount string
	// 	if len(t.Params) >= 2 {
	// 		to = fmt.Sprintf("%v", t.Params[0])
	// 		amount = fmt.Sprintf("%v", t.Params[1])
	// 	} else {
	// 		retError.Message = "params error"
	// 		response.Error = retError
	// 		return
	// 	}
	// 	h, err := s.SendBtc(to, amount)
	// 	log.Println("hash:", h)
	// 	if err != nil {
	// 		retError.Message = err.Error()
	// 		response.Error = retError
	// 	} else {
	// 		bytes, _ := json.Marshal(h)
	// 		response.Result = bytes
	// 		response.Error = nil
	// 	}
	default:
		retError.Message = "Method not found"
		response.Error = retError
	}

}
func IotaHandler(w http.ResponseWriter, r *http.Request) {
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
			dealwithIotaMethod(t, &response, &retError)
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
