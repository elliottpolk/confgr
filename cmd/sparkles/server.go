package main

import (
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend"
	cli "gopkg.in/urfave/cli.v2"
)

//	server flags
var (
	stdListenPortFlag = cli.StringFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Value:   "8080",
		Usage:   "HTTP port to listen on",
		EnvVars: []string{"PSPARKLES_HTTP_PORT"},
	}

	tlsListenPortFlag = cli.StringFlag{
		Name:    "tls-port",
		Value:   "8443",
		Usage:   "HTTPS port to listen on",
		EnvVars: []string{"PSPARKLES_HTTPS_PORT"},
	}

	tlsCertFlag = cli.StringFlag{
		Name:    "tls-cert",
		Usage:   "TLS certificate file for HTTPS",
		EnvVars: []string{"PSPARKLES_TLS_CERT"},
	}

	tlsKeyFlag = cli.StringFlag{
		Name:    "tls-key",
		Usage:   "TLS key file for HTTPS",
		EnvVars: []string{"PSPARKLES_TLS_KEY"},
	}

	datastoreTypeFlag = cli.StringFlag{
		Name:    "datastore-type",
		Aliases: []string{"dst"},
		Value:   backend.File,
		Usage:   "backend type to be used for storage",
		EnvVars: []string{"PSPARKLES_DS_TYPE"},
	}

	datastoreFileFlag = cli.StringFlag{
		Name:    "datastore-file",
		Aliases: []string{"dsf"},
		Value:   "/var/lib/peppermint-sparkles/psparkles.db",
		Usage:   "name / location of file for storing secrets",
		EnvVars: []string{"PSPARKLES_DS_FILE"},
	}

	datastoreAddrFlag = cli.StringFlag{
		Name:    "datastore-addr",
		Aliases: []string{"dsa"},
		Value:   "localhost:6379",
		Usage:   "address for the remote datastore",
		EnvVars: []string{"PSPARKLES_DS_ADDR"},
	}
)
