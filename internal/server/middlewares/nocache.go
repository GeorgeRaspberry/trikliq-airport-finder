package middlewares

import (
	"github.com/gin-gonic/gin"
)

// NoCache disables caching
func NoCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer middlewareRecovery()

		ctx.Writer.Header().Set("expires", "Sat, 01 Jan 2000 00:00:00 GMT")
		ctx.Writer.Header().Set("cache-control", "private, no-cache, no-store, must-revalidate")

		ctx.Next()
	}
}
