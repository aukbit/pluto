package jwt

import (
	"context"
	"testing"

	"github.com/paulormart/assert"
)

func TestToken(t *testing.T) {

	prv, err := LoadPrivateKey("")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, PrivateKeyContextKey, prv)

	cs := &ClaimSet{
		Identifier: "identifier",
		Audience:   "bearer",
		Scope:      "users jobs",
		Jti:        "123",
		Principal:  "principal",
		Expiration: 3650,
	}
	token, err := NewToken(ctx, cs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, token != "")

	pub, err := LoadPublicKey("")
	if err != nil {
		t.Fatal(err)
	}
	ctx = context.WithValue(ctx, PublicKeyContextKey, pub)

	err = Verify(ctx, token)
	if err != nil {
		t.Fatal(err)
	}

	i := Identifier(token)
	if i != "identifier" {
		t.Fatal("invalid identifier")
	}

	a := Audience(token)
	if a != "bearer" {
		t.Fatal("invalid audience")
	}

	s := Scope(token)
	if s != "users jobs" {
		t.Fatal("invalid scope")
	}

	j := Jti(token)
	if j != "123" {
		t.Fatal("invalid scope")
	}

	p := Principal(token)
	if p != "principal" {
		t.Fatal("invalid scope")
	}
}
