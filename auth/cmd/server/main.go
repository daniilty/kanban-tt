package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/daniilty/kanban-tt/auth/internal/core"
	"github.com/daniilty/kanban-tt/auth/internal/jwt"
	"github.com/daniilty/kanban-tt/auth/internal/kafka"
	"github.com/daniilty/kanban-tt/auth/internal/server"
	"github.com/daniilty/kanban-tt/schema"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	exitCodeInitError = 2
)

func run() error {
	cfg, err := loadEnvConfig()
	if err != nil {
		return err
	}

	manager, err := jwt.NewManagerImpl([]byte(cfg.pubKey), []byte(cfg.privKey), int64(cfg.tokenExpiry))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	conn, err := grpc.DialContext(ctx, cfg.usersGRPCAddr, grpc.WithInsecure())
	if err != nil {
		cancel()

		return err
	}

	client := schema.NewUsersClient(conn)
	kafkaProducer := kafka.NewProducer(cfg.kafkaTopic, []string{cfg.kafkaBroker}, cfg.kafkaGroupID)
	service := core.NewServiceImpl(client, manager, kafkaProducer)

	loggerCfg := zap.NewProductionConfig()

	logger, err := loggerCfg.Build()
	if err != nil {
		cancel()

		return err
	}

	httpServer := server.NewHTTP(cfg.httpAddr, logger.Sugar(), service)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(ctx context.Context) {
		httpServer.Run(ctx)
		wg.Done()
	}(ctx)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
	cancel()

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(exitCodeInitError)
	}
}
