package backend

import "testing"

func TestDatastore(t *testing.T) {
	tests := map[string]func(t *testing.T){
		"new": func(t *testing.T) {

		},
	}

	// TODO:
	// - write a test Repo interface
	// - include proper tests for the testing repo

	for name, test := range tests {
		t.Run(name, test)
	}
}
