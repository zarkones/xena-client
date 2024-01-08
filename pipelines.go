package c2api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// GetPipelines asks the C2 for the list of pipelines.
func GetPipelines() (pipelines []Pipeline, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/pipelines", nil)
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

	if err := json.Unmarshal(respBody, &pipelines); err != nil {
		return nil, err
	}

	return pipelines, nil
}

// GetPipelineRuns asks the C2 for the list of historical runs (executions) of a pipeline.
func GetPipelineRuns(pipelineID string) (runs []PipelineRun, err error) {
	req, err := http.NewRequest(http.MethodGet, *BaseURL+"/v1/pipelines/runs/"+pipelineID, nil)
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

	if err := json.Unmarshal(respBody, &runs); err != nil {
		return nil, err
	}

	return runs, nil
}

// UpsertPipeline will insert or update a pipeline.
func UpsertPipeline(pipeline Pipeline) (err error) {
	jsonPayload, err := json.Marshal(&pipeline)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/pipelines", bytes.NewBuffer(jsonPayload))
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

type ExecPipelineReqCtx struct {
	AgentIDs   []string `json:"agentIds"`
	PipelineID string   `json:"pipelineId"`
}

// ExecutePipeline will start a pipeline given its name.
// Each agent referenced by ID would execute the pipeline.
func ExecutePipeline(pipelineID string, agentIDs []string) (err error) {
	payload := ExecPipelineReqCtx{
		AgentIDs:   agentIDs,
		PipelineID: pipelineID,
	}

	jsonPayload, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/pipelines/exec", bytes.NewBuffer(jsonPayload))
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

type SetPipelineSettingsReqCtx struct {
	PipelineID       string           `json:"pipelineId"`
	PipelineSettings PipelineSettings `json:"pipelineSettings"`
}

// SetPipelineSettings will update the pipeline.
func SetPipelineSettings(pipelineID string, pipelineSettings PipelineSettings) (err error) {
	payload := SetPipelineSettingsReqCtx{
		PipelineID:       pipelineID,
		PipelineSettings: pipelineSettings,
	}

	jsonPayload, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *BaseURL+"/v1/pipelines/settings", bytes.NewBuffer(jsonPayload))
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

func DeletePipeline(pipelineID string) (err error) {
	req, err := http.NewRequest(http.MethodDelete, *BaseURL+"/v1/pipelines/"+pipelineID, nil)
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
