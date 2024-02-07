package c2api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	cry "github.com/zarkones/xena-crypto"
)

func ListFiles() (files []File, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/files", nil)
	if err != nil {
		return nil, err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return files, json.Unmarshal(respBody, &files)
}

func DownloadFile(fileID string) (fileContent []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/files/"+fileID, nil)
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

	return respBody, nil
}

func UploadFile(fileID string, fileContent *[]byte) (err error) {
	if TrustedPubKey == nil {
		return ErrKeyIsNil
	}

	encryptedFileContent, err := cry.SecureWrap(TrustedPubKey, string(*fileContent))
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, *BaseURL+"/v1/files/"+fileID, bytes.NewBufferString(encryptedFileContent))
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrUnexpectedStatusCode
	}

	return nil
}