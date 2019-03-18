package ripple

import (
	"github.com/rubblelabs/ripple/data"
)

type transaction struct {
	Account         string `json:",omitempty"`
	Amount          int    `json:",omitempty"`
	Destination     string `json:",omitempty"`
	TransactionType string `json:",omitempty"`
	SettleDelay     int    `json:",omitempty"`
	PublicKey       string `json:",omitempty"`
	DestinationTag  uint32 `json:",omitempty"`
}

type param struct {
	Offline     bool         `json:"offline,omitempty"`
	Secret      string       `json:"secret,omitempty"`
	TxJSON      *transaction `json:"tx_json,omitempty"`
	TxBlob      string       `json:"tx_blob,omitempty"`
	FeeMultMax  int          `json:"fee_mult_max,omitempty"`
	Account     string       `json:"account,omitempty"`
	DestAccount string       `json:"destination_account,omitempty"`
	ChannelID   string       `json:"channel_id,omitempty"`
	Amount      int          `json:"amount,omitempty"`

	///for account_tx
	Binary           bool  `json:"binary,omitempty"`
	Forward          bool  `json:"forward,omitempty"`
	Ledger_index_max int64 `json:"ledger_index_max,omitempty"`

	Transaction string `json:"transaction,omitempty"`
	Limit       int    `json:"limit,omitempty"`

	Marker           *Mark `json:"marker,omitempty"`
	Ledger_index_min int64 `json:"ledger_index_min,omitempty"`
	Offset           int64 `json:"offset,omitempty"`
}

// type accoutTxParam struct {
// 	// "account": "r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59",
// 	// "binary": false,
// 	// "forward": false,
// 	// "ledger_index_max": -1,
// 	// "ledger_index_min": -1,
// 	// "limit": 2
// 	Account          string `json:"account,omitempty"`
// 	Binary           bool   `json:"binary,omitempty"`
// 	Forward          bool   `json:"forward,omitempty"`
// 	Ledger_index_max int64  `json:"ledger_index_max,omitempty"`
// 	Ledger_index_min int64  `json:"ledger_index_min,omitempty"`

// 	Limit int `json:"limit,omitempty"`
// 	// Result *AccountTxResult `json:"result,omitempty"`
// }
type request struct {
	Method string  `json:"method"`
	Params []param `json:"params,omitempty"`
}

//Channel when querying for channels between accounts, this is the information provided by the server
type Channel struct {
	Account            string `json:"account"`
	Amount             string `json:"amount"`
	Balance            string `json:"balance"`
	ChannelID          string `json:"channel_id"`
	DestinationAccount string `json:"destination_account"`
	PublicKey          string `json:"public_key"`
	PublicKeyHex       string `json:"public_key_hex"`
	SettleDelay        int    `json:"settle_delay"`
}

