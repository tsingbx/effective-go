//go:build !skip

package main

import (
	"bytes"
	"flag"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type testEnv struct {
	args           string
	stdout, stderr bytes.Buffer
}

func (e *testEnv) run() error {
	s := flag.NewFlagSet("hit", flag.ContinueOnError)
	s.SetOutput(&e.stderr)
	return run(s, strings.Fields(e.args), &e.stdout)
}

func TestRunSad(t *testing.T) {
	t.Parallel()
	sad := map[string]string{
		"url/missing": "",
		"url/err":     "://foo",
		"url/host":    "http://",
		"url/scheme":  "ftp://",
		"c/err":       "-c=x http://foo",
		"n/err":       "-n=x http://foo",
		"c/neg":       "-c=-1 http://foo",
		"n/neg":       "-n=-1 http://foo",
		"c/zero":      "-c=0 http://foo",
		"n/zero":      "-n=0 http://foo",
		"c/greater":   "-n=1 -c=2 http://foo",
		"t/zero":      "-t=0s http://foo",
		"t/neg":       "-t=-1s http://foo",
		"m/invalid":   "-m=go http://foo",
		"m/empty":     "-m=\"\" http://foo",
		"m/none":      "-m= http://foo",
		"H/empty":     "-H= http://foo",
		"H/invalid":   "-H=go http://foo",
	}
	for name, in := range sad {
		in := in
		t.Run(name+fmt.Sprintf("=>%q", in), func(t *testing.T) {
			t.Parallel()
			e := &testEnv{args: in}
			if e.run() == nil {
				t.Fatal("got nil; want err")
			}
			if e.stderr.Len() == 0 {
				t.Fatal("stderr = 0 bytes; want > 0")
			}
		})
	}
}

func TestRunHappy(t *testing.T) {
	t.Parallel()
	happy := map[string]struct {
		in, out string
	}{
		"url": {
			"http://foo",
			"100 requests to http://foo with a concurrency level of " + strconv.Itoa(runtime.NumCPU()),
		},
		"n_c": {
			"-n=20 -c=5 http://foo",
			"20 requests to http://foo with a concurrency level of 5",
		},
		"m=GET": {
			"-m=GET http://foo",
			"100 requests to http://foo with a concurrency level of " + strconv.Itoa(runtime.NumCPU()),
		},
		"m=PUT": {
			"-m=PUT http://foo",
			"100 requests to http://foo with a concurrency level of " + strconv.Itoa(runtime.NumCPU()),
		},
		"m=POST": {
			"-m=POST http://foo",
			"100 requests to http://foo with a concurrency level of " + strconv.Itoa(runtime.NumCPU()),
		},
		"t>0s": {
			"-t=5s http://foo",
			"100 requests to http://foo with a concurrency level of " + strconv.Itoa(runtime.NumCPU()),
		},
		"-H:Accept": {
			"-H='Accept:text/json' http://foo",
			"100 requests to http://foo with a concurrency level of " + strconv.Itoa(runtime.NumCPU()),
		},
		"-H:User-Agent": {
			"-H='User-Agent:hit' http://foo",
			"100 requests to http://foo with a concurrency level of " + strconv.Itoa(runtime.NumCPU()),
		},
		"-H multi set": {
			"-H=User-Agent:hit' -H='Accept:text/json' http://foo",
			"100 requests to http://foo with a concurrency level of " + strconv.Itoa(runtime.NumCPU()),
		},
	}
	for name, tt := range happy {
		tt := tt
		t.Run(name+fmt.Sprintf("=>%q", tt.in), func(t *testing.T) {
			t.Parallel()
			e := &testEnv{args: tt.in}
			if err := e.run(); err != nil {
				t.Fatalf("got %q;\nwant nil err", err)
			}
			if out := e.stdout.String(); !strings.Contains(out, tt.out) {
				t.Errorf("got:\n%s\nwant %q", out, tt.out)
			}
		})
	}
}
