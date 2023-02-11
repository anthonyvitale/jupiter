package moneta

import (
	"context"
	"testing"

	"github.com/anthonyvitale/jupiter/pkg/moneta/mocks"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fortytw2/leaktest"
	"github.com/golang/mock/gomock"
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

	suite.store = New(suite.s3Mock, suite.bucket)
}

func (suite *MonetaSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func TestMonetaSuite(t *testing.T) {
	defer leaktest.Check(t)
	suite.Run(t, new(MonetaSuite))
}

func (suite *MonetaSuite) Test_MonetaPing() {
	// TODO: TDD here with ping error and no error
	suite.s3Mock.
		EXPECT().
		HeadBucket(suite.ctx, &s3.HeadBucketInput{Bucket: aws.String(suite.bucket)}).
		Return(nil, nil)

	err := suite.store.Ping(suite.ctx)
	suite.NoError(err)
}
