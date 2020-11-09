package backend

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Value map[string]string

type Repo interface {
	Open() error
	Close() error
	Keys() []string
	Set(key, value string) error
	Get(key string) string
	Remove(key string) error
	Info() string
}

type Datastore struct {
	primary Repo
	audit   Repo
}

func (v *Value) String() string {
	out, _ := json.Marshal(v)
	return string(out)
}

//  Keys iterates over the available keys and returns as a list for the primary
//  datastore.
func (ds *Datastore) Keys(c context.Context) []string {
	if err := ds.AddAuditEvent(NewEvent(ListKeysOp, ds.primary.Info(), "", c.Value("user").(string))); err != nil {
		log.Error(errors.Wrap(err, "unable to write audit event"))
	}

	log.Debug("retrieving primary datastore keys")
	return ds.primary.Keys()
}

//  Set adds a new entry into the key/value store. If the key exists, the old
//  value will be overwritten.
func (ds *Datastore) Set(c context.Context, k, v string) error {
	if ds.primary == nil {
		return ErrInvalidDatastore
	}

	ae := NewEvent(SetOp, ds.primary.Info(), (&Value{"key": k, "value": v}).String(), c.Value("user").(string))
	if err := ds.AddAuditEvent(ae); err != nil {
		log.Error(errors.Wrap(err, "unable to write audit event"))
	}

	log.Debug("writing key/value pair to datastore")
	if err := ds.primary.Set(k, v); err != nil {
		if err := ds.AddAuditEvent(ae.Fail()); err != nil {
			log.Error(errors.Wrap(err, "unable to write audit event"))
		}

		return err
	}

	return nil
}

//  Get retrieves the relevant content for the provided key.
func (ds *Datastore) Get(c context.Context, k string) string {
	if ds.primary == nil {
		return ""
	}

	ae := NewEvent(GetOp, ds.primary.Info(), (&Value{"key": k}).String(), c.Value("user").(string))
	if err := ds.AddAuditEvent(ae); err != nil {
		log.Error(errors.Wrap(err, "unable to write audit event"))
	}

	log.Debug("retrieving value for provided key from datastore")
	return ds.primary.Get(k)
}

//  Remove deletes the content for the provided key. No error is returned if the
//  provided key does not exist.
func (ds *Datastore) Remove(c context.Context, k string) error {
	if ds.primary == nil {
		return ErrInvalidDatastore
	}

	ae := NewEvent(RemoveOp, ds.primary.Info(), (&Value{"key": k}).String(), c.Value("user").(string))
	if err := ds.AddAuditEvent(ae); err != nil {
		log.Error(errors.Wrap(err, "unable to write audit event"))
	}

	log.Debug("removing value for provided key from datastore")
	if err := ds.primary.Remove(k); err != nil {
		if err := ds.AddAuditEvent(ae.Fail()); err != nil {
			log.Error(errors.Wrap(err, "unable to write audit event"))
		}

		return err
	}

	return nil
}
