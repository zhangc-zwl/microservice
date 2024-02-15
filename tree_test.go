package microservice

import (
	"fmt"
	"testing"
)

func TestTreeNode(t *testing.T) {
	root := &TreeNode{
		name:     "/",
		children: make([]*TreeNode, 0),
	}
	root.Put("/user/get/:id")
	root.Put("/user/create/hello")
	root.Put("/user/create/aaa")
	root.Put("/order/get/:id")
	root.Put("/order/create/hello")
	root.Put("/admin/**")

	node := root.Get("/user/get/111")
	fmt.Println(node)
	node = root.Get("/order/get/111")
	fmt.Println(node)
	node = root.Get("/user/create")
	fmt.Println(node)
	node = root.Get("/order/create/hello")
	fmt.Println(node)
	node = root.Get("/admin/create/hello")
	fmt.Println(node)
	//fmt.Println(root.children[0].children[1])
}
