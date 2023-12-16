package hit

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// Client sends HTTP requests and returns an aggregated performance
// result. The fields should not be changed after initializing.
type Client struct {
	C       int           // C is the concurrency level
	RPS     int           // RPS throttles the requests per second
	Timeout time.Duration // Timeout per request
}

// Do sends n HTTP requests and returns an aggregated result.
func (c *Client) Do(ctx context.Context, r *http.Request, n int) *Result {
	t := time.Now()
	sum := c.do(ctx, r, n)
	return sum.Finalize(time.Since(t))
}

func (c *Client) do(ctx context.Context, r *http.Request, n int) *Result {
	p := produce(ctx, n, func() *http.Request {
		return r.Clone(ctx)
	})

	if c.RPS > 0 {
		p = throttle(ctx, p, time.Second/time.Duration(c.RPS*c.C))
	}

	var (
		sum    Result
		client = c.client()
	)
	defer client.CloseIdleConnections()
	for result := range split(p, c.C, c.send(client)) {
		sum.Merge(result)
	}
	return &sum
}

// SendFunc is the type of the function called by Client.Do
// to send an HTTP request and return a performance result.
type SendFunc func(*http.Request) *Result

func (c *Client) send(client *http.Client) SendFunc {
	return func(r *http.Request) *Result {
		return Send(client, r)
	}
}

func (c *Client) client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: c.concurrency(),
		},
		Timeout: c.Timeout,
	}
}

func (c *Client) concurrency() int {
	if c.C > 0 {
		return c.C
	}
	return runtime.NumCPU()
}

// Option allows changes Client's behavior.
type Option func(*Client) Option

// Concurrency changes the Client's concurrency level.
func Concurrency(n int) Option {
	return func(c *Client) Option {
		prev := c.C
		c.C = n
		return Concurrency(prev)
	}
}

func RPS(rps int) Option {
	return func(c *Client) Option {
		prev := c.RPS
		c.RPS = rps
		return RPS(prev)
	}
}

// Timeout changes the Client's timeout per request.
func Timeout(d time.Duration) Option {
	return func(c *Client) Option {
		prev := c.Timeout
		c.Timeout = d
		return Timeout(prev)
	}
}

func Do(ctx context.Context, url string, n int, opts ...Option) (*Result, error) {
	var c Client
	for _, o := range opts {
		o(&c)
	}
	r, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("new http request: %w", err)
	}
	return c.Do(ctx, r, n), nil
}
