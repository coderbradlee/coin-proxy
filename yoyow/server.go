package yoyow

import (
	"bytes"
	"fmt"
	// "crypto/sha256"
	"encoding/json"
	"errors"
	// "fmt"
	// "math/big"
	"net"
	"net/http"
	"strconv"
	"time"
	// "strings"
	"sync"
	// "github.com/ethereum/go-ethereum/common"
	// "util"
	"crypto/rand"
	"encoding/binary"
	// "reflect"
	"log"
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

const receiptStatusSuccessful = "0x1"

type TxReceipt struct {
	TxHash    string `json:"transactionHash"`
	GasUsed   string `json:"gasUsed"`
	BlockHash string `json:"blockHash"`
	Status    string `json:"status"`
}

func (r *TxReceipt) Confirmed() bool {
	return len(r.BlockHash) > 0
}

// Use with previous method
func (r *TxReceipt) Successful() bool {
	if len(r.Status) > 0 {
		return r.Status == receiptStatusSuccessful
	}
	return true
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

func NewRPCClient(name, url, timeout string) *RPCClient {
	rpcClient := &RPCClient{Name: name, Url: url}
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

// func (r *RPCClient) GetWork() ([]string, error) {
// 	rpcResp, err := r.doPost(r.Url, "eth_getWork", []string{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	var reply []string
// 	err = json.Unmarshal(*rpcResp.Result, &reply)
// 	return reply, err
// }
func (r *RPCClient) GetAccounts(uid string) ([]Account, error) {
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	// var uidArray [1]string
	// uidArray[0] = uid

	// // var dataSlice []int = foo()
	// var interfaceSlice []interface{} = make([]interface{}, 3)
	// // for i, d := range dataSlice {
	// // 	interfaceSlice[i] = d
	// // }
	// interfaceSlice[0] = 0
	// interfaceSlice[1] = "get_accounts_by_uid"
	// interfaceSlice[2] = uidArray

	// rpcResp, err := r.doPost(r.Url, "call", []interface{}{0, "get_accounts_by_uid", []interface{uidArray}})
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "get_accounts_by_uid", [["` + uid + `"]]], "id": ` + reqIdString + `}`
	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		return nil, err
	}
	if rpcResp.Result != nil {
		var reply []Account
		err = json.Unmarshal(*rpcResp.Result, &reply)
		return reply, err
	}
	return nil, nil
}

//通过getfullaccount获取各个余额
func (r *RPCClient) GetBalance(uid string) (change, points, balance uint64, err error) {
	// curl --data '{"jsonrpc": "2.0", "method": "call", "params": [0, "get_full_account", [["250926091"]]], "id": 1}' http://localhost:8091/rpc
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "get_full_account", ["` + uid + `"]], "id": ` + reqIdString + `}`

	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		return
	}
	if rpcResp.Result != nil {
		var reply FullAccount
		err = json.Unmarshal(*rpcResp.Result, &reply)
		change = reply.Statistics.Prepaid
		points = reply.Statistics.Csaf
		balance = reply.Statistics.Core_balance
		return
	}
	return
}

// func (r *RPCClient) GetBalanceOfCoin(uid, coin string) (balance string, err error) {
// 	// curl --data '{"jsonrpc": "2.0", "method": "call", "params": [0, "get_account_balances", ["250926091", [0,1]]], "id": 1}' http://47.52.155.181:10011/rpc
// 	balance = "0"
// 	var reqId uint16
// 	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
// 	// fmt.Printf("%x\n", reqId)
// 	reqIdString := fmt.Sprintf("%d", reqId)
// 	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "get_account_balances", ["` + uid + `", [` + coin + `]]], "id": ` + reqIdString + `}`

// 	rpcResp, err := r.doPost2(r.Url+"/rpc", content)
// 	if err != nil {
// 		return
// 	}
// 	// {
// 	// 	"amount": 1098850704,
// 	// 	"asset_id": 0
// 	//   }
// 	type Ret struct {
// 		Amount   string `json:"amount"`
// 		Asset_id int64  `json:"asset_id"`
// 	}
// 	if rpcResp.Result != nil {
// 		var reply []Ret
// 		err = json.Unmarshal(*rpcResp.Result, &reply)
// 		if err != nil {
// 			return
// 		}

