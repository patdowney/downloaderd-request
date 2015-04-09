package api

// IncomingDownload ...
type IncomingDownload struct {
	RequestID    string `json:"request_id"`
	URL          string `json:"url"`
	Checksum     string `json:"checksum",omitempty`
	ChecksumType string `json:"checksum_type",omitempty`
	Callback     string `json:"callback",omitempty`
	ETag         string `json:"etag",omitempty`
}
