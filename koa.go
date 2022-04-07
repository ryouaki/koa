package koa

import (
	"fmt"
	"net/http"
	"reflect"
)

// Application Object
type Application struct {
	handles []Handle
}

type Handle struct {
	Url     string
	Handler []Handler
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

	_url := "/"
	if reflect.TypeOf(argus[0]).String() == "string" {
		_url = reflect.ValueOf(argus[0]).String()
	}

	_handle := Handle{
		Url:     _url,
		Handler: make([]Handler, 0, 8),
	}

	for i := 1; i < len(argus); i++ {
		_fb := argus[i]
		_handle.Handler = append(_handle.Handler, _fb.(Handler))
	}
}

// // Use func
// func (app *Application) Use(argus ...interface{}) {
// 	if len(argus) <= 0 {
// 		return
// 	}

// 	var firstArgu string

// 	var middleware []interface{}
// 	if reflect.TypeOf(argus[0]).String() == "string" {
// 		firstArgu = reflect.ValueOf(argus[0]).String()
// 		middleware = argus[1:]
// 	} else {
// 		firstArgu = "/"
// 		middleware = argus
// 	}

// 	for _, fb := range middleware {
// 		app.middlewares = append(app.middlewares, MiddlewareHandler{
// 			path:    firstArgu,
// 			handler: fb.(func(error, *Context, NextCb)),
// 		})
// 	}
// }

// func (app *Application) initRouter(method string) []RouterHandler {
// 	if _, ok := app.route[method]; !ok {
// 		app.route[method] = make([]RouterHandler, 0, 16)
// 	}

// 	return app.route[method]
// }

// func (app *Application) appendRouter(method string, path string, cbs []Handler) error {
// 	routers := app.initRouter(method)

// 	for _, router := range routers {
// 		if router.path == path {
// 			return errors.New("Router is exist")
// 		}
// 	}

// 	app.route[method] = append(app.route[method], RouterHandler{
// 		path:    path,
// 		handler: cbs,
// 	})

// 	return nil
// }

// // Get func
// func (app *Application) Get(path string, cbFunc ...Handler) error {
// 	app.appendRouter("get", path, cbFunc)
// 	return nil
// }

// // Post func
// func (app *Application) Post(path string, cbFunc ...Handler) error {
// 	app.appendRouter("post", path, cbFunc)
// 	return nil
// }

// // Delete func
// func (app *Application) Delete(path string, cbFunc ...Handler) error {
// 	app.appendRouter("delete", path, cbFunc)
// 	return nil
// }

// // Patch func
// func (app *Application) Patch(path string, cbFunc ...Handler) error {
// 	app.appendRouter("patch", path, cbFunc)
// 	return nil
// }

// // Put func
// func (app *Application) Put(path string, cbFunc ...Handler) error {
// 	app.appendRouter("put", path, cbFunc)
// 	return nil
// }

// // Options func
// func (app *Application) Options(path string, cbFunc ...Handler) error {
// 	app.appendRouter("options", path, cbFunc)
// 	return nil
// }

// // Head func
// func (app *Application) Head(path string, cbFunc ...Handler) error {
// 	app.appendRouter("head", path, cbFunc)
// 	return nil
// }

// Run func
func (app *Application) Run(port int) error {
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, app)
}

// // RunTLS func
// func (app *Application) RunTLS(port int, certFile string, keyFile string) error {
// 	addr := fmt.Sprintf(":%d", port)
// 	return http.ListenAndServeTLS(addr, certFile, keyFile, app)
// }

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
