package download

import (
	"fmt"
	"net/http"

	"github.com/patdowney/downloaderd-request/common"
)

// RequestService ...
type RequestService struct {
	Clock          common.Clock
	IDGenerator    IDGenerator
	requestStore   RequestStore
	downloadClient Client
}

// NewRequestService ...
func NewRequestService(requestStore RequestStore, downloadClient Client) *RequestService {
	s := RequestService{
		IDGenerator:    &UUIDGenerator{},
		Clock:          &common.RealClock{},
		requestStore:   requestStore,
		downloadClient: downloadClient}

	return &s
}

// ProcessNewRequest ...
func (s *RequestService) ProcessNewRequest(downloadRequest *Request) (*Request, error) {
	id, err := s.IDGenerator.GenerateID()
	if err != nil {
		return nil, err
	}

	downloadRequest.ID = id
	downloadRequest.TimeRequested = s.Clock.Now()

	m, err := GetMetadataFromHead(s.Clock.Now(), downloadRequest)
	if err != nil {
		downloadRequest.AddError(err, s.Clock.Now())
	} else {
		downloadRequest.Metadata = m
		if m.StatusCode == http.StatusOK {
			download, err := s.downloadClient.ProcessRequest(downloadRequest)
			if err != nil {
				downloadRequest.AddError(err, s.Clock.Now())
			}
			if download != nil {
				downloadRequest.DownloadID = download.ID
			}
		} else {
			err := fmt.Errorf("non-200 response from source")

			downloadRequest.AddError(err, s.Clock.Now())
		}
	}

	err = s.requestStore.Add(downloadRequest)
	if err != nil {
		downloadRequest.AddError(err, s.Clock.Now())
	}

	return downloadRequest, err
}

// ListAll ...
func (s *RequestService) ListAll() ([]*Request, error) {
	return s.requestStore.FindAll(0, 100)
}

// FindByID ...
func (s *RequestService) FindByID(id string) (*Request, error) {
	return s.requestStore.FindByID(id)
}
