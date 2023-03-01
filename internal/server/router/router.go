package router

import (
	"trikliq-airport-finder/internal/server/middlewares"

	"github.com/gin-gonic/gin"
)

var (
	Router = gin.New()
)

func init() {
	Router.Use(middlewares.NoCache())
	Router.Use(middlewares.Session())
	Router.Use(middlewares.CORS())
	Router.Use(middlewares.Security())

}
