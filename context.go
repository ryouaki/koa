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
