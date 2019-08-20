package teams

import (
	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
)

type datastore struct {
	*bolt.DB
}

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

func (ds *datastore) save(b, k, v string) error {
	return ds.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(b)).Put([]byte(k), []byte(v))
	})
}

func (ds *datastore) retrieve(b, k string) []byte {
	var val []byte
	ds.View(func(tx *bolt.Tx) error {
		val = tx.Bucket([]byte(b)).Get([]byte(k))
		return nil
	})

	return val
}

func (ds *datastore) del(b, k string) error {
	return ds.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(b)).Delete([]byte(k))
	})
}

func (ds *datastore) list(b string) []map[string][]byte {
	vals := make([]map[string][]byte, 0)
	for _, k := range ds.keys(b) {
		vals = append(vals, map[string][]byte{k: ds.retrieve(b, k)})
	}
	return vals
}
