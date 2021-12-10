package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"artemkv.net/winaday/app"
	"artemkv.net/winaday/health"
	"artemkv.net/winaday/reststats"
	"artemkv.net/winaday/server"
	"github.com/gin-gonic/gin"
)

var version = "1.2"

func main() {
	// setup logging
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// load .env
	LoadDotEnv()

	// initialize session encryption key
	sessionEncryptionPassphrase := GetMandatoryString("WINADAY_SESSION_ENCRYPTION_PASSPHRASE")
	app.SetEncryptionPassphrase(sessionEncryptionPassphrase)

	// initialize REST stats
	reststats.Initialize(version)

	// configure router
	allowedOrigin := GetMandatoryString("WINADAY_ALLOW_ORIGIN")
	router := gin.New()
	app.SetupRouter(router, allowedOrigin)

	// determine whether to use HTTPS
	useTls := GetBoolean("WINADAY_TLS")
	certFile := ""
	keyFile := ""
	if useTls {
		certFile = GetMandatoryString("WINADAY_CERT_FILE")
		keyFile = GetMandatoryString("WINADAY_KEY_FILE")
	}

	serverConfig := &server.ServerConfiguration{
		UseTls:   useTls,
		CertFile: certFile,
		KeyFile:  keyFile,
	}

	// determine port
	port := GetOptionalString("WINADAY_PORT", ":8700")

	// start the server
	server.Serve(router, port, serverConfig, func() {
		health.SetIsReadyGlobally()
	})
}
