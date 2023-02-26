package auth

import (
	"bookbox-backend/internal/config"
	"bookbox-backend/internal/database"
	"bookbox-backend/internal/model"
	"bookbox-backend/pkg/logger"
	"bookbox-backend/pkg/redis"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func GetIssuer(ctx *gin.Context) (issuer *model.User, err error) {
	return ValidateToken(ctx)
}

func ValidateToken(ctx *gin.Context) (issuer *model.User, err error) {
	tokenToValidate, err := ctx.Cookie("access_token")
	if err != nil {
		return RefreshToken(ctx)
	}

	return ValidateTokenClaims(ctx, tokenToValidate, false)
}

func ValidateTokenClaims(
	ctx *gin.Context,
	tokenToValidate string,
	isRefresh bool,
) (*model.User, error) {

	parsedToken, err := jwt.ParseWithClaims(
		tokenToValidate,
		&jwt.RegisteredClaims{},
		func(jwt *jwt.Token) (interface{}, error) {
			if isRefresh {
				return config.JWT.RefreshPrivateKey.Public(), nil
			} else {
				return config.JWT.AccessPrivateKey.Public(), nil
			}
		},
	)
	if err != nil {
		logger.Log.Error("Failed to bind input data",
			zap.Error(err),
		)
		return nil, err
	}

	claims := parsedToken.Claims.(*jwt.RegisteredClaims)

	cacheJSON, err := redis.Client.Get(fmt.Sprintf("token_pair-%s", claims.Issuer)).
		Result()
	if err != nil {
		logger.Log.Error("Failed to get data from Redis",
			zap.Error(err),
		)
		return nil, err
	}

	cachedTokenPair := TokenPairCache{}

	err = json.Unmarshal([]byte(cacheJSON), &cachedTokenPair)
	if err != nil {
		return nil, err
	}

	var tokenUID string
	if isRefresh {
		tokenUID = cachedTokenPair.RefreshTokenUUID
	} else {
		tokenUID = cachedTokenPair.AccessTokenUUID
	}

	if tokenUID != claims.ID {
		userID, err := strconv.Atoi(claims.Issuer)
		if err != nil {
			return nil, err
		}

		err = DestroyToken(ctx, uint(userID))
		if err != nil {
			logger.Log.Error("This token most likely stolen and refreshed by another user. User will be logged out for security reasons.",
				zap.Error(err),
			)
		}

	}

	var issuer model.User
	resp := database.DB.Where("ID = ?", claims.Issuer).
		First(&issuer)

	if resp.RowsAffected < 1 || resp.Error != nil {
		logger.Log.Error("Issuer not found",
			zap.Error(err),
		)
	}

	return &issuer, nil
}
