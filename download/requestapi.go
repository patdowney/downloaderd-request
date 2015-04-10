package download

import (
	"github.com/patdowney/downloaderd-common/common"
	"github.com/patdowney/downloaderd-request/api"
)

// ToAPIRequestList ...
func ToAPIRequestList(origList *[]*Request) *[]*api.Request {
	rs := make([]*api.Request, len(*origList))

	for i, r := range *origList {
		rs[i] = ToAPIRequest(r)
	}

	return &rs
}

// ToAPIRequest ...
func ToAPIRequest(orig *Request) *api.Request {
	r := &api.Request{
		ID:                   orig.ID,
		URL:                  orig.URL,
		ExpectedChecksum:     orig.Checksum,
		ExpectedChecksumType: orig.ChecksumType,
		DownloadID:           orig.DownloadID,
		TimeRequested:        orig.TimeRequested,
		Callback:             orig.Callback,
		Errors:               make([]*api.Error, 0, len(orig.Errors)),
		Links:                make([]api.Link, 0)}

	if orig.Metadata != nil {
		r.Metadata = ToAPIMetadata(orig.Metadata)
	}

	if len(orig.Errors) > 0 {
		for _, e := range orig.Errors {
			if e.OriginalError != "" {
				r.Errors = append(r.Errors, ToAPIError(&e.TimestampedError))
			}
		}
	}

	return r
}

// FromAPIIncomingRequest ...
func FromAPIIncomingRequest(air *api.IncomingRequest) *Request {
	downloadReq := &Request{
		URL:          air.URL,
		Checksum:     air.Checksum,
		ChecksumType: air.ChecksumType,
		Callback:     air.Callback,
		Errors:       make([]*RequestError, 0)}

	return downloadReq
}

// ToAPIError ...
func ToAPIError(e *common.TimestampedError) *api.Error {
	err := &api.Error{Time: e.Time}
	if e.OriginalError != "" {
		err.Error = e.OriginalError
	} else {
		err.Error = "error missing - weird."
	}

	return err
}
