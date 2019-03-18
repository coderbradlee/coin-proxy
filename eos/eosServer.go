package eos

import (
	// "encoding/json"
	"encoding/json"
	// "io"
	// "log"
	// "net"
	"net/http"

	"strings"
	// "sync"
	// "sync/atomic"
	// "time"
	"encoding/hex"
	"fmt"
	// "fmt"
	"crypto/rand"
	"encoding/binary"
	"errors"
	eosgo "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"io/ioutil"
	"log"
	"strconv"
)

type EosServer struct {
	API_METHOD string
	API_URL    string
	API_PORT   string
	api        *eosgo.API
}

func New(API_METHOD, API_URL, API_PORT string) *EosServer {
	api := &EosServer{
		API_METHOD: API_METHOD,
		API_URL:    API_URL,
		API_PORT:   API_PORT,
		api:        eosgo.New(API_METHOD + "://" + API_URL + ":" + API_PORT),
	}

	return api
}
func GetNewAddress(pass string) (string, error) {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	fmt.Printf("%x\n", n)
	//成功后保存文件，地址为文件名，内容为地址+密码的hash值
	hash := sha3.NewKeccak256()
	address := fmt.Sprintf("%x", n)
	var buf []byte
	//hash.Write([]byte{0xcc})
	hash.Write([]byte(address + pass))
	buf = hash.Sum(buf)

	// fmt.Println(hex.EncodeToString(buf))
	// d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile("keys/eos/"+address, []byte(hex.EncodeToString(buf)), 0644)

	if err != nil {
		log.Println("write file err:", err)
		return "", err
	}
	return address, nil
}

