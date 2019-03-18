package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "bytes"
	// "time"
	"./rpcclient"
	"./yoyow"
	"log"
	"strconv"
)

func dealwithMethod(t rpcclient.JsonRequest, response *rpcclient.Response2, retError *rpcclient.ErrorStruct) {
	// s := eos.New(cfg.Eos.API_METHOD, cfg.Eos.API_URL, strconv.Itoa(cfg.Eos.API_PORT))
	// "http://127.0.0.1:8091/rpc"
	url := cfg.Yoyow.API_METHOD + "://" + cfg.Yoyow.API_URL + ":" + strconv.Itoa(cfg.Yoyow.API_PORT)
	s := yoyow.NewRPCClient("yoyow", url, "10s")
	switch t.Method {
	case "getbalance":
		fmt.Println("CheckSync():", s.CheckSync())
		// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
		var addr string
		if len(t.Params) > 0 {
			addr = fmt.Sprintf("%v", t.Params[0])
		}
		type retStruct struct {
			Change  string `json:"change"`
			Points  string `json:"points"`
			Balance string `json:"balance"`
		}
		var b retStruct
		change, points, balance, err := s.GetBalance(addr)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			b.Change = fmt.Sprintf("%d", change)
			b.Points = fmt.Sprintf("%d", points)
			b.Balance = fmt.Sprintf("%d", balance)
			log.Println("balance:", b)

			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "getbalanceofcoin":
		fmt.Println("CheckSync():", s.CheckSync())
		// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
		var addr string
		var coin string
		if len(t.Params) > 1 {
			addr = fmt.Sprintf("%v", t.Params[0])
			coin = fmt.Sprintf("%v", t.Params[1])
		}
		b, err := s.GetBalanceOfCoin(addr, coin)
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
		fmt.Println("CheckSync():", s.CheckSync())
		if !s.CheckSync() {
			retError.Message = "node sync error"
			response.Error = retError
			return
		}
		// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
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
		err := s.Unlocks(cfg.Yoyow.WALLET_PRIV_KEY)
		if err != nil {
			log.Println("unlock:", err)
		}
		b, err := s.TransferRaw(from, to, amount, userMemo, symbol)
		s.Locks()

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("balance:", b)

			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "verifyIn":
		var account string
		var seq string
		var symbol string
		if len(t.Params) >= 3 {
			account = fmt.Sprintf("%v", t.Params[0])
			seq = fmt.Sprintf("%v", t.Params[1])
			symbol = fmt.Sprintf("%v", t.Params[2])
		}
		err := s.Unlocks(cfg.Yoyow.WALLET_PRIV_KEY)
		if err != nil {
			log.Println("unlock:", err)
		}
		h, err := s.VerifyIn(account, seq, symbol)
		s.Locks()
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			bytes, errs := json.Marshal(h)
			if errs != nil {
				// fmt.Printf("errs")
				retError.Message = errs.Error()
				response.Error = retError
			} else {
				// fmt.Printf("!errs")
				response.Result = bytes
				response.Error = nil
			}

		}
	case "verifyOut":
		var account string
		var seq string
		var symbol string
		if len(t.Params) >= 3 {
			account = fmt.Sprintf("%v", t.Params[0])
			seq = fmt.Sprintf("%v", t.Params[1])
			symbol = fmt.Sprintf("%v", t.Params[2])
		}
		err := s.Unlocks(cfg.Yoyow.WALLET_PRIV_KEY)
		if err != nil {
			log.Println("unlock:", err)
		}
		h, err := s.VerifyOut(account, seq, symbol)
		s.Locks()
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			bytes, errs := json.Marshal(h)
			if errs != nil {
				// fmt.Printf("errs")
				retError.Message = errs.Error()
				response.Error = retError
			} else {
				// fmt.Printf("!errs")
				response.Result = bytes
				response.Error = nil
			}

		}
	case "verify":
		var account string
		var seq string
		if len(t.Params) >= 2 {
			account = fmt.Sprintf("%v", t.Params[0])
			seq = fmt.Sprintf("%v", t.Params[1])
		}
		err := s.Unlocks(cfg.Yoyow.WALLET_PRIV_KEY)
		if err != nil {
			log.Println("unlock:", err)
		}
		h, err := s.VerifyAll(account, seq)
		s.Locks()
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			bytes, errs := json.Marshal(h)
			if errs != nil {
				// fmt.Printf("errs")
				retError.Message = errs.Error()
				response.Error = retError
			} else {
				// fmt.Printf("!errs")
				response.Result = bytes
				response.Error = nil
			}

		}
	default:
		retError.Message = "Method not found"
		response.Error = retError
	}
	// if t.Method == "getnewaddress" {

	// } else if t.Method == "getbalance" {

	// } else if t.Method == "send" {
	// 	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	// 	// curl --user 7MoETepmSthUn5BN:9Kb-t7w-zRV-pS7 -X POST --data-binary '{"jsonrpc":"1.0","id":"curltext","method":"send","params":["mseraccounts","mseraccount2","1.0000","usermemo","EOS"]}' -H 'content-type:text/plain;' http://127.0.0.1:18888/eos
	// 	var from string
	// 	var to string
	// 	var amount string
	// 	var userMemo string
	// 	var symbol string
	// 	if len(t.Params) >= 5 {
	// 		from = fmt.Sprintf("%v", t.Params[0])
	// 		to = fmt.Sprintf("%v", t.Params[1])
	// 		amount = fmt.Sprintf("%v", t.Params[2])
	// 		userMemo = fmt.Sprintf("%v", t.Params[3])
	// 		symbol = fmt.Sprintf("%v", t.Params[4])
	// 	}
	// 	localeosurl := cfg.Eos.API_METHOD + "://" + cfg.Eos.LOCAL_API_URL + ":" + strconv.Itoa(cfg.Eos.LOCAL_API_PORT)
	// 	h, err := s.Send(from, to, amount, userMemo, localeosurl, cfg.Eos.WALLET_NAME, cfg.Eos.WALLET_PRIV_KEY, symbol)
	// 	log.Println("hash:", h)
	// 	if err != nil {
	// 		retError.Message = err.Error()
	// 		response.Error = &retError
	// 	} else {
	// 		bytes, _ := json.Marshal(h)
	// 		response.Result = bytes
	// 		response.Error = nil
	// 	}
	// } else if t.Method == "verifyIn" { //get action
	// 	// "mseraccounts","12345678","EOS"
	// 	var account string
	// 	var seq string
	// 	var symbol string
	// 	if len(t.Params) >= 3 {
	// 		account = fmt.Sprintf("%v", t.Params[0])
	// 		seq = fmt.Sprintf("%v", t.Params[1])
	// 		symbol = fmt.Sprintf("%v", t.Params[2])
	// 	}
	// 	// type tempStruct struct {
	// 	// 	TransactionId string `json:"transactionId"`
	// 	// 	Quantity      string `json:"quantity"`
	// 	// }
	// 	h, err := s.VerifyIn(account, seq, symbol)
	// 	if err != nil {
	// 		retError.Message = err.Error()
	// 		response.Error = &retError
	// 	} else {
	// 		// fmt.Println("transactionId:", h)
	// 		// fmt.Println("quantity", quantity)
	// 		// outTransaction := tempStruct{
	// 		// 	TransactionId: h,
	// 		// 	Quantity:      quantity,
	// 		// }
	// 		// fmt.Println("outTransaction:", outTransaction)
	// 		bytes, errs := json.Marshal(h)
	// 		// log.Println("bytes:", json.RawMessage(bytes))
	// 		// fmt.Printf("Raw Message : %s\n", bytes)
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
	// } else if t.Method == "gettransaction" { //get transaction
	// 	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	// 	var transactionId string

	// 	if len(t.Params) >= 1 {
	// 		transactionId = fmt.Sprintf("%v", t.Params[0])
	// 	}
	// 	h, err := s.VerifyOut(transactionId)
	// 	if err != nil {
	// 		retError.Message = err.Error()
	// 		response.Error = &retError
	// 	} else {
	// 		bytes, _ := json.Marshal(h)
	// 		response.Result = bytes
	// 		response.Error = nil
	// 	}
	// } else if t.Method == "getaccount" { //get transaction
	// 	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	// 	var account string

	// 	if len(t.Params) >= 1 {
	// 		account = fmt.Sprintf("%v", t.Params[0])
	// 	}
	// 	h, err := s.GetAccount(account)
	// 	if err != nil {
	// 		retError.Message = err.Error()
	// 		response.Error = &retError
	// 	} else {
	// 		bytes, _ := json.Marshal(h)
	// 		response.Result = bytes
	// 		response.Error = nil
	// 	}
	// } else {
	// 	retError.Message = "Method not found"
	// 	response.Error = &retError
	// }
}
func YoyowHandler(w http.ResponseWriter, r *http.Request) {
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
			dealwithMethod(t, &response, &retError)
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
