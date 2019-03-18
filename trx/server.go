package trx

import (
	"bytes"
	"encoding/hex"
	// "crypto/rand"
	// "encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sasaxie/go-client-api/common/base58"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	// "strconv"
	"github.com/sasaxie/go-client-api/common/crypto"
	// "github.com/sasaxie/go-client-api/core"
	// "github.com/sasaxie/go-client-api/util"
	"crypto/ecdsa"
	"crypto/sha256"
	// "github.com/golang/protobuf/proto"
	// "github.com/ethereum/go-ethereum/crypto/sha3"
	// "github.com/sasaxie/go-client-api/common/global"
	// "github.com/sasaxie/go-client-api/service"
	"sync"
	"time"
	// "unicode/utf8"
)

type RPCClient struct {
	sync.RWMutex
	Url         string
	LocalUrl    string
	Name        string
	sick        bool
	sickRate    int
	successRate int
	client      *http.Client
}

type JSONRpcResp struct {
	Id     *json.RawMessage       `json:"id"`
	Result *json.RawMessage       `json:"result"`
	Error  map[string]interface{} `json:"error,omitempty"`
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}
func NewRPCClient(name, url, localurl, timeout string) *RPCClient {
	rpcClient := &RPCClient{Name: name, Url: url, LocalUrl: localurl}
	timeoutIntv := MustParseDuration(timeout)

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
		},
	}
	return rpcClient
}

func (r *RPCClient) GetAccounts(addr string) (balance float64, err error) {
	a := base58.DecodeCheck(addr)
	if len(a) == 0 {
		err = errors.New("address format error")
		return
	}
	content := `{"address": "` + hex.EncodeToString(a) + `"}`
	rpcResp, err := r.doPost4(r.Url+"/wallet/getaccount", content)
	if err != nil {
		return
	}
	// {"account_name":"54567a7152773468365362573371484566486b38384b70363845373442634348435a","address":"41dbb157467c1b206494fc977c008c25f5901dfe4a","balance":95000000000000000,"allowance":5472000000,"is_witness":true,"account_resource":{}}
	if rpcResp == nil {
		err = errors.New("ret is nil")
		return
	}
	ret := rpcResp.(map[string]interface{})
	bla := ret["balance"]
	if bla == nil {
		// err = errors.New("balance is nil")
		log.Println("balance err:", err)
		return
	}
	balance = bla.(float64)

	// var reply core.Account
	// err = json.Unmarshal(*rpcResp, &reply)
	// balance = reply.Balance
	bytes, err := json.Marshal(rpcResp)
	log.Println("GetAccounts:", string(bytes))
	return

}
func (r *RPCClient) writeFile(dir, password, filename, privateKey string) (err error) {
	key := []byte(password)
	result, err := AesEncrypt([]byte(privateKey), key)
	if err != nil {
		return
	}
	// fmt.Println(base64.StdEncoding.EncodeToString(result))
	err = ioutil.WriteFile(dir+"/"+filename, result, 0644)

	if err != nil {
		log.Println("write file err:", err)
		return err
	}
	return nil
}
func (r *RPCClient) GetPrivate(dir, password, address string) (priva string, err error) {
	a := base58.DecodeCheck(address)
	if len(a) == 0 {
		err = errors.New("address format error")
		return
	}
	// log.Println("11")
	addr := fmt.Sprintf("%x", a)
	b, err := ioutil.ReadFile(dir + "/" + addr)
	if err != nil {
		// fmt.Print(err)
		log.Println("read file:", err)
		return
	}
	// log.Println("12")

	origData, err := AesDecrypt(b, []byte(password))
	if err != nil {
		log.Println("AesDecrypt:", err)
		return
	}
	// log.Println("13")
	if len(origData) == 0 {
		err = errors.New("decrpt is null")
		return
	}
	// log.Println("14")
	priva = string(origData)
	// log.Println("15")
	return
}
func (r *RPCClient) GenerateAddressapi(dir, password string) (addr string, err error) {
	// curl -X POST  http://127.0.0.1:8090/wallet/generateaddress
	// {"privateKey":"d033b04225e125ff884efbef3a12d5d8eb4d4cc154a27d643b3b97e6d76fca38","address":"TUoXEKdRkWG4v2U8u4is2B6AR4DyiUupXS","hexAddress":"41ce957261b1697ae095f3482986945806b7d35b3f"}
	// type addrStruct struct {
	// 	PrivateKey string `json:"privateKey"`
	// 	Address    string `json:"address"`
	// 	HexAddress string `json:"hexAddress"`
	// }

	rpcResp, err := r.doPost(r.Url+"/wallet/generateaddress", "")
	if err != nil {
		return
	}
	// var reply addrStruct
	// err = json.Unmarshal([]byte(rpcResp.(map[string]interface{})), &reply)
	// if err != nil {
	// 	return
	// }
	ret := rpcResp.(map[string]interface{})
	addr = ret["address"].(string)
	hexAddress := ret["hexAddress"].(string)
	prikey := ret["privateKey"].(string)
	err = r.writeFile(dir, password, hexAddress, prikey)
	log.Printf("password:%s addr:%s hexAddress:%s prikey:%s", password, addr, hexAddress, prikey)
	// base58.DecodeCheck(address)
	return
}
func (r *RPCClient) GenerateAddress(dir, password, net string) (addr string, err error) {
	prikey, pub, hexAddress, addr, err := r.generateAddress(net)

	err = r.writeFile(dir, password, hexAddress, prikey)
	log.Printf("password:%s addr:%s hexAddress:%s prikey:%s pub:%s", password, addr, hexAddress, prikey, pub)
	// base58.DecodeCheck(address)
	return
}

