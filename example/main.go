package main

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/ryouaki/koa"
	"github.com/ryouaki/koa/catch"
	"github.com/ryouaki/koa/example/plugin"
	"github.com/ryouaki/koa/log"
	"github.com/ryouaki/koa/session"
)

func init() {
	catch.Try(func() interface{} {
		return "test1"
	}).Then(func(ret interface{}) {
		fmt.Println(ret.(string))
		panic(errors.New("test error"))
	}).Then(func(ret interface{}) {
		fmt.Println("不会被执行", ret.(string))
	}).Catch(func(err interface{}) {
		fmt.Println(err)
	}).Catch(func(err interface{}) {
		panic("test error 2")
		fmt.Println(err)
	}).Catch(func(err interface{}) {
		fmt.Println(err)
	}).Finally(func() {
		fmt.Println("end")
	})
}

func main() {
	app := koa.New()

	log.New(&log.Config{
		Level:   log.LevelInfo,
		Mode:    log.LogFile,
		MaxDays: 1,
		LogPath: "./../logs",
	})
	app.Use(plugin.Duration)
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
