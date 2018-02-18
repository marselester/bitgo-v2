package bitgo_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/marselester/bitgo-v2"
)

func TestConsolidate(t *testing.T) {
	filename := filepath.Join("testdata", "consolidateunspents.json")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	want := bitgo.TxInfo{
		TxID:   "5885a7e6c7802206f69655ed763d14f101cf46501aef38e275c67c72cfcedb75",
		Tx:     "010000001945f506cad0d2b8be8e36ed5d4dbb7502cd3a2822d01f7afe411f9773c6fd381cad0 ...",
		Status: "signed",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer srv.Close()

	c := bitgo.NewClient(
		bitgo.WithBaseURL(srv.URL),
	)
	tx, err := c.Wallet.Consolidate(context.Background(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if *tx != want {
		t.Fatalf("should be %#v, not %#v", want, tx)
	}
}

func TestUnspents(t *testing.T) {
	filename := filepath.Join("testdata", "unspents.json")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	want := bitgo.Unspent{
		ID:           "952ac7fd9c1a5df8380e0e305fac8b42db:0",
		Address:      "2NEqutgZ741a5df8380e0e30gkrM9vAyn3",
		Value:        203125000,
		BlockHeight:  999999999,
		Date:         "2017-02-23T06:59:21.538Z",
		Wallet:       "58ae81a5df8380e0e307e876",
		FromWallet:   "",
		Chain:        0,
		Index:        0,
		RedeemScript: "522102f601b186b23d6c7b1fc3a3363a7e47b1a48e13e559601c9cf22c98b249c288bf210385dd4200926a87b1363667d50b8e46d17f811ee7bed3c5c29607545f231233d521036ed3744f71e371796b8dfea84dbeeb49a270339ec34eb9a92b87b6874674ecb357ae",
		IsSegwit:     false,
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer srv.Close()

	c := bitgo.NewClient(
		bitgo.WithBaseURL(srv.URL),
	)
	err = c.Wallet.Unspents(context.Background(), "", nil, func(list *bitgo.UnspentList) {
		got := list.Unspents[0]
		if got != want {
			t.Fatalf("should be %#v, not %#v", want, got)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkUnspents(b *testing.B) {
	filename := filepath.Join("testdata", "unspents.json")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		b.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer srv.Close()

	c := bitgo.NewClient(
		bitgo.WithBaseURL(srv.URL),
	)
	f := func(list *bitgo.UnspentList) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Wallet.Unspents(context.Background(), "", nil, f)
	}
}
