package teams

import (
	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
)

type datastore struct {
	*bolt.DB
}

type meta struct {
	Created   int64  `json:"created"`
	CreatedBy string `json:"created_by,omitempty"`
	Updated   int64  `json:"updated"`
	UpdatedBy string `json:"updated_by,omitempty"`
	Status    string `json:"status"`
	Version   int    `json:"version"`
}

const (
	active   string = "active"
	inactive string = "inactive"
	archived string = "archived"
)

func open(name string, opts *bolt.Options) (*datastore, error) {
	db, err := bolt.Open(name, 0600, opts)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open datastore file")
	}

	return &datastore{db}, nil
}

func (ds *datastore) close() error {
	if ds.DB != nil {
		return ds.DB.Close()
	}
	return nil
}

func (ds *datastore) keys(b string) []string {
	vals := make([]string, 0)
	if ds.DB != nil {
		ds.View(func(tx *bolt.Tx) error {
			curs := tx.Bucket([]byte(b)).Cursor()
			for k, _ := curs.First(); k != nil; k, _ = curs.Next() {
				vals = append(vals, string(k))
			}

			return nil
		})
	}

	return vals
}

func (ds *datastore) ensureBucket(b string) error {
	return ds.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(b)); err != nil {
			return err
		}
		return nil
	})
}

func (ds *datastore) save(b, k, v string) error {
	if err := ds.ensureBucket(b); err != nil {
		return errors.Wrapf(err, "unable to ensure bucket '%s' exists", b)
	}

	return ds.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(b)).Put([]byte(k), []byte(v))
	})
}

func (ds *datastore) retrieve(b, k string) []byte {
	if err := ds.ensureBucket(b); err != nil {
		return nil
	}

	var val []byte
	ds.View(func(tx *bolt.Tx) error {
		val = tx.Bucket([]byte(b)).Get([]byte(k))
		return nil
	})

	return val
}

func (ds *datastore) del(b, k string) error {
	if err := ds.ensureBucket(b); err != nil {
		return errors.Wrapf(err, "unable to ensure bucket '%s' exists", b)
	}

	return ds.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(b)).Delete([]byte(k))
	})
}

func (ds *datastore) list(b string) []map[string][]byte {
	vals := make([]map[string][]byte, 0)
	if err := ds.ensureBucket(b); err != nil {
		return vals
	}

	for _, k := range ds.keys(b) {
		vals = append(vals, map[string][]byte{k: ds.retrieve(b, k)})
	}
	return vals
}
