package c2api

import "net/http"

func setAuth(req *http.Request) {
	if *authToken != "" {
		req.Header.Add("Authorization", *authToken)
	}
}
