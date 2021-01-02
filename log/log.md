# Koa.Log 介绍
一个基本的日志库，支持控制台输出与文件输出，并支持最大保留天数以及多种日志分类存储。

## 如何使用Koa.Log处理日志信息呢？
使用起来非常方便
```go
    package main

    import (
        "errors"
        "fmt"

        "github.com/ryouaki/koa"
        "github.com/ryouaki/koa/log"
    )

    func main() {
        app := koa.New()

        // 初始化日志系统
        log.New(&log.Config{
            Level:   log.LevelInfo,     // 日志级别info
            Mode:    log.LogFile,       // 以文件形式保存日志
            MaxDays: 1,                 // 最多保留一天
            LogPath: "./../logs",       // 日志文件存储位置
        })
        
        app.Use("/", func(err error, ctx *koa.Context, next koa.NextCb) {
            log.Info("Request in") // 打印日志
            next(err)
            log.Info("Request out")// 打印日志
        })

        app.Get("/test/:var/p", func(err error, ctx *koa.Context, next koa.NextCb) {
            ctx.Write([]byte("Hello World"))
        })

        err := app.Run(8080)
        if err != nil {
            fmt.Println(err)
        }
    }
```