package main

import (
	"fmt"
	"os"
	"strconv"
)

type envConfig struct {
	httpAddr      string
	usersGRPCAddr string
	pubKey        string
	privKey       string
	tokenExpiry   int
}

func loadEnvConfig() (*envConfig, error) {
	var err error

	cfg := &envConfig{}

	cfg.httpAddr, err = lookupEnv("HTTP_SERVER_ADDR")
	if err != nil {
		return nil, err
	}

	cfg.usersGRPCAddr, err = lookupEnv("USERS_GRPC_ADDR")
	if err != nil {
		return nil, err
	}

	cfg.pubKey, err = lookupEnv("PUBKEY")
	if err != nil {
		return nil, err
	}

	cfg.privKey, err = lookupEnv("PRIVKEY")
	if err != nil {
		return nil, err
	}

	var expiryString string

	expiryString, err = lookupEnv("TOKEN_EXPIRY")
	if err != nil {
		return nil, err
	}

	cfg.tokenExpiry, err = strconv.Atoi(expiryString)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func lookupEnv(name string) (string, error) {
	const provideEnvErrorMsg = `please provide "%s" environment variable`

	val, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf(provideEnvErrorMsg, name)
	}

	return val, nil
}
