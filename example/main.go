package main

import (
	"fmt"

	"github.com/ryouaki/koa"
	"github.com/ryouaki/koa/example/plugin"
	"github.com/ryouaki/koa/log"
)

func main() {
	app := koa.New()

	log.New(&log.Config{
		Level:   log.LevelInfo,
		Mode:    log.LogFile,
		MaxDays: 1,
		LogPath: "./../logs",
	})
	app.Use(plugin.Duration)
	app.Use("/", func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("test1")
		log.Info("Request in")
		next(err)
		fmt.Println("test1")
	})
	app.Use(func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("test2")
		next(err)
		fmt.Println("test2")
	}, func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("test3")
		ctx.SetData("test", ctx.Query["c"][0])
		next(nil)
		fmt.Println("test3")
	})

	app.Get("/test/:var/p", func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("test", ctx.Params)
		next(nil)
	}, func(err error, ctx *koa.Context, next koa.NextCb) {
		ctx.Write([]byte(ctx.GetData("test").(string)))
	})

	err := app.Run(8080)
	if err != nil {
		fmt.Println(err)
	}
}
