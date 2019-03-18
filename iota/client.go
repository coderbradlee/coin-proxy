package iota

import (
	// "bytes"
	"encoding/json"
	"errors"
	"fmt"
	// "github.com/iotaledger/iota.go/address"
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/bundle"
	"github.com/iotaledger/iota.go/consts"
	// "github.com/iotaledger/iota.go/converter"
	"../aes"
	"context"
	"github.com/iotaledger/iota.go/pow"
	"github.com/iotaledger/iota.go/transaction"
	"github.com/iotaledger/iota.go/trinary"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	// *rpc.Client
	// httpClient *http.Client
	client *api.API
	// User     string
	// Password string
	URL  string
	Seed string
}

// New return new rpc client
func New(connect string, port, method, seed string) (*Client, error) {

	url := fmt.Sprintf("%s://%s:%s", method, connect, port)
	// url := "http://" + Host
	_, proofOfWorkFunc := pow.GetFastestProofOfWorkImpl()
	httpClient := &http.Client{Timeout: time.Duration(20) * time.Second}
	iotaAPI, err := api.ComposeAPI(api.HTTPClientSettings{URI: url, LocalProofOfWorkFunc: proofOfWorkFunc, Client: httpClient})
	// iotaAPI, err := api.ComposeAPI(api.HTTPClientSettings{URI: url, Client: httpClient})
	if err != nil {
		return nil, err
	}

	// iotaAPI, err = api.ComposeAPI(api.HTTPClientSettings{
	// 	URI:                  endpoint,
	// 	LocalProofOfWorkFunc: powFunc,
	// 	Client:               httpClient,
	// })
	ret := &Client{
		client: iotaAPI,
		// User:     user,
		// Password: password,
		URL:  url,
		Seed: seed,
	}
	return ret, nil
}
func (c *Client) GetNewAddress(ind, password, dir, seed string) (ret string, err error) {
	// seed := "HZVEINVKVIKGFRAWRTRXWD9JLIYLCQNCXZRBLDETPIQGKZJRYKZXLTV9JNUVBIAHAGUZVIQWIAWDZ9ACW"
	index, err := strconv.ParseUint(ind, 10, 64)
	if err != nil {
		return
	}
	log.Println("index:", index)
	addr, err := c.client.GetNewAddress(c.Seed, api.GetNewAddressOptions{Index: index, Security: consts.SecurityLevelMedium, Checksum: true, ReturnAll: true})
	if len(addr) > 0 {
		ret = addr[0]
	}
	err = c.writeFile(dir, seed, ret, password)
	return

}
func (r *Client) writeFile(dir, password, filename, content string) (err error) {
	key := []byte(password)
	result, err := aes.AesEncrypt([]byte(content), key)
	if err != nil {
		return
	}
	// fmt.Println(base64.StdEncoding.EncodeToString(result))
	err = ioutil.WriteFile(dir+"/"+filename, result, 0644)

	if err != nil {
		log.Println("write file err:", err)
		return err
	}
	return nil
}
func (c *Client) GetBalance(addr string) (ret uint64, err error) {
	balances, err := c.client.GetBalances(trinary.Hashes{addr}, 100)
	if err != nil {
		// handle error
		return
	}
	balan := balances.Balances
	log.Println(addr, " balance:", balan)
	if len(balan) > 0 {
		ret = balan[0]
	}
	return
}
func (c *Client) GetBalance2(addr string) (ret uint64, milestonesIndex int64, err error) {
	balances, err := c.client.GetBalances(trinary.Hashes{addr}, 100)
	if err != nil {
		// handle error
		return
	}
	balan := balances.Balances
	log.Println(addr, " balance:", balan)
	if len(balan) > 0 {
		ret = balan[0]
	}
	milestonesIndex = balances.MilestoneIndex
	return
}
func (c *Client) ListTransaction(addr string) (hash []string, err error) {
	// balances, err := c.client.GetBalances(trinary.Hashes{addr}, 100)
	hash, err = c.client.FindTransactions(api.FindTransactionsQuery{Addresses: []string{addr}})

	return
}

type RetBundle struct {
	Hash      string `json:"hash,omitempty"`
	Bundle    string `json:"bundle,omitempty"`
	Amount    int64  `json:"amount,omitempty"`
	Timestamp uint64 `json:"timestamp,omitempty"`
}

