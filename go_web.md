```shell
(base) ➜  go tree 
.
├── gee
│   ├── context.go
│   ├── gee.go
│   ├── go.mod
│   └── router.go
├── go.mod
├── go_web.md
├── main.go
└── test.go
```

这里的gee里的文件都是用一个包里的

```go
package gee
```

只有首字母大写的标识符才是公开的（可导出的），可以被其他包导入并使用；而首字母小写的标识符是私有的（未导出的），只能在同一个包内部访问！！！

```shell
go build .
#auto compile whole package
go run .
#auto compile and run main function!

```

### go.mod

```shell
go mod init gee
```

### %s & %q

`%s` 用于格式化字符串类型的值。它会按原样输出字符串的内容。

`%q` 用于格式化字符串类型的值，但会将字符串用双引号括起来，同时会对特殊字符进行转义

```go
str:="Hello, World!\n"
Using %s: Hello, World!
Using %q: "Hello, World!\n"
```

### Syntax

```go
s1 := []int{1, 2}
s2 := []int{3, 4}
s1 = append(s1, s2...)
//...is expand element in s2
//for loop in go
for _,x := range s1{
    //...
}
```

```go
func A(numbers ...int){
    a:=[]int{}
    a=append(a,numbers...)
}
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
    group.middlewares = append(group.middlewares, middlewares...)
}

```

 `...int` 表示这个函数可以接受任意数量的 `int` 参数

### Data Structure

`container/heap`

```go
package main

import (
    "container/heap"
    "fmt"
)

// 定义一个整数堆，继承了heap.Interface接口
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func main() {
    h := &IntHeap{2, 1, 5}
    heap.Init(h)               // 初始化堆
    heap.Push(h, 3)            // 插入一个元素
    fmt.Printf("最小元素: %d\n", (*h)[0]) // 查看最小元素

    for h.Len() > 0 {
        fmt.Printf("%d ", heap.Pop(h)) // 依次取出最小元素
    }
}

```

`container/list`

https://pkg.go.dev/container/list

```go

import (
	"container/list"
	"fmt"
)

func main() {
	// Create a new list and put some numbers in it.
	l := list.New()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

}
```

`container/ring`

```go
package main

import (
    "container/ring"
    "fmt"
)

func main() {
    r := ring.New(3)  // 创建一个大小为3的循环链表

    // 给每个元素赋值
    for i := 0; i < r.Len(); i++ {
        r.Value = i
        r = r.Next()
    }

    // 打印循环链表中的所有元素
    r.Do(func(p interface{}) {
        fmt.Println(p.(int))
    })

    // 移动指针并再次打印
    r = r.Move(1)
    fmt.Println("After moving 1 step:")
    r.Do(func(p interface{}) {
        fmt.Println(p.(int))
    })
}
```

### Test

#### Run Test

```shell
$go test
$go test -v				#显示每个用例的测试结果
$go test -cover			#输出每个被测试的函数的覆盖率信息
$go test -run TestAdd -v
```

name:  `XX_test.go` 在待测试文件同一目录下

#### Benchmark:  `*testing.B` 

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = add(1, 2)
    }
}

func add(a, b int) int {
    return a + b
}
```

#### TestMain:  `*testing.M`

```go
// TestMain 是测试的主函数
func TestMain(m *testing.M) {
    // 在测试运行前执行的初始化代码
    fmt.Println("Setting up tests")

    // 运行所有的测试
    exitCode := m.Run()

    // 在测试运行后执行的清理代码
    fmt.Println("Cleaning up")

    // 根据测试结果退出
    os.Exit(exitCode)
}

