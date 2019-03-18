package ripple

import (
	"bytes"
	// "fmt"
	// "crypto/ecdsa"
	// "crypto/sha256"
	// "encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	// "io/ioutil"
	"encoding/base64"
	"log"
	// "math/big"
	// "github.com/rubblelabs/ripple/data"
	// "github.com/rubblelabs/ripple/websockets"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RPCClient struct {
	sync.RWMutex
	Url         string
	Name        string
	successRate int
	client      *http.Client
	user        string
	password    string
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}
func NewRPCClient(name, url, timeout, user, password string) *RPCClient {
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
	rpcClient.user = user
	rpcClient.password = password
	return rpcClient
}
func (cli *RPCClient) GetBalance(account string) (balance string, err error) {
	r := request{Method: "account_info"}
	p := param{
		Account: account,
	}
	r.Params = make([]param, 0)
	r.Params = append(r.Params, p)

	reply, err := cli.doPost(cli.Url, &r)
	if err != nil {
		return
	}
	log.Println("raw message:", string(*reply))
	var ret AccountInfoResult
	err = json.Unmarshal(*reply, &ret)
	if err != nil {
		return
	}
	if ret.Result.Status != "success" {
		err = errors.New(ret.Result.Error_message)
		return
	}
	bal := ret.Result.AccountData.Balance
	if bal == nil {
		err = errors.New("account not exists")
		return
	}
	// balance = ret.Result.AccountData.Balance.Float()
	balance = ret.Result.AccountData.Balance.Rat().RatString()
	return
}

// SignXaction sign a payment transaction using a rippled server.
func (cli *RPCClient) signXaction(secret string, account string, dest string, amount, memo int) (signed string, err error) {

	r := request{Method: "sign"}
	t := transaction{
		Account:         account,
		Amount:          amount,
		Destination:     dest,
		TransactionType: "Payment",
		DestinationTag:  uint32(memo),
	}
	p := param{
		Offline: false,
		Secret:  secret,
		TxJSON:  &t,
	}

	r.Params = make([]param, 0)
	r.Params = append(r.Params, p)

	reply, err := cli.doPost(cli.Url, &r)
	if err != nil {
		return
	}
	log.Println("raw message:", string(*reply))
	var ret SignedResult
	err = json.Unmarshal(*reply, &ret)
	if err != nil {
		return
	}
	if ret.Result.Status != "success" {
		err = errors.New(ret.Result.Error_message)
		return
	}
	// signed = ret.Result.Tx_json.TxnSignature
	signed = ret.Result.Tx_blob

	return
}

