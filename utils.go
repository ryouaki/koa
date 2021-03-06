package koa

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

// Bytes2String func
func Bytes2String(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

// compare path middleware prefix, target request path
func compare(path string, target string, isRouter bool) bool {
	if len(path) == 0 || path == "/" {
		return true
	}

	pathArr := strings.Split(path, "/")
	targetArr := strings.Split(target, "/")
	pathLen := len(pathArr)
	targetLen := len(targetArr)

	if pathLen > targetLen {
		return false
	}

	if isRouter && pathLen != targetLen {
		return false
	}

	for idx, val := range pathArr {
		if val != targetArr[idx] {
			if !strings.HasPrefix(val, ":") {
				return false
			}

			variable := strings.TrimPrefix(val, ":")
			if matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", variable); !matched {
				return false
			}
		}
	}

	return true
}

func compose(ctx *Context, handlers []Handler) func(error) {
	currentPoint := int32(0)
	var next func(err error)
	var cbMax = len(handlers)

	next = func(err error) {
		_ctx := ctx
		_router := handlers
		bFound := false

		for int(currentPoint) < cbMax && bFound == false {
			bFound = true
			currRouterHandler := _router[currentPoint]
			atomic.AddInt32(&currentPoint, 1)
			currRouterHandler(err, _ctx, next)
		}
	}

	return next
}

func formatQuery(values map[string]([]string)) map[string]([]string) {
	result := make(map[string]([]string))
	for key, data := range values {
		if strings.HasSuffix(key, "[]") {
			key = strings.Replace(key, "[]", "", -1)
		}
		result[key] = data
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

// GetIPAddr func
func GetIPAddr() string {
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

// GetMD5ID func
func GetMD5ID(b []byte) string {
	res := md5.Sum(b)
	return hex.EncodeToString(res[:])
}

// StructAssign func
func StructAssign(binding interface{}, value interface{}) {
	bVal := reflect.ValueOf(binding).Elem()
	vVal := reflect.ValueOf(value).Elem()
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		name := vTypeOfT.Field(i).Name
		if ok := bVal.FieldByName(name).IsValid(); ok {
			bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
		}
	}
}
