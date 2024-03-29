package main

import (
	"net/http"

	"gee"
)
/*
使用New()创建 gee 的实例，使用 GET()方法添加路由，
最后使用Run()启动Web服务。这里的路由，只是静态路由，
不支持/hello/:name这样的动态路由，动态路由我们将在下一次实现
*/
func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}