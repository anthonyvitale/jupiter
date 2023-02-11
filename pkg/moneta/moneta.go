package moneta

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3API is an interface for the functionality needed from our blob storage host.
type S3API interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// Moneta describes the S3-like interface.
type Moneta interface {
	// Ping is used to check connection health
	func Ping(ctx context.Context) error
	// UploadImage uploads the given content to the provided path. It returns an error indicating whether or not the
	// upload was a success. TODO: might need to return more than the error?
	func UploadImage(ctx context.Context, path string, body io.Reader) error
}

// Store implements the Moneta interface and is the way to interface with an S3-like product.
type Store struct {
	client S3API
	bucket string
}

// New creates a Store.
func New(client S3API, bucket string) *Store {
	return &Store{
		client: client,
		bucket: bucket,
	}
}

// Ping checks if a bucket exists. Useful as a health check.
func (s *Store) Ping(ctx context.Context) error {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: s.bucket})
	return err
}

// UploadImage uploads an image to the backing blob storage.
func (s *Store) UploadImage(ctx context.Context, key string, body io.Reader) error {
	return nil
}
