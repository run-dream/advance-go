### 基本需求
当dao层遇到 sql.NoRows 应该如何处理

### 分析
sql.NoRows 查不到属于正常情况，所以不需要包装，直接返回。其他类型的错误如 connection timeout 等需要 wrap 抛出
1. 获取单个时，直接返回空指针。
2. 获取多个时，直接返回空数组。