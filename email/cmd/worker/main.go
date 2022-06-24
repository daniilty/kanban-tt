package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/daniilty/kanban-tt/email/internal/emails"
	"github.com/daniilty/kanban-tt/email/internal/kafka"
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

	loggerCfg := zap.NewProductionConfig()

	logger, err := loggerCfg.Build()
	if err != nil {
		return err
	}

	splittedHost := strings.Split(cfg.smtpHost, ":")
	if len(splittedHost) != 2 {
		return errors.New("invalid host name")
	}

	host := splittedHost[0]
	port, err := strconv.Atoi(splittedHost[1])
	if err != nil {
		return err
	}

	consumer := kafka.NewConsumerImpl(cfg.kafkaTopic, []string{cfg.kafkaBroker}, cfg.kafkaGroupID)
	emailSender := emails.NewSender(host, port, &emails.User{
		Login:    cfg.smtpUser,
		Password: cfg.smtpPassword,
	}, cfg.smtpUser)
	emailsHandler := emails.NewEventsListener(logger.Sugar(), time.Duration(cfg.eventsTimeout)*time.Second, consumer, emailSender)

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		emailsHandler.Listen(ctx)
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
