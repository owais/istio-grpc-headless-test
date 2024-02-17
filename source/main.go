package main

import (
	"os"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.InfoLevel,
	))
	defer logger.Sync()

	args := os.Args[1:]
	if len(args) == 0 {
		logger.Fatal("available commands: server, client")
	}

	switch args[0] {
	case "server":
		mainServer(logger)
	case "client":
		mainClient(logger)
	default:
		logger.Fatal("available commands: server, client")
	}
}
