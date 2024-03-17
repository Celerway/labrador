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
	Active bool     `json:"active"`
	Images []string `json:"images"`
}

type StorageImage struct {
	Origin string `json:"origin"`
	Lun    string `json:"lun"`
	Size   int    `json:"size"`
	Error  string `json:"error"` // set if the last Control request failed
}
