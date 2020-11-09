// +build docker

package redis

import (
	"fmt"
	"math/rand"
	"os/exec"
	"sort"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func prep(t *testing.T) (*Repo, func()) {
	var (
		name = fmt.Sprintf("%d_redis", time.Now().UnixNano())
		port = port()
	)

	n := fmt.Sprintf("--name=%s", name)
	p := fmt.Sprintf("-p=%s:6379", port)

	log.Debugf("booting redis docker container %s and exposing port %s", name, port)
	assert.NoError(t, exec.Command("docker", "run", "-d", n, p, "redis:alpine").Run())

	r := NewRepo(fmt.Sprintf("localhost:%s", port), "", PrimaryBucket)
	assert.NotNil(t, r)
	assert.NoError(t, r.Open())

	return r, func() {
		log.Debugf("stopping redic docker container %s", name)
		assert.NoError(t, exec.Command("docker", "stop", name).Run())

		log.Debugf("removing redis docker container %s", name)
		assert.NoError(t, exec.Command("docker", "rm", name).Run())
	}
}

func port() string {
	rand.Seed(time.Now().UnixNano())
	var min, max int64 = 2000, 9999
	return fmt.Sprintf("%d", rand.Int63n(max-min)+min)
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
