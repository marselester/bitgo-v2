// List bitcoin unspent transaction outputs (UTXOs).
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marselester/bitgo-v2"
)

func main() {
	baseURL := flag.String("host", "http://0.0.0.0:3080", "BitGo API server base URL.")
	accessToken := flag.String("token", "", "BitGo access token.")
	coin := flag.String("coin", "btc", "Coin identifier.")
	walletID := flag.String("wallet", "", "BitGo wallet ID.")
	prevID := flag.String("prev-id", "", "Continue iterating unspents from this ID as provided by nextBatchPrevId in the previous list.")
	minSize := flag.Float64("min-size", 0, "Ignore unspents smaller than this amount of bitcoins.")
	maxSize := flag.Float64("max-size", 0, "Ignore unspents larger than this amount of bitcoins.")
	minHeight := flag.Int("min-height", 0, "Ignore unspents confirmed at a lower block height than the given height.")
	minConfirms := flag.Int("min-confirms", 0, "Ignore unspents that have fewer than the given confirmations.")
	waitSeconds := flag.Int("wait", 15, "How many seconds to wait after failed download attempt.")
	debug := flag.Bool("debug", false, "Enable debug mode.")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Listen to Ctrl+C and kill/killall to gracefully stop listing unspents.
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan

		log.Print("utxo: stopping...")
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

	params := url.Values{}
	if *prevID != "" {
		params.Set("prevId", *prevID)
	}
	if *minSize > 0 {
		params.Set("minValue", fmt.Sprintf("%d", bitgo.ToSatoshis(*minSize)))
	}
	if *maxSize > 0 {
		params.Set("maxValue", fmt.Sprintf("%d", bitgo.ToSatoshis(*maxSize)))
	}
	if *minHeight > 0 {
		params.Set("minHeight", fmt.Sprintf("%d", *minHeight))
	}
	if *minConfirms > 0 {
		params.Set("minConfirms", fmt.Sprintf("%d", *minConfirms))
	}

	downloaded := 0
	nextBatchPrevID := ""
	for {
		err := client.Wallet.Unspents(ctx, *walletID, params, func(list *bitgo.UnspentList) {
			downloaded += len(list.Unspents)
			log.Printf("utxo: fetched %d unspents", downloaded)

			for _, utxo := range list.Unspents {
				fmt.Printf("%0.8f\n", bitgo.ToBitcoins(utxo.Value))
			}

			nextBatchPrevID = list.NextBatchPrevID
		})
		// Stop when we downloaded everything without errors or
		// when a context was cancelled (user hit Ctrl+C).
		if err == nil || ctx.Err() != nil {
			break
		}

		if apiErr, ok := err.(bitgo.Error); ok {
			log.Printf("utxo: failed to list unspents, %d: %v", apiErr.HTTPStatusCode, apiErr)
		} else {
			log.Printf("utxo: failed to list unspents: %v", err)
		}

		// We shall wait a bit and then try again.
		log.Printf("utxo: retrying in %d seconds...", *waitSeconds)
		time.Sleep(time.Duration(*waitSeconds) * time.Second)
		if nextBatchPrevID != "" {
			params.Set("prevId", nextBatchPrevID)
		}
	}
}

// StdLogger prints logs to standard error.
func StdLogger(keyvals ...interface{}) error {
	log.Printf("%q", keyvals)
	return nil
}
