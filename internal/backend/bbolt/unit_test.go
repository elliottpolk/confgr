package bbolt

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	// log.SetLevel(log.DebugLevel)
}

func prep(t *testing.T) (*Repo, func()) {
	what := fmt.Sprintf("%d_psparkles_testing", time.Now().UnixNano())

	r := NewRepo(what, PrimaryBucket)
	assert.NotNil(t, r)
	assert.NoError(t, r.Open())

	return r, func() {
		assert.NoError(t, r.Close())
		assert.NoError(t, os.RemoveAll(what))
	}
}

func TestOpen(t *testing.T) {
	r, cleanup := prep(t)
	defer cleanup()

	key, want := "foo", "bar"
	assert.NoError(t, r.Set(key, want))
	assert.Equal(t, want, r.Get(key))
}

func TestClose(t *testing.T) {
	r, cleanup := prep(t)
	defer cleanup()

	// Steps:
	// - write arbitrary data to ensure it's open
	// - close to later test if data will be retrieved
	// - test to make sure the above close worked

	key := "foo"
	assert.NoError(t, r.Set(key, "bar"))
	assert.NoError(t, r.Close())
	assert.Empty(t, r.Get(key))
}

func TestKeys(t *testing.T) {
	r, cleanup := prep(t)
	defer cleanup()

	assert.Empty(t, r.Keys())

	wants := []string{"foo", "bar", "baz"}
	for _, w := range wants {
		_r := rand.NewSource(time.Now().UnixNano()).Int63()
		assert.NoError(t, r.Set(w, fmt.Sprintf("%d", _r)))
	}

	keys := r.Keys()
	assert.Equal(t, len(wants), len(keys))

	sort.Strings(wants)
	sort.Strings(keys)
	assert.Equal(t, wants, keys)
}

func TestSetGet(t *testing.T) {
	r, cleanup := prep(t)
	defer cleanup()

	key := "foo"
	assert.Empty(t, r.Get(key))

	// set the initial value
	assert.NoError(t, r.Set(key, "bar"))
	assert.Equal(t, "bar", r.Get(key))

	// test an update of the value
	assert.NoError(t, r.Set(key, "baz"))
	assert.Equal(t, "baz", r.Get(key))
}

func TestRemove(t *testing.T) {
	r, cleanup := prep(t)
	defer cleanup()

	// Steps:
	// - write the record
	// - ensure it was written
	// - remove the record
	// - confirm removed

	key := "foo"
	assert.NoError(t, r.Set(key, "bar"))
	assert.Equal(t, "bar", r.Get(key))
	assert.NoError(t, r.Remove(key))
	assert.Empty(t, r.Get(key))
}
