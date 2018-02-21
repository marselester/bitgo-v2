package bitgo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// walletService communicates with the wallet API endpoints.
type walletService struct {
	client *Client
}

// TxInfo is a response we get from consolidateunspents API endpoint.
type TxInfo struct {
	// TxID is an id of the transaction.
	TxID string `json:"txid"`
	// Tx is the serialized transaction.
	Tx string `json:"tx"`
	// Status if the transaction was signed.
	Status string `json:"status"`
}

// WalletConsolidateParams represents API parameters used when coalescing UTXOs.
// For more details, see https://www.bitgo.com/api/v2/#consolidate-wallet-unspents.
type WalletConsolidateParams struct {
	// Passphrase to decrypt the wallet's private key.
	WalletPassphrase string `json:"walletPassphrase,omitempty"`
	// Number of outputs created by the consolidation transaction (defaults to 1).
	NumUnspentsToMake int `json:"numUnspentsToMake,omitempty"`
	// Number of unspents to select (defaults to 25, max is 200).
	Limit int `json:"limit,omitempty"`
	// Ignore unspents smaller than this amount of satoshis.
	MinValue int64 `json:"minValue,omitempty"`
	// Ignore unspents larger than this amount of satoshis.
	MaxValue int64 `json:"maxValue,omitempty"`
	// The minimum height of unspents on the block chain to use.
	MinHeight int `json:"minHeight,omitempty"`
	// The desired fee rate for the transaction in satoshis/KB.
	FeeRate int `json:"feeRate,omitempty"`
	// Fee rate is automatically chosen by targeting a transaction confirmation
	// in this number of blocks (only available on BTC, FeeRate takes precedence if also set).
	FeeTxConfirmTarget int `json:"feeTxConfirmTarget,omitempty"`
	// Maximum percentage of an unspent's value to be used for fees. Cannot be combined with MinValue.
	MaxFeePercentage int `json:"maxFeePercentage,omitempty"`
	// The required number of confirmations for each transaction input.
	MinConfirms int `json:"minConfirms,omitempty"`
	// Apply the required confirmations set in MinConfirms for change outputs.
	EnforceMinConfirmsForChange bool `json:"enforceMinConfirmsForChange,omitempty"`
}

// Consolidate coalesces UTXOs currently held in a wallet to a smaller number.
func (s *walletService) Consolidate(ctx context.Context, walletID string, bodyParams *WalletConsolidateParams) (*TxInfo, error) {
	path := fmt.Sprintf("wallet/%s/consolidateunspents", walletID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil, bodyParams)
	if err != nil {
		return nil, err
	}

	tx := TxInfo{}
	_, err = s.client.Do(req, &tx)
	return &tx, err
}

// Unspent is an unspent transaction output (UTXO).
type Unspent struct {
	// The outpoint of the unspent (txid:vout). For example, "952ac7fd9c1a5df8380e0e305fac8b42db:0".
	ID string
	// The address that owns this unspent.
	Address string
	// Value of the unspent in satoshis.
	Value int64
	// The height of the block that created this unspent.
	BlockHeight int64
	// The date the unspent was created.
	Date string
	// The id of the wallet the unspent is in.
	Wallet string
	// The id of the wallet the unspent came from (if it was sent from a BitGo wallet you're a member on)
	FromWallet string
	// The address type and derivation path of the unspent
	// (0 = normal unspent, 1 = change unspent, 10 = segwit unspent, 11 = segwit change unspent).
	Chain int
	// The position of the address in this chain's derivation path.
	Index int
	// The script defining the criteria to be satisfied to spend this unspent.
	RedeemScript string
	// A flag indicating whether this is a segwit unspent.
	IsSegwit bool
}

// ListMeta is a pagination metadata.
type ListMeta struct {
	// Can be used to iterate the next batch of results.
	NextBatchPrevID string
	// The digital currency of the unspents.
	Coin string
}

// UnspentList is a list of unspents as retrieved from unspents endpoint.
type UnspentList struct {
	ListMeta
	Unspents []Unspent `json:"unspents"`
}

// Unspents gets a list of unspent transaction outputs (UTXOs) for a wallet.
// It invokes f for each page of results.
// You can filter unspents using query parameters as described in the docs
// https://www.bitgo.com/api/v2/#list-wallet-unspents.
func (s *walletService) Unspents(ctx context.Context, walletID string, queryParams url.Values, f func(*UnspentList)) error {
	path := fmt.Sprintf("wallet/%s/unspents", walletID)

	for {
		req, err := s.client.NewRequest(ctx, http.MethodGet, path, queryParams, nil)
		if err != nil {
			return err
		}

		v := UnspentList{}
		_, err = s.client.Do(req, &v)
		if err != nil {
			return err
		}
		f(&v)

		if v.NextBatchPrevID == "" {
			break
		}
		queryParams.Set("prevId", v.NextBatchPrevID)
	}

	return nil
}
