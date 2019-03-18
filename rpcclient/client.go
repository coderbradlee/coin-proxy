package rpcclient

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"io/ioutil"
	"net/http"
	"strconv"
	// "reflect"
	"log"
	"net"
	// "net/http"
	// "strconv"
	"time"
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
func (c *Client) OmniSend(from, to string, propertyid int64, amount string) (string, error) {
	cmd, err := btcjson.NewCmd("omni_send", from, to, propertyid, amount)
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

type JSONRpcResp struct {
	Id     *json.RawMessage       `json:"id"`
	Result *json.RawMessage       `json:"result"`
	Error  map[string]interface{} `json:"error,omitempty"`
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}
func (c *Client) doPost2(url string, content string) (*JSONRpcResp, error) {
	// jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	// data, _ := json.Marshal(jsonReq)
	log.Println("content:", content)
	data := []byte(content)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "text/plain")
	// req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.User, c.Password)
	timeoutIntv := MustParseDuration("30s")
	cli := &http.Client{
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
			TLSHandshakeTimeout:   30 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
			DisableKeepAlives:     true,
		},
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp *JSONRpcResp
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		return nil, err
	}
	if rpcResp.Error != nil {
		return nil, errors.New(rpcResp.Error["message"].(string))
	}
	return rpcResp, err
}
func (c *Client) SendBtcMany(url, fromAccount string, to string) (string, error) {
	var reqId uint16
	binary.Read(rand.Reader, binary.LittleEndian, &reqId)
	// fmt.Printf("%x\n", reqId)
	reqIdString := fmt.Sprintf("%d", reqId)
	// {"jsonrpc": "1.0", "id":"curltest", "method": "sendmany", "params": ["mr-usdt", {"mt71cPxGGAgsyqH1h1xpVRwR6q3quX2VcQ":0.0001,"mgQ8aK8FLbKju1fFoUCfeKn5rxJ95WXN2d":0.0001}, 6] }
	content := `{"jsonrpc": "1.0", "id":"` + reqIdString + `", "method": "sendmany", "params": ["` + fromAccount + `", ` + to + `, 6] }`

	rpcResp, err := c.doPost2(url, content)
	if err != nil {
		return "", err
	}
	if rpcResp.Result != nil {
		var reply string
		err = json.Unmarshal(*rpcResp.Result, &reply)
		return reply, err
	}
	return "", errors.New("result is nil")
}

// func (c *Client) SendBtcMany(fromAccount string, to map[string]float64) (string, error) {
// 	// func (c *Client) SendBtcMany(fromAccount string, to string) (string, error) {
// 	// cmd, err := btcjson.NewCmd("omni_send", from, to, propertyid, amount)
// 	// f, err := strconv.ParseFloat(amount, 64)
// 	// if err != nil {
// 	// 	return "", err
// 	// }
// 	// &btcjson.SendManyCmd{
// 	// 	FromAccount: "from",
// 	// 	Amounts:     map[string]float64{"1Address": 0.5},
// 	// 	MinConf:     btcjson.Int(1),
// 	// 	Comment:     nil,
// 	// cmd, err := btcjson.NewCmd("sendtoaddress", to, f)
// 	cmd, err := btcjson.NewCmd("sendmany", fromAccount, to)
// 	// NewSendManyCmd(fromAccount string, amounts map[string]float64, minConf *int, comment *string) *SendManyCmd
// 	// cmd := btcjson.NewSendManyCmd(fromAccount, to, btcjson.Int(6), nil)
// 	if err == nil {
// 		// return "", errors.New("cmd is nil")
// 		return "", err
// 	}
// 	resBytes, err := c.sendCmd(cmd)
// 	if err != nil {
// 		return "", err
// 	}
// 	// {"result":"1LXdpMTAQbUvV3EStdCb36Py6mMseYYXfd","error":null,"id":"curltext"}
// 	var result string
// 	err = json.Unmarshal(resBytes, &result)
// 	if err != nil {
// 		return "", err
// 	}
// 	return result, nil
// }

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
func (c *Client) WalletLock() (bool, error) {
	cmd, err := btcjson.NewCmd("walletlock")
	if err != nil {
		return false, err
	}
	_, err = c.sendCmd(cmd)
	if err != nil {
		return false, err
	}
	return true, nil
}

// OmniGetNewAddress ...
func (c *Client) OmniGetNewAddress(account string) (string, error) {
	cmd, err := btcjson.NewCmd("getnewaddress", account)
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
	// cmd, err := btcjson.NewCmd("getbalance", account)
	cmd, err := btcjson.NewCmd("getbalance")
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

// OmniGettransaction ...
func (c *Client) OmniGettransaction(txHash string) (*Transaction, error) {
	cmd := NewGetTransactionCmd(txHash)
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return nil, err
	}

	var result Transaction
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// OmniListBlockTransactions ...
func (c *Client) OmniListBlockTransactions(index int64) ([]string, error) {
	cmd := NewListBlockTransactionsCmd(index)
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return nil, err
	}

	var result []string
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// OmniGetBalance ...
func (c *Client) OmniGetBalance(address string, propertyID int64) (*Balance, error) {
	cmd := NewGetBalanceCmd(address, propertyID)
	resBytes, err := c.sendCmd(cmd)
	if err != nil {
		return nil, err
	}

	var result Balance
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
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
