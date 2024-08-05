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