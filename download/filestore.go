package download

import (
	"io"
)

type FileStore interface {
	GetWriter(*Download) (io.WriteCloser, error)
	GetReader(*Download) (io.ReadCloser, error)
}