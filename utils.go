package koa

import (
	"bytes"
	"net"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// compare path middleware prefix, target request path
func compare(path string, target string) bool {
	if len(path) == 0 || path == "*" {
		return true
	}

	pathArr := strings.Split(path, "/")
	targetArr := strings.Split(target, "/")
	pathLen := len(pathArr)
	targetLen := len(targetArr)

	if pathLen > targetLen {
		return false
	}

	reg, _ := regexp.Compile("^[a-zA-Z0-9_-]+$")

	for idx, val := range pathArr {
		if val != targetArr[idx] {
			if !strings.HasPrefix(val, ":") {
				return false
			}

			variable := strings.TrimPrefix(val, ":")
			if !reg.MatchString(variable) {
				return false
			}
		}
	}

	return true
}

func compose(handles []Handle) Handler {
	return func(ctx *Context, n Next) {
		_curr := int32(0)
		_max := len(handles)
		var _next Next

		_next = func() {
			_ctx := ctx
			_handles := handles

			if int(_curr) < _max {
				_currHandler := _handles[_curr]
				_curr += 1
				if (_currHandler.method == USE || _currHandler.method == ctx.Method) &&
					compare(_currHandler.path, ctx.Path) {
					ctx.Params = formatParams(_currHandler.path, ctx.Path)
					_currHandler.handler(_ctx, _next)
				} else {
					_next()
				}
			}
		}
		_next()
	}
}

func formatQuery(values map[string]([]string)) map[string]interface{} {
	result := make(map[string]interface{})
	for key, data := range values {
		if strings.HasSuffix(key, "[]") {
			key = strings.Replace(key, "[]", "", -1)
			result[key] = data
		} else {
			result[key] = data[0]
		}
	}
	return result
}

func formatParams(path string, target string) map[string]string {
	result := make(map[string]string)
	pathArr := strings.Split(path, "/")
	targetArr := strings.Split(target, "/")

	for idx, val := range pathArr {
		if val != targetArr[idx] {
			variable := strings.TrimPrefix(val, ":")
			result[variable] = targetArr[idx]
		}
	}

	return result
}

// GetLocalAddrIp func
func GetLocalAddrIp() string {
	ret := ""

	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, value := range addrs {
			if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ret += ipnet.IP.String()
				}
			}
		}
	}

	return ret
}

// GetGoroutineID func
func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
