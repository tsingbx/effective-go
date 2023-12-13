package main

import "errors"

const (
	GET = 1 + iota
	PUT
	POST
)

// method is one of http methods(GET,POST,PUT)
type method int

func toMethod(p *int) *method {
	return (*method)(p)
}

func (m *method) Set(s string) error {
	switch {
	case s == "GET":
		*m = GET
	case s == "PUT":
		*m = PUT
	case s == "POST":
		*m = POST
	default:
		return errors.New("invalid method")
	}
	return nil
}

func (m *method) String() string {
	v := int(*m)
	switch v {
	case GET:
		return "GET"
	case PUT:
		return "PUT"
	case POST:
		return "POST"
	}
	return ""
}
