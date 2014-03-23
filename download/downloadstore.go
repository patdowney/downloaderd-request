package download

type DownloadStore interface {
	Add(*Download) error
	Update(*Download) error
	FindByID(string) (*Download, error)
	FindByResourceKey(ResourceKey) (*Download, error)
	FindAll(uint, uint) ([]*Download, error)
	FindFinished(uint, uint) ([]*Download, error)
	FindNotFinished(uint, uint) ([]*Download, error)
	FindInProgress(uint, uint) ([]*Download, error)
	FindWaiting(uint, uint) ([]*Download, error)
}
