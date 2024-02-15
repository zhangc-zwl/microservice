package microservice

import "strings"

type TreeNode struct {
	isEndNode  bool
	name       string
	routerName string
	children   []*TreeNode
}

func (t *TreeNode) Put(path string) {
	root := t
	strs := strings.Split(path, "/")
	for index, name := range strs {
		if index == 0 {
			continue
		}
		isMatch := false
		for _, node := range root.children {
			if node.name == name {
				isMatch = true
				root = node
				break
			}
		}
		if !isMatch {
			node := &TreeNode{
				name:       name,
				routerName: root.routerName + "/" + name,
				children:   make([]*TreeNode, 0),
			}
			root.children = append(root.children, node)
			root = node
		}

		if index == len(strs)-1 {
			root.isEndNode = true
		}
	}
}

func (t *TreeNode) Get(path string) *TreeNode {
	root := t
	strs := strings.Split(path, "/")
	for index, name := range strs {
		if index == 0 {
			continue
		}
		isMatch := false
		for _, node := range root.children {
			if node.name == name || node.name == "*" || strings.Contains(node.name, ":") {
				if index == len(strs)-1 {
					return node
				}
				isMatch = true
				root = node
				break
			} else if node.name == "**" {
				return node
			}
		}
		if !isMatch {
			return nil
		}
	}
	return nil
}
