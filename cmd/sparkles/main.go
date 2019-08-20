package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"time"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/go-common/pcf/vcap"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend"
	fileds "git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend/file"
	redisds "git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend/redis"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"
	bolt "github.com/coreos/bbolt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	cli "gopkg.in/urfave/cli.v2"
)

var version string

func main() {
	log.Init(version)

	app := cli.App{
		Copyright: "Copyright © 2018",
		Usage:     "Server and client for managing super special secrets 🦄",
		Version:   version,
		Commands: []*cli.Command{
			&cli.Command{
				Aliases: []string{"ls", "list"},
				Flags: []cli.Flag{
					&addrFlag,
					&appNameFlag,
					&appEnvFlag,
					&secretIdFlag,
					&decryptFlag,
					&tokenFlag,
					&insecureFlag,
				},
				Usage: "retrieves secrets",
				Action: func(context *cli.Context) error {
					addr := context.String(addrFlag.Name)
					if len(addr) < 1 {
						cli.ShowCommandHelpAndExit(context, context.Command.FullName(), 1)
						return nil
					}

					token := context.String(tokenFlag.Name)
					decrypt := context.Bool(decryptFlag.Name)

					if decrypt && len(token) < 1 {
						return cli.Exit(errors.New("decrypt token must be specified in order to decrypt"), 1)
					}

					params := &url.Values{
						service.AppParam: []string{context.String(appNameFlag.Name)},
						service.EnvParam: []string{context.String(appEnvFlag.Name)},
					}

					insecure := context.Bool(insecureFlag.Name)

					s, err := get(decrypt, insecure, token, addr, context.String(secretIdFlag.Name), params)
					if err != nil {
						return cli.Exit(errors.Wrap(err, "unable to retrieve secert"), 1)
					}

					log.Infof("\n%s\n", s.MustString())
					return nil
				},
			},
			&cli.Command{
				Name:    "set",
				Aliases: []string{"add", "create", "new", "update"},
				Flags: []cli.Flag{
					&addrFlag,
					&secretFlag,
					&secretFileFlag,
					&encryptFlag,
					&tokenFlag,
					&secretIdFlag,
					&insecureFlag,
				},
				Usage: "adds or updates a secret",
				Action: func(context *cli.Context) error {
					addr := context.String(addrFlag.Name)
					if len(addr) < 1 {
						cli.ShowCommandHelpAndExit(context, context.Command.FullName(), 1)
						return nil
					}

					raw, f := context.String(secretFlag.Name), context.String(secretFileFlag.Name)
					if len(raw) > 0 && len(f) > 0 {
						return cli.Exit(errors.New("only 1 input method is allowed"), 1)
					}

					//	raw should not have anything if this is true
					if len(f) > 0 {
						info, err := os.Stat(f)
						if err != nil {
							return cli.Exit(errors.Wrap(err, "uanble to access secrets file"), 1)
						}

						if info.Size() > int64(MaxData) {
							return cli.Exit(errors.New("secret must be less than 3MB"), 1)
						}

						r, err := ioutil.ReadFile(f)
						if err != nil {
							return cli.Exit(errors.Wrap(err, "unable to read in secret file"), 1)
						}

						raw = string(r)
					}

					//	if raw is still empty at this point, attempt to read in piped data
					tick := 0
					for len(raw) < 1 {
						if tick > 0 {
							return cli.Exit(errors.New("a valid secret must be specified"), 1)
						}

						r, err := pipe()
						if err != nil {
							switch err {
							case ErrNoPipe:
								return cli.Exit(errors.New("a valid secret must be specified"), 1)
							case ErrDataTooLarge:
								return cli.Exit(errors.New("secret must be less than 3MB"), 1)
							default:
								return cli.Exit(errors.Wrap(err, "unable to read piped in data"), 1)
							}
						}
						raw, tick = r, +1
					}

					encrypt := context.Bool(encryptFlag.Name)
					token := context.String(tokenFlag.Name)
					if encrypt {
						if len(token) < 1 {
							//	attempt to generate a token if one not provided, erroring and exiting
							//	if unable. This attempts to prevent encrypting with empty string
							t, err := crypto.NewToken()
							if err != nil {
								return cli.Exit(errors.Wrap(err, "unable to generate encryption token"), 1)
							}
							token = t
						}
					}

					// get current logged in user
					u, err := user.Current()
					if err != nil {
						return cli.Exit(errors.Wrap(err, "unable to retrieve current, logged-in user"), 1)
					}

					insecure := context.Bool(insecureFlag.Name)

					s, err := set(encrypt, insecure, token, u.Username, raw, addr)
					if err != nil {
						return cli.Exit(errors.Wrap(err, "unable to set secret"), 1)
					}

					//	ensure to display encryption token, since it may have been generated
					if encrypt {
						log.Infof(tag, "token: %s", token)
					}
					log.Infof(tag, "secret:\n%s", s.MustString())

					return nil
				},
			},
			&cli.Command{
				Name:    "delete",
				Aliases: []string{"del", "rm"},
				Flags: []cli.Flag{
					&addrFlag,
					&appNameFlag,
					&appEnvFlag,
					&secretIdFlag,
					&insecureFlag,
				},
				Usage: "deletes a secret",
				Action: func(context *cli.Context) error {
					addr := context.String(addrFlag.Name)
					id := context.String(secretIdFlag.Name)

					u, err := user.Current()
					if err != nil {
						return cli.Exit(errors.Wrap(err, "unable to retrieve current, logged-in user"), 1)
					}

					params := &url.Values{
						service.UserParam: []string{u.Username},
						service.AppParam:  []string{context.String(appNameFlag.Name)},
						service.EnvParam:  []string{context.String(appEnvFlag.Name)},
					}

					insecure := context.Bool(insecureFlag.Name)

					if err := rm(insecure, id, addr, params); err != nil {
						return cli.Exit(errors.Wrap(err, "unable to remove secret"), 1)
					}

					return nil
				},
			},
			&cli.Command{
				Name:    "server",
				Aliases: []string{"serve"},
				Flags: []cli.Flag{
					&stdListenPortFlag,
					&tlsListenPortFlag,
					&tlsCertFlag,
					&tlsKeyFlag,
					&datastoreAddrFlag,
					&datastoreFileFlag,
					&datastoreTypeFlag,
				},
				Usage: "start the server",
				Action: func(context *cli.Context) error {
					var (
						ds  backend.Datastore
						err error
					)

					dst := context.String(datastoreTypeFlag.Name)

					switch dst {
					case backend.Redis:
						opts := &redis.Options{Addr: context.String(datastoreAddrFlag.Name)}

						//	check if running in PCF pull the vcap services if available
						services, err := vcap.GetServices()
						if err != nil {
							return cli.Exit(errors.Wrap(err, "unable to retrieve vcap services"), 1)
						}

						if services != nil {
							if i := services.Tagged(dst); i != nil {
								creds := i.Credentials
								opts = &redis.Options{
									Addr:     fmt.Sprintf("%s:%d", creds["host"].(string), int(creds["port"].(float64))),
									Password: creds["password"].(string),
								}
							}
						}

						if ds, err = redisds.Open(opts); err != nil {
							return cli.Exit(errors.Wrap(err, "unable to open connection to datastore"), 1)
						}

					case backend.File:

						//	FIXME ... include / handle additional bolt options (e.g. timeout, etc)
						fname := context.String(datastoreFileFlag.Name)
						if ds, err = fileds.Open(fname, bolt.DefaultOptions); err != nil {
							return cli.Exit(errors.Wrap(err, "unable to open connection to datastore"), 1)
						}

					default:
						return cli.Exit(errors.Errorf("%s is not a supported datastore type", dst), 1)
					}

					defer ds.Close()
					log.Debug(tag, "datastore opened")

					mux := http.NewServeMux()

					//	attach current service handler
					mux = service.Handle(mux, &service.Handler{Backend: ds})

					//	start HTTPS listener in a seperate go routine since it is a blocking func
					go func() {
						cert, key := context.String(tlsCertFlag.Name), context.String(tlsKeyFlag.Name)
						if len(cert) < 1 || len(key) < 1 {
							return
						}

						if _, err := os.Stat(cert); err != nil {
							log.Error(tag, err, "unable to access TLS cert file")
							return
						}

						if _, err := os.Stat(key); err != nil {
							log.Error(tag, err, "unable to access TLS key file")
							return
						}

						svr := &http.Server{
							Addr:    fmt.Sprintf(":%s", context.String(tlsListenPortFlag.Name)),
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

						log.Debug(tag, "starting HTTPS listener")
						log.Fatal(tag, svr.ListenAndServeTLS(cert, key))
					}()

					log.Debug(tag, "starting HTTP listener")

					addr := fmt.Sprintf(":%s", context.String(stdListenPortFlag.Name))
					log.Fatal(tag, http.ListenAndServe(addr, mux))

					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}
