package ipfs

import (
	// "bytes"
	// "fmt"
	// "crypto/sha256"
	// "encoding/json"
	"errors"
	// "fmt"
	// "github.com/ethereum/go-ethereum/common"
	"log"
	// "math/big"
	// "net/http"
	// "strconv"
	// "strings"
	// "sync"
	// "time"
	//"github.com/ethereumproject/go-ethereum/common"
	// "91pool/util"
	"../aes"
	shell "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	// "regexp"
)

type RPCClient struct {
	Url    string
	PeerId string
	Client *shell.Shell
}

func NewRPCClient(url, peerId string) *RPCClient {
	rpcClient := &RPCClient{Url: url, PeerId: peerId}
	rpcClient.Client = shell.NewShell(url)
	return rpcClient
}

func (r *RPCClient) Get(mnemonicHash, dir string) (ret string, err error) {
	//首先调用ipfs name resolve QmZnrtGXvFzwTXjtykQRMhrxKApZcy8ghUaYi9gRnzUY1b
	// err = r.getFromChain(dir)
	// if err != nil {
	// 	return
	// }

	b, err := ioutil.ReadFile(dir + "/" + mnemonicHash)
	if err != nil {
		// fmt.Print(err)
		log.Println("read file:", err)
		err = errors.New("there's no contents under this key")
		return
	}
	// log.Println(string(b))
	origData, err := aes.AesDecrypt(b, []byte(r.PeerId))
	if err != nil {
		log.Println("AesDecrypt:", err)
		return
	}
	// log.Println("13")
	if len(origData) == 0 {
		err = errors.New("decrpt is null")
		return
	}
	// log.Println("14")
	ret = string(origData)
	return
}
func (r *RPCClient) getFromChain(dir string) (err error) {
	paths := "/ipns/" + r.PeerId
	dirHash, err := r.Client.ResolvePath(paths)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("dirHash get from chain:", dirHash)
	err = r.Client.Get(dirHash, dir)
	if err != nil {
		log.Println(err)

		return
	}
	return
}
func (r *RPCClient) Add(key, value, dir string) (ret bool, err error) {
	// err = r.getFromChain(dir)
	// if err != nil {
	// 	return
	// }
	err = r.writeFile(dir, r.PeerId, key, value)
	/////调用命令
	// ipfs add -r website/
	// QmRt7Vpe3r7z8TnNLApfVNdXzGKzomVr8NUisUya87nY95
	// ipfs name publish QmRt7Vpe3r7z8TnNLApfVNdXzGKzomVr8NUisUya87nY95
	go func() {
		final, err := r.Client.AddDir(dir)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("final:", final)
		///////////////////////////////////
		pubResp, err := r.Client.Publish(final, "use for what")
		if err != nil {
			log.Println("Publish:", err)
			return
		}
		log.Println(pubResp)
	}()

	ret = true
	return
}
func (r *RPCClient) writeFile(dir, password, filename, privateKey string) (err error) {
	key := []byte(password)
	result, err := aes.AesEncrypt([]byte(privateKey), key)
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
