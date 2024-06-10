package msgs

import "time"

type PowerControl struct {
	Power bool `json:"power"`
}

type PowerStatus struct {
	Power bool   `json:"power"`
	Error string `json:"error"` // set if the last Control request failed
}

type StorageControl struct {
	Action string   `json:"action"` // activate, deactivate, status
	Images []string `json:"images"` // list of images to activate
}

type StorageStatus struct {
	Status           string        `json:"status"`
	Images           []string      `json:"images"`
	LastStatusChange time.Time     `json:"last_status_change"`
	Hostname         string        `json:"hostname"`
	Error            string        `json:"error"`
	ErrorAt          time.Time     `json:"error_at"`
	Uptime           time.Duration `json:"uptime"`
}

type Status struct {
	Status string `json:"status"`
}

type Firmware struct {
	Firmware string `json:"firmware"`
}
