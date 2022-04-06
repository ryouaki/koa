package static

import (
	"net/http"
	"path/filepath"
	"strings"

	koa "github.com/ryouaki/koa"
)

// Static func
// Params path // 路径
// Params prefix // 访问前缀
func Static(path string, prefix string) func(err error, ctx *koa.Context, next koa.NextCb) {
	_path := path
	_prefix := prefix

	return func(err error, ctx *koa.Context, next koa.NextCb) {
		_p := ctx.Path
		if !strings.HasPrefix(_p, _prefix) {
			next(nil)
			return
		}

		_staticFilePath := strings.Replace(_p, _prefix, "", len(_prefix))
		_staticFile, _ := filepath.Abs(_path)

		http.ServeFile(ctx.Res, ctx.Req, filepath.Join(_staticFile, _staticFilePath))
	}
}
