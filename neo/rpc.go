package neo

import (
	// "encoding/hex"
	// "encoding/json"
	"fmt"
	// "math/big"
	"errors"
	// "github.com/stretchr/testify/require"
	// "github.com/dynamicgo/config"
	neorpc "github.com/CityOfZion/neo-go/pkg/rpc"
	// "github.com/stretchr/testify/assert"
	"../aes"
	"context"
	"io/ioutil"
	"log"
	"strconv"
)

func GetBalance(url, addr, assetid string) (ret string, err error) {
	if assetid[:2] != "0x" {
		assetid = "0x" + assetid
	}
	ret = "0"
	opts := neorpc.ClientOptions{}

	client, err := neorpc.NewClient(context.TODO(), url, opts)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := client.GetAccountState(addr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("script:", resp.Result.ScriptHash)
	log.Println("balance:", resp.Result.Balances)
	// if len(resp.Result.Balances) == 0 {
	// 	ret = "0"
	// 	return
	// }
	for _, v := range resp.Result.Balances {
		if v.Asset == assetid {
			ret = v.Value
			return
		}
	}
	return
}

func GetNewAddress(url, userpass, dir, walletpassword string) (ret string, err error) {

	opts := neorpc.ClientOptions{}

	client, err := neorpc.NewClient(context.TODO(), url, opts)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.GetNewAddress()
	if err != nil {
		log.Println(err)
		return
	}
	ret = resp.Result
	log.Println("address:", resp.Result)
	if len(ret) == 0 {
		err = errors.New("returned address is null")
		return
	}
	err = writeFile(dir, walletpassword, ret, userpass)
	return
}
func GetHeight(url string) (ret string, err error) {

	opts := neorpc.ClientOptions{}

	client, err := neorpc.NewClient(context.TODO(), url, opts)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.GetBlockCount()
	if err != nil {
		log.Println(err)
		return
	}
	ret = fmt.Sprintf("%d", resp.Result)
	log.Println("height:", ret)
	return
}
func SendFrom(url, walletPass, from, to, amount, assetid, userPass, neo_dir string) (ret string, errs error) {
	if assetid[:2] != "0x" {
		assetid = "0x" + assetid
	}
	{
		////首先验证密码
		b, err := ioutil.ReadFile(neo_dir + "/" + from)
		if err != nil {
			// fmt.Print(err)
			log.Println("read file:", err)
			errs = errors.New("wrong pass")
			return
		}
		origData, err := aes.AesDecrypt(b, []byte(walletPass))
		if err != nil {
			log.Println("AesDecrypt:", err)
			errs = errors.New("wrong pass")
			return
		}
		// log.Println("13")
		if len(origData) == 0 {
			err = errors.New("decrpt is null")
			errs = err
			return
		}
		// log.Println("14")
		if userPass != string(origData) {
			log.Println(string(origData))
			err = errors.New("wrong password")
			errs = err
			return
		}
	}
	opts := neorpc.ClientOptions{}

	client, err := neorpc.NewClient(context.TODO(), url, opts)
	if err != nil {
		log.Println(err)
		return
	}
	// resp, err := client.GetNewAddress()
	// SendToAddress(to, assetid, value, fee, coin string) (*GetTransactionResponse, error)
	fee := "0"
	resp, err := client.SendFrom(from, to, assetid, amount, fee, "NEO")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Tx:", resp.Result)
	ret = resp.Result.Txid
	if len(ret) == 0 {
		errs = errors.New("transaction maybe failed")
		return
	}
	return
}
func Send(url, walletPass, to, amount, assetid, userPass, neo_dir string) (ret string, errs error) {
	if assetid[:2] != "0x" {
		assetid = "0x" + assetid
	}
	{
		////首先验证密码
		b, err := ioutil.ReadFile(neo_dir + "/" + to)
		if err != nil {
			// fmt.Print(err)
			log.Println("read file:", err)
			errs = err
			return
		}
		origData, err := aes.AesDecrypt(b, []byte(walletPass))
		if err != nil {
			log.Println("AesDecrypt:", err)
			errs = err
			return
		}
		// log.Println("13")
		if len(origData) == 0 {
			err = errors.New("decrpt is null")
			errs = err
			return
		}
		// log.Println("14")
		if userPass != string(origData) {
			log.Println(string(origData))
			err = errors.New("wrong password")
			errs = err
			return
		}
	}
	opts := neorpc.ClientOptions{}

	client, err := neorpc.NewClient(context.TODO(), url, opts)
	if err != nil {
		log.Println(err)
		return
	}
	// resp, err := client.GetNewAddress()
	// SendToAddress(to, assetid, value, fee, coin string) (*GetTransactionResponse, error)
	fee := "0"
	resp, err := client.SendToAddress(to, assetid, amount, fee, "NEO")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Tx:", resp.Result)
	ret = resp.Result.Txid

	return
}
func writeFile(dir, password, filename, content string) (err error) {
	key := []byte(password)
	result, err := aes.AesEncrypt([]byte(content), key)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(dir+"/"+filename, result, 0644)

	if err != nil {
		log.Println("write file err:", err)
		return err
	}
	return nil
}

type GetTransactionRes struct {
	Vin  []*neorpc.Input  `json:"vin,omitempty"`
	Vout []*neorpc.Output `json:"vout,omitempty"`
}

func GetTransaction(url, hash string) (res GetTransactionRes, err error) {

	opts := neorpc.ClientOptions{}

	client, err := neorpc.NewClient(context.TODO(), url, opts)
	if err != nil {
		log.Println(err)
		return
	}
	// GetRawTransaction(hash string, verbose bool) (*response, error)
	resp, err := client.GetRawTransaction(hash, true)
	if err != nil {
		log.Println(err)
		return
	}
	ret := resp.Result
	log.Println("rawtransaction block:", ret)
	res.Vin = ret.Inputs
	res.Vout = ret.Outputs
	return
}
func GetBlock(url, height string) (ret []neorpc.TransactionInBlock, err error) {

	opts := neorpc.ClientOptions{}

	client, err := neorpc.NewClient(context.TODO(), url, opts)
	if err != nil {
		log.Println(err)
		return
	}
	//
	hei, err := strconv.ParseInt(height, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.GetBlock(hei, true)
	if err != nil {
		log.Println(err)
		return
	}
	// ret = resp.Result
	// bl, err := json.Marshal(resp.Result)
	// if err != nil {
	// 	log.Printf("getblock json marshal Error: %s", err)
	// 	return
	// }
	// log.Println("getblock:", string(bl))
	trans := resp.Result.Transactions
	if len(trans) == 0 {
		err = errors.New("no transaction in this block")
		return
	}
	for _, tran := range trans {

		if tran.Type == "ContractTransaction" {
			// tr, errs := json.Marshal(tran)
			// if errs != nil {
			// 	log.Printf("getblock tran json marshal Error: %s", err)
			// 	err = errs
			// 	return
			// }
			// ret = string(tr)
			log.Println("tran.Type:", tran.Type)
			ret = append(ret, tran)
			// log.Println("getblock:", ret)
		}
	}
	if len(ret) == 0 {
		err = errors.New("no ContractTransaction in this block")
	}
	return
}
