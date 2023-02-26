package middlewares

import (
	"github.com/gin-gonic/gin"
)

// Security will inject HTTP headers related to security
func Security() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer middlewareRecovery()

		ctx.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		ctx.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
		ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		ctx.Writer.Header().Set("Referrer-Policy", "same-origin")
		ctx.Writer.Header().Set("X-Xss-Protection", "1; mode=block")

		ctx.Next()
	}
}
