# koa.go
这是一款Golang版本的koa框架，完整的实现了koa.js的中间件机制，洋葱模型，错误优先机制，并且提供了一个抽象的ctx.SetData及ctx.GetData可以向下传递数据。

## 如何使用
使用起来极其简单，初始化应用对象，添加中间件，设置路由回调，启动程序

```go
  package main

  import "github.com/ryouaki/koa.go"

  func main() {
    app := koa.New() // 初始化服务对象

    app.Use(plugin.Duration) // 使用中间件，这是一个打印当前请求耗费时长的中间件

    // 自定义中间件，为了验证洋葱模型
    app.Use(func(err error, ctx *koa.Context, next koa.NextCb) {
      fmt.Println("test2")
      next(err) // 中间件必须显式调用next使请求进入下一个处理
      fmt.Println("test2")
    }, func(err error, ctx *koa.Context, next koa.NextCb) {
      fmt.Println("test3")
      // 在数据缓存区设置参数，因为golang无法像js那样动态添加属性，所以需要实现一个map[string]interface{}接口缓存数据
      ctx.SetData("test", ctx.Query["c"][0])
      next(nil) // 中间件必须显式调用next使请求进入下一个处理
      fmt.Println("test3")
    })

    // 设置api路由，其中var为url传参
    app.Get("/test/:var/p", func(err error, ctx *koa.Context, next koa.NextCb) {
      fmt.Println("test", ctx.Params) // 打印url参数
      next(nil) // 执行下一个回调
    }, func(err error, ctx *koa.Context, next koa.NextCb) {
      // 将query内的参数回传给客户端
      ctx.Write([]byte(ctx.GetData("test").(string)))
    })

    err := app.Run(8080) // 启动
    if err != nil { // 是否发生错误
      fmt.Println(err)
    }
  }
```

计时中间件的实现
```go
  package plugin

  import (
    "fmt"
    "time"

    "github.com/ryouaki/koa.go"
  )

  // 中间件和路由api接口实现都必须是func (err error, ctx *koa.Context, next koa.NextCb)类型
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

其中Context提供了多种接口。
```go
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
    data     map[string]interface{}
  }

  // 获取请求中header内容
  func (ctx *Context) GetHeader(key string) []string 
  // 设置响应的header内容
  func (ctx *Context) SetHeader(key string, value string)
  // 读取请求中的Cookie
  func (ctx *Context) GetCookie(key string) *http.Cookie
  // 设置响应的Cookie
  func (ctx *Context) SetCookie(cookie *http.Cookie)
  // 向请求上下文中设置数据
  func (ctx *Context) SetData(key string, value interface{})
  // 从请求上下文中读取缓存数据
  func (ctx *Context) GetData(key string) interface{}
  // 设置响应返回内容
  func (ctx *Context) Write(data []byte) (int, error)
  // 该请求是否已经结束，这个非常重要。一个请求只能由一个地方进行结束。否则无法保证返回内容的可预测性
  func (ctx *Context) IsFinished() bool
```

# 组件
- github.com/ryouaki/koa/log 日志组件。支持控制台输出，文件输出，文件保留日期的设置