// 		for _, v := range reply {
// 			asid := fmt.Sprintf("%d", v.Asset_id)
// 			if asid == coin {
// 				balance = v.Amount
// 			}
// 		}
// 		return
// 	}
// 	return
// }
func (r *RPCClient) GetBalanceOfCoin(uid, coin string) (balance string, err error) {
	// curl --data '{"jsonrpc": "2.0", "method": "call", "params":[0,"list_account_balances",["250926091"]], "id": 1}' http://localhost:8091
	// curl --data '{"jsonrpc": "2.0", "method": "call", "params": [0, "get_account_balances", ["250926091", [0,1]]], "id": 1}' http://47.52.155.181:10011/rpc
	balance = "0"
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "list_account_balances", ["` + uid + `"]], "id": ` + reqIdString + `}`

	rpcResp, err := r.doPost2(r.Url, content)
	log.Println("GetBalanceOfCoin:", rpcResp)
	if err != nil {
		return
	}
	// {
	// 	"amount": 1098850704,
	// 	"asset_id": 0
	//   }
	type Ret struct {
		Amount   interface{} `json:"amount"`
		Asset_id int64       `json:"asset_id"`
	}
	if rpcResp.Result != nil {
		var reply []Ret
		err = json.Unmarshal(*rpcResp.Result, &reply)
		if err != nil {
			return
		}

		for _, v := range reply {
			asid := fmt.Sprintf("%d", v.Asset_id)
			if asid == coin {

				switch amountType := v.Amount.(type) {
				case float64:
					balance = fmt.Sprintf("%.0f", amountType)
				case string:
					balance = amountType
				default:
					log.Println("amount type not int and string")
				}
				// balance = v.Amount
			}
		}
		return
	}
	return
}
func (r *RPCClient) GetInfo() (last_irreversible_block_num uint64, err error) {
	// curl --data '{"jsonrpc": "2.0", "method": "call", "params": [0, "info", [[]]], "id": 1}' http://localhost:8091/rpc
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "info", []], "id": ` + reqIdString + `}`

	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		return
	}
	if rpcResp.Result != nil {
		reply := make(map[string]interface{})
		err = json.Unmarshal(*rpcResp.Result, &reply)
		if _, ok := reply["last_irreversible_block_num"]; !ok {
			//存在
			err = errors.New("last_irreversible_block_num is not exist")
			return
		}
		// fmt.Println("type:", reply["last_irreversible_block_num"].(type))

		// v := reflect.ValueOf(reply["last_irreversible_block_num"])
		// //获取传递参数类型
		// v_t := v.Type()

		// //类型名称对比
		// fmt.Println("type:", v_t.String())

		switch value := reply["last_irreversible_block_num"].(type) {
		case float64:
			// fmt.Printf("uint32 %v\n", value)
			last_irreversible_block_num = uint64(value)
		default:
			fmt.Printf("default %#v\n", value)
		}
		// ret, ok := reply["last_irreversible_block_num"].(int)
		// if !ok {

		// 	err = errors.New("last_irreversible_block_num is not uint64")
		// 	return
		// }

	}
	return
}

// 安全起见，我们只在节点状态正常的情况下处理提现。 通过info命令/API来检验。

