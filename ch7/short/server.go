package short

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tsingbx/effective-go/ch7/internal/httpio"
	"github.com/tsingbx/effective-go/ch7/linkit"
)

const maxKeyLen = 16

type link struct {
	uri      string
	shortKey string
}

func checkShortKey(k string) error {
	if strings.TrimSpace(k) == "" {
		return errors.New("empty key")
	}
	if len(k) > maxKeyLen {
		return fmt.Errorf("key too long (max %d)", maxKeyLen)
	}
	return nil
}

func checkLink(ln link) error {
	if err := checkShortKey(ln.shortKey); err != nil {
		return err
	}
	u, err := url.ParseRequestURI(ln.uri)
	if err != nil {
		return err
	}
	if u.Host == "" {
		return errors.New("empty host")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme must be http or https")
	}
	return nil
}

const (
	shorteningRoute  = "/s"
	resolveRoute     = "/r/"
	healthCheckRoute = "/health"
)

type mux http.Handler

type Server struct {
	mux
}

func NewServer() *Server {
	var s Server
	s.registerRoutes()
	return &s
}

func (s *Server) registerRoutes() {
	mux := http.NewServeMux()
	mux.Handle(shorteningRoute, httpio.Handler(s.shorteningHandler))
	mux.Handle(resolveRoute, httpio.Handler(s.resolveHandler))
	mux.Handle(healthCheckRoute, httpio.Handler(s.healthCheckHandler))
	s.mux = mux
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) http.Handler {
	return httpio.JSON(http.StatusOK, "OK")
}

func (s *Server) shorteningHandler(w http.ResponseWriter, r *http.Request) http.Handler {
	if r.Method != http.MethodPost {
		return httpio.Error(http.StatusMethodNotAllowed, "method not allowed")
	}
	var input struct {
		URL string
		Key string
	}
	err := httpio.Decode(http.MaxBytesReader(w, r.Body, 4_096), &input)
	if err != nil {
		return httpio.Error(http.StatusBadRequest, "cannot decode JSON")
	}
	ln := link{
		uri:      input.URL,
		shortKey: input.Key,
	}
	if err := checkLink(ln); err != nil {
		return httpio.Error(http.StatusBadRequest, err.Error())
	}
	_ = httpio.Encode(w, http.StatusCreated, map[string]any{"key": ln.shortKey})
	return httpio.JSON(http.StatusCreated, map[string]any{
		"key": ln.shortKey,
	})
}

func (s *Server) resolveHandler(w http.ResponseWriter, r *http.Request) http.Handler {
	key := r.URL.Path[len(resolveRoute):]
	if err := checkShortKey(key); err != nil {
		return httpio.Error(http.StatusBadRequest, err.Error())
	}
	if key == "fortesting" {
		return httpio.Error(http.StatusBadRequest, "db at IP ... failed")
	}
	if key != "go" {
		return httpio.Error(http.StatusNotFound, linkit.ErrNotExists.Error())
	}
	const uri = "http://go.dev"
	http.Redirect(w, r, uri, http.StatusFound)
	return nil
}
