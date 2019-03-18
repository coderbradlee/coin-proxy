package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "bytes"
	// "time"
	"./ada"
	"./rpcclient"
	// "errors"
	"log"
	"strconv"
)

func dealwithAdaMethod(t rpcclient.JsonRequest, response *rpcclient.Response2, retError *rpcclient.ErrorStruct) {
	// s := eos.New(cfg.Eos.API_METHOD, cfg.Eos.API_URL, strconv.Itoa(cfg.Eos.API_PORT))
	// "http://127.0.0.1:8091/rpc"
	url := cfg.Ada.API_METHOD + "://" + cfg.Ada.API_URL + ":" + strconv.Itoa(cfg.Ada.API_PORT)
	s := ada.NewRPCClient("Ada", url, "10s", cfg.Ada.Capath)
	switch t.Method {
	case "getinfo":
		ret, err := s.GetInfo()
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			bytes, _ := json.Marshal(ret)
			response.Result = bytes
			response.Error = nil
		}
	case "getnewaddress":
		var pass string

		if len(t.Params) > 0 {
			pass = fmt.Sprintf("%v", t.Params[0])
		}

		// change, points, balance, err := s.GetBalance(addr)
		b, err := s.GetNewAddress(pass, cfg.Ada.AccountIndex, cfg.Ada.WalletId, cfg.Ada.SpendingPassword, ada_dir)

		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("newaddress:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "getbalanceaddr":
		var addr string
		if len(t.Params) > 0 {
			addr = fmt.Sprintf("%v", t.Params[0])
		}
		b, err := s.GetBalance(cfg.Ada.WalletId, addr)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("balance:", b)

			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "getbalance":
		b, err := s.GetBalance2(cfg.Ada.WalletId)
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
		if len(t.Params) < 1 {
			retError.Message = "param length <1"
			response.Error = retError
			return
		}
		// ob, ok := t.Params[0].(map[string]interface{})
		// if !ok {
		// 	retError.Message = "param convert error"
		// 	response.Error = retError
		// 	return
		// }
		// bolB, _ := json.Marshal(ob)
		b, err := s.Send(ada_dir, cfg.Ada.AccountIndex, cfg.Ada.WalletId, cfg.Ada.SpendingPassword, t.Params)
		if err != nil {
			retError.Message = err.Error()
			response.Error = retError
		} else {
			log.Println("send hash:", b)
			bytes, _ := json.Marshal(b)
			response.Result = bytes
			response.Error = nil
		}
	case "verifyIn":
		var addr string
		var page string
		var perpage string
		if len(t.Params) >= 3 {
			addr = fmt.Sprintf("%v", t.Params[0])
			page = fmt.Sprintf("%v", t.Params[1])
			perpage = fmt.Sprintf("%v", t.Params[2])
		}

		h, err := s.VerifyIn(cfg.Ada.AccountIndex, cfg.Ada.WalletId, addr, page, perpage)
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
	case "gettransaction":
		var id string
		if len(t.Params) >= 1 {
			id = fmt.Sprintf("%v", t.Params[0])
		}

		h, err := s.GetTransaction(cfg.Ada.AccountIndex, cfg.Ada.WalletId, id)
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
}
func AdaHandler(w http.ResponseWriter, r *http.Request) {
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
			dealwithAdaMethod(t, &response, &retError)
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
