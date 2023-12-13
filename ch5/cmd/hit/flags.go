package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"
)

type flags struct {
	url  string
	n, c int
	t    time.Duration
	m    int
	H    string
}

const usageText = `
Usage:
	hit [options] url
Options:`

func validateURL(s string) error {
	u, err := url.Parse(s)
	switch {
	case strings.TrimSpace(s) == "":
		err = errors.New("required")
	case err != nil:
		err = errors.New("parse error")
	case u.Scheme != "http":
		err = errors.New("only supported scheme is http")
	case u.Host == "":
		err = errors.New("missing host")
	}
	return err
}

func (f *flags) validate(fs *flag.FlagSet) error {
	if err := validateURL(f.url); err != nil {
		return fmt.Errorf("url: %w", err)
	}
	if f.c > f.n {
		return fmt.Errorf("-c=%d: should be less than or equal to -n=%d", f.c, f.n)
	}

	exists := make([]string, 15)
	fs.Visit(func(f *flag.Flag) {
		exists = append(exists, f.Name)
	})

	i := slices.Index(exists, "m")
	if i != -1 && f.m != GET && f.m != PUT && f.m != POST {
		m := method(f.m)
		return fmt.Errorf("-m=%v: should be one of the valid http methods(GET,PUT,POST)", m.String())
	}

	j := slices.Index(exists, "t")
	if j != -1 && f.t.Nanoseconds() <= 0 {
		return fmt.Errorf("-t=%v: should be > 0", f.t)
	}

	k := slices.Index(exists, "H")
	if k != -1 && isValidHead(f.H) {
		return fmt.Errorf("-H=%v: invalid http head", f.H)
	}

	return nil
}

func (f *flags) parse(s *flag.FlagSet, args []string) error {
	s.Usage = func() {
		fmt.Fprintln(s.Output(), usageText[1:])
		s.PrintDefaults()
	}

	s.Var(toNumber(&f.n), "n", "Number of requests to make")
	s.Var(toNumber(&f.c), "c", "Concurrency level")
	s.Var(toMethod(&f.m), "m", "Should be GET,PUT,POST")
	s.Var(toReqTimeout(&f.t), "t", "Timeout should be > 0s")
	s.Var(toHttpHead(&f.H), "H", "Http head shoud valid")

	if err := s.Parse(args); err != nil {
		return err
	}

	f.url = s.Arg(0)

	if err := f.validate(s); err != nil {
		fmt.Fprintln(s.Output(), err)
		s.Usage()
		return err
	}
	return nil
}
