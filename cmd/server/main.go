package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/elliottpolk/peppermint-sparkles/server"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	altsrc "github.com/urfave/cli/v2/altsrc"
)

var (
	version  string
	compiled string = fmt.Sprint(time.Now().Unix())
	githash  string

	// log levels
	levels = map[string]log.Level{
		// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
		// logging level is set to Panic.
		"fatal": log.FatalLevel,
		// ErrorLevel level. Logs. Used for errors that should definitely be noted.
		// Commonly used for hooks to send errors to an error tracking service.
		"error": log.ErrorLevel,
		// InfoLevel level. General operational entries about what's going on inside the
		// application.
		"info": log.InfoLevel,
		// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
		"debug": log.DebugLevel,
	}

	cfgFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "server.config.file",
		Aliases: []string{"config-file", "cfg", "c"},
		Usage:   "optional path and filename to server config",
	})

	logLevelFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "server.config.log.level",
		Aliases: []string{"log-level", "ll"},
		Value:   "info",
		Usage:   "log level output",
		EnvVars: []string{
			"PSPARKLES_LOG_LEVEL",
			"SPARKLES_LOG_LEVEL",
		},
	})

	httpPortFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "server.http.port",
		Aliases: []string{"http-port", "http.port", "port", "p"},
		Value:   "8080",
		Usage:   "HTTP port to listen on",
		EnvVars: []string{
			"PSPARKLES_HTTP_PORT",
			"SPARKLES_HTTP_PORT",
		},
	})

	httpsPortFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "server.https.port",
		Aliases: []string{"tls-port", "https-port", "https.port"},
		Value:   "8443",
		Usage:   "HTTPS port to listen on",
		EnvVars: []string{
			"PSPARKLES_HTTPS_PORT",
			"SPARKLES_HTTPS_PORT",
		},
	})

	tlsCertFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "server.https.cert",
		Aliases: []string{"tls-cert", "https-cert", "https.cert"},
		Usage:   "TLS certificate file for HTTPS",
		EnvVars: []string{
			"PSPARKLES_TLS_CERT",
			"SPARKLES_TLS_CERT",
		},
	})

	tlsKeyFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "server.https.key",
		Aliases: []string{"tls-key", "https-key", "https.key"},
		Usage:   "TLS key file for HTTPS",
		EnvVars: []string{
			"PSPARKLES_TLS_KEY",
			"SPARKLES_TLS_KEY",
		},
	})

	datastoreTypeFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "datastore.type",
		Aliases: []string{"datastore-type", "dst"},
		//Value:   backend.File,
		Usage: "backend type to be used for storage",
		EnvVars: []string{
			"PSPARKLES_DS_TYPE",
			"SPARKLES_DS_TYPE",
		},
	})

	datastoreFileFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "datastore.file",
		Aliases: []string{"datastore-file", "dsf"},
		Value:   "/var/lib/peppermint-sparkles/psparkles.db",
		Usage:   "name / location of file for storing secrets",
		EnvVars: []string{
			"PSPARKLES_DS_FILE",
			"SPARKLES_DS_FILE",
		},
	})

	datastoreAddrFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "datastore.address",
		Aliases: []string{"datastore-addr", "dsa"},
		Value:   "localhost:6379",
		Usage:   "address for the remote datastore",
		EnvVars: []string{
			"PSPARKLES_DS_ADDR",
			"SPARKLES_DS_ADDR",
		},
	})
)

func main() {
	// this is needed to get the compile time which will be set at build.
	// if the build did not set this (i.e. `make` was not used) , it will
	// default to the current timestamp
	ct, err := strconv.ParseInt(compiled, 0, 0)
	if err != nil {
		panic(err) // panic because this wasn't set properly
	}

	app := cli.App{
		Name:      "sparklesd",
		Usage:     "",
		Copyright: fmt.Sprintf("Copyright © 2018-%s Elliott Polk", time.Now().Format("2006")),
		Version:   fmt.Sprintf("%s | compiled %s | commit %s", version, time.Unix(ct, -1).Format(time.RFC3339), githash),
		Compiled:  time.Unix(ct, -1),
		Flags: []cli.Flag{
			cfgFlag,
			logLevelFlag,
			httpPortFlag,
			httpsPortFlag,
			tlsCertFlag,
			tlsKeyFlag,
			datastoreTypeFlag,
			datastoreFileFlag,
			datastoreAddrFlag,
		},
		Before: func(ctx *cli.Context) error {
			if len(ctx.String(cfgFlag.Name)) > 0 {
				return altsrc.InitInputSourceWithContext(ctx.App.Flags, altsrc.NewYamlSourceFromFlagFunc(cfgFlag.Name))(ctx)
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {
			log.SetLevel(levels[ctx.String(logLevelFlag.Name)])
			log.Debug("logger level set")

			// TODO:
			// setup the datastore
			// attach service handler(s)
			// setup HTTPS in gofunc
			// listen on HTTP
			mux := http.NewServeMux()
			mux = server.Handle(mux, &server.Handler{})

			go func() {
				var (
					cert = ctx.String(tlsCertFlag.Name)
					key  = ctx.String(tlsKeyFlag.Name)
				)

				if len(cert) < 1 || len(key) < 1 {
					return
				}

				if _, err := os.Stat(cert); err != nil {
					log.Errorf("unable to access TLS cert file: %s", cert)
					return
				}

				if _, err := os.Stat(key); err != nil {
					log.Errorf("unable to access TLS key file: %s", key)
					return
				}

				svr := &http.Server{
					Addr:    fmt.Sprintf(":%s", ctx.String(httpsPortFlag.Name)),
					Handler: mux,
					TLSConfig: &tls.Config{
						PreferServerCipherSuites: true,
						CurvePreferences: []tls.CurveID{
							tls.CurveP256,
							tls.X25519,
						},
						MinVersion: tls.VersionTLS12,
						CipherSuites: []uint16{
							tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
							tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
							tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
							tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
							tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
							tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

							// excluding due to no forward secrecy, but leaving
							// as it might be necessary for some clients
							// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
							// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
						},
					},
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  20 * time.Second,
				}

				log.Debug("starting HTTPS listener")
				log.Fatal(svr.ListenAndServeTLS(cert, key))
			}()

			log.Debug("starting HTTP listener")
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ctx.String(httpPortFlag.Name)), mux))
			return nil
		},
	}

	app.Run(os.Args)
}
