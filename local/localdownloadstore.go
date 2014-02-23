package local

import (
	"github.com/patdowney/downloaderd/download"
	"sync"
)

type LocalDownloadStore struct {
	LocalJSONStore
	sync.RWMutex
	repository []*download.Download
}

func (s *LocalDownloadStore) Add(download *download.Download) error {
	s.Lock()
	s.repository = append(s.repository, download)
	s.Unlock()

	err := s.Commit()

	return err
}

func (s *LocalDownloadStore) Update(download *download.Download) error {
	d, err := s.FindByID(download.ID)
	if err == nil {
		*d = *download
	}

	err = s.Commit()

	return err
}

func (s *LocalDownloadStore) Commit() error {
	return s.SaveToDisk(s.repository)
}

func (s *LocalDownloadStore) purgeUnfinished() error {
	newRepository := make([]*download.Download, 0, len(s.repository))

	for _, download := range s.repository {
		if download.Finished {
			newRepository = append(newRepository, download)
		}
	}
	s.repository = newRepository
	return nil
}

func (s *LocalDownloadStore) load() error {
	err := s.LoadFromDisk(&s.repository)

	s.purgeUnfinished()

	return err
}

func (s *LocalDownloadStore) FindByID(downloadID string) (*download.Download, error) {
	s.RLock()
	defer s.RUnlock()

	for _, download := range s.repository {
		if download.ID == downloadID {
			return download, nil
		}
	}
	return nil, nil
}

func (s *LocalDownloadStore) FindByResourceKey(resourceKey download.ResourceKey) (*download.Download, error) {
	s.RLock()
	defer s.RUnlock()

	for _, download := range s.repository {
		if download.URL == resourceKey.URL {
			if download.Metadata != nil {
				if download.Metadata.ETag == resourceKey.ETag {
					return download, nil

				}
			}
			return download, nil
		}
	}
	return nil, nil
}

func (s *LocalDownloadStore) ListAll() ([]*download.Download, error) {
	s.RLock()
	defer s.RUnlock()

	tmpRepository := make([]*download.Download, len(s.repository), len(s.repository))
	copy(tmpRepository, s.repository)

	return tmpRepository, nil
}

func NewDownloadStore(dataFile string) (download.DownloadStore, error) {
	downloadStore := &LocalDownloadStore{
		repository: make([]*download.Download, 0)}

	downloadStore.DataFile = dataFile
	err := downloadStore.load()

	return downloadStore, err
}