func TestAdd(t *testing.T) {
    result := add(1, 2)
    if result != 3 {
        t.Errorf("Expected 3 but got %d", result)
    }
}
```

#### `testing.T` 

```go
func TestXXX(t *testing.T){
	result := Add(1, 2)
    expected := 3
    if result != expected {
        t.Errorf("Expected %d but got %d", expected, result)
    }
    
    if result > 10 {
        t.Fatal("Result exceeds 10")
    }
}
```

`t.Error` 和 `t.Fatal`: 用于报告测试错误，`t.Error` 继续执行测试，`t.Fatal` 停止测试执行

#### Subtests

```go
func TestMul(t *testing.T) {
	t.Run("pos", func(t *testing.T) {
		if Mul(2, 3) != 6 {
			t.Fatal("fail")
		}

	})
	t.Run("neg", func(t *testing.T) {
		if Mul(2, -3) != -6 {
			t.Fatal("fail")
		}
	})
}
```

```shell
go test -run TestMul/pos
go test -run TestMul/neg
go test -run TestMul
```



# Gee

## HTTP

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

// Engine is the uni handler for all requests
type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}
```

**启动服务器**：`http.ListenAndServe(":9999", engine)`启动一个HTTP服务器，监听本地的9999端口。这个函数需要两个参数：监听的地址和处理所有请求的处理器。在这里，`engine`实例作为处理器，因为它实现了`http.Handler`接口。

`ServeHTTP`相当是结构体Engine的方法

```go
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
```

当你在`http.ListenAndServe`函数调用中传入`nil`作为第二个参数时，这意味着你希望使用Go标准库中的默认路由器（默认的多路复用器，`http.DefaultServeMux`）。`http.HandleFunc`函数实际上是在这个默认的多路复用器上注册路径和处理函数（handlers）的快捷方式。所以`nil`和通过`http.HandleFunc`注册的路径及其处理函数是结合在一起工作的。

### net/http

```go
package http

type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}

func ListenAndServe(address string, h Handler) error
```

第二个参数的类型是什么呢？通过查看`net/http`的源码可以发现，`Handler`是一个接口，需要实现方法 ServeHTTP ，也就是说，只要传入任何实现了 ServerHTTP 接口的实例，所有的HTTP请求，就都交给了该实例处理了。

**实现了路由映射表，提供了用户注册静态路由的方法，包装了启动服务的函数**

我们对一个结构体Engine里建立一个router map里面对应的是method和handler func

```go
// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}
```

然后对Engine里添加方法，GET，POST，AddRoute，Run，

- 首先定义了类型`HandlerFunc`，这是提供给框架用户的，用来定义路由映射的处理方法。我们在`Engine`中，添加了一张路由映射表`router`，key 由请求方法和静态路由地址构成，例如`GET-/`、`GET-/hello`、`POST-/hello`，这样针对相同的路由，如果请求方法不同,可以映射不同的处理方法(Handler)，value 是用户映射的处理方法。

- 当用户调用`(*Engine).GET()`方法时，会将路由和处理方法注册到映射表 router 中，`(*Engine).Run()`方法，是 ListenAndServe 的包装。

  - ```go
    // Run defines the method to start a http server
    func (engine *Engine) Run(addr string) (err error) {
    	return http.ListenAndServe(addr, engine)
    }
    ```

- `Engine`实现的 ServeHTTP 方法的作用就是，解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND 。

### http.Request

`http.Request`是Go语言`net/http`包中定义的一个结构体，代表了一个HTTP请求。这个结构体包含了一个HTTP请求的所有数据，比如请求行（方法、URL）、请求头（Headers）、请求体（Body）、查询参数等。`http.Request`还提供了一些方法来方便地操作这些数据和处理HTTP请求。以下是一些常用的`http.Request`方法和属性的介绍：

#### 常用属性

- **Method**：字符串，表示HTTP请求的方法（如"GET", "POST", "PUT", "DELETE"等）。
- **URL**：`*url.URL`类型，表示请求的URL。它包含了路径、查询字符串等组成部分。
- **Header**：`http.Header`类型，表示请求的头部。它是一个映射，映射的键是头部字段名，值是字段值的切片。
- **Body**：`io.ReadCloser`接口，表示请求的主体。对于"GET"请求，这通常是`nil`；对于"POST"或"PUT"等请求，这里包含了发送的数据。读取Body后，你需要关闭它。

#### 常用方法