// head_block_time 应当在15秒以内
// participation 应大于 80, 意味着80%的区块生产者在线且状态正常。
func (r *RPCClient) CheckSync() bool {
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "info", []], "id": ` + reqIdString + `}`

	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		// fmt.Println("232")
		return false
	}
	if rpcResp.Result != nil {
		reply := make(map[string]interface{})
		err = json.Unmarshal(*rpcResp.Result, &reply)
		if _, ok := reply["head_block_time"]; !ok {
			//存在
			// fmt.Println("239")
			return false
		}

		switch value := reply["head_block_time"].(type) {
		case string:
			format := "2006-01-02T15:04:05"
			headtime, err := time.Parse(format, value)
			if err != nil {
				fmt.Println("247:", err)
				return false
			}
			//比较
			fmt.Printf("%d\n", time.Now().Unix())
			fmt.Printf("%d\n", headtime.Unix())
			if time.Now().Unix()-headtime.Unix() <= 15 {
				return true
			}
		default:
			return false
		}
	}
	return false
}
func (r *RPCClient) ColletcPoints(from, to, amount string) (reply CollectPoints, err error) {
	// curl --data '{"jsonrpc": "2.0", "method": "call", "params":[0, "collect_csaf",["250926091", "250926091", 1, "YOYO", true]], "id": 1}' http://localhost:8091/rpc
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	// content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "collect_csaf", ["` + uid + `"]], "id": ` + reqIdString + `}`
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "collect_csaf", ["` + from + `","` + to + `","` + amount + `","YOYO",true]], "id": ` + reqIdString + `}`
	// {"jsonrpc": "2.0", "method": "call", "params":[0, "collect_csaf",["244958118","244958118",0.22,"YOYO",true]], "id": 2603870902}

	// map[ref_block_num:6949 ref_block_prefix:2.18008016e+08 expiration:2018-10-17T02:24:36 operations:[[6 map[fee:map[total:map[amount:100000 asset_id:0] options:map[from_csaf:map[amount:100000 asset_id:0]]] from:2.44958118e+08 to:2.44958118e+08 amount:map[amount:165000 asset_id:0] time:2018-10-17T02:24:00]]] signatures:[20090706bbc5b7ead8c6e2f527b6ef80c04e50b8883ef10d661a499fdcaa8b685229541afb274cdd19553d7cf7aaf6c52a012aca30f81ad24e23de06bebe034e0b 1f54e30b9d98e4f309d20340e159cb30d770120027f332d63f520a405a6402bca138ef3a4979ef88864a51656e0060ef46c9e29bc3ae417b028cba1e9eb0fcf5be]]

	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		return
	}
	if rpcResp.Result != nil {
		err = json.Unmarshal(*rpcResp.Result, &reply)
		return
	}
	return
}

// func (r *RPCClient) GetBalance(uid string) ([]AccountBalance, error) {
// 	// curl --data '{"jsonrpc": "2.0", "method": "call", "params": [0, "get_account_balances", ["244958118", [0,1]]], "id": 1}' http://127.0.0.1:8090/rpc
// 	var reqId uint32
// 	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
// 	// fmt.Printf("%x\n", reqId)
// 	reqIdString := fmt.Sprintf("%d", reqId)
// 	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "get_account_balances", ["` + uid + `",[0,1]]], "id": ` + reqIdString + `}`

