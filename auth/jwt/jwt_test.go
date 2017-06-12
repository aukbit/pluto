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

	token, err := NewToken("identifier", 3650, pk)
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
}
