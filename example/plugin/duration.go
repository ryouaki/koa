package plugin

import (
	"fmt"
	"time"

	"koa.go"
)

// Duration func
func Duration(err error, ctx *koa.Context, next koa.NextCb) {
	startTime := time.Now()
	next(nil)
	d := time.Now().Sub(startTime)
	fmt.Println("Request cost: ", float64(d)/float64(time.Millisecond), "ms")
}
