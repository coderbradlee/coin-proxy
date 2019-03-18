package trx

import (
	// "bytes"
	// "encoding/hex"
	// "crypto/rand"
	// "encoding/binary"
	// "encoding/json"
	"errors"
	// "fmt"
	"context"
	"github.com/sasaxie/go-client-api/common/base58"
	// "io/ioutil"
	"log"
	// "net"
	// "net/http"
	// "strconv"
	// "crypto/ecdsa"
	// "crypto/sha256"
	// "github.com/golang/protobuf/proto"
	"github.com/sasaxie/go-client-api/common/crypto"
	"github.com/sasaxie/go-client-api/core"
	"github.com/sasaxie/go-client-api/util"
	// "github.com/ethereum/go-ethereum/crypto/sha3"
	// "github.com/sasaxie/go-client-api/common/global"
	"github.com/sasaxie/go-client-api/service"
	// "sync"
	// "time"
	// "unicode/utf8"
	"strconv"
)

func (r *RPCClient) Sendtest(from, to, amount, password, dir, privatenotindir string) (result BroadcastTransactionReturn, err error) {
	f := base58.DecodeCheck(from)
	if len(f) == 0 {
		err = errors.New("address format error")
		return
	}
	t := base58.DecodeCheck(to)
	if len(t) == 0 {
		err = errors.New("address format error")
		return
	}
	key, err := r.GetPrivate(dir, password, from)
	if err != nil {
		key = privatenotindir
		// return
		log.Println("GetPrivate:", err)
		return
	}
	a, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		log.Println("ParseInt:", err)
		return
	}
	log.Println("sign privatekey:", key)
	priKey, err := crypto.GetPrivateKeyByHexString(key)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("205 start:")

	transferContract := new(core.TransferContract)
	transferContract.OwnerAddress = f
	transferContract.ToAddress = t
	transferContract.Amount = a
	client := service.NewGrpcClient("localhost:50051")
	client.Start()
	defer client.Conn.Close()
	transferTransaction, err := client.Client.CreateTransaction(context.
		Background(), transferContract)

	if err != nil {
		log.Println("transfer error: ", err)
		return
	}
	log.Println(*transferTransaction)
	if transferTransaction == nil || len(transferTransaction.
		GetRawData().GetContract()) == 0 {
		log.Println("transfer error: invalid transaction")
	}

	util.SignTransaction(transferTransaction, priKey)
	log.Println(transferTransaction)
	return
}
