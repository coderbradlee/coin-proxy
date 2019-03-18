package trx

import (
// "fmt"
)

type Vstruct struct {
	Amount        int64  `json:"amount"`
	Owner_address string `json:"owner_address"`
	To_address    string `json:"to_address"`
}
type P struct {
	Value Vstruct `json:"value"`
	// Value    []byte `json:"value"`
	Type_url string `json:"type_url"`
}

// type CT struct {
// 	Parameter P      `json:"parameter"`
// 	Type      string `json:"type"`
// }

type CT struct {
	Parameter    *P     `json:"parameter,omitempty"`
	Type         string `json:"type,omitempty"`
	Provider     []byte `json:"provider,omitempty"`
	ContractName []byte `json:"ContractName,omitempty"`
}
type Rd struct {
	Auths    []*Acuthrity `json:"auths,omitempty"`
	Data     []byte       `json:"data,omitempty"`
	Contract []CT         `json:"contract"`
	// Ref_block_bytes string `json:"ref_block_bytes"`
	Ref_block_bytes []byte `json:"ref_block_bytes"`
	RefBlockNum     int64  `json:"ref_block_num,omitempty"`
	// Ref_block_hash  string `json:"ref_block_hash"`
	Ref_block_hash []byte `json:"ref_block_hash"`
	Expiration     int64  `json:"expiration"`
	Scripts        []byte `json:"script,omitempty"`
	Timestamp      int64  `json:"timestamp"`
}
type Acuthrity struct {
	Account        *AccountId `json:"account,omitempty"`
	PermissionName []byte     `json:"permissionName,omitempty"`
}
type AccountId struct {
	Name    []byte `json:"name,omitempty"`
	Address []byte `json:"address,omitempty"`
}

type GetTransactionResp struct {
	Raw_data TransactionRaw `json:"raw_data,omitempty"`
	// only support size = 1,  repeated list here for muti-sig extension
	Signature []string      `json:"signature,omitempty"`
	Ret       []interface{} `json:"ret,omitempty"`
	TxID      string        `json:"txID,omitempty"`
}

type TransactionRaw struct {
	Ref_block_bytes []byte      `json:"ref_block_bytes,omitempty"`
	Ref_block_num   int64       `json:"ref_block_num,omitempty"`
	Ref_block_hash  []byte      `json:"ref_block_hash,omitempty"`
	Expiration      int64       `json:"expiration,omitempty"`
	Auths           interface{} `json:"auths,omitempty"`
	// data not used

	Fee_limit int64  `json:"fee_limit,omitempty"`
	Data      []byte `json:"data,omitempty"`
	// only support size = 1,  repeated list here for extension
	Contract []Transaction_Contract `json:"contract,omitempty"`
	// scripts not used
	Scripts   []byte `json:"scripts,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
type Transaction_Contract struct {
	Type         string          `json:"type,omitempty"`
	Parameter    ParameterStruct `json:"parameter,omitempty"`
	Provider     []byte          `json:"provider,omitempty"`
	Contractname []byte          `json:"contractname,omitempty"`
}
type ParameterStruct struct {
	Type_url string      `json:"type_url,omitempty"`
	Value    ValueStruct `json:"value,omitempty"`
}
type ValueStruct struct {
	Owner_address    string `json:"owner_address,omitempty"`
	To_address       string `json:"to_address,omitempty"`
	Amount           uint64 `json:"amount,omitempty"`
	Data             string `json:"data,omitempty"`
	Contract_address string `json:"contract_address,omitempty"`
}
type BroadcastTransactionReturn struct {
	Result        bool   `json:"result,omitempty"`
	Code          string `json:"code,omitempty"`
	Message       []byte `json:"message,omitempty"`
	TransactionID string `json:"transactionID,omitempty"`
	Error         string `json:"error,omitempty"`
}

// func decodeUTF8(b []byte) rune {
// 	switch b0 := b[0]; {
// 	case b0 < 0x80:
// 		return rune(b0)
// 	case b0 < 0xE0:
// 		return rune(b0&b2Mask)<<6 |
// 			rune(b[1]&mbMask)
// 	case b0 < 0xF0:
// 		return rune(b0&b3Mask)<<12 |
// 			rune(b[1]&mbMask)<<6 |
// 			rune(b[2]&mbMask)
// 	default:
// 		return rune(b0&b4Mask)<<18 |
// 			rune(b[1]&mbMask)<<12 |
// 			rune(b[2]&mbMask)<<6 |
// 			rune(b[3]&mbMask)
// 	}
// }
// func (u *BroadcastTransactionReturn) MarshalJSON() ([]byte, error) {
// 	me := []rune{u.Message}
// 	ret := fmt.Sprintf(`{"result":%v,"code":%v,"message":%v,"transactionID":%v,"error":%v,}`, u.Result, u.Code, string(me), u.TransactionID, u.Error)
// 	return []byte(ret), nil
// }

type Transaction_Result struct {
	Fee int64 `protobuf:"varint,1,opt,name=fee" json:"fee,omitempty"`
	Ret int32 `protobuf:"varint,2,opt,name=ret,enum=protocol.Transaction_ResultCode" json:"ret,omitempty"`
}

type Transaction struct {
	// {
	// 	"txID": "5b3f2fd7376f9051bb2482575d85f7a1f34684456f43045a8a1fb5ad35018ef0",
	// Signature [][]byte `json:"signature,omitempty"`
	TxID string `json:"txID,omitempty"`
	// Raw_data  Rd       `json:"raw_data"`
	Raw_data *TransactionRaw `protobuf:"bytes,1,opt,name=raw_data,json=rawData" json:"raw_data,omitempty"`
	// only support size = 1,  repeated list here for muti-sig extension
	Signature []string              `protobuf:"bytes,2,rep,name=signature,proto3" json:"signature,omitempty"`
	Ret       []*Transaction_Result `protobuf:"bytes,5,rep,name=ret" json:"ret,omitempty"`
	// 	"raw_data": {
	// 		"contract": [
	// 			{
	// 				"parameter": {
	// 					"value": {
	// 						"amount": 1000,
	// 						"owner_address": "41dbb157467c1b206494fc977c008c25f5901dfe4a",
	// 						"to_address": "41195b896cf7364d50aa274dc749b94d1a0b63f9f4"
	// 					},
	// 					"type_url": "type.googleapis.com/protocol.TransferContract"
	// 				},
	// 				"type": "TransferContract"
	// 			}
	// 		],
	// 		"ref_block_bytes": "03b8",
	// 		"ref_block_hash": "627c8f545cb1fb81",
	// 		"expiration": 1540630533000,
	// 		"timestamp": 1540630475290
	// 	}
	// }
	Error string `json:"error,omitempty"`
}

// func (u *Transaction) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(&Transaction{
// 		LastSeen: u.LastSeen.Unix(),
// 		MyUser:   u,
// 	})
// }
// func (u *Transaction) UnmarshalJSON(data []byte) error {
// 	type Alias MyUser
// 	aux := &Transaction{
// 		Alias: (*Alias)(u),
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	u.LastSeen = time.Unix(aux.LastSeen, 0)
// 	return nil
// }
