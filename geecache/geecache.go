package geecache
import (
	"fmt"
	"log"
	"sync"
	"geecache/singleflight"
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
	peers  	    PeerPicker
	loader      *singleflight.Group
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
		loader:     &singleflight.Group{},

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


func (g *Group)load(key string)(value ByteView,err error){
	viewi,err:=g.loader.Do(key,func()(interface{},error){
		if g.peers!=nil{
			if peer,ok:=g.peers.PickPeer(key);ok{
				if bytes,err:=g.getFromPeer(peer,key);err==nil{
					return bytes,nil
				}
				log.Println("[GeeCache]Failed to get from peer!")
			}
		}
		return g.getLocally(key)
	})
	if err==nil{
		return viewi.(ByteView),nil
	}
	return
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

func (g *Group) RegisterPeers(peers PeerPicker){
	if g.peers!=nil{
		panic("RegisterPeers called more than once")
	}
	g.peers=peers
}

func (g *Group) getFromPeer(peer PeerGetter,key string)(ByteView,error){
	bytes,err:=peer.Get(g.name,key)
	if err!=nil{
		return ByteView{},err
	}
	return ByteView{bytes},nil
}