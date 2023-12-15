package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/tsingbx/effective-go/ch5/hit"
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
	if f.rps > 0 {
		fmt.Fprintf(out, "(RPS: %d)\n", f.rps)
	}
	request, err := http.NewRequest(http.MethodGet, f.url, http.NoBody)
	if err != nil {
		return err
	}
	const timeout = 60 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	defer stop()
	client := &hit.Client{C: f.c, RPS: f.rps, Timeout: 10 * time.Second}
	sum := client.Do(ctx, request, f.n)
	sum.Fprint(out)
	if err := ctx.Err(); errors.Is(err, context.DeadlineExceeded) {
		err = fmt.Errorf("error occurred: timed out in %s", timeout)
		fmt.Fprintf(out, "%v", err)
		return err
	}
	return nil
}