- **FormValue(key string) string**：返回URL的查询参数或POST表单数据中键为`key`的第一个值。如果键不存在，则返回空字符串。这个方法在内部调用`ParseForm`方法，因此不需要事先调用`ParseForm`。
- **PostFormValue(key string) string**：类似于`FormValue`，但仅返回POST请求体中的表单参数值。
- **Context() context.Context**：返回请求的`context.Context`。这个上下文在请求处理过程中携带跨API边界和goroutines的请求范围的值、取消信号、截止日期等。
- **WithContext(ctx context.Context) \*Request**：返回一个新的`Request`指针，它的上下文被改为`ctx`。原始请求的副本被修改；原始请求体不会被改变。
- **ParseForm() error**：解析URL中的查询字符串和请求体中的表单数据（仅适用于"POST", "PUT"和"PATCH"方法，且内容类型为`application/x-www-form-urlencoded`）。解析后的数据可以通过`Form`和`PostForm`字段访问。
- **ParseMultipartForm(maxMemory int64) error**：解析多部分表单数据，`maxMemory`参数指定了系统使用的最大内存，超过这个值的文件内容会被写入临时文件中。解析后的文件数据可以通过`MultipartForm`字段访问。

## Context

- `Handler`的参数变成成了`gee.Context`，提供了查询Query/PostForm参数的功能。
- `gee.Context`封装了`HTML/String/JSON`函数，能够快速构造HTTP响应。

`curl "http://localhost:9999/login" -X POST -d 'username=geektutu&password=1234'`

`-X POST` : method is POST

`-d 'username=geektutu&password=1234'` : data

相当于创建一个router，然后一个request url 在router里对应一个handlerFunc

然后context管理request和response的格式

## Dynamic Route

map can only store static router

动态路由实现

方法1：开源的路由实现`gorouter`支持在路由规则中嵌入正则表达式，例如`/p/[0-9A-Za-z]+`即路径中的参数仅匹配数字和字母

2：另一个开源实现`httprouter`就不支持正则表达式

动态路由具备以下两个功能

`:/lang`是占位符

- 参数匹配`:`。例如 `/p/:lang/doc`，可以匹配 `/p/c/doc` 和 `/p/go/doc`。
- 通配`*`。例如 `/static/*filepath`，可以匹配`/static/fav.ico`，也可以匹配`/static/js/jQuery.js`，这种模式常用于静态服务器，能够递归地匹配子路径。(`filepath`匹配`fav.ico`)

路由规则`/assets/*filepath`，可以匹配`/assets/`开头的所有的地址。例如`/assets/js/geektutu.js`，匹配后，参数`filepath`就赋值为`js/geektutu.js`

#### Parse

例如`/p/go/doc`匹配到`/p/:lang/doc`，解析结果为：`{lang: "go"}`，`/static/css/geektutu.css`匹配到`/static/*filepath`，解析结果为`{filepath: "css/geektutu.css"}`。

```shell
#gee has router_test.go
go test gee # will auto run test file
```

## Route Group Control

- 以`/post`开头的路由匿名可访问。
- 以`/admin`开头的路由需要鉴权。
- 以`/api`开头的路由是 RESTful 接口，可以对接第三方平台，需要三方平台鉴权。

路由分组不仅支持简单的分组，还可以嵌套分组

`/post` 是一个分组，它包含了 `/post/a` 和 `/post/b` 这样的子分组

`/admin` 可能包含了 `/admin/user` 和 `/admin/settings` 这样的子分组

为了实现分组功能，框架需要定义一个 `Group` 对象。这个对象通常包含以下几个属性：

1. **前缀（prefix）**：
   - 表示这个分组的 URL 前缀，比如 `/`、`/api` 或者 `/admin`。所有在这个分组下定义的路由，都会自动加上这个前缀。
2. **父分组（parent）**：
   - 如果支持分组嵌套，那么每个分组都可能有一个父分组。比如 `/admin/user` 分组的父分组就是 `/admin`。通过父分组，可以构建出分组的层次结构。
3. **中间件（middlewares）**：
   - 存储应用在该分组上的中间件。这些中间件会在处理分组下的所有路由之前执行。
4. **引擎（Engine）**：
   - 这是整个框架的核心，管理所有的路由和中间件。`Group` 对象需要访问 `Engine`，以便在分组中定义路由时能够正确地将它们注册到框架中。

