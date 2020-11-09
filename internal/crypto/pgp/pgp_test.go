package pgp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPgp(t *testing.T) {
	var (
		t1 Crypter = []byte("YTZkNDUxMjQtMDk5Ny00NTc2LTg5YTUtMzVlZGExOTlmOTk0Cg==")
		t2 Crypter = []byte("NDM1MkM4NTItQjc1QS00NzJCLUI3RDktMTBFOEZDNkMzMzRFCg==")
	)

	//	generated at http://fillerama.io/
	const (
		filler1 = `Wow, you got that off the Internet? In my day, the Internet was
		only used to download pornography. No, she'll probably make me do it. Oh,
		I always feared he might run off like this. Why, why, why didn't I break his
		legs?
		`
		filler2 = `Anyone who laughs is a communist! But I've never been to the moon!
		You, a bobsleder!? That I'd like to see! Who are you, my warranty?! Is that a
		cooking show? Take me to your leader! That could be 'my' beautiful soul sitting
		naked on a couch. If I could just learn to play this stupid thing. Who said
		that? SURE you can die! You want to die?!
		`
	)

	cypher1, err := t1.Encrypt([]byte(filler1))
	assert.NoError(t, err)

	cypher2, err := t2.Encrypt([]byte(filler2))
	assert.NoError(t, err)

	// test encrypting
	t.Run("encrypting", func(t *testing.T) {
		for _, c := range [][]byte{cypher1, cypher2} {
			assert.NotEmpty(t, c)
			assert.NotEqual(t, filler1, string(c))
			assert.NotEqual(t, filler2, string(c))
		}

		assert.NotEqual(t, string(cypher1), string(cypher2))
	})

	// test decrypting
	t.Run("decrypting", func(t *testing.T) {
		// ensure a different token can not produce the original string
		_, err := t2.Decrypt(cypher1)
		assert.Error(t, err, ErrInvalidToken)

		p1, err := t1.Decrypt(cypher1)
		assert.NoError(t, err)
		assert.NotEmpty(t, p1)
		assert.Equal(t, filler1, string(p1))
		assert.NotEqual(t, filler2, string(p1))

		p2, err := t2.Decrypt(cypher2)
		assert.NoError(t, err)
		assert.NotEmpty(t, p2)
		assert.NotEqual(t, filler1, string(p2))
		assert.Equal(t, filler2, string(p2))
	})
}
