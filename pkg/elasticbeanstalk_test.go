package pkg

import (
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk/elasticbeanstalkiface"
	"github.com/stretchr/testify/assert"
)

type mockedEBEnv struct {
	elasticbeanstalkiface.ElasticBeanstalkAPI
	respDescribeEnvironments elasticbeanstalk.EnvironmentDescriptionsMessage
	respListTagsForResource  elasticbeanstalk.ListTagsForResourceOutput
}

func (m *mockedEBEnv) DescribeEnvironments(*elasticbeanstalk.DescribeEnvironmentsInput) (*elasticbeanstalk.EnvironmentDescriptionsMessage, error) {
	return &m.respDescribeEnvironments, nil
}

func (m *mockedEBEnv) ListTagsForResource(*elasticbeanstalk.ListTagsForResourceInput) (*elasticbeanstalk.ListTagsForResourceOutput, error) {
	return &m.respListTagsForResource, nil
}

func TestEBEnvironments(t *testing.T) {
	cases := []*mockedEBEnv{
		{
			respDescribeEnvironments: describeEnvironmentsResponse,
		},
	}

	expectedResult := &describeEnvironmentsResponse

	for _, c := range cases {
		t.Run("getEBEnvironments", func(t *testing.T) {
			result := getEBEnvironments(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseEBEnvironmentsTags(t *testing.T) {
	cases := []*mockedEBEnv{
		{
			respDescribeEnvironments: describeEnvironmentsResponse,
			respListTagsForResource:  listTagsForEBEnvironments,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:elasticbeanstalk:us-east-1:12345678:environment/test-app/test-env", "test-eb1", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseEBTags", func(t *testing.T) {
			result := ParseEBTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeEnvironmentsResponse = elasticbeanstalk.EnvironmentDescriptionsMessage{
	Environments: []*elasticbeanstalk.EnvironmentDescription{
		{
			EnvironmentArn: aws.String("arn:aws:elasticbeanstalk:us-east-1:12345678:environment/test-app/test-env"),
		},
	},
}

var listTagsForEBEnvironments = elasticbeanstalk.ListTagsForResourceOutput{
	ResourceTags: []*elasticbeanstalk.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("test-eb1"),
		},
		{
			Key:   aws.String("Owner"),
			Value: aws.String("mpostument"),
		},
	},
}
