package api

// IncomingDownload ...
type IncomingDownload struct {
	RequestID    string `json:"request_id"`
	URL          string `json:"url"`
	Checksum     string `json:"checksum"`
	ChecksumType string `json:"checksum_type"`
	Callback     string `json:"callback"`
	ETag         string `json:"etag"`
}