// Response a response from the rippled server
type Response struct {
	Result *struct {
		Role             string `json:"role"`
		Status           string `json:"status"`
		EngineResult     string `json:"engine_result"`
		EngineResultCode *int   `json:"engine_result_code"` // pointer to int allows nil (instead of 0 as default value)
		EngineResultMsg  string `json:"engine_result_message"`
		TxBlob           string `json:"tx_blob"`
		ErrorMsg         string `json:"error_message"`
		Account          string `json:"account"`
		Signature        string `json:"signature"`

		Tx *struct {
			Account         string
			Amount          string
			Destination     string
			Fee             string
			Flags           int
			Sequence        int
			SigningPubKey   string
			TransactionType string
			TxnSignature    string
			Hash            string `json:"hash"`
		} `json:"tx_json"`
		Channels *[]Channel `json:"channels"`
	}
}
type AccountInfoResult struct {
	Result AccountInfoResultIn `json:"result,omitempty"`
}
type AccountInfoResultIn struct {
	LedgerSequence uint32           `json:"ledger_current_index,omitempty"`
	AccountData    data.AccountRoot `json:"account_data,omitempty"`
	Status         string           `json:"status,omitempty"`
	Validated      bool             `json:"validated,omitempty"`
	// "error": "actNotFound",
	//     "error_code": 19,
	// 	"error_message": "Account not found.",

	Error         string `json:"error,omitempty"`
	Error_code    uint32 `json:"error_code,omitempty"`
	Error_message string `json:"error_message,omitempty"`
}
type SignedResult struct {
	Result SignedResultIn `json:"result,omitempty"`
}
type TxStruct struct {
	Account     string     `json:"account,omitempty"`
	Amount      data.Value `json:"amount,omitempty"`
	Destination string     `json:"destination,omitempty"`
	Fee         string     `json:"fee,omitempty"`
	Flags       uint32     `json:"flags,omitempty"`
	Sequence    uint32     `json:"sequence,omitempty"`

	SigningPubKey   string `json:"signingPubKey,omitempty"`
	TransactionType string `json:"transactionType,omitempty"`
	TxnSignature    string `json:"txnSignature,omitempty"`
	Hash            string `json:"hash,omitempty"`
	// "Account": "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn",
	//             "Amount": {
	//                 "currency": "USD",
	//                 "issuer": "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn",
	//                 "value": "1"
	//             },
	//             "Destination": "ra5nK24KXen9AHvsdFTKHSANinZseWnPcX",
	//             "Fee": "10000",
	//             "Flags": 2147483648,
	//             "Sequence": 360,
	//             "SigningPubKey": "03AB40A0490F9B7ED8DF29D246BF2D6269820A0EE7742ACDD457BEA7C7D0931EDB",
	//             "TransactionType": "Payment",
	//             "TxnSignature": "304402200E5C2DD81FDF0BE9AB2A8D797885ED49E804DBF28E806604D878756410CA98B102203349581946B0DDA06B36B35DBC20EDA27552C1F167BCF5C6ECFF49C6A46F8580",
	//             "hash": "4D5D90890F8D49519E4151938601EF3D0B30B16CD6A519D9C99102C9FA77F7E0"
}
type SignedResultIn struct {
	Tx_json TxStruct `json:"tx_json,omitempty"`
	Status  string   `json:"status,omitempty"`
	Tx_blob string   `json:"tx_blob,omitempty"`

	Error         string `json:"error,omitempty"`
	Error_code    uint32 `json:"error_code,omitempty"`
	Error_message string `json:"error_message,omitempty"`
}
type SubmitResult struct {
	Result SubmitResultIn `json:"result,omitempty"`
}
type SubmitResultIn struct {
	Tx_json TxStruct `json:"tx_json,omitempty"`
	Status  string   `json:"status,omitempty"`
	Tx_blob string   `json:"tx_blob,omitempty"`

	Engine_result         string `json:"engine_result,omitempty"`
	Engine_result_code    uint32 `json:"engine_result_code,omitempty"`
	Engine_result_message string `json:"engine_result_message,omitempty"`

	Error         string `json:"error,omitempty"`
	Error_code    uint32 `json:"error_code,omitempty"`
	Error_message string `json:"error_message,omitempty"`
}
type AccountTxResult struct {
	// Id     uint32            `json:"id,omitempty"`
	// Status string            `json:"status,omitempty"`
	// Type   string            `json:"type,omitempty"`
	Result AccountTxResultIn `json:"result,omitempty"`
}
type Mark struct {
	Ledger int64 `json:"ledger"`
	Seq    int64 `json:"seq"`
}
type AccountTxResultIn struct {
	Account          string `json:"account,omitempty"`
	Ledger_index_max int64  `json:"ledger_index_max,omitempty"`
	Ledger_index_min int64  `json:"ledger_index_min,omitempty"`
	Limit            int    `json:"limit,omitempty"`
	// Validated        bool                      `json:"validated,omitempty"`
	Marker        Mark                      `json:"marker,omitempty"`
	Transactions  []TransactionWithMetaData `json:"transactions,omitempty"`
	Error         string                    `json:"error,omitempty"`
	Error_code    uint32                    `json:"error_code,omitempty"`
	Error_message string                    `json:"error_message,omitempty"`
	Status        string                    `json:"status,omitempty"`
}
type AccountTxTransaction struct {
	Account        string `json:"account,omitempty"`
	Amount         string `json:"amount,omitempty"`
	Destination    string `json:"destination,omitempty"`
	DestinationTag uint32 `json:"DestinationTag,omitempty"`
	Fee            string `json:"fee,omitempty"`
	Flags          uint32 `json:"flags,omitempty"`
	Sequence       uint32 `json:"sequence,omitempty"`

	SigningPubKey   string `json:"signingPubKey,omitempty"`
	TransactionType string `json:"transactionType,omitempty"`
	TxnSignature    string `json:"txnSignature,omitempty"`
	Date            uint32 `json:"date,omitempty"`
	Hash            string `json:"hash,omitempty"`
	InLedger        uint32 `json:"inLedger,omitempty"`
	Ledger_index    uint32 `json:"ledger_index,omitempty"`
}
type TransactionWithMetaData struct {
	Tx        AccountTxTransaction `json:"tx,omitempty"`
	MetaData  MetaData             `json:"meta,omitempty"`
	Validated bool                 `json:"validated,omitempty"`
}
type MetaData struct {
	AffectedNodes     []NodeEffect `json:"AffectedNodes,omitempty"`
	TransactionIndex  uint32       `json:"TransactionIndex,omitempty"`
	TransactionResult string       `json:"TransactionResult,omitempty"`
	Delivered_amount  string       `json:"delivered_amount,omitempty"`
}
type NodeEffect struct {
	ModifiedNode *data.AffectedNode `json:",omitempty"`
	CreatedNode  *data.AffectedNode `json:",omitempty"`
	DeletedNode  *data.AffectedNode `json:",omitempty"`
}

