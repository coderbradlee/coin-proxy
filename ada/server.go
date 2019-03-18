package ada

import (
	// "github.com/golang/tools/godoc/util"
	"bytes"
	// "github.com/rubblelabs/ripple/data"
	// "crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	// "math/big"
	"net"
	"net/http"
	// "strconv"
	"time"
	// "strings"
	// utils "../util"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"io/ioutil"
	"log"
	"sync"
)

type RPCClient struct {
	sync.RWMutex
	Url         string
	Name        string
	sick        bool
	sickRate    int
	successRate int
	client      *http.Client
}

type GetBlockReply struct {
	Number       string   `json:"number"`
	Hash         string   `json:"hash"`
	Nonce        string   `json:"nonce"`
	Miner        string   `json:"miner"`
	Difficulty   string   `json:"difficulty"`
	GasLimit     string   `json:"gasLimit"`
	GasUsed      string   `json:"gasUsed"`
	Transactions []Tx     `json:"transactions"`
	Uncles       []string `json:"uncles"`
	// https://github.com/ethereum/EIPs/issues/95
	SealFields []string `json:"sealFields"`
}

type GetBlockReplyPart struct {
	Number     string `json:"number"`
	Difficulty string `json:"difficulty"`
}

type Tx struct {
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Hash     string `json:"hash"`
}

type JSONRpcResp struct {
	Id     *json.RawMessage       `json:"id"`
	Result *json.RawMessage       `json:"result"`
	Error  map[string]interface{} `json:"error,omitempty"`
}

func NewRPCClient(name, url, timeout, path string) *RPCClient {
	rpcClient := &RPCClient{Name: name, Url: url}
	timeoutIntv := MustParseDuration(timeout)

	pool := x509.NewCertPool()
	caCertPath := path + "/ca.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Println("ReadFile err:", err)
		return nil
	}
	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair(path+"/client.crt", path+"/client.key")
	if err != nil {
		log.Println("Loadx509keypair err:", err)
		return nil
	}

	rpcClient.client = &http.Client{
		Timeout: timeoutIntv,
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
			// TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
			},
			// TLSClientConfig: &tls.Config{RootCAs: loadCA("server.crt")}
		},
	}
	return rpcClient
}