func (s *EosServer) GetBalance(account, symbol string) (string, error) {
	ass, err := s.api.GetCurrencyBalance(eosgo.AccountName(account), symbol, "eosio.token")
	if err != nil {
		log.Println(err)
		return "", err
	}
	// for i, v := range ass {
	// 	fmt.Println(i, ":", v)
	// }
	if len(ass) > 0 {
		// return fmt.Sprintf("%d", ass[0].Amount), nil
		return ass[0].String(), nil
	} else {
		return "0", errors.New("reply is empty")
	}
}
func (s *EosServer) GetAccount(account string) (*eosgo.AccountResp, error) {
	return s.api.GetAccount(eosgo.AccountName(account))

}
func (s *EosServer) dealWithAmount(amount, symbol string) string {
	if !strings.Contains(amount, ".") {
		return amount + ".0000 " + symbol
	}
	out := strings.Split(amount, ".")
	le := len(out[1])
	if le < 4 {
		for i := 0; i < 4-le; i++ {
			amount += "0"
		}
	}
	return amount + " " + symbol
}
func (s *EosServer) transfer(from, to, amount, userMemo, localurl, localwallet, localpriv, symbol string) (string, error) {
	m := make(eosgo.M)
	// "from": "eosio",
	// "to": "noprom",
	// "quantity": "1.0000 EOS",
	// "memo": "created by noprom"
	m["from"] = from //100eos
	m["to"] = to     //0eos
	amountOut := s.dealWithAmount(amount, symbol)
	m["quantity"] = amountOut // amount + " EOS"
	m["memo"] = userMemo
	ass, err := s.api.ABIJSONToBin("eosio.token", "transfer", m)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// fmt.Println(hex.EncodeToString(ass))

	//////////////////////////////////////////////////////////////
	// Account       AccountName       `json:"account"`
	// Name          ActionName        `json:"name"`
	// Authorization []PermissionLevel `json:"authorization,omitempty"`
	// ActionData
	//
	actionData := eosgo.NewActionDataFromHexData(ass)
	pl, err := eosgo.NewPermissionLevel(from + "@active")
	if err != nil {
		log.Println(err)
		return "", err
	}
	action := eosgo.Action{
		Account:       "eosio.token",
		Name:          "transfer",
		Authorization: []eosgo.PermissionLevel{pl},
		ActionData:    actionData,
	}
	//////(api *API) SetSigner(s Signer)
	// NewWalletSigner(api *API, walletName string) *WalletSigner
	// localeosurl := cfg.Eos.API_METHOD + "://" + cfg.Eos.LOCAL_API_URL + ":" + strconv.Itoa(cfg.Eos.LOCAL_API_PORT)
	localapi := eosgo.New(localurl)
	// err = localapi.WalletUnlock(cfg.Eos.WALLET_NAME, cfg.Eos.WALLET_PRIV_KEY)
	err = localapi.WalletUnlock(localwallet, localpriv)
	if err != nil {
		log.Println(err)
		// return "", err
	}
	signer := eosgo.NewWalletSigner(localapi, localwallet)
	s.api.SetSigner(signer)
	///////SignPushActions(a ...*Action) (out *PushTransactionFullResp///
	resp, err := s.api.SignPushActions(&action)
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println("TransactionID:", resp.TransactionID)
	///////////////////////////////////////
	err = localapi.WalletLock(localwallet)
	if err != nil {
		log.Println(err)
		// return "", err
	}
	return resp.TransactionID, nil
}
func (s *EosServer) Delegatebw(from, to, stakeCPU, stakeNet, localurl, localwallet, walletPassword string) (ret string, err error) {
	// NewDelegateBW(from, receiver eos.AccountName, stakeCPU, stakeNet eos.Asset, transfer bool) *eos.Action
	fromA := eosgo.AccountName(from)
	toA := eosgo.AccountName(to)
	cpu, err := eosgo.NewEOSAssetFromString(stakeCPU)
	if err != nil {
		return
	}
	net, err := eosgo.NewEOSAssetFromString(stakeNet)
	if err != nil {
		return
	}
	action := system.NewDelegateBW(fromA, toA, cpu, net, false) //最后一个参数是指以租借还是赠送的方式，true为赠送方式，这时from不能等于to

	localapi := eosgo.New(localurl)
	// err = localapi.WalletUnlock(cfg.Eos.WALLET_NAME, cfg.Eos.WALLET_PRIV_KEY)
	err = localapi.WalletUnlock(localwallet, walletPassword)
	if err != nil {
		log.Println(err)
		// return "", err
	}
	signer := eosgo.NewWalletSigner(localapi, localwallet)
	s.api.SetSigner(signer)
	///////SignPushActions(a ...*Action) (out *PushTransactionFullResp///
	resp, err := s.api.SignPushActions(action)
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println("TransactionID:", resp.TransactionID)
	///////////////////////////////////////
	err = localapi.WalletLock(localwallet)
	if err != nil {
		log.Println(err)
		// return "", err
	}
	ret = resp.TransactionID
	return
}
func (s *EosServer) Undelegatebw(from, to, stakeCPU, stakeNet, localurl, localwallet, walletPassword string) (ret string, err error) {
	// NewDelegateBW(from, receiver eos.AccountName, stakeCPU, stakeNet eos.Asset, transfer bool) *eos.Action
	fromA := eosgo.AccountName(from)
	toA := eosgo.AccountName(to)
	cpu, err := eosgo.NewEOSAssetFromString(stakeCPU)
	if err != nil {
		return
	}
	net, err := eosgo.NewEOSAssetFromString(stakeNet)
	if err != nil {
		return
	}
	action := system.NewUndelegateBW(fromA, toA, cpu, net)
	localapi := eosgo.New(localurl)
	// err = localapi.WalletUnlock(cfg.Eos.WALLET_NAME, cfg.Eos.WALLET_PRIV_KEY)
	err = localapi.WalletUnlock(localwallet, walletPassword)
	if err != nil {
		log.Println(err)
		// return "", err
	}
	signer := eosgo.NewWalletSigner(localapi, localwallet)
	s.api.SetSigner(signer)
	///////SignPushActions(a ...*Action) (out *PushTransactionFullResp///
	resp, err := s.api.SignPushActions(action)
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println("TransactionID:", resp.TransactionID)
	///////////////////////////////////////
	err = localapi.WalletLock(localwallet)
	if err != nil {
		log.Println(err)
		// return "", err
	}
	ret = resp.TransactionID
	return
}
func (s *EosServer) checkPass(memo, pass string) error {
	b, err := ioutil.ReadFile("keys/eos/" + memo)
	if err != nil {
		// fmt.Print(err)
		log.Println("read file:", err)
		return errors.New("open file error!")
	}
	str := string(b)
	log.Println(memo+pass, ":", str)

	hash := sha3.NewKeccak256()

	var buf []byte
	//hash.Write([]byte{0xcc})
	hash.Write([]byte(memo + pass))
	buf = hash.Sum(buf)
	if str != hex.EncodeToString(buf) {
		log.Println("password error:", str, "!=", hex.EncodeToString(buf))
		return errors.New("wrong password!")
	}
	return nil
}

type ActionResult struct {
	From          string `json:"from,omitempty"`
	To            string `json:"to,omitempty"`
	Quantity      string `json:"quantity,omitempty"`
	TransactionId string `json:"transactionId,omitempty"`
	Memo          string `json:"memo"`
	Seq           string `json:"seq,omitempty"`
	Packed        bool   `json:"packed,omitempty"`
}

