package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// serverInitialize will return new instance of configured server with given router
func Initialize(ip net.IP, port int, router *gin.Engine) (server *http.Server) {
	// Enable TLS1.3 (https://golang.org/pkg/crypto/tls/)
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")

	// If httpsRecord is enabled, disable HTTP2
	nextProtos := []string{
		"h2",
		"http/1.1",
	}
	// Initialize TLS Configuration
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			// TLS1.2
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			// TLS1.3
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
		NextProtos: nextProtos,
	}

	// Instantiate new HTTP Server with tlsConfig
	var (
		writeTimeout      = 5 * time.Second
		readTimeout       = 5 * time.Second
		readHeaderTimeout = 5 * time.Second
		idleTimeout       = 5 * time.Second
	)

	// Make silent error log
	errorLog := log.New(ioutil.Discard, "", 0)

	server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", ip, port),
		WriteTimeout:      writeTimeout,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
		TLSConfig:         tlsConfig,
		Handler:           router,
		ErrorLog:          errorLog,
	}
	server.SetKeepAlivesEnabled(false)

	return server
}

// Wait will block processes and wait for signal to stop/shut-down given server
func Wait(httpServer *http.Server, logger *zap.Logger) {
	log := logger.WithOptions(zap.Fields(
		zap.String("address", httpServer.Addr),
		zap.Duration("readTimeout", httpServer.ReadTimeout),
		zap.Duration("writeTimeout", httpServer.WriteTimeout),
		zap.Duration("readHeaderTimeout", httpServer.ReadHeaderTimeout),
		zap.Duration("idleTimeout", httpServer.IdleTimeout),
	))
	log.Info("waiting for HTTPS server signals")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	exitChannel := make(chan string)
	go func() {
		for {
			s := <-signalChannel
			switch s {
			case syscall.SIGHUP:
				logger.Warn("received SIGHUP")
				// exitChannel <- "SIGHUP"
			case syscall.SIGINT:
				logger.Warn("received SIGINT")
				exitChannel <- "SIGINT"
			case syscall.SIGTERM:
				logger.Warn("received SIGTERM")
				exitChannel <- "SIGTERM"
			case syscall.SIGQUIT:
				logger.Warn("received SIGQUIT")
				exitChannel <- "SIGQUIT"
			default:
				logger.Warn("received SIGHUP")
				// exitChannel <- "UNKNOWN"
			}
		}
	}()

	exitSignal := <-exitChannel
	log.Warn("attempting to stop HTTPS server", zap.String("exitSignal", exitSignal))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	log.Warn("HTTPS server stopped successfully", zap.String("exitSignal", exitSignal))
	httpServer.Shutdown(ctx)
	os.Exit(0)
}