func (r *RPCClient) GetTransaction(hash string) (ret ValueStruct, err error) {
	// curl -X POST  http://127.0.0.1:8090/wallet/gettransactionbyid -d '{"value": "3d878c5bd911f8499b72812b71eef2eb70027c810a9442751040d3a935354460"}'
	// {"ret":[{"contractRet":"SUCCESS"}],"signature":["5b57db0af9b43fabb14de13341e0f9023ade47c42e1fe24c93b582e73ad1e0c55fba148eb1027498a7065ac60c731e5eb0463ecb20508f0372132c268195702a01"],"txID":"3d878c5bd911f8499b72812b71eef2eb70027c810a9442751040d3a935354460","raw_data":{"contract":[{"parameter":{"value":{"amount":1050000000,"owner_address":"41eba3361c2e8772827c50fbaa40c736601d75b442","to_address":"4165bf23e2fee14c0208e3047e1f8ec40e3fe94e9d"},"type_url":"type.googleapis.com/protocol.TransferContract"},"type":"TransferContract"}],"ref_block_bytes":"4973","ref_block_hash":"e52c84c8d42e86f1","expiration":1540483794000}}
	content := `{"value": "` + hash + `"}`
	rets, err := r.doPost3(r.Url+"/wallet/gettransactionbyid", content)
	if err != nil {
		return
	}
	log.Println("transaction:", string(*rets), ":", err)
	var reply GetTransactionResp
	err = json.Unmarshal(*rets, &reply)
	if err != nil {
		log.Println(err)
		return
	}
	// ret = reply.String()
	// if reply.Raw_data == nil {
	// 	err = errors.New("contract RawData is null")
	// 	return
	// }
	con := reply.Raw_data.Contract
	if len(con) == 0 {
		err = errors.New("contract is null")
		return
	}
	log.Println(len(con))
	v := con[0]
	// for _, v := range con {
	log.Println(v)
	ret = v.Parameter.Value
	// }
	// con0 := con[0]

	// ret = con0.Parameter.Value
	return
}

// func (r *RPCClient) Send(from, to, amount, password, dir, privatenotindir string) (result BroadcastTransactionReturn, err error) {
// 	client := service.NewGrpcClient()
// 	client.Start()
// 	defer client.Conn.Close()

// 	key, err := crypto.GetPrivateKeyByHexString(*ownerPrivateKey)

// 	if err != nil {
// 		log.Fatalf("get private key by hex string error: %v", err)
// 	}

// 	result := client.Transfer(key, *address, *amount)

// 	fmt.Printf("result: %v\n", result)

// }

