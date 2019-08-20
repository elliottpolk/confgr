package main

import cli "gopkg.in/urfave/cli.v2"

var (
	appNameFlag = cli.StringFlag{
		Name:    "app-name",
		Aliases: []string{"a", "app"},
		Usage:   "app name of secret",
	}

	appEnvFlag = cli.StringFlag{
		Name:    "app-env",
		Aliases: []string{"e", "env"},
		Usage:   "environment of secret (e.g. PROD, DEV, TEST, etc.)",
	}

	secretIdFlag = cli.StringFlag{
		Name:    "secret-id",
		Aliases: []string{"id", "sid"},
		Usage:   "generated ID of secret",
	}

	secretFlag = cli.StringFlag{
		Name:    "secret",
		Aliases: []string{"s"},
		Usage:   "secret to be stored",
	}

	secretFileFlag = cli.StringFlag{
		Name:    "secret-file",
		Aliases: []string{"f"},
		Usage:   "filepath to secret",
	}

	encryptFlag = cli.BoolFlag{
		Name:  "encrypt",
		Value: true,
		Usage: "encrypt secrets",
	}

	decryptFlag = cli.BoolFlag{
		Name:  "decrypt",
		Usage: "decrypt secrets",
	}

	tokenFlag = cli.StringFlag{
		Name:    "token",
		Aliases: []string{"t", "tok"},
		Usage:   "token used to encrypt / decrypt secrets",
	}

	addrFlag = cli.StringFlag{
		Name:    "addr",
		Usage:   "secrets service address",
		EnvVars: []string{"PSPARKLES_ADDR"},
	}

	insecureFlag = cli.BoolFlag{
		Name:    "insecure",
		Aliases: []string{"k"},
		Value:   false,
		Usage:   "(TLS) this option explicitly allows to perform \"insecure\" SSL connections",
	}
)
