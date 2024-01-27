package c2api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type AttackTargetCtx struct {
	AgentID  string `json:"agentId"`
	TargetID string `json:"targetId"`
}

func GetTargets() (targets []Target, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/targets", nil)
	if err != nil {
		return nil, err
	}

	setAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&targets); err != nil {
		return nil, err
	}

	return targets, nil
}

func UpsertTargets(target Target) (err error) {
	jsonTarget, err := json.Marshal(&target)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/targets", bytes.NewReader(jsonTarget))
	if err != nil {
		return err
	}

	setAuth(req)

	if _, err := c.Do(req); err != nil {
		return err
	}

	return nil
}

func RemoveTarget(targetID string) (err error) {
	req, err := http.NewRequest(http.MethodDelete, *BaseURL+"/v1/targets/"+targetID, nil)
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
