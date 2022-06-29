package main

import (
	"fmt"

	"github.com/ryouaki/koa"
)

// func init() {
// 	catch.Try(func() interface{} {
// 		return "test1"
// 	}).Then(func(ret interface{}) {
// 		fmt.Println(ret.(string))
// 		panic(errors.New("test error"))
// 	}).Then(func(ret interface{}) {
// 		fmt.Println("不会被执行", ret.(string))
// 	}).Catch(func(err interface{}) {
// 		fmt.Println(err)
// 	}).Catch(func(err interface{}) {
// 		panic("test error 2")
// 		fmt.Println(err)
// 	}).Catch(func(err interface{}) {
// 		fmt.Println(err)
// 	}).Finally(func() {
// 		fmt.Println("end")
// 	})
// }

func main() {
	app := koa.New()

	fmt.Println(koa.GetIPAddr())

	handler1 := func(ctx *koa.Context, n koa.Next) {
		fmt.Println("handler1")
		n()
		fmt.Println("handler1")
	}

	handler2 := func(ctx *koa.Context, n koa.Next) {
		fmt.Println("handler2")
		ctx.SetBody([]byte("Hello world"))
		fmt.Println("handler2")
	}

	app.Get("/test", handler1, handler2)

	app.Get("/test1/:a/:c", handler1, handler2)

	app.Run(8080)
}

// log.New(&log.Config{
// 	Level:   log.LevelInfo,
// 	Mode:    log.LogFile,
// 	MaxDays: 1,
// 	LogPath: "./../logs",
// })
// app.Use(plugin.Duration)
// rds := redis.NewUniversalClient(&redis.UniversalOptions{
// 	Addrs: []string{"42.192.194.38:6001"},
// })

// // store := session.NewRedisStore(rds)
// // app.Use(session.Session(&session.Config{
// // 	Store:  store,
// // 	MaxAge: 100,
// // }))

// // store := session.NewMemStore()
// // app.Use(session.Session(&session.Config{
// // 	Store:  store,
// // 	MaxAge: 1000,
// // }))
// app.Use(static.Static("./static", "/static/"))
// app.Use(func(err error, ctx *koa.Context, next koa.Next) {
// 	fmt.Println("b-use1")
// 	next(nil)
// 	fmt.Println("b-us2e")
// })
// app.Use("/a", func(err error, ctx *koa.Context, next koa.Next) {
// 	fmt.Println("a")
// 	next(nil)
// })
// app.Get("/b", func(err error, ctx *koa.Context, next koa.Next) {
// 	fmt.Println("b1")
// 	ctx.SetBody([]byte("Hello world"))
// 	fmt.Println("b2")
// })
// app.Get("/c", func(err error, ctx *koa.Context, next koa.Next) {
// 	ctx.Status = 500
// 	fmt.Println("c")
// 	// next(nil)
// })

// app.Get("/d", func(err error, ctx *koa.Context, next koa.Next) {
// 	fmt.Println("d1")
// 	next(nil)
// }, func(err error, ctx *koa.Context, next koa.Next) {
// 	fmt.Println("d2")
// 	next(nil)
// })

// app.Get("/e/:f", func(err error, ctx *koa.Context, next koa.Next) {
// 	fmt.Println("c")
// 	// next(nil)
// })

// app.Get("/json", func(err error, ctx *koa.Context, next koa.Next) {
// 	ctx.SetHeader("Content-Type", "application/json")
// 	data := make(map[string]interface{})
// 	data["test"] = "test"
// 	ret, _ := json.Marshal(data)
// 	ctx.SetBody(ret)
// })

// app.Use(func(err error, ctx *koa.Context, next koa.Next) {
// 	ctx.Status = 404
// 	ctx.SetBody([]byte("Request not found"))
// })
// err := app.Run(8080)
// if err != nil {
// 	fmt.Println(err)
// }
// }