// {
//     "result": {
//         "engine_result": "tesSUCCESS",
//         "engine_result_code": 0,
//         "engine_result_message": "The transaction was applied. Only final in a validated ledger.",
//         "status": "success",
//         "tx_blob": "1200002280000000240000016961D4838D7EA4C6800000000000000000000000000055534400000000004B4E9C06F24296074F7BC48F92A97916C6DC5EA9684000000000002710732103AB40A0490F9B7ED8DF29D246BF2D6269820A0EE7742ACDD457BEA7C7D0931EDB74473045022100A7CCD11455E47547FF617D5BFC15D120D9053DFD0536B044F10CA3631CD609E502203B61DEE4AC027C5743A1B56AF568D1E2B8E79BB9E9E14744AC87F38375C3C2F181144B4E9C06F24296074F7BC48F92A97916C6DC5EA983143E9D4A2B8AA0780F682D136F7A56D6724EF53754",
//         "tx_json": {
//             "Account": "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn",
//             "Amount": {
//                 "currency": "USD",
//                 "issuer": "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn",
//                 "value": "1"
//             },
//             "Destination": "ra5nK24KXen9AHvsdFTKHSANinZseWnPcX",
//             "Fee": "10000",
//             "Flags": 2147483648,
//             "Sequence": 361,
//             "SigningPubKey": "03AB40A0490F9B7ED8DF29D246BF2D6269820A0EE7742ACDD457BEA7C7D0931EDB",
//             "TransactionType": "Payment",
//             "TxnSignature": "3045022100A7CCD11455E47547FF617D5BFC15D120D9053DFD0536B044F10CA3631CD609E502203B61DEE4AC027C5743A1B56AF568D1E2B8E79BB9E9E14744AC87F38375C3C2F1",
//             "hash": "5B31A7518DC304D5327B4887CD1F7DC2C38D5F684170097020C7C9758B973847"
//         }
//     }
// }
type TransactionDetailIn struct {
	Account     string `json:"Account,omitempty"`
	Amount      string `json:"Amount,omitempty"`
	Destination string `json:"Destination,omitempty"`
	Fee         string `json:"Fee,omitempty"`
	Flags       uint32 `json:"Flags,omitempty"`
	Sequence    uint32 `json:"Sequence,omitempty"`

	SigningPubKey   string   `json:"SigningPubKey,omitempty"`
	TransactionType string   `json:"TransactionType,omitempty"`
	TxnSignature    string   `json:"TxnSignature,omitempty"`
	Hash            string   `json:"hash,omitempty"`
	Status          string   `json:"status,omitempty"`
	Validated       bool     `json:"validated,omitempty"`
	MetaData        MetaData `json:"meta,omitempty"`
	Date            uint32   `json:"date,omitempty"`
	InLedger        uint32   `json:"inLedger,omitempty"`
	Ledger_index    uint32   `json:"ledger_index,omitempty"`
	//     "date": 594701280,
	//     "inLedger": 62,
	//     "ledger_index": 62,
	//     "meta": {
	//         "AffectedNodes": [
	//             {
	//                 "ModifiedNode": {
	//                     "FinalFields": {
	//                         "Account": "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
	//                         "Balance": "99999999798999920",
	//                         "Flags": 0,
	//                         "OwnerCount": 0,
	//                         "Sequence": 9
	//                     },
	//                     "LedgerEntryType": "AccountRoot",
	//                     "LedgerIndex": "2B6AC232AA4C4BE41BF49D2459FA4A0347E1B543A4C92FCEE0821C0201E2E9A8",
	//                     "PreviousFields": {
	//                         "Balance": "99999999798999930",
	//                         "Sequence": 8
	//                     },
	//                     "PreviousTxnID": "23AD613D0C3ED00BFB175C2A65EBA9457AC77E9014A8F587A250D01293F32DED",
	//                     "PreviousTxnLgrSeq": 62
	//                 }
	//             }
	//         ],
	//         "TransactionIndex": 1,
	//         "TransactionResult": "tecNO_DST_INSUF_XRP"
	//     },
	//     "status": "success",
	//     "validated": true
	// }

	Error         string `json:"error,omitempty"`
	Error_code    uint32 `json:"error_code,omitempty"`
	Error_message string `json:"error_message,omitempty"`
}
type TransactionDetail struct {
	Result TransactionDetailIn `json:"result,omitempty"`
}
type VerifyResultIn struct {
	From              string `json:"from,omitempty"`
	To                string `json:"to,omitempty"`
	TransactionId     string `json:"transactionId,omitempty"`
	Memo              string `json:"memo,omitempty"`
	Seq               string `json:"seq,omitempty"`
	Quantity          string `json:"quantity,omitempty"`
	Packed            bool   `json:"packed,omitempty"`
	TransactionResult string `json:"transactionResult,omitempty"`
}
type VerifyResult struct {
	Trans  []VerifyResultIn `json:"detail,omitempty"`
	Marker Mark             `json:"mark,omitempty"`
}
