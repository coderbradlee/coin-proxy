package yoyow

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

// {
//     "id": 1,
//     "jsonrpc": "2.0",
//     "result": {
//         "head_block_num": 10621811,
//         "head_block_id": "00a21373202b690647ee56b32fa1033eeefa1950",
//         "head_block_time": "2018-09-22T00:00:48",
//         "head_block_age": "23 days old",
//         "last_irreversible_block_num": 10621801,
//         "chain_id": "3505e367fe6cde243f2a1c39bd8e58557e23271dd6cbf4b29a8dc8c44c9af8fe",
//         "participation": "100.00000000000000000",
//         "active_witnesses": [
//             [
//                 25997,
//                 "scheduled_by_vote_top"
//             ],

//         ],
//         "active_committee_members": [
//             25997,
//             26264,
//             26460,
//             26861,
//             27027
//         ]
//     }
// }
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

// type Witness struct {
// 	Account_uid uint32 `json:"account_uid"`
// }

// enum scheduled_witness_type
//    {
//       scheduled_by_vote_top  = 0,
//       scheduled_by_vote_rest = 1,
//       scheduled_by_pledge    = 2
//    };
type Witness struct {
	Head_block_num  uint32   `json:"head_block_num"`
	Head_block_id   string   `json:"head_block_id"`
	Head_block_time JSONTime `json:"head_block_time"`
	Head_block_age  string   `json:"head_block_age"`

	Last_irreversible_block_num uint32            `json:"last_irreversible_block_num"`
	Chain_id                    SHA256Bytes       `json:"chain_id"`
	Participation               string            `json:"participation"`
	Active_witnesses            map[uint64]string `json:"active_witnesses"`
	active_committee_members    []uint64          `json:"active_committee_members"`
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

type KeyAuthStruct struct {
	Weight_threshold  uint32        `json:"weight_threshold"`
	Account_uid_auths []interface{} `json:"account_uid_auths"`
	Key_auths         []interface{} `json:"key_auths"`
}
type allStruct struct {
	amount   uint32 `json:"amount"`
	asset_id uint32 `json:"asset_id"`
}
type RegStruct struct {
	Registrar             uint32    `json:"registrar"`
	Referrer              uint32    `json:"referrer"`
	Registrar_percent     uint32    `json:"registrar_percent"`
	Referrer_percent      uint32    `json:"referrer_percent"`
	Allowance_per_article allStruct `json:"allowance_per_article"`
	Max_share_per_article allStruct `json:"max_share_per_article"`
	Max_share_total       allStruct `json:"max_share_total"`
	Buyout_percent        uint32    `json:"buyout_percent"`
}
type Account struct {
	Id               string        `json:"id"`
	Uid              uint32        `json:"uid"`
	Name             string        `json:"name"`
	Owner            KeyAuthStruct `json:"owner"`
	Active           KeyAuthStruct `json:"active"`
	Secondary        KeyAuthStruct `json:"secondary"`
	Memo_key         string        `json:"memo_key"`
	Reg_info         RegStruct     `json:"reg_info"`
	Can_post         bool          `json:"can_post"`
	Can_reply        bool          `json:"can_reply"`
	Can_rate         bool          `json:"can_rate"`
	Is_full_member   bool          `json:"is_full_member"`
	Is_registrar     bool          `json:"is_registrar"`
	Is_admin         bool          `json:"is_admin"`
	Create_time      JSONTime      `json:"create_time"`
	Last_update_time JSONTime      `json:"last_update_time"`
	Active_data      string        `json:"active_data"`
	Secondary_data   string        `json:"secondary_data"`
	Statistics       string        `json:"statistics"`
}

// {"amount":1200000000,"asset_id":0}
type AccountBalance struct {
	Amount   uint64 `json:"amount"`
	Asset_id uint32 `json:"asset_id"`
}
type StatisticStruct struct {
	Id                              string   `json:"id"`
	Owner                           uint64   `json:"owner"`
	Total_ops                       uint64   `json:"total_ops"`
	Removed_ops                     uint64   `json:"removed_ops"`
	Prepaid                         uint64   `json:"prepaid"`
	Csaf                            uint64   `json:"csaf"`
	Core_balance                    uint64   `json:"core_balance"`
	Core_leased_in                  uint64   `json:"core_leased_in"`
	Core_leased_out                 uint64   `json:"core_leased_out"`
	Average_coins                   uint64   `json:"average_coins"`
	Average_coins_last_update       JSONTime `json:"average_coins_last_update"`
	Coin_seconds_earned             string   `json:"coin_seconds_earned"`
	Coin_seconds_earned_last_update JSONTime `json:"coin_seconds_earned_last_update"`

	Total_witness_pledge                uint64 `json:"total_witness_pledge"`
	Releasing_witness_pledge            uint64 `json:"releasing_witness_pledge"`
	Witness_pledge_release_block_number uint64 `json:"witness_pledge_release_block_number"`
	Last_witness_sequence               uint64 `json:"last_witness_sequence"`
	Uncollected_witness_pay             uint64 `json:"uncollected_witness_pay"`
	Witness_last_confirmed_block_num    uint64 `json:"witness_last_confirmed_block_num"`
	Witness_last_aslot                  uint64 `json:"witness_last_aslot"`

	Witness_total_produced                       uint64 `json:"witness_total_produced"`
	Witness_total_missed                         uint64 `json:"witness_total_missed"`
	Witness_last_reported_block_num              uint64 `json:"witness_last_reported_block_num"`
	Witness_total_reported                       uint64 `json:"witness_total_reported"`
	Total_committee_member_pledge                uint64 `json:"total_committee_member_pledge"`
	Releasing_committee_member_pledge            uint64 `json:"releasing_committee_member_pledge"`
	Committee_member_pledge_release_block_number uint64 `json:"committee_member_pledge_release_block_number"`

	Last_committee_member_sequence uint64 `json:"last_committee_member_sequence"`
	Can_vote                       bool   `json:"can_vote"`
	Is_voter                       bool   `json:"is_voter"`

	Last_voter_sequence                  uint64 `json:"last_voter_sequence"`
	Last_platform_sequence               uint64 `json:"last_platform_sequence"`
	Total_platform_pledge                uint64 `json:"total_platform_pledge"`
	Releasing_platform_pledge            uint64 `json:"releasing_platform_pledge"`
	Platform_pledge_release_block_number uint64 `json:"platform_pledge_release_block_number"`
	Last_post_sequence                   uint64 `json:"last_post_sequence"`
}
type FullAccount struct {
	Acc        Account         `json:"account"`
	Statistics StatisticStruct `json:"statistics"`
	Balances   []interface{}   `json:"balances"`
}
type CollectPoints struct {
	Ref_block_num    uint64        `json:"ref_block_num"`
	Ref_block_prefix uint64        `json:"ref_block_prefix"`
	Expiration       JSONTime      `json:"expiration"`
	Operations       []interface{} `json:"operations"`
	// map[ref_block_num:6949 ref_block_prefix:2.18008016e+08 expiration:2018-10-17T02:24:36 operations:[[6 map[fee:map[total:map[amount:100000 asset_id:0] options:map[from_csaf:map[amount:100000 asset_id:0]]] from:2.44958118e+08 to:2.44958118e+08 amount:map[amount:165000 asset_id:0] time:2018-10-17T02:24:00]]] signatures:[20090706bbc5b7ead8c6e2f527b6ef80c04e50b8883ef10d661a499fdcaa8b685229541afb274cdd19553d7cf7aaf6c52a012aca30f81ad24e23de06bebe034e0b 1f54e30b9d98e4f309d20340e159cb30d770120027f332d63f520a405a6402bca138ef3a4979ef88864a51656e0060ef46c9e29bc3ae417b028cba1e9eb0fcf5be]]

}
type TransferResponse struct {
	Ref_block_num    uint64        `json:"ref_block_num"`
	Ref_block_prefix uint64        `json:"ref_block_prefix"`
	Expiration       JSONTime      `json:"expiration"`
	Operations       []interface{} `json:"operations"`
	Signatures       []interface{} `json:"signatures"`
}
type AmountS struct {
	Amount   uint64 `json:"amount"`
	Asset_id uint64 `json:"asset_id"`
}
type OpInOp1 struct {
	Fee        interface{} `json:"fee"`
	From       uint64      `json:"from"`
	To         uint64      `json:"to"`
	Amount     AmountS     `json:"amount"`
	Memo       interface{} `json:"memo"`
	Extensions interface{} `json:"extensions"`
}
type HistoryOp struct {
	Id              string        `json:"id"`
	Result          []interface{} `json:"result"`
	Op              []interface{} `json:"op"`
	Block_timestamp JSONTime      `json:"block_timestamp"`
	Block_num       uint64        `json:"block_num"`
	Trx_in_block    uint64        `json:"trx_in_block"`
	Op_in_trx       uint64        `json:"op_in_trx"`
	Virtual_op      uint64        `json:"virtual_op"`
}
type AccountHistory struct {
	Memo        string    `json:"memo"`
	Description string    `json:"description"`
	Sequence    uint64    `json:"sequence"`
	Op          HistoryOp `json:"op"`
}

// {
// 	"id": 1,
// 	"jsonrpc": "2.0",
// 	"result": {
// 	  "ref_block_num": 19646,
// 	  "ref_block_prefix": 555752677,
// 	  "expiration": "2018-04-16T02:39:09",
// 	  "operations": [
// 		[
// 		  6,
// 		  {
// 			"fee": {
// 			  "total": {
// 				"amount": 100000,
// 				"asset_id": 0
// 			  },
// 			  "options": {
// 				"from_csaf": {
// 				  "amount": 100000,
// 				  "asset_id": 0
// 				}
// 			  }
// 			},
// 			"from": 250926091,
// 			"to": 250926091,
// 			"amount": {
// 			  "amount": 100000,
// 			  "asset_id": 0
// 			},
// 			"time": "2018-04-16T02:37:00"
// 		  }
// 		]
// 	  ],
// 	  "signatures": [
// 		"203a417b25f10110d8143d7476976abbcbb3490f13432630366e5b0d1d8d7580573c8595e93109af4a55282756b8b4916ae055147cceae1bc7b85f2b0a7f2fa042",
// 		"2054d3b25618ddaeae499297a483d5490bac77f35bac7dd850645400d7f8001a2265cd997ff62db54740e9fcda52b0bbbaf5aa6d12d3fbcd65a71e2ccf6baa1e1a"
// 	  ]
// 	}
//   }
