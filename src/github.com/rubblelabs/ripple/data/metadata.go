package data

import (
	"fmt"
	"sort"
)

type LedgerEntryState uint8

const (
	Created LedgerEntryState = iota
	Modified
	Deleted
)

type AffectedNode struct {
	FinalFields       LedgerEntry `json:",omitempty"`
	LedgerEntryType   LedgerEntryType
	LedgerIndex       *Hash256    `json:",omitempty"`
	PreviousFields    LedgerEntry `json:",omitempty"`
	NewFields         LedgerEntry `json:",omitempty"`
	PreviousTxnID     *Hash256    `json:",omitempty"`
	PreviousTxnLgrSeq *uint32     `json:",omitempty"`
}

type NodeEffect struct {
	ModifiedNode *AffectedNode `json:",omitempty"`
	CreatedNode  *AffectedNode `json:",omitempty"`
	DeletedNode  *AffectedNode `json:",omitempty"`
}

type NodeEffects []NodeEffect

type MetaData struct {
	AffectedNodes     NodeEffects
	TransactionIndex  uint32
	TransactionResult TransactionResult
	DeliveredAmount   *Amount `json:"delivered_amount,omitempty"`
}

type TransactionSlice []*TransactionWithMetaData

func (s TransactionSlice) Len() int      { return len(s) }
func (s TransactionSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TransactionSlice) Less(i, j int) bool {
	if s[i].LedgerSequence == s[j].LedgerSequence {
		return s[i].MetaData.TransactionIndex < s[j].MetaData.TransactionIndex
	}
	return s[i].LedgerSequence < s[j].LedgerSequence
}

func (s TransactionSlice) Sort() { sort.Sort(s) }

type TransactionWithMetaData struct {
	Transaction
	MetaData       MetaData   `json:"meta"`
	Date           RippleTime `json:"date"`
	LedgerSequence uint32     `json:"ledger_index"`
	Id             Hash256    `json:"-"`
}
type Transaction2 struct {
	Result AccountTxResultIn `json:"result,omitempty"`
}
type AccountTxResultIn struct {
	Account          string `json:"account,omitempty"`
	Ledger_index_max int64  `json:"ledger_index_max,omitempty"`
	Ledger_index_min int64  `json:"ledger_index_min,omitempty"`
	Limit            int    `json:"limit,omitempty"`
	// Validated        bool                      `json:"validated,omitempty"`
	Marker        map[string]interface{}     `json:"marker,omitempty"`
	Transactions  []TransactionWithMetaData2 `json:"transactions,omitempty"`
	Error         string                     `json:"error,omitempty"`
	Error_code    uint32                     `json:"error_code,omitempty"`
	Error_message string                     `json:"error_message,omitempty"`
	Status        string                     `json:"status,omitempty"`
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
type TransactionWithMetaData2 struct {
	Tx        AccountTxTransaction `json:"tx,omitempty"`
	MetaData  MetaData2            `json:"meta,omitempty"`
	Validated bool                 `json:"validated,omitempty"`
}
type MetaData2 struct {
	AffectedNodes     []NodeEffect2 `json:"AffectedNodes,omitempty"`
	TransactionIndex  uint32        `json:"TransactionIndex,omitempty"`
	TransactionResult string        `json:"TransactionResult,omitempty"`
	Delivered_amount  string        `json:"delivered_amount,omitempty"`
}
type NodeEffect2 struct {
	ModifiedNode *AffectedNode `json:",omitempty"`
	CreatedNode  *AffectedNode `json:",omitempty"`
	DeletedNode  *AffectedNode `json:",omitempty"`
}

func (t *TransactionWithMetaData) GetType() string    { return t.Transaction.GetType() }
func (t *TransactionWithMetaData) Prefix() HashPrefix { return HP_TRANSACTION_NODE }
func (t *TransactionWithMetaData) NodeType() NodeType { return NT_TRANSACTION_NODE }
func (t *TransactionWithMetaData) Ledger() uint32     { return t.LedgerSequence }
func (t *TransactionWithMetaData) NodeId() *Hash256   { return &t.Id }

func (t *TransactionWithMetaData) Affects(account Account) bool {
	for _, effect := range t.MetaData.AffectedNodes {
		if _, final, _, _ := effect.AffectedNode(); final.Affects(account) {
			return true
		}
	}
	return false
}

func NewTransactionWithMetadata(typ TransactionType) *TransactionWithMetaData {
	return &TransactionWithMetaData{Transaction: TxFactory[typ]()}
}

// AffectedNode returns the AffectedNode, the current LedgerEntry,
// the previous LedgerEntry (which might be nil) and the LedgerEntryState
func (effect *NodeEffect) AffectedNode() (*AffectedNode, LedgerEntry, LedgerEntry, LedgerEntryState) {
	var (
		node            *AffectedNode
		final, previous LedgerEntry
		state           LedgerEntryState
	)
	switch {
	case effect.CreatedNode != nil && effect.CreatedNode.NewFields != nil:
		node, final, state = effect.CreatedNode, effect.CreatedNode.NewFields, Created
	case effect.DeletedNode != nil && effect.DeletedNode.FinalFields != nil:
		node, final, state = effect.DeletedNode, effect.DeletedNode.FinalFields, Deleted
	case effect.ModifiedNode != nil && effect.ModifiedNode.FinalFields != nil:
		node, final, state = effect.ModifiedNode, effect.ModifiedNode.FinalFields, Modified
	case effect.ModifiedNode != nil && effect.ModifiedNode.FinalFields == nil:
		node, final, state = effect.ModifiedNode, LedgerEntryFactory[effect.ModifiedNode.LedgerEntryType](), Modified
	default:
		panic(fmt.Sprintf("Unknown LedgerEntryState: %+v", effect))
	}
	previous = node.PreviousFields
	if previous == nil {
		previous = LedgerEntryFactory[final.GetLedgerEntryType()]()
	}
	return node, final, previous, state
}
