package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/daniilty/kanban-tt/schema"
	"github.com/daniilty/kanban-tt/users/internal/core"
	"github.com/daniilty/kanban-tt/users/internal/pg"
	"github.com/daniilty/kanban-tt/users/internal/server"
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

	d, err := pg.Connect(context.Background(), cfg.pgConnString)
	if err != nil {
		return err
	}

	service := core.NewServiceImpl(d)

	loggerCfg := zap.NewProductionConfig()

	logger, err := loggerCfg.Build()
	if err != nil {
		return err
	}

	sugaredLogger := logger.Sugar()

	wg := &sync.WaitGroup{}
	listener, err := net.Listen("tcp", cfg.grpcAddr)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	grpcService := server.NewGRPC(service)
	schema.RegisterUsersServer(grpcServer, grpcService)

	sugaredLogger.Infow("GRPC server is starting.", "addr", listener.Addr())

	wg.Add(1)
	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			sugaredLogger.Errorw("Server failed to start.", "err", err)
		}
		wg.Done()
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	sugaredLogger.Info("Gracefully stopping GRPC server.")
	grpcServer.GracefulStop()

	wg.Wait()

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(exitCodeInitError)
	}
}
