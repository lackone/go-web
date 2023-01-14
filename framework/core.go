package framework

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Core struct {
	Router map[string]*Tree `json:"router"`
}

func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()

	return &Core{Router: router}
}

func (this *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//封装自定义context
	ctx := NewContext(w, r)

	json, err := json.Marshal(this.Router)
	fmt.Println(err)
	fmt.Println(string(json))

	handler := this.FindRouteHandler(r)
	if handler == nil {
		ctx.Json(404, "not found")
		return
	}

	if err := handler(ctx); err != nil {
		ctx.Json(500, "server error")
		return
	}
}

func (this *Core) Get(url string, handler ControllerHandler) {
	if err := this.Router["GET"].AddRouter(url, handler); err != nil {
		log.Fatalln(err)
	}
}

func (this *Core) Post(url string, handler ControllerHandler) {
	if err := this.Router["POST"].AddRouter(url, handler); err != nil {
		log.Fatalln(err)
	}
}

func (this *Core) Put(url string, handler ControllerHandler) {
	if err := this.Router["PUT"].AddRouter(url, handler); err != nil {
		log.Fatalln(err)
	}
}

func (this *Core) Delete(url string, handler ControllerHandler) {
	if err := this.Router["DELETE"].AddRouter(url, handler); err != nil {
		log.Fatalln(err)
	}
}

func (this *Core) Group(prefix string) IGroup {
	return NewGroup(this, prefix)
}

func (this *Core) FindRouteHandler(req *http.Request) ControllerHandler {
	path := strings.ToUpper(req.URL.Path)
	method := strings.ToUpper(req.Method)

	if m, ok := this.Router[method]; ok {
		return m.FindHandler(path)
	}
	return nil
}
