package handler

import (
	"fmt"
	"net/http"
)

// TestHandler возвращает Metod, Host, Path
func TestHandler(w http.ResponseWriter, r *http.Request) {

	s := fmt.Sprintf("Method: %s\nHost: %s\nPath: %s",
		r.Method, r.Host, r.URL.Path)
	w.Write([]byte(s))
}
