package key

import "testing"

func TestGenerate(t *testing.T) {
	wants := map[string][][]string{
		"97df3588b5a3f24babc3851b372f0ba71a9dcdded43b14b9d06961bfc1707d9d632d5637a51e187ff57022d960706d4213a4113e6bd72f62b61c09262ae1148d": {
			{"foo", "bar", "baz"},
		},
		"ebf3b019bb7e36bdc0fbc4159345c04af54193ccf43ae4572922f6d4aa94bd5b632d5637a51e187ff57022d960706d4213a4113e6bd72f62b61c09262ae1148d": {
			{"foo", "bar", "baz", "biz"},
		},
		"2c60dbf3773104dce76dfbda9b82a729e98a42a7a0b3f9bae5095c7bed752b90632d5637a51e187ff57022d960706d4213a4113e6bd72f62b61c09262ae1148d": {
			{"foo", "bar", "bazz"},
			{"foo", "bar", "baz", "z"},
		},
		"796362b8b4289fca4d666ab486487d6699e828f9c098fc1c91566c291ef682f6632d5637a51e187ff57022d960706d4213a4113e6bd72f62b61c09262ae1148d": {
			{"foo", "bar", "baz z"},
			{"foo", "bar", "baz", " z"},
		},
	}

	var (
		app = "test"
		env = "testing"
	)

	for want, vals := range wants {
		for _, val := range vals {
			if got := Generate(app, env, val...); got != want {
				t.Errorf("\nwant %s\ngot %s\nwith %+v", want, got, val)
			}
		}
	}
}
