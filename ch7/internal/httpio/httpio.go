package httpio

import (
	"encoding/json"
	"io"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) http.Handler

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if next := h(w, r); next != nil {
		next.ServeHTTP(w, r)
	}
}

func Decode(r io.Reader, v any) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func Encode(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
