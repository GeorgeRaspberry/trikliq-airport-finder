package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Session() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := fmt.Sprintf("%v", time.Now().UnixNano())
		ctx.Set("id", id)
		ctx.Next()
	}
}
