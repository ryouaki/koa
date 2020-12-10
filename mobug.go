package mkbug

import (
	"fmt"
	"net/http"
	"reflect"
)

// Application Object
type Application struct {
	prefix      string
	middlewares map[string]([]Handler)
}

// Handler Func
type Handler func(err error, ctx *Context, next NextCb)

// NextCb Func
type NextCb func(err error)

// New for a Mkbug instance
func New() *Application {
	return &Application{
		middlewares: make(map[string]([]Handler)),
	}
}

// Use func
func (app *Application) Use(argus ...interface{}) {
	if len(argus) <= 0 {
		return
	}

	firstArgu := reflect.ValueOf(argus[0]).String()
	var middleware []interface{}
	if reflect.TypeOf(firstArgu).String() != "string" {
		firstArgu = "/"
		middleware = argus
	} else {
		middleware = argus[1:]
	}

	if _, ok := app.middlewares[firstArgu]; ok {
		app.middlewares[firstArgu] = make([]Handler, 0, 16)
	}

	tmp := make([]Handler, 0, 16)
	for _, fb := range middleware {
		tmp = append(tmp, fb.(func(error, *Context, NextCb)))
		// fmt.Println(reflect.TypeOf(fb))
	}

	app.middlewares[firstArgu] = append(app.middlewares[firstArgu], tmp...)
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
	// var err error = nil
	// var body []uint8 = nil

	// method := strings.ToLower(req.Method)
	// // fullUri := req.RequestURI
	// // path := req.URL.Path

	// if method != "get" && method != "options" && method != "head" {
	// 	body, err = ioutil.ReadAll(req.Body)
	// }

	// ctx := &Context{
	// 	Header: req.Header,
	// 	Res:    w,
	// 	Req:    req,
	// 	URL:    req.RequestURI,
	// 	Path:   req.URL.Path,
	// 	Query:  req.URL.RawQuery,
	// 	Body:   body,
	// 	Method: method,
	// 	Status: 200,
	// }

	// currentPoint := int32(-1)
	// var next func(err error)

	// next = func(err error) {
	// 	atomic.AddInt32(&currentPoint, 1)

	// 	_ctx := ctx
	// 	_middlewares := app.middlewares

	// 	cb := _middlewares[currentPoint]
	// 	cb(err, _ctx, next)
	// }

	// next(err)

	// fmt.Println(Bytes2String(body))
	// w.Write([]byte(ctx.Get("Content-Type")[0]))
}
