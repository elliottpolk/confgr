package peppermintsparkles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecret(t *testing.T) {
	wants := []struct {
		want   string
		secret *Secret
	}{
		{
			want: `{"id":"7604c73e-7595-4215-ae37-c99a7e464eb1","content":"notSuperS3cret","app_info":{"app_name":"testing","env":"test"},"metadata":{"team":"sample.test.team","username":"user"}}`,
			secret: &Secret{
				Id:      "7604c73e-7595-4215-ae37-c99a7e464eb1",
				Content: "notSuperS3cret",
				App:     &AppInfo{"testing", "test"},
				Metadata: &Meta{
					Team: "sample.test.team",
					User: "user",
				},
			},
		},
		{
			want: `{"id":"0623c465-5c3b-478b-bb9f-4ddef703454a","content":"notSuperS3cret","app_info":{"app_name":"testing","env":"test"},"metadata":{"team":"sample.test.team","username":"user","expires":1604936854032347000}}`,
			secret: &Secret{
				Id:      "0623c465-5c3b-478b-bb9f-4ddef703454a",
				Content: "notSuperS3cret",
				App:     &AppInfo{"testing", "test"},
				Metadata: &Meta{
					Team:    "sample.test.team",
					User:    "user",
					Expires: 1604936854032347000,
				},
			},
		},
		{
			want: `{"id":"1019074c-9655-4d54-9dd2-202cdd2f28d0","content":{"password":"not@ReallyS3cret!","username":"bob"},"app_info":{"app_name":"testing","env":"test"},"metadata":{"team":"sample.test.team","username":"user"}}`,
			secret: &Secret{
				Id: "1019074c-9655-4d54-9dd2-202cdd2f28d0",
				Content: map[string]interface{}{
					"username": "bob",
					"password": "not@ReallyS3cret!",
				},
				App: &AppInfo{"testing", "test"},
				Metadata: &Meta{
					Team: "sample.test.team",
					User: "user",
				},
			},
		},
		{
			want: `{"id":"d0b0ecf5-71e0-4e03-aba7-c6d6d76caf88","content":null,"app_info":null,"metadata":null}`,
			secret: &Secret{
				Id: "d0b0ecf5-71e0-4e03-aba7-c6d6d76caf88",
			},
		},
	}

	for _, want := range wants {
		assert.Equal(t, want.want, want.secret.String())
	}
}
