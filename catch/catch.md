> Golang虽然也有捕获异常的机制，但是却不是很人性化，很容易忽略对异常的处理。从而导致很多问题。

经验丰富的Go工程师会选用defer的方式捕获异常。但是这种碎片式的使用方式难免会导致疏忽。而传统的try - catch -finally的模式因为在一个执行片段内，并且易于维护，非常符合大众的开发习惯。作为一个不正经的前端工程师。又开始手痒痒了。因为我是前端，所以我习惯于写Promise处理异常。所以。。。。。Koa.Catch库就诞生了

## 如何用Koa.Catch安全的处理这些意外呢？
使用起来非常简单。有点像Nodejs的Promise。
```go
    package main

    import (
        "errors"
        "fmt"
        catch "github.com/ryouaki/koa/catch" 
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

## 如何实现的呢？这里是源代码
其实很简单，golang的defer语法可以保证在函数结束的时候一定会执行。所以，我要做的就是在回调函数里面使用defer然后进行捕获异常。

以下是代码：
```go
    package catch

    const (
	    PENDING   int = 0
	    FULFILLED int = 1
	    REJECTED  int = 2
	    FINALLY   int = 3
    )

    // 为了像Promise一样使用的顺手，肯定要有状态和结果了。这样才更Promise
    type Exec struct { 
	    status int
	    ret    interface{}
    }

    // Promise必备的3个接口呀。当然还有All什么的。但是Go中用不上。就算了
    type Exception interface { 
	    Then(func(interface{})) Exception
	    Catch(func(err interface{})) Exception
	    Finally(func()) Exception
    }

    // 这是一个非常关键的函数，用于捕获异常信息。因为try，catch，then中都要用。所以就抽离出来了。
    func (exec *Exec) deferHelper(cb func()) {
	    defer func() {
		    defer func() {
			    r := recover()
			    if r != nil && exec.status != FINALLY {
				    exec.status = REJECTED
				    exec.ret = r

			    }
		    }()

		    if exec.status == PENDING {
			    exec.status = FULFILLED
		    }
		    cb()
	    }()
    }

    // Try函数，主要是用于包含执行体，并初始化Catch的状态，如果不发生意外的话，就在deferHelper中把状态置成FULFILLED，表示成功执行结束，如果发生了意外，就是REJECTED
    func Try(cb func() interface{}) Exception {
	    exec := &Exec{
		    status: PENDING,
	    }

	    exec.deferHelper(func() {
		    exec.ret = cb()
	    })

	    return exec
    }

    // 如果执行成功了。就可以通过Then回调拿到结果。为什么要这么搞，只是为了希望写起来像Promise (*_*)
    func (exec *Exec) Then(cb func(interface{})) Exception {
	    if exec.status == FULFILLED {
		    exec.deferHelper(func() {
			    cb(exec.ret)
		    })
	    }
	    return exec
    }

    // 如果如果发生了意外，那么就是Catch了。
    func (exec *Exec) Catch(cb func(err interface{})) Exception {
	    if exec.status == REJECTED {
		    exec.deferHelper(func() {
			    cb(exec.ret)
		    })
	    }
	    return exec
    }

    // 当然，无论你发不发生意外，是否执行成功，最终都会执行Finally。这也是这个库的最终目标
    func (exec *Exec) Finally(cb func()) Exception {
	    if exec.status != FINALLY {
		    exec.status = FINALLY
		    exec.deferHelper(cb)
	    }

	return exec
    }
```

## 最后

哈哈哈哈，是不是很好玩，把Golang搞成Nodejs的样子，这就是[Koa.go](https://github.com/ryouaki/koa)的使命，目前各种插件，中间件在逐步开发中。由于只有周末，进度比较慢。有没有小伙伴愿意加入把它打造成Golang届的Egg.go