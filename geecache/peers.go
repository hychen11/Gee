package geecache
//locate peers that has key!
type PeerPicker interface {
	PickPeer(key string)(peer PeerGetter,ok bool)
}

type PeerGetter interface {
	Get(group string, key string) ([]byte,error)
}

