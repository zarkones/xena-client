package c2api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// FetchMessages will reach out to C2 server and fetch messages.
func FetchMessages(agentID string) (messages []Message, err error) {
	req, err := http.NewRequest("GET", *BaseURL+"/v1/messages/"+agentID, nil)
	if err != nil {
		return nil, err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// AgentFetchMessages will reach out to C2 server and fetch messages to which it has not reponded.
func AgentFetchMessages(agentID string) (messages []Message, err error) {
	req, err := http.NewRequest("GET", *BaseURL+"/v1/messages/"+agentID, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// AgentRespondToMessage allows an agent to respond to a message.
func AgentRespondToMessage(messageID, response string) (err error) {
	var msgRespCtx struct {
		MessageID string `json:"messageId"`
		Response  string `json:"response"`
	}
	msgRespCtx.MessageID = messageID
	msgRespCtx.Response = response

	jsonMsgResp, err := json.Marshal(&msgRespCtx)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/respond", bytes.NewBuffer(jsonMsgResp))
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

// InsertMessage allows an operator to insert a message.
func InsertMessage(message Message) (err error) {
	jsonMsgResp, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/messages", bytes.NewBuffer(jsonMsgResp))
	if err != nil {
		return err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return ErrUnexpectedStatusCode
	}

	return nil
}
