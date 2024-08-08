package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern 	string 	 	//left part of string to be matched
	part 		string		//:lang
	children    []*node
	isWild 		bool 		//is strict match? part has : or * is true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

//return first match child node (used for insert)
func (n *node)matchChild(part string) *node{
	for _,child :=range n.children{
		if child.pattern==part || child.isWild{
			return child
		}
	}
	return nil
}

//match all child nodes!
func (n *node)matchChildren(part string) []*node{
	nodes:=make([]*node,0)
	for _,child :=range n.children{
		child.String()
		if child.part==part || child.isWild{
			nodes=append(nodes,child)
		}
	}
	return nodes
}

//parts is the url formatted!
func (n *node)insert(pattern string,parts []string,height int){
	if len(parts)==height{
		n.pattern=pattern
		return
	}
	part:=parts[height]
	child:=n.matchChild(part)
	if child==nil{
		child=&node{part:part,isWild:part[0]==':'||part[0]=='*'}
		n.children=append(n.children,child)
	}
	child.insert(pattern,parts,height+1)
}

//parts is the url formatted!
func (n *node)search(parts []string,height int)*node{
	if len(parts)==height||strings.HasPrefix(n.part,"*"){
		if n.pattern==""{
			return nil
		}
		return n
	}
	part:=parts[height]
	children:=n.matchChildren(part)
	for _,child:=range children{
		result:=child.search(parts,height+1)
		if result!=nil{
			return result
		}	
	}
	return nil
}


func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}