//返回此笔交易的转入金额
// VerifyIn(account, seq, symbol)
func (s *EosServer) VerifyAll(inAccount, seq, symbol string) ([]ActionResult, error) {
	bestPosString := seq
	bestPos, err := strconv.ParseInt(bestPosString, 10, 64)
	if err != nil {
		return nil, err
	}
	{
		//获取最新的交易sequence
		req := eosgo.GetActionsRequest{
			AccountName: eosgo.AccountName(inAccount),
			Pos:         -1,
			Offset:      -1,
		}

		actions, err := s.api.GetActions(req)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		for _, v := range actions.Actions {
			// bestPos = v.Trace.Receipt.GlobalSequence
			if v.AccountSeq > bestPos {
				bestPos = v.AccountSeq
			}
		}
	}
	{
		seqence := seq

		seqInt, err := strconv.ParseInt(seqence, 10, 64)
		if err != nil {
			return nil, err
		}
		if bestPos < seqInt {
			bestPos = seqInt
		}
		offset := bestPos - seqInt
		if offset > 99 {
			offset = 99
		}
		log.Println("seqInt:", seqInt, " bestpos:", bestPos, " offset:", offset)
		// 	pos INT                     sequence number of action for this account, -1 for last
		//   offset INT                  get actions [pos,pos+offset] for positive offset or [pos-offset,pos) for negative offset
		req := eosgo.GetActionsRequest{
			AccountName: eosgo.AccountName(inAccount),
			// Pos:         -1,
			// Offset:      seqInt - 1,
			Pos:    seqInt,
			Offset: offset,
		}

		actions, err := s.api.GetActions(req)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		log.Println("actions len:", len(actions.Actions))
		var ret []ActionResult
		for _, v := range actions.Actions {
			actData := v.Trace.Action.ActionData.Data
			log.Println("one action:", actData)
			if actData == nil {
				log.Println("actData is nil")
				continue
			}
			dataMap, ok := actData.(map[string]interface{})
			if !ok {
				log.Println("dataMap not ok")
				continue
			}
			from, ok := dataMap["from"].(string)
			if !ok {
				log.Println("from not ok:")
				continue
			}
			to, ok := dataMap["to"].(string)
			if !ok {
				log.Println("to not ok:")
				continue
			}
			quantity, ok := dataMap["quantity"].(string)
			if !ok {
				log.Println("quantity not ok:")
				continue
			}
			memos, ok := dataMap["memo"].(string)
			if !ok {
				log.Println("memo not ok:")
				continue
			}
			log.Println("from:", from, " to:", to, " quantity:", quantity, " memo:", memos)
			out := strings.Split(quantity, " ")

			// log.Println("symbol:", out[1])

			if out[1] == symbol {
				// return quantity, nil
				// bestTransactionID = v.Trace.TransactionID.String()
				// bestQuantity = quantity
				temp := ActionResult{
					From:          from,
					To:            to,
					Quantity:      quantity,
					TransactionId: v.Trace.TransactionID.String(),
					Memo:          memos,
					Seq:           fmt.Sprintf("%d", v.AccountSeq),
				}
				retLen := len(ret)

				if (retLen > 0) && (temp.TransactionId == ret[retLen-1].TransactionId) {
					// ret = delete(ret) 去除重复的交易
					ret = ret[:retLen-1]
				}
				ret = append(ret, temp)
			}
		}
		// if bestTransactionID != "" && bestQuantity != "" {
		// 	return bestTransactionID, bestQuantity, nil
		// }
		if len(ret) == 0 {
			return nil, errors.New("not found transaction")
		}
		return ret, nil
	}
}
func (s *EosServer) Verify(url, inAccount, page string) (result []RetTransactionHistoryDetail, err error) {
	log.Println("url in config:", url)
	all := fmt.Sprintf(url, inAccount, page)
	ret, err := s.doPost(all)
	if err != nil {
		return
	}
	result = ret.Data.Trace_list
	return
}

