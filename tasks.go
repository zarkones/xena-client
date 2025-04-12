package c2api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// GetTasks asks the C2 for the list of tasks.
func GetTasks() (tasks []Task, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/tasks", nil)
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

	if err := json.Unmarshal(respBody, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// UpsertTask will insert or update a pipeline.
func UpsertTask(task Task) (err error) {
	jsonPayload, err := json.Marshal(&task)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/tasks", bytes.NewBuffer(jsonPayload))
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
