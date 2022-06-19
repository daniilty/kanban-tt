package main

import (
	"fmt"
	"os"
)

type envConfig struct {
	pgConnString string
	httpAddr     string
}

func loadEnvConfig() (*envConfig, error) {
	const (
		provideEnvErrorMsg = `please provide "%s" environment variable`

		pgConnStringEnv = "PG_CONN"
		httpAddrEnv     = "HTTP_ADDR"
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

	return cfg, nil
}
