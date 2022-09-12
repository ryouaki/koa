package koa

import (
	"testing"

	"github.com/issue9/assert"
)

func TestComparea(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(true, compare("/", "/"))
	assert.Equal(false, compare("/a", "/"))
	assert.Equal(true, compare("/*", "/"))

	assert.Equal(true, compare("/*", "/1/2/3"))
	assert.Equal(true, compare("/1/*", "/1/2"))
	assert.Equal(true, compare("/1/:a", "/1/2"))

	assert.Equal(false, compare("/1/:a", "/1/2/1"))
	assert.Equal(true, compare("/1/:a/1", "/1/2/1"))
	assert.Equal(true, compare("/1/:a/*", "/1/2/1"))
	assert.Equal(true, compare("/1/:a/*", "/1/2/1"))
}