func (r *RPCClient) Send(from, to, amount, password, dir, privatenotindir string) (result BroadcastTransactionReturn, err error) {
	f := base58.DecodeCheck(from)
	if len(f) == 0 {
		err = errors.New("address format error")
		return
	}
	t := base58.DecodeCheck(to)
	if len(t) == 0 {
		err = errors.New("address format error")
		return
	}
	ret, err := r.createTransaction(hex.EncodeToString(f), hex.EncodeToString(t), amount)
	if err != nil {
		return
	}
	// fmt.Println(string(ret))
	j, err := json.Marshal(ret)
	if err != nil {
		// panic(err)
		log.Println("Marshal:", err)
		return
	}
	log.Println("createTransaction:", string(j))
	////本地对交易签名
	// SignTransaction(transaction *core.Transaction, key *ecdsa.PrivateKey)

	var tr Transaction
	// var tr core.Transaction
	err = json.Unmarshal(j, &tr)
	if err != nil {
		log.Println("unmarshal:", err)
		return
	}
	if tr.Error != "" {
		err = errors.New(tr.Error)
		return
	}
	log.Println("205 start:")
	key, err := r.GetPrivate(dir, password, from)
	if err != nil {
		key = privatenotindir
		// return
		log.Println("GetPrivate:", err)
		return
	}
	log.Println("sign privatekey:", key)
	// signedtr, err := r.SignTransaction(&tr, key)
	signedtr, err := r.SignTransactionFromApiLocal(&tr, key)
	// signedtr, err := r.SignTransactionUseCode(&tr, key)
	if err != nil {
		log.Println("SignTransaction:", err)
		return
	}
	// log.Println("broadcastTransaction start:")
	reply, err := r.broadcastTransaction(signedtr)
	if err != nil {
		log.Println("broadcastTransaction:", err)
		return
	}
	out, err := reply.MarshalJSON()
	log.Println("broadcastTransaction reply:", string(out), ":", err)
	// var temp BroadcastTransactionReturn
	err = json.Unmarshal(*reply, &result)
	if err != nil {

		return
	}
	result.TransactionID = signedtr.TxID
	log.Println("ret:", result)
	if result.Error != "" {
		log.Println("ret error:", result.Error)
		err = errors.New(result.Error + ", probably password is error")
	}
	// log.Println("message:", string(result.Message))

	// ret: {"Error":"class org.tron.core.services.http.JsonFormat$ParseException : 1:15: INVALID hex String"}
	return
}
func (r *RPCClient) SignTransactionFromApi(transaction *Transaction, key string) (ret Transaction, err error) {
	//////签名的api调用
	// 	/wallet/gettransactionsign
	// 作用：对交易签名，该api有泄漏private key的风险，请确保在安全的环境中调用该api
	// demo: curl -X POST  http://127.0.0.1:8090/wallet/gettransactionsign -d '{
	// "transaction" : {"txID":"454f156bf1256587ff6ccdbc56e64ad0c51e4f8efea5490dcbc720ee606bc7b8","raw_data":{"contract":[{"parameter":{"value":{"amount":1000,"owner_address":"41e552f6487585c2b58bc2c9bb4492bc1f17132cd0","to_address":"41d1e7a6bc354106cb410e65ff8b181c600ff14292"},"type_url":"type.googleapis.com/protocol.TransferContract"},"type":"TransferContract"}],"ref_block_bytes":"267e","ref_block_hash":"9a447d222e8de9f2","expiration":1530893064000,"timestamp":1530893006233}}, "privateKey": "your private key"
	// }'
	rawData, err := json.Marshal(transaction)

	if err != nil {
		log.Printf("sign transaction error: %v", err)
		return
	}
	content := `{"transaction":` + string(rawData) + `, "privateKey": "` + key + `" }`
	reply, err := r.doPost3(r.Url+"/wallet/gettransactionsign", content)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(*reply, &ret)
	log.Println("ret:", string(*reply))
	return
}
func (r *RPCClient) SignTransactionFromApiLocal(transaction *Transaction, key string) (ret Transaction, err error) {

	rawData, err := json.Marshal(transaction)

	if err != nil {
		log.Printf("sign transaction error: %v", err)
		return
	}
	content := `{"transaction":` + string(rawData) + `, "privateKey": "` + key + `" }`
	reply, err := r.doPost3(r.LocalUrl+"/wallet/gettransactionsign", content)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(*reply, &ret)
	log.Println("ret:", string(*reply))
	return
}
func (r *RPCClient) SignTransactionUseCode(transaction *Transaction, key string) (ret Transaction, err error) {
	priKey, err := crypto.GetPrivateKeyByHexString(key)
	if err != nil {
		log.Println(err)
		return
	}
	err = r.signTransaction(transaction, priKey)
	if err != nil {
		log.Println(err)
		return
	}
	ret = *transaction
	return
}
func (r *RPCClient) signTransaction(transaction *Transaction, key *ecdsa.PrivateKey) (err error) {
	// SignTransaction(transaction *core.Transaction, key *ecdsa.PrivateKey) {
	// transaction.GetRawData().Timestamp = time.Now().UnixNano() / 1000000
	transaction.Raw_data.Timestamp = time.Now().UnixNano() / 1000000
	rawData, err := json.Marshal(transaction.Raw_data)
	// rawData, err := proto.Marshal(transaction.Raw_data)
	if err != nil {
		log.Printf("sign transaction error: %v", err)
		return
	}
	header := "\x19TRON Signed Message:\n"
	leng := fmt.Sprintf("%d", len(rawData))
	header += leng
	header += string(rawData)

	h256h := sha256.New()
	// h256h := sha3.NewKeccak256()
	h256h.Write([]byte(header))
	hash := h256h.Sum(nil)

	contractList := transaction.Raw_data.Contract

	for range contractList {
		signature, errs := crypto.Sign(hash, key)
		if errs != nil {
			log.Printf("sign transaction error: %v", errs)
			err = errs
			return
		}
		// transaction.Signature = append(transaction.Signature, hex.EncodeToString(signature))
		// transaction.Signature = append(transaction.Signature, hex.EncodeToString(signature))
		transaction.Signature = append(transaction.Signature, hex.EncodeToString(signature))
	}

	return

}
func (r *RPCClient) createTransaction(from, to, amount string) (reply *json.RawMessage, err error) {
	// curl -X POST  http://127.0.0.1:8090/wallet/createtransaction -d '{"to_address": "41e9d79cc47518930bc322d9bf7cddd260a0260a8d", "owner_address": "41D1E7A6BC354106CB410E65FF8B181C600FF14292", "amount": 1000 }'
	content := `{"to_address": "` + to + `", "owner_address": "` + from + `", "amount": ` + amount + ` }`
	reply, err = r.doPost3(r.Url+"/wallet/createtransaction", content)
	if err != nil {
		return
	}
	return
}
func (r *RPCClient) broadcastTransaction(t Transaction) (reply *json.RawMessage, err error) {
	// curl -X POST  http://127.0.0.1:8090/wallet/broadcasttransaction -d '{"signature":["97c825b41c77de2a8bd65b3df55cd4c0df59c307c0187e42321dcc1cc455ddba583dd9502e17cfec5945b34cad0511985a6165999092a6dec84c2bdd97e649fc01"],"txID":"454f156bf1256587ff6ccdbc56e64ad0c51e4f8efea5490dcbc720ee606bc7b8","raw_data":{"contract":[{"parameter":{"value":{"amount":1000,"owner_address":"41e552f6487585c2b58bc2c9bb4492bc1f17132cd0","to_address":"41d1e7a6bc354106cb410e65ff8b181c600ff14292"},"type_url":"type.googleapis.com/protocol.TransferContract"},"type":"TransferContract"}],"ref_block_bytes":"267e","ref_block_hash":"9a447d222e8de9f2","expiration":1530893064000,"timestamp":1530893006233}}'
	content, err := json.Marshal(t)
	if err != nil {
		return
	}
	log.Println("broadcastTransaction")
	// content := `{"signature":["` + sign + `"],"txID":"` + txid + `","raw_data":` + raw + `}`
	reply, err = r.doPost3(r.Url+"/wallet/broadcasttransaction", string(content))
	if err != nil {
		return
	}
	return
}

