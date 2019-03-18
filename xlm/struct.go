package xlm

import (
	"time"
)

type OffersPage struct {
	Links    Links `json:"_links,omitempty"`
	Embedded struct {
		Records []Offer `json:"records,omitempty"`
	} `json:"_embedded,omitempty"`
}
type Links struct {
	Self Link `json:"self,omitempty"`
	Next Link `json:"next,omitempty"`
	Prev Link `json:"prev,omitempty"`
}
type Link struct {
	Href      string `json:"href,omitempty"`
	Templated bool   `json:"templated,omitempty"`
}
type Offer struct {
	Links struct {
		Self Link `json:"self,omitempty"`
		// OfferMaker  Link `json:"offer_maker,omitempty"`
		Transaction Link `json:"transaction,omitempty"`
		Effects     Link `json:"effects,omitempty"`
		Succeeds    Link `json:"succeeds,omitempty"`
		Precedes    Link `json:"precedes,omitempty"`
	} `json:"_links,omitempty"`

	ID string `json:"id,omitempty"`
	PT string `json:"paging_token,omitempty"`

	Source_account string `json:"source_account,omitempty"`
	Type           string `json:"type,omitempty"`
	Type_i         int32  `json:"type_i,omitempty"`

	Created_at       *time.Time `json:"created_at,omitempty"`
	Transaction_hash string     `json:"transaction_hash,omitempty"`
	Starting_balance string     `json:"starting_balance,omitempty"`
	Funder           string     `json:"funder,omitempty"`
	Account          string     `json:"account,omitempty"`
	Asset_type       string     `json:"asset_type,omitempty"`
	From             string     `json:"from,omitempty"`
	To               string     `json:"to,omitempty"`

	Seller string `json:"seller,omitempty"`
	// Selling            Asset      `json:"selling,omitempty"`
	// Buying             Asset      `json:"buying,omitempty"`
	Amount string `json:"amount,omitempty"`
	// PriceR             Price      `json:"price_r,omitempty"`
	Price              string     `json:"price,omitempty"`
	LastModifiedLedger int32      `json:"last_modified_ledger,omitempty"`
	LastModifiedTime   *time.Time `json:"last_modified_time,omitempty"`
	// Memo               struct {
	// 	Type  string `json:"memo_type"`
	// 	Value string `json:"memo"`
	// }
}
type Asset struct {
	Type   string `json:"asset_type,omitempty"`
	Code   string `json:"asset_code,omitempty"`
	Issuer string `json:"asset_issuer,omitempty"`
}

type Price struct {
	N int32 `json:"n,omitempty"`
	D int32 `json:"d,omitempty"`
}
type TransactionResponse struct {
	Transaction_hash string     `json:"transaction_hash,omitempty"`
	From             string     `json:"from,omitempty"`
	To               string     `json:"to,omitempty"`
	Amount           string     `json:"amount,omitempty"`
	Created_at       *time.Time `json:"created_at,omitempty"`
	// Memo             struct {
	// 	Type  string `json:"memo_type"`
	// 	Value string `json:"memo"`
	// }
	Memo   string `json:"memo,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}
type Memo struct {
	Type  string `json:"memo_type"`
	Value string `json:"memo"`
}
