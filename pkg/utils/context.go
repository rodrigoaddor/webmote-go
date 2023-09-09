package utils

import (
	"github.com/gin-gonic/gin"
)

func Get[K any](ctx *gin.Context, key string) (*K, bool) {
	maybeValue, ok := ctx.Get(key)
	if !ok {
		return nil, false
	}

	value, ok := maybeValue.(K)
	if !ok {
		return nil, false
	}

	return &value, true
}
