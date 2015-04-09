package rethinkdb

import (
	r "github.com/dancannon/gorethink"
	"github.com/patdowney/downloaderd-request/download"
)

type RequestStore struct {
	GeneralStore
}

func ResourceKeyIndex(row r.Term) interface{} {
	return []interface{}{row.Field("URL"), row.Field("Metadata").Field("ETag")}
}

func (s *RequestStore) createIndexes() error {
	err := s.IndexCreateWithFunc("ResourceKey", ResourceKeyIndex)
	if err != nil {
		return err
	}

	s.IndexWait()
	return nil
}

func (s *RequestStore) Init() error {
	return s.createIndexes()
}

func NewRequestStoreWithSession(s *r.Session, dbName string, tableName string) (*RequestStore, error) {

	generalStore, err := NewGeneralStoreWithSession(s, dbName, tableName)
	if err != nil {
		return nil, err
	}

	requestStore := &RequestStore{}
	requestStore.GeneralStore = *generalStore

	err = requestStore.Init()
	if err != nil {
		return nil, err
	}
	return requestStore, nil
}

func NewRequestStore(c Config) (*RequestStore, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address: c.Address,
		MaxIdle: c.MaxIdle,
		MaxOpen: c.MaxOpen,
	})
	if err != nil {
		return nil, err
	}

	return NewRequestStoreWithSession(session, c.Database, "RequestStore")
}

func (s *RequestStore) Add(request *download.Request) error {
	err := s.Insert(request)
	return err
}

func (s *RequestStore) FindByID(requestID string) (*download.Request, error) {
	idLookup := s.Get(requestID)

	return s.getSingleRequest(idLookup)
}

func (s *RequestStore) FindByResourceKey(resourceKey download.ResourceKey, offset uint, count uint) ([]*download.Request, error) {
	resourceKeyLookup := s.GetAllByIndex("ResourceKey", []interface{}{resourceKey.URL, resourceKey.ETag})

	return s.getMultiRequest(resourceKeyLookup, offset, count)
}

func (s *RequestStore) FindAll(offset uint, count uint) ([]*download.Request, error) {
	allLookup := s.BaseTerm()
	return s.getMultiRequest(allLookup, offset, count)
}

func (s *RequestStore) getMultiRequest(term r.Term, offset uint, count uint) ([]*download.Request, error) {
	var results []*download.Request

	rows, err := term.Slice(offset, (offset + count)).Run(s.Session)
	if err != nil {
		return nil, err
	}

	err = rows.All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *RequestStore) getSingleRequest(term r.Term) (*download.Request, error) {
	row, err := term.Run(s.Session)

	if err != nil {
		return nil, err
	}

	if row.IsNil() {
		return nil, nil
	}

	var request download.Request
	err = row.One(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}