func (r *RPCClient) GetInfo() (ret interface{}, err error) {
	// url https://127.0.0.1:8090/api/v1/node-info
	url := r.Url + "/api/v1/node-info"
	ret, err = r.get(url)

	return
}
func (c *RPCClient) getNewAddress(accountIndex, walletId, spendingPassword string) (addr string, err error) {
	// curl -X POST https://localhost:8090/api/v1/addresses \
	// -H "Accept: application/json; charset=utf-8" \
	//   -H "Content-Type: application/json; charset=utf-8" \
	//  --cacert /root/ada/cardano-sl/state-wallet-mainnet/tls/client/ca.crt --cert /root/ada/cardano-sl/state-wallet-mainnet/tls/client/client.crt --key /root/ada/cardano-sl/state-wallet-mainnet/tls/client/client.key  \
	//   -d '{
	//   "accountIndex": 2149473377,
	//   "walletId": "Ae2tdPwUPEZK3qR3RkcPfdeYAVWw3EsM7SZTX3P6h9BtQCM5NdFP47poz34",
	//   "spendingPassword": "6383ba0ac13d92957745725008907acaf4613c910764f04959030e1b81c603b9"
	// }'
	content := fmt.Sprintf(`{"accountIndex":%s,"walletId":"%s","spendingPassword":"%s"}`, accountIndex, walletId, spendingPassword)
	url := c.Url + "/api/v1/addresses"
	ret, err := c.doPost(url, content)
	if err != nil {
		return
	}
	log.Println("getnewaddress:", ret)
	if ret.Status != "success" {
		err = errors.New(ret.Diagnostic.Msg)
		return
	}
	dd := ret.Data.(map[string]interface{})
	if data, ok := dd["id"]; ok {
		addrs, stringok := data.(string)
		if !stringok {
			err = errors.New("id convert error")
			return
		}
		addr = addrs
	}
	return
}
func (s *RPCClient) checkPass(addr, pass, path string) error {
	b, err := ioutil.ReadFile(path + "/" + addr)
	if err != nil {
		// fmt.Print(err)
		log.Println("read file:", err)
		return errors.New("open file error!")
	}
	str := string(b)
	log.Println(addr+pass, ":", str)

	hash := sha3.NewKeccak256()

	var buf []byte
	//hash.Write([]byte{0xcc})
	hash.Write([]byte(addr + pass))
	buf = hash.Sum(buf)
	if str != hex.EncodeToString(buf) {
		log.Println("password error:", str, "!=", hex.EncodeToString(buf))
		return errors.New("wrong password!")
	}
	return nil
}
func (c *RPCClient) Send(path, accountIndex, walletId, spendingPassword string, to []interface{}) (hash string, reterr error) {
	// /api/v1/transactions
	// {
	// 	"groupingPolicy": null,
	// 	"destinations": [
	// 	],
	// 	"source": {
	// 	  "accountIndex": 2,
	// 	  "walletId": "J7rQqaLLHBFPrgJXwpktaMB1B1kQBXAyc2uRSfRPzNVGiv6TdxBzkPNBUWysZZZdhFG9gRy3sQFfX5wfpLbi4XTFGFxTg"
	// 	},
	// 	"spendingPassword": "0200020002020202010000010101020100020001020201020001000101010101"
	//   }
	// log.Println("enter send:", time.Now())
	type TransferData struct {
		Amount   float64 `json:"amount,omitempty"`
		Address  string  `json:"address,omitempty"`
		Password string  `json:"password,omitempty"`
	}
	var toData []TransferData
	/////每一个都验证下密码，若密码错误则返回此地址
	for _, v := range to {
		ob, ok := v.(map[string]interface{})
		if !ok {
			reterr = errors.New("convert to map error")
			return
		}
		addr, ok := ob["address"]
		if !ok {
			reterr = errors.New("address is not exist error")
			return
		}
		amount, ok := ob["amount"]
		if !ok {
			reterr = errors.New("amount is not exist error")
			return
		}
		// password, ok := ob["password"]
		// if !ok {
		// 	reterr = errors.New("password is not exist error")
		// 	return
		// }
		addre, addreok := addr.(string)
		if !addreok {
			reterr = errors.New("address convert error")
			return
		}
		// passw, passwok := password.(string)
		// if !passwok {
		// 	reterr = errors.New("password convert error")
		// 	return
		// }

		amoun, amounok := amount.(float64)
		if !amounok {
			reterr = errors.New("amount convert error")
			return
		}

		// err := c.checkPass(addre, passw, path)
		// if err != nil {
		// 	reterr = errors.New(addre + ":" + err.Error())
		// 	return
		// }
		t := TransferData{
			Amount:  amoun,
			Address: addre,
		}
		toData = append(toData, t)

	}
	// log.Println("224:", time.Now())
	if len(toData) == 0 {
		reterr = errors.New("make data error")
		return
	}
	marshalledData, reterr := json.Marshal(toData)
	if reterr != nil {
		return
	}
	content := fmt.Sprintf(`{"groupingPolicy": "OptimizeForSecurity","destinations": %s,"source": {"accountIndex": %s,"walletId": "%s"},"spendingPassword": "%s"}`, string(marshalledData), accountIndex, walletId, spendingPassword)
	url := c.Url + "/api/v1/transactions"
	// log.Println("255:", time.Now())
	ret, reterr := c.doPost(url, content)
	// log.Println("264:", time.Now())
	if reterr != nil {
		return
	}
	// log.Println("transaction hash response:", ret)
	if ret.Status != "success" {
		reterr = errors.New(ret.Diagnostic.Msg)
		return
	}

	dd := ret.Data.(map[string]interface{})
	if data, ok := dd["id"]; ok {
		hashs, stringok := data.(string)
		if !stringok {
			reterr = errors.New("id convert error")
			return
		}
		hash = hashs
	}
	return
}
func (c *RPCClient) VerifyIn(accountIndex, walletId, address, page, perpage string) (results []GetTransactionResponse, reterr error) {
	url := fmt.Sprintf(`%s/api/v1/transactions?wallet_id=%s&accountIndex=%s&page=%s&per_page=%s&address=%s&sort_by=created_at`, c.Url, walletId, accountIndex, page, perpage, address)

	ret, err := c.get2(url)
	log.Println("VerifyIn:", string(ret))
	var tranData Transaction
	err = json.Unmarshal(ret, &tranData)
	if err != nil {
		return
	}
	if tranData.Status != "success" {
		err = errors.New(tranData.Diagnostic.Msg)
		return
	}
	for _, datas := range tranData.Data {
		var result GetTransactionResponse
		result.Hash = datas.Id
		result.CreationTime = datas.CreationTime
		for _, v := range datas.Outputs {
			// log.Println("addr:", v.Address)
			if v.Address == address {
				result.Amount = fmt.Sprintf("%d", v.Amount)
				result.Address = address
				// results = append(results, result)
				results = append([]GetTransactionResponse{result}, results...)
			}
		}
	}
	if len(results) == 0 {
		err = errors.New("no transaction found")
	}
	return
}

