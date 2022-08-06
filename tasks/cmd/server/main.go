package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/daniilty/kanban-tt/tasks/internal/core"
	"github.com/daniilty/kanban-tt/tasks/internal/pg"
	"github.com/daniilty/kanban-tt/tasks/internal/server"
	"github.com/daniilty/kanban-tt/tasks/internal/worker"
	"go.uber.org/zap"
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

	service := core.NewService(d)

	loggerCfg := zap.NewProductionConfig()

	logger, err := loggerCfg.Build()
	if err != nil {
		return err
	}

	sugaredLogger := logger.Sugar()
	httpServer := server.NewHTTP(cfg.httpAddr, sugaredLogger, service)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		httpServer.Run(ctx)
		wg.Done()
	}()

	cleaner := worker.NewTasksCleaner(24*time.Hour, d, sugaredLogger)
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
