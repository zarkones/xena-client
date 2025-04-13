package c2api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// GetChat asks the C2 for the list of chats.
func GetChat() (chat []ChatMessage, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/chat", nil)
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

	if err := json.Unmarshal(respBody, &chat); err != nil {
		return nil, err
	}

	return chat, nil
}

// UpsertChatMessage will insert or update a pipeline.
func UpsertChatMessage(chat *ChatMessage) (err error) {
	jsonPayload, err := json.Marshal(&chat)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/chat", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Join(ErrUnexpectedStatusCode, errors.New(resp.Status))
	}

	return nil
}
