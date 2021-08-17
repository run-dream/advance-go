package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

// GrpcServer
type GrpcServer struct {
	*grpc.Server
	lis *net.Listener
}

func NewGrpcServer(s *grpc.Server, lis *net.Listener) *GrpcServer {
	return &GrpcServer{s, lis}
}

func (s *GrpcServer) Start(c context.Context) error {
	log.Println("正在启动 grpc 服务")
	return s.Serve(*s.lis)
}

func (s *GrpcServer) Stop(c context.Context) error {
	log.Println("正在关闭 grpc 服务")
	s.GracefulStop()
	return nil
}
