package bitgo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marselester/bitgo-v2"
)

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		statusCode int
		want       bitgo.Error
	}{
		{
			name:       "401 unauthorized error",
			body:       `{"error":"unauthorized","name":"","requestId":"bj9h0dap1723kadrsnfkvsinz","message":"unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			want: bitgo.Error{
				Type:           bitgo.ErrorTypeAuthentication,
				HTTPStatusCode: http.StatusUnauthorized,
				Body:           `{"error":"unauthorized","name":"","requestId":"bj9h0dap1723kadrsnfkvsinz","message":"unauthorized"}` + "\n",
				Message:        "unauthorized",
				RequestID:      "bj9h0dap1723kadrsnfkvsinz",
			},
		},
		{
			name:       "500 temporary API error",
			body:       "some internal server error",
			statusCode: http.StatusInternalServerError,
			want: bitgo.Error{
				Type:           bitgo.ErrorTypeAPI,
				HTTPStatusCode: 500,
				Body:           "some internal server error\n",
				Message:        "",
				RequestID:      "",
			},
		},
		{
			name:       "503 temporary API error",
			body:       "server is overloaded",
			statusCode: http.StatusServiceUnavailable,
			want: bitgo.Error{
				Type:           bitgo.ErrorTypeAPI,
				HTTPStatusCode: 503,
				Body:           "server is overloaded\n",
				Message:        "",
				RequestID:      "",
			},
		},
		{
			name:       "504 temporary API error",
			body:       "",
			statusCode: http.StatusGatewayTimeout,
			want: bitgo.Error{
				Type:           bitgo.ErrorTypeAPI,
				HTTPStatusCode: 504,
				Body:           "\n",
				Message:        "",
				RequestID:      "",
			},
		},
		{
			name:       "catch all errors",
			body:       "Cloudflare",
			statusCode: 522,
			want: bitgo.Error{
				Type:           bitgo.ErrorTypeAPI,
				HTTPStatusCode: 522,
				Body:           "Cloudflare\n",
				Message:        "",
				RequestID:      "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, test.body, test.statusCode)
			}))
			defer srv.Close()

			c := bitgo.NewClient(
				bitgo.WithBaseURL(srv.URL),
			)
			_, err := c.Wallet.Consolidate(context.Background(), "", nil)
			if err != test.want {
				t.Fatalf("should be %#v, not %#v", test.want, err)
			}
		})
	}
}
