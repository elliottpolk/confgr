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
	"strings"

	"github.com/elliottpolk/confgr/config"
	"github.com/elliottpolk/confgr/pgp"
	"github.com/elliottpolk/confgr/server"
	"github.com/elliottpolk/confgr/uuid"
	"github.com/urfave/cli"
)

func Set(c *cli.Context) {
	c.Command.VisibleFlags()

	app := c.String(AppFlag)
	env := c.String(EnvFlag)
	val := c.String(CfgFlag)

	encrypt := c.Bool(EncryptFlag)
	token := c.String(TokenFlag)
	if encrypt {
		tok, enc, err := encryptCfg(token, val)
		if err != nil {
			fmt.Printf("unable to encrypt config for app %s: %v\n", app, err)
			return
		}

		token = tok
		val = enc
	}

	cfg := &config.Config{app, env, val}
	out, err := json.Marshal(cfg)
	if err != nil {
		fmt.Printf("unable to marshal config to json for app %s: %v\n", app, err)
		return
	}

	if err := saveCfg(string(out)); err != nil {
		fmt.Printf("unable to save config for app %s: %v\n", app, err)
		return
	}

	if encrypt {
		fmt.Printf("token:           %s\n", token)
		fmt.Printf("token as base64: %s\n", base64.StdEncoding.EncodeToString([]byte(token)))
	}

	fmt.Printf("stored config:\n%s\n", string(out))
}

func encryptCfg(token, value string) (string, string, error) {
	if len(token) < 1 {
		if token = uuid.GetV4(); len(token) < 1 {
			return token, value, fmt.Errorf("UUID produced an empty string")
		}
	} else {
		t, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return token, value, err
		}

		token = string(t)
	}

	cypher, err := pgp.Encrypt([]byte(token), []byte(value))
	if err != nil {
		return token, value, err
	}

	return token, string(cypher), nil
}

func saveCfg(cfg string) error {
	addr := server.GetConfgrAddr()

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/set", addr), strings.NewReader(cfg))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if res.Body != nil {
			res.Body.Close()
		}

		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if code := res.StatusCode; code != http.StatusOK {
		fmt.Printf("server responded with status code %d\n", code)
		return fmt.Errorf("%s", string(body))
	}

	return nil
}
