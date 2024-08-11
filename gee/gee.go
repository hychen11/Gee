package gee

import (
	"log"
	"net/http"
	"strings"
)
type HandlerFunc func(*Context)

type( 
	RouterGroup struct{
		prefix 		string
		middleware 	[]HandlerFunc
		parent     	*RouterGroup
		engine 		*Engine				// all groups share a Engine instance
	}
	Engine struct{
		*RouterGroup
		router *router
		groups []*RouterGroup
	}
)

func New() *Engine {
	engine:=&Engine{router: newRouter()}
	engine.RouterGroup=&RouterGroup{engine:engine}
	engine.groups= []*RouterGroup{engine.RouterGroup}

	return engine
}

func (group *RouterGroup)Group(prefix string) *RouterGroup {
	engine:=group.engine
	newGroup:=&RouterGroup{
		prefix:group.prefix+prefix,
		parent:group,
		engine:engine,
	}
	engine.groups=append(engine.groups,newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern:=group.prefix+comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method,pattern,handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}


func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method,pattern,handler)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	middlewares:=[]HandlerFunc{}
	for _,group:=range engine.groups{
		if strings.HasPrefix(req.URL.Path,group.prefix){
			middlewares=append(middlewares,group.middleware...)
		}
	}

	c:=newContext(w,req)
	c.handlers=middlewares
	engine.router.handle(c)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc){
	group.middleware=append(group.middleware,middlewares...)
}