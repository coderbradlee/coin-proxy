package xmr

import (
	// "../aes"
	// "bytes"
	// "crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	// "io/ioutil"
	"log"
	// "net/http"
	// "regexp"
	// "../digest"
	// "encoding/base64"
	// "strings"
	// "golang.org/x/net/context"
	// "encoding/hex"
	// "io"
	"github.com/gabstv/go-monero/walletrpc"
	"strconv"
	"time"
)

type RPCClient struct {
	Url  string
	Name string
	// client   *http.Client
	Username string
	Password string
	// t := NewTransport("myUserName", "myP@55w0rd")
	// Tranp *digest.Transport
	// Tranp *digest.DigestRequest
	Client walletrpc.Client
}

type JSONRpcResp struct {
	Id      *json.RawMessage       `json:"id,omitempty"`
	Result  *json.RawMessage       `json:"result,omitempty"`
	Error   map[string]interface{} `json:"error,omitempty"`
	Jsonrpc string                 `json:"jsonrpc,omitempty"`
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}
func NewRPCClient(name, url, username, password, timeout string) *RPCClient {
	rpcClient := &RPCClient{Name: name, Url: url, Username: username, Password: password}
	// timeoutIntv := MustParseDuration(timeout)
	// rpcClient.client = &http.Client{
	// 	Timeout: timeoutIntv,
	// }
	// rpcClient.Tranp = digest.NewTransport(username, password)
	// log.Println("tranp:", *(rpcClient.Tranp))
	// ctx := context.Background()
	// rpcClient.Tranp = digest.New(ctx, username, password) // username & password
	client := walletrpc.New(walletrpc.Config{
		Address: url,
	})
	rpcClient.Client = client
	return rpcClient
}
func (r *RPCClient) GetAddress(account_index string) (ret string, err error) {
	// balance, unlocked, err := r.Client.GetBalance()
	return r.Client.CreateAddress(account_index)
}

type TransResponse struct {
	MaxHeight    uint64               `json:"maxHeight`
	Transactions []walletrpc.Transfer `json:"transactions`
}

func (r *RPCClient) GetTransaction(account_index, minHeight, maxHeight string) (ret TransResponse, err error) {
	ret.MaxHeight, err = r.Client.GetHeight()
	if err != nil {
		log.Println(err)
		return
	}
	acc, err := strconv.ParseUint(account_index, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	minh, err := strconv.ParseUint(minHeight, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	maxh, err := strconv.ParseUint(maxHeight, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	var req walletrpc.GetTransfersRequest
	req.In = true
	req.Account_index = acc
	req.FilterByHeight = true
	req.MinHeight = minh
	req.MaxHeight = maxh
	resp, err := r.Client.GetTransfers(req)
	if err != nil {
		log.Println(err)
		return
	}
	// for _, v := range resp.In {

	// }
	// return r.Client.CreateAddress(account_index)
	ret.Transactions = resp.In
	return
}
func (r *RPCClient) Transfer(fromIndex, to, amount string) (ret, fee string, err error) {
	// 	// Make a transfer
	amountI, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	acc, err := strconv.ParseUint(fromIndex, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	res, err := r.Client.Transfer(walletrpc.TransferRequest{
		Account_index: acc,
		Destinations: []walletrpc.Destination{
			{
				Address: to,
				Amount:  amountI, // 0.01 XMR
			},
		},
		Priority: walletrpc.PriorityUnimportant,
		Mixin:    1,
	})
	if err != nil {
		if iswerr, werr := walletrpc.GetWalletError(err); iswerr {
			// insufficient funds return a monero wallet error
			// walletrpc.ErrGenericTransferError
			message := fmt.Sprintf("Wallet error (id:%v) %v\n", werr.Code, werr.Message)
			err = errors.New(message)
			log.Println(err)
			return
		}
		log.Println("Error:", err.Error())
		return
	}
	log.Println("Transfer success! Fee:", walletrpc.XMRToDecimal(res.Fee), "Hash:", res.TxHash)
	ret = res.TxHash
	fee = fmt.Sprintf("%d", res.Fee)
	return
}

type Balance struct {
	Balance  string `json:"balance"`
	Unlocked string `json:"unlocked"`
}

func (r *RPCClient) GetBalance(account_index string) (ret Balance, err error) {
	balance, unlocked, err := r.Client.GetBalance(account_index)

	if err != nil {
		if iswerr, werr := walletrpc.GetWalletError(err); iswerr {
			// it is a monero wallet error
			message := fmt.Sprintf("Wallet error (id:%v) %v\n", werr.Code, werr.Message)
			err = errors.New(message)
			log.Println(err)
			return
		}
		log.Println("Error:", err.Error())
		return
	}

	log.Println("Balance:", walletrpc.XMRToDecimal(balance))
	log.Println("Unlocked balance:", walletrpc.XMRToDecimal(unlocked))
	ret.Balance = walletrpc.XMRToDecimal(balance)
	ret.Unlocked = walletrpc.XMRToDecimal(unlocked)
	return
}
