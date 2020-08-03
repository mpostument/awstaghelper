package cloudWatchLib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockedCloudWatchLog struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
	respDescribeLogGroups cloudwatchlogs.DescribeLogGroupsOutput
	respListTagsLogGroup  cloudwatchlogs.ListTagsLogGroupOutput
}

func (m *mockedCloudWatchLog) DescribeLogGroupsPages(input *cloudwatchlogs.DescribeLogGroupsInput, pageFunc func(*cloudwatchlogs.DescribeLogGroupsOutput, bool) bool) error {
	pageFunc(&m.respDescribeLogGroups, true)
	return nil
}

func (m *mockedCloudWatchLog) ListTagsLogGroup(*cloudwatchlogs.ListTagsLogGroupInput) (*cloudwatchlogs.ListTagsLogGroupOutput, error) {
	return &m.respListTagsLogGroup, nil
}

func TestGetInstances(t *testing.T) {
	cases := []*mockedCloudWatchLog{
		{
			respDescribeLogGroups: describeLogGroupsResponse,
		},
	}

	expectedResult := describeLogGroupsResponse.LogGroups

	for _, c := range cases {
		t.Run("getCWLogGroups", func(t *testing.T) {
			result := getCWLogGroups(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseCwLogGroupTags(t *testing.T) {
	cases := []*mockedCloudWatchLog{
		{
			respDescribeLogGroups: describeLogGroupsResponse,
			respListTagsLogGroup:  listTagsLogGroupResponse,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:logs:us-east-1:666666666666:log-group:/test-log-group:*", "test-log-group", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseCwLogGroupTags", func(t *testing.T) {
			result := ParseCwLogGroupTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeLogGroupsResponse = cloudwatchlogs.DescribeLogGroupsOutput{
	LogGroups: []*cloudwatchlogs.LogGroup{
		{
			Arn: aws.String("arn:aws:logs:us-east-1:666666666666:log-group:/test-log-group:*"),
		},
	},
}

var listTagsLogGroupResponse = cloudwatchlogs.ListTagsLogGroupOutput{
	Tags: map[string]*string{
		"Name":  aws.String("test-log-group"),
		"Owner": aws.String("mpostument"),
	},
}
