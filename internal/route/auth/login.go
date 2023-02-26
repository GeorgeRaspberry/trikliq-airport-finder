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

func Login(ctx *gin.Context) {
	var (
		loginRequest    = model.Request{}
		loginResponse   = model.Response{}
		eMail, password string
		ok              bool
	)

	requestID, _ := ctx.Get("id")

	// Parse the posted JSON data.
	err := ctx.ShouldBind(&loginRequest)
	if err != nil {
		if err != nil {
			logger.Log.Error("Failed to bind input data",
				zap.Error(err),
			)

			error.ReturnError(ctx, loginResponse, error.SystemError(requestID), logger.Log)
			return
		}
	}

	log := logger.Log.WithOptions(zap.Fields(
		zap.Any("data", loginRequest.Data),
	))

	log.Info("login started")

	if eMail, ok = loginRequest.Data["email"].(string); !ok {
		err := fmt.Errorf("email is missing")
		log.Error("auth/login failed")

		error.ReturnError(ctx, loginResponse, err.Error(), log)
		return
	}
	if password, ok = loginRequest.Data["password"].(string); !ok {
		err := fmt.Errorf("password is missing")
		log.Error("auth/login failed")

		error.ReturnError(ctx, loginResponse, err.Error(), log)
		return
	}

	var user model.User
	// Look for the provided user.
	database.DB.Where("email = ?", eMail).First(&user)

	// User not found.
	if user.ID == 0 {
		err := fmt.Errorf("user not found")

		log.Error("auth/login failed",
			zap.Error(err),
		)

		error.ReturnError(ctx, loginResponse, err.Error(), log)
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

		error.ReturnError(ctx, loginResponse, err.Error(), log)
		return
	}

	_, err = GetIssuer(ctx)
	if err != nil {
		log.Error("Login failed",
			zap.Error(err),
		)

		error.ReturnError(ctx, loginResponse, err.Error(), log)
		return
	}

	log.Info("login finished")

	loginResponse.Status = true
	ctx.JSON(200, loginResponse)
}

func init() {
	router.Router.Handle("POST", "auth/login", Login)
}
