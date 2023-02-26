package auth

import (
	"bookbox-backend/internal/config"
	"bookbox-backend/internal/model"
	"bookbox-backend/pkg/redis"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type SignedTokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenPairCache struct {
	AccessTokenUUID  string    `json:"accessTokenUUID"`
	RefreshTokenUUID string    `json:"refreshTokenUUID"`
	EntryTime        time.Time `json:"entryTime"`
}

func DestroyToken(ctx *gin.Context, userID uint) error {

	err := redis.Client.Del(
		fmt.Sprintf("token_pair-%d", userID),
	).Err()
	if err != nil {
		return err
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	maxAge := int(expiration.Unix())
	domain := ctx.Request.Host
	domain = strings.Split(domain, ":")[0]

	ctx.SetCookie("access_token", "", maxAge, "/", domain, true, false)
	ctx.SetCookie("refresh_token", "", maxAge, "/", domain, true, false)

	return nil
}

func RefreshToken(ctx *gin.Context) (issuer *model.User, err error) {
	tokenToValidate, err := ctx.Cookie("refresh_token")
	if err != nil {
		return nil, err
	}

	return ValidateTokenClaims(ctx, tokenToValidate, true)
}

func NewTokenPair(ctx *gin.Context, userID uint) error {
	refreshTokenUUID := uuid.New().String()
	accessTokenUUID := uuid.New().String()

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.RegisteredClaims{
		Issuer:    strconv.FormatUint(uint64(userID), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWT.RefreshExpiry)),
		ID:        refreshTokenUUID,
	}).
		SignedString(config.JWT.RefreshPrivateKey)
	if err != nil {
		return err
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.RegisteredClaims{
		Issuer:    strconv.FormatUint(uint64(userID), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWT.AccessExpiry)),
		ID:        accessTokenUUID,
	}).
		SignedString(config.JWT.AccessPrivateKey)
	if err != nil {
		return err
	}

	cacheJSON, err := json.Marshal(TokenPairCache{
		AccessTokenUUID:  accessTokenUUID,
		RefreshTokenUUID: refreshTokenUUID,
		EntryTime:        time.Now(),
	})
	if err != nil {
		return err
	}

	err = redis.Client.Set(
		fmt.Sprintf("token_pair-%d", userID),
		string(cacheJSON),
		config.JWT.RefreshExpiry,
	).Err()
	if err != nil {
		return err
	}

	domain := ctx.Request.Host
	domain = strings.Split(domain, ":")[0]

	ctx.SetCookie("access_token", accessToken, 3600, "/", domain, true, false)
	ctx.SetCookie("refresh_token", refreshToken, 7776000, "/", domain, true, false)

	return nil
}

// Used to refresh tokens.
func NewAccessToken(ctx *gin.Context, userID uint) error {

	accessTokenUUID := uuid.New().String()

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.RegisteredClaims{
		Issuer:    strconv.FormatUint(uint64(userID), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWT.AccessExpiry)),
		ID:        accessTokenUUID,
	}).
		SignedString(config.JWT.AccessPrivateKey)
	if err != nil {
		return err
	}

	cacheJSON, _ := redis.Client.Get(fmt.Sprintf("token_pair-%d", userID)).
		Result()

	redis.Client.Get(fmt.Sprintf("token_pair-%d", userID))

	cachedTokenPair := TokenPairCache{}
	err = json.Unmarshal([]byte(cacheJSON), &cachedTokenPair)
	if err != nil {
		return err
	}

	cachedTokenPair.AccessTokenUUID = accessTokenUUID
	if err != nil {
		return err
	}

	raw, err := json.Marshal(cachedTokenPair)
	if err != nil {
		return err
	}

	elapsed := time.Since(cachedTokenPair.EntryTime)

	err = redis.Client.Set(
		fmt.Sprintf("token_pair-%d", userID),
		string(raw),
		config.JWT.RefreshExpiry-elapsed,
	).Err()

	if err != nil {
		return err
	}

	domain := ctx.Request.Host
	domain = strings.Split(domain, ":")[0]

	ctx.SetCookie("access_token", accessToken, 3600, "/", domain, true, false)
	//refresh -> cookie ?

	return nil
}
