package koa

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Context Request
type Context struct {
	Res    http.ResponseWriter    // instance for response
	Req    *http.Request          // instance for request
	Url    string                 // Url for request
	Path   string                 // Path for request
	Method string                 // Method for request
	Status int                    // Status for response
	Body   []uint8                // Body from request
	Query  map[string]([]string)  // Query from request
	Params map[string](string)    // Params from request
	data   map[string]interface{} // cache for context
	body   []byte                 // body for response
}

// New a context for request
func NewContext(res http.ResponseWriter, req *http.Request) *Context {
	var body []uint8 = nil

	if req.Body != nil {
		body, _ = ioutil.ReadAll(req.Body)
	}

	return &Context{
		Res:    res,
		Req:    req,
		Url:    req.RequestURI,
		Path:   req.URL.Path,
		Method: strings.ToLower(req.Method),
		Status: 200,
		Body:   body,
		Query:  formatQuery(req.URL.Query()),
		data:   make(map[string]interface{}),
		body:   nil,
	}
}

// Get Header from request func
func (ctx *Context) GetHeader(key string) []string {
	if data, ok := ctx.Req.Header[key]; ok {
		return data
	}
	return nil
}

// Set header item for Response func
func (ctx *Context) SetHeader(key string, value string) {
	ctx.Res.Header().Set(key, value)
}

// Get cookie from request func
func (ctx *Context) GetCookie(key string) *http.Cookie {
	cookie, ok := ctx.Req.Cookie(key)
	if ok != nil {
		return nil
	}
	return cookie
}

// Set cookie for response func
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	http.SetCookie(ctx.Res, cookie)
}

// Set data for context func
func (ctx *Context) SetData(key string, value interface{}) {
	ctx.data[key] = value
}

// Get data from contexxt  func
func (ctx *Context) GetData(key string) interface{} {
	if data, ok := ctx.data[key]; ok {
		return data
	}
	return nil
}

// Set data for response
func (ctx *Context) SetBody(data []byte) {
	ctx.body = data
}
