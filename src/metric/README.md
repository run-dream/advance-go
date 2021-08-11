## 描述
参考 hystrix 实现一个滑动窗口

## 思路

### 需求分析
把需求进行分解：
1. 滑动窗口实现
2. 定义指标
3. 更新指标

一步步的来看
1. 滑动窗口实现，需要确定
   - 窗口大小
   - 如何存储数据
   - 如何更新窗口

2. 定义指标
   - CPU 使用率 可以通过 
   - Inflight 当前服务中正在进行的请求的数量
   - Pass&RT 采样窗口内成功请求的数量 单个采样窗口中平均响应时间

3. 更新指标
   不同的指标，触发的方式不一样。
   - CPU 可以通过新建一个 go routine 来定时采样
   - Inflight 需要在每个请求开始和结束的时候手动更新
   - Pass&RT 也是在每个请求的生命周期内进行更新的

这样的话，对于每个请求都要进行处理，比较麻烦。


### hystrix 实现
我们来看看, hystrix 是怎么做的, 找到 [hystrix-go](https://github.com/afex/hystrix-go) 的实现。
git clone 以后, tree /F, 大致过一遍代码
```
├─hystrix
│  │  circuit.go                        // CircuitBreaker 熔断器
│  │  circuit_test.go					
│  │  doc.go
│  │  eventstream.go					// StreamHandler 将每个请求的 metrics 发给 http 仪表盘
│  │  eventstream_test.go
│  │  hystrix.go						// API入口，提供 Go, Do 方法给用户调用
│  │  hystrix_test.go
│  │  logger.go							// logger 日志的抽象， 默认未实现
│  │  metrics.go						// metricExchange 系统的执行情况
│  │  metrics_test.go
│  │  pool.go							// executorPool command 对应的执行池
│  │  pool_metrics.go					// poolMetrics 执行池的指标
│  │  pool_test.go
│  │  settings.go						// Settings command 的配置
│  │  settings_test.go
│  │  
│  ├─metric_collector					// Registry 收集器
│  │      default_metric_collector.go
│  │      metric_collector.go
│  │      
│  └─rolling							// 滑动窗口
│          rolling.go
│          rolling_test.go
│          rolling_timing.go
│          rolling_timing_test.go
│          
├─loadtest
│  │  README.md
│  │  
│  └─service
│          main.go
│          
├─plugins
│      datadog_collector.go
│      graphite_aggregator.go
│      statsd_collector.go
│      statsd_collector_test.go
│      
└─scripts
        vagrant.sh
```



看懵了，我们来看看别的文档,Hystrix是如何实现的:

- 命令模式 command

  将所有请求外部系统（或者叫依赖服务）的逻辑封装到 Go 和 Run 方法里

- 隔离策略 

  为每个 command 生成一个  pool 

- 观察者模式

  - Hystrix通过观察者模式对服务进行状态监听
  - 每个任务都包含有一个对应的Metrics，所有Metrics都由 Registry 来进行维护，
  - 在任务的不同阶段会往Metrics中写入不同的信息，Metrics会对统计到的历史信息进行统计汇总，供熔断器以及Dashboard监控时使用


### 模仿练习(Ctrl C + Ctrl V)

按照题目要求，我们参考 hystrix  把 metric 相关的代码进行处理,只实现了
1. rolling 滑动窗口
   Number 数字类型的滑动窗口，可以进行均值的统计
   Timing 时间类型的滑动窗口，一个 bucket 里可能有多个时长，可以额外提供 P分位的统计

2. collector 数据收集
   Registry 注册采样方式
   MetricResult 统计指标
   MetricCollector 统计方式接口
   RequestCollector 请求相关的收集器
   还可以扩展 ProfileCollector 但是 hystrix 的指标里不包含 CPU 的部分

3. metric 收集注册和变更
   MetricExchange 某个 topic 的 exchange
   Event          触发更新的事件

最后在 main 函数里，提供 /profile 来查看 exchange 的实时信息, /test 用来模拟实际请求。
   

## 参考资料
[hystrix-go 源码分析](https://cloud.tencent.com/developer/article/1454740)

[一文彻底读懂 hystrix-go 源码](https://learnku.com/articles/53090)

[服务容错与保护方案 — Hystrix](https://kiswo.com/article/1030)