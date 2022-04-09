package plugin

import (
	"fmt"
	"time"

	"github.com/ryouaki/koa"
)

// Duration func
func Duration(ctx *koa.Context, next koa.Next) {
	startTime := time.Now()
	next()
	d := time.Since(startTime)
	fmt.Println(time.Date(startTime.Year(),
		startTime.Month(),
		startTime.Day(),
		startTime.Hour(),
		startTime.Minute(),
		startTime.Second(), 0, time.Local),
		ctx.Url,
		"Request cost: ",
		float64(d)/float64(time.Millisecond), "ms")
}
