package main

import (
	"fmt"
	"log"
	"os"
)

type envConfig struct {
	pgConnString string
	grpcAddr     string
}

func loadEnvConfig() (*envConfig, error) {
	const (
		provideEnvErrorMsg = `please provide "%s" environment variable`

		pgConnStringEnv = "PG_CONN"
		grpcAddrEnv     = "GRPC_SERVER_ADDR"
	)

	var ok bool

	cfg := &envConfig{}

	cfg.pgConnString, ok = os.LookupEnv(pgConnStringEnv)
	if !ok {
		return nil, fmt.Errorf(provideEnvErrorMsg, pgConnStringEnv)
	}

	log.Println(cfg.pgConnString)

	cfg.grpcAddr, ok = os.LookupEnv(grpcAddrEnv)
	if !ok {
		return nil, fmt.Errorf(provideEnvErrorMsg, grpcAddrEnv)
	}

	return cfg, nil
}
