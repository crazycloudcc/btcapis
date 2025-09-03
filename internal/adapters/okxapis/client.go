package okxapis

import (
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	base *url.URL
	http *http.Client
}

func New(baseURL string, timeout int) *Client {
	u, _ := url.Parse(baseURL)
	return &Client{
		base: u,
		http: &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}
