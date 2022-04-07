package koa

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

// Application Object
type Application struct {
	handles []Handle
}

type Handle struct {
	url     string
	method  string
	handler []Handler
}

// Handler Func
type Handler func(err error, ctx *Context, n Next)

// Next Func
type Next func(err error)

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
	if reflect.TypeOf(argus[0]).String() == "string" {
		_url = reflect.ValueOf(argus[0]).String()
	}

	_handle := Handle{
		url:     _url,
		method:  "Use",
		handler: make([]Handler, 0, 8),
	}

	for i := 1; i < len(argus); i++ {
		_fb := argus[i]
		_handle.handler = append(_handle.handler, _fb.(func(error, *Context, Next)))
	}

	app.handles = append(app.handles, _handle)
}

func (app *Application) appendRouter(method string, path string, cbs []Handler) error {
	for _, handle := range app.handles {
		if handle.url == path && handle.method == method {
			return errors.New("router is exist")
		}
	}

	_handle := Handle{
		url:     path,
		method:  method,
		handler: make([]Handler, 0, 8),
	}

	_handle.handler = append(_handle.handler, cbs...)

	app.handles = append(app.handles, _handle)

	return nil
}

// Get func
func (app *Application) Get(path string, cbFunc ...Handler) error {
	app.appendRouter("Get", path, cbFunc)
	return nil
}

// Post func
func (app *Application) Post(path string, cbFunc ...Handler) error {
	app.appendRouter("Post", path, cbFunc)
	return nil
}

// Delete func
func (app *Application) Delete(path string, cbFunc ...Handler) error {
	app.appendRouter("Delete", path, cbFunc)
	return nil
}

// Patch func
func (app *Application) Patch(path string, cbFunc ...Handler) error {
	app.appendRouter("Patch", path, cbFunc)
	return nil
}

// Put func
func (app *Application) Put(path string, cbFunc ...Handler) error {
	app.appendRouter("Put", path, cbFunc)
	return nil
}

// Options func
func (app *Application) Options(path string, cbFunc ...Handler) error {
	app.appendRouter("Options", path, cbFunc)
	return nil
}

// Head func
func (app *Application) Head(path string, cbFunc ...Handler) error {
	app.appendRouter("Head", path, cbFunc)
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
	ctx := NewContext(w, req)
	fmt.Println(ctx.Url)
	// ctx.data["session"] = make(map[string]interface{})

	// var routerHandler []Handler
	// for _, middleware := range app.middlewares {
	// 	if ok := compare(middleware.path, ctx.Path, false); ok {
	// 		routerHandler = append(routerHandler, middleware.handler)
	// 	}
	// }

	// for _, router := range app.route[ctx.Method] {
	// 	if ok := compare(router.path, ctx.Path, true); ok {
	// 		ctx.RequestNotFound = false
	// 		ctx.MatchURL = router.path
	// 		ctx.Params = formatParams(router.path, ctx.Path)
	// 		routerHandler = append(routerHandler, router.handler...)
	// 	}
	// }

	// fb := compose(ctx, routerHandler)
	// fb(err)
}
