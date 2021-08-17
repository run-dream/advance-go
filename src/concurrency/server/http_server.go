package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

// HttpServer
type HttpServer struct {
	*http.Server
}

func NewHttpServer(server *http.Server) *HttpServer {
	return &HttpServer{server}
}

func (s *HttpServer) Start(c context.Context) error {
	log.Println("正在启动 http 服务")
	return s.ListenAndServe()
}

func (s *HttpServer) Stop(c context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	log.Println("正在关闭 http 服务")
	return s.Shutdown(timeoutCtx)
}
