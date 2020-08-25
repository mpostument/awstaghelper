package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
)

type mockedS3 struct {
	s3iface.S3API
	respListBuckets      s3.ListBucketsOutput
	respGetBucketTagging s3.GetBucketTaggingOutput
}

func (m *mockedS3) ListBuckets(*s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return &m.respListBuckets, nil
}

func (m *mockedS3) GetBucketTagging(*s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	return &m.respGetBucketTagging, nil
}

func TestGetBuckets(t *testing.T) {
	cases := []*mockedS3{
		{
			respListBuckets: ListBucketsResponse,
		},
	}

	expectedResult := &ListBucketsResponse

	for _, c := range cases {
		t.Run("GetBuckets", func(t *testing.T) {
			result := getBuckets(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseS3Tags(t *testing.T) {
	cases := []*mockedS3{
		{
			respGetBucketTagging: GetBucketTaggingResponse,
			respListBuckets:      ListBucketsResponse,
		},
	}

	expectedResult := [][]string{
		{"Name", "Owner", "Environment"},
		{"test-bucket-1", "mpostument", "Test"},
	}

	for _, c := range cases {
		t.Run("ParseS3Tags", func(t *testing.T) {
			result := ParseS3Tags("Owner,Environment", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var ListBucketsResponse = s3.ListBucketsOutput{
	Buckets: []*s3.Bucket{
		{
			Name: aws.String("test-bucket-1"),
		},
	},
}

var GetBucketTaggingResponse = s3.GetBucketTaggingOutput{
	TagSet: []*s3.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("test-bucket-1"),
		},
		{
			Key:   aws.String("Owner"),
			Value: aws.String("mpostument"),
		},
		{
			Key:   aws.String("Environment"),
			Value: aws.String("Test"),
		},
	},
}
