## 评论系统架构设计



### 功能模块

不要做需求的翻译机，先理解业务背后的**本质**，事情的初衷。

- 发布评论
- 读取评论
- 删除评论

- 管理评论 作者置顶、后台运营管理(搜索、删除、审核等)。

**在动手设计前，反复思考，真正编码的时间只有5%**

### 架构设计

- BFF: comment

  复杂评论业务的服务编排，比如访问账号服务进行等级判定，同时需要在 BFF 面向移动端/WEB场景来设计 API，这一层抽象把评论的本身的内容列表处理(加载、分页、排序等)进行了隔离，关注在业务平台化逻辑上。

- Service: comment-service

  服务层，去平台业务的逻辑，专注在评论功能的 API 实现上，比如发布、读取、删除等，关注在稳定性、可用性上，这样让上游可以灵活组织逻辑把基础能力和业务能力剥离。

- Job: comment-job

  消息队列的最大用途是消峰处理

- Admin: comment-admin

    管理平台，按照安全等级划分服务，尤其划分运营平台，他们会共享服务层的存储层(MySQL、Redis)。运营体系的数据大量都是检索，我们使用 canal 进行同步到 ES 中，整个数据的展示都是通过 ES，再通过业务主键更新业务数据层，这样运营端的查询压力就下方给了独立的 fulltext search 系统。

- Dependency: account-service、filter-service

  整个评论服务依赖的一些外部 gRPC 服务，统一的平台业务逻辑在 comment BFF 层收敛，这里 account-service 主要是账号服务，filter-service 是敏感词过滤服务。

``` 
http/gRPC -> BFF: comment -> otherService
			/
		  Service: comment-service
		  /		  |	 	 \
		Kafka	Redis   MySql -> Cannal -> ElasticSearch <- comment-admin
		 \		  |		/
		    Job: comment-job    
```

- comment-service，专注在评论数据处理


 **架构设计等同于数据设计，梳理清楚数据的走向和逻辑。尽量避免环形依赖、数据双向请求等。**



#### 从cache-aside到 write behind

读的核心逻辑: 

Cache-Aside 模式，先读取缓存，再读取存储。

早期 cache rebuild 是做到服务里的，对于重建逻辑，一般会使用 read ahead 的思路，即预读，用户访问了第一页，很有可能访问第二页，所以缓存会超前加载，避免频繁 cache miss。当缓存抖动时候，特别容易引起集群 thundering herd 现象，大量的请求会触发 cache rebuild，因为使用了**预加载**，容易导致服务 OOM。

所以我们看到回源的逻辑里，我们使用了消息队列来进行逻辑异步化，对于当前请求只返回 mysql 中部分数据即止。

写的核心逻辑: 
写和读相比较，写可以认为是**透穿到存储层**的，系统的瓶颈往往就来自于存储层，或者有状态层。对于写的设计上，我们认为刚发布的评论有极短的延迟(通常小于几 ms)对用户可见是可接受的，把**对存储的直接冲击下放到消息队列**，按照消息反压的思路，即如果存储 latency 升高，消费能力就下降，自然消息容易堆积，系统始终以最大化方式消费。



Kafka 是存在 partition 概念的，可以认为是物理上的一个小队列，一个 topic 是由一组 partition 组成的，所以 Kafka 的吞吐模型理解为: **全局并行，局部串行**的生产消费方式。对于入队的消息，可以按照 hash(comment_subject) % N(partitions) 的方式进行分发。那么某个 partition 中的 评论主题的数据一定都在一起，这样方便我们串行消费。
mysql binlog 中的数据被 **canal 中间件流式消费**，获取到业务的原始 CRUD 操作，需要回放录入到 es 中，但是 es 中的数据最终是面向运营体系提供服务能力，需要检索的数据维度比较多，在入 es 前需要做一个异构的 joiner，把单表变宽预处理好 join 逻辑，然后倒入到 es 中。
一般来说，运营后台的检索条件都是组合的，使用 es 的好处是避免依赖 mysql 来做多条件组合检索，同时 mysql 毕竟是 oltp 面向线上联机事务处理的。通过冗余数据的方式，使用其他引擎来实现。
es 一般会存储检索、展示、primary key 等数据，当我们操作编辑的时候，找到记录的 primary key，最后交由 comment-admin 进行运营测的 CRUD 操作。



