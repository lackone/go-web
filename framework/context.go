package framework

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Context struct {
	res          http.ResponseWriter
	req          *http.Request
	resLock      *sync.RWMutex       //控制res的锁
	isTimeout    bool                //是否超时
	handlers     []ControllerHandler //当前请求的handler链条
	handlerIndex int                 //当前链条在哪个节点
	params       map[string]string   //uri参数
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		res:          w,
		req:          r,
		resLock:      &sync.RWMutex{},
		isTimeout:    false,
		handlers:     []ControllerHandler{},
		handlerIndex: -1,
		params:       map[string]string{},
	}
}

func (this *Context) WriterMux() *sync.RWMutex {
	return this.resLock
}

func (this *Context) GetRequest() *http.Request {
	return this.req
}

func (this *Context) GetResponse() http.ResponseWriter {
	return this.res
}

func (this *Context) SetHandlers(handlers []ControllerHandler) {
	this.handlers = handlers
}

func (this *Context) SetIsTimeout() {
	this.isTimeout = true
}

func (this *Context) GetIsTimeout() bool {
	return this.isTimeout
}

func (this *Context) BaseContext() context.Context {
	return this.req.Context()
}

func (this *Context) Next() error {
	this.handlerIndex++
	if this.handlerIndex < len(this.handlers) {
		if err := this.handlers[this.handlerIndex](this); err != nil {
			return err
		}
	}
	return nil
}

func (this *Context) SetParams(params map[string]string) {
	this.params = params
}

// 开始实现context.Context接口
func (this *Context) Deadline() (deadline time.Time, ok bool) {
	return this.BaseContext().Deadline()
}

func (this *Context) Done() <-chan struct{} {
	return this.BaseContext().Done()
}

func (this *Context) Err() error {
	return this.BaseContext().Err()
}

func (this *Context) Value(key any) any {
	return this.BaseContext().Value(key)
}

// 结束实现context.Context接口