// 	rpcResp, err := r.doPost2(r.Url, content)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if rpcResp.Result != nil {
// 		var reply []AccountBalance
// 		err = json.Unmarshal(*rpcResp.Result, &reply)
// 		return reply, err
// 	}
// 	return nil, nil
// }
func (r *RPCClient) Unlocks(pass string) error {
	// curl --data '{"jsonrpc": "2.0", "method": "call", "params":[0, "unlock", ["yoyow-pass"]], "id": 1}' http://127.0.0.1:8091/rpc
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "unlock",["` + pass + `"]], "id": ` + reqIdString + `}`

	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		return err
	}
	if rpcResp.Result != nil {
		var reply interface{}
		err = json.Unmarshal(*rpcResp.Result, &reply)
		return errors.New(fmt.Sprintf("%v", reply))
	}
	return nil
}
func (r *RPCClient) Locks() error {
	// {"jsonrpc": "2.0", "method": "call", "params": [0, "lock", []], "id": 1}
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "lock",[]], "id": ` + reqIdString + `}`

	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		return err
	}
	if rpcResp.Result != nil {
		var reply interface{}
		err = json.Unmarshal(*rpcResp.Result, &reply)
		return errors.New(fmt.Sprintf("%v", reply))
	}
	return nil
}
func (r *RPCClient) Transfer(from, to, amount, memo, coin string) (*TransferResponse, error) {
	// curl -d '{"jsonrpc": "2.0", "method": "transfer", "params": [123456789,123456789,"1","YOYO",null,true], "id": 1}' http://127.0.0.1:8091/rpc
	// curl --data '{"jsonrpc": "2.0", "method": "call", "params":[0, "transfer",[250926091, 209414065, "10", "YOYO", "feho", true]], "id": 1}' http://localhost:8091/rpc
	// {"jsonrpc": "2.0", "method": "call", "params": [0, "transfer",[244958118,226369314,"10","YOYO","test",true]], "id": 2933752513}
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "transfer", [` + from + `,` + to + `,"` + amount + `","` + coin + `","` + memo + `",true]], "id": ` + reqIdString + `}`

	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		return nil, err
	}
	if rpcResp.Result != nil {
		var reply TransferResponse
		err = json.Unmarshal(*rpcResp.Result, &reply)
		return &reply, err
	}
	return nil, nil
}
func (r *RPCClient) TransferRaw(from, to, amount, memo, coin string) (ret string, err error) {
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "call", "params":[0, "transfer",[` + from + `,` + to + `,"` + amount + `","` + coin + `","` + memo + `",true]], "id": ` + reqIdString + `}`
	// {"jsonrpc": "2.0", "method": "call", "params": [0, "transfer", [278137833,278705639,"100","KFC","usermemo",true]], "id": 61101}
	// curl --data '{"jsonrpc": "2.0", "method": "call", "params":[0, "transfer",[250926091, 209414065, "10", "YOYO", "feho", true]], "id": 1}' http://localhost:8091/rpc
	rpcResp, err := r.doPost2(r.Url+"/rpc", content)
	if err != nil {
		return
	}
	var transferReturnJson string
	if rpcResp.Result != nil {

		{
			// var reply TransferResponse
			j, errs := json.Marshal(rpcResp.Result)
			if errs != nil {
				err = errs
				return
			}
			transferReturnJson = string(j)
		}

		// {
		// temp := make(map[string]json.RawMessage)
		// err = json.Unmarshal(*rpcResp.Result, &temp)
		// 	operations, ok := temp["operations"]
		// 	if ok {
		// 		j, errs := json.Marshal(&operations)
		// 		if errs != nil {
		// 			err = errs
		// 			return
		// 		}
		// 		ret = string(j)
		// 	}
		// }

		log.Println("transfer:", transferReturnJson)
		tranId, errs := r.GenerateTransactionId(transferReturnJson)
		if errs != nil {
			log.Println("GenerateTransactionId:", errs)
			err = errs
			return
		}
		ret = tranId
		log.Println("tranID:", tranId)
	}
	return
}
func (r *RPCClient) GetAccountHistory(account, start, limit, end string) (reply []AccountHistory, err error) {
	fmt.Println(start, ":", limit, ":", end)
	// curl --data '{"jsonrpc": "2.0", "method": "call", "params":[0,"get_relative_account_history",["250926091",null,10,10,0]], "id": 1}' http://localhost:8091
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "get_relative_account_history", ["` + account + `",0,` + start + `,` + limit + `,` + end + `]], "id": ` + reqIdString + `}`

	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		return
	}
	if rpcResp.Result != nil {
		err = json.Unmarshal(*rpcResp.Result, &reply)
		return
	}
	return
}
func (r *RPCClient) GenerateTransactionId(oper string) (transactionid string, err error) {
	// curl --data '{"jsonrpc": "2.0", "method": "call", "params":[0,"get_transaction_id",[{"operations":[[0,{"fee":{"total":{"amount":100000,"asset_id":0}},"from":250926091,"to":223331844,"amount":{"amount":100000,"asset_id":0},"extensions":{}}]]}]], "id": 1}' http://localhost:8091

	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "call", "params": [0, "get_transaction_id", [` + oper + `]], "id": ` + reqIdString + `}`

	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		return
	}
	if rpcResp.Result != nil {
		err = json.Unmarshal(*rpcResp.Result, &transactionid)
		return
	}
	return
}
func (r *RPCClient) GetTransactionId(blockNum, trx_in_block string) (transactionId string, err error) {
	// curl -d '{"jsonrpc": "2.0", "method": "get_block", "params": [160000], "id": 1}' http://127.0.0.1:8091/rpc
	// {"id":1,"result":{"previous":"00ad36fa9b87592e0a0da371d88d3d79be9af1ef","timestamp":"2018-10-17T08:20:24","witness":28465,"transaction_merkle_root":"cf1663139cd34ba6adcaa4b218a43550baea7ba6","witness_signature":"2057bcc6b3dece2a96ebbbccf5cb5e57d11940c93edcc2e381e8f7d0324d1dfc190dbfd55cc0d647e36562f3a01ff3d74c12d72d87690715627aa17688e35f9a09","transactions":[{"ref_block_num":14074,"ref_block_prefix":777619355,"expiration":"2018-10-17T08:20:51","operations":[[0,{"fee":{"total":{"amount":28984,"asset_id":0},"options":{"from_csaf":{"amount":28984,"asset_id":0}}},"from":244958118,"to":226369314,"amount":{"amount":100000,"asset_id":0},"memo":{"from":"YYW6BThRL3arZjTHNSpDVjegMkkAaiYcMiuy6p9o6KSZVbQvBLiYF","to":"YYW4uqFfSWcBwDCuLNjruwMeba5MrGkaiJqGKznJn3Qh6GEJwH94g","nonce":"14701288948159977843","message":"10799a3f9284de45db214f292d04d51a"}}]],"signatures":["205d2bc026ef98826428083c1d34aec7dc0e4cd2f0892a477ded74f6e168bc1b7575ecfbfe3963a591e9c325bdbad55838faa3a28add149e366383426b587cd6b4"],"operation_results":[[0,{}]]}],"block_id":"00ad36fb5e76898464e28d1d3b68f232b8a4b2c0","signing_key":"YYW71suPihtG7jJAGiVBCkd63ECHYebQaPa894oy3r54zk3eM1itt","transaction_ids":["266c1b0105267553c46552387ba3d238c8aa6196"]}}
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	content := `{"jsonrpc": "2.0", "method": "get_block", "params": [` + blockNum + `], "id": ` + reqIdString + `}`

	rpcResp, err := r.doPost2(r.Url, content)
	if err != nil {
		return
	}
	if rpcResp.Result != nil {
		reply := make(map[string]interface{})
		err = json.Unmarshal(*rpcResp.Result, &reply)
		if err != nil {
			return
		}
		// transaction_ids
		transaction_ids, ok := reply["transaction_ids"].([]interface{})
		if ok {
			// aString := make([]string, len(transaction_ids))
			for i, v := range transaction_ids {
				// aString[i] = v.(string)
				iString := fmt.Sprintf("%d", i)
				if iString == trx_in_block {
					transactionId = v.(string)
					return
				}
			}
		}

	}
	return
}
func (r *RPCClient) GetBestSequenceNumber(account string) (maxNumber uint64, err error) {
	rep, err := r.GetAccountHistory(account, "0", "1", "0")
	if err != nil {
		return
	}
	if len(rep) == 0 {
		err = errors.New("GetBestSequenceNumber ret len is 0")
		return
	}
	maxNumber = rep[0].Sequence
	return
}

