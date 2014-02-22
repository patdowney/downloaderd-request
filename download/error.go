package download

import (
	"github.com/patdowney/downloaderd/common"
	"time"
)

type DownloadError struct {
	common.ErrorWrapper
	DownloadId string
}

func NewDownloadError(id string, err error, errorTime time.Time) *DownloadError {
	downloadErr := &DownloadError{DownloadId: id}
	downloadErr.Time = errorTime
	downloadErr.OriginalError = err.Error()

	return downloadErr
}

type RequestError struct {
	common.ErrorWrapper
}

func NewRequestError(err error, errorTime time.Time) *RequestError {
	reqErr := &RequestError{}
	reqErr.Time = errorTime
	reqErr.OriginalError = err.Error()
	return reqErr
}

type CallbackError struct {
	DownloadError
	RequestId string
}