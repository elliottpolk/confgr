package bbolt

import (
	"fmt"

	"github.com/elliottpolk/peppermint-sparkles/internal/backend"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

const (
	PrimaryBucket string = "mints"
	AuditBucket   string = "slendermints"
)

type Repo struct {
	name   string
	bucket string
	db     *bbolt.DB
}

func NewRepo(n, b string) *Repo {
	return &Repo{
		name:   n,
		bucket: b,
	}
}

func (r *Repo) Open() error {
	var err error

	log.Debug("opening bbolt repo")
	r.db, err = bbolt.Open(r.name, 0600, &bbolt.Options{})
	if err != nil {
		return err
	}

	// ensure that the bucket exists
	log.Debug("ensuring the bucket exists for the repo")
	err = r.db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(r.bucket)); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return errors.Wrapf(err, "unable to ensure creation of repo bucket %s", r.bucket)
	}

	return nil
}

func (r *Repo) Close() error {
	log.Debug("closing the repo")
	if r.db != nil {
		return r.db.Close()
	}

	return nil
}

func (r *Repo) Keys() []string {
	vals := make([]string, 0)
	if r.db != nil {
		log.Debug("retrieving the repo keys")
		r.db.View(func(tx *bbolt.Tx) error {
			curs := tx.Bucket([]byte(r.bucket)).Cursor()
			for k, _ := curs.First(); k != nil; k, _ = curs.Next() {
				vals = append(vals, string(k))
			}

			return nil
		})
	}

	log.Debugf("retrieved %d keys", len(vals))
	return vals
}

func (r *Repo) Set(k, v string) error {
	if r.db == nil {
		return backend.ErrInvalidDatastore
	}

	log.Debug("writing provided key and value to repo")
	return r.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(r.bucket)).Put([]byte(k), []byte(v))
	})
}

func (r *Repo) Get(k string) string {
	if r.db == nil {
		return ""
	}

	log.Debug("retrieving value for provided key from repo")

	var val string
	r.db.View(func(tx *bbolt.Tx) error {
		val = string(tx.Bucket([]byte(r.bucket)).Get([]byte(k)))
		return nil
	})

	return val
}

func (r *Repo) Remove(k string) error {
	if r.db == nil {
		return backend.ErrInvalidDatastore
	}

	log.Debug("removing value for provided key from repo")
	return r.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(r.bucket)).Delete([]byte(k))
	})
}

func (r *Repo) Info() string {
	log.Debug("providing repo info string")
	return fmt.Sprintf("%s.%s", r.name, r.bucket)
}
