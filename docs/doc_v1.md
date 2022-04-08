> 从我的工作经历就能够看出我注定是一个不正经非主流前端工程师，从一开始的`c`语言写网络应用，用`Java`开发`Swing`，好不容易赶上移动浪潮去搞了1年`android`却通过了`ios`开发工程师认证但是每天在写`js`开发后端接口。然后就一直挂着前端工程师的`title`写`java`，写`node`，写`shell`。来滴滴后，发现大家都在写`golang`，于是乎，这条筋又不对了。开始研究了`golang`，俗话说，知己知彼百战百胜，只有了解了`RD`的武器，才能更好地（撕）合作。

`Koa.js`是`Nodejs`技术一个非常注明的框架，由`JS`大神`TJ`创建，其中洋葱模型和`co`在`nodejs4.0` - `nodejs10`解决了大家非常多的困扰。虽然我不喜欢`koa.js`。但是也非常喜欢洋葱模型，`co`以及它的中间件机制。虽然在`nodejs10`+以后`co`已经没什么卵用了。

### 什么是洋葱模型呢？
![](example/static/2892151181-5ab48de7b5013_articlex.png)

说白了就是请求进来，一层一层的通过中间件执行`next`函数进入到你设置的下一个中间件中，并且可以通过`context`对象一直向下传递下去，当到达最后一个中间件的时候，又向上返回到最初的地方。是不是有点像`dom`事件的捕获冒泡？确实是这个样子。

有人说，不就是递归吗，这有什么？确实，从原理角度来看确实是递归但是又有很多不同：

- 由于我们需要通过第三参数`next`进行下一个回调的执行，并且在`next`调用的下一个回调执行后，还能继续执行剩余代码，并且可以在任意位置执行下一次调用的函数。所以这个实现就不仅仅是递归那么简单了。因为执行的不是同一个函数，且位置也不同。
- 传参和引用，我们需要保障整个执行过程，各个中间执行的函数中`context`指向的是同一个上下文对象。

为了实现以上2个功能，首先语言必须具备1个基本能力 - 闭包，`golang`因为具备闭包，因此从理论上来讲是可以实现洋葱模型的。实际上是否可行呢？确实也是可行的。

### `koa`一个实现了洋葱模型和中间件机制的`golang`框架
我就把它叫做`koa`吧。目前大多数`golang`的框架都是支持中间件的，但是中间件都是只支持顺序执行，并不支持洋葱模型。所以。这个也算是一个挑战吧。

话不多说先看看执行效果
```go
package main

import (
	"fmt"

	"github.com/ryouaki/koa"
	"github.com/ryouaki/koa/example/plugin"
)

func main() {
	app := koa.New()

	app.Use(plugin.Duration)
	app.Use("/", func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("test1")
		next(err)
		fmt.Println("test1")
	})

	app.Get("/test/:var/p", func(err error, ctx *koa.Context, next koa.NextCb) {
		fmt.Println("test", ctx.Params)
		ctx.SetData("test", ctx.Query["c"][0])
		next(nil)
	}, func(err error, ctx *koa.Context, next koa.NextCb) {
		ctx.Write([]byte(ctx.GetData("test").(string)))
	})

	err := app.Run(8080)
	if err != nil {
		fmt.Println(err)
	}
}
```
这段代码包含一个中间件`Duration`用于输出请求耗时，`test1`中间件是观察洋葱模型的执行顺序，`test/:var/p`是我们的测试接口
```sh
    $ curl localhost:8080/test/test/p?c=Hello World

    test1
    test map[var:test]
    test1
    2020-12-17 15:04:28 +0800 CST /test/test/p?c=Hello%20World Request cost:  0.040756 ms
```
通过日志我们可以看到，我们的`test/:var/p`接口函数像嵌套在`test1`中间件`next`的位置一样被执行了。那么是如何实现的呢？是不是很好玩？

### 洋葱模型的核心代码
整个`koa`核心代码不足300行，其中最核心的代码只有20行，也就是实现洋葱模型的核心代码：
```go
func compose(ctx *Context, handlers []Handler) func(error) {
	currentPoint := int32(0) // 记录当前执行到第几个函数
	var next func(err error)    // 缓存退栈函数，
	var cbMax = len(handlers)    // 记录最大执行次数

        // 整个框架最核心的代码，退栈执行函数
	next = func(err error) {
		_ctx := ctx    // 闭包缓存当前执行上下文指针
		_router := handlers    // 缓存需要执行的函数
		bFound := false     // 当前轮训是否匹配到已执行函数。保证每次只执行一个函数

		for int(currentPoint) < cbMax && bFound == false {
			bFound = true
                        // 取出当前需要执行的函数
			currRouterHandler := _router[currentPoint] 
			atomic.AddInt32(¤tPoint, 1)
                        // 执行，并将退栈执行函数传入。
			currRouterHandler(err, _ctx, next) 
		}
	}
        // 返回退栈函数
	return next
}
```
#### 解读如何使用`compose`实现洋葱模型
为了保障各个中间件的执行顺序，并且保障执行上下文的一致，我们需要将所有符合请求路由规则的`Handler`进行压栈。
```go
// ServeHTTP interface func
func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var err error = nil
	var body []uint8 = nil

	method := strings.ToLower(req.Method)

	if req.Body != nil {
		body, err = ioutil.ReadAll(req.Body)
	}
        // 创建一个新的执行上下文
	ctx := &Context{
		Header:   req.Header,
		Res:      w,
		Req:      req,
		URL:      req.RequestURI,
		Path:     req.URL.Path,
		Query:    formatQuery(req.URL.Query()),
		Body:     body,
		Method:   method,
		Status:   200,
		IsFinish: false,
		data:     make(map[string]interface{}),
	}

	var routerHandler []Handler
        // 将符合规则的中间件进行压栈
	for _, middleware := range app.middlewares {
		if ok := compare(middleware.path, ctx.Path, false); ok {
			routerHandler = append(routerHandler, middleware.handler)
		}
	}
        // 将符合规则的路由函数进行压栈
	for _, router := range app.route[ctx.Method] {
		if ok := compare(router.path, ctx.Path, true); ok {
			ctx.MatchURL = router.path
			ctx.Params = formatParams(router.path, ctx.Path)
			routerHandler = append(routerHandler, router.handler...)
		}
	}
        // 开始执行因为需要保障每一个中间件和路由函数访问同一个执行上下文，所以传入的是指针。
	fb := compose(ctx, routerHandler)
	fb(err)
}
```    

压栈后的函数，统一推入`compose`中，进行退栈执行。就是这么简单。

### 总结
通过洋葱模型，使我们可以更容易的监控整个路由请求的全链路，并且可以在中间件做更多的事情。比如统计响应耗时：
```go
package plugin

import (
	"fmt"
	"time"

	"github.com/ryouaki/koa"
)

// Duration func
func Duration(err error, ctx *koa.Context, next koa.NextCb) {
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
```

可以看到，在洋葱模型的支持下，我们可以将整个中间件逻辑写在一起即可。而next就是具体的路由逻辑，我们并不需要关心它具体做了什么，只需要关心我们自己需要做什么。这样可以将同一个业务逻辑的代码聚合到一起，更容易维护，不是么？

## 最后 
[求一波follow和star 项目地址](https://github.com/ryouaki/koa)