func (c *Client) ListBundle(addr string) (ret []RetBundle, err error) {
	// balances, err := c.client.GetBalances(trinary.Hashes{addr}, 100)
	bundles, err := c.client.GetBundlesFromAddresses(trinary.Hashes{addr}, true)
	var allhash trinary.Hashes
	for _, v := range bundles {
		for _, iv := range v {
			allhash = append(allhash, iv.Hash)
		}
	}
	confirmed, err := c.client.GetLatestInclusion(allhash)
	if err != nil {
		return
	}
	confirmedMap := make(map[string]bool)
	for i, v := range allhash {
		confirmedMap[v] = confirmed[i]
	}
	bundlesbytes, _ := json.Marshal(bundles)
	log.Println("bundles:", string(bundlesbytes))
	for _, v := range bundles {
		for _, iv := range v {
			// log.Println("v:", iv)
			if (iv.Address == addr) && (iv.Value > 0) {
				if confirmedMap[iv.Hash] {
					temp := RetBundle{Hash: iv.Hash, Bundle: iv.Bundle, Amount: iv.Value, Timestamp: iv.Timestamp}
					ret = append(ret, temp)
				}
			}
		}
	}
	if len(ret) == 0 {
		err = errors.New("no bundle with this address")
	}
	return
}
func (c *Client) GetTransaction(txid string) (trx transaction.Transactions, err error) {
	// balances, err := c.client.GetBalances(trinary.Hashes{addr}, 100)
	// hash, err = c.client.FindTransactions(api.FindTransactionsQuery{Addresses: []string{addr}})
	// trytes, err := c.client.GetTrytes(txid)
	// if err != nil {
	// 	// handle error
	// 	return
	// }
	// tryte = trytes[0]
	// log.Println("txid:", txid, ",trytes:", tryte)
	// // log.Println(len(trytes))
	// ascii, err := converter.TrytesToASCII(tryte)
	// if err != nil {
	// 	// handle error
	// 	return
	// }
	// log.Println(ascii) // output: IOTA
	// func (api *API) GetTransactionObjects(hashes ...Hash) (transaction.Transactions, error) {
	// 	if err := Validate(ValidateTransactionHashes(hashes...)); err != nil {
	// 		return nil, err
	// 	}
	// 	trytes, err := api.GetTrytes(hashes...)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return transaction.AsTransactionObjects(trytes, hashes)
	// }
	trx, err = c.client.GetTransactionObjects(txid)
	if err != nil {
		// handle error
		return
	}
	return
}
func (c *Client) ReplayBundle(tailTxHash string) (hash string, err error) {
	// ReplayBundle(tailTxHash Hash, depth uint64, mwm uint64, reference ...Hash) (bundle.Bundle, error)
	confirmed, err := c.IsConfirmed(tailTxHash)
	if err != nil {
		log.Println(err)
		return
	}
	if confirmed {
		log.Println("already confirmed!")
		return
	}
	log.Println("not confirmed,need replaybundle")
	bndle, err := c.client.ReplayBundle(tailTxHash, 3, 14)
	if err != nil {
		log.Println(err)
		return
	}
	b := bndle[0].Bundle
	log.Println(b)
	hash = bundle.TailTransactionHash(bndle)
	log.Println("replay bundle with tail tx hash: ", hash)
	return
}
func (c *Client) IsConfirmed(tailTxHash string) (ret bool, err error) {
	// GetLatestInclusion(txHashes Hashes) ([]bool, error)
	rets, err := c.client.GetLatestInclusion([]string{tailTxHash})
	if err != nil {
		log.Println(err)
		return
	}
	if len(rets) <= 0 {
		log.Println("len rets is null")
		err = errors.New("len rets is null")
		return
	}
	ret = rets[0]
	return
}
func (c *Client) CheckPassword(seeds, from, dir, password string) (err error) {
	b, err := ioutil.ReadFile(dir + "/" + from)
	if err != nil {
		// fmt.Print(err)
		log.Println("read file:", err)
		return
	}
	origData, err := aes.AesDecrypt(b, []byte(seeds))
	if err != nil {
		log.Println("AesDecrypt:", err)
		return
	}
	// log.Println("13")
	if len(origData) == 0 {
		err = errors.New("decrpt is null")
		log.Println(err)
		return
	}
	// log.Println("14")
	if password != string(origData) {
		log.Println(string(origData))
		err = errors.New("wrong password")
		return
	}
	return nil
}

