package rpcclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// GetTransactionCmd omni_gettransaction
type GetTransactionCmd struct {
	Txid string
}

// NewGetTransactionCmd ...
func NewGetTransactionCmd(txHash string) *GetTransactionCmd {
	return &GetTransactionCmd{
		Txid: txHash,
	}
}

// ListBlockTransactionsCmd omni_listblocktransactions
type ListBlockTransactionsCmd struct {
	Index int64
}

// NewListBlockTransactionsCmd ...
func NewListBlockTransactionsCmd(index int64) *ListBlockTransactionsCmd {
	return &ListBlockTransactionsCmd{
		Index: index,
	}
}

// GetBalanceCmd omni_getbalance
type GetBalanceCmd struct {
	Address    string
	PropertyID int64
}
type SendCmd struct {
	From       string
	To         string
	PropertyID int64
	Amount     string
}

// NewGetBalanceCmd ...
func NewGetBalanceCmd(address string, propertyID int64) *GetBalanceCmd {
	return &GetBalanceCmd{
		Address:    address,
		PropertyID: propertyID,
	}
}

func init() {
	flags := btcjson.UFWalletOnly
	btcjson.MustRegisterCmd("omni_gettransaction", (*GetTransactionCmd)(nil), flags)
	btcjson.MustRegisterCmd("omni_listblocktransactions", (*ListBlockTransactionsCmd)(nil), flags)
	btcjson.MustRegisterCmd("omni_getbalance", (*GetBalanceCmd)(nil), flags)
	btcjson.MustRegisterCmd("omni_send", (*SendCmd)(nil), flags)
}

//Get BTC balance Info
type ChainBalanceInfo struct {
	Status string           `json:"status"`
	Data   ChainBalanceData `json:"data"`
}

type ChainBalanceData struct {
	Network             string `json:"network"`
	Address             string `json:"address"`
	Confirmed_balance   string `json:"confirmed_balance"`
	Unconfirmed_balance string `json:"unconfirmed_balance"`
}

func GetBTCBalanceByAddr(address string, net string) (balance string, err error) {
	// /api/v2/get_address_balance/{NETWORK}/{ADDRESS}[/{MINIMUM CONFIRMATIONS}]
	//https://chain.so/api/v2/get_address_balance/BTCTEST/%s/%d
	//https://chain.so/api/v2/get_address_balance/BTC/%s/%d
	if len(address) == 0 {
		return "", errors.New("The Addres is Empty!!!")
	}
	var _url string
	if net == "mainnet" {
		_url = fmt.Sprintf("https://chain.so/api/v2/get_address_balance/BTC/%s/%d", address, 1)
		// GET /api/v2/get_address_balance/{NETWORK}/{ADDRESS}[/{MINIMUM CONFIRMATIONS}]
	} else {
		_url = fmt.Sprintf("https://chain.so/api/v2/get_address_balance/BTCTEST/%s/%d", address, 1)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	var rest ChainBalanceInfo
	resp, err := client.Get(_url)
	if err != nil {
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("The get Balance info is ", string(bs))
	err = json.Unmarshal(bs, &rest)
	if err != nil {
		//fmt.Println("There are some errors:", err)
		return getBtcExtendBalance(address, net)
	}
	if rest.Status != "success" {
		//err = errors.New("The Get Balance is wrong!!!!")
		return getBtcExtendBalance(address, net)
	}
	return rest.Data.Confirmed_balance, err
}
func getBtcExtendBalance(address string, net string) (balance string, err error) {
	// /api/v2/get_address_balance/{NETWORK}/{ADDRESS}[/{MINIMUM CONFIRMATIONS}]
	var _url string
	if net == "mainnet" {
		_url = fmt.Sprintf("https://blockchain.info/q/addressbalance/%s", address)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(_url)
	if err != nil {
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("The Btc get Balance info is ", string(bs))
	balanceValue, err := strconv.ParseFloat(string(bs), 10)
	btc := balanceValue / (100 * 1000 * 1000)
	retBalance := fmt.Sprintf("%.8f", btc)
	if err != nil {
		return "", err
	} else {
		return retBalance, nil
	}
}
