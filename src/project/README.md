## 需求

按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包括对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 wire 进行依赖构建。
可以使用自己熟悉的框架。

## 介绍

这里选择 kratos 来进行处理。 因为之前没有接触过 golang 完整的项目，把这次作业当做 对 [kratos](https://github.com/go-kratos/kratos) 框架的学习。

[文档地址](https://go-kratos.dev/docs/)

## 安装和初始化

```bash
# 安装 go, mac brew 的方式
brew install go

# 开启GO111MODULE
go env -w GO111MODULE=on

# 配置 GOPROXY
go env -w GOPROXY=https://goproxy.cn,direct

# 安装 kratos
go get -u github.com/go-kratos/kratos/cmd/kratos/v2@latest

# 注意 kratos 需要安装 protoc
# 从 https://github.com/protocolbuffers/protobuf/releases 下载解压，配置环境变量

# 安装protoc-gen-go
# 官方的protoc编译器中并不支持Go语言，需要安装一个插件才能生成Go代码, 安装protoc-gen-go
go get -u github.com/golang/protobuf/protoc-gen-go

# 新建项目模板
kratos new project

# 生成所有proto源码、wire等等
go generate ./...

# 运行项目
kratos run
```

## 目录结构

```bash
# windows 下执行 tree /f 或者 mac 下执行 tree

│  .gitignore 			# git 忽略掉二进制文件、临时文件以及开发工具环境
│  Dockerfile				# Docker 的配置文件
│  generate.go 			# go generate, 生成客户端 Api
│  go.mod						# go mod
│  go.sum						# go mod
│  LICENSE					# 开源许可
│  Makefile					# 支持可以通过 make 命令来构建 
│  README.md				# 项目说明
│  
├─api								# 接口，维护了微服务使用的 proto 文件以及根据它们所生成的go文件
│  └─helloworld			
│      └─v1					# 接口版本
│              error_reason.pb.go							
│              error_reason.pb.validate.go		
│              error_reason.proto							
│              error_reason.swagger.json
│              error_reason_errors.pb.go
│              greeter.pb.go									# protoc 根据 proto 文件生成相关的定义
│              greeter.pb.validate.go					# protoc 根据 proto 文件生成相关的校验代码
│              greeter.proto									# proto 文件
│              greeter.swagger.json						# swagger 文档定义
│              greeter_grpc.pb.go
│              greeter_http.pb.go
│              
├─cmd							  		# 整个项目启动的入口文件
│  └─project						# 项目名
│          main.go			# 入口文件
│          wire.go			# wire 定义
│          wire_gen.go	# wire 生成的文件
│          
├─configs								# 配置文件
│      config.yaml
│      
├─internal							# 业务逻辑，不对外暴露
│  ├─biz								# 业务逻辑的组装层，类似 DDD 的 domain 层。 暴露出的结构体是 *XXXUsecase*
│  │      biz.go				# 使用 wire 来整合所有的 biz。 
│  │      greeter.go		# 具体的 greeter 的业务组装
│  │      README.md
│  │      
│  ├─conf								# 内部使用的config的结构定义，使用 proto格式生成
│  │      conf.pb.go
│  │      conf.proto
│  │      
│  ├─data								# 业务数据访问，包含 cache、db 等封装，实现了 biz 的 repo 接口。暴露出的结构体是 *Data* 和 *XXXRepo*
│  │      data.go
│  │      greeter.go
│  │      README.md
│  │      
│  ├─server							# http和grpc实例的创建和配置
│  │      grpc.go
│  │      http.go
│  │      server.go
│  │      
│  └─service						# 实现了 api 定义的服务层，类似 DDD 的 application 层，处理 DTO 到 biz 领域实体的转换(DTO -> DO)，同时协同各类 biz 交互，但是不应处理复杂逻辑。暴露出的结构体是 *XXXService*
│          greeter.go
│          README.md
│          service.go
│          
└─third_party						#  api 依赖的第三方 proto
    │  README.md
    │  
    ├─errors
    │      errors.proto
    │      
    ├─google
    │  └─api
    │          annotations.proto
    │          http.proto
    │          httpbody.proto
    │          
    ├─protoc-gen-openapiv2
    │  └─options
    │          annotations.proto
    │          openapiv2.proto
    │          
    └─validate
            README.md
            validate.proto
```



## 构建工具

- 生成 Proto 代码

  ```bash
  kratos proto client api/helloworld/demo.proto
  ```

- 生成 service 代码

  ``` bash
  kratos proto server api/helloworld/demo.proto -t internal/service
  ```

- 生成 swagger 文件

   ``` bash
  make swagger
   ```

  

## 提供的插件

- 服务注册和服务发现

  - consul
  - etcd
  - kubernetes
  - nacos

- 日志

  - fluent

- 配置中心

  - consul

  - etcd
  - kubernetes
  - nacos

- 监控警告

  - datadog
  - prometheus
  - sentry

- API 文档

  - swagger-api



### 使用到的包和作用

- 依赖注入

  `github.com/google/wire`

- protobuf 

  `google.golang.org/protobuf`



### kratos 提供了什么功能

Kratos 的定位是 一个优雅的 Go 微服务的工具包。

我们看看 go-kratos/kratos 到底提供了些什么？

- `cmd`

  - kratos

    基于 cobra 提供 CLI 给用户使用，包括

    - CmdNew
    - CmdProto
    - CmdUpgrade
    - CmdChange
    - CmdRun

    internal 包括来这些命令的实现。

    - base 基础的功能的实现，如
      - go install 
      - go mod
      - repo  git 管理模版的基础layout
    - project 通过 repo 去初始化仓库
    - change 查看 kratos 的变更日志
    - proto 包括了proto 相关的一些处理
      - CmdAdd 新增 proto 文件
      - CmdClient  生成客户端 go 代码
      - CmdServer 生成服务端 go 代码 
    - run  允许 kratos 项目
    - upgrade 升级 kratos 工具

  - protocol-gen-go-errors

    通过 proto 预定义定义错误码，然后通过 proto-gen-go 生成帮助代码，直接返回 error

  - Protocol-gen-go-http

    处理 http api，生成 http.Handler，然后可以注册到 HTTPServer 中。

-  `config`

  配置源可以指定多个，并且 config 会进行合并成 key/value

- `encoding`

  序列化，支持

  - json
  - protobuf
  - xml
  - yaml

- `errors`

  错误处理。

- `log`

  只提供了日志接口，具体怎么实现需要自己指定。

- `metadata`

  元信息传递。

- `metrics`

  监控。

- `middleware`

  中间件。

  - logging: 用于请求日志的记录。
  - metrics: 用于启用metric。
  - recovery: 用于recovery panic。
  - tracing: 用于启用trace。
  - validate: 用于处理参数校验。
  - metadata: 用于启用元信息传递

- `registry`

  服务注册和发现。

- `third_party`

  第三方的一些 proto 文件

- `transport`

  传输层, 提供了以下的实现

  - grpc
  - http



## 参考资料

[kratos 官方文档](https://go-kratos.dev/docs/)

[make](http://www.ruanyifeng.com/blog/2015/02/make.html)

