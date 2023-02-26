package middlewares

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"bookbox-backend/pkg/logger"

	"go.uber.org/zap"
)

var (
	environment  = os.Getenv("ENVIRONMENT")
	serverDomain = os.Getenv("SERVER_DOMAIN")
)

func init() {
	environment = strings.ToLower(environment)
}

func SetOrigin(ip net.IP, port int) {
	if environment == "development" {
		Origin = "*"
	} else {
		Origin = fmt.Sprintf("https://%s:%d", serverDomain, port)
	}
}

// middlewareRecovery recovers middleware from a problem
func middlewareRecovery() {
	log := logger.Log.WithOptions(zap.Fields())

	if err := recover(); err != nil {
		_, file, _, _ := runtime.Caller(2)
		file = filepath.Base(file)
		file = strings.Split(file, ".")[0]
		file = strings.Title(file)

		log.Error(fmt.Sprintf("panic recovered in %s Middleware", file),
			zap.String("recover", fmt.Sprintf("%v", err)),
		)
	}
}

var rootDir string

func init() {
	rootDir = os.Getenv("rootDir")
	if rootDir == "" {
		rootDir = "/tmp"
	}

	filePath := filepath.Join(rootDir, "db.json")
	_, err := os.Stat(filePath)
	if err != nil {
		user := []map[string]interface{}{
			{
				"email":     "admin",
				"id":        "admin",
				"firstname": "admin",
				"lastname":  "admin",
				"role":      "admin",
				"password":  "1234.dbe9787aaf4002c6662e490b3f1f7512807459b6dee2e1c2e56738e1cbbd993c",
			},
		}

		raw, _ := json.Marshal(user)
		ioutil.WriteFile(filePath, raw, 0744)
	}
}
