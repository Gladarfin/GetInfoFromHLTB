package client

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

const (
	hltbBaseURL = "https://howlongtobeat.com"
)

type Client struct {
	httpClient *http.Client
	authToken  string
	searchURL  string
	userAgent  string
}

// New create new Client for HLTB
func New() *Client {
	jar, _ := cookiejar.New(nil)

	return &Client{
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: false,
				ForceAttemptHTTP2:  true,
			},
		},
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}
