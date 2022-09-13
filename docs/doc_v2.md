# Go版本Koa的1.0 到 2.0 — 不再只是单纯的模仿

> 1年前，我还在滴滴外卖的时候，在学习`Go`期间，尝试实现了`Go`版本的`Koa`。也在一些小的场景中落地了，效果很不错，非常适合`Nodejs`开发者。虽然我的版本完整实现了`Koajs`的最具特色的功能 — `洋葱模型`和`中间件系统`，但是也依稀存在一些不足。毕竟，只是出于兴趣和对学习的总结`卷`出来的轮子，仅此而已，所以并没有过多在意。

### 是Koajs还是Koa？
每一个写过`Nodejs`的人都对`Koajs`并不陌生，`Koajs`的洋葱模型和中间件系统基本上是必知必会的。那么`Koajs`仅仅有洋葱模型么？那肯定会想起来相对于`Expressjs`的极简。

当然这种极简的`Koajs`过于极端，自身除了洋葱模型和中间件系统其它的都被舍弃了，必须以中间件的形式加回来。正是由于这种特性，我们很容易通过中间件机制以及洋葱模型实现非常复杂的`过滤`效果。可以很容易通过中间件实现`解耦`，`隔离`，`复用`。每一个逻辑都可以做到很`纯粹`。

但是也反而带来一些问题：它是一个`Web`框架。那么一个`Web`极简框架也要具备最基本的`路由`能力吧？由于`Koajs`是把路由对象当做中间件处理的，路由的能力都集中在一个中间件里面，也就导致使用了`router`模块以后性能下降非常明显。当然了这在集群模式下还是可以忽略不计的。

我还是很赞同`Koajs`的思路，如果路由能力整合到`Koajs`中，那么聚合中间件的`compose`就要耦合路径匹配能力。导致`compose`不再纯粹。对于整个中间件系统是侵入式的。因此我们能够在`Koajs`身上看到非常纯粹框架。但是我们也不得不面对一些麻烦，就是反复的加各种中间件，比如这个不可或缺的`router`。

所以在`V1.x`的时候，我就直接把路由能力集成进去了。但是为了解耦，只能在响应处理过程中进行路径匹配筛选出符合条件的中间件和路由函数，再进行`compose`操作。这直接导致一个问题就是 — 需要手动结束也就是需要调`ctx.Write`接口。这很不`Koa`。于是乎`V2.x`肯定要提上日程的。

### 造轮子应该从使用者角度触发
`Koajs`不得不说是一个非常巧妙的框架，但是在实际使用中`99%`的情况下我们依然需要把`router`模块加回来。所以操作貌似这有点多余。内置`router`模式依然是必须解决的问题。当然既然已经定位是一个`Web`框架了。所以我也`看破红尘`了。并不需要像`Koajs`做到那么纯粹。

我的目标是实现一个可以像`Koa`风格那样使用`Go`语言开发`Web`的框架。就像当年我创造`mkbugjs`一样只是希望`Java`开发者更容易接受`Nodejs`开发，虽然后面出现的`nestjs`做的更好，但是在那个空白期我做到了。何不继续挑战一下？

