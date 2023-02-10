package moneta

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3API interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// This feels off.. I'm doing something wrong here, but it's late and I'm tired..
type ObjectStore struct {
	client *s3.Client
}

func New(client *s3.Client) *ObjectStore {
	return &ObjectStore{
		client: client,
	}
}

func (m *ObjectStore) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return
}
