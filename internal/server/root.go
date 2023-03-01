package server

import (
	"os"
	"strconv"

	"net"
	"net/http"

	_ "trikliq-airport-finder/internal/config"
	_ "trikliq-airport-finder/pkg/redis"

	"trikliq-airport-finder/internal/server/ca"
	"trikliq-airport-finder/internal/server/middlewares"
	"trikliq-airport-finder/internal/server/router"
	"trikliq-airport-finder/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	ip         = net.IPv4(127, 0, 0, 1)
	port       = 4443
	httpServer *http.Server
)

func init() {
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort != "" {
		port, _ = strconv.Atoi(serverPort)
	}

	serverIP := os.Getenv("SERVER_IP")
	if serverIP != "" {
		parsedIP, _, err := net.ParseCIDR(serverIP + "/32")
		if err == nil {
			ip = parsedIP
		}
	}

	middlewares.SetOrigin(ip, port)
}

func Start() {
	log := logger.Log.WithOptions(zap.Fields(
		zap.String("ip", ip.String()),
	))

	log.Info("starting engine")

	// Sets default size only if not set
	gin.SetMode(gin.ReleaseMode)

	// Auto generate self signed certificate and private key
	ca.Setup(logger.Log)
	certFile := ca.GetCertificate()
	keyFile := ca.GetPrivateKey()
	log = log.WithOptions(zap.Fields(
		zap.String("certFile", certFile),
		zap.String("keyFile", keyFile),
	))

	// Start HTTPS server
	httpServer = Initialize(ip, port, router.Router)
	go func() {
		errorsStartingUp := 0
		var err error

		for errorsStartingUp < 5 {
			log.Info("attempting to start HTTPS server",
				zap.Int("port", port),
				zap.Int("errorsStartingUp", errorsStartingUp),
			)

			httpServer = Initialize(ip, port, router.Router)
			//err = httpServer.ListenAndServeTLS(certFile, keyFile)
			err = httpServer.ListenAndServe()
			if err != nil {
				log.Info("retrying to start HTTPS server, because of an error",
					zap.Int("port", port),
					zap.Int("errorsStartingUp", errorsStartingUp),
					zap.Error(err),
				)

				port++
				errorsStartingUp++
				continue
			}
			break
		}

		log.Panic("failed to start HTTPS server",
			zap.Int("port", port),
			zap.Int("errorsStartingUp", errorsStartingUp),
			zap.Error(err),
		)
	}()

	Wait(httpServer, log)
}
