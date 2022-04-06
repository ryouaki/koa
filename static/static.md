# Koa.Static介绍
Static是一个koa的中间件，用于实现静态资源服务。

## 如何用Koa.Static管理Static呢？
Static本身就是一个中间件，所以直接Use即可：
```go
  package main

  import (
    "errors"
    "fmt"

    "github.com/ryouaki/koa"
    "github.com/ryouaki/koa/static"
  )

  func main() {
    app := koa.New()

    app.Use(static.Static("./statics", "/static/"))

    err := app.Run(8080)
    if err != nil {
      fmt.Println(err)
    }
  }
```

当我们访问/static/1.png的时候会请求到./statics/1.png
```sh
  $ curl localhost:8080/static/1.png
```