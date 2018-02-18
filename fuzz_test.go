package bitgo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoCrashers(t *testing.T) {
	crashers := []string{
		"{\n    \"txid\": \"5885a" +
			"7e6c7802206f69655ed7" +
			"63d14f101cf46501aef3" +
			"8e275c67c72cfcedb75\"" +
			",\n    \"tx\": \"0100000" +
			"01945f506cfd381cad0 " +
			"...\",\n    \"status\": " +
			"\"signed\"\n}\n",
		"\"\\u\\uE.88\\uE644\\uE9&" +
			"8\\uT",
	}

	for _, data := range crashers {
		t.Run("", func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(data))
			}))
			defer srv.Close()

			c := NewClient(
				WithBaseURL(srv.URL),
			)
			ctx := context.Background()
			c.Do(c.NewRequest(ctx, http.MethodPost, "", nil, nil))
		})
	}
}
