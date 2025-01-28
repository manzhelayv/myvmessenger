package drivers

import "context"

// DataStore is an abstraction for data manipulation.
type ImagesStorage interface {
	Close(ctx context.Context) error
	Connect() error

	// repositories
	ImagesRepository() ImagesRepository
}

type FilesStorage interface {
	Connect() error
	Close(ctx context.Context) error

	FilesRepository() FilesRepository
}
