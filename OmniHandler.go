package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "bytes"
	// "time"
	"./rpcclient"
	// "errors"
	"log"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
		} else if t.Method == "getnewaddress" {
			var pass string
			if len(t.Params) > 0 {
				// log.Println("t.Params:",t.Params[0])
				// pass=t.Params[0]
				pass = fmt.Sprintf("%v", t.Params[0])
			}
			// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
			addr, err := getNewAddress(pass)
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				bytes, _ := json.Marshal(addr)
				response.Result = bytes
				response.Error = nil
			}
		} else if t.Method == "getbalance" {
			// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
			var addr string
			var pass string
			if len(t.Params) > 1 {
				// log.Println("t.Params:",t.Params[0])
				// addr=t.Params[0]
				addr = fmt.Sprintf("%v", t.Params[0])
				pass = fmt.Sprintf("%v", t.Params[1])
			}
			log.Println("addr:", addr)
			// json.Unmarshal(t.Params,&addr)
			b, _, err := getBalance(addr, pass)
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				log.Println("balance:", b)

				bytes, _ := json.Marshal(b)
				response.Result = bytes
				response.Error = nil
			}
		} else if t.Method == "getbtcbalance" {
			// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
			var addr string
			if len(t.Params) > 0 {
				// log.Println("t.Params:",t.Params[0])
				// addr=t.Params[0]
				addr = fmt.Sprintf("%v", t.Params[0])
			}
			log.Println("addr:", addr)
			// json.Unmarshal(t.Params,&addr)
			b, err := getBtcBalance(addr)
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				log.Println("balance:", b)

				bytes, _ := json.Marshal(b)
				response.Result = bytes
				response.Error = nil
			}
		} else if t.Method == "getaccountbalance" {
			b, err := getAccountBalance()
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				log.Println("balance:", b)

				bytes, _ := json.Marshal(b)
				response.Result = bytes
				response.Error = nil
			}
		} else if t.Method == "send" {
			// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
			var from string
			var to string
			var amount string
			var fee string
			var pass string
			if len(t.Params) >= 5 {
				// log.Println("t.Params:",t.Params[0])
				// from=t.Params[0]
				// to=t.Params[1]
				// amount=t.Params[2]
				// fee=t.Params[3]
				// pass=t.Params[4]
				from = fmt.Sprintf("%v", t.Params[0])
				to = fmt.Sprintf("%v", t.Params[1])
				amount = fmt.Sprintf("%v", t.Params[2])
				fee = fmt.Sprintf("%v", t.Params[3])
				pass = fmt.Sprintf("%v", t.Params[4])
			}
			log.Println("from:", from)
			log.Println("to:", to)
			log.Println("amount:", amount)
			log.Println("fee:", fee)
			log.Println("pass:", pass)
			// json.Unmarshal(t.Params,&addr)
			h, err := send(from, to, amount, fee, pass)
			log.Println("hash:", h)
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				bytes, _ := json.Marshal(h)
				response.Result = bytes
				response.Error = nil
			}
		} else if t.Method == "sendbtcmany" {
			// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
			var to string
			var fee string
			var h string
			var err error
			if len(t.Params) >= 2 {
				// to = fmt.Sprintf("%v", t.Params[0])
				fee = fmt.Sprintf("%v", t.Params[1])
			}
			log.Println("to:", to)
			log.Println("fee:", fee)
			var ob map[string]interface{}
			// err = json.Unmarshal([]byte(t.Params[0].(string)), &ob)
			ob, ok := t.Params[0].(map[string]interface{})
			if !ok {
				// retError.Message = err.Error()
				retError.Message = "interface to map error"
				response.Error = &retError
			} else {
				bolB, _ := json.Marshal(ob)
				h, err = sendbtcmany(string(bolB), fee)
				if err != nil {
					retError.Message = err.Error()
					response.Error = &retError
				} else {
					log.Println("hash:", h)
					bytes, _ := json.Marshal(h)
					response.Result = bytes
					response.Error = nil
				}

			}
		} else {
			retError.Message = "Method not found"
			response.Error = &retError
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
