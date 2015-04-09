package download

import (
	"github.com/patdowney/downloaderd-request/api"
)

// ToAPIMetadata ...
func ToAPIMetadata(dm *Metadata) *api.Metadata {
	m := &api.Metadata{
		TimeRequested: dm.TimeRequested,
		MimeType:      dm.MimeType,
		Size:          dm.Size,
		Server:        dm.Server,
		LastModified:  dm.LastModified,
		ETag:          dm.ETag,
		Expires:       dm.Expires,
		StatusCode:    dm.StatusCode}

	return m
}