```go
r := gee.New()       // 创建一个新的框架实例
v1 := r.Group("/v1") // 创建一个以 "/v1" 为前缀的分组

// 在分组 v1 下定义一个 GET 请求的路由
v1.GET("/", func(c *gee.Context) {
    c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
})
```

v1下面的router都包含/v1为前缀

## Middlewares

框架需要有一个插口，允许用户自己定义功能，嵌入到框架中，仿佛这个功能是框架原生支持的一样

Gee 的中间件的定义与路由映射的 Handler 一致，处理的输入是`Context`对象。插入点是框架接收到请求初始化`Context`对象后，允许用户使用自己定义的中间件做一些额外的处理，例如记录日志等，以及对`Context`进行二次加工

`c.Next()`表示等待执行其他的中间件或用户的`Handler`	

当接收到请求后，匹配路由，该请求的所有信息都保存在`Context`中。中间件也不例外，接收到请求后，应查找所有应作用于该路由的中间件，保存在`Context`中，依次进行调用。

```go
type Context{
    //other 
    ...
	//middleware
	handlers []HandlerFunc	//array list to record HandlerFunc
    index 	int				//index of array list
}
```

为什么依次调用后，还需要在`Context`中保存呢？因为在设计中，中间件不仅作用在**处理流程前**，也可以作用在**处理流程后**，即在用户定义的 Handler 处理完毕后，还可以执行剩下的操作。

`c.Next()` 表示进入下一个中间件，等下一个中间件完成后返回结果

```go
func A(c *Context){
    part1
    c.Next()
    part2
}

func B(c *Context){
    part3
    c.Next()
    part4
}
```

`func A -> part1 -> part3 ->HandlerFunc -> part2 -> part4`

具体实现过程就是

一个RouterGroup里有自己的middleware array list

```go
func (group *RouterGroup)Use(middleware ...HandlerFunc){
    group.middleware=append(group.middleware,middleware)
}
```

这里有请求时，对应的url从group里取出相应的HandlerFunc，存入对应的router的middleware list里，然后调用`c.Next()`进行处理:`c.index`遍历调用

注意！每个middleware里只能调用一个Next() 然后下面的Context Next每次调用index++，只用一个handlers

```go
func (c *Context) Next(){
	c.index++
	s:=len(c.handlers)
	for ;c.index < s; c.index++{
		c.handlers[c.index](c)
	}
}
```



## HTML Template

#### 静态文件(Serve Static Files)

要做到服务端渲染，第一步便是要支持 JS、CSS 等静态文件

如果将所有的静态文件放在`/usr/web`目录下，那么`filepath`的值即是该目录下文件的相对地址。映射到真实的文件后，将文件返回，静态服务器就实现了

gee 框架要做的，仅仅是解析请求的地址，映射到服务器上文件的真实地址，交给`http.FileServer`处理就好了

在 `Context` 中添加了成员变量 `engine *Engine`，这样就能够通过 Context 访问 Engine 中的 HTML 模板

## Panic Recover

#### 主动触发panic

```go
	panic("crash")
```

#### defer

panic 会导致程序被中止，但是在退出前，会先处理完当前协程上已经defer 的任务，执行完成后再退出。效果类似于 java 语言的 `try...catch`。

出现异常在退出前先把defer的内容运行了

```go
func main(){
	defer func(){
		fmt.Println("defer")
	}()
	panic("crash")
}
//defer
//crash
```

当 `main` 函数开始执行时，首先注册了一个 `defer` 函数，该函数会在 `main` 函数结束时打印 `"defer"`。

接着，程序遇到了 `panic("crash")`，这会引发一个 `panic`，导致程序的正常执行流程被中断。

在 `panic` 被引发后，Go 会立刻开始执行所有已注册的 `defer` 语句。因此，`fmt.Println("defer")` 被执行，打印出 `"defer"`。

最后，程序由于 `panic` 而终止，并输出 `panic: crash`，表明程序崩溃并提供了 `panic` 的信息。

#### recover

