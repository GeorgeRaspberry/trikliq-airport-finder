package auth

import (
	"bookbox-backend/internal/model"
	"bookbox-backend/internal/route/error"
	"bookbox-backend/internal/server/router"
	"bookbox-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Refresh(ctx *gin.Context) {

	var (
		refreshResponse = model.Response{}
	)

	log := logger.Log.WithOptions(zap.Fields())

	log.Info("refresh started")

	issuer, err := RefreshToken(ctx)
	if err != nil {
		log.Error("Failed to validate token token",
			zap.Error(err),
		)

		error.ReturnError(ctx, refreshResponse, err.Error(), log)
		return
	}

	// Refresh tokens.
	err = NewAccessToken(ctx, issuer.ID)
	if err != nil {
		log.Error("Failed to refresh token",
			zap.Error(err),
		)

		error.ReturnError(ctx, refreshResponse, err.Error(), log)
		return
	}

	log.Info("refresh finished")

	refreshResponse.Status = true
	ctx.JSON(200, refreshResponse)
}

func init() {
	router.Router.Handle("POST", "auth/refresh", Refresh)
}
