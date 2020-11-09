package backend

import "github.com/pkg/errors"

var ErrInvalidDatastore error = errors.New("no valid datastore")