```go
func test_recover() {
	defer func() {
		fmt.Println("defer func")
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
	fmt.Println("after panic")
}

func main() {
	test_recover()
	fmt.Println("after recover")
}
```

recover 捕获了 panic，程序正常结束。*test_recover()* 中的 *after panic* 没有打印，这是正确的，当 panic 被触发时，控制权就被交给了 defer

就像在 java 中，`try`代码块中发生了异常，控制权交给了 `catch`，接下来执行 catch 代码块中的代码

main() 中打印了 after recover，说明程序已经恢复正常，继续往下执行直到结束

`runtime.Callers(3, pcs[:])` Caller返回调用栈， `skip first 3 caller`

第 0 个 Caller 是 Callers 本身，第 1 个是上一层 trace，第 2 个是再上一层的 `defer func`

`runtime.FuncForPC(pc)` 获取对应的函数，在通过 `fn.FileLine(pc)` 获取到调用该函数的文件名和行号

```go
func trace(message string)string{
	var pcs [32]uintptr
	n:=runtime.Callers(3,pcs[:])// skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	
	for _,pc :=range pcs[:n] {
		fn:=runtime.FuncForPC(pc)
		file,line:=fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d",file,line))
	}
	return str.String()
}

func Recovery() HandlerFunc{
	return func(c *Context){
		defer func(){
			if err:=recover();err!=nil{
				message:=fmt.Sprintf("%s",err)
				log.Printf("%s \n\n",trace(message))
				c.Fail(http.StatusInternalServerError,"Internal Server Error")
			}
		}()
		c.Next()
	}
}
```

# GeeCache

访问一个网页，网页和引用的 JS/CSS 等静态文件，根据不同的策略，会缓存在浏览器本地或是 CDN 服务器，那在第二次访问的时候，就会觉得网页加载的速度快了不少

微博的点赞的数量，不可能每个人每次访问，都从数据库中查找所有点赞的记录再统计，数据库的操作是很耗时的，很难支持那么大的流量，所以一般点赞这类数据是缓存在 Redis 服务集群中的

#### scale horizontally

利用多台计算机的资源，并行处理提高性能就要缓存应用能够支持分布式，这称为水平扩展(scale horizontally)

#### scale vertically

增加单个节点的计算、存储、带宽等，来提高系统的性能，硬件的成本和性能并非呈线性关系，大部分情况下，分布式系统是一个更优的选择

## LRU

Linkedlist+Hashmap

```go
type Cache struct {
	maxBytes 	int64 // max capacity
	nBytes 		int64 // number current bytes	
	ll 			*list.List
	cache 		map[string]*list.Element
	OnEvicted	func(key string,value Value)
}

type entry struct {
	key 	string
	value   Value
}

//interface ( count how many bytes it takes )
type Value interface{
	Len()	int
}
```

### 类型断言语法

处理 `interface{}` 类型时非常有用的工具

```go
x.(T)
//x interface
//T type

var x interface{} = "Hello, Go!"

// 类型断言为 string
s := x.(string)
fmt.Println(s)  // 输出: Hello, Go!
```

如果类型断言失败（例如 `x` 的实际类型不是 `string`），程序会发生 `panic`

```go
kv:=ele.Value(*entry)
//Value is interface
//*entry is type
//judge if ele is *entry type, if true, convert to *entry type
```

## sync.Mutex 

`sync.Mutex`互斥锁，如果goroutine获取锁的占有权，别的goroutine请求会阻塞在`Lock()`，直到调用`Unlock()`

````go
func A(){
    m.Lock()
    defer m.Unlock()
    //other operation
}
````

在 `add` 方法中，判断了 `c.lru` 是否为 nil，如果等于 nil 再创建实例。这种方法称之为延迟初始化(**Lazy Initialization**)，一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。主要用于提高性能，并减少程序内存要求。

## Group

```
                           是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
```

回调函数就是`Getter`

## HTTP Server

```go
package http

type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}
```

`HTTPPool` 存URL.PATH

