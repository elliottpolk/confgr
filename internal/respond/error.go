package respond

import (
	"fmt"
	"net/http"
)

func Error(w http.ResponseWriter, statuscode int, format string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(format, args...), statuscode)
}

func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
