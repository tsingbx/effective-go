package httpio

import (
	"net/http"

	"github.com/tsingbx/effective-go/ch7/linkit"
)

func JSON(code int, v any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := Encode(w, code, v); err != nil {
			Log(r.Context(), "%s: JSON.Encode: %v", r.URL.Path, err)
		}
	}
}

func Error(code int, message string) Handler {
	return func(_ http.ResponseWriter, r *http.Request) http.Handler {
		if code == http.StatusInternalServerError {
			Log(r.Context(), "%s: %v", r.URL.Path, message)
			message = linkit.ErrInternal.Error()
		}
		return JSON(code, map[string]string{
			"error": message,
		})
	}
}
