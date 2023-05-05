package moneta

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/anthonyvitale/jupiter/pkg/moneta/mocks"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fortytw2/leaktest"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MonetaSuite struct {
	suite.Suite
	*require.Assertions
	ctx  context.Context
	ctrl *gomock.Controller

	s3Mock *mocks.MockS3API
	store  *Store
	bucket string
}

func (suite *MonetaSuite) SetupTest() {
	suite.Assertions = suite.Suite.Require()
	suite.ctx = context.Background()
	suite.bucket = "test_bucket"

	suite.ctrl = gomock.NewController(suite.T())
	suite.s3Mock = mocks.NewMockS3API(suite.ctrl)

	store, err := New(suite.s3Mock, suite.bucket)
	suite.NoError(err)
	suite.store = store
}

func (suite *MonetaSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func TestMonetaSuite(t *testing.T) {
	defer leaktest.Check(t)
	suite.Run(t, new(MonetaSuite))
}

func (suite *MonetaSuite) Test_MonetaNew() {
	type args struct {
		bucket string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty string",
			args: args{
				bucket: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		store, err := New(suite.s3Mock, tt.args.bucket)
		if tt.wantErr {
			suite.Error(err, tt.name)
		} else {
			suite.NoError(err, tt.name)
		}
		suite.Nil(store)
	}
}

func (suite *MonetaSuite) Test_MonetaPing() {
	type fields struct {
		s3Mock *mocks.MockS3API
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		wantErr bool
	}{
		{
			name: "ping successful, no error",
			prepare: func(f *fields) {
				f.s3Mock.EXPECT().HeadBucket(suite.ctx, &s3.HeadBucketInput{Bucket: aws.String(suite.bucket)}).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "ping failure, error",
			prepare: func(f *fields) {
				f.s3Mock.EXPECT().HeadBucket(suite.ctx, &s3.HeadBucketInput{Bucket: aws.String(suite.bucket)}).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		f := fields{
			s3Mock: suite.s3Mock,
		}
		if tt.prepare != nil {
			tt.prepare(&f)
		}

		err := suite.store.Ping(suite.ctx)
		if tt.wantErr {
			suite.Error(err, tt.name)
		} else {
			suite.NoError(err, tt.name)
		}
	}
}

func (suite *MonetaSuite) Test_MonetaUploadImage() {
	type fields struct {
		s3Mock *mocks.MockS3API
	}
	type args struct {
		key  string
		body io.Reader
	}
	tests := []struct {
		name    string
		args    args
		prepare func(f *fields)
		wantErr bool
	}{
		{
			name: "uploadImage successful, no error",
			args: args{
				key:  "key_1",
				body: bytes.NewBufferString("my_body"),
			},
			prepare: func(f *fields) {
				f.s3Mock.EXPECT().PutObject(suite.ctx, &s3.PutObjectInput{
					Bucket: aws.String(suite.bucket),
					Key:    aws.String("key_1"),
					Body:   bytes.NewBufferString("my_body"),
				}).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "uploadImage unsuccessful, error",
			args: args{
				key:  "key_1",
				body: bytes.NewBufferString("my_body"),
			},
			prepare: func(f *fields) {
				f.s3Mock.EXPECT().PutObject(suite.ctx, &s3.PutObjectInput{
					Bucket: aws.String(suite.bucket),
					Key:    aws.String("key_1"),
					Body:   bytes.NewBufferString("my_body"),
				}).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "uploadImage missing key",
			args: args{
				key:  "",
				body: bytes.NewBufferString("my_body"),
			},
			prepare: func(f *fields) {},
			wantErr: true,
		},
		{
			name: "uploadImage missing body",
			args: args{
				key:  "key_1",
				body: nil,
			},
			prepare: func(f *fields) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		f := fields{
			s3Mock: suite.s3Mock,
		}
		if tt.prepare != nil {
			tt.prepare(&f)
		}

		err := suite.store.UploadImage(suite.ctx, tt.args.key, tt.args.body)
		if tt.wantErr {
			suite.Error(err, tt.name)
		} else {
			suite.NoError(err, tt.name)
		}
	}
}
