package main

import (
	"fmt"

	"koa.go"
)

func main() {
	app := koa.New()

	app.Use("/", func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("test1")
		next(err)
		fmt.Println("test1")
	})
	app.Use(func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("test2")
		next(err)
		fmt.Println("test2")
	})

	app.Use(func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("test3")
		next(nil)
		fmt.Println("test3")
	})

	app.Use("/aaa", func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("test4")
		next(nil)
		fmt.Println("test4")
	})

	app.Get("/aaa", func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("aaa")
	})
	err := app.Run(8080)
	if err != nil {
		fmt.Println(err)
	}
}
