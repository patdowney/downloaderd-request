package api

// IncomingRequest ...
type IncomingRequest struct {
	URL          string `json:"url"`
	Checksum     string `json:"checksum,omitempty"`
	ChecksumType string `json:"checksum_type,omitempty"`
	Callback     string `json:"callback,omitempty"`
}
