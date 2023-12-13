package hit

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Result is a request's result
type Result struct {
	RPS      float64       // RPS is the requests per second
	Requests int           // Requests is the number of requests
	Errors   int           // Errors is the number of errors occur
	Bytes    int64         // Bytes is the number of bytes downloaded
	Duration time.Duration // Duration is a single or all request
	Fatest   time.Duration // Fatest request result duration among
	Slowest  time.Duration // Slowest request result duration among
	Status   int           // Status is a request's HTTP status code
	Error    error         // Error is not nil if the request is failed
}

func (r *Result) Merge(o *Result) {
	r.Requests++
	r.Bytes += o.Bytes
	if r.Fatest == 0 || o.Duration < r.Fatest {
		r.Fatest = o.Duration
	}
	if o.Duration > r.Slowest {
		r.Slowest = o.Duration
	}
	switch {
	case o.Error != nil:
		fallthrough
	case o.Status >= http.StatusBadRequest:
		r.Errors++
	}
}

// Finalize the total duration and calculate RPS.
func (r *Result) Finalize(total time.Duration) *Result {
	r.Duration = total
	r.RPS = float64(r.Requests) / total.Seconds()
	return r
}

// Fprint the result to an io.Writer
func (r *Result) Fprint(out io.Writer) {
	p := func(format string, args ...any) {
		fmt.Fprintf(out, format, args...)
	}
	p("\nSummary:\n")
	p("\tSuccess	: %.0f%%\n", r.success())
	p("\tRPS		: %.1f\n", r.RPS)
	p("\tRequests 	: %d\n", r.Requests)
	p("\tErrors 	: %d\n", r.Errors)
	p("\tBytes 		: %d\n", r.Bytes)
	p("\tDuration 	: %s\n", round(r.Duration))
	if r.Requests > 1 {
		p("\tFatest 	: %s\n", round(r.Fatest))
		p("\tSlowest 	: %s\n", round(r.Slowest))
	}
}

func (r *Result) success() float64 {
	rr, e := float64(r.Requests), float64(r.Errors)
	return (rr - e) / rr * 100
}

func round(t time.Duration) time.Duration {
	return t.Round(time.Microsecond)
}
