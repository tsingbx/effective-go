package hit

import (
	"context"
	"net/http"
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
		p = throttle(p, time.Second/time.Duration(c.RPS*c.C))
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
			MaxIdleConnsPerHost: c.C,
		},
		Timeout: c.Timeout,
	}
}
