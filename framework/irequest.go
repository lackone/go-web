package framework

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"mime/multipart"
	"strconv"
)

const defaultMultipartMemory = 32 << 20 // 32 MB

type IRequest interface {
	//查询URL中的参数
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)
	Query(key string) interface{}

	//路由匹配中的参数
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	Param(key string) interface{}

	//form表单中带的参数
	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStringSlice(key string, def []string) ([]string, bool)
	FormFile(key string) (*multipart.FileHeader, error)
	Form(key string) interface{}

	//绑定JSON
	BindJson(obj interface{}) error
	//绑定XML
	BindXml(obj interface{}) error

	//获取原始数据
	GetRawData() ([]byte, error)

	//基本信息
	Uri() string
	Method() string
	Host() string
	ClientIp() string

	//头信息
	Headers() map[string][]string
	Header(key string) (string, bool)

	//cookie
	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

func (c *Context) QueryAll() map[string][]string {
	if c.req != nil {
		return c.req.URL.Query()
	}
	return map[string][]string{}
}

func (c *Context) QueryInt(key string, def int) (int, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.Atoi(val[len-1])
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) QueryInt64(key string, def int64) (int64, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseInt(val[len-1], 10, 64)
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) QueryFloat32(key string, def float32) (float32, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseFloat(val[len-1], 32)
			if err != nil {
				return def, false
			}
			return float32(v), true
		}
	}
	return def, false
}

func (c *Context) QueryFloat64(key string, def float64) (float64, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseFloat(val[len-1], 64)
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) QueryBool(key string, def bool) (bool, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseBool(val[len-1])
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) QueryString(key string, def string) (string, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val[len-1], true
		}
	}
	return def, false
}

func (c *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val, true
		}
	}
	return def, false
}

func (c *Context) Query(key string) interface{} {
	params := c.QueryAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val
		}
	}
	return nil
}

func (c *Context) GetParam(key string) string {
	if c.params != nil {
		if val, ok := c.params[key]; ok {
			return val
		}
	}
	return ""
}

func (c *Context) ParamInt(key string, def int) (int, bool) {
	val := c.GetParam(key)
	if val != "" {
		v, err := strconv.Atoi(val)
		if err != nil {
			return def, false
		}
		return v, true
	}
	return def, false
}

func (c *Context) ParamInt64(key string, def int64) (int64, bool) {
	val := c.GetParam(key)
	if val != "" {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return def, false
		}
		return v, true
	}
	return def, false
}

func (c *Context) ParamFloat32(key string, def float32) (float32, bool) {
	val := c.GetParam(key)
	if val != "" {
		v, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return def, false
		}
		return float32(v), true
	}
	return def, false
}

func (c *Context) ParamFloat64(key string, def float64) (float64, bool) {
	val := c.GetParam(key)
	if val != "" {
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return def, false
		}
		return v, true
	}
	return def, false
}

func (c *Context) ParamBool(key string, def bool) (bool, bool) {
	val := c.GetParam(key)
	if val != "" {
		v, err := strconv.ParseBool(val)
		if err != nil {
			return def, false
		}
		return v, true
	}
	return def, false
}

func (c *Context) ParamString(key string, def string) (string, bool) {
	val := c.GetParam(key)
	if val != "" {
		return val, true
	}
	return def, false
}

func (c *Context) Param(key string) interface{} {
	return c.GetParam(key)
}

func (c *Context) FormAll() map[string][]string {
	if c.req != nil {
		c.req.ParseForm()
		return c.req.PostForm
	}
	return map[string][]string{}
}

func (c *Context) FormInt(key string, def int) (int, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.Atoi(val[len-1])
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) FormInt64(key string, def int64) (int64, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseInt(val[len-1], 10, 64)
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) FormFloat32(key string, def float32) (float32, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseFloat(val[len-1], 32)
			if err != nil {
				return def, false
			}
			return float32(v), true
		}
	}
	return def, false
}

func (c *Context) FormFloat64(key string, def float64) (float64, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseFloat(val[len-1], 64)
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) FormBool(key string, def bool) (bool, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			v, err := strconv.ParseBool(val[len-1])
			if err != nil {
				return def, false
			}
			return v, true
		}
	}
	return def, false
}

func (c *Context) FormString(key string, def string) (string, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val[len-1], true
		}
	}
	return def, false
}

func (c *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val, true
		}
	}
	return def, false
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if c.req.MultipartForm == nil {
		if err := c.req.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	file, header, err := c.req.FormFile(key)
	if err != nil {
		return nil, err
	}
	file.Close()
	return header, nil
}

func (c *Context) Form(key string) interface{} {
	params := c.FormAll()
	if val, ok := params[key]; ok {
		len := len(val)
		if len > 0 {
			return val
		}
	}
	return nil
}

func (c *Context) BindJson(obj interface{}) error {
	if c.req != nil {
		all, err := io.ReadAll(c.req.Body)
		if err != nil {
			return err
		}
		//body只能读一次，读出来后需要重置下body
		c.req.Body = io.NopCloser(bytes.NewBuffer(all))

		err = json.Unmarshal(all, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("request empty")
	}
	return nil
}

func (c *Context) BindXml(obj interface{}) error {
	if c.req != nil {
		all, err := io.ReadAll(c.req.Body)
		if err != nil {
			return err
		}
		//body只能读一次，读出来后需要重置下body
		c.req.Body = io.NopCloser(bytes.NewBuffer(all))

		err = xml.Unmarshal(all, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("request empty")
	}
	return nil
}

func (c *Context) GetRawData() ([]byte, error) {
	if c.req != nil {
		all, err := io.ReadAll(c.req.Body)
		if err != nil {
			return nil, err
		}
		//body只能读一次，读出来后需要重置下body
		c.req.Body = io.NopCloser(bytes.NewBuffer(all))
		return all, nil
	}
	return nil, errors.New("request empty")
}

func (c *Context) Uri() string {
	return c.req.RequestURI
}

func (c *Context) Method() string {
	return c.req.Method
}

func (c *Context) Host() string {
	return c.req.URL.Host
}

func (c *Context) ClientIp() string {
	ip := c.req.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = c.req.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.req.RemoteAddr
	}
	return ip
}

func (c *Context) Headers() map[string][]string {
	return c.req.Header
}

func (c *Context) Header(key string) (string, bool) {
	values := c.req.Header.Values(key)
	if values == nil || len(values) <= 0 {
		return "", false
	}
	return values[0], true
}

func (c *Context) Cookies() map[string]string {
	cookies := c.req.Cookies()
	ret := map[string]string{}
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}

func (c *Context) Cookie(key string) (string, bool) {
	cookies := c.Cookies()
	if val, ok := cookies[key]; ok {
		return val, true
	}
	return "", false
}
