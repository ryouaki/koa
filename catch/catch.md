# Koa.Catch 介绍
Golang本身不支持异常处理，但是在日常开发中，难免写出一些带有错误的代码。合并panic的时候更是如此，如果没有异常捕获，会导致资源无法回收等问题，严重甚至导致死锁。

所以我们要保障我的代码可以有一个安全的地方回收资源。如果可以try-catch-finally那么就容易的多了。但是Golang却不这么做。因此也带来很多麻烦。

## 如何用Koa.Catch安全的处理这些意外呢？
使用起来非常简单。有点像Nodejs的Promise。
```go
    package main

    import (
        "errors"
        "fmt"
        catch "github.com/mkbug-com/pc-api/koa" // 暂时在这里，稍后会移出去。
    )

    func main() {
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
```

运行结果如下：
```go
    $ go run main.go

    test1
    test error
    test error 2
    end
```