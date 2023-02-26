package auth

import (
	"bookbox-backend/internal/model"
	"bookbox-backend/internal/route/error"
	"bookbox-backend/internal/server/router"
	"bookbox-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logout(ctx *gin.Context) {

	var (
		logoutResponse = model.Response{}
	)

	log := logger.Log.WithOptions(zap.Fields())

	log.Info("logout started")

	issuer, err := GetIssuer(ctx)
	if err != nil {
		log.Error("No issuer found",
			zap.Error(err),
		)

		error.ReturnError(ctx, logoutResponse, err.Error(), log)
		return
	}

	err = DestroyToken(ctx, issuer.ID)
	if err != nil {
		log.Error("Failed to logout user",
			zap.Error(err),
		)

		error.ReturnError(ctx, logoutResponse, err.Error(), log)
		return
	}

	log.Info("logout finished")

	logoutResponse.Status = true
	ctx.JSON(200, logoutResponse)
}

func init() {
	router.Router.Handle("POST", "auth/logout", Logout)
}
