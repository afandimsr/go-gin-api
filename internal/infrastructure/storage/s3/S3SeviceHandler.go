package s3

import (
	"context"
	"io"

	"github.com/afandimsr/go-gin-api/internal/domain/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Service implements the IS3Service interface for S3 storage operations.
type S3Service struct {
	client *s3.Client
	bucket string
}

// NewUploader creates a new instance of S3Service.
func NewUploader(client *s3.Client, bucket string) *S3Service {
	return &S3Service{
		client: client,
		bucket: bucket,
	}
}

func (s *S3Service) Upload(
	ctx context.Context,
	key string,
	body io.Reader,
) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
		Body:   body,
	})
	return err
}

var _ storage.IS3Service = (*S3Service)(nil)
