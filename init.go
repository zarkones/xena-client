package c2api

import (
	"net/http"
	"time"
)

var BaseURL string
var authToken *string

var c *http.Client

func Init(c2Host string, timeout time.Duration) {
	BaseURL = c2Host
	c = &http.Client{
		Timeout: timeout,
	}
}

func SetAuthToken(newAuthToken *string) {
	authToken = newAuthToken
}
