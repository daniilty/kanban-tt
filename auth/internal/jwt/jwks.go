package jwt

import (
	"encoding/json"

	"github.com/lestrrat-go/jwx/jwk"
)

type jwks struct {
	Keys []jwk.Key `json:"keys"`
}

func getPubSetBytes(pubKey jwk.Key) ([]byte, error) {
	pubKeySet := &jwks{
		Keys: []jwk.Key{pubKey},
	}

	return json.MarshalIndent(pubKeySet, "", "  ")
}

func (m *ManagerImpl) JWKS() []byte {
	return m.jwksBytes
}
