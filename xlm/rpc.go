package xlm

import (
	// "../aes"
	"fmt"
	// "github.com/stellar/go/keypair"
	"github.com/stellar/go/clients/horizon"
	// "io/ioutil"
	"log"
	"net/http"
	// "strconv"
	"github.com/stellar/go/build"
	// hProtocol "github.com/stellar/go/protocols/horizon"
	// "context"
	"encoding/json"
	// "github.com/stellar/go/xdr"
	"bytes"
	// "errors"
	"time"

	"io"
)

type RPCClient struct {
	Url    string
	Client *horizon.Client
	Seed   string
	Public string
	client *http.Client
}

func NewRPCClient(url, seed, public string) *RPCClient {
	c := &horizon.Client{
		URL:  url,
		HTTP: http.DefaultClient,
	}
	rpcClient := &RPCClient{Url: url, Client: c, Seed: seed, Public: public}
	timeoutIntv := MustParseDuration("20s")
	rpcClient.client = &http.Client{
		Timeout: timeoutIntv,
	}
	return rpcClient
}
func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}
func (r *RPCClient) GetBalance(address string) (balance string, err error) {
	if address == "" {
		address = r.Public
	}
	account, err := r.Client.LoadAccount(address)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Balances for account:", address)

	for _, b := range account.Balances {
		log.Println(b)
		balance = b.Balance
	}
	return
}
func (r *RPCClient) GetTransactions(address, cursor string) (ret []TransactionResponse, err error) {
	url := fmt.Sprintf("%s/accounts/%s/payments?cursor=%s&limit=100", r.Url, address, cursor)
	// limit,order
	var payment OffersPage
	err = r.doPost(url, &payment)
	if err != nil {
		log.Println(err)
		return
	}
	j, err := json.Marshal(payment)
	if err != nil {
		log.Println("marshal:", err)
		return
	}
	log.Println(string(j))
	for _, v := range payment.Embedded.Records {
		var temp TransactionResponse
		if v.Type == "payment" {
			temp.From = v.From
			temp.To = v.To
			temp.Amount = v.Amount
			temp.Created_at = v.Created_at
			temp.Transaction_hash = v.Transaction_hash
			temp.Cursor = v.ID
			res, errs := r.Client.HTTP.Get(v.Links.Transaction.Href)
			if err != nil {
				err = errs
				return
			}
			defer res.Body.Close()

			var tempmemo Memo
			json.NewDecoder(res.Body).Decode(&tempmemo)
			temp.Memo = tempmemo.Value
			ret = append(ret, temp)
		}
	}
	return
}
func (r *RPCClient) doPost(url string, out interface{}) error {
	log.Println("get content:", url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// var rpcResp []byte
	// err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	// if err != nil {
	// 	return nil, err
	// }
	// if len(rpcResp) == 0 {
	// 	return nil, errors.New("body is null")
	// }
	// log.Println(rpcResp)
	var cnt bytes.Buffer
	_, err = io.Copy(&cnt, resp.Body)
	if err != nil {
		return fmt.Errorf("Copy: %s", err)
	}

	if err := json.Unmarshal(cnt.Bytes(), &out); err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}
	return err
}

func (r *RPCClient) Send(fromseed, public, pass, to, amount, memo string) (ret string, err error) {
	// var mut build.TransactionMutator
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: fromseed},
		// build.Network{pass},
		build.Network{pass},
		build.AutoSequence{SequenceProvider: r.Client},
		build.Payment(
			build.Destination{AddressOrSeed: to},
			build.NativeAmount{Amount: amount},
		),
	)
	if err != nil {
		log.Println(err)
		return
	}

	mut := build.MemoText{memo}
	err = tx.Mutate(mut)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("seq:", tx.TX.SeqNum)
	log.Println("memo:", *(tx.TX.Memo.Text))
	txe, err := tx.Sign(fromseed)
	if err != nil {
		log.Println(err)
		return
	}

	txeB64, err := txe.Base64()

	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("tx base64: %s", txeB64)
	// Output: tx base64: AAAAADZY/nWY0gx6beMpf4S8Ur0qHsjA8fbFtBzBx1cbQzHwAAAAZAAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAALSRpLtCLv2eboZlEiHDSGR6Hb+zZL92fbSdNpObeE0EAAAAAAAAAAB3NZQAAAAAAAAAAARtDMfAAAABA2oIeQxoJl53RMRWFeLB865zcky39f2gf2PmUubCuJYccEePRSrTC8QQrMOgGwD8a6oe8dgltvezdDsmmXBPyBw==

	// client := DefaultPublicNetClient
	// transactionEnvelopeXdr := "AAAAABSxFjMo7qcQlJBlrZQypSqYsHA5hHaYxk5hFXwiehh6AAAAZAAIdakAAABZAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAFLEWMyjupxCUkGWtlDKlKpiwcDmEdpjGTmEVfCJ6GHoAAAAAAAAAAACYloAAAAAAAAAAASJ6GHoAAABAp0FnKOQ9lJPDXPTh/a91xoZ8BaznwLj59sdDGK94eGzCOk7oetw7Yw50yOSZg2mqXAST6Agc9Ao/f5T9gB+GCw=="

	response, err := r.Client.SubmitTransaction(txeB64)
	if err != nil {
		log.Println(err)
		herr, isHorizonError := err.(*horizon.Error)
		if isHorizonError {
			resultCodes, errs := herr.ResultCodes()
			if errs != nil {
				log.Println("failed to extract result codes from horizon response")
				err = errs
				return
			}
			log.Println("resultCodes:", resultCodes)
		}
		return
	}

	log.Println(response)
	ret = response.Hash
	return
}
