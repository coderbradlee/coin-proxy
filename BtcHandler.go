package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	// "errors"
	// "./bch"
	// "./rpcclient"
	"log"
	// "strconv"
	// "./ipfs"
	// "github.com/sasaxie/go-client-api/service"
	// "encoding/hex"
)

func BtcHandler(w http.ResponseWriter, r *http.Request) {
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		response, err := doPost3(body)
		var content string
		if b, err := json.Marshal(response); err == nil {
			content = string(b)
		}
		fmt.Fprint(w, content)
		log_str := fmt.Sprintf("Started %s %s for %s:%s response:%s", r.Method, r.URL.Path, addr, body, content)
		log.Println(log_str)
	}

}
func doPost3(data []byte) (*json.RawMessage, error) {
	url := "http://" + cfg.Btc.Host + fmt.Sprintf(":%d", cfg.Btc.Port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	// req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "text/plain")
	// req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(cfg.Btc.Username, cfg.Btc.Password)
	client := &http.Client{
		// Timeout: timeoutIntv,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp json.RawMessage
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &rpcResp, err
}
