package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// newServer 初始化一个http 服务
func newServer() (http.Server, <-chan struct{}) {
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
	return http.Server{
		Handler: mux,
		Addr:    ":8001",
	}, shutdown
}

func main() {
	// 用来管理 goroutine
	group, ctx := errgroup.WithContext(context.Background())

	// 启动 http 服务
	server, shutdown := newServer()
	group.Go(func() error {
		return server.ListenAndServe()
	})

	// 监听操作系统信号
	group.Go(func() error {
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		select {
		// 这里通过
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return errors.Errorf("接收到操作系统信号: %v", sig)
		}
	})

	// 处理shutdown 信号 和 errgroup Done 的信号
	group.Go(func() error {
		select {
		case <-ctx.Done():
			log.Printf("errgroup 退出: %+v", ctx.Err())
		case <-shutdown:
			log.Println("接收到 http server 关闭信号")
		}
		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		log.Println("正在关闭 http 服务")
		return server.Shutdown(timeoutCtx)
	})

	fmt.Printf("errgroup 退出: %+v\n", group.Wait())

}
