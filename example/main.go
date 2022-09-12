package main

import (
	"github.com/ryouaki/koa"
)

func main() {
	app := koa.New()

	app.Get("/1/:a", func(ctx *koa.Context, n koa.Next) {
		ctx.SetBody([]byte("Hello World"))
	})

	app.Run(8080)
}
