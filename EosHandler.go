package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "bytes"
	// "time"
	"./eos"
	"./rpcclient"
	"log"
	"strconv"
)

func EosHandler(w http.ResponseWriter, r *http.Request) {
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
		s := eos.New(cfg.Eos.API_METHOD, cfg.Eos.API_URL, strconv.Itoa(cfg.Eos.API_PORT))
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
			addr, err := eos.GetNewAddress(pass)
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
			var symbol string
			if len(t.Params) > 1 {
				// log.Println("t.Params:",t.Params[0])
				// addr=t.Params[0]
				addr = fmt.Sprintf("%v", t.Params[0])
				symbol = fmt.Sprintf("%v", t.Params[1])
			}
			log.Println("account:", addr)
			// json.Unmarshal(t.Params,&addr)
			// New(API_METHOD, API_URL, API_PORT string)
			// cfg.Eos.API_METHOD + "://" + cfg.Eos.API_URL + ":" + strconv.Itoa(cfg.Eos.API_PORT)

			b, err := s.GetBalance(addr, symbol)
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
			// curl --user 7MoETepmSthUn5BN:9Kb-t7w-zRV-pS7 -X POST --data-binary '{"jsonrpc":"1.0","id":"curltext","method":"send","params":["mseraccounts","mseraccount2","1.0000","usermemo","EOS"]}' -H 'content-type:text/plain;' http://127.0.0.1:18888/eos
			var from string
			var to string
			var amount string
			var userMemo string
			var symbol string
			if len(t.Params) >= 5 {
				from = fmt.Sprintf("%v", t.Params[0])
				to = fmt.Sprintf("%v", t.Params[1])
				amount = fmt.Sprintf("%v", t.Params[2])
				userMemo = fmt.Sprintf("%v", t.Params[3])
				symbol = fmt.Sprintf("%v", t.Params[4])
			}
			localeosurl := cfg.Eos.API_METHOD + "://" + cfg.Eos.LOCAL_API_URL + ":" + strconv.Itoa(cfg.Eos.LOCAL_API_PORT)
			h, err := s.Send(from, to, amount, userMemo, localeosurl, cfg.Eos.WALLET_NAME, cfg.Eos.WALLET_PRIV_KEY, symbol)
			log.Println("hash:", h)
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				bytes, _ := json.Marshal(h)
				response.Result = bytes
				response.Error = nil
			}
		} else if t.Method == "delegatebw" {
			var from string
			var to string
			var stakeCPU string
			var stakeNet string
			if len(t.Params) >= 4 {
				from = fmt.Sprintf("%v", t.Params[0])
				to = fmt.Sprintf("%v", t.Params[1])
				stakeCPU = fmt.Sprintf("%v", t.Params[2])
				stakeNet = fmt.Sprintf("%v", t.Params[3])
			}
			localeosurl := cfg.Eos.API_METHOD + "://" + cfg.Eos.LOCAL_API_URL + ":" + strconv.Itoa(cfg.Eos.LOCAL_API_PORT)
			// h, err := s.Send(from, to, amount, userMemo, localeosurl, cfg.Eos.WALLET_NAME, cfg.Eos.WALLET_PRIV_KEY, symbol)
			h, err := s.Delegatebw(from, to, stakeCPU, stakeNet, localeosurl, cfg.Eos.WALLET_NAME, cfg.Eos.WALLET_PRIV_KEY)
			log.Println("hash:", h)
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				bytes, _ := json.Marshal(h)
				response.Result = bytes
				response.Error = nil
			}
		} else if t.Method == "undelegatebw" {
			var from string
			var to string
			var stakeCPU string
			var stakeNet string
			if len(t.Params) >= 4 {
				from = fmt.Sprintf("%v", t.Params[0])
				to = fmt.Sprintf("%v", t.Params[1])
				stakeCPU = fmt.Sprintf("%v", t.Params[2])
				stakeNet = fmt.Sprintf("%v", t.Params[3])
			}
			localeosurl := cfg.Eos.API_METHOD + "://" + cfg.Eos.LOCAL_API_URL + ":" + strconv.Itoa(cfg.Eos.LOCAL_API_PORT)
			// h, err := s.Send(from, to, amount, userMemo, localeosurl, cfg.Eos.WALLET_NAME, cfg.Eos.WALLET_PRIV_KEY, symbol)
			h, err := s.Undelegatebw(from, to, stakeCPU, stakeNet, localeosurl, cfg.Eos.WALLET_NAME, cfg.Eos.WALLET_PRIV_KEY)
			log.Println("hash:", h)
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				bytes, _ := json.Marshal(h)
				response.Result = bytes
				response.Error = nil
			}
		} else if t.Method == "verifyall" { //get action
			// "mseraccounts","12345678","EOS"
			var account string
			var seq string
			var symbol string
			if len(t.Params) >= 3 {
				account = fmt.Sprintf("%v", t.Params[0])
				seq = fmt.Sprintf("%v", t.Params[1])
				symbol = fmt.Sprintf("%v", t.Params[2])
			}
			// type tempStruct struct {
			// 	TransactionId string `json:"transactionId"`
			// 	Quantity      string `json:"quantity"`
			// }
			h, err := s.VerifyAll(account, seq, symbol)
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				// fmt.Println("transactionId:", h)
				// fmt.Println("quantity", quantity)
				// outTransaction := tempStruct{
				// 	TransactionId: h,
				// 	Quantity:      quantity,
				// }
				// fmt.Println("outTransaction:", outTransaction)
				bytes, errs := json.Marshal(h)
				// log.Println("bytes:", json.RawMessage(bytes))
				// fmt.Printf("Raw Message : %s\n", bytes)
				if errs != nil {
					fmt.Printf("errs")
					retError.Message = errs.Error()
					response.Error = &retError
				} else {
					fmt.Printf("!errs")
					response.Result = bytes
					response.Error = nil
				}

			}
			// } else if t.Method == "verify" { //get action
			// 	// "mseraccounts","12345678","EOS"
			// 	var account string
			// 	var page string
			// 	if len(t.Params) >= 2 {
			// 		account = fmt.Sprintf("%v", t.Params[0])
			// 		page = fmt.Sprintf("%v", t.Params[1])
			// 	}
			// 	h, err := s.Verify(cfg.Eos.Account_History, account, page)
			// 	if err != nil {
			// 		retError.Message = err.Error()
			// 		response.Error = &retError
			// 	} else {
			// 		bytes, errs := json.Marshal(h)
			// 		if errs != nil {
			// 			fmt.Printf("errs")
			// 			retError.Message = errs.Error()
			// 			response.Error = &retError
			// 		} else {
			// 			fmt.Printf("!errs")
			// 			response.Result = bytes
			// 			response.Error = nil
			// 		}

			// 	}
		} else if t.Method == "gettransaction" { //get transaction
			// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
			var transactionId string

			if len(t.Params) >= 1 {
				transactionId = fmt.Sprintf("%v", t.Params[0])
			}
			h, err := s.GetTransaction(transactionId)
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				bytes, _ := json.Marshal(h)
				response.Result = bytes
				response.Error = nil
			}
		} else if t.Method == "getaccount" { //get transaction
			// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
			var account string

			if len(t.Params) >= 1 {
				account = fmt.Sprintf("%v", t.Params[0])
			}
			h, err := s.GetAccount(account)
			if err != nil {
				retError.Message = err.Error()
				response.Error = &retError
			} else {
				bytes, _ := json.Marshal(h)
				response.Result = bytes
				response.Error = nil
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
