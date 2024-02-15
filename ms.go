package microservice

import (
	"fmt"
	"log"
	"net/http"
)

const ANY = "ANY"

type HandlerFunc func(ctx *Context)

type router struct {
	routerGroups []*routerGroup
}

type routerGroup struct {
	name             string
	handlerFuncMap   map[string]map[string]HandlerFunc
	handlerMethodMap map[string][]string
	treeNode         *TreeNode
}

func (r *router) Group(name string) *routerGroup {
	group := &routerGroup{
		name:             name,
		handlerFuncMap:   make(map[string]map[string]HandlerFunc),
		handlerMethodMap: make(map[string][]string),
		treeNode:         &TreeNode{name: "/", children: make([]*TreeNode, 0)},
	}
	r.routerGroups = append(r.routerGroups, group)
	return group
}

func (r *routerGroup) handle(name, method string, handlerFunc HandlerFunc) {
	if _, ok := r.handlerFuncMap[name]; !ok {
		r.handlerFuncMap[name] = make(map[string]HandlerFunc)
	}
	if _, ok := r.handlerFuncMap[name][method]; ok {
		panic("有重复的路由:" + name)
	}
	r.handlerFuncMap[name][method] = handlerFunc

	r.treeNode.Put(name)
}

func (r *routerGroup) Any(name string, handlerFunc HandlerFunc) {
	r.handle(name, ANY, handlerFunc)
}

func (r *routerGroup) Get(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodGet, handlerFunc)
}

func (r *routerGroup) Post(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPost, handlerFunc)
}

func (r *routerGroup) Delete(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodDelete, handlerFunc)
}

func (r *routerGroup) Put(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPut, handlerFunc)
}

func (r *routerGroup) Patch(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPatch, handlerFunc)
}

func (r *routerGroup) Options(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodOptions, handlerFunc)
}

func (r *routerGroup) Head(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodHead, handlerFunc)
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router{},
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	for _, group := range e.routerGroups {
		routerName := SubStringLast(r.RequestURI, "/"+group.name)
		node := group.treeNode.Get(routerName)
		if node != nil && node.isEndNode {
			ctx := &Context{
				w,
				r,
			}
			if handle, ok := group.handlerFuncMap[node.routerName][ANY]; ok {
				handle(ctx)
				return
			}
			if handle, ok := group.handlerFuncMap[node.routerName][method]; ok {
				handle(ctx)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "%s %s not allowed \n", r.RequestURI, method)
			return
		}
		//for name, methodHandle := range group.handlerFuncMap {
		//	url := "/" + group.name + name
		//	if r.RequestURI == url {
		//	}
		//}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s not allowed \n", r.RequestURI)
}

func (e *Engine) Run() {
	//for _, group := range e.routerGroups {
	//	for key, value := range group.handlerFuncMap {
	//		http.HandleFunc("/"+group.name+key, value)
	//	}
	//}
	http.Handle("/", e)
	if err := http.ListenAndServe(":8111", nil); err != nil {
		log.Fatal(err)
	}
}
