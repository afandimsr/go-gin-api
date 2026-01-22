package storage

import (
	"context"
	"io"
)

// IS3Service defines the interface for S3 storage operations.
type IS3Service interface {
	Upload(
		ctx context.Context,
		key string,
		body io.Reader,
	) error
}
