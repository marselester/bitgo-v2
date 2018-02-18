// +build gofuzz

package bitgo

import (
	"context"
	"net/http"
	"net/http/httptest"
)

func Fuzz(data []byte) int {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer srv.Close()

	c := NewClient(
		WithBaseURL(srv.URL),
	)
	ctx := context.Background()
	if _, err := c.Do(c.NewRequest(ctx, http.MethodPost, "", nil, nil)); err != nil {
		return 0
	}
	return 1
}
