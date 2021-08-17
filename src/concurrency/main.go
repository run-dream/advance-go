package main

import (
	"context"
	"fmt"
	"go-advance/src/concurrency/api/v1"
	"go-advance/src/concurrency/server"
	"go-advance/src/concurrency/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// newHttpServer 初始化一个http 服务
func newHttpServer() (server.Server, <-chan struct{}) {
	mux := http.NewServeMux()
	shutdown := make(chan struct{})
	mux.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		log.Println("http server 还存活")
		w.Write([]byte("still alive"))
	})
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		log.Println("http server 收到了 shutdown 的请求")
		shutdown <- struct{}{}
	})
	httpSever := http.Server{
		Handler: mux,
		Addr:    ":8001",
	}
	return server.NewHttpServer(&httpSever), shutdown
}

func newGrpcServer() server.Server {
	rpcServer := grpc.NewServer()

	api.RegisterHelloServiceServer(rpcServer, new(service.HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	return server.NewGrpcServer(rpcServer, &lis)
}

func main() {
	// 用来管理 goroutine
	group, ctx := errgroup.WithContext(context.Background())

	// 启动 http 服务
	httpServer, shutdown := newHttpServer()

	group.Go(func() error {
		return httpServer.Start(ctx)
	})

	// 启动 grpc 服务
	grpcServer := newGrpcServer()
	group.Go(func() error {
		return grpcServer.Start(ctx)
	})

	// 监听操作系统信号
	group.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		select {
		// 这里通过
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return errors.Errorf("接收到操作系统信号: %v", sig)
		}
	})

	// httpServer stop
	group.Go(func() error {
		select {
		case <-ctx.Done():
			log.Printf("errgroup 退出: %+v", ctx.Err())
		case <-shutdown:
			log.Println("接收到 http server 关闭信号")
		}
		return httpServer.Stop(ctx)
	})

	// grpcServer stop
	group.Go(func() error {
		<-ctx.Done()
		log.Printf("errgroup 退出: %+v", ctx.Err())
		return grpcServer.Stop(ctx)
	})

	fmt.Printf("errgroup 退出: %+v\n", group.Wait())
}
