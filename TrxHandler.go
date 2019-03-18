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
	"./trx"
	// "github.com/sasaxie/go-client-api/service"
	// "encoding/hex"
)

func dealwithTrxMethod(t rpcclient.JsonRequest, response *rpcclient.Response2, retError *rpcclient.ErrorStruct) {
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

	// s := eos.New(cfg.Eos.API_METHOD, cfg.Eos.API_URL, strconv.Itoa(cfg.Eos.API_PORT))
	// "http://127.0.0.1:8091/rpc"
	// url := cfg.Bch.API_METHOD + "://" + cfg.Bch.API_URL + ":" + strconv.Itoa(cfg.Bch.API_PORT)
	// s := yoyow.NewRPCClient("yoyow", url, "10s")
	// s, err := bch.New(cfg.Bch.Host, cfg.Bch.Port, cfg.Bch.Username, cfg.Bch.Password)
	// s := service.NewGrpcClient(cfg.Trx.Host + ":" + cfg.Trx.Port)
	// err := s.Start()
	// defer s.Conn.Close()
	url := cfg.Trx.Method + "://" + cfg.Trx.Host + ":" + cfg.Trx.Port
	s := trx.NewRPCClient("trx", url, cfg.Trx.Localurl, "10s")
	// if err != nil {
	// 	log.Println("new client:", err)
	// 	retError.Message = err.Error()
	// 	response.Error = retError
	// 	return
	// }
	switch t.Method {
	// case "getbalance":
	// 	var account string
	// 	if len(t.Params) > 0 {
	// 		account = fmt.Sprintf("%v", t.Params[0])
	// 	}

	// 	// change, points, balance, err := s.GetBalance(addr)
	// 	b, err := s.GetBlance(account)
	// 	if err != nil {
	// 		retError.Message = err.Error()
	// 		response.Error = retError
	// 	} else {
	// 		log.Println("balance:", b)
	// 		bytes, _ := json.Marshal(fmt.Sprintf("%d", b))
	// 		response.Result = bytes
	// 		response.Error = nil
	// 	}

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
		var pass string
		if len(t.Params) > 0 {
			pass = fmt.Sprintf("%v", t.Params[0])
		}
		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.GenerateAddress(trx_dir, pass, cfg.Trx.Net)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("newaddress:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "testgetnewaddress":
		// var pass string
		// if len(t.Params) > 0 {
		// 	pass = fmt.Sprintf("%v", t.Params[0])
		// }
		// change, points, balance, err := s.GetBalance(addr)
		_, _, _, b, err := s.TestGenerateAddress(cfg.Trx.Net)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("newaddress:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "getbalanceofaddress":
		var addr string
		if len(t.Params) > 0 {
			addr = fmt.Sprintf("%v", t.Params[0])
		}
		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.GetAccounts(addr)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("balance:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}

	case "getprivate":
		var addr string
		var pass string
		if len(t.Params) > 1 {
			addr = fmt.Sprintf("%v", t.Params[0])
			pass = fmt.Sprintf("%v", t.Params[1])
		}
		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.GetPrivate(trx_dir, pass, addr)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("private:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "gettransaction":
		var hash string
		if len(t.Params) >= 1 {
			hash = fmt.Sprintf("%v", t.Params[0])
		}
		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.GetTransaction(hash)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("transaction:", b)
			bytes, err := json.Marshal(b)
			if err != nil {
				retError.Message = err.Error()
				response.Error = retError
				return
			}
			response.Result = bytes
			response.Error = nil
		}
	case "send":
		var from string
		var to string
		var amount string
		var password string
		if len(t.Params) >= 3 {
			from = fmt.Sprintf("%v", t.Params[0])
			to = fmt.Sprintf("%v", t.Params[1])
			amount = fmt.Sprintf("%v", t.Params[2])
			password = fmt.Sprintf("%v", t.Params[3])
		}
		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.Send(from, to, amount, password, trx_dir, cfg.Trx.Private)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			// log.Println("broadtransaction:", hex.EncodeToString(*b))
			// bytes, err := json.Marshal(b)
			bytes, err := json.Marshal(b)

			if err != nil {
				retError.Message = err.Error()
				response.Error = retError
				return
			}
			response.Result = bytes
			response.Error = nil
		}
	// case "listtransaction":
	// 	var addr string
	// 	if len(t.Params) > 0 {
	// 		addr = fmt.Sprintf("%v", t.Params[0])
	// 	}

	// 	// change, points, balance, err := s.GetBalance(addr)
	// 	b, err := s.ListTransaction(addr)

	// 	if err != nil {
	// 		retError.Message = err.Error()
	// 		response.Error = retError
	// 	} else {
	// 		log.Println("transaction list:", b)
	// 		bytes, _ := json.Marshal(b)
	// 		response.Result = bytes
	// 		response.Error = nil
	// 	}
	// case "gettransaction":
	// 	var transactionId string
	// 	if len(t.Params) > 0 {
	// 		transactionId = fmt.Sprintf("%v", t.Params[0])
	// 	}

	// 	// change, points, balance, err := s.GetBalance(addr)
	// 	b, err := s.GetTransaction(transactionId)

	// 	if err != nil {
	// 		retError.Message = err.Error()
	// 		response.Error = retError
	// 	} else {
	// 		log.Println("transaction list:", b)
	// 		bytes, _ := json.Marshal(b)
	// 		response.Result = bytes
	// 		response.Error = nil
	// 	}
	// case "send":
	// 	var from string
	// 	var to string
	// 	var amount string
	// 	var fee string
	// 	if len(t.Params) >= 4 {
	// 		from = fmt.Sprintf("%v", t.Params[0])
	// 		to = fmt.Sprintf("%v", t.Params[1])
	// 		amount = fmt.Sprintf("%v", t.Params[2])
	// 		fee = fmt.Sprintf("%v", t.Params[3])
	// 	}
	// 	log.Println("from:", from)
	// 	log.Println("to:", to)
	// 	log.Println("amount:", amount)
	// 	log.Println("fee:", fee)
	// 	h, err := s.Send(from, to, amount, fee, cfg.Bch.WalletPass)
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
func TrxHandler(w http.ResponseWriter, r *http.Request) {
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
			dealwithTrxMethod(t, &response, &retError)
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
