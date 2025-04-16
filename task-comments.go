package c2api

import (
	"encoding/json"
	"io"
	"net/http"
)

// GetComments asks the C2 for the comments of tasks.
func GetComments(taskID string) (tasks []TaskComment, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/tasks/comments/"+taskID, nil)
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
