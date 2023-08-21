package c2api

import (
	"net/http"
	"time"
)

var baseURL string

var c *http.Client

func Init(c2Host string, timeout time.Duration) {
	baseURL = c2Host
	c = &http.Client{
		Timeout: timeout,
	}
}
