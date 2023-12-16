package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

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
	f := &flags{}
	if err := f.parse(s, args); err != nil {
		return err
	}
	fmt.Fprintln(out, banner())
	fmt.Fprintf(out, "Making %d requests to %s with a concurrency level of %d.\n",
		f.n, f.url, f.c)
	if f.rps > 0 {
		fmt.Fprintf(out, "(RPS: %d)\n", f.rps)
	}

	ctx, cancel := context.WithTimeout(context.Background(), f.t)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	defer stop()
	sum, _ := hit.Do(ctx, f.url, f.n, hit.Concurrency(f.c), hit.Timeout(f.t), hit.RPS(f.rps))
	sum.Fprint(out)
	if err := ctx.Err(); errors.Is(err, context.DeadlineExceeded) {
		err = fmt.Errorf("error occurred: timed out in %s", f.t)
		fmt.Fprintf(out, "%v", err)
		return err
	}
	return nil
}
