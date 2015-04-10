package download

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/patdowney/downloaderd-request/api"
)

// Client ...
type Client interface {
	ProcessRequest(*Request) (*Download, error)
}

// HTTPClient ...
type HTTPClient struct {
	URL *url.URL
}

// ProcessRequest ...
func (c *HTTPClient) ProcessRequest(r *Request) (*Download, error) {

	p, _ := json.MarshalIndent(r, "", "  ")
	log.Printf("p:%v", string(p))

	rr := api.IncomingDownload{
		RequestID:    r.ID,
		URL:          r.URL,
		Checksum:     r.Checksum,
		ChecksumType: r.ChecksumType,
		Callback:     r.Callback,
		ETag:         r.Metadata.ETag,
	}

	return c.postRequest(rr)
}

func (c *HTTPClient) postRequest(r api.IncomingDownload) (*Download, error) {
	p, _ := json.MarshalIndent(r, "", "  ")
	log.Printf("p:%v", string(p))

	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	byteReader := bytes.NewReader(jsonBytes)
	res, err := http.Post(c.URL.String(), "application/json", byteReader)
	defer func() {
		err := res.Body.Close()
		log.Printf("%v\n", err)
	}()

	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("postRequest: status: %v", res.StatusCode)
}

// NewHTTPClient ...
func NewHTTPClient(url *url.URL) (Client, error) {
	return &HTTPClient{URL: url}, nil
}
