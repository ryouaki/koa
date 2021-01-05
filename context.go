package koa

import (
	"errors"
	"net/http"
)

// Context Request
type Context struct {
	Header   map[string]([]string)
	Res      http.ResponseWriter
	Req      *http.Request
	URL      string
	Path     string
	Method   string
	Status   int
	MatchURL string
	Body     []uint8
	Query    map[string]([]string)
	Params   map[string](string)
	IsFinish bool
	data     map[string]interface{}
}

// GetHeader func
func (ctx *Context) GetHeader(key string) []string {
	if data, ok := ctx.Header[key]; ok {
		return data
	}
	return nil
}

// SetHeader func
func (ctx *Context) SetHeader(key string, value string) {
	ctx.Res.Header().Set(key, value)
}

// GetCookie func
func (ctx *Context) GetCookie(key string) *http.Cookie {
	cookie, ok := ctx.Req.Cookie(key)
	if ok != nil {
		return nil
	}
	return cookie
}

// SetCookie func
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	http.SetCookie(ctx.Res, cookie)
}

// SetData func
func (ctx *Context) SetData(key string, value interface{}) {
	if value == nil {
		return
	}
	ctx.data[key] = value
}

// GetData func
func (ctx *Context) GetData(key string) interface{} {
	if data, ok := ctx.data[key]; ok {
		return data
	}
	return nil
}

func (ctx *Context) Write(data []byte) (int, error) {
	if ctx.IsFinish {
		return -1, errors.New("Do not write data to response after sended")
	}

	ctx.Res.WriteHeader(ctx.Status)
	code, err := ctx.Res.Write(data)

	if err == nil {
		ctx.IsFinish = true
	}
	return code, err
}

// IsFinished func
func (ctx *Context) IsFinished() bool {
	return ctx.IsFinish == true
}
