### goroutine

#### 进程和线程

操作系统会为应用程序创建一个进程，进程是**分配资源**的基本单位，包括地址空间，文件句柄，设备和线程。

线程是进程调度的一个**执行路径**，用于在处理器执行我们在函数中编写的代码。一个进程从一个线程开始，即主线程，在主线程终止时，进程终止。主线程可以依次启动更多的线程，这些线程也可以启动更多的线程。

无论线程属于哪个进程，操作系统都会安排线程在可用处理器上运行。由操作系统控制调度。

#### goroutine 和 并行

##### 并发和并行

从CUP资源占用角度考虑:

**并发**：单个CPU的操作系统在进行多线程操作时，将CPU运行时间划分成若干个时间段，再将时间段分配给各个线程执行，在一个时间段内只运行一个线程，其他线程处于挂起状态。

**并行**：对于多CPU系统进行多线程操作，一个CPU在执行一个线程时，另一个CPU执行另一个线程。两个线程互不抢CPU资源，可以**同时**进行。

区别：并发是指两个或者多个时间在同一**时间间隔**内发生；并行是指两个或者多个时间在同一**时刻**发生。

##### goroutine

可以认为 goroutine 是由 go runtime 管理的用户态的线程。

##### goroutine 的使用原则

1. **keep yourself busy or do the work yourself** 

   如果你的 goroutine 在从另一个 goroutine 的执行结果前无法继续执行，那么就不该使用这个 goroutine

2. **leave concurrency to the caller**

   让调用者来决定是否并行。 同时提供并行和非并行的接口

   返回 chan 的问题:

   - 无法区分 错误 和 空
   - 必须读取完所有 channel

   通过回调函数

3. **never start a go routine without knowing when it will stop**

   避免 go routine 泄漏。使用的时候首先问自己：

   1. 什么时候会停止
   2. 是否提供了停止的入口

   使用 WaitGroup 避免 incomplete work

### Memory model

如何保证在一个 goroutine 中看到另一个 goroutine 修改的变量的值， 如果程序中修改数据时有其他 goroutine 同时读取，那么必须将读取**串行化**，需要使用 channel 或者其他同步原语如 sync或者 atomic

#### [底层原理](https://bbs.huaweicloud.com/blogs/140746)

- happen before

  - 定义

    是一种**偏序关系**。如果存在 `happen-before(a,b)`，那么操作a及a之前在内存上面所做的操作（如赋值操作等）都对操作b可见，即操作 a 影响了操作 b 。

  - 背景

    cpu的运行极快，而读取主存对于cpu而言有点慢了，在读取主存的过程中cpu一直闲着（也没数据可以运行），这对资源来说造成极大的浪费。所以慢慢的cpu演变成了多级cache结构，cpu在读cache的速度比读内存快了n倍。

    当线程在执行时，会保存临界资源的副本到私有work memory中，这个memory在**CPU cache**中，修改这个临界资源会更新work memory但并不一定立刻刷到**主存**中, 也就无法被其他线程看到，但是其他线程又可能对这个值进行操作，所以需要 happen-before 来指导重排序。

  在一个线程或者goroutine内，按照程序代码顺序，书写在前面的操作 Happen-Before书写在后面的操作。

  由于重排，不同的 goroutine 可能会看到不同的执行顺序。

- memory reordering

  - 定义

    用户代码，要先编译成汇编代码，也就是各种指令，包括读写内存的指令。CPU为了**提升性能**，加入了流水线，分支预测等等。其中，为了提高读写内存的效率，会对指令进行重排序，这就是内存重排。

  - 要求

    对于多线程的程序，CPU提供了 barrier 或者 fence 。要求 barrier 指令要求对所有的内存的操作必须扩散到其他的CPU缓存才能继续执行。

- single machine word 复制是原子的

  指针赋值是原子的。

- copy on write

### Package sync

#### Share Memory By Communicating

传统的线程模型的线程之间通信需要使用共享内存。通过共享数据结构需要使用**锁**保护。在某些情况下，可以使用**线程安全**的数据结构，这样更加高效。

而 *Go* 提供了 并发原语 goroutines 和 channel， **不显示的使用锁**，而是鼓励使用 chan 在 goroutine 之间传递对数据的引用，保证在特定的时间只有一个 goroutine 可以访问数据。

#### Detecting Race Conditions With Go

```bash
go build -race // race detecting
go build -s   // 输出汇编
```

产生原因：

- i++ 不是原子操作
- interface 有两个指针

原则：

没有安全的data race

- 原子性
- 可见性

#### sync.atomic

使用方式:

每次构建一个新的对象，然后atomic.Store,使用 atomic.Load 来访问，适合读很多写极少的情况。

性能测试:

使用 bench mark 来测试。

原理：

Copy On Write.

#### mutex

分类：

- rwlock
- mutex

[实现原理](https://medium.com/a-journey-with-go/go-mutex-and-starvation-3f4f4e75ad50)：

- barging 为了提高吞吐量，当锁被释放的时候，会唤醒第一个等待者，然后把锁给第一个等待者或者第一个请求锁的人
- handsoff 当锁释放的时候，锁会一直持有到第一个等待者准备好获取锁。它降低了吞吐量，因为锁被持有，即使另一个 goroutine 准备获取锁。
- spinning 自旋会一直等待循环获取锁，在等待队列为空或者应用程序重度使用锁时效果不错。

go 1.8 使用了 barging 和 spinning 的结合实现。当试图获取已被持有的锁时，如果本地队列为空，且P的数量大于1，goroutine 会自旋几次，自旋后，go routing 会展厅。

go 1.9 go使用了一个新的饥饿模式来解决先前问题。在释放的时候会触发handsoff。所有等待锁操作1ms的goroutine会被诊断为饥饿。当被标记为饥饿时，unlock方法会handsoff把锁直接扔给第一个等待者。在饥饿模式下，自旋也被停用，因为传入的goroutines没有机会获取到为下一个等待者保留的锁。

#### errgroup

使用场景：

分解为多个小任务并发执行，最终等全部执行完毕。

核心原理：

sync.WaitGroup 管理并行执行的 goroutine

#### sync.Pool

使用场景:

保存和复用临时对象，减少内存分配，降低GC压力

核心原理：

ring buffer + 双向链表

### Package context

#### channel

channel 是一种类型安全的消息队列，充当两个 goroutine 之间的管道，来交互任意资源。

分类：

- 无缓冲

  阻塞的，本质是保证同步。

  receive 先于 send 发生。

  好处： 100%保证可以收到

  代价： 延迟时间未知

- 有缓冲

  有容量。

  send 先于 receive 发生。

  好处：延迟更小

  代价：不保证数据到达，越大的buffer，越小的保证到达。

close:

只有在通知接收方goroutine所有的数据都发送完毕的时候才需要关闭通道

 1.对一个关闭的通道再发送值就会导致panic。
 2.对一个关闭的通道进行接收会一直获取值直到通道为空。
 3.对一个关闭的并且没有值的通道执行接收操作会得到对应类型的零值。
 4.关闭一个已经关闭的通道会导致panic。



#### context

请求级别的上下文。

用途：

- 超时处理
- 取消
- 传递数据

使用：

- 第一个参数 context
- 请求结构的可选context 对象
- 不要在结构体里使用 context， 除非它是 request

实现原理:

- WithValue 新建一个context

  ``` go
  type valueCtx struct {
      Contex
      key, value interface{}
  }
  ```

  基于链表，依次查询

  注意,这些数据必须是安全的，不能篡改。而是用context的方法来构建新的context。

  为了挂载多个字段，可以使用 `map[string]interface{}`

- 级联取消



