package jws

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
)

func TestSignAndVerify(t *testing.T) {
	header := &Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}
	payload := &ClaimSet{
		Iss: "http://google.com/",
		Aud: "",
		Exp: time.Now().Unix() + 2,
		Iat: time.Now().Unix(),
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	token, err := Encode(header, payload, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	err = Verify(token, &privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	c, err := Decode(token)
	if err != nil {
		t.Fatal(err)
	}
	if c.Iss != "http://google.com/" {
		t.Error("invalid decode")
	}
	if time.Now().Unix() > c.Exp {
		t.Error("token as expired")
	}
}

func TestVerifyFailsOnMalformedClaim(t *testing.T) {
	err := Verify("abc.def", nil)
	if err == nil {
		t.Error("got no errors; want improperly formed JWT not to be verified")
	}
}
