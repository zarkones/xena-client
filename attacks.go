package c2api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func GetOngoingAttacks() (attacks []Attack, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/attacks", nil)
	if err != nil {
		return nil, err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&attacks); err != nil {
		return nil, err
	}

	return attacks, nil
}

func AttackTarget(agentId string, targetId string) (err error) {
	payload := AttackTargetCtx{
		AgentID:  agentId,
		TargetID: targetId,
	}

	jsonPayload, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/targets/attack", bytes.NewReader(jsonPayload))
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
