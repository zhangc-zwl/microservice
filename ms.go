package microservice

import (
	"fmt"
	"log"
	"net/http"
)

const ANY = "ANY"

type HandlerFunc func(ctx *Context)

type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

type router struct {
	routerGroups []*routerGroup
}

type routerGroup struct {
	name              string
	handlerFuncMap    map[string]map[string]HandlerFunc
	middlewareFuncMap map[string]map[string][]MiddlewareFunc
	handlerMethodMap  map[string][]string
	treeNode          *TreeNode
	middlewares       []MiddlewareFunc
}

func (r *router) Group(name string) *routerGroup {
	group := &routerGroup{
		name:              name,
		handlerFuncMap:    make(map[string]map[string]HandlerFunc),
		handlerMethodMap:  make(map[string][]string),
		middlewareFuncMap: make(map[string]map[string][]MiddlewareFunc),
		treeNode:          &TreeNode{name: "/", children: make([]*TreeNode, 0)},
	}
	r.routerGroups = append(r.routerGroups, group)
	return group
}

func (r *routerGroup) Use(middlewareFunc ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewareFunc...)
}

func (r *routerGroup) methodHandle(name string, method string, h HandlerFunc, ctx *Context) {
	//通用中间件
	if r.middlewares != nil {
		for _, middleware := range r.middlewares {
			h = middleware(h)
		}
	}

	//路由中间件
	middlewareFunc := r.middlewareFuncMap[name][method]
	if middlewareFunc != nil {
		for _, middleware := range middlewareFunc {
			h = middleware(h)
		}
	}

	h(ctx)
}

func (r *routerGroup) handle(name, method string, handlerFunc HandlerFunc, middleFunc ...MiddlewareFunc) {
	if _, ok := r.handlerFuncMap[name]; !ok {
		r.handlerFuncMap[name] = make(map[string]HandlerFunc)
		r.middlewareFuncMap[name] = make(map[string][]MiddlewareFunc)
	}
	if _, ok := r.handlerFuncMap[name][method]; ok {
		panic("有重复的路由:" + name)
	}
	r.handlerFuncMap[name][method] = handlerFunc
	r.middlewareFuncMap[name][method] = append(r.middlewareFuncMap[name][method], middleFunc...)

	r.treeNode.Put(name)
}

func (r *routerGroup) Any(name string, handlerFunc HandlerFunc, middleFunc ...MiddlewareFunc) {
	r.handle(name, ANY, handlerFunc, middleFunc...)
}

func (r *routerGroup) Get(name string, handlerFunc HandlerFunc, middleFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodGet, handlerFunc, middleFunc...)
}

func (r *routerGroup) Post(name string, handlerFunc HandlerFunc, middleFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handlerFunc, middleFunc...)
}

func (r *routerGroup) Delete(name string, handlerFunc HandlerFunc, middleFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodDelete, handlerFunc, middleFunc...)
}

func (r *routerGroup) Put(name string, handlerFunc HandlerFunc, middleFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPut, handlerFunc, middleFunc...)
}

func (r *routerGroup) Patch(name string, handlerFunc HandlerFunc, middleFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPatch, handlerFunc, middleFunc...)
}

func (r *routerGroup) Options(name string, handlerFunc HandlerFunc, middleFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodOptions, handlerFunc, middleFunc...)
}

func (r *routerGroup) Head(name string, handlerFunc HandlerFunc, middleFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodHead, handlerFunc, middleFunc...)
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
	e.httpRequestHandler(w, r)
}

func (e *Engine) httpRequestHandler(w http.ResponseWriter, r *http.Request) {
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
				group.methodHandle(node.routerName, ANY, handle, ctx)
				return
			}
			if handle, ok := group.handlerFuncMap[node.routerName][method]; ok {
				group.methodHandle(node.routerName, method, handle, ctx)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "%s %s not allowed \n", r.RequestURI, method)
			return
		}
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
