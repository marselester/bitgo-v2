# Go client for [BitGo.com API v2](https://www.bitgo.com/api/v2/)

[![Documentation](https://godoc.org/github.com/marselester/bitgo-v2?status.svg)](https://godoc.org/github.com/marselester/bitgo-v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/marselester/bitgo-v2)](https://goreportcard.com/report/github.com/marselester/bitgo-v2)
[![Travis CI](https://travis-ci.org/marselester/bitgo-v2.png)](https://travis-ci.org/marselester/bitgo-v2)

This is unofficial API client. There are no plans to implement all resources.

## [Consolidate Wallet Unspents](https://www.bitgo.com/api/v2/#consolidate-wallet-unspents)

This API call will consolidate Bitcoin Cash of `585951a5df8380e0e3063e9f` wallet using max `0.001` BCH unspents
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

## [List Wallet Unspents](https://www.bitgo.com/api/v2/#list-wallet-unspents)

This API call will retrieve the unspent transaction outputs (UTXOs) within a wallet.
For example, we want to request `58ae81a5df8380e0e307e876` wallet's confirmed unspents and
print amounts in BTC. You can stop pagination by cancelling a `ctx` context.

```go
c := bitgo.NewClient(
	bitgo.WithCoin("bch"),
	bitgo.WithAccesToken("swordfish"),
)
params := url.Values{}
params.Set("minConfirms", "1")
err := c.Wallet.Unspents(ctx, "58ae81a5df8380e0e307e876", params, func(list *bitgo.UnspentList) {
	for _, utxo := range list.Unspents {
		fmt.Printf("%0.8f\n", toBitcoins(utxo.Value))
	}
})
```

There is a CLI program to list all unspensts of a wallet.

```sh
$ go build ./cmd/utxo/
$ ./utxo -token=swordfish -coin=bch -wallet=58ae81a5df8380e0e307e876 -min-confirms=1
0.00000117
0.00000001
0.00000001
0.00000562
0.00000001
0.00000562
```

You can use it to get a rough idea about unspents available in the wallet.

```sh
$ ./utxo -token=swordfish -coin=bch -wallet=58ae81a5df8380e0e307e876 > unspents.txt
$ cat unspents.txt | sort | uniq -c | sort -n -r
   3 0.00000001
   2 0.00000562
   1 0.00000117
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

## Testing

Quick tutorial on [how to fuzz](https://medium.com/@dgryski/go-fuzz-github-com-arolek-ase-3c74d5a3150c) by Damian Gryski.
Copy JSON files from `testdata` into `workdir/corpus` as sample inputs.

```sh
$ go-fuzz-build github.com/marselester/bitgo-v2
$ go-fuzz -bin=bitgo-fuzz.zip -workdir=workdir
```
