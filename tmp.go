package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default() //创建gin
	r.GET("/", index) //绑定路由
	r.Run(":8001") //运行绑定端口
}

func index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "go!go!go!",
	})
}
  
