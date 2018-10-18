package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/aukbit/pluto/auth/jws"
)

var (
	privKeyPath = "./keys/auth.rsa"
	pubKeyPath  = "./keys/auth.rsa.pub"
)

var (
	ErrExpiredToken           = errors.New("token has expired")
	ErrInvalidAudience        = errors.New("token has invalid audience")
	ErrInvalidIdentifier      = errors.New("token has invalid identifier")
	ErrPrivateKeyNotAvailable = errors.New("private key not available in context")
	ErrPublicKeyNotAvailable  = errors.New("public key not available in context")
)

// LoadPublicKey loads a public key from PEM encoded data.
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	if path == "" {
		path = pubKeyPath
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		err = errors.New("Invalid PublicKey format")
		return nil, err
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

// LoadPrivateKey loads a private key from PEM encoded data.
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	if path == "" {
		path = privKeyPath
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		err = errors.New("Invalid PrivateKey format")
		return nil, err
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

type ClaimSet struct {
	Identifier, Audience, Scope, Jti, Principal string
	Expiration                                  int64
}

// NewToken returns a JWT token signed with the given RSA private key.
func NewToken(ctx context.Context, cs *ClaimSet) (string, error) {
	prv, ok := PrivateKeyFromContext(ctx)
	if !ok {
		return "", ErrPrivateKeyNotAvailable
	}
	header := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}
	payload := &jws.ClaimSet{
		Iss:   cs.Identifier,
		Aud:   cs.Audience,
		Scope: cs.Scope,
		Exp:   time.Now().Unix() + cs.Expiration,
		Iat:   time.Now().Unix(),
		Sub:   cs.Jti,
		Prn:   cs.Principal,
	}
	t, err := jws.Encode(header, payload, prv)
	if err != nil {
		return "", err
	}

	return t, nil
}

// Verify tests whether the provided JWT token's signature was produced by the
// private key associated with the supplied public key.
// Also verifies if Token as expired
func Verify(ctx context.Context, token string) error {
	pub, ok := PublicKeyFromContext(ctx)
	if !ok {
		return ErrPublicKeyNotAvailable
	}
	err := jws.Verify(token, pub)
	if err != nil {
		return err
	}
	c, err := jws.Decode(token)
	if err != nil {
		return err
	}
	if time.Now().Unix() > c.Exp {
		return ErrExpiredToken
	}
	return nil
}

// Identifier the "iss" (issuer) claim identifies the principal that issued the JWT.
func Identifier(token string) string {
	c, _ := jws.Decode(token)
	return c.Iss
}

// Scope space-delimited list of the permissions the application requests.
func Scope(token string) string {
	c, _ := jws.Decode(token)
	return c.Scope
}

// Audience The "aud" (audience) claim identifies the audience that the JWT is
// intended for.
func Audience(token string) string {
	c, _ := jws.Decode(token)
	return c.Aud
}

// Principal The "prn" (principal) claim identifies the subject of the JWT.
func Principal(token string) string {
	c, _ := jws.Decode(token)
	return c.Prn
}

// Jti The "jti" (JWT ID) claim provides a unique identifier for the JWT.
func Jti(token string) string {
	c, _ := jws.Decode(token)
	return c.Sub
}

// BearerAuth returns the token provided in the request's
// Authorization header, if the request uses HTTP Bearer Authentication.
func BearerAuth(r *http.Request) (token string, ok bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return
	}
	return parseBearerAuth(auth)
}

// parseBearerAuth parses an HTTP Bearer Authentication string.
// "Bearer QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("QWxhZGRpbjpvcGVuIHNlc2FtZQ==", true).
func parseBearerAuth(auth string) (token string, ok bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	return auth[len(prefix):], true
}
