package lru

import (
	"container/list"
)

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

func (c *Cache) Len() int { 
	return c.ll.Len()
}

func New(maxBytes int64, onEvicted func(string,Value))*Cache{
	return &Cache{
		maxBytes: 	maxBytes,
		nBytes:   	0,
		ll: 		list.New(), 	//func New() *List
		cache: 		make(map[string]*list.Element),
		OnEvicted: 	onEvicted,
	}
}

//front is most recently used entry!
func (c *Cache) Get(key string)(value Value,ok bool){
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv:=ele.Value.(*entry)
		return kv.value,true
	}
	return
}

func (c *Cache) RemoveOldest(){
	ele:=c.ll.Back()
	if ele!=nil{
		c.ll.Remove(ele)
		
		//convert ele to *entry struct
		kv:=ele.Value.(*entry)

		delete(c.cache,kv.key)
		c.nBytes-=int64(len(kv.key))+int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}	
}

func (c *Cache) Add(key string,value Value){
	//exists already
	if ele,ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv:=ele.Value.(*entry)
		c.nBytes-=int64(kv.value.Len())
		c.nBytes+=int64(value.Len())
		kv.value=value
	}else{
		c.ll.PushFront(&entry{key:key,value:value})
		c.cache[key]=c.ll.Front()
		c.nBytes+=int64(len(key))+int64(value.Len())
	}
	for c.maxBytes!=0 && c.maxBytes<c.nBytes{
		c.RemoveOldest()
	}
}
