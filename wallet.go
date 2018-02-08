package bitgo

import (
	"context"
	"fmt"
	"net/http"
)

// Satoshi is the smallest unit of bitcoin.
const Satoshi = 0.00000001

// ToBitcoins converts satoshis to bitcoins.
func ToBitcoins(amount int64) float64 {
	return float64(amount) * Satoshi
}

// ToSatoshis converts bitcoins to satoshis.
func ToSatoshis(amount float64) int64 {
	return int64(amount / Satoshi)
}

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
