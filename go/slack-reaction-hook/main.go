package main

import (
	"os"

	"go.uber.org/zap"
)

const (
	AppToken        = "xapp-xxxxx"
	BotToken        = "xoxb-xxxxx"
	Port            = "8080"
	TriggerReaction = "bankai"
)

var (
	log Logger
)

// Logger ... zap logger
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func init() {
	// logger初期化
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	s := logger.Sugar()
	log = Logger{logger, s}
}

func main() {
	if err := run(); err != nil {
		log.sugar.Errorf("%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	log.sugar.Infof("Server is starting ... :%s", Port)
	// start server
	if err := Server(Port); err != nil {
		return err
	}

	return nil
}
