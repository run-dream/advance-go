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

# 安装protoc-gen-go
# go get -u github.com/golang/protobuf/protoc-gen-go

# 新建项目模板
kratos new project
```

### 目录结构

```bash
# windows 下执行 tree /f

│  .gitignore # git 忽略掉二进制文件、临时文件以及开发工具环境
│  Dockerfile	# Docker的配置文件
│  generate.go # 
│  go.mod
│  go.sum
│  LICENSE
│  Makefile
│  README.md
│  tree.txt
│  
├─api
│  └─helloworld
│      └─v1
│              error_reason.pb.go
│              error_reason.pb.validate.go
│              error_reason.proto
│              error_reason.swagger.json
│              error_reason_errors.pb.go
│              greeter.pb.go
│              greeter.pb.validate.go
│              greeter.proto
│              greeter.swagger.json
│              greeter_grpc.pb.go
│              greeter_http.pb.go
│              
├─cmd
│  └─project
│          main.go
│          wire.go
│          wire_gen.go
│          
├─configs
│      config.yaml
│      
├─internal
│  ├─biz
│  │      biz.go
│  │      greeter.go
│  │      README.md
│  │      
│  ├─conf
│  │      conf.pb.go
│  │      conf.proto
│  │      
│  ├─data
│  │      data.go
│  │      greeter.go
│  │      README.md
│  │      
│  ├─server
│  │      grpc.go
│  │      http.go
│  │      server.go
│  │      
│  └─service
│          greeter.go
│          README.md
│          service.go
│          
└─third_party
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


### 配置管理
