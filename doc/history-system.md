## 功能模块
- 变更功能

  添加记录、删除记录、清空历史

- 读取功能

  按照 timeline 返回 top N，点查获取进度信息。

- 其他记录

  暂停/恢复记录，首次观察增加经验等

历史记录类型的业务，是一个**极高 tps 写入，高 qps 读取**的业务服务。分析清楚系统的 **hot path**，投入优化，而不是哪哪都去优化。

## 架构设计

``` 
Client  		APP
   \  			 |
BFF: app-interface/hitory -- Other Service: archive,  commic
	|
  history-service
   /	 | 		\
kafka	Redis	HBase
  \		|		/
  	 history-job
```

- BFF: app-interface

## 存储设计

## 可用性设计

