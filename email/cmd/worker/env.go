package main

import (
	"fmt"
	"os"
	"strconv"
)

type envConfig struct {
	kafkaBroker   string
	kafkaTopic    string
	kafkaGroupID  string
	smtpHost      string
	smtpUser      string
	smtpPassword  string
	eventsTimeout int
}

func loadEnvConfig() (*envConfig, error) {
	var err error

	cfg := &envConfig{}

	cfg.kafkaBroker, err = lookupEnv("KAFKA_BROKER")
	if err != nil {
		return nil, err
	}

	cfg.kafkaTopic, err = lookupEnv("KAFKA_TOPIC")
	if err != nil {
		return nil, err
	}

	cfg.smtpHost, err = lookupEnv("SMTP_HOST")
	if err != nil {
		return nil, err
	}

	cfg.smtpUser, err = lookupEnv("SMTP_USER")
	if err != nil {
		return nil, err
	}

	cfg.smtpPassword, err = lookupEnv("SMTP_PASSWORD")
	if err != nil {
		return nil, err
	}

	cfg.kafkaGroupID, err = lookupEnv("KAFKA_GROUP_ID")
	if err != nil {
		return nil, err
	}

	timeoutString, err := lookupEnv("TIMEOUT")
	if err != nil {
		return nil, err
	}

	cfg.eventsTimeout, err = strconv.Atoi(timeoutString)
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
