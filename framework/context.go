package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	res       http.ResponseWriter
	req       *http.Request
	resLock   *sync.RWMutex //控制res的锁
	isTimeout bool          //是否超时
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		res:       w,
		req:       r,
		resLock:   &sync.RWMutex{},
		isTimeout: false,
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

func (this *Context) SetHandler(handler ControllerHandler) {

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

func (this *Context) QueryInt(key string, def int) int {
	all := this.QueryAll()
	if val, ok := all[key]; ok {
		//注意val类型是[]string
		len := len(val)
		if len > 0 {
			//取val最后一个字符串进行转换
			v, err := strconv.Atoi(val[len-1])
			if err != nil {
				return def
			}
			return v
		}
	}
	return def
}

func (this *Context) QueryString(key string, def string) string {
	all := this.QueryAll()
	if val, ok := all[key]; ok {
		len := len(val)
		if len > 0 {
			return val[len-1]
		}
	}
	return def
}

func (this *Context) QueryArray(key string, def []string) []string {
	all := this.QueryAll()
	if val, ok := all[key]; ok {
		return val
	}
	return def
}

func (this *Context) QueryAll() map[string][]string {
	return this.req.URL.Query()
}

func (this *Context) FormInt(key string, def int) int {
	all := this.FormAll()
	if val, ok := all[key]; ok {
		//注意val类型是[]string
		len := len(val)
		if len > 0 {
			//取val最后一个字符串进行转换
			v, err := strconv.Atoi(val[len-1])
			if err != nil {
				return def
			}
			return v
		}
	}
	return def
}

func (this *Context) FormString(key string, def string) string {
	all := this.FormAll()
	if val, ok := all[key]; ok {
		len := len(val)
		if len > 0 {
			return val[len-1]
		}
	}
	return def
}

func (this *Context) FormArray(key string, def []string) []string {
	all := this.FormAll()
	if val, ok := all[key]; ok {
		return val
	}
	return def
}

func (this *Context) FormAll() map[string][]string {
	this.req.ParseForm()
	return this.req.PostForm
}

func (this *Context) BindJson(data any) error {
	all, err := io.ReadAll(this.req.Body)
	if err != nil {
		return err
	}
	//body只能读一次，读出来后需要重置下body
	this.req.Body = io.NopCloser(bytes.NewBuffer(all))
	err = json.Unmarshal(all, data)
	if err != nil {
		return err
	}
	return nil
}

func (this *Context) Json(statusCode int, data any) error {
	if this.GetIsTimeout() {
		return nil
	}
	this.res.WriteHeader(statusCode)
	this.res.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = this.res.Write(json)
	if err != nil {
		return err
	}
	return nil
}

func (this *Context) Html(statusCode int, data any, template string) error {
	if this.GetIsTimeout() {
		return nil
	}
	this.res.WriteHeader(statusCode)
	this.res.Header().Set("Content-Type", "text/html")
	return nil
}

func (this *Context) Text(statusCode int, data string) error {
	if this.GetIsTimeout() {
		return nil
	}
	this.res.WriteHeader(statusCode)
	this.res.Header().Set("Content-Type", "text/plain")
	_, err := this.res.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}
