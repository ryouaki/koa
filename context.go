package koa

import "net/http"

// Context Request
type Context struct {
	Header   map[string]([]string)
	Res      http.ResponseWriter
	Req      *http.Request
	URL      string
	Path     string
	Query    map[string](string)
	Params   map[string](string)
	Method   string
	Status   int
	MatchURL string
	Body     []uint8
}

// Get func
func (ctx *Context) Get(key string) []string {
	return ctx.Header[key]
}

// Set func
func (ctx *Context) Set(key string, value string) {
	ctx.Res.Header().Set(key, value)
}

// IsFinish func
// func (ctx *Context) IsFinish() bool {
// 	return ctx.Res.handlerDone == 1
// }
