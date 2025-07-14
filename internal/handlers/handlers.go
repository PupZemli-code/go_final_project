package handler

import (
	"fmt"
	"net/http"
)

// HomeHandler обрабатывает корневой путь http://localhost:7540/ и отдает индексный файл
func TestHandler(w http.ResponseWriter, r *http.Request) {

	s := fmt.Sprintf("Method: %s\nHost: %s\nPath: %s",
		r.Method, r.Host, r.URL.Path)
	w.Write([]byte(s))
}
