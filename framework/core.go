package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router      map[string]*Tree    `json:"router"`
	middlewares []ControllerHandler `json:"-"` //中间件处理函数
}

func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()

	return &Core{router: router}
}

func (this *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//封装自定义context
	ctx := NewContext(w, r)

	//导找路由
	node := this.FindRouteNode(r)
	if node == nil {
		ctx.SetStatus(404).Json("not found")
		return
	}

	ctx.SetHandlers(node.handlers)

	//解析参数
	params := node.ParseParamsFromEndNode(r.URL.Path)
	ctx.SetParams(params)

	if err := ctx.Next(); err != nil {
		ctx.SetStatus(500).Json("server error")
		return
	}
}

func (this *Core) Use(middlewares ...ControllerHandler) {
	this.middlewares = append(this.middlewares, middlewares...)
}

func (this *Core) Get(url string, handlers ...ControllerHandler) {
	allHandlers := append(this.middlewares, handlers...)
	if err := this.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatalln(err)
	}
}

func (this *Core) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(this.middlewares, handlers...)
	if err := this.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatalln(err)
	}
}

func (this *Core) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(this.middlewares, handlers...)
	if err := this.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatalln(err)
	}
}

func (this *Core) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(this.middlewares, handlers...)
	if err := this.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatalln(err)
	}
}

func (this *Core) Group(prefix string) IGroup {
	return NewGroup(this, prefix)
}

func (this *Core) FindRouteNode(req *http.Request) *Node {
	path := strings.ToUpper(req.URL.Path)
	method := strings.ToUpper(req.Method)

	if m, ok := this.router[method]; ok {
		return m.FindNode(path)
	}
	return nil
}