// func (s *EosServer) GetTransaction(transactionId string) (trs *eosgo.TransactionResp2, err error) {
func (s *EosServer) GetTransaction(transactionId string) (ret ActionResult, err error) {
	// GetTransaction(id string) (out *TransactionResp, err error)
	trs, err := s.api.GetTransaction2(transactionId)
	// outtrx, _ := json.Marshal(trs)
	// fmt.Println("json:", string(outtrx))
	if err != nil {
		log.Println(err)
		err = errors.New("transaction not found!")
		return
	}
	ret.Packed = false
	if trs.Trans.Receipt.Status != eosgo.TransactionStatusExecuted {
		err = errors.New("transaction still not executed!")
		return
	}
	ret.Packed = true
	// trs.trx
	// test := trs.Receipt.PackedTransaction.ID.String()
	// fmt.Println("test:", test)
	// return trs.ID.String(), nil
	actData := trs.Trans.TrxInTrx.Actions[0].Data
	dataMap, ok := actData.(map[string]interface{})
	if !ok {
		log.Println("dataMap not ok")
		err = errors.New("dataMap not ok")
		return
	}
	ret.From, ok = dataMap["from"].(string)
	if !ok {
		log.Println("from not ok:")
		err = errors.New("from not ok")
		return
	}
	ret.To, ok = dataMap["to"].(string)
	if !ok {
		log.Println("to not ok:")

		err = errors.New("to not ok")
		return
	}
	ret.Quantity, ok = dataMap["quantity"].(string)
	if !ok {
		log.Println("quantity not ok:")
		err = errors.New("quantity not ok")
		return
	}
	ret.Memo, ok = dataMap["memo"].(string)
	if !ok {
		log.Println("memo not ok:")
		err = errors.New("memo not ok")
		return
	}
	return
}

// Send(to, memo, pass, amount)
// (from, to, amount, userMemo, localeosurl, cfg.Eos.WALLET_NAME, cfg.Eos.WALLET_PRIV_KEY, symbol)
func (s *EosServer) Send(from, to, amount, userMemo, localurl, localwallet, localpriv, symbol string) (string, error) {
	// err := s.checkPass(memo, pass)
	// if err != nil {
	// 	return "", err
	// }

	/////转账
	return s.transfer(from, to, amount, userMemo, localurl, localwallet, localpriv, symbol)
}

// {"errno":3,"errmsg":"error"}
// {"errno":0,"errmsg":"Success","data":{"trace_count":0,"trace_list":[{"trx_id":"ac1febc1649fd27a8c26fc57ceb2294be135685b6b956915d3961e5616eb8206","timestamp":"2018-11-09T06:55:51.000","receiver":"monstereos12","sender":"gateiowallet","code":"eosio.token","quantity":"9.9000","memo":"","symbol":"EOS","status":"executed"},{"trx_id":"34d3fd3fe4ff21d2a770730fe4dd0c372ae754754f3b873925a3c2bfd79cc528","timestamp":"2018-11-09T06:33:09.500","receiver":"gateiowallet","sender":"monstereos12","code":"eosio.token","quantity":"0.0100","memo":"000867ea449f73c3","symbol":"EOS","status":"executed"},{"trx_id":"99dd8880518bb57a0f9edf0fee48315b451955bd050d12857b7d5a8dfd911585","timestamp":"2018-11-09T06:29:12.000","receiver":"monstereos12","sender":"gateiowallet","code":"eosio.token","quantity":"0.9000","memo":"","symbol":"EOS","status":"executed"}]}}
type RetTransactionHistoryDetail struct {
	Trx_id string `json:"trx_id,omitempty"`

	Timestamp string `json:"timestamp,omitempty"`
	Receiver  string `json:"receiver,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Code      string `json:"code,omitempty"`
	Quantity  string `json:"quantity,omitempty"`
	Memo      string `json:"memo,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Status    string `json:"status,omitempty"`
	// "trx_id":"ac1febc1649fd27a8c26fc57ceb2294be135685b6b956915d3961e5616eb8206","timestamp":"2018-11-09T06:55:51.000",
	// "receiver":"monstereos12",
	// "sender":"gateiowallet",
	// "code":"eosio.token",
	// "quantity":"9.9000",
	// "memo":"",
	// "symbol":"EOS",
	// "status":"executed"
}
type RetTransactionHistoryd struct {
	Trace_count int                           `json:"trace_count,omitempty"`
	Trace_list  []RetTransactionHistoryDetail `json:"trace_list,omitempty"`
}
type RetTransactionHistory struct {
	Errno  int                    `json:"errno,omitempty"`
	Errmsg string                 `json:"errmsg,omitempty"`
	Data   RetTransactionHistoryd `json:"data,omitempty"`
}

func (r *EosServer) doPost(url string) (reply RetTransactionHistory, err error) {
	log.Println("url:", url)
	// jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	// data, _ := json.Marshal(jsonReq)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println(err)
		return
	}

	cli := &http.Client{}
	// req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	//log.Println("++++++",req.Body)
	resp, err := cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return
	}
	if reply.Errno != 0 {
		err = errors.New(reply.Errmsg)
	}
	return
}
