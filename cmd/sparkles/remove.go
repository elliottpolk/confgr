package main

import (
	"fmt"
	"net/url"

	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"

	"github.com/pkg/errors"
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
