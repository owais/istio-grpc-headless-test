package main

import (
	"context"
	"math/rand/v2"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func mainClient(logger *zap.Logger) {
	servers := []string{}
	for _, s := range strings.Split(os.Getenv("SERVERS"), ",") {
		s = strings.Trim(s, " ")
		if s != "" {
			servers = append(servers, s)
		}
	}
	if len(servers) == 0 {
		logger.Error("no SERVERS env var found")
	}

	for _, s := range servers {
		go spawnClient(logger, strings.Trim(s, " "))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(1)

}

func spawnClient(logger *zap.Logger, server_address string) {
	conn, err := grpc.Dial(server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect", zap.Error(err))
		return
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		if r, err := c.SayHello(ctx, &pb.HelloRequest{Name: server_address}); err != nil {
			logger.Error("could not greet", zap.Error(err))
		} else {
			logger.Info("Received Greeting", zap.String("message", r.GetMessage()))
		}
		// sleep for 5-10 seconds
		time.Sleep(5 + rand.N(10*time.Second))
	}
}
