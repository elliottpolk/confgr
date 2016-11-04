// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elliottpolk/confgr/config"
	"github.com/elliottpolk/confgr/pgp"
	"github.com/elliottpolk/confgr/server"
	"github.com/urfave/cli"
)

func Get(c *cli.Context) {
	c.Command.VisibleFlags()

	token := c.String(TokenFlag)
	decrypt := c.Bool(DecryptFlag)
	if decrypt && len(token) < 1 {
		fmt.Println("decryption token must be provided if decryption flag is set to true")
		return
	}

	addr := server.GetConfgrAddr()
	app := c.String(AppFlag)
	env := c.String(EnvFlag)

	res, err := http.Get(fmt.Sprintf("%s/get?app=%s&env=%s", addr, app, env))
	if err != nil {
		if res.Body != nil {
			res.Body.Close()
		}

		fmt.Printf("unable to retrieve config for app %s: %v\n", app, err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("unable to read config data for app %s: %v\n", app, err)
		return
	}

	fmt.Printf("\n%s\n", string(body))

	if decrypt {
		if ptxt := decryptCfg(body, token); len(ptxt) > 0 {
			fmt.Println("decrypted config:")
			fmt.Printf("%s\n", ptxt)
		}
	}
}

func decryptCfg(data []byte, token string) string {
	cfg := &config.Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("unable to decrypt results: %v\n", err)
		return ""
	}

	t, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		fmt.Printf("unable to decrypt results: %v\n", err)
		return ""
	}

	plaintxt, err := pgp.Decrypt(t, []byte(cfg.Value))
	if err != nil {
		fmt.Printf("unable to decrypt results: %v\n", err)
		return ""
	}

	return string(plaintxt)
}