type ActionResult struct {
	From          string `json:"from,omitempty"`
	To            string `json:"to,omitempty"`
	TransactionId string `json:"transactionId"`
	Memo          string `json:"memo"`
	Seq           string `json:"seq"`
	Quantity      string `json:"quantity"`
	Packed        bool   `json:"packed"`
	Symbol        string `json:"symbol,omitempty"`
}

func (s *RPCClient) VerifyAll(acc, seq string) (ret []ActionResult, err error) {

	out, err := s.GetBestSequenceNumber(acc)
	if err != nil {
		return
	}
	start, err := strconv.ParseInt(seq, 10, 64)
	if err != nil {
		return
	}
	var end uint64
	if out-uint64(start) > 100 {
		end = uint64(start) + 100 - 1
	}

	accountHis, err := s.GetAccountHistory(acc, seq, "100", fmt.Sprintf("%d", end)) //返回 10 9 8
	if err != nil {
		log.Println("GetAccountHistory err:", err)
		return
	}
	// 取 result[N]["op"]["block_num"], 若比 last_irreversible_block_num小, 意味着是可信的, 需要被处理.
	// 取 result[N]["op"]["op"][0], 如 == 0, 则是一个transfer请求. (实际上肯定会等于0，因为请求参数里已约定只取transfer记录)
	// 取 result[N]["op"]["op"][1]["to"], 验证是否与自身account ID 相同。若相同，即为一个充值请求。
	// 取 result[N]["op"]["op"][1]["amount"]["asset_id"] 验证其是否 == 0, 若是，即表示该资产为YOYO
	// 取 result[N]["op"]["op"][1]["amount"]["amount"], 该值为充值数量. 切记：该数字精度为5.
	// 取 result[N]["memo"], 为该转账记录的MEMO信息。此处已经被解密过了，它可以用来作为充值客户的识别符。
	// 保存 result[N]["sequence"] 作为最新的sequence序号，下一次循环时需使用。
	// 保存 result[N]["op"]["trx_in_block"] 留作后续使用。
	// 使用 get_block 命令/API 来获取此次转账的txid ,如下：(记得修改params值为 block_num):
	// curl -d '{"jsonrpc": "2.0", "method": "get_block", "params": [160000], "id": 1}' http://127.0.0.1:8091/rpc
	// 记录返回值为 new_response, 保存 new_response["result"]["transaction_ids"][trx_in_block] 作为本次充值的txid 留作后用.
	last, err := s.GetInfo()
	if err != nil {
		return
	}
	for _, v := range accountHis {
		var temp ActionResult
		temp.Memo = v.Memo
		temp.Seq = fmt.Sprintf("%d", v.Sequence)
		temp.Packed = false
		if v.Op.Block_num < last {
			log.Println("block num:", v.Op.Block_num, " last:", last)
			temp.Packed = true
		}
		p, ok := (v.Op.Op[1]).(map[string]interface{})
		if ok {
			fl := p["to"].(float64)
			f2 := p["from"].(float64)
			temp.From = fmt.Sprintf("%d", uint64(f2))
			temp.To = fmt.Sprintf("%d", uint64(fl))

			amo, ok := (p["amount"]).(map[string]interface{})
			if !ok {
				log.Println("amount convert to map error")
				continue
			}
			temp.Symbol = fmt.Sprintf("%v", amo["asset_id"])
			// if temp.Symbol == "0" {
			// 	temp.Symbol = "YOYO"
			// } else if temp.Symbol == "46" {
			// 	temp.Symbol = "KFC"
			// }
			switch value := amo["amount"].(type) {
			case float64:
				a := fmt.Sprintf("%d", uint64(value))
				temp.Quantity = a
			case string:
				temp.Quantity = value
			default:
				temp.Quantity = fmt.Sprintf("%v", value)
			}
		}
		tid, err := s.GetTransactionId(fmt.Sprintf("%d", v.Op.Block_num), fmt.Sprintf("%d", v.Op.Trx_in_block))
		if err != nil {
			log.Println("GetTransactionId:", err)
			continue
		}
		temp.TransactionId = tid
		ret = append([]ActionResult{temp}, ret...)
	}
	if len(ret) == 0 {
		err = errors.New("no transaction in this range")
	}
	return
}

