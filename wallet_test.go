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
