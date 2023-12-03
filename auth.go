package c2api

import "net/http"

func setAuth(req *http.Request) {
	if AuthToken == nil {
		return
	}
	if *AuthToken == "" {
		return
	}
	req.Header.Add("Authorization", *AuthToken)
}
