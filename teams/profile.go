package teams

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/pkg/errors"
)

const BucketProfile string = ""

type Profile struct {
	Team string `json:"team"`
	Key  []byte `json:"key"`
}

const DefaultKeyBits = 5096

func generateKey() ([]byte, error) {
	pk, err := rsa.GenerateKey(rand.Reader, DefaultKeyBits)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate key")
	}

	return x509.MarshalPKCS1PrivateKey(pk), nil
}

func (p *Profile) Write() error {
	return nil
}
