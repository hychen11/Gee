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
- 通配`*`。例如 `/static/*filepath`，可以匹配`/static/fav.ico`，也可以匹配`/static/js/jQuery.js`，这种模式常用于静态服务器，能够递归地匹配子路径。

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



## Syntax

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
