package etc

import (
	"bytes"
	// "fmt"
	// "crypto/sha256"
	"encoding/json"
	"errors"
	// "fmt"
	// "github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"net/http"
	// "strconv"
	// "strings"
	"sync"
	"time"
	//"github.com/ethereumproject/go-ethereum/common"
	// "91pool/util"
	"../aes"
	"io/ioutil"
	"regexp"
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

//dao
type GetBlockReply2 struct {
	Number           string   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            string   `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	Difficulty       string   `json:"difficulty"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	Size             string   `json:"size"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Timestamp        string   `json:"timestamp"`
	Transactions     []Tx2    `json:"transactions"`
	Uncles           []string `json:"uncles"`
	// https://github.com/ethereum/EIPs/issues/95
	SealFields []string `json:"sealFields"`
}

type GetBlockReply struct {
	Number           string   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            string   `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	Difficulty       string   `json:"difficulty"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	Size             string   `json:"size"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Timestamp        string   `json:"timestamp"`
	Transactions     []Tx     `json:"transactions"`
	Uncles           []string `json:"uncles"`
	// https://github.com/ethereum/EIPs/issues/95
	SealFields []string `json:"sealFields"`
}

type GetBlockReplyPart struct {
	Number     string `json:"number"`
	Difficulty string `json:"difficulty"`
}

type TxReceipt struct {
	TxHash  string `json:"transactionHash"`
	GasUsed string `json:"gasUsed"`
}

type Tx struct {
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Hash     string `json:"hash"`
}

//dao
type Tx2 struct {
	Hash string `json:"hash"` //"hash":"0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b",
	// "nonce":"0x",
	BlockHash        string `json:"blockHash"`        // "blockHash": "0xbeab0aa2411b7ab17f30a99d3cb9c6ef2fc5426d6ad6fd9e2a26a6aed1d1055b",
	BlockNumber      string `json:"blockNumber"`      // "blockNumber": "0x15df", // 5599
	TransactionIndex string `json:"transactionIndex"` // "transactionIndex":  "0x1", // 1  "transactionIndex":  "0x1", // 1
	From             string `json:"from"`             //  "from":"0x407d73d8a49eeb85d32cf465507dd71d507100c1",
	To               string `json:"to"`               //  "to":"0x85h43d8a49eeb85d32cf465507dd71d507100c1",
	Value            string `json:"value"`            // "value":"0x7f110" // 520464
	Gas              string `json:"gas"`              // "0x7f110" // 520464
	GasPrice         string `json:"gasPrice"`         //"0x09184e72a000",
	// "input":"0x603880600c6000396000f300603880600c6000396000f3603880600c6000396000f360",
}

type JSONRpcResp struct {
	Id     *json.RawMessage       `json:"id"`
	Result *json.RawMessage       `json:"result"`
	Error  map[string]interface{} `json:"error"`
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}
func NewRPCClient(name, url, timeout string) *RPCClient {
	rpcClient := &RPCClient{Name: name, Url: url}
	timeoutIntv := MustParseDuration(timeout)
	rpcClient.client = &http.Client{
		Timeout: timeoutIntv,
	}
	return rpcClient
}

func (r *RPCClient) GetWork() ([]string, error) {
	rpcResp, err := r.doPost(r.Url, "eth_getWork", []string{})
	if err != nil {
		return nil, err
	}
	var reply []string
	err = json.Unmarshal(*rpcResp.Result, &reply)
	return reply, err
}
func (r *RPCClient) GetNewAddress(userpass, realpass, dir string) (addr string, err error) {
	rpcResp, err := r.doPost(r.Url, "personal_newAccount", []string{realpass})
	if err != nil {
		return
	}
	err = json.Unmarshal(*rpcResp.Result, &addr)
	//以产生的地址为文件名，保存userpass的密文，此密文以realpass为密码用aes加密
	if err != nil {
		return
	}
	err = r.writeFile(dir, realpass, addr, userpass)
	return
}
func (r *RPCClient) writeFile(dir, password, filename, content string) (err error) {
	key := []byte(password)
	result, err := aes.AesEncrypt([]byte(content), key)
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

// func (r *RPCClient) GetBlockNumber() (int64, error) {
// 	rpcResp, err := r.doPost(r.Url, "eth_blockNumber", []string{})
// 	if err != nil {
// 		return int64(0), err
// 	}
// 	var reply string
// 	err = json.Unmarshal(*rpcResp.Result, &reply)
// 	if err != nil {
// 		return int64(0), err
// 	}
// 	height, err := strconv.ParseInt(reply, 0, 64)
// 	return int64(height), err
// }

// func (r *RPCClient) GetPendingBlock() (*GetBlockReplyPart, error) {
// 	rpcResp, err := r.doPost(r.Url, "eth_getBlockByNumber", []interface{}{"pending", false})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if rpcResp.Result != nil {
// 		var reply *GetBlockReplyPart
// 		err = json.Unmarshal(*rpcResp.Result, &reply)
// 		return reply, err
// 	}
// 	return nil, nil
// }

// //dao
// func (r *RPCClient) GetBlockByHeight2(height int64) (*GetBlockReply2, error) {
// 	params := []interface{}{fmt.Sprintf("0x%x", height), true}
// 	return r.getBlockBy2("eth_getBlockByNumber", params)
// }

//dao
// func (r *RPCClient) getBlockBy2(method string, params []interface{}) (*GetBlockReply2, error) {
// 	rpcResp, err := r.doPost(r.Url, method, params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if rpcResp.Result != nil {
// 		var reply *GetBlockReply2
// 		err = json.Unmarshal(*rpcResp.Result, &reply)
// 		return reply, err
// 	}
// 	return nil, nil
// }

// func (r *RPCClient) GetBlockByHeight(height int64) (*GetBlockReply, error) {
// 	params := []interface{}{fmt.Sprintf("0x%x", height), true}
// 	return r.getBlockBy("eth_getBlockByNumber", params)
// }

// func (r *RPCClient) GetBlockByHash(hash string) (*GetBlockReply, error) {
// 	params := []interface{}{hash, true}
// 	return r.getBlockBy("eth_getBlockByHash", params)
// }

// func (r *RPCClient) GetUncleByBlockNumberAndIndex(height int64, index int) (*GetBlockReply, error) {
// 	params := []interface{}{fmt.Sprintf("0x%x", height), fmt.Sprintf("0x%x", index)}
// 	return r.getBlockBy("eth_getUncleByBlockNumberAndIndex", params)
// }

// func (r *RPCClient) getBlockBy(method string, params []interface{}) (*GetBlockReply, error) {
// 	rpcResp, err := r.doPost(r.Url, method, params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if rpcResp.Result != nil {
// 		var reply *GetBlockReply
// 		err = json.Unmarshal(*rpcResp.Result, &reply)
// 		return reply, err
// 	}
// 	return nil, nil
// }

// func (r *RPCClient) GetTxReceipt(hash string) (*TxReceipt, error) {
// 	rpcResp, err := r.doPost(r.Url, "eth_getTransactionReceipt", []string{hash})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if rpcResp.Result != nil {
// 		var reply *TxReceipt
// 		err = json.Unmarshal(*rpcResp.Result, &reply)
// 		return reply, err
// 	}
// 	return nil, nil
// }

//dao
func (r *RPCClient) GetTransaction(hash string) (*Tx2, error) {
	if hash[:2] != "0x" {
		hash = "0x" + hash
	}
	rpcResp, err := r.doPost(r.Url, "eth_getTransactionByHash", []string{hash})
	if err != nil {
		return nil, err
	}
	if rpcResp.Result != nil {
		var reply *Tx2
		err = json.Unmarshal(*rpcResp.Result, &reply)
		return reply, err
	}
	return nil, nil
}

// func (r *RPCClient) SubmitBlock(params []string) (bool, error) {
// 	rpcResp, err := r.doPost(r.Url, "eth_submitWork", params)
// 	if err != nil {
// 		return false, err
// 	}
// 	var reply bool
// 	err = json.Unmarshal(*rpcResp.Result, &reply)
// 	return reply, err
// }

//锁定帐号
// func (r *RPCClient) LockAccount(address string) (bool, error) {
// 	rpcResp, err := r.doPost(r.Url, "personal_lockAccount", []string{address})
// 	if err != nil {
// 		return false, err
// 	}
// 	var reply bool
// 	err = json.Unmarshal(*rpcResp.Result, &reply)
// 	return reply, err
// }

// func (r *RPCClient) UnlockAccount(params []interface{}) (bool, error) {
// 	rpcResp, err := r.doPost(r.Url, "personal_unlockAccount", params)
// 	if err != nil {
// 		return false, err
// 	}
// 	var reply bool
// 	err = json.Unmarshal(*rpcResp.Result, &reply)
// 	return reply, err
// }

func (r *RPCClient) GetBalance(address string) (reply string, err error) {
	if address[0:2] != "0x" {
		address = "0x" + address
	}
	rpcResp, err := r.doPost(r.Url, "eth_getBalance", []string{address, "latest"})
	if err != nil {
		return
	}
	var ret string
	err = json.Unmarshal(*rpcResp.Result, &ret)
	if err != nil {
		return
	}
	m := new(big.Int)
	n, ok := m.SetString(ret[2:], 16)
	if !ok {
		err = errors.New("to big int error")
		return
	}
	reply = n.Text(10)
	return
}

// func (r *RPCClient) Sign(from string, s string) (string, error) {
// 	hash := sha256.Sum256([]byte(s))
// 	rpcResp, err := r.doPost(r.Url, "eth_sign", []string{from, common.ToHex(hash[:])})
// 	var reply string
// 	if err != nil {
// 		return reply, err
// 	}
// 	err = json.Unmarshal(*rpcResp.Result, &reply)
// 	if err != nil {
// 		return reply, err
// 	}
// 	if util.IsZeroHash(reply) {
// 		err = errors.New("Can't sign message, perhaps account is locked")
// 	}
// 	return reply, err
// }

// func (r *RPCClient) GetPeerCount() (int64, error) {
// 	rpcResp, err := r.doPost(r.Url, "net_peerCount", nil)
// 	if err != nil {
// 		return 0, err
// 	}
// 	var reply string
// 	err = json.Unmarshal(*rpcResp.Result, &reply)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return strconv.ParseInt(strings.Replace(reply, "0x", "", -1), 16, 64)
// }

// func (r *RPCClient) SendTransaction(from, to, gas, gasPrice, value string, autoGas bool) (string, error) {

// 	params := map[string]string{
// 		"from":  from,
// 		"to":    to,
// 		"value": value,
// 	}
// 	if !autoGas {
// 		params["gas"] = gas
// 		params["gasPrice"] = gasPrice
// 	}

// 	rpcResp, err := r.doPost(r.Url, "eth_sendTransaction", []interface{}{params})
// 	var reply string
// 	if err != nil {
// 		return reply, err
// 	}
// 	err = json.Unmarshal(*rpcResp.Result, &reply)
// 	if err != nil {
// 		return reply, err
// 	}
// 	/* There is an inconsistence in a "standard". Geth returns error if it can't unlock signer account,
// 	 * but Parity returns zero hash 0x000... if it can't send tx, so we must handle this case.
// 	 * https://github.com/ethereum/wiki/wiki/JSON-RPC#returns-22
// 	 */
// 	if util.IsZeroHash(reply) {
// 		err = errors.New("transaction is not yet available")
// 	}
// 	return reply, err
// }

func (r *RPCClient) SendTransactionParity(from, to, gas, gasPrice, value string, autoGas bool, pwd, dir, realpass string) (hash string, err error) {
	b, err := ioutil.ReadFile(dir + "/" + from)
	if err != nil {
		// fmt.Print(err)
		log.Println("read file:", err)
		return
	}
	origData, err := aes.AesDecrypt(b, []byte(realpass))
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
	if pwd != string(origData) {
		log.Println(string(origData))
		err = errors.New("wrong password")
		return
	}
	m := new(big.Int)
	n, ok := m.SetString(value, 10)
	if !ok {
		err = errors.New("to big int error")
		return
	}
	values := "0x" + n.Text(16)
	params := map[string]string{
		"from":  from,
		"to":    to,
		"value": values,
	}
	if !autoGas {
		params["gas"] = gas
		params["gasPrice"] = gasPrice
	}

	rpcResp, err := r.doPost(r.Url, "personal_sendTransaction", []interface{}{params, realpass})
	var reply string
	if err != nil {
		return reply, err
	}
	err = json.Unmarshal(*rpcResp.Result, &reply)
	if err != nil {
		return reply, err
	}
	// 	/* There is an inconsistence in a "standard". Geth returns error if it can't unlock signer account,
	// 	 * but Parity returns zero hash 0x000... if it can't send tx, so we must handle this case.
	// 	 * https://github.com/ethereum/wiki/wiki/JSON-RPC#returns-22
	// 	 */
	if IsZeroHash(reply) {
		err = errors.New("transaction is not yet available")
	}
	return reply, err
}
func IsZeroHash(s string) bool {
	var zeroHash = regexp.MustCompile("^0?x?0+$")
	return zeroHash.MatchString(s)
}
func (r *RPCClient) doPost(url string, method string, params interface{}) (*JSONRpcResp, error) {
	log.Println("post content:", url, ":", method, ":", params)
	jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	data, _ := json.Marshal(jsonReq)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	//log.Println("++++++",req.Body)
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

func (r *RPCClient) Check() bool {
	_, err := r.GetWork()
	if err != nil {
		return false
	}
	r.markAlive()
	return !r.Sick()
}

func (r *RPCClient) Sick() bool {
	r.RLock()
	defer r.RUnlock()
	return r.sick
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

func (r *RPCClient) markAlive() {
	r.Lock()
	r.successRate++
	if r.successRate >= 5 {
		r.sick = false
		r.sickRate = 0
		r.successRate = 0
	}
	r.Unlock()
}
