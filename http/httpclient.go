package http

import (
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	useragent string
	client    *http.Client
	policy    FetchPolicy
}

// A ClientOption is a functional option.
type ClientOption func(*clientOptions)

type clientOptions struct {
	clientTimeout time.Duration
	useragent     string
	policy        FetchPolicy
}

// WithClientTimeout sets the client timeout.
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(opts *clientOptions) { opts.clientTimeout = timeout }
}

// WithUserAgent sets the client user-agent.
func WithUserAgent(useragent string) ClientOption {
	return func(opts *clientOptions) { opts.useragent = useragent }
}

// WithFetchPolicy applies a URL fetch policy. A zero policy keeps the
// historical client behavior.
func WithFetchPolicy(policy FetchPolicy) ClientOption {
	return func(opts *clientOptions) { opts.policy = policy }
}

// NewClient returns a new http.Client with a default timeout.
func NewClient(options ...ClientOption) *Client {
	opts := clientOptions{clientTimeout: 10 * time.Second}

	for i := range options {
		options[i](&opts)
	}
	client := &Client{
		client: &http.Client{Timeout: opts.clientTimeout},
		policy: opts.policy,
	}

	if opts.useragent != "" {
		client.useragent = opts.useragent
	}
	if client.policy.Active() {
		if client.policy.BlockPrivateNetworks {
			client.client.Transport = client.policy.transport()
		}
		client.client.CheckRedirect = client.policy.checkRedirect
	}
	return client
}

// ValidateURL checks u against the client's fetch policy.
func (c *Client) ValidateURL(u *url.URL) error {
	if c == nil {
		return nil
	}
	return c.policy.ValidateURL(u)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.useragent != "" {
		req.Header.Set("User-Agent", c.useragent)
	}
	return c.client.Do(req)
}
