package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
)

const (
	bannerText = ` 
__ __ __ ______
/\ \_\ \ /\ \ /\__ _\
\ \ __ \ \ \ \ \/_/\ \/
\ \_\ \_\ \ \_\ \ \_\
\/_/\/_/ \/_/ \/_/
`
)

func banner() string {
	return bannerText[1:]
}

func main() {
	if err := run(flag.CommandLine, os.Args[1:], os.Stdout); err != nil {
		os.Exit(1)
	}
}

func run(s *flag.FlagSet, args []string, out io.Writer) error {
	f := &flags{
		n: 100,
		c: runtime.NumCPU(),
	}
	if err := f.parse(s, args); err != nil {
		return err
	}
	fmt.Fprintln(out, banner())
	m := method(f.m)
	fmt.Fprintf(out, "Making %d requests to %s with a concurrency level of %d timeout %v, method:%s, headers: %s\n", f.n, f.url, f.c, f.t, m.String(), headersString())
	return nil
}
