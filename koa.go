package koa

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

// Application Object
type Application struct {
	prefix      string
	middlewares []MiddlewareHandler
	route       map[string]([]RouterHandler)
}

// MiddlewareHandler struct
type MiddlewareHandler struct {
	path    string
	handler Handler
}

// RouterHandler struct
type RouterHandler struct {
	path    string
	handler []Handler
}

// Handler Func
type Handler func(err error, ctx *Context, next NextCb)

// NextCb Func
type NextCb func(err error)

// New for a koa instance
func New() *Application {
	return &Application{
		middlewares: make([]MiddlewareHandler, 0, 16),
		route:       make(map[string]([]RouterHandler)),
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

func (app *Application) initRouter(method string) []RouterHandler {
	if _, ok := app.route[method]; ok {
		app.route[method] = make([]RouterHandler, 0, 16)
	}

	return app.route[method]
}

func (app *Application) appendRouter(method string, path string, cbs []Handler) {
	app.route[method] = append(app.route[method], RouterHandler{
		path:    path,
		handler: cbs,
	})
}

// Get func
func (app *Application) Get(path string, cbFunc ...Handler) error {
	routers := app.initRouter("get")
	for _, router := range routers {
		if router.path == path {
			return errors.New("Router is exist")
		}
	}

	app.appendRouter("get", path, cbFunc)
	return nil
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

	if req.Body != nil {
		body, err = ioutil.ReadAll(req.Body)
	}

	ctx := &Context{
		Header: req.Header,
		Res:    w,
		Req:    req,
		URL:    req.RequestURI,
		Path:   req.URL.Path,
		Body:   body,
		Method: method,
		Status: 200,
	}

	var routerHandler []Handler
	for _, middleware := range app.middlewares {
		if ok := compare(middleware.path, ctx.Path); ok {
			routerHandler = append(routerHandler, middleware.handler)
		}
	}

	for _, router := range app.route[ctx.Method] {
		if ok := compare(router.path, ctx.Path); ok {
			routerHandler = append(routerHandler, router.handler...)
		}
	}

	fb := compose(ctx, routerHandler)
	fb(err)
}
