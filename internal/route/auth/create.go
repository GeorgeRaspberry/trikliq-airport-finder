package auth

import (
	"bookbox-backend/internal/database"
	"bookbox-backend/internal/model"
	"bookbox-backend/internal/route/error"
	"bookbox-backend/internal/server/router"
	"bookbox-backend/pkg/logger"
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func Create(ctx *gin.Context) {
	var (
		authRequest     = model.Request{}
		authResponse    = model.Response{}
		eMail, password string
		ok              bool
	)

	requestID, _ := ctx.Get("id")

	// Parse the posted JSON data.
	err := ctx.ShouldBind(&authRequest)
	if err != nil {
		if err != nil {
			logger.Log.Error("Failed to bind input data",
				zap.Error(err),
			)

			error.ReturnError(ctx, authResponse, error.SystemError(requestID), logger.Log)
			return
		}
	}

	log := logger.Log.WithOptions(zap.Fields(
		zap.Any("data", authRequest.Data),
	))

	log.Info("create started")

	if eMail, ok = authRequest.Data["email"].(string); !ok {
		err := fmt.Errorf("email is missing")
		log.Error("auth/create failed")

		error.ReturnError(ctx, authResponse, err.Error(), log)
		return
	}
	if password, ok = authRequest.Data["password"].(string); !ok {
		err := fmt.Errorf("password is missing")
		log.Error("auth/create failed")

		error.ReturnError(ctx, authResponse, err.Error(), log)
		return
	}

	var user model.User
	// Look for the provided user.
	database.DB.Where("email = ?", eMail).First(&user)

	// User not found.
	if user.ID == 0 {
		err := fmt.Errorf("user not found")

		log.Error("auth/create failed",
			zap.Error(err),
		)

		error.ReturnError(ctx, authResponse, err.Error(), log)
		return
	}

	passDecoded, _ := base64.StdEncoding.DecodeString(user.Password)

	// Compare the passwords.
	err = bcrypt.CompareHashAndPassword(
		passDecoded,
		[]byte(password),
	)
	if err != nil {
		log.Error("Incorrect credentials",
			zap.Error(err),
		)

		error.ReturnError(ctx, authResponse, err.Error(), log)
		return
	}

	// Generate JWT pair.
	err = NewTokenPair(ctx, user.ID)
	if err != nil {
		log.Error("Failed to generate token",
			zap.Error(err),
		)

		error.ReturnError(ctx, authResponse, error.SystemError(requestID), log)
		return
	}

	log.Info("create finished")

	authResponse.Status = true
	ctx.JSON(200, authResponse)
}

func init() {
	router.Router.Handle("POST", "auth/create", Create)
}
