package c2api

import (
	"errors"
	"io"
	"net/http"
)

var ErrKeyIsNil = errors.New("key is nil")
var ErrTrustedKeyIsNil = errors.New("trusted public key is not set")

func GetC2PublicKey() (pubKeyPEM []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/public-key", nil)
	if err != nil {
		return nil, err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}
