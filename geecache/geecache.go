package geecache
import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get (key string) ([]byte,error)
}

type GetterFunc func(key string) ([]byte,error)

func (f GetterFunc) Get(key string) ([]byte,error) {
	return f(key)
}

type Group struct{
	name 		string
	getter      Getter
	maincache   cache
}

//RWMutex include RLock and Lock 
var mu sync.RWMutex
var groups = make(map[string]*Group)

func NewGroup(name string,cacheBytes int64,getter Getter)*Group{
	if getter == nil {
		panic("nil getter")
	}

	mu.Lock()
	defer mu.Unlock()
	
	g:= &Group{
		name: 		name,
		getter: 	getter,
		maincache:  cache{cacheBytes:cacheBytes},
	}
	groups[name]=g
	return g
}

func GetGroup(name string) *Group{
	mu.RLock()
	defer mu.RUnlock()
	g:=groups[name]
	return g
}

func (g *Group)Get(key string)(ByteView,error){
	if key==""{
		return ByteView{},fmt.Errorf("key nil!")
	}
	if v,ok:=g.maincache.get(key);ok{
		log.Println("hit")
		return v,nil
	}
	return g.load(key)
}


func (g *Group)load(key string)(ByteView,error){
	return g.getLocally(key)
}

func (g *Group)getLocally(key string)(ByteView,error){
	bytes,err:=g.getter.Get(key)
	if err!=nil{
		return ByteView{},err
	}
	value :=ByteView{b:cloneBytes(bytes)} 
	g.maincache.add(key,value)
	return value,nil
}
