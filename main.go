// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"os"

	"github.com/elliottpolk/confgr/cmd"
	"github.com/elliottpolk/confgr/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "confgr"
	app.Usage = "a simple configuration service"
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list all available configurations",
			Action:  cmd.List,
		},
		{
			Name:   "get",
			Usage:  "get a specific app configuration",
			Action: cmd.Get,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, a", cmd.AppFlag),
					Usage: "app name of configuration",
				},
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, e", cmd.EnvFlag),
					Usage: "configuration environment (e.g. PROD, DEV, TEST)",
				},
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, t", cmd.TokenFlag),
					Usage: "token to decrypt configuration",
				},
				cli.BoolFlag{
					Name:  cmd.DecryptFlag,
					Usage: "decrypt configuration",
				},
			},
		},
		{
			Name:   "set",
			Usage:  "set an app configuration",
			Action: cmd.Set,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, a", cmd.AppFlag),
					Usage: "app name of configuration",
				},
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, e", cmd.EnvFlag),
					Usage: "configuration environment (e.g. PROD, DEV, TEST)",
				},
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, c", cmd.CfgFlag),
					Usage: "configuration to store",
				},
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, t", cmd.TokenFlag),
					Usage: "token to encrypt configuration",
				},
				cli.BoolFlag{
					Name:  cmd.EncryptFlag,
					Usage: "encrypt configuration",
				},
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"del"},
			Usage:   "remove an app configuration",
			Action:  cmd.Remove,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, a", cmd.AppFlag),
					Usage: "app name of configuration",
				},
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, e", cmd.EnvFlag),
					Usage: "configuration environment (e.g. PROD, DEV, TEST)",
				},
			},
		},
		{
			Name:   "server",
			Usage:  "confgr server for storing app configurations",
			Action: server.Start,
		},
	}

	app.Run(os.Args)
}
