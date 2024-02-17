package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	default_port = "9090"
)

type server struct {
	logger *zap.Logger
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	s.logger.Info("Received greeting", zap.String("name", in.GetName()))
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func mainServer(logger *zap.Logger) {
	port := strings.Trim(os.Getenv("BIND_PORT"), " ")
	if port == "" {
		port = default_port
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Error("failed to listen: %v", zap.Error(err))
		return
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{logger: logger})
	grpc_health_v1.RegisterHealthServer(s, &healthChecker{time.Now()})
	logger.Info("server listening", zap.String("address", lis.Addr().String()))
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
