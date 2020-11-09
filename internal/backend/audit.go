package backend

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// AuditKeys returns a list of available keys for the audit repo
// WARNING -- this does not have pagination/ranging so it could
// and could return a large list
func (ds *Datastore) AuditKeys() []string {
	return ds.audit.Keys()
}

// AuditRecord returns a string representation of audit record
func (ds *Datastore) AuditRecord(k string) string {
	return ds.audit.Get(k)
}

// AddAuditEvent inserts the event to the audit repo
func (ds *Datastore) AddAuditEvent(e Event) error {
	if ds.audit == nil {
		return ErrInvalidDatastore
	}
	log.Debugf("prepping to add a new audit event: %+v", e)

	log.Debug("reading in random data to generate head of event key")
	buf := make([]byte, 2048)
	if _, err := rand.Read(buf); err != nil {
		return errors.Wrap(err, "unable to read in random data to generate key")
	}

	// SHA256 together the content of the event + the previously collected random data to
	// genearte a psuedo random key to attempt and prevent key collisions
	log.Debug("writing new event log to datastore")
	return ds.audit.Set(fmt.Sprintf("%x", sha256.Sum256(append([]byte(e.String()), buf...))), e.String())
}
