package geecache

// import (
// 	"string"
// )

type ByteView struct {
	b []byte	//byte can support multiple data bytes
}

func (v ByteView) Len() int { 
	return len(v.b) 
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c:=make([]byte, len(b))
	copy(c,b)
	return c
}

func (v ByteView) String()	string { 
	return string(v.b)
}