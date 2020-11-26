package webpack

import (
	"fmt"
	"strings"
)

/*
 Example:
 		     /
 	     /       \
      /:lang     about
   /    |    \
/intro /doc /tutorial
*/

type node struct {
	pattern  string  // 待匹配路由,例如 /:lang/doc
	part     string  // 路由的一部分，例如 :lang
	children []*node //子节点列表，例如 :lang 的子节点为 [intro, doc, tutorial]
	isWlid   bool    //是否采用模糊匹配，例如匹配串中含有 * 或者 : 时为 true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWlid=%t}", n.pattern, n.part, n.isWlid)
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		// the pattern only store in the leaf node.
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWlid: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
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

// return the first matched node for inserting
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWlid {
			return child
		}
	}
	return nil
}

// return all the matched nodes for searching
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWlid {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
