package c2api

import "fyne.io/fyne/v2"

type Attack struct {
	ID       string `json:"id"`
	AgentID  string `json:"agentId"`
	TargetID string `json:"targetId"`
	Comment  string `json:"comment"`
}

type Agent struct {
	ID        string `json:"id"`
	Hostname  string `json:"hostname"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	IpAddress string `json:"ipAddress"`
}

type Message struct {
	ID            string `json:"id"`
	AgentID       string `json:"agentId"`
	FriendlyTitle string `json:"friendlyTitle"`
	Request       string `json:"request"`
	Response      string `json:"response"`
}

type Pipeline struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Desc     string `json:"description"`
	Category string `json:"category"`
	Settings string `json:"settings"` // JSON stringified "PipelineSettings" variable.
}

type PipelineSettings struct {
	Input map[string]string       `json:"input"`
	Steps map[string]PipelineStep `json:"steps"`
}

type PipelineStep struct {
	ID          string        `json:"id"`
	Position    fyne.Position `json:"position"`
	Name        string        `json:"name"`
	Cmd         string        `json:"cmd"`
	ToOutputTag string        `json:"toOutput"`
	Tool        Tool          `json:"tool"`
	LinkedTo    string        `json:"linkedTo"` // ID of a step it is linked towards.
}

type Tool struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	ToolCategoryName string                `json:"toolCategoryName"`
	Inputs           map[string]ToolInput  `json:"inputs"`
	Outputs          map[string]ToolOutput `json:"outputs"`
}

type ToolInput struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

type ToolOutput struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// Refers to ToolInput.Type and ToolOutput.Type
const TOOL_IO_TYPE_STRING = "STRING"

type Target struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
