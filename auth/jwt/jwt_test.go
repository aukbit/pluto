package jwt

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestToken(t *testing.T) {

	pk, err := LoadPrivateKey("")
	if err != nil {
		t.Fatal(err)
	}

	token, err := NewToken("identifier", "bearer", "users jobs", 3650, pk)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, len(token) > 0)

	err = Verify(token, &pk.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	pubk, err := LoadPublicKey("")
	if err != nil {
		t.Fatal(err)
	}

	err = Verify(token, pubk)
	if err != nil {
		t.Fatal(err)
	}

	i := GetIdentifier(token)
	if i != "bearer" {
		t.Fatal("invalid identifier")
	}

	a := GetAudience(token)
	if a != "bearer" {
		t.Fatal("invalid audience")
	}

	s := GetScope(token)
	if s != "users jobs" {
		t.Fatal("invalid scope")
	}
}
