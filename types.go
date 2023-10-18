package c2api

type Agent struct {
	ID        string `json:"id"`
	Hostname  string `json:"hostname"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	IpAddress string `json:"ipAddress"`
}

type Message struct {
	ID      string `json:"id"`
	Request string `json:"request"`
}

type Pipeline struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Desc     string `json:"description"`
	Category string `json:"category"`
	Settings string `json:"settings"` // JSON stringified "PipelineSettings" variable.
}

type PipelineSettings struct {
	Input map[string]string `json:"input"`
	Steps []PipelineCommand `json:"steps"`
}

type PipelineCommand struct {
	Name        string `json:"name"`
	Cmd         string `json:"cmd"`
	ToOutputTag string `json:"toOutput"`
}
