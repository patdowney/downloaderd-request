package download

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
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

	return nil, fmt.Errorf("not implemented yet")
}

// NewHTTPClient ...
func NewHTTPClient(url *url.URL) (Client, error) {
	return &HTTPClient{URL: url}, nil
}
