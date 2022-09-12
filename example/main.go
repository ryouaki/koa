package main

import (
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

	// rds := redis.NewUniversalClient(&redis.UniversalOptions{
	// 	Addrs: []string{"42.192.194.38:6001"},
	// })
	// store := session.NewRedisStore(rds)
	// app.Use(session.Session(session.SessionConf{
	// 	MaxAge: 1000,
	// }, store))

	// app.Get("/", func(ctx *koa.Context, next koa.Next) {
	// 	sess := ctx.GetData("session")
	// 	_sess := sess.(map[string]interface{})
	// 	_sess["data"] = "hello world"
	// 	ctx.SetData("session", _sess)
	// })
	// app.Get("/ping/:aaa", func(ctx *koa.Context, next koa.Next) {
	// 	ctx.SetBody([]byte("Hello World"))
	// })

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
	// app.Use(func(ctx *koa.Context, next koa.Next) {
	// 	fmt.Println("b-use1")
	// 	next()
	// 	fmt.Println("b-us2e")
	// })
	// app.Use("/a", func(ctx *koa.Context, next koa.Next) {
	// 	fmt.Println("a")
	// 	next()
	// })
	// app.Get("/b", func(ctx *koa.Context, next koa.Next) {
	// 	fmt.Println("b1")
	// 	ctx.SetBody([]byte("Hello world"))
	// 	fmt.Println("b2")
	// })
	// app.Get("/c", func(ctx *koa.Context, next koa.Next) {
	// 	ctx.Status = 500
	// 	fmt.Println("c")
	// 	// next(nil)
	// })

	// app.Get("/d", func(ctx *koa.Context, next koa.Next) {
	// 	fmt.Println("d1")
	// 	next()
	// }, func(ctx *koa.Context, next koa.Next) {
	// 	fmt.Println("d2")
	// 	next()
	// })

	// app.Get("/e/:f", func(ctx *koa.Context, next koa.Next) {
	// 	fmt.Println("c", ctx.Params)
	// 	// next(nil)
	// })

	// app.Get("/json", func(ctx *koa.Context, next koa.Next) {
	// 	ctx.SetHeader("Content-Type", "application/json")
	// 	data := make(map[string]interface{})
	// 	data["test"] = "test"
	// 	ret, _ := json.Marshal(data)
	// 	ctx.SetBody(ret)
	// })

	// app.Use(func(ctx *koa.Context, next koa.Next) {
	// 	ctx.Status = 404
	// 	ctx.SetBody([]byte("Request not found"))
	// })
	// err := app.Run(8080)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// }

	app.Run(8080)
}
