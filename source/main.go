package main

import (
	"os"

	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
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