- `HTTPPool` 只有 2 个参数，一个是 self，用来记录自己的地址，包括主机名/IP 和端口。
- 另一个是 basePath，作为节点间通讯地址的前缀，默认是 `/_geecache/`，那么 http://example.com/_geecache/ 开头的请求，就用于节点间的访问。因为一个主机上还可能承载其他的服务，加一段 Path 是一个好习惯。比如，大部分网站的 API 接口，一般以 `/api` 作为前缀。

先从url parse出groupName和key

然后GetGroup后Get(key)

## Consistent Hashing

一致性哈希算法将 key 映射到 2^32 的空间中，将这个数字首尾相连，形成一个环

![一致性哈希添加节点 consistent hashing add peer](https://geektutu.com/post/geecache-day4/add_peer.jpg)

`key11`，`key2`，`key27` 均映射到 peer2，`key23` 映射到 peer4。此时，如果新增节点/机器 peer8，假设它新增位置如图所示，那么只有 `key27` 从 peer2 调整到 peer8，其余的映射均没有发生改变。

也就是说，一致性哈希算法，在新增/删除节点时，只需要重新定位该节点附近的一小部分数据，而不需要重新定位所有的节点，这就解决了上述的问题。

#### 虚拟节点

一个真实节点对应多个虚拟节点

假设 1 个真实节点对应 3 个虚拟节点，那么 peer1 对应的虚拟节点是 peer1-1、 peer1-2、 peer1-3

为了均衡负载，我们引入虚拟节点的概念。假设每个物理节点有 3 个虚拟节点：

- `peer1` 对应 `peer1-1`、`peer1-2`、`peer1-3`
- `peer2` 对应 `peer2-1`、`peer2-2`、`peer2-3`
- `peer3` 对应 `peer3-1`、`peer3-2`、`peer3-3`

通过虚拟节点，数据在不同的物理节点之间分布得更加均匀，避免了数据倾斜的问题。虚拟节点扩充了环上的节点数量，增加了映射的可能性，从而使得负载更加均衡。

```go
type Hash func(data []byte) uint32

type Map struct{
	hash 		Hash
	replicas	int
	keys        []int
	hashMap     map[int]string
}

func New(replicas int,fn Hash)*Map{
	m:= &Map(
		hash:fn,
		replicas:replicas,
		keys:make([]int),
		hashMap:make(map[int]string),
	)

	if m.hash==nil{
		m.hash=crc32.ChecksumIEEE
	}
	return m
}
```

`replicas` 是虚拟节点的倍数

`Map`存所有的hash keys

```go
import (
	"fmt"
	"hash/crc32"
)

func main() {
	data := []byte("hello")
	checksum := crc32.ChecksumIEEE(data)
	fmt.Printf("CRC-32 Checksum: %08x\n", checksum)
}
```

```go
import "sort"
numbers:=[]int{1,3,2,5,4}
sort.Ints(numbers)
fmt.Println("Sorted numbers:", numbers)


package main
import (
	"fmt"
	"sort"
)
type Person struct {
	Name string
	Age  int
}
type People []Person
func (p People) Len() int {
	return len(p)
}
func (p People) Less(i, j int) bool {
	return p[i].Age < p[j].Age
}
func (p People) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func main() {
	people := People{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	sort.Sort(people)
	for _, person := range people {
		fmt.Printf("%s: %d\n", person.Name, person.Age)
	}
}
```

这里Lambda自定义就是通过`Less(i,j int)`实现的

`sort.Sort(people)` 

#### Binary Search (with Lambda)

```go
idx := sort.Search(len(m.keys), func(i int) bool {
    return m.keys[i] >= hash
})
```

假设 `m.keys` 是一个已排序的切片，`hash` 是你要查找的位置。这个代码的作用是找到 `m.keys` 中第一个大于或等于 `hash` 的元素的位置 `idx`。

- **`len(m.keys)`**：指定了搜索的范围，即从 `0` 到 `len(m.keys)-1`。
- **`func(i int) bool { return m.keys[i] >= hash }`**：这是一个匿名函数，用来判断 `m.keys[i]` 是否大于或等于 `hash`。在 `sort.Search` 的过程中，这个函数会被多次调用，通过不断缩小搜索范围来找到符合条件的第一个索引 `i`。

```go
import (
	"fmt"
	"strconv"
)

func main() {
	// 字符串转整数
	num, err := strconv.Atoi("123")
	if err == nil {
		fmt.Println("Integer:", num)
	}

	// 整数转字符串
	str := strconv.Itoa(456)
	fmt.Println("String:", str)
}
```

这里的`Add`和`Get` 就是增加虚拟节点，查询最近的虚拟节点

## Register Peers

在分布式缓存系统中，缓存数据通常不会存储在一个单一的服务器上，而是分布在多个节点上。每个节点都是一个独立的服务器，保存着整个缓存系统的一部分数据。**远程节点**就是指除了当前正在处理请求的节点以外的其他节点。

```go
type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)
```

这里的`url.QueryEscape()`是 Go 语言中 `net/url` 包提供的一个函数，用于对字符串进行 URL 编码。具体来说，它会将字符串中的特殊字符转换为百分号编码格式，以便在 URL 中安全传输。

`http.Get(u)`：这个函数是 Go 语言标准库中的一个函数，用于发送一个 HTTP GET 请求。`u` 是目标 URL（通常是一个字符串），程序会尝试连接到这个 URL 并获取响应。

这里的`res` 是一个 `*http.Response` 类型的指针，表示 HTTP 请求的响应。它包含了 HTTP 请求返回的所有信息，包括状态码、响应头、响应体等

```go
type Response struct {
    Status     string // 服务器返回的状态行，例如 "200 OK"
    StatusCode int    // 状态码，例如 200, 404, 500 等
    Proto      string // HTTP协议版本，例如 "HTTP/1.1"
    ProtoMajor int    // 主版本号，例如 1
    ProtoMinor int    // 次版本号，例如 1
    Header     Header // 响应头，一个 map[string][]string，保存了键值对形式的响应头信息
    Body       io.ReadCloser // 响应体，用于读取服务器返回的数据
    ContentLength int64 // 响应体的长度，值为 -1 时表示未知长度
    TransferEncoding []string // 传输编码，例如 "chunked"
    Close      bool // 是否在响应后关闭连接
    Uncompressed bool // 是否已解压缩响应体
    Trailer    Header // 可能的 Trailer 响应头，类似于 Header
    Request    *Request // 产生该响应的请求
    TLS        *tls.ConnectionState // 如果使用了 TLS，保存相关的 TLS 信息
}
```

使用 `io.ReadAll` 从 `res.Body` 读取响应体内容，然后将其转换为字符串或其他所需格式。需要注意的是，读取后 `res.Body` 会被耗尽，无法再次读取

## singleflight

> **缓存雪崩**：缓存在同一时刻全部失效，造成瞬时DB请求量大、压力骤增，引起雪崩。缓存雪崩通常因为缓存服务器宕机、缓存的 key 设置了相同的过期时间等引起。

> **缓存击穿**：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到 DB ，造成瞬时DB请求量大、压力骤增。

> **缓存穿透**：查询一个不存在的数据，因为不存在则不会写到缓存中，所以每次都会去请求 DB，如果瞬间流量过大，穿透到 DB，导致宕机。

并发协程之间不需要消息传递，非常适合 `sync.WaitGroup`。

- wg.Add(1) 锁加1
- wg.Wait() 阻塞，直到锁被释放
- wg.Done() 锁减1

我们并发了 N 个请求 `?key=Tom`，8003 节点向 8001 同时发起了 N 次请求。假设对数据库的访问没有做任何限制的，很可能向数据库也发起 N 次请求，容易导致缓存击穿和穿透。即使对数据库做了防护，HTTP 请求是非常耗费资源的操作，针对相同的 key，8003 节点向 8001 发起三次请求也是没有必要的。那这种情况下，我们如何做到只向远端节点发起一次请求呢

```go
package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call
}
```

* `call` 代表正在进行中，或已经结束的请求。使用 `sync.WaitGroup` 锁避免重入。
* `Group` 是 singleflight 的主数据结构，管理不同 key 的请求(call)

## protobuf