### 重构compose
在`V1.x`中我犯的最大的错误就是完全参考`Koajs`的实现去抄袭。所以为了实现洋葱模型也就导致了中间件的定义和路由的定义是不同的。必须区别对待。但是在`V2.x`中。为了更好的使路由能力继承进来。对路由和中间件的定义进行了抹平处理。使用同一种定义：
```go
// 应用实例
type Application struct {
  handles []Handle // 存储中间件和路由的堆栈
  _cb     Handler // 中间件和路由聚合后的响应函数
}

// 中间件结构
type Handle struct {
  path     string // 注册的url，无论是中间件还是路由，都可以设置地址
  method  string // 方法，中间件默认为 use
  handler Handler // 响应体函数
}

// 响应体函数结构定义
type Handler func(ctx *Context, n Next)

// 同koa的next
type Next func()

// 方法定义
const (
  GET     = "get"
  POST    = "post"
  PUT     = "put"
  DELETE  = "delete"
  PATCH   = "patch"
  OPTIONS = "options"
  HEAD    = "head"
  USE     = "use"
)
```
这样，将所有的路由和中间件都视为统一的接口，通过`Next`去决定是否透传，当然这很`Koa`。但是这也就导致了`compose`的不同，我需要在`compose`操作的同时进行路径匹配:
```go
// 由于Go同样具备闭包的能力，因此可以实现高阶函数，这也是能够完整复刻Koajs语法的必要条件
func compose(handles []Handle) Handler {
  return func(ctx *Context, n Next) { // 返回一个高阶函数，用于作为请求响应回调
    _curr := int32(0) // 计数器，不能超过最大层数
    _max := len(handles) // 设置的中间件和路由的个数
    var _next Next // 下传回调

    _next = func() { // 为了起到一个reduce的效果，需要无限回调自己
      _ctx := ctx // 通过闭包记住上下文
      _handles := handles // 通过闭包记住中间件和路由栈

      if int(_curr) < _max { // 不能超过最大堆栈数量
        _currHandler := _handles[_curr] // 取出当前需要处理的中间件或者路由
        _curr += 1 // 堆栈前移，在下一次回调处理下一个
        if (_currHandler.method == USE || _currHandler.method == ctx.Method) &&
          compare(_currHandler.path, ctx.Path) { // 中间件和方法都需要执行handler，但是url为”*”或者””的必须执行
          ctx.Params = formatParams(_currHandler.path, ctx.Path) // 处理URL传参
          _currHandler.handler(_ctx, _next) // 执行匹配的中间件或者路由
        } else {
          _next() // 如果没有匹配的中间件或者路由则进行下一次自调用
        }
      }
    }
    _next()
  }
}
```
这里我移除了在`V1.x`版本的`err`参数。在使用的时候，既然都发生错误了。还`next`传递`error`干什么呢？哪发生哪处理就好了。所以之前为了实现`Koa`设计了太多实际工作当中毫无用处的东西。

_*不能为了造轮子而造轮子，一定要关注我们造的轮子要解决什么问题，用户希望有什么样的一个轮子去解决这个问题。而不仅仅是我觉得这个轮子可以这样解决这个问题。要以最终用户角度出发去解决问题。*_

> 这个处理看起来是非常低效的。但是目前`V2.x`在`mac air m1`平台在只返回`hello world`的情况下轻轻松松完成了`22w`的`qps`。后端绝大多数情况下，单机的负载是`200-400/分`左右,按峰值三倍计算也就是`600-1200/分`左右，当然也可能更高一些。所以这个性能对于任何后端都是严重溢出的。影响都是微乎其微的，所以也就可以忽略不计了。

最终就可以很容易的实现我们期望的效果，让`Go`也有了`Nodejs`的身影，有了`Koajs`的影子：
```go
package main

import (
  "fmt"

  "github.com/ryouaki/koa"
)

func main() {
  app := koa.New()

  handler1 := func(ctx *koa.Context, n koa.Next) {
    fmt.Println("handler1")
    n()
    fmt.Println("handler1")
  }

  handler2 := func(ctx *koa.Context, n koa.Next) {
    fmt.Println("handler2")
    ctx.SetBody([]byte("Hello world"))
    fmt.Println("handler2")
  }

  app.Get("/test", handler1, handler2)

  app.Run(8080)
}
```
测试结果：
```go
  $ curl localhost:8080/test   // Hello world

  handler1
  handler2
  handler2
  handler1
```

### 离不开的“模式”
实现洋葱模型的核心是中间件系统，而实现中间件系统使用微内核架构模式再适合不过了。在这里依然使用解释器模式实现了微内核架构模式。这里有之前的文章专门用于介绍微内核架构[微内核架构模式在前端的实践与原则](https://zhuanlan.zhihu.com/p/443982576)感兴趣可以去了解一下。

### 最后:卷非卷，似卷非卷
实现完`V2.x`之后让我非常兴奋，这是多么无聊的一晚啊。既然已经存在了`Koajs`为什么又要搞一版`Go`版本的`Koa`呢？但是谁让互联网人就是这么“卷”呢？

造轮子无非为了`KPI`，为了知名度，为了影响力，为了更好的`Offer`，或者为了满足自己虚荣心？当然，也可能像我一样仅仅是为了喜欢写代码，作为自己学习`Go`语言过程中的一个阶段性总结。

[源码Go版Koa](https://github.com/ryouaki/koa)