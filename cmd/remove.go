package main

import (
	"fmt"
	"net/url"
	"os/user"

	"github.com/manulife-gwam/peppermint-sparkles/service"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var (
	Remove = &cli.Command{
		Name:    "delete",
		Aliases: []string{"del", "rm"},
		Flags: []cli.Flag{
			&AddrFlag,
			&AppNameFlag,
			&AppEnvFlag,
			&SecretIdFlag,
			&InsecureFlag,
		},
		Usage: "deletes a secret",
		Action: func(context *cli.Context) error {
			addr := context.String(AddrFlag.Name)
			id := context.String(SecretIdFlag.Name)

			u, err := user.Current()
			if err != nil {
				return cli.Exit(errors.Wrap(err, "unable to retrieve current, logged-in user"), 1)
			}

			params := &url.Values{
				service.UserParam: []string{u.Username},
				service.AppParam:  []string{context.String(AppNameFlag.Name)},
				service.EnvParam:  []string{context.String(AppEnvFlag.Name)},
			}

			insecure := context.Bool(InsecureFlag.Name)

			if err := rm(insecure, id, addr, params); err != nil {
				return cli.Exit(errors.Wrap(err, "unable to remove secret"), 1)
			}

			return nil
		},
	}
)

func rm(insecure bool, id, addr string, params *url.Values) error {
	if len(id) < 1 {
		return errors.New("a valid secret ID must be provided")
	}

	if len(params.Get(service.AppParam)) < 1 {
		return errors.New("a valid secret app name must be provided")
	}

	if len(params.Get(service.EnvParam)) < 1 {
		return errors.New("a valid secret environment must be provided")
	}

	if _, err := del(asURL(addr, fmt.Sprintf("%s/%s", service.PathSecrets, id), params.Encode()), insecure); err != nil {
		return err
	}

	return nil
}
