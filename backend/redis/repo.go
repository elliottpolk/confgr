package redis

import (
	"fmt"

	"github.com/elliottpolk/peppermint-sparkles/internal/backend"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	PrimaryBucket int = 0
	AuditBucket   int = 1
)

type Repo struct {
	addr   string
	passwd string
	db     int
	client *redis.Client
}

func NewRepo(a, p string, db int) *Repo {
	return &Repo{
		addr:   a,
		passwd: p,
		db:     db,
	}
}

func (r *Repo) Open() error {
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.addr,
		Password: r.passwd,
		DB:       r.db,
	})

	//	ensure a valid connection prior to returning
	log.Debug("validating redis connection")
	cr, err := r.client.Ping().Result()
	if err != nil {
		return errors.Wrap(err, "unable to ping redis datastore")
	}

	log.Debugf("client ping result %s", cr)
	return nil
}

//  Close attempts to close the redis connection of the datastore
func (r *Repo) Close() error {
	log.Debug("closing the repo")
	if r.client != nil {
		return r.client.Close()
	}

	return nil
}

func (r *Repo) Keys() []string {
	vals := make([]string, 0)
	if r.client != nil {
		log.Debug("retrieving the repo keys")
		keys, err := r.client.Keys("*").Result()
		if err != nil {
			return []string{}
		}

		vals = make([]string, len(keys))
		copy(vals, keys)
	}

	log.Debugf("retrieved %d keys", len(vals))
	return vals
}

func (r *Repo) Set(k, v string) error {
	if r.client == nil {
		return backend.ErrInvalidDatastore
	}

	log.Debug("writing provided key and value to repo")
	return r.client.Set(k, v, 0).Err()
}

func (r *Repo) Get(k string) string {
	if r.client == nil {
		return ""
	}

	log.Debug("retrieving value for given key from repo")
	res, err := r.client.Get(k).Result()
	if err != nil && err != redis.Nil {
		//	log error but still return the empty string
		log.Error(errors.Wrap(err, "unable to retrieve result for key"))
		return ""
	}

	return res
}

func (r *Repo) Remove(k string) error {
	if r.client == nil {
		return backend.ErrInvalidDatastore
	}

	return r.client.Del(k).Err()
}

func (r *Repo) Info() string {
	log.Debug("providing repo info string")
	return fmt.Sprintf("%s.%d", r.addr, r.db)
}