func (c *RPCClient) GetTransaction(accountIndex, walletId, id string) (results GetTransactionRet, reterr error) {
	url := fmt.Sprintf(`%s/api/v1/transactions?wallet_id=%s&accountIndex=%s&id=%s`, c.Url, walletId, accountIndex, id)
	// url := r.Url + "/api/v1/transactions?wallet_id=" + walletid&
	ret, err := c.get2(url)
	log.Println("GetTrasaction:", string(ret))
	var tranData Transaction
	err = json.Unmarshal(ret, &tranData)
	if err != nil {
		return
	}
	if tranData.Status != "success" {
		err = errors.New(tranData.Diagnostic.Msg)
		return
	}
	if len(tranData.Data) == 0 {
		err = errors.New("no transaction found")
		return
	}
	results.Id = id
	var in uint64
	in = 0
	var out uint64
	out = 0
	for _, v := range tranData.Data {
		for _, iv := range v.Inputs {
			in += iv.Amount
		}
		for _, iv := range v.Outputs {
			out += iv.Amount
		}
	}
	results.Fee = in - out
	return
}

func (r *RPCClient) GetBalance(walletid, addr string) (balance string, err error) {
	url := r.Url + "/api/v1/transactions?wallet_id=" + walletid
	ret, err := r.get2(url)
	log.Println("GetBalance:", string(ret))
	var tranData Transaction
	err = json.Unmarshal(ret, &tranData)
	if err != nil {
		return
	}
	if tranData.Status != "success" {
		err = errors.New(tranData.Diagnostic.Msg)
		return
	}
	for _, datas := range tranData.Data {
		for _, v := range datas.Outputs {
			log.Println("addr:", v.Address)
			if v.Address == addr {
				balance = fmt.Sprintf("%d", v.Amount)
				return
			}
		}
	}
	err = errors.New("no this address")
	return
}
func (r *RPCClient) GetBalance2(walletid string) (balance uint64, err error) {
	url := r.Url + "/api/v1/wallets/" + walletid
	ret, err := r.get2(url)
	log.Println("GetBalance of wallets:", string(ret))
	var tranData BalanceResponse
	err = json.Unmarshal(ret, &tranData)
	if err != nil {
		return
	}
	if tranData.Status != "success" {
		err = errors.New(tranData.Diagnostic.Msg)
		return
	}
	// for _, datas := range tranData.Data {
	// 	for _, v := range datas.Outputs {
	// 		log.Println("addr:", v.Address)
	// 		if v.Address == addr {
	// 			balance = fmt.Sprintf("%d", v.Amount)
	// 			return
	// 		}
	// 	}
	// }
	// err = errors.New("no this address")
	balance = tranData.Data.Balance
	return
}
func (c *RPCClient) GetNewAddress(pass, accountIndex, walletId, spendingPassword, path string) (addr string, err error) {
	//成功后保存文件，地址为文件名，内容为地址+密码的hash值
	hash := sha3.NewKeccak256()
	addr, err = c.getNewAddress(accountIndex, walletId, spendingPassword)
	if err != nil {
		return
	}
	var buf []byte
	//hash.Write([]byte{0xcc})
	hash.Write([]byte(addr + pass))
	buf = hash.Sum(buf)

	// fmt.Println(hex.EncodeToString(buf))
	// d1 := []byte("hello\ngo\n")
	err = ioutil.WriteFile(path+addr, []byte(hex.EncodeToString(buf)), 0644)

	if err != nil {
		log.Println("write file err:", err)
		return
	}
	return
}

// func (r *RPCClient) CreateWallet() (ret interface{}, err error) {
// 	// url /api/v1/wallets
// 	url := r.Url + "/api/v1/wallets"
// 	param := `{
// 		"operation": "create",
// 		"backupPhrase": ["squirrel", "material", "silly", "twice", "direct", "slush", "pistol", "razor", "become", "junk", "kingdom", "flee"],
// 		"assuranceLevel": "normal",
// 		"name": "MyFirstWallet",
// 		"spendingPassword": "5416b2988745725998907addf4613c9b0764f04959030e1b81c603b920a115d0"
// 	  }`
// 	ret, err = r.get(url)

// 	return
// }
func (r *RPCClient) get(url string) (ret interface{}, err error) {

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Close = true

	resp, err := r.client.Do(req)
	if err != nil {
		r.markSick()
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		r.markSick()
		return
	}
	return
}
func (r *RPCClient) get2(url string) (ret json.RawMessage, err error) {
	log.Println("get2 url:", url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Close = true

	resp, err := r.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return
	}
	return
}
func (r *RPCClient) doPost(url, content string) (response ResponseStruct, err error) {
	log.Println("437:", time.Now())
	log.Println(url, " doPost content:", content)
	data := []byte(content)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Close = true
	log.Println("445:", time.Now())
	resp, err := r.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	log.Println("451:", time.Now())
	// var rpcResp *JSONRpcResp
	err = json.NewDecoder(resp.Body).Decode(&response)
	// if err != nil {
	// 	return
	// }
	log.Println("457:", time.Now())
	log.Println("doPost response:", response)
	log.Println("459:", time.Now())
	return
}

func (r *RPCClient) markSick() {
	r.Lock()
	r.sickRate++
	r.successRate = 0
	if r.sickRate >= 5 {
		r.sick = true
	}
	r.Unlock()
}
