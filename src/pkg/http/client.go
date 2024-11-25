package http

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	clientTimeout   = 10 * time.Second
	clientUserAgent = "cloudstatus-mon/0.1"
)

//nolint:mnd,gochecknoglobals // Magic numbers in a global variable :(
var clientTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: defaultTransportDialContext(&net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 30 * time.Second,
	}),
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          10,
	IdleConnTimeout:       30 * time.Second,
	TLSHandshakeTimeout:   5 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

type Client struct {
	client *http.Client
}

//nolint:wrapcheck // Should be understood and wrapped by the caller.
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func NewRequest(
	ctx context.Context,
	method, url string,
	headers map[string]string,
	body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", clientUserAgent)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Transport: clientTransport,
			Timeout:   clientTimeout,
		},
	}
}

func NewClientFromClient(client *http.Client) *Client {
	client.Timeout = clientTimeout
	return &Client{
		client: client,
	}
}

//nolint:wrapcheck // Should be understood and wrapped by the caller.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *Client) GetString(ctx context.Context, url string, headers map[string]string) (string, error) {
	body, err := c.get(ctx, url, headers)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func GetJSON[T any](ctx context.Context, c *Client, url string, headers map[string]string) (*T, error) {
	body, err := c.get(ctx, url, headers)
	if err != nil {
		return nil, err
	}

	res := new(T)
	if err := json.Unmarshal(body, res); err != nil {
		//nolint:wrapcheck // Should be understood and wrapped by the caller.
		return nil, err
	}

	return res, nil
}

//nolint:wrapcheck // Should be understood and wrapped by the caller.
func (c *Client) get(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	req, err := NewRequest(ctx, "GET", url, headers, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, NewError(res)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
