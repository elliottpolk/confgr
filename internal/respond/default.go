package respond

import (
	"fmt"
	"net/http"
)

func DefaultOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=us-ascii")
	fmt.Fprint(w, "ok")
}
