package key

import (
	"crypto/sha256"
	"fmt"
)

func convert(args ...string) string {
	val := make([]byte, 0)
	for _, a := range args {
		val = append(val, []byte(a)...)
	}
	return fmt.Sprintf("%x", sha256.Sum256(val))
}

func Generate(app, env string, values ...string) string {
	return fmt.Sprintf("%s%s", convert(values...), convert(app, env))
}
