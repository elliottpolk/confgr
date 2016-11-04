// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elliottpolk/confgr/server"
	"github.com/urfave/cli"
)

func Remove(c *cli.Context) {
	c.Command.VisibleFlags()

	app := c.String(AppFlag)
	env := c.String(EnvFlag)

	if len(env) < 1 {
		fmt.Println("an environment value must be specified")
		return
	}

	addr := server.GetConfgrAddr()
	res, err := http.Get(fmt.Sprintf("%s/remove?app=%s&env=%s", addr, app, env))
	if err != nil {
		if res.Body != nil {
			res.Body.Close()
		}

		fmt.Printf("unable to remove configuration for app %s: %v\n", app, err)
		return
	}
	defer res.Body.Close()

	if code := res.StatusCode; code != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("unable to read remove error response for app %s: %v\n", app, err)
			return
		}

		fmt.Printf("remove API responded with a status code other than OK: %d\n", code)
		fmt.Printf("response: %s\n", string(body))
	}
}