comment 作为 BFF，是面向端，面向平台，面向业务组合的服务。所以平台扩展的能力，我们都在 comment 服务来实现，方便统一和准入平台，以统一的接口形式提供平台化的能力。



### 存储设计

#### 数据库设计

- 数据写入: 事务更新 comment_subject，comment_index，comment_content 三张表，其中 content 属于**非强制需要一致性考虑**的。可以先写入 content，之后事务更新其他表。即便 content 先成功，后续失败仅仅存在一条 ghost 数据。
- 数据读取: 基于 obj_id + obj_type 在 comment_index 表找到评论列表，WHERE root = 0 ORDER BY floor。之后根据 comment_index 的 id 字段捞出 comment_content 的评论内容。对于二级的子楼层，WHERE parent/root IN (id...)。
  因为产品形态上只存在二级列表，因此只需要迭代查询两次即可。对于嵌套层次多的，产品上，可以通过二次点击支持。

comment_index: 评论楼层的索引组织表，实际并不包含内容。

comment_content: 评论内容的表，包含评论的具体内容。其中 comment_index 的 id 字段和comment_content 是1对1的关系，这里面包含几种设计思想。

- 表都有**主键**，即 cluster index，是物理组织形式存放的，comment_content 没有 id，是为了减少一次 二级索引查找，直接基于主键检索，同时 comment_id 在写入要尽可能的顺序自增。
- 索引、内容分离，方便 mysql datapage 缓存更多的 row，如果和 content 耦合，会导致更大的 IO。长远来看 content 信息可以直接使用 KV storage 存储。

#### 缓存设计

- comment_subject_cache: 对应主题的缓存，value 使用 protobuf 序列化的方式存入。我们早期使用 memcache 来进行缓存，因为 redis 早期单线程模型，吞吐能力不高。
- comment_index_cache: 使用 **redis sortedset** 进行索引的缓存，索引即数据的组织顺序，而非数据内容。参考过百度的贴吧，他们使用自己研发的拉链存储来组织索引，我认为 mysql 作为主力存储，利用 redis 来做加速完全足够，因为 cache miss 的构建，我们前面讲过使用 kafka 的消费者中处理，预加载少量数据，通过增量加载的方式逐渐**预热**填充缓存，而 redis sortedset skiplist 的实现，可以做到 O(logN) + O(M) 的时间复杂度，效率很高。
  sorted set 是要增量追加的，因此必须判定 key 存在，才能 zadd。
  comment_content_cache: 对应评论内容数据，使用 protobuf 序列化的方式存入。类似的我们早期使用 memcache 进行缓存。
  增量加载 + lazy 加载

### 可用性设计

#### singleflight

对于热门的主题，如果存在缓存穿透的情况，会导致大量的同进程、跨进程的数据回源到存储层，可能会引起存储过载的情况，如何只交给同进程内，一个人去做加载存储?
使用归并回源的思路:
https://pkg.go.dev/golang.org/x/sync/singleflight
同进程只交给一个人去获取 mysql 数据，然后批量返回。同时这个 lease owner 投递一个 kafka 消息，做 index cache 的 recovery 操作。这样可以大大减少 mysql 的压力，以及大量透穿导致的密集写 kafka 的问题。
更进一步的，后续连续的请求，仍然可能会短时 cache miss，我们可以在进程内设置一个 short-lived flag，标记最近有一个人投递了 cache rebuild 的消息，直接 drop。

#### 热点

流量热点是因为突然热门的主题，被高频次的访问，因为底层的 cache 设计，一般是按照主题 key 进行一致性 hash 来进行分片，但是热点 key 一定命中某一个节点，这时候 remote cache 可能会变为瓶颈，因此做 cache 的**升级 local cache** 是有必要的，我们一般使用单进程自适应发现热点的思路，附加一个短时的 ttl local cache，可以在进程内吞掉大量的读请求。
在内存中使用 hashmap 统计每个 key 的访问频次，这里可以使用**滑动窗口统计**，即每个窗口中，维护一个 hashmap，之后统计所有未过去的 bucket，汇总所有 key 的数据。
之后使用小堆计算 TopK 的数据，**自动进行热点识别**。