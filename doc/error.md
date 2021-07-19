### Error vs Exception

- go 标准库 errors

  ```go
  package errors
  
  // New returns an error that formats as the given text.
  // Each call to New returns a distinct error value even if the text is identical.
  func New(text string) error {
  	return &errorString{text}
  }
  
  // errorString is a trivial implementation of error.
  type errorString struct {
  	s string
  }
  
  func (e *errorString) Error() string {
  	return e.s
  }
  ```

  注意两点：

  1. error 返回的是指针，目的是让每次都返回新的错误 
  2. go 的 == 符号 对于结构体来说 如果所有的字段都是可比较的，就会比较是否所有的字段是否相等

- 错误/异常的表示方式

  1. C  int  单参数返回
  2. C++ exception 无法知道方法是否会抛出异常
  3. Java checked exception 需要声明 throws Exception 无法区分异常和灾难 其实也有runtime error
  4. go 多参数返回 error / panic + recover 缺点 大量的 if err != nil 优点可以区分良性和恶性异常

  注意 panic 只用来处理 fatal 错误。



### Error Type

- sentinel error 预定义的错误 表示不可能进行进一步处理的做法 如 io.EOF,

  缺点:

  - 不应该检查 error.Error 的输出 因为会被 error.Errorf 破坏
  - 成为API公共部分
  - 在两个包里创建了依赖 import loop

  结论： 尽可能避免 sentinel errors。

- error types 自己实现 error 接口，并添加额外的信息，

  ```go
  // PathError records an error and the operation and file path that caused it.
  type PathError struct {
  	Op   string
  	Path string
  	Err  error
  }
  
  func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }
  ```

  优点:

  - 可以包装底层错误以提供更多上下文

  缺点：

  调用者需要使用类型断言和类型switch，需要将自定义的error变成public。导致和调用者产生强耦合。

  结论： 尽量避免Error types称为Api的一部分

- opaque errors 不透明错误

  ```go
  // An Error represents a network error.
  type Error interface {
  	error
  	Timeout() bool   // Is the error a timeout?
  	Temporary() bool // Is the error temporary?
  }
  // DNSError represents a DNS lookup error.
  type DNSError struct {
  	Err         string // description of the error
  	Name        string // name looked for
  	Server      string // server used
  	IsTimeout   bool   // if true, timed out; not all timeouts set this
  	IsTemporary bool   // if true, error is temporary; not all errors set this
  	IsNotFound  bool   // if true, host could not be found
  }
  ```

  只返回错误，不返回具体错误信息

  缺点:

  无法提供更多的信息，只能处理二分错误

  优化：

  用方法来暴露错误相关的行为，而不是特定的值或类型



### Handing Error

- indented flow is for errors

  避免代码缩进

- eliminate error handing by eliminating errors

  直接返回原始错误

  使用scanner

  ```
  // Err returns the first non-EOF error that was encountered by the Scanner.
  func (s *Scanner) Err() error {
  	if s.err == io.EOF {
  		return nil
  	}
  	return s.err
  }
  ```


- wrap errors 类似于解决错误堆栈的问题，避免多次处理

  - go 内置方法 不是很强大
  - 第三方包 https://github.com/pkg/errors

  ```go
  type withMessage struct {
  	cause error
  	msg   string
  }
  type withStack struct {
  	error
  	*stack
  }
  // Unwrap provides compatibility for Go 1.13 error chains.
  func (w *withStack) Unwrap() error { return w.error }
  ```

  - 使用errors.New 或者 errors.Errorf 返回错误
  - 调用其他包的函数，通过简单的直接返回
  - 使用其他库进行协作的时候，考虑errors.Wrap
  - 直接返回错误
  - 在应用底层使用 errors.Cause 获取 root error，在和 sentinel error判断
  - 基础库的错误不应该使用 errors.Wrap	

- Error Inspection
  - handle check

