package download

import (
	"time"

	"github.com/patdowney/downloaderd-request/common"
)

// Error ...
type Error struct {
	common.ErrorWrapper
	DownloadID string
}

// NewError ...
func NewError(id string, err error, errorTime time.Time) *Error {
	downloadErr := &Error{DownloadID: id}
	downloadErr.Time = errorTime
	downloadErr.OriginalError = err.Error()

	return downloadErr
}

// RequestError ...
type RequestError struct {
	common.ErrorWrapper
}

// NewRequestError ...
func NewRequestError(err error, errorTime time.Time) *RequestError {
	reqErr := &RequestError{}
	reqErr.Time = errorTime
	reqErr.OriginalError = err.Error()
	return reqErr
}
