package c2api

import "net/http"

func setAuth(req *http.Request) {
	if authToken == nil {
		return
	}
	if *authToken == "" {
		return
	}
	req.Header.Add("Authorization", *authToken)
}
