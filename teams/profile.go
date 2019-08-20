package teams

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const BucketProfiles string = "profiles"

type Profile struct {
	Id       string `json:"id"`
	Team     string `json:"team"`
	Username string `json:"username"`
	Key      []byte `json:"key"`

	// datastore metadata
	Meta *meta `json:"meta"`
}

const DefaultKeyBits = 5096

func generateKey() ([]byte, error) {
	pk, err := rsa.GenerateKey(rand.Reader, DefaultKeyBits)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate key")
	}

	return x509.MarshalPKCS1PrivateKey(pk), nil
}

func (p *Profile) String() string {
	out, _ := json.MarshalIndent(p, "", " ")
	return string(out)
}

func (p *Profile) Write(ds *datastore) error {
	now := time.Now().UnixNano()

	if p.Meta == nil {
		p.Meta = &meta{
			Created: now,
			Status:  active,
		}
	}

	// this may be a new entry, so generate a new id
	if len(p.Id) < 1 {
		p.Id = uuid.New().String()
	}

	// ensure that the updated value is always set to now
	p.Meta.Updated = now
	return ds.save(BucketProfiles, p.Id, p.String())
}

func GetProfile(ds *datastore, id string) (*Profile, error) {
	p := &Profile{}
	if err := json.Unmarshal(ds.retrieve(BucketProfiles, id), &p); err != nil {
		return nil, err
	}

	return p, nil
}