type SendRet struct {
	Bundle   string `json:"bundle"`
	TailHash string `json:"tailHash"`
}

func (c *Client) Send(seeds, from, recipientAddress, fromkeyindexs, toamounts, changeAddress, password, dir string) (ret SendRet, err error) {
	err = c.CheckPassword(seeds, from, dir, password)
	if err != nil {
		return
	}
	frombalance, err := c.GetBalance(from)
	if err != nil {
		return
	}
	fromkeyindex, err := strconv.ParseUint(fromkeyindexs, 10, 64)
	if err != nil {
		return
	}
	toamount, err := strconv.ParseUint(toamounts, 10, 64)
	if err != nil {
		return
	}
	if frombalance < toamount {
		err = errors.New("insufficient balance")
		return
	}
	var seed = trinary.Trytes(seeds)
	transfers := bundle.Transfers{
		{
			Address: recipientAddress,
			Value:   toamount,
		},
	}
	// inputs := c.client.GetInputObjects([]string{from}, []uint64{frombalance}, fromkeyindex, consts.SecurityLevelMedium)
	inputs := []api.Input{
		{
			Address:  from,
			Security: consts.SecurityLevelMedium,
			KeyIndex: fromkeyindex,
			Balance:  frombalance,
		},
	}
	prepTransferOpts := api.PrepareTransfersOptions{Inputs: inputs, RemainderAddress: &changeAddress}

	sendoption := api.SendTransfersOptions{PrepareTransfersOptions: prepTransferOpts}

	bndl, err := c.client.SendTransfer(seed, 3, 14, transfers, &sendoption)
	if err != nil {
		log.Println(err)
		return
	}
	bytes, err := json.Marshal(bndl)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(bytes))
	hash := bundle.TailTransactionHash(bndl)
	log.Println("broadcasted bundle with tail tx hash: ", hash)
	ret.TailHash = hash
	if len(bndl) > 0 {
		ret.Bundle = bndl[0].Bundle
	}
	// ReplayBundle(tailTxHash Hash, depth uint64, mwm uint64, reference ...Hash) (bundle.Bundle, error)
	// bndle, err := c.client.ReplayBundle(hash, 3, 14)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// log.Println(bndle[0].Bundle)

	//send, promote and then replay if pending for half an hour
	go func() {
		promotionTransfers := bundle.Transfers{bundle.EmptyTransfer}

		// options for promotion
		delay := time.Duration(5) * time.Second
		// stop promotion after one minute
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(1)*time.Minute)
		opts := api.PromoteTransactionOptions{
			Ctx: ctx,
			// wait for 5 seconds before each promotion
			Delay: &delay,
		}

		// this promotion will stop until the passed in Context is done
		promotionTx, err := c.client.PromoteTransaction(hash, 3, 14, promotionTransfers, opts)
		if err != nil {
			// handle error
			return
		}
		log.Println("promoted tx with new tx:", promotionTx[0].Hash)
	}()

	return
}

// func (c *Client) Send(seeds, from, recipientAddress string, fromkeyindex, frombalance, toamount uint64) (hash string, err error) {
// func (c *Client) Send(seeds, from, recipientAddress, fromkeyindexs, toamounts, changeAddress string) (ret string, err error) {
// 	balan, err := c.GetBalance(from)
// 	if err != nil {
// 		return
// 	}
// 	// if len(balan) == 0 {
// 	// 	err = errors.New("len balan is zero")
// 	// 	return
// 	// }
// 	frombalance := balan
// 	fromkeyindex, err := strconv.ParseUint(fromkeyindexs, 10, 64)
// 	if err != nil {
// 		return
// 	}
// 	toamount, err := strconv.ParseUint(toamounts, 10, 64)
// 	if err != nil {
// 		return
// 	}
// 	if frombalance < toamount {
// 		err = errors.New("insufficient balance")
// 		return
// 	}
// 	var seed = trinary.Trytes(seeds)
// 	const mwm = 9
// 	const depth = 3
// 	transfers := bundle.Transfers{
// 		{
// 			Address: recipientAddress,
// 			Value:   toamount,
// 		},
// 	}

