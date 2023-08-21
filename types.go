package c2api

type Agent struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
}

type Message struct {
	ID      string `json:"id"`
	Request string `json:"request"`
}
