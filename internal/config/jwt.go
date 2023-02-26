package config

import (
	"bookbox-backend/pkg/logger"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type JWTConfig struct {
	RefreshPrivateKey ed25519.PrivateKey
	RefreshPublicKey  ed25519.PublicKey
	RefreshExpiry     time.Duration

	AccessPrivateKey ed25519.PrivateKey
	AccessPublicKey  ed25519.PublicKey
	AccessExpiry     time.Duration
}

func readED25519PrivateKey(envVar string) (ed25519.PrivateKey, error) {

	keyFilePEM, err := ReadFromRelativePath(
		os.Getenv(envVar),
	)
	if err != nil {
		return nil, err
	}

	keyDecodedPEM, _ := pem.Decode(keyFilePEM)

	keyParsed, err := x509.ParsePKCS8PrivateKey(
		keyDecodedPEM.Bytes,
	)
	if err != nil {
		return nil, err
	}

	return keyParsed.(ed25519.PrivateKey), nil

}

var JWT JWTConfig

func init() {

	log := logger.Log

	// Read refresh token's private key file.
	refreshPrivateKey, err := readED25519PrivateKey(
		"JWT_REFRESH_PRIVATE_KEY",
	)
	if err != nil {
		log.Error("readED25519PrivateKey failed",
			zap.Error(err),
		)
		return
	}
	refreshExpiry, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRY"))
	if err != nil {
		log.Error("JWT_REFRESH_EXPIRY reading failed",
			zap.Error(err),
		)
		return
	}

	// Read access token's private key file.
	accessPrivateKey, err := readED25519PrivateKey("JWT_ACCESS_PRIVATE_KEY")
	if err != nil {
		log.Error("readED25519PrivateKey failed",
			zap.Error(err),
		)
		return
	}
	accessExpiry, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY"))
	if err != nil {
		log.Error("JWT_ACCESS_EXPIRY reading failed",
			zap.Error(err),
		)
		return
	}

	JWT = JWTConfig{

		// Refresh key.
		RefreshPrivateKey: refreshPrivateKey,
		RefreshExpiry:     time.Minute * time.Duration(refreshExpiry),

		// Access key.
		AccessPrivateKey: accessPrivateKey,
		AccessExpiry:     time.Minute * time.Duration(accessExpiry),
	}
}

func ReadFromRelativePath(revPath string) ([]byte, error) {

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path = path + "/" + revPath[2:]

	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return content, nil
}
