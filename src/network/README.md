1. 总结几种 socket 粘包的解包方式: fix length/delimiter based/length field based frame decoder。尝试举例其应用
- 定长分隔(每个数据包最大为该长度，不足时使用特殊字符填充) ，但是数据不足时会浪费传输资源
- 使用特定字符来分割数据包，但是若数据中含有分割字符则会出现Bug
- 在数据包中添加长度字段，弥补了以上两种思路的不足，推荐使用
参考：https://juejin.cn/post/6844903882108174343
2. 实现一个从 socket connection 中解码出 goim 协议的解码器。
首先查看协议原文: http://www.tony.wiki/development/2016/09/04/goim-protocol.html
开始实现。