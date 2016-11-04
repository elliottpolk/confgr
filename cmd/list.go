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

func List(c *cli.Context) error {
	res, err := http.Get(fmt.Sprintf("%s/list", server.GetConfgrAddr()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
