package ada

import (
	"encoding/json"
	"time"
	// "io"
	// "log"
	// "net"
	// "net/http"

	// "strings"
	// "sync"
	// "sync/atomic"
	// "time"
	"encoding/hex"
	"fmt"
	"math/big"
)

type JsonRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	// Params []json.RawMessage `json:"params"`
	Params []interface{} `json:"params"`
}
type JSONTime struct {
	time.Time
}

const JSONTimeFormat = "2006-01-02T15:04:05"

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", t.Format(JSONTimeFormat))), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}

	t.Time, err = time.Parse(`"`+JSONTimeFormat+`"`, string(data))
	return err
}

// ParseJSONTime will parse a string into a JSONTime object
func ParseJSONTime(date string) (JSONTime, error) {
	var t JSONTime
	var err error
	t.Time, err = time.Parse(JSONTimeFormat, string(date))
	return t, err
}

type SHA256Bytes []byte

func (t SHA256Bytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(t))
}

func (t *SHA256Bytes) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*t, err = hex.DecodeString(s)
	return
}
func (t SHA256Bytes) String() string {
	return hex.EncodeToString(t)
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}

func String2Big(num string) *big.Int {
	n := new(big.Int)
	n.SetString(num, 0)
	return n
}

type Diag struct {
	Msg string `json:"msg,omitempty"`
}
type ResponseStruct struct {
	Status     string      `json:"status,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Diagnostic Diag        `json:"diagnostic,omitempty"`
	// "diagnostic":{"msg":"Request error (No account with id 2cWKMJemoBak8MPHidBtVJG63rD76dD21vTk4dLzxoajfZSEgrLRuPTXuwe8M8vvdJCnG@0 found)"},"message":"UnknownError"
	Message string `json:"message,omitempty"`
}
type Inputstruct struct {
	Address string `json:"address,omitempty"`

	// 			"address": "DdzFFzCqrhtBatWqyFge4w6M6VLgNUwRHiXTAg3xfQCUdTcjJxSrPHVZJBsQprUEc5pRhgMWQaGciTssoZVwrSKmG1fneZ1AeCtLgs5Y",
	// 			"amount": 51541623
	Amount uint64 `json:"amount,omitempty"`
}
type Outputstruct struct {
	Address string `json:"address,omitempty"`
	// "address": "DdzFFzCqrhtBatWqyFge4w6M6VLgNUwRHiXTAg3xfQCUdTcjJxSrPHVZJBsQprUEc5pRhgMWQaGciTssoZVwrSKmG1fneZ1AeCtLgs5Y",
	// // 			"amount": 49369962
	Amount uint64 `json:"amount,omitempty"`
}
type Statustruct struct {
	// "tag": "persisted",
	Tag string `json:"tag,omitempty"`
	// 		"data": {}
	Data interface{} `json:"data,omitempty"`
}
type TransactionData struct {
	Id            string `json:"id,omitempty"`
	Confirmations uint64 `json:"confirmations,omitempty"`

	amount       uint64         `json:"amount,omitempty"`
	Inputs       []Inputstruct  `json:"inputs,omitempty"`
	Outputs      []Outputstruct `json:"outputs,omitempty"`
	Type         string         `json:"type,omitempty"`
	Direction    string         `json:"direction,omitempty"`
	CreationTime JSONTime       `json:"creationTime,omitempty"`
	Status       Statustruct    `json:"status,omitempty"`
}
type Transaction struct {
	Status     string            `json:"status,omitempty"`
	Data       []TransactionData `json:"data,omitempty"`
	Meta       interface{}       `json:"meta,omitempty"`
	Diagnostic Diag              `json:"diagnostic,omitempty"`
	// "diagnostic":{"msg":"Request error (No account with id 2cWKMJemoBak8MPHidBtVJG63rD76dD21vTk4dLzxoajfZSEgrLRuPTXuwe8M8vvdJCnG@0 found)"},"message":"UnknownError"
	Message string `json:"message,omitempty"`
}
type GetTransactionResponse struct {
	Hash         string   `json:"hash,omitempty"`
	Address      string   `json:"address,omitempty"`
	Amount       string   `json:"amount,omitempty"`
	CreationTime JSONTime `json:"creationTime,omitempty"`
}
type BalanceResponse struct {
	Status     string      `json:"status,omitempty"`
	Data       BalanceData `json:"data,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Diagnostic Diag        `json:"diagnostic,omitempty"`
	Message    string      `json:"message,omitempty"`
}
type BalanceData struct {
	Id                         string      `json:"id,omitempty"`
	Name                       string      `json:"name,omitempty"`
	Balance                    uint64      `json:"balance,omitempty"`
	HasSpendingPassword        bool        `json:"hasSpendingPassword,omitempty"`
	SpendingPasswordLastUpdate JSONTime    `json:"spendingPasswordLastUpdate,omitempty"`
	CreatedAt                  JSONTime    `json:"createdAt,omitempty"`
	AssuranceLevel             string      `json:"assuranceLevel,omitempty"`
	SyncState                  interface{} `json:"syncState,omitempty"`
}
type GetTransactionRet struct {
	Id  string `json:"id,omitempty"`
	Fee uint64 `json:"fee,omitempty"`
}
