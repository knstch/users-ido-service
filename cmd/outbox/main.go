package main

import (
	"fmt"
	defaultLog "log"
	"os"
	"path/filepath"

	"github.com/knstch/knstch-libs/log"

	"github.com/knstch/knstch-ido-kafka/outbox"

	"users-service/config"
)

func main() {
	if err := run(); err != nil {
		defaultLog.Println(err)
	}
}

func run() error {
	args := os.Args

	dir, err := filepath.Abs(filepath.Dir(args[0]))
	if err != nil {
		return fmt.Errorf("filepath.Abs: %w", err)
	}

	if err = config.InitENV(dir); err != nil {
		return fmt.Errorf("config.InitENV: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("config.GetConfig: %w", err)
	}

	logger := log.NewLogger(cfg.ServiceName, log.InfoLevel)

	listener, err := outbox.NewOutboxListener(cfg.KafkaAddr, cfg.GetDSN(), logger)
	if err != nil {
		return fmt.Errorf("outbox.NewOutboxListener: %w", err)
	}

	listener.Start()

	select {}
}
