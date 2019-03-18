package ont

import (
	// "../aes"
	// "bytes"
	// "encoding/json"
	// "errors"
	sdk "github.com/ontio/ontology-go-sdk"
	// "io/ioutil"
	"log"
	// "math/big"
	// "net/http"
	// "regexp"
	// "sync"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"strconv"
)

type RPCClient struct {
	Url    string
	Wallet string
}

func NewRPCClient(url, wallet string) *RPCClient {
	rpcClient := &RPCClient{Url: url, Wallet: wallet}
	return rpcClient
}

func (r *RPCClient) GetDefaultAccount(pass string) (addr string, err error) {
	OntSdk := sdk.NewOntologySdk()
	OntSdk.NewRpcClient().SetAddress(r.Url)

	Wallet, err := OntSdk.OpenWallet(r.Wallet)
	if err != nil {
		log.Printf("account.Open error:%s\n", err)
		return
	}
	DefAcc, err := Wallet.GetDefaultAccount([]byte(pass))
	if err != nil {
		log.Printf("GetDefaultAccount error:%s\n", err)
		return
	}
	addr = DefAcc.Address.ToBase58()
	return
}

func (r *RPCClient) GetNewAddress(userpass, encryptfilepass string) (addr string, err error) {

	// NewDefaultSettingAccount(passwd []byte) (*Account, error)
	OntSdk := sdk.NewOntologySdk()
	OntSdk.NewRpcClient().SetAddress(r.Url)

	Wallet, err := OntSdk.OpenWallet(r.Wallet)
	if err != nil {
		log.Printf("OpenWallet error:%s\n", err)
		return
	}
	DefAcc, err := Wallet.NewDefaultSettingAccount([]byte(userpass))
	if err != nil {
		log.Printf("NewDefaultSettingAccount error:%s\n", err)
		return
	}
	err = Wallet.Save()
	if err != nil {
		log.Printf("Save error:%s\n", err)
		return
	}
	addr = DefAcc.Address.ToBase58()
	log.Printf("newAddress:%s,pass:%s\n", addr, userpass)
	// err = r.writeFile(dir, encryptfilepass, addr, userpass)
	return
}

// func (r *RPCClient) writeFile(dir, password, filename, content string) (err error) {
// 	key := []byte(password)
// 	result, err := aes.AesEncrypt([]byte(content), key)
// 	if err != nil {
// 		return
// 	}
// 	// fmt.Println(base64.StdEncoding.EncodeToString(result))
// 	err = ioutil.WriteFile(dir+"/"+filename, result, 0644)

// 	if err != nil {
// 		log.Println("write file err:", err)
// 		return err
// 	}
// 	return nil
// }

func (r *RPCClient) GetBalance(address string) (balance uint64, err error) {
	// testOntSdk.Native.Ont.
	OntSdk := sdk.NewOntologySdk()
	OntSdk.NewRpcClient().SetAddress(r.Url)

	addr, err := common.AddressFromBase58(address)
	if err != nil {
		log.Printf("AddressFromBase58 error:%s\n", err)
		return
	}
	balance, err = OntSdk.Native.Ont.BalanceOf(addr)

	return
}
func (r *RPCClient) GetTransaction(txHash string) (event *sdkcom.SmartContactEvent, err error) {

	OntSdk := sdk.NewOntologySdk()
	OntSdk.NewRpcClient().SetAddress(r.Url)

	event, err = OntSdk.GetSmartContractEvent(txHash)

	return
}
func (r *RPCClient) Send(from, pass, to, amount, gasprice string) (ret string, err error) {
	// testOntSdk.Native.Ont.
	OntSdk := sdk.NewOntologySdk()
	OntSdk.NewRpcClient().SetAddress(r.Url)
	// GetAccountByAddress(address string, passwd []byte) (*Account, error)
	Wallet, err := OntSdk.OpenWallet(r.Wallet)
	if err != nil {
		log.Printf("OpenWallet error:%s\n", err)
		return
	}
	fromAccount, err := Wallet.GetAccountByAddress(from, []byte(pass))
	if err != nil {
		log.Printf("GetAccountByAddress error:%s\n", err)
		return
	}
	addr, err := common.AddressFromBase58(to)
	if err != nil {
		log.Printf("AddressFromBase58 error:%s\n", err)
		return
	}
	amon, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		log.Printf("ParseUint error:%s\n", err)
		return
	}
	gasPrice, err := strconv.ParseUint(gasprice, 10, 64)
	if err != nil {
		log.Printf("ParseUint error:%s\n", err)
		return
	}
	// gasPrice := uint64(0)
	gasLimit := uint64(20000)
	// Transfer(gasPrice, gasLimit uint64, from *Account, to common.Address, amount uint64) (common.Uint256, error)
	hash, err := OntSdk.Native.Ont.Transfer(gasPrice, gasLimit, fromAccount, addr, amon)
	if err != nil {
		log.Printf("Transfer error:%s\n", err)
		return
	}
	ret = hash.ToHexString()
	return
}
