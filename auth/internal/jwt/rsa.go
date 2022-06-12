package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func parseRSAKeyPair(privateBytes []byte, publicBytes []byte) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateBlock, _ := pem.Decode(privateBytes)
	publicBlock, _ := pem.Decode(publicBytes)

	privateRSAKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	parsedPublicRSAKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse public key: %w", err)
	}

	publicRSAKey, ok := parsedPublicRSAKey.(*rsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected public key format")
	}

	return privateRSAKey, publicRSAKey, nil
}
