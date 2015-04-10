package local

import (
	"sync"

	"github.com/patdowney/downloaderd-common/local"
	"github.com/patdowney/downloaderd-request/download"
)

// RequestStore ...
type RequestStore struct {
	local.JSONStore
	sync.RWMutex
	repository []*download.Request
}

// NewRequestStore ...
func NewRequestStore(dataFile string) (*RequestStore, error) {
	requestStore := &RequestStore{
		repository: make([]*download.Request, 0)}

	requestStore.DataFile = dataFile

	err := requestStore.LoadFromDisk(&requestStore.repository)

	return requestStore, err
}

// Add ...
func (s *RequestStore) Add(request *download.Request) error {
	s.Lock()
	defer s.Unlock()
	s.repository = append(s.repository, request)

	err := s.SaveToDisk(s.repository)

	return err
}

// FindByID ...
func (s *RequestStore) FindByID(requestID string) (*download.Request, error) {
	s.RLock()
	defer s.RUnlock()
	for _, request := range s.repository {
		if request.ID == requestID {
			return request, nil
		}
	}
	return nil, nil
}

// FindByResourceKey ...
func (s *RequestStore) FindByResourceKey(resourceKey download.ResourceKey, offset uint, count uint) ([]*download.Request, error) {
	s.RLock()
	defer s.RUnlock()
	results := make([]*download.Request, 0, len(s.repository))
	for _, request := range s.repository {
		if request.ResourceKey() == resourceKey {
			results = append(results, request)
		}
	}
	return results, nil
}

// FindAll ...
func (s *RequestStore) FindAll(offset uint, count uint) ([]*download.Request, error) {
	s.RLock()
	defer s.RUnlock()

	tmpRepository := make([]*download.Request, len(s.repository), len(s.repository))
	copy(tmpRepository, s.repository)

	return tmpRepository, nil
}
