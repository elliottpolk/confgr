//  +build unit

// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/backend/redis/unit_test.go
//
package redis

import "testing"

func TestToKey(t *testing.T) {
	wants := map[string][][]string{
		"97df3588b5a3f24babc3851b372f0ba71a9dcdded43b14b9d06961bfc1707d9d": {{"foo", "bar", "baz"}},
		"2c60dbf3773104dce76dfbda9b82a729e98a42a7a0b3f9bae5095c7bed752b90": {{"foo", "bar", "bazz"}, {"foo", "bar", "baz", "z"}},
		"796362b8b4289fca4d666ab486487d6699e828f9c098fc1c91566c291ef682f6": {{"foo", "bar", "baz", " z"}},
	}

	ds := &Datastore{}
	for want, vals := range wants {
		for _, val := range vals {
			if got := ds.ToKey(val...); want != got {
				t.Errorf("\nwant: %s\ngot: %s", want, got)
			}
		}
	}
}
