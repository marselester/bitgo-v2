// Consolidate the unspents currently held in a wallet to a smaller number.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/marselester/bitgo-v2"
)

func main() {
	baseURL := flag.String("host", "http://0.0.0.0:3080", "BitGo Express API server base URL.")
	accessToken := flag.String("token", "", "BitGo access token.")
	coin := flag.String("coin", "btc", "Coin identifier.")
	walletID := flag.String("wallet", "", "BitGo wallet ID.")
	walletPassphrase := flag.String("passphrase", "", "Passphrase of the wallet.")
	numUnspentsToMake := flag.Int("target", 1, "Number of outputs created by the consolidation transaction.")
	limit := flag.Int("limit", 25, "Number of unspents to select (max is 200).")
	minValue := flag.Float64("min-value", 0, "Ignore unspents smaller than this amount of bitcoins.")
	maxValue := flag.Float64("max-value", 0, "Ignore unspents larger than this amount of bitcoins.")
	minHeight := flag.Int("min-height", 0, "The minimum height of unspents on the block chain to use.")
	feeRate := flag.Int("fee-rate", 0, "The desired fee rate for the transaction in satoshis/KB.")
	feeTxConfirmTarget := flag.Int(
		"fee-tx-confirm-target",
		0,
		`Fee rate is automatically chosen by targeting a transaction confirmation in this number of blocks
(only available on BTC, fee-rate takes precedence if also set).`,
	)
	maxFeePercentage := flag.Int("max-fee-percentage", 0, "Maximum percentage of an unspent's value to be used for fees. Cannot be combined with min-value.")
	minConfirms := flag.Int("min-confirms", 0, "The required number of confirmations for each transaction input.")
	enforceMinConfirmsForChange := flag.Bool("enforce-min-confirms-for-change", false, "Apply the required confirmations set in min-confirms for change outputs.")
	maxIter := flag.Int("max-iter", 1, "Maximum number of consolidation iterations to perform.")
	debug := flag.Bool("debug", false, "Enable debug mode.")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Listen to Ctrl+C and kill/killall to gracefully stop consolidation.
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan

		log.Print("consolidate: stopping...")
		cancel()
	}()

	var logger bitgo.Logger
	if *debug {
		logger = bitgo.LoggerFunc(StdLogger)
	} else {
		logger = &bitgo.NoopLogger{}
	}
	client := bitgo.NewClient(
		bitgo.WithBaseURL(*baseURL),
		bitgo.WithCoin(*coin),
		bitgo.WithAccesToken(*accessToken),
		bitgo.WithLogger(logger),
	)

	params := &bitgo.WalletConsolidateParams{
		WalletPassphrase:            *walletPassphrase,
		NumUnspentsToMake:           *numUnspentsToMake,
		Limit:                       *limit,
		MinValue:                    bitgo.ToSatoshis(*minValue),
		MaxValue:                    bitgo.ToSatoshis(*maxValue),
		MinHeight:                   *minHeight,
		FeeRate:                     *feeRate,
		FeeTxConfirmTarget:          *feeTxConfirmTarget,
		MaxFeePercentage:            *maxFeePercentage,
		MinConfirms:                 *minConfirms,
		EnforceMinConfirmsForChange: *enforceMinConfirmsForChange,
	}

	for i := 0; i < *maxIter; i++ {
		tx, err := client.Wallet.Consolidate(ctx, *walletID, params)
		// Print consolidated transaction ID.
		if err == nil {
			fmt.Printf("%s\n", tx.TxID)
			continue
		}

		// Stop when a context was cancelled (user hit Ctrl+C).
		if ctx.Err() != nil {
			break
		}

		if apiErr, ok := err.(bitgo.Error); ok {
			log.Fatalf("consolidate: failed to coalesce unspents, %d: %v", apiErr.HTTPStatusCode, apiErr)
		}
		log.Fatalf("consolidate: failed to coalesce unspents: %v", err)
	}
}

// StdLogger prints logs to standard error.
func StdLogger(keyvals ...interface{}) error {
	log.Printf("%q", keyvals)
	return nil
}