func (r *RPCClient) doPost(url string, content string) (interface{}, error) {
	// jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	// data, _ := json.Marshal(jsonReq)
	log.Println("content:", content)
	data := []byte(content)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	// req.Header.Set("Content-Length", (string)(len(data)))
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		r.markSick()
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp interface{}
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		r.markSick()
		return nil, err
	}
	return rpcResp, err
}
func (r *RPCClient) doPost2(url string, content string) (*JSONRpcResp, error) {
	// jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	// data, _ := json.Marshal(jsonReq)
	log.Println("content 2:", content)
	data := []byte(content)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		r.markSick()
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp *JSONRpcResp
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		r.markSick()
		return nil, err
	}
	if rpcResp.Error != nil {
		r.markSick()
		return nil, errors.New(rpcResp.Error["message"].(string))
	}
	return rpcResp, err
}
func (r *RPCClient) doPost3(url string, content string) (*json.RawMessage, error) {
	// jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	// data, _ := json.Marshal(jsonReq)
	log.Println("content 3:", content)
	data := []byte(content)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
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
func (r *RPCClient) doPost4(url string, content string) (interface{}, error) {
	// jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	// data, _ := json.Marshal(jsonReq)
	log.Println("content 4:", content)
	data := []byte(content)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		r.markSick()
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp interface{}
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		r.markSick()
		return nil, err
	}
	return rpcResp, err
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