// 	// create inputs for the transfer
// 	inputs := []api.Input{
// 		{
// 			Address:  from,
// 			Security: consts.SecurityLevelMedium,
// 			KeyIndex: fromkeyindex,
// 			Balance:  frombalance,
// 		},
// 	}
// 	prepTransferOpts := api.PrepareTransfersOptions{Inputs: inputs, RemainderAddress: &changeAddress}
// 	// prepare the transfer by creating a bundle with the given transfers and inputs.
// 	// the result are trytes ready for PoW.
// 	trytes, err := c.client.PrepareTransfers(seed, transfers, prepTransferOpts)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	// you can decrease your chance of sending to a spent address by checking the address before
// 	// broadcasting your bundle.
// 	spent, err := c.client.WereAddressesSpentFrom(transfers[0].Address)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	if spent[0] {
// 		err = errors.New("recipient address is spent from, aborting transfer")
// 		log.Println(err)
// 		return
// 	}

// 	// at this point the bundle trytes are signed.
// 	// now we need to:
// 	// 1. select two tips
// 	// 2. do proof-of-work
// 	// 3. broadcast the bundle
// 	// 4. store the bundle
// 	// SendTrytes() conveniently does the steps above for us.
// 	// bndl, err := c.client.SendTrytes(trytes, depth, mwm)
// 	// SendTransfer(seed Trytes, depth uint64, mwm uint64, transfers bundle.Transfers, options *SendTransfersOptions) (bundle.Bundle, error)
// 	opts := c.client.SendTransfersOptions{Reference: &tailTxHash}
// 	opts.PrepareTransfersOptions = getPrepareTransfersDefaultOptions(opts.PrepareTransfersOptions)
// 	bndl, err := c.client.SendTransfer(seed, depth, mwm, transfers, &opts)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	bytes, err := json.Marshal(bndl)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	log.Println(string(bytes))
// 	hash := bundle.TailTransactionHash(bndl)
// 	log.Println("broadcasted bundle with tail tx hash: ", hash)
// 	if len(bndl) > 0 {
// 		ret = bndl[0].Bundle
// 	}
// 	return
// }

// func (c *Client) Send(seeds, from, recipientAddress, fromkeyindexs, toamounts, changeAddress, password, dir string) (ret string, err error) {
// 	// fromshort := from[:81]
// 	// recipientAddressshort := recipientAddress[:81]
// 	// changeAddressshort := changeAddress[:81]
// 	// log.Println("from:", fromshort)
// 	// log.Println("recipientAddressshort:", recipientAddressshort)
// 	// log.Println("changeAddressshort:", changeAddressshort)
// 	b, err := ioutil.ReadFile(dir + "/" + from)
// 	if err != nil {
// 		// fmt.Print(err)
// 		log.Println("read file:", err)
// 		return
// 	}
// 	origData, err := aes.AesDecrypt(b, []byte(seeds))
// 	if err != nil {
// 		log.Println("AesDecrypt:", err)
// 		return
// 	}
// 	// log.Println("13")
// 	if len(origData) == 0 {
// 		err = errors.New("decrpt is null")
// 		return
// 	}
// 	// log.Println("14")
// 	if password != string(origData) {
// 		log.Println(string(origData))
// 		err = errors.New("wrong password")
// 		return
// 	}
// 	// balan, milestonesIndex, err := c.GetBalance2(from)
// 	balan, err := c.GetBalance(from)
// 	if err != nil {
// 		return
// 	}
// 	// if len(balan) == 0 {
// 	// 	err = errors.New("len balan is zero")
// 	// 	return
// 	// }
// 	frombalance := balan
// 	fromkeyindex, err := strconv.ParseUint(fromkeyindexs, 10, 64)
// 	if err != nil {
// 		return
// 	}
// 	toamount, err := strconv.ParseUint(toamounts, 10, 64)
// 	if err != nil {
// 		return
// 	}
// 	if frombalance < toamount {
// 		err = errors.New("insufficient balance")
// 		return
// 	}
// 	// var endpoint = "<node-url>"

