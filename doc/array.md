### array

定义: `var a [len]int`

注意:

1. 是**同一种数据类型**的**固定长度**的序列

2. 数组是**值类型**，赋值和传参会复制整个数组，而不是指针。因此改变副本的值，不会改变本身的值。此外

   值拷贝行为会造成性能问题，通常会建议使用 slice，或数组指针。

3. 指针数组 `[n]*T`，数组指针 `*[n]T`

初始化:

1.  一维数组

   ```go
   	a := [3]int{1, 2}           // 未初始化元素值为 0。
       b := [...]int{1, 2, 3, 4}   // 通过初始化值确定数组长度。
       c := [5]int{2: 100, 4: 200} // 使用索引号初始化元素。
    	d := [...]struct {
           name string
           age  uint8
       }{
           {"user1", 10}, // 可省略元素类型。
           {"user2", 20}, // 别忘了最后一行的逗号。
       }
   ```

2. 多维数组

   ```go
   var arr1 [2][3]int = [...][3]int{{1, 2, 3}, {7, 8, 9}} // 只有第一个维度的才能使用 ...
   ```


### slice

定义:  `var slice []type

注意:

1. 切片是数组的一个引用，因此切片是**引用**类型。但自身是结构体，值拷贝传递
2. 切片的长度可以改变，因此，切片是一个**可变**的数组
3.  如果 `slice == nil`，那么 len、cap 结果都等于 0

初始化：

1. 索引

   ```
   nums := [...]int{1,2,3}
   slice := nums[0:2]
   ```

2. make

   ```
   slice := make([]int, 2)
   ```



append:

```go
nums = append(nums, 1)
```

底层数据结构:

```go
type slice struct {
    array unsafe.Pointer
    len   int
    cap   int
}
```



### 指针

定义: `var ptr *type`

注意:

1. 区别于C/C++中的指针，Go语言中的指针不能进行偏移和运算，是安全指针。
2. new用于类型的内存分配，并且内存对应的值为类型零值，返回的是指向类型的指针
3. make只用于slice、map以及channel的初始化，返回的还是这三个引用类型本身
4. 在Go语言中支持对结构体指针直接使用.来访问结构体的成员

初始化:

```
num := 1
ptr := &p
*p = 2

ptrStr := new(string)
ptrSlice ：= make([]int, 10)
	
```



### 切片、指针和slice的比较

1. 把一个大数组传递给函数会消耗很多内存，采用切片的方式传参可以避免上述问题。

   切片是引用传递，所以它们不需要使用额外的内存并且比使用数组更有效率。

2. 传指针会有一个弊端，两个指针地址都是同一个，万一原数组的指针指向更改了，那么函数里面的指针指向都会跟着更改。

3. 并非所有时候都适合用切片代替数组，因为切片底层数组可能会在堆上分配内存，而且小数组在栈上拷贝的消耗也未必比 make 消耗大
