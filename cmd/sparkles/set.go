package main

import (
	"bufio"
	"io"
	"net/url"
	"os"

	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto/pgp"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/models"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func pipe() (string, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", errors.Wrap(err, "unable to stat stdin")
	}

	if fi.Mode()&os.ModeCharDevice != 0 {
		return "", ErrNoPipe
	}

	buf, res := bufio.NewReader(os.Stdin), make([]byte, 0)
	for {
		in, _, err := buf.ReadLine()
		if err != nil && err == io.EOF {
			break
		}
		res = append(res, in...)

		if len(res) > MaxData {
			return "", ErrDataTooLarge
		}
	}

	return string(res), nil
}

func set(encrypt, insecure bool, token, usr, raw, addr string) (*models.Secret, error) {
	s, err := models.ParseSecret(raw)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse secret")
	}

	// ensure the secret has an ID set
	if len(s.Id) < 1 {
		s.Id = uuid.New().String()
	}

	if encrypt {
		c := &pgp.Crypter{Token: []byte(token)}

		// encrypt the content of the secret
		cypher, err := c.Encrypt([]byte(s.Content))
		if err != nil {
			return nil, errors.Wrap(err, "unable to encrypt secret content")
		}

		// set the content to the encrypted text
		s.Content = string(cypher)
	}

	params := url.Values{
		service.UserParam: []string{usr},
		service.AppParam:  []string{s.App},
		service.EnvParam:  []string{s.Env},
		service.IdParam:   []string{s.Id},
	}

	res, err := send(asURL(addr, service.PathSecrets, params.Encode()), s.MustString(), insecure)
	if err != nil {
		return nil, errors.Wrap(err, "unable to send secret")
	}

	in, err := models.ParseSecret(string(res))
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse in service response")
	}

	return in, nil
}
