package framework

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

type IResponse interface {
	Json(obj interface{}) IResponse

	Jsonp(obj interface{}) IResponse

	Xml(obj interface{}) IResponse

	Html(file string, obj interface{}) IResponse

	Text(format string, values ...interface{}) IResponse

	Redirect(path string) IResponse

	SetHeader(key string, val string) IResponse

	SetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	SetStatus(code int) IResponse

	SetOkStatus() IResponse
}

func (c *Context) Json(obj interface{}) IResponse {
	json, err := json.Marshal(obj)
	if err != nil {
		return c.SetStatus(http.StatusInternalServerError)
	}
	c.SetHeader("Content-Type", "application/json")
	c.res.Write(json)
	return c
}

func (c *Context) Jsonp(obj interface{}) IResponse {
	callbackFunc, _ := c.QueryString("callback", "callback_func")
	c.SetHeader("Content-Type", "application/javascript")
	callbackFunc = template.JSEscapeString(callbackFunc)

	_, err := c.res.Write([]byte(callbackFunc))
	if err != nil {
		return c
	}
	_, err = c.res.Write([]byte("("))
	if err != nil {
		return c
	}
	json, err := json.Marshal(obj)
	if err != nil {
		return c
	}
	_, err = c.res.Write(json)
	if err != nil {
		return c
	}
	_, err = c.res.Write([]byte(")"))
	if err != nil {
		return c
	}
	return c
}

func (c *Context) Xml(obj interface{}) IResponse {
	xml, err := xml.Marshal(obj)
	if err != nil {
		return c.SetStatus(http.StatusInternalServerError)
	}
	c.SetHeader("Content-Type", "application/xml")
	c.res.Write(xml)
	return c
}

func (c *Context) Html(file string, obj interface{}) IResponse {
	files, err := template.New("output").ParseFiles(file)
	if err != nil {
		return c
	}
	if err := files.Execute(c.res, obj); err != nil {
		return c
	}
	c.SetHeader("Content-Type", "application/html")
	return c
}

func (c *Context) Text(format string, values ...interface{}) IResponse {
	text := fmt.Sprintf(format, values...)
	c.SetHeader("Content-Type", "application/text")
	c.res.Write([]byte(text))
	return c
}

func (c *Context) Redirect(path string) IResponse {
	http.Redirect(c.res, c.req, path, http.StatusMovedPermanently)
	return c
}

func (c *Context) SetHeader(key string, val string) IResponse {
	c.res.Header().Add(key, val)
	return c
}

func (c *Context) SetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.res, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return c
}

func (c *Context) SetStatus(code int) IResponse {
	c.res.WriteHeader(code)
	return c
}

func (c *Context) SetOkStatus() IResponse {
	c.res.WriteHeader(http.StatusOK)
	return c
}
