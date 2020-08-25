package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/stretchr/testify/assert"
)

type mockedCloudFront struct {
	cloudfrontiface.CloudFrontAPI
	respListDistributions   cloudfront.ListDistributionsOutput
	respListTagsForResource cloudfront.ListTagsForResourceOutput
}

func (m *mockedCloudFront) ListDistributions(*cloudfront.ListDistributionsInput) (*cloudfront.ListDistributionsOutput, error) {
	return &m.respListDistributions, nil
}

func (m *mockedCloudFront) ListTagsForResource(*cloudfront.ListTagsForResourceInput) (*cloudfront.ListTagsForResourceOutput, error) {
	return &m.respListTagsForResource, nil
}

func TestGetCloudFrontDistributions(t *testing.T) {
	cases := []*mockedCloudFront{
		{
			respListDistributions: describeDistributionsResponse,
		},
	}

	expectedResult := &describeDistributionsResponse

	for _, c := range cases {
		t.Run("getDistributions", func(t *testing.T) {
			result := getDistributions(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseCloudFrontDistributionTags(t *testing.T) {
	cases := []*mockedCloudFront{
		{
			respListDistributions:   describeDistributionsResponse,
			respListTagsForResource: listDistributionsTags,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:cloudfront::084888172679:distribution/E2ZG1054Z6EMOE", "test-distribution", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseDistributionsTags", func(t *testing.T) {
			result := ParseDistributionsTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeDistributionsResponse = cloudfront.ListDistributionsOutput{
	DistributionList: &cloudfront.DistributionList{
		Items: []*cloudfront.DistributionSummary{
			{
				ARN: aws.String("arn:aws:cloudfront::084888172679:distribution/E2ZG1054Z6EMOE"),
			},
		},
	},
}

var listDistributionsTags = cloudfront.ListTagsForResourceOutput{
	Tags: &cloudfront.Tags{
		Items: []*cloudfront.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("test-distribution"),
			},
			{
				Key:   aws.String("Owner"),
				Value: aws.String("mpostument"),
			},
		},
	},
}
