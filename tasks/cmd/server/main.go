package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/daniilty/kanban-tt/schema"
	"github.com/daniilty/kanban-tt/tasks/internal/core"
	"github.com/daniilty/kanban-tt/tasks/internal/pg"
	"github.com/daniilty/kanban-tt/tasks/internal/server"
	"github.com/daniilty/kanban-tt/tasks/internal/worker"
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

	loggerCfg := zap.NewProductionConfig()

	logger, err := loggerCfg.Build()
	if err != nil {
		return err
	}

	sugaredLogger := logger.Sugar()

	ctx, cancel := context.WithCancel(context.Background())

	conn, err := grpc.DialContext(ctx, cfg.usersGRPCAddr, grpc.WithInsecure())
	if err != nil {
		cancel()

		return err
	}

	client := schema.NewUsersClient(conn)
	service := core.NewService(d, client)

	httpServer := server.NewHTTP(cfg.httpAddr, sugaredLogger, service)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		httpServer.Run(ctx)
		wg.Done()
	}()

	cleaner := worker.NewTasksCleaner(24*time.Hour, service, sugaredLogger)
	wg.Add(1)
	go func() {
		cleaner.Run(ctx)
		wg.Done()
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
	cancel()

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
