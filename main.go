package main

import (
	// "fmt"
	"net/http"
	"geecache"
	"gee"
	"time"
	"log"
	"fmt"
	// "html/template"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := geecache.NewHttpPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))

	// r:=gee.Default()

	// // r.SetFuncMap(template.FuncMap{
	// // 	"FormatAsDate":FormatAsDate,
	// // })
	// // r.LoadHTMLGlob("templates/*")
	// // r.Static("/assets","./static")

	// // stu1 := &student{Name: "Geektutu", Age: 20}
	// // stu2 := &student{Name: "Jack", Age: 22}

	// r.GET("/", func(c *gee.Context) {
	// 	// c.HTML(http.StatusOK, "css.tmpl",nil)
	// 	c.String(http.StatusOK, "hello,hychen\n")
	// })
	
	// r.GET("/panic",func(c *gee.Context){
	// 	names:=[]string{"hychen"}
	// 	c.String(http.StatusOK, names[2])
	// })
	// // r.GET("/students", func(c *gee.Context) {
	// // 	c.HTML(http.StatusOK, "arr.tmpl", gee.H{
	// // 		"title":  "gee",
	// // 		"stuArr": [2]*student{stu1, stu2},
	// // 	})
	// // })

	// // r.GET("/date", func(c *gee.Context) {
	// // 	c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
	// // 		"title": "gee",
	// // 		"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
	// // 	})
	// // })

	// // v2 := r.Group("/v2")
	// // v2.Use(onlyForV2())// v2 midlleware
	// // {
	// // 	v2.GET("/hello/:name", func(c *gee.Context) {
	// // 		// expect /hello?name=geektutu
	// // 		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	// // 	})
	// // }

	// r.Run(":9999")
}