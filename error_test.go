package bitgo_test

import (
	"testing"

	"github.com/marselester/bitgo-v2"
)

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		err  bitgo.Error
		want bool
	}{
		{bitgo.Error{Type: bitgo.ErrorTypeAuthentication}, true},
		{bitgo.Error{Type: bitgo.ErrorTypeInvalidRequest}, false},
		{bitgo.Error{Type: bitgo.ErrorTypeRateLimit}, false},
		{bitgo.Error{Type: bitgo.ErrorTypeAPI}, false},
		{bitgo.Error{}, false},
	}
	for _, test := range tests {
		got := test.err.IsUnauthorized()
		if got != test.want {
			t.Errorf("IsUnauthorized(%#v) = %v, want %v", test.err, got, test.want)
		}
	}
}
