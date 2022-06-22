package jwt

import (
	"fmt"

	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
)

var _ Manager = (*ManagerImpl)(nil)

// Manager - jwt token signer, refresher, etc...
type Manager interface {
	// Generate - generate rs256 jwt token with sub.
	Generate(*claims.Subject) (string, error)
	// Refresh - generate new token if old one is expired.
	Refresh(string) (string, error)
	// JWKS - get public key jwks so we can use it for proxy-side validation.
	JWKS() []byte
	// ParseRawToken - parse raw access token string with manager public key set.
	ParseRawToken(string) (*claims.Subject, error)
}

type ManagerImpl struct {
	alg        jwa.SignatureAlgorithm
	privateKey jwk.Key
	publicKey  jwk.Key
	tokenExp   int64
	jwksBytes  []byte
	pubSet     jwk.Set
}

func NewManagerImpl(pubKeyBytes []byte, privKeyBytes []byte, exp int64) (*ManagerImpl, error) {
	const (
		alg = jwa.RS256
		kid = "kanban_kid"
	)

	privateRSAKey, publicRSAKey, err := parseRSAKeyPair(privKeyBytes, pubKeyBytes)
	if err != nil {
		return nil, err
	}

	publicJWTKey, err := jwk.New(publicRSAKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public cert: %v", err)
	}

	privateJWTKey, err := jwk.New(privateRSAKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private cert: %v", err)
	}

	err = publicJWTKey.Set(jwk.KeyIDKey, kid)
	if err != nil {
		return nil, fmt.Errorf("failed to set token kid: %w", err)
	}

	err = publicJWTKey.Set(jwk.AlgorithmKey, alg)
	if err != nil {
		return nil, fmt.Errorf("failed to set token alg: %w", err)
	}

	err = privateJWTKey.Set(jwk.KeyIDKey, kid)
	if err != nil {
		return nil, fmt.Errorf("failed to set token kid: %w", err)
	}

	pubJWKSBytes, err := getPubSetBytes(publicJWTKey)
	if err != nil {
		return nil, fmt.Errorf("get public bytes: %w", err)
	}

	pubJWKS := jwk.NewSet()
	pubJWKS.Add(publicJWTKey)

	return &ManagerImpl{
		alg:        alg,
		privateKey: privateJWTKey,
		publicKey:  publicJWTKey,
		tokenExp:   exp,
		jwksBytes:  pubJWKSBytes,
		pubSet:     pubJWKS,
	}, nil
}
