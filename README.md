# koa.go
Expressive HTTP middleware framework for Golang to make web applications and APIs more enjoyable to write like Koa.js. Koa's middleware stack flows in a stack-like manner, allowing you to perform actions downstream then filter and manipulate the response upstream.

Koa is not bundled with any middleware.

## Benchmarks
ryouaki/koa has the best performance:

| name | qps | transactions/s |
|--|--|--|
|koa(golang) | 12598080 | 42014.05 |
|gin | 11936723 | 39812.96 |
|echo | 12347680 | 41179.66	 |
|beego | 11051408 | 36856.09 |
|koa(nodejs) | 8042021 | 26818.72	 |

## Installation
```go
  $ go get github.com/ryouaki/koa
```

## Hello Koa
```go
  package main

  import (
    "fmt"

    "github.com/ryouaki/koa"
  )

  func main() {
    app := koa.New() // 初始化服务对象

    // 设置api路由，其中var为url传参
    app.Get("/", func(ctx *koa.Context, next koa.Next) {
      ctx.SetBody([]byte("Hello Koa"))
    })

    err := app.Run(8080) // 启动
    if err != nil {      // 是否发生错误
      fmt.Println(err)
    }
  }
```

## Middleware
Koa is a middleware framework, Here is an example of logger middleware with each of the different functions:

```go
  package plugin

  import (
    "fmt"
    "time"

    "github.com/ryouaki/koa"
  )

  // Log out about the duration for request.
  func Duration(ctx *koa.Context, next koa.Next) {
    startTime := time.Now() // 开始计时
    next(nil) // 执行后续操作
    d := time.Now().Sub(startTime) // 计算耗时
    // 打印结果
    fmt.Println(time.Date(startTime.Year(),
      startTime.Month(),
      startTime.Day(),
      startTime.Hour(),
      startTime.Minute(),
      startTime.Second(), 0, time.Local),
      ctx.URL,
      "Request cost: ",
      float64(d)/float64(time.Millisecond), "ms")
  }
```




## Koa Application
```go
  type Application struct {
    handles []Handle
    _cb     Handler
  }

  type Handle struct {
    path     string
    method  string
    handler Handler
  }

  // Handler Func
  type Handler func(ctx *Context, n Next)

  // Next Func
  type Next func()

  // New a Koa instance
  func New() *Application
  // Add a middleware for koa application
  /**
  * params: path<string|option> path for request
  * params: callback<koa.Handler|option> cb for request
  * params: callback ...
  */
  func (app *Application) Use(argus ...interface{})
  // Get func
  func (app *Application) Get(path string, cbFunc ...Handler) error
  // Post func
  func (app *Application) Post(path string, cbFunc ...Handler) error
  // Delete func
  func (app *Application) Delete(path string, cbFunc ...Handler) error
  // Patch func
  func (app *Application) Patch(path string, cbFunc ...Handler) error 
  // Put func
  func (app *Application) Put(path string, cbFunc ...Handler) error
  // Options func
  func (app *Application) Options(path string, cbFunc ...Handler) error
  // Head func
  func (app *Application) Head(path string, cbFunc ...Handler) error
  // Run func
  func (app *Application) Run(port int) error 
  // RunTLS func
  func (app *Application) RunTLS(port int, certFile string, keyFile string) error 
```

## About Context
```go
  type Context struct {
    Res      http.ResponseWriter    // Res the response object
    Req      *http.Request          // Req the request object from client
    Url      string                 // Url
    Path     string                 // Rqeust path
    Method   string                 // Method like Get, Post and others
    Status   int                    // Status you want to let client konw for the request
    MatchURL string                 // For the router path
    Body     []uint8                // The body from client
    Query    map[string]([]string)  // The Query from request's path
    Params   map[string](string)    // The Params from request's path
  }

  // New a context for request
  func NewContext(res http.ResponseWriter, req *http.Request) *Context
  // Get Header from request func
  func (ctx *Context) GetHeader(key string) []string 
  // Set header item for Response func
  func (ctx *Context) SetHeader(key string, value string)
  // Get cookie from request func
  func (ctx *Context) GetCookie(key string) *http.Cookie
  // Set cookie for response func
  func (ctx *Context) SetCookie(cookie *http.Cookie)
  // Set data for context func
  func (ctx *Context) GetData(key string) interface{}
  // Get data from contexxt  func
  func (ctx *Context) SetData(key string, value interface{})
  // Set data for response
  func (ctx *Context) SetBody(data []byte)
```

## Components
- [github.com/ryouaki/koa/log](https://github.com/ryouaki/koa/blob/main/log/log.md) Logger plugin
- [github.com/ryouaki/koa/session](https://github.com/ryouaki/koa/blob/main/session/session.md) The middleware for session. support two way to save data - memory and redis base on go-redis
- [github.com/ryouaki/koa/static](https://github.com/ryouaki/koa/blob/main/static/static.md) The middleware for Static.
