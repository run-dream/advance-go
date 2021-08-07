## 需求
按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包括对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 wire 进行依赖构建。
可以使用自己熟悉的框架。

## 介绍
这里选择 kratos 来进行处理。 因为之前没有接触过 golang 完整的项目，把这次作业当做 对 [kratos](https://github.com/go-kratos/kratos) 框架的学习。

[文档地址](https://go-kratos.dev/docs/)

### 安装和初始化

#### 安装

```bash
# 开启GO111MODULE
go env -w GO111MODULE=on

# 配置 GOPROXY
go env -w GOPROXY=https://goproxy.cn,direct

# 安装 kratos
go get -u github.com/go-kratos/kratos/cmd/kratos/v2@latest

# 注意 kratos 需要安装 protoc
# 从 https://github.com/protocolbuffers/protobuf/releases 下载解压，配置环境变量

# 官方的protoc编译器中并不支持Go语言，需要安装一个插件才能生成Go代码, 安装protoc-gen-go
go get -u github.com/golang/protobuf/protoc-gen-go

# 新建项目模板
kratos new project

# 生成所有proto源码、wire等等
go generate ./...

# 运行项目
kratos run
```

### 目录结构

```bash
# windows 下执行 tree /f

│  .gitignore   # git 忽略掉二进制文件、临时文件以及开发工具环境
│  Dockerfile	# Docker的配置文件
│  generate.go  
│  go.mod       # go mod
│  go.sum       # go mod
│  LICENSE      # 开源许可
│  Makefile		# linux 下使用 make 命令来构建项目 
│  README.md	# 项目说明
│  
├─api			# 下面维护了微服务使用的proto文件以及根据它们所生成的go文件
│  └─helloworld	# 项目名称
│      └─v1		# 接口版本
│              error_reason.pb.go			# protoc 生成的结构相关的文件
│              error_reason.pb.validate.go	# protoc 生成的message中属性的验证规则
│              error_reason.proto			# proto   文件
│              error_reason.swagger.json	# swagger 文档
│              error_reason_errors.pb.go
│              greeter.pb.go
│              greeter.pb.validate.go
│              greeter.proto
│              greeter.swagger.json
│              greeter_grpc.pb.go
│              greeter_http.pb.go
│              
├─cmd			# 整个项目启动的入口文件
│  └─project
│          main.go		# 入口文件
│          wire.go		
│          wire_gen.go	# wire 生成
│          
├─configs				# 配置文件
│      config.yaml
│      
├─internal				# 业务逻辑		
│  ├─biz				# 业务逻辑的组装层，类似 DDD 的 domain 层，data 类似 DDD 的 repo
│  │      biz.go
│  │      greeter.go
│  │      README.md
│  │      
│  ├─conf				# 内部使用的config的结构定义，使用proto格式生成
│  │      conf.pb.go
│  │      conf.proto
│  │      
│  ├─data				#  业务数据访问，包含 cache、db 等封装，实现了 biz 的 repo 接口
│  │      data.go
│  │      greeter.go
│  │      README.md
│  │      
│  ├─server				# http和grpc实例的创建和配置
│  │      grpc.go
│  │      http.go
│  │      server.go
│  │      
│  └─service			# 实现了 api 定义的服务层，类似 DDD 的 application 层，处理 DTO 到 biz 领域实体的转换(DTO -> DO)，同时协同各类 biz 交互，但是不应处理复杂逻辑
│          greeter.go
│          README.md
│          service.go
│          
└─third_party			# api 依赖的第三方proto
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



### 使用到的包和作用

#### gRPC 相关

- `google.golang.org/grpc`

#### 依赖注入相关

- `github.com/google/wire`

#### 

### 配置管理





### 参考文档

[make](https://www.ruanyifeng.com/blog/2015/02/make.html)

[protoc](https://segmentfault.com/a/1190000020418571)