//返回此笔交易的转入金额
// VerifyIn(account, seq, symbol)
func (s *RPCClient) Verify(outAccount, inAccount, seq, symbol string) (ret []ActionResult, err error) {
	var acc string
	if outAccount != "" {
		acc = outAccount
	}
	if inAccount != "" {
		acc = inAccount
	}
	fmt.Println("acc:", acc)
	out, err := s.GetBestSequenceNumber(acc)
	if err != nil {
		return
	}
	start, err := strconv.ParseInt(seq, 10, 64)
	if err != nil {
		return
	}
	var end uint64
	if out-uint64(start) > 100 {
		end = uint64(start) + 100
	}

	accountHis, err := s.GetAccountHistory(acc, seq, "100", fmt.Sprintf("%d", end)) //返回 10 9 8
	if err != nil {
		log.Println("GetAccountHistory err:", err)
		return
	}
	// 取 result[N]["op"]["block_num"], 若比 last_irreversible_block_num小, 意味着是可信的, 需要被处理.
	// 取 result[N]["op"]["op"][0], 如 == 0, 则是一个transfer请求. (实际上肯定会等于0，因为请求参数里已约定只取transfer记录)
	// 取 result[N]["op"]["op"][1]["to"], 验证是否与自身account ID 相同。若相同，即为一个充值请求。
	// 取 result[N]["op"]["op"][1]["amount"]["asset_id"] 验证其是否 == 0, 若是，即表示该资产为YOYO
	// 取 result[N]["op"]["op"][1]["amount"]["amount"], 该值为充值数量. 切记：该数字精度为5.
	// 取 result[N]["memo"], 为该转账记录的MEMO信息。此处已经被解密过了，它可以用来作为充值客户的识别符。
	// 保存 result[N]["sequence"] 作为最新的sequence序号，下一次循环时需使用。
	// 保存 result[N]["op"]["trx_in_block"] 留作后续使用。
	// 使用 get_block 命令/API 来获取此次转账的txid ,如下：(记得修改params值为 block_num):
	// curl -d '{"jsonrpc": "2.0", "method": "get_block", "params": [160000], "id": 1}' http://127.0.0.1:8091/rpc
	// 记录返回值为 new_response, 保存 new_response["result"]["transaction_ids"][trx_in_block] 作为本次充值的txid 留作后用.
	last, err := s.GetInfo()
	if err != nil {
		return
	}
	for _, v := range accountHis {
		var temp ActionResult
		temp.Memo = v.Memo
		temp.Seq = fmt.Sprintf("%d", v.Sequence)
		temp.Packed = false
		if v.Op.Block_num < last {
			log.Println("block num:", v.Op.Block_num, " last:", last)
			temp.Packed = true
		}
		p, ok := (v.Op.Op[1]).(map[string]interface{})
		if ok {
			fl := p["to"].(float64)
			f2 := p["from"].(float64)
			from := fmt.Sprintf("%d", uint64(f2))
			to := fmt.Sprintf("%d", uint64(fl))
			if inAccount != "" && to != inAccount {
				log.Println("transfer to account:", to)
				continue
			}
			if outAccount != "" && from != outAccount {
				log.Println("transfer to account:", to)
				continue
			}
			amo, ok := (p["amount"]).(map[string]interface{})
			if !ok {
				log.Println("amount convert to map error")
				continue
			}
			assetId := fmt.Sprintf("%v", amo["asset_id"])
			if symbol == "YOYO" {
				if assetId != "0" {
					log.Println("assetid!=0:", assetId)
					continue
				}
			}
			switch value := amo["amount"].(type) {
			case float64:
				a := fmt.Sprintf("%d", uint64(value))
				temp.Quantity = a
			case string:
				temp.Quantity = value
			default:
				temp.Quantity = fmt.Sprintf("%v", value)
			}
			// aa, ok := amo["amount"].(float64)
			// if ok {
			// 	a := fmt.Sprintf("%d", uint64(aa))
			// 	temp.Quantity = a
			// }

		}
		// var opinop OpInOp1
		// opinop1 := v.Op.Op[1]
		// if _, ok := opinop1["to"]; !ok {
		// 	//存在
		// 	err = errors.New("opinop1[\"to\"] is not exist")
		// 	continue
		// }
		// fmt.Println(v.Sequence, ":", v.Memo, ":", v.Description)
		// fmt.Println("opinop1", opinop1)
		//获取交易id
		tid, err := s.GetTransactionId(fmt.Sprintf("%d", v.Op.Block_num), fmt.Sprintf("%d", v.Op.Trx_in_block))
		if err != nil {
			log.Println("GetTransactionId:", err)
			continue
		}
		temp.TransactionId = tid
		ret = append(ret, temp)
	}
	return
}
func (s *RPCClient) VerifyIn(inAccount, seq, symbol string) (ret []ActionResult, err error) {
	return s.Verify("", inAccount, seq, symbol)
}
func (s *RPCClient) VerifyOut(outAccount, seq, symbol string) (ret []ActionResult, err error) {
	return s.Verify(outAccount, "", seq, symbol)
}
func (r *RPCClient) doPost(url string, method string, params interface{}) (*JSONRpcResp, error) {
	jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	data, _ := json.Marshal(jsonReq)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Close = true

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
func (r *RPCClient) doPost2(url string, content string) (*JSONRpcResp, error) {
	// jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	// data, _ := json.Marshal(jsonReq)
	log.Println("content:", content)
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

func (r *RPCClient) markSick() {
	r.Lock()
	r.sickRate++
	r.successRate = 0
	if r.sickRate >= 5 {
		r.sick = true
	}
	r.Unlock()
}
