package main

import (
	// "encoding/json"
	// "io"
	// "log"
	// "net"
	// "net/http"

	// "strings"
	// "sync"
	// "sync/atomic"
	// "time"
	"./rpcclient"
	"encoding/hex"
	"fmt"
	// "fmt"
	"errors"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"io/ioutil"
	"log"
)

// func decodeHex(s string) []byte {
// 	b, err := hex.DecodeString(s)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return b
// }
func getNewAddress(pass string) (string, error) {
	c, err := rpcclient.New(cfg.Omni.Host, cfg.Omni.Port, cfg.Omni.Username, cfg.Omni.Password)
	if err != nil {
		return "", err
	}
	{
		success, err := c.Walletpassphrase(cfg.Omni.WalletPass)
		if err != nil {
			return "", err
		}
		if !success {
			return "", errors.New("WalletPass error!")
		}
	}

	// h, err := c.SendBtc(cfg.Omni.Account, to, amount)

	address, err := c.OmniGetNewAddress(cfg.Omni.Account)
	c.WalletLock()
	if err != nil {
		return "", err
	}
	//成功后保存文件，地址为文件名，内容为地址+密码的hash值
	hash := sha3.NewKeccak256()

	var buf []byte
	//hash.Write([]byte{0xcc})
	hash.Write([]byte(address + pass))
	buf = hash.Sum(buf)

	// fmt.Println(hex.EncodeToString(buf))
	// d1 := []byte("hello\ngo\n")
	err = ioutil.WriteFile("keys/"+address, []byte(hex.EncodeToString(buf)), 0644)

	if err != nil {
		log.Println("write file err:", err)
		return "", err
	}
	return address, nil
}
func getBalance(addr, pass string) (string, string, error) {
	b, err := ioutil.ReadFile("keys/" + addr)
	if err != nil {
		// fmt.Print(err)
		log.Println("read file:", err)
		return "", "", errors.New("wrong password!")
	}
	str := string(b)
	log.Println(addr, ":", str)

	hash := sha3.NewKeccak256()

	var buf []byte
	//hash.Write([]byte{0xcc})
	hash.Write([]byte(addr + pass))
	buf = hash.Sum(buf)
	if str != hex.EncodeToString(buf) {
		log.Println("password error:", str, "!=", hex.EncodeToString(buf))
		return "", "", errors.New("wrong password!")
	}
	c, err := rpcclient.New(cfg.Omni.Host, cfg.Omni.Port, cfg.Omni.Username, cfg.Omni.Password)
	if err != nil {
		return "", "", err
	}
	blan, err := c.OmniGetBalance(addr, cfg.Omni.Tokenid)
	if err != nil {
		return "", "", err
	}
	return blan.Balance, blan.Reserved, nil
}
func getBtcBalance(addr string) (string, error) {

	b, err := rpcclient.GetBTCBalanceByAddr(addr, cfg.Omni.Net)
	if err != nil {
		return "", err
	}
	return b, nil
}
func getAccountBalance() (string, error) {

	c, err := rpcclient.New(cfg.Omni.Host, cfg.Omni.Port, cfg.Omni.Username, cfg.Omni.Password)
	if err != nil {
		return "", err
	}
	b, err := c.GetBalance(cfg.Omni.Account)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", b), nil
}

// OmniSend(from,to string,propertyid int,amount string) (string, error)
func send(from, to, amount, fee, pass string) (string, error) {
	b, err := ioutil.ReadFile("keys/" + from)
	if err != nil {
		// fmt.Print(err)
		log.Println("read file:", err)
		return "", errors.New("wrong password!")
	}
	str := string(b)
	log.Println(from, ":", str)

	hash := sha3.NewKeccak256()

	var buf []byte
	//hash.Write([]byte{0xcc})
	hash.Write([]byte(from + pass))
	buf = hash.Sum(buf)
	if str != hex.EncodeToString(buf) {
		log.Println("password error:", str, "!=", hex.EncodeToString(buf))
		return "", errors.New("wrong password!")
	}
	c, err := rpcclient.New(cfg.Omni.Host, cfg.Omni.Port, cfg.Omni.Username, cfg.Omni.Password)
	if err != nil {
		return "", err
	}

	{
		success, err := c.Settxfee(fee)
		if err != nil {
			return "", err
		}
		if !success {
			return "", errors.New("settxfee error!")
		}
	}
	{
		success, err := c.Walletpassphrase(cfg.Omni.WalletPass)
		if err != nil {
			return "", err
		}
		if !success {
			return "", errors.New("WalletPass error!")
		}
	}

	h, err := c.OmniSend(from, to, cfg.Omni.Tokenid, amount)
	if err != nil {
		return "", err
	}
	return h, nil
}
func sendbtc(to, amount, fee string) (string, error) {

	c, err := rpcclient.New(cfg.Omni.Host, cfg.Omni.Port, cfg.Omni.Username, cfg.Omni.Password)
	if err != nil {
		return "", err
	}

	{
		success, err := c.Settxfee(fee)
		if err != nil {
			return "", err
		}
		if !success {
			return "", errors.New("settxfee error!")
		}
	}
	{
		success, err := c.Walletpassphrase(cfg.Omni.WalletPass)
		if err != nil {
			return "", err
		}
		if !success {
			return "", errors.New("WalletPass error!")
		}
	}

	h, err := c.SendBtc(cfg.Omni.Account, to, amount)
	c.WalletLock()
	{
		success, err := c.Settxfee("0.00001")
		if err != nil {
			return "", err
		}
		if !success {
			return "", errors.New("settxfee error!")
		}
	}
	if err != nil {
		return "", err
	}
	return h, nil
}
func sendbtcmany(to string, fee string) (string, error) {

	c, err := rpcclient.New(cfg.Omni.Host, cfg.Omni.Port, cfg.Omni.Username, cfg.Omni.Password)
	if err != nil {
		return "", err
	}

	{
		success, err := c.Settxfee(fee)
		if err != nil {
			return "", err
		}
		if !success {
			return "", errors.New("settxfee error!")
		}
	}
	{
		success, err := c.Walletpassphrase(cfg.Omni.WalletPass)
		if err != nil {
			return "", err
		}
		if !success {
			return "", errors.New("WalletPass error!")
		}
	}

	// h, err := c.SendBtc(cfg.Omni.Account, to, amount)
	h, err := c.SendBtcMany("http://"+cfg.Omni.Host+":"+fmt.Sprintf("%d", cfg.Omni.Port), cfg.Omni.Account, to)
	c.WalletLock()
	{
		success, err := c.Settxfee("0.00001")
		if err != nil {
			return "", err
		}
		if !success {
			return "", errors.New("settxfee error!")
		}
	}
	if err != nil {
		return "", err
	}
	return h, nil
}
