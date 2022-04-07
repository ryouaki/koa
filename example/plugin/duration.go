package plugin

import (
	"fmt"
	"time"

	"github.com/ryouaki/koa"
)

// Duration func
func Duration(err error, ctx *koa.Context, next koa.Next) {
	startTime := time.Now()
	next(nil)
	d := time.Now().Sub(startTime)
	fmt.Println(time.Date(startTime.Year(),
		startTime.Month(),
		startTime.Day(),
		startTime.Hour(),
		startTime.Minute(),
		startTime.Second(), 0, time.Local),
		ctx.URL,
		"Request cost: ",
		float64(d)/float64(time.Millisecond), "ms")
}
