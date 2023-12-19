package short

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func jsonReader(t *testing.T, v any) *bytes.Reader {
	t.Helper()
	body, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewReader(body)
}

func TestShortening(t *testing.T) {
	t.Parallel()
	reader := jsonReader(t, map[string]any{
		"url": "https://go.dev",
		"key": "go",
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, shorteningRoute, reader)

	srv := NewServer()
	srv.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("got status code = %d, want %d", w.Code, http.StatusCreated)
	}
	if !strings.Contains(w.Body.String(), `"go"`) {
		t.Errorf("got body = %s\twant contains %s", w.Body.String(), `"go"`)
	}
}
