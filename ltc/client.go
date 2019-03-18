package ltc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Client ...
type Client struct {
	*rpcclient.Client
	httpClient *http.Client
	User       string
	Password   string
	URL        string
}

// New return new rpc client
func New(connect string, port int, user, password string) (*Client, error) {
	conn := &rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", connect, port),
		User:         user,
		Pass:         password,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	c, err := rpcclient.New(conn, nil)
	if err != nil {
		return nil, err
	}
	return &Client{c, &http.Client{}, user, password, fmt.Sprintf("http://%s:%d", connect, port)}, nil
}

func (c *Client) SendBtc(fromAccount, to string, amount string) (string, error) {
	// cmd, err := btcjson.NewCmd("omni_send", from, to, propertyid, amount)
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", err
	}
	// cmd, err := btcjson.NewCmd("sendtoaddress", to, f)
	cmd, err := btcjson.NewCmd("sendfrom", fromAccount, to, f)
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
}
func (c *Client) ListTransaction(addr string) ([]string, error) {
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
	for _, v := range result {
		if v.Address == addr {
			return v.TxIDs, nil
		}
	}
	return nil, errors.New("not found trasaction")
}
func (c *Client) GetTransaction(transactionId string) (*btcjson.GetTransactionResult, error) {
	// btcjson.NewListReceivedByAddressCmd(btcjson.Int(6), nil, nil)
	// cmd, err := btcjson.NewCmd("sendtoaddress", to, f)
	// cmd, err := btcjson.NewCmd("sendfrom", fromAccount, to, f)
	cmd, err := btcjson.NewCmd("gettransaction", transactionId) //minconf include_empty include_watchonly
	if err != nil {
		return nil, err
	}
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return nil, err
	}
	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	var result btcjson.GetTransactionResult
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func (c *Client) GetBalanceOfAddr(addr string) (balance float64, err error) {
	cmd, err := btcjson.NewCmd("listunspent", 1, 9999999, []string{addr})
	if err != nil {
		return
	}
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return
	}
	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	var result []btcjson.ListUnspentResult
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return
	}
	for _, v := range result {
		balance += v.Amount
	}
	return
}

//settxfee
func (c *Client) Settxfee(account string) (bool, error) {
	f, err := strconv.ParseFloat(account, 64)
	if err != nil {
		return false, err
	}
	cmd, err := btcjson.NewCmd("settxfee", f)
	if err != nil {
		return false, err
	}
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return false, err
	}
	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	var result bool
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

//walletpassphrase "passphrase" timeout
func (c *Client) Walletpassphrase(pass string) (bool, error) {
	cmd, err := btcjson.NewCmd("walletpassphrase", pass, 3)
	if err != nil {
		return false, err
	}
	_, err = c.sendCmd(cmd)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (c *Client) Send(fromAccount, to, amount, fee, walletPass string) (string, error) {
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
		success, err := c.Walletpassphrase(walletPass)
		if err != nil {
			return "", err
		}
		if !success {
			return "", errors.New("WalletPass error!")
		}
	}

	h, err := c.SendBtc(fromAccount, to, amount)
	{
		c.WalletLock()
		c.Settxfee("0.001")
	}
	if err != nil {
		return "", err
	}
	return h, nil
}

// OmniGetNewAddress ...
func (c *Client) GetNewAddress(account string) (string, error) {
	cmd, err := btcjson.NewCmd("getnewaddress", account) //每个账号分配一个地址
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
}
func (c *Client) GetBalance(account string) (float64, error) {
	// cmd, err := btcjson.NewCmd("getnewaddress", account)
	cmd, err := btcjson.NewCmd("getbalance", account)
	// cmd, err := btcjson.NewCmd("getbalance")
	if err != nil {
		return 0, err
	}
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return 0, err
	}
	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
	var result float64
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
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
