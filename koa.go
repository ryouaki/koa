package koa

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// Application Object
type Application struct {
	handles []Handle
	_cb     Handler
}

type Handle struct {
	url     string
	length  int
	method  string
	handler Handler
}

// Handler Func
type Handler func(err error, ctx *Context, n Next)

// Next Func
type Next func(err error)

const (
	GET     = "get"
	POST    = "post"
	PUT     = "put"
	DELETE  = "delete"
	PATCH   = "patch"
	OPTIONS = "options"
	HEAD    = "head"
	USE     = "use"
)

// New a Koa instance
func New() *Application {
	return &Application{
		handles: make([]Handle, 0, 8),
	}
}

// Add a middleware for koa application
/**
 * params: path<string|option> url for request
 * params: callback<koa.Handler|option> cb for request
 * params: callback ...
 */
func (app *Application) Use(argus ...interface{}) {
	if len(argus) == 0 {
		return
	}

	_url := ""
	_cbs := argus
	if reflect.TypeOf(argus[0]).String() == "string" {
		_url = reflect.ValueOf(argus[0]).String()
		_cbs = argus[1:]
	}

	var _handlers []Handler
	for _, v := range _cbs {
		_handlers = append(_handlers, v.(func(error, *Context, Next)))
	}

	app.appendRouter(USE, _url, _handlers)
}

func (app *Application) appendRouter(method string, path string, cbs []Handler) error {
	_sarr := strings.Split(path, "/")

	for _, v := range cbs {
		_handle := Handle{
			url:     path,
			length:  len(_sarr),
			method:  strings.ToLower(method),
			handler: v,
		}

		app.handles = append(app.handles, _handle)
	}

	return nil
}

// Get func
func (app *Application) Get(path string, cbFunc ...Handler) error {
	app.appendRouter(GET, path, cbFunc)
	return nil
}

// Post func
func (app *Application) Post(path string, cbFunc ...Handler) error {
	app.appendRouter(POST, path, cbFunc)
	return nil
}

// Delete func
func (app *Application) Delete(path string, cbFunc ...Handler) error {
	app.appendRouter(DELETE, path, cbFunc)
	return nil
}

// Patch func
func (app *Application) Patch(path string, cbFunc ...Handler) error {
	app.appendRouter(PATCH, path, cbFunc)
	return nil
}

// Put func
func (app *Application) Put(path string, cbFunc ...Handler) error {
	app.appendRouter(PUT, path, cbFunc)
	return nil
}

// Options func
func (app *Application) Options(path string, cbFunc ...Handler) error {
	app.appendRouter(OPTIONS, path, cbFunc)
	return nil
}

// Head func
func (app *Application) Head(path string, cbFunc ...Handler) error {
	app.appendRouter(HEAD, path, cbFunc)
	return nil
}

// Run func
func (app *Application) Run(port int) error {
	app._cb = compose(app.handles)
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, app)
}

// RunTLS func
func (app *Application) RunTLS(port int, certFile string, keyFile string) error {
	app._cb = compose(app.handles)
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServeTLS(addr, certFile, keyFile, app)
}

// ServeHTTP interface func
func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(w, req)
	app._cb(nil, ctx, nil)

	ctx.Res.WriteHeader(ctx.Status)
	ctx.Res.Write(ctx.body)
}
