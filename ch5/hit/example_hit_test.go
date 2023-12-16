package hit_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/tsingbx/effective-go/ch5/hit"
)

func ExampleDo() {
	timeout := time.Second * 15
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	defer stop()
	sum, _ := hit.Do(ctx, "http://localhost:80", 100000, hit.Concurrency(10), hit.Timeout(timeout), hit.RPS(2))
	sum.Fprint(os.Stdout)
	if err := ctx.Err(); errors.Is(err, context.DeadlineExceeded) {
		err = fmt.Errorf("error occurred: timed out in %s", timeout)
		fmt.Fprintf(os.Stdout, "%v", err)
	}
}
