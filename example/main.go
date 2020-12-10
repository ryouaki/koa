package main

import (
	"fmt"

	"mkbug.go"
)

func main() {
	app := mkbug.New()

	app.Use("/", func(err error, ctx *mkbug.Context, next mkbug.NextCb) {
		fmt.Println("test1")
		next(err)
		fmt.Println("test1")
	})
	app.Use(func(err error, ctx *mkbug.Context, next mkbug.NextCb) {
		fmt.Println("test1")
		next(err)
		fmt.Println("test1")
	})

	app.Use(func(err error, ctx *mkbug.Context, next mkbug.NextCb) {
		fmt.Println("test2")
		next(nil)
		fmt.Println("test2")
	})
	err := app.Run(8080)
	if err != nil {
		fmt.Println(err)
	}
}
