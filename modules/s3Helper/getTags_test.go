package s3Helper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockedListTagsForResource struct {
	s3iface.S3API
	resp     s3.ListBucketsOutput
	respTags s3.GetBucketTaggingOutput
}

func (m *mockedListTagsForResource) ListBuckets(*s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return &m.resp, nil
}

func (m *mockedListTagsForResource) GetBucketTagging(*s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	return &m.respTags, nil
}

func TestGetBuckets(t *testing.T) {
	cases := []*mockedListTagsForResource{
		{
			resp: ListBucketsResponse,
		},
	}

	expectedResult := &ListBucketsResponse

	for _, c := range cases {
		t.Run("GetBuckets", func(t *testing.T) {
			result := getBuckets(c)
			assert := assert.New(t)
			assert.EqualValues(expectedResult, result)
		})

	}
}

func TestParseS3Tags(t *testing.T) {
	cases := []*mockedListTagsForResource{
		{
			respTags: GetBucketTaggingResponse,
			resp:     ListBucketsResponse,
		},
	}

	expectedResult := [][]string{
		{"Name", "Owner", "Environment"},
		{"test-bucket-1", "mpostument", "Test"},
	}

	for _, c := range cases {
		t.Run("ParseS3Tags", func(t *testing.T) {
			result := ParseS3Tags("Owner,Environment", c)
			assert := assert.New(t)
			assert.EqualValues(expectedResult, result)
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
