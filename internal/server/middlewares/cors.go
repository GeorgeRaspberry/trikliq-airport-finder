package middlewares

import (
	"github.com/gin-gonic/gin"
)

var (
	Origin = ""
)

// CORS will inject HTTP header Access-Control-Allow-Origin
func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer middlewareRecovery()

		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Next()
	}
}
