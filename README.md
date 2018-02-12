# Go client for [BitGo.com API v2](https://www.bitgo.com/api/v2/)

[![Documentation](https://godoc.org/github.com/marselester/bitgo-v2?status.svg)](https://godoc.org/github.com/marselester/bitgo-v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/marselester/bitgo-v2)](https://goreportcard.com/report/github.com/marselester/bitgo-v2)
[![Travis CI](https://travis-ci.org/marselester/bitgo-v2.png)](https://travis-ci.org/marselester/bitgo-v2)

This is unofficial API client. There are no plans to implement all resources.

## [Consolidate Wallet Unspents](https://www.bitgo.com/api/v2/#consolidate-wallet-unspents)

Consolidates Bitcoin Cash of `585951a5df8380e0e3063e9f` wallet using max `0.001` BCH unspents
with 5000 satoshis/kb fee rate. You can stop consolidation by cancelling a `ctx` context.

```go
c := bitgo.NewClient(
    bitgo.WithCoin("bch"),
    bitgo.WithAccesToken("swordfish"),
)
tx, err := c.Wallet.Consolidate(ctx, "585951a5df8380e0e3063e9f", &bitgo.WalletConsolidateParams{
    WalletPassphrase: "root",
    MaxValue:         100000,
    FeeRate:          5000,
})
if err != nil {
    log.Fatalf("Failed to coalesce unspents: %v", err)
}
fmt.Printf("Consolidated transaction ID: %s", tx.TxID)
```

There is a CLI program to consolidate unspensts of a wallet.

```sh
$ go build ./cmd/consolidate/
$ ./consolidate -token=swordfish -coin=bch -wallet=585951a5df8380e0e3063e9f -passphrase=root -max-value=0.001 -fee-rate=5000
5885a7e6c7802206f69655ed763d14f101cf46501aef38e275c67c72cfcedb75
```

## Error Handling

Dave Cheney recommends
[asserting errors for behaviour](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully), not type.

```go
package main

import (
	"fmt"

	"github.com/marselester/bitgo-v2"
	"github.com/pkg/errors"
)

// IsUnauthorized returns true if err caused by authentication problem.
func IsUnauthorized(err error) bool {
	e, ok := errors.Cause(err).(interface {
		IsUnauthorized() bool
	})
	return ok && e.IsUnauthorized()
}

func main() {
	err := bitgo.Error{Type: bitgo.ErrorTypeAuthentication}
	fmt.Println(IsUnauthorized(err))
	fmt.Println(IsUnauthorized(fmt.Errorf("")))
	fmt.Println(IsUnauthorized(nil))
	// Output: true
	// false
	// false
}
```
