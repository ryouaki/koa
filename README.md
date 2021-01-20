# koa.go
Expressive HTTP middleware framework for Golang to make web applications and APIs more enjoyable to write like Koa.js. Koa's middleware stack flows in a stack-like manner, allowing you to perform actions downstream then filter and manipulate the response upstream.

Koa is not bundled with any middleware.

[中文](README_CN.md)

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
    app.Get("/", func(err error, ctx *koa.Context, next koa.NextCb) {
      ctx.Write([]byte("Hello Koa"))
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

    "github.com/ryouaki/koa.go"
  )

  // Log out about the duration for request.
  func Duration(err error, ctx *koa.Context, next koa.NextCb) {
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

## About Context
```go
  type Context struct {
    Header   map[string]([]string)  // Header which you get from client
    Res      http.ResponseWriter    // Res the response object
    Req      *http.Request          // Req the request object from client
    URL      string                 // Url
    Path     string                 // Rqeust path
    Method   string                 // Method like Get, Post and others
    Status   int                    // Status you want to let client konw for the request
    MatchURL string                 // For the router path
    Body     []uint8                // The body from client
    Query    map[string]([]string)  // The Query from request's url
    Params   map[string](string)    // The Params from request's path
    IsFinish bool                   // One request only can be done by one time
    data     map[string]interface{} // ...
  }

  // Get the information from request's header 
  func (ctx *Context) GetHeader(key string) []string 
  // Set the information to response's header
  func (ctx *Context) SetHeader(key string, value string)
  // Get the information from request's cookie 
  func (ctx *Context) GetCookie(key string) *http.Cookie
  // Set the information to response's cookie
  func (ctx *Context) SetCookie(cookie *http.Cookie)
  // Get the information from context
  func (ctx *Context) GetData(key string) interface{}
  // Set the information to context
  func (ctx *Context) SetData(key string, value interface{})
  // Set the information to session, but you should use session middleware first
  func (ctx *Context) SetSession(key string, value interface{}) 
  // Update the session, but you should use session middleware first
  func (ctx *Context) UpdateSession(sess map[string]interface{})
  // Get the information from session, but you should use session middleware first
  func (ctx *Context) GetSession() map[string]interface{}
  // Send the data for response
  func (ctx *Context) Write(data []byte) (int, error)
  // Check if the response is done, it's very important for middleware.
  func (ctx *Context) IsFinished() bool
```

## Components
- [github.com/ryouaki/koa/log](https://github.com/ryouaki/koa/blob/main/log/log.md) Logger plugin
- [github.com/ryouaki/koa/catch](https://github.com/ryouaki/koa/blob/main/catch/catch.md) A library for catch exception like try - then - catch - finally
- [github.com/ryouaki/koa/session](https://github.com/ryouaki/koa/blob/main/session/session.md) The middleware for session. support two way to save data - memory and redis base on go-redis
