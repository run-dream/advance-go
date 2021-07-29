### 需求
基于 errgroup 实现一个http server 的启动和关闭，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。 

### 分析
要做的事情包括:
1. 启动一个 http 服务
2. 启动一个 goroutine 监听操作系统的终止信号。注意程序无法捕获信号 SIGKILL 和 SIGSTOP，因此 os/signal 包对这两个信号无效。
3. http 服务提供 shutdown 的接口，关闭自身，并通知 goroutine
4. goroutine 收到终止信号时，关闭自身，并通知 http 服务。

可以 使用 两个 channel 来同步信号。

但是题目要求使用 errgroup， 看errgroup 的文档 得知
数据结构
``` go
type Group struct {
  cancel  func()             //context cancel()
    wg      sync.WaitGroup         
    errOnce sync.Once          //只会传递第一个出现错的协程的 error
    err     error              //传递子协程错误
}
```
1. errgroup.Go 会启动一个 goroutine, 当返回第一个非空 error 时，会调用 group 传进来的 cancel 方法， 通知 ctx.Done()
``` go
func (g *Group) Go(f func() error) {
    g.wg.Add(1)

    go func() {
        defer g.wg.Done()
        if err := f(); err != nil {
            g.errOnce.Do(func() {       
                g.err = err            
                
                if g.cancel != nil {
                    g.cancel()
                }
            })
        }
    }()
}
```
也就是说 Go 方法生成的 goroutine 会自动终止。
2. errgroup.Wait 底层通过 wg.wait 来阻塞当前 goroutine。

因为主 goroutine 要调用 errgroup.Wait 阻塞。所以另起一个 goroutine 来终止 httpserver。


