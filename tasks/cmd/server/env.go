package main

import (
	"fmt"
	"os"
)

type envConfig struct {
	pgConnString  string
	usersGRPCAddr string
	httpAddr      string
}

func loadEnvConfig() (*envConfig, error) {
	const (
		provideEnvErrorMsg = `please provide "%s" environment variable`

		pgConnStringEnv  = "PG_CONN"
		httpAddrEnv      = "HTTP_ADDR"
		usersGRPCAddrEnv = "USERS_GRPC_ADDR"
	)

	var ok bool

	cfg := &envConfig{}

	cfg.pgConnString, ok = os.LookupEnv(pgConnStringEnv)
	if !ok {
		return nil, fmt.Errorf(provideEnvErrorMsg, pgConnStringEnv)
	}

	cfg.httpAddr, ok = os.LookupEnv(httpAddrEnv)
	if !ok {
		return nil, fmt.Errorf(provideEnvErrorMsg, httpAddrEnv)
	}

	cfg.usersGRPCAddr, ok = os.LookupEnv(usersGRPCAddrEnv)
	if !ok {
		return nil, fmt.Errorf(provideEnvErrorMsg, usersGRPCAddrEnv)
	}

	return cfg, nil
}
