package teams

import (
	"io/ioutil"
	"os"
	"testing"

	bolt "github.com/coreos/bbolt"
	"github.com/stretchr/testify/assert"
)

/*
func TestGenKey(t *testing.T) {
	generateKey()
}
*/

func TestString(t *testing.T) {
	want := `{
 "id": "",
 "team": "testing",
 "username": "tester",
 "key": null,
 "meta": null
}`
	got := (&Profile{
		Team:     "testing",
		Username: "tester",
	}).String()

	assert.Equal(t, want, got)
}

func TestWrite(t *testing.T) {
	tmp, err := ioutil.TempFile(".", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp.Name())

	ds, err := open(tmp.Name(), bolt.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}
	defer ds.close()

	p := &Profile{
		Team:     "testing",
		Username: "tester",
	}

	if err := p.Write(ds); err != nil {
		t.Fatal(err)
	}

	// retrieve direct from datastore to ensure what is actually written
	assert.Equal(t, p.String(), string(ds.retrieve(BucketProfiles, p.Id)))
}

func TestGetProfile(t *testing.T) {
	tmp, err := ioutil.TempFile(".", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp.Name())

	ds, err := open(tmp.Name(), bolt.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}
	defer ds.close()

	p := &Profile{
		Team:     "testing",
		Username: "tester",
	}

	if err := p.Write(ds); err != nil {
		t.Fatal(err)
	}

	got, err := GetProfile(ds, p.Id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, p.String(), got.String())
}
