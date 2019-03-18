package zcash

import (
	// "bytes"
	"encoding/json"
	// "errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	// "github.com/btcsuite/btcd/rpcclient"
	// "io/ioutil"
	"net/http"
	// "strconv"
	"bytes"
	// "encoding/json"
	"errors"
	rpc "github.com/arithmetric/zcashrpcclient"
	// "github.com/btcsuite/btcd/rpcclient"
	"io/ioutil"
	// "strconv"
	"github.com/arithmetric/zcashrpcclient/zcashjson"
	"log"
	"strconv"
)

type Client struct {
	*rpc.Client
	httpClient *http.Client
	User       string
	Password   string
	URL        string
}

// New return new rpc client
func New(connect string, port int, user, password string) (*Client, error) {
	conn := &rpc.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", connect, port),
		User:         user,
		Pass:         password,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	c, err := rpc.New(conn, nil)
	if err != nil {
		return nil, err
	}
	return &Client{c, &http.Client{}, user, password, fmt.Sprintf("http://%s:%d", connect, port)}, nil
}
func (c *Client) GetNewAddress() (addr string, err error) {
	// addr, err = c.ZGetNewAddress()
	cmd, err := btcjson.NewCmd("getnewaddress", "") //每个账号分配一个地址
	if err != nil {
		return "", err
	}
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return "", err
	}
	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	var result string
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return "", err
	}
	return result, nil
	return
}

func (c *Client) GetBalance(addr string) (amount int64, err error) {
	amounts, err := c.ZGetBalance(addr)
	amount = int64(amounts)
	return
}

// func (c *Client) GetTransaction(operationid string) (ret []zcashjson.ZGetOperationStatusResult, err error) {
// 	ret, err = c.ZGetOperationResult(operationid)
// 	return

// }

func (c *Client) SendBtc(to, amount string) (ret SendResp, err error) {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return
	}
	// cmd, err := btcjson.NewCmd("sendtoaddress", to, f)
	cmd, err := btcjson.NewCmd("sendfrom", "", to, f)
	if err != nil {
		return
	}
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return
	}
	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	var result string
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return
	}
	rets, err := c.GetTransaction(result)
	if err != nil {
		return
	}
	ret.Hash = result
	ret.Fee = rets.Fee
	return
}
func (c *Client) GetTransaction(transactionId string) (ret GetTransactionResp, err error) {
	// btcjson.NewListReceivedByAddressCmd(btcjson.Int(6), nil, nil)
	// cmd, err := btcjson.NewCmd("sendtoaddress", to, f)
	// cmd, err := btcjson.NewCmd("sendfrom", fromAccount, to, f)
	cmd, err := btcjson.NewCmd("gettransaction", transactionId) //minconf include_empty include_watchonly
	if err != nil {
		return
	}
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return
	}
	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	var result btcjson.GetTransactionResult
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return
	}
	// ret.Fee =
	fee := fmt.Sprintf("%.8f", result.Fee)
	if fee[:1] == "-" {
		ret.Fee = fee[1:]
	} else {
		ret.Fee = fee
	}
	ret.BlockTime = result.BlockTime
	for _, v := range result.Details {
		if v.Category == "receive" {
			ret.Address = v.Address
			ret.Amount = v.Amount
		}
	}
	return
}
func (c *Client) ListTransaction(addr string) (ret []btcjson.ListReceivedByAddressResult, err error) {
	// btcjson.NewListReceivedByAddressCmd(btcjson.Int(6), nil, nil)
	// cmd, err := btcjson.NewCmd("sendtoaddress", to, f)
	// cmd, err := btcjson.NewCmd("sendfrom", fromAccount, to, f)
	cmd, err := btcjson.NewCmd("listreceivedbyaddress", 6, true, false) //minconf include_empty include_watchonly
	if err != nil {
		return nil, err
	}
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return nil, err
	}
	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	var result []btcjson.ListReceivedByAddressResult
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return nil, err
	}
	log.Println("listreceivedbyaddress:", result)
	for _, v := range result {
		if v.Address == addr {
			ret = append(ret, v)
		}
	}
	if len(ret) == 0 {
		err = errors.New("no transaction in this address")
	}
	return
	// return nil, errors.New("not found trasaction")
}
func (c *Client) Send(from, to, amount string) (hash SendResp, err error) {
	// package zcashjson

	// type ZSendManyEntry struct {
	// 	Address string  `json:"address"`
	// 	Amount  float64 `json:"amount"`
	// 	Memo    *string `json:"memo"`
	// }
	// ZSendMany(fromAccount string, amounts []zcashjson.ZSendManyEntry) (string, error)
	amo, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return
	}
	entry := zcashjson.ZSendManyEntry{
		Address: to,
		Amount:  amo,
		// Memo:    &memo,
	}

	hashs, err := c.ZSendMany(from, []zcashjson.ZSendManyEntry{entry})
	if err != nil {
		return
	}
	ret, err := c.GetTransaction(hashs)
	if err != nil {
		return
	}
	hash.Hash = hashs
	hash.Fee = ret.Fee
	return
}
func (c *Client) sendCmd(cmd interface{}) ([]byte, error) {
	// Get the method associated with the command.
	method, err := btcjson.CmdMethod(cmd)
	if err != nil {
		return nil, err
	}

	// Marshal the command.
	id := c.NextID()
	marshalledJSON, err := btcjson.MarshalCmd(id, cmd)
	if err != nil {
		return nil, err
	}

	jReq := &jsonRequest{
		id:             id,
		method:         method,
		cmd:            cmd,
		marshalledJSON: marshalledJSON,
	}
	return c.sendRequest(jReq)
}

func (c *Client) sendRequest(jReq *jsonRequest) ([]byte, error) {
	bodyReader := bytes.NewReader(jReq.marshalledJSON)
	httpReq, err := http.NewRequest("POST", c.URL, bodyReader)
	if err != nil {
		return nil, err
	}
	httpReq.Close = true
	httpReq.Header.Set("Content-Type", "application/json")

	// Configure basic access authorization.
	httpReq.SetBasicAuth(c.User, c.Password)

	return c.sendPostRequest(httpReq, jReq)
}

func (c *Client) sendPostRequest(httpReq *http.Request, jReq *jsonRequest) ([]byte, error) {
	httpResponse, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	// Read the raw bytes and close the response.
	respBytes, err := ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		err = fmt.Errorf("error reading json reply: %v", err)
		return nil, err
	}

	var res Response
	err = json.Unmarshal(respBytes, &res)
	if err != nil {
		return nil, err
	}
	if res.Error.Code != 0 || res.Error.Message != "" {
		return nil, errors.New(res.Error.Message)
	}

	return res.Result, nil
}
