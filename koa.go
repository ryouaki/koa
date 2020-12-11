package koa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync/atomic"
)

// Application Object
type Application struct {
	prefix      string
	middlewares []MiddlewareHandler
}

// MiddlewareHandler struct
type MiddlewareHandler struct {
	path    string
	handler Handler
}

// Handler Func
type Handler func(err error, ctx *Context, next NextCb)

// NextCb Func
type NextCb func(err error)

// New for a koa instance
func New() *Application {
	return &Application{
		middlewares: make([]MiddlewareHandler, 0, 16),
	}
}

// Use func
func (app *Application) Use(argus ...interface{}) {
	if len(argus) <= 0 {
		return
	}

	var firstArgu string

	var middleware []interface{}
	if reflect.TypeOf(argus[0]).String() == "string" {
		firstArgu = reflect.ValueOf(argus[0]).String()
		middleware = argus[1:]
	} else {
		firstArgu = "/"
		middleware = argus
	}

	for _, fb := range middleware {
		app.middlewares = append(app.middlewares, MiddlewareHandler{
			path:    firstArgu,
			handler: fb.(func(error, *Context, NextCb)),
		})
	}
}

// Get func
func (app *Application) Get(path string, cbFunc ...Handler) {

}

// Run func
func (app *Application) Run(port int) error {
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, app)
}

// RunTLS func
func (app *Application) RunTLS(port int, certFile string, keyFile string) error {
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServeTLS(addr, certFile, keyFile, app)
}

// ServeHTTP interface func
func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err error = nil
	var body []uint8 = nil

	method := strings.ToLower(req.Method)
	// // fullUri := req.RequestURI
	// // path := req.URL.Path

	if req.Body != nil {
		body, err = ioutil.ReadAll(req.Body)
	}

	ctx := &Context{
		Header: req.Header,
		Res:    w,
		Req:    req,
		URL:    req.RequestURI,
		Path:   req.URL.Path,
		Query:  req.URL.RawQuery,
		Body:   body,
		Method: method,
		Status: 200,
	}

	currentPoint := int32(0)
	var next func(err error)
	var cbMax = len(app.middlewares)

	next = func(err error) {
		_ctx := ctx
		_middlewares := app.middlewares
		bFound := false

		for int(currentPoint) < cbMax && bFound == false {
			curMiddleware := _middlewares[currentPoint]
			atomic.AddInt32(&currentPoint, 1)

			if compare(curMiddleware.path, ctx.Path) {
				bFound = true
				cb := curMiddleware.handler
				cb(err, _ctx, next)
			}
		}
	}

	next(err)
}

func compare(path string, target string) bool {
	return true
}