// 	// // must be 81 trytes long and truly random
// 	var seed = trinary.Trytes(seeds)

// 	// // difficulty of the proof of work required to attach a transaction on the tangle
// 	// const mwm = 14
// 	const mwm = 14

// 	// // how many milestones back to start the random walk from
// 	const depth = 3

// 	// // can be 90 trytes long (with checksum)
// 	// const recipientAddress = "BBBB....."

// 	// func main() {

// 	// get the best available PoW implementation
// 	// _, proofOfWorkFunc := pow.GetFastestProofOfWorkImpl()

// 	// create a new API instance
// 	// api, err := ComposeAPI(HTTPClientSettings{
// 	// 	URI: endpoint,
// 	// 	// (!) if no PoWFunc is supplied, then the connected node is requested to do PoW for us
// 	// 	// via the AttachToTangle() API call.
// 	// 	LocalProofOfWorkFunc: proofOfWorkFunc,
// 	// })
// 	// must(err)

// 	// create a transfer to the given recipient address
// 	// optionally define a message and tag
// 	transfers := bundle.Transfers{
// 		{
// 			Address: recipientAddress,
// 			Value:   toamount,
// 		},
// 	}

// 	// create inputs for the transfer
// 	// inputs := []api.Input{
// 	// 	{
// 	// 		Address:  fromshort,
// 	// 		Security: consts.SecurityLevelMedium,
// 	// 		// KeyIndex: uint64(milestonesIndex),
// 	// 		KeyIndex: fromkeyindex,
// 	// 		Balance:  frombalance,
// 	// 	},
// 	// }
// 	// GetInputObjects(addresses Hashes, balances []uint64, start uint64, secLvl SecurityLevel) Inputs
// 	inputs := c.client.GetInputObjects([]string{from}, []uint64{frombalance}, fromkeyindex, consts.SecurityLevelMedium)
// 	// inputs := []api.Input{input}
// 	// log.Println("milestoneIndex:", milestonesIndex)
// 	// create an address for the remainder.
// 	// in this case we will have 20 iotas as the remainder, since we spend 100 from our input
// 	// address and only send 80 to the recipient.
// 	// remainderAddress, err := address.GenerateAddress(seed, 1, consts.SecurityLevelMedium)

// 	// we don't need to set the security level or timestamp in the options because we supply
// 	// the input and remainder addresses.
// 	// prepTransferOpts := api.PrepareTransfersOptions{Inputs: inputs, RemainderAddress: &remainderAddress}
// 	prepTransferOpts := api.PrepareTransfersOptions{Inputs: inputs.Inputs, RemainderAddress: &changeAddress}
// 	// prepare the transfer by creating a bundle with the given transfers and inputs.
// 	// the result are trytes ready for PoW.
// 	trytes, err := c.client.PrepareTransfers(seed, transfers, prepTransferOpts)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	// you can decrease your chance of sending to a spent address by checking the address before
// 	// broadcasting your bundle.
// 	spent, err := c.client.WereAddressesSpentFrom(transfers[0].Address)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	if spent[0] {
// 		err = errors.New("recipient address is spent from, aborting transfer")
// 		log.Println(err)
// 		return
// 	}

// 	// at this point the bundle trytes are signed.
// 	// now we need to:
// 	// 1. select two tips
// 	// 2. do proof-of-work
// 	// 3. broadcast the bundle
// 	// 4. store the bundle
// 	// SendTrytes() conveniently does the steps above for us.
// 	bndl, err := c.client.SendTrytes(trytes, depth, mwm)
// 	// SendTransfer(seed Trytes, depth uint64, mwm uint64, transfers bundle.Transfers, options *SendTransfersOptions) (bundle.Bundle, error)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	bytes, err := json.Marshal(bndl)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	log.Println(string(bytes))
// 	hash := bundle.TailTransactionHash(bndl)
// 	log.Println("broadcasted bundle with tail tx hash: ", hash)
// 	if len(bndl) > 0 {
// 		ret = bndl[0].Bundle
// 	}
// 	log.Println(ret)
// 	bndl2, err := c.client.ReplayBundle(bndl[0].Hash, 3, 14)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	log.Println(bndl2[0].Bundle)

// 	return
// }
