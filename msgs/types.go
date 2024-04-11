package msgs

type PowerControl struct {
	Power bool `json:"power"`
}

type PowerStatus struct {
	Power bool   `json:"power"`
	Error string `json:"error"` // set if the last Control request failed
}

type StorageControl struct {
	Active bool     `json:"active"`
	Images []string `json:"images"`
}

type StorageStatus struct {
	Status string   `json:"status"`
	Images []string `json:"images"`
	Error  string   `json:"error"`
}

type Status struct {
	Status string `json:"status"`
}