// SubmitSignedXaction submits a signed transaction blob for entry into the ledger.
func (cli *RPCClient) SubmitSignedXaction(txBlob string) (hash string, err error) {
	r := request{Method: "submit"}
	p := param{
		TxBlob: txBlob,
	}
	r.Params = make([]param, 0)
	r.Params = append(r.Params, p)

	// return queryServer(RippledURL, &r)
	reply, err := cli.doPost(cli.Url, &r)
	if err != nil {
		return
	}
	log.Println("raw message:", string(*reply))
	var ret SubmitResult
	err = json.Unmarshal(*reply, &ret)
	if err != nil {
		return
	}
	if ret.Result.Status != "success" {
		err = errors.New(ret.Result.Error_message)
		return
	}
	hash = ret.Result.Tx_json.Hash
	// signed = ret.Result.Tx_blob

	return
}
func (cli *RPCClient) Send(secret string, account string, dest string, amount, memo string) (hash string, err error) {
	amo, err := strconv.Atoi(amount)
	if err != nil {
		return
	}
	memoInt, err := strconv.Atoi(memo)
	if err != nil {
		return
	}
	signed, err := cli.signXaction(secret, account, dest, amo, memoInt)
	log.Println("signed:", signed, " err:", err)
	if err != nil {
		return
	}

	return cli.SubmitSignedXaction(signed)
}
func (cli *RPCClient) Verify(addr, ledger, seq string) (retAll VerifyResult, err error) {
	r := request{Method: "account_tx"}
	// min, err := strconv

	Ledger, err := strconv.ParseInt(ledger, 10, 64)
	if err != nil {
		return
	}
	Seq, err := strconv.ParseInt(seq, 10, 64)
	if err != nil {
		return
	}

	p := param{
		Account:          addr,
		Ledger_index_min: -1,
		Ledger_index_max: -1,
		Limit:            2, //需要验证是不是最新的100，要取start到start+100，而不是最新的100条
		// Offset:           startInt,
		// Marker: m,
		Forward: true,
	}
	p.Marker = nil
	if Ledger != 0 || Seq != 0 {
		var m Mark
		m.Ledger = Ledger
		m.Seq = Seq
		p.Marker = &m
	}
	r.Params = make([]param, 0)
	r.Params = append(r.Params, p)

	// return queryServer(RippledURL, &r)
	reply, err := cli.doPost(cli.Url, &r)
	if err != nil {
		return
	}
	log.Println("raw message:", string(*reply))
	var ret AccountTxResult
	err = json.Unmarshal(*reply, &ret)
	if err != nil {
		return
	}
	// log.Println("202")
	if ret.Result.Status != "success" {
		err = errors.New(ret.Result.Error_message)
		return
	}
	// {"from":"274769226","to":"278137833","transactionId":"dc4938c335c2b6a77fcffa9dd72a4df04502746c","memo":"","seq":"4","quantity":"10000000000","packed":true,"symbol":"46"}
	allTx := ret.Result.Transactions
	tranCount := len(allTx)
	if tranCount == 0 {
		err = errors.New("transaction is empty")
		return
	}
	// signed = ret.Result.Tx_blob
	for _, v := range allTx {
		var temp VerifyResultIn
		temp.Memo = fmt.Sprintf("%d", v.Tx.DestinationTag)
		temp.Seq = fmt.Sprintf("%d", v.Tx.Sequence)
		temp.Packed = v.Validated

		temp.From = v.Tx.Account
		temp.To = v.Tx.Destination
		temp.Quantity = v.Tx.Amount
		temp.TransactionId = v.Tx.Hash
		temp.TransactionResult = v.MetaData.TransactionResult
		// retAll = append([]VerifyResult{temp}, retAll...)
		retAll.Trans = append(retAll.Trans, temp)
	}
	retAll.Marker = ret.Result.Marker
	if retAll.Marker.Ledger == 0 && retAll.Marker.Seq == 0 {
		retAll.Marker.Ledger = Ledger
		retAll.Marker.Seq = Seq + int64(tranCount)
	}
	return
}
func (cli *RPCClient) GetTransaction(txid string) (retResult TransactionDetailIn, err error) {
	r := request{Method: "tx"}

	p := param{
		Transaction: txid,
	}
	r.Params = make([]param, 0)
	r.Params = append(r.Params, p)

	// return queryServer(RippledURL, &r)
	reply, err := cli.doPost(cli.Url, &r)
	if err != nil {
		return
	}
	log.Println("raw message:", string(*reply))
	var ret TransactionDetail
	err = json.Unmarshal(*reply, &ret)
	if err != nil {
		return
	}
	if ret.Result.Status != "success" {
		err = errors.New(ret.Result.Error_message)
		return
	}
	retResult = ret.Result
	// signed = ret.Result.Tx_blob

	return
}
func (cli *RPCClient) doPost(url string, r *request) (rpcResp *json.RawMessage, err error) {
	if url == "" {
		err = errors.New("URL is empty")
		return
	}
	payload, err := json.Marshal(r)
	if err != nil {
		return
	}
	p := bytes.NewBuffer(payload)
	log.Println("content:", p.String())
	req, err := http.NewRequest("POST", url, p)
	// req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if cli.user != "" {
		req.Header.Add("Authorization", "Basic "+basicAuth(cli.user, cli.password))
	}

	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		return
	}
	return
}
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
