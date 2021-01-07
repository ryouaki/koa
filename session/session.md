# Koa.Session介绍
Session是一个koa的中间件，用于实现session的操作。目前支持内存存储和redis存储两种形式。

* redis存储主要依赖于[go-redis/redis](github.com/go-redis/redis)

## 如何用Koa.Session管理Session呢？
Session本身就是一个中间件，所以直接Use即可：
```go
  package main

  import (
    "errors"
    "fmt"

    "github.com/ryouaki/koa"
    "github.com/ryouaki/koa/session"
  )

  func main() {
    app := koa.New()

    rds := redis.NewUniversalClient(&redis.UniversalOptions{
      Addrs: []string{"42.192.194.38:6001"},
    })

    store := session.NewRedisStore(rds)
    app.Use(session.Session(&session.Config{
      Store:  store,
      MaxAge: 100,
    }))

    // store := session.NewMemStore()
    // app.Use(session.Session(&session.Config{
    // 	Store:  store,
    // 	MaxAge: 1000,
    // }))
    app.Use("/a", func(err error, ctx *koa.Context, next koa.NextCb) {
      ctx.SetSession("test", "kkkk")
      next(nil)
    })
    app.Get("/", func(err error, ctx *koa.Context, next koa.NextCb) {
      data := ctx.GetSession()
      if data["test"] == nil {
        data["test"] = "nil"
      }
      ctx.Write([]byte(data["test"].(string)))
    })

    err := app.Run(8080)
    if err != nil {
      fmt.Println(err)
    }
  }
```

当我们访问/a的时候，在session中设置了test=kkkk，当我们请求任意路径会将该信息返回：
```sh
  $ curl localhost:8080/a
  kkkk
  $ curl localhost:8080/
  kkkk

  // 操作100s后
  $ curl localhost:8080/
  nil
```