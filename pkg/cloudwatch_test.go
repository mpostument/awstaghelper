package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/stretchr/testify/assert"
)

type mockedCloudWatchAlarm struct {
	cloudwatchiface.CloudWatchAPI
	respDescribeAlarmsGroups cloudwatch.DescribeAlarmsOutput
	respListTagsForResource  cloudwatch.ListTagsForResourceOutput
}

func (m *mockedCloudWatchAlarm) DescribeAlarmsPages(input *cloudwatch.DescribeAlarmsInput, pageFunc func(*cloudwatch.DescribeAlarmsOutput, bool) bool) error {
	pageFunc(&m.respDescribeAlarmsGroups, true)
	return nil
}

func (m *mockedCloudWatchAlarm) ListTagsForResource(*cloudwatch.ListTagsForResourceInput) (*cloudwatch.ListTagsForResourceOutput, error) {
	return &m.respListTagsForResource, nil
}

func TestGetCWAlarm(t *testing.T) {
	cases := []*mockedCloudWatchAlarm{
		{
			respDescribeAlarmsGroups: describeAlarmResponse,
		},
	}

	expectedResult := describeAlarmResponse.MetricAlarms

	for _, c := range cases {
		t.Run("getCWAlarm", func(t *testing.T) {
			result := getCWAlarm(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseCwAlarmTags(t *testing.T) {
	cases := []*mockedCloudWatchAlarm{
		{
			respDescribeAlarmsGroups: describeAlarmResponse,
			respListTagsForResource:  listCloudWatchAlarmsResp,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:cloudwatch:us-east-1:6666666666:alarm:test-alarm", "test-alarm", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseCwAlarmTags", func(t *testing.T) {
			result := ParseCwAlarmTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeAlarmResponse = cloudwatch.DescribeAlarmsOutput{
	MetricAlarms: []*cloudwatch.MetricAlarm{
		{
			AlarmArn: aws.String("arn:aws:cloudwatch:us-east-1:6666666666:alarm:test-alarm"),
		},
	},
}

var listCloudWatchAlarmsResp = cloudwatch.ListTagsForResourceOutput{
	Tags: []*cloudwatch.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("test-alarm"),
		},
		{
			Key:   aws.String("Owner"),
			Value: aws.String("mpostument"),
		},
	},
}

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

func TestGetCWLogGroups(t *testing.T) {
	cases := []*mockedCloudWatchLog{
		{
			respDescribeLogGroups: describeCloudWatchLogGroupsResponse,
		},
	}

	expectedResult := describeCloudWatchLogGroupsResponse.LogGroups

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
			respDescribeLogGroups: describeCloudWatchLogGroupsResponse,
			respListTagsLogGroup:  listCloudWatchLogsTagResponse,
		},
	}

	expectedResult := [][]string{
		{"LogGroupName", "Name", "Owner"},
		{"test-log-group", "test-log-group", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseCwLogGroupTags", func(t *testing.T) {
			result := ParseCwLogGroupTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeCloudWatchLogGroupsResponse = cloudwatchlogs.DescribeLogGroupsOutput{
	LogGroups: []*cloudwatchlogs.LogGroup{
		{
			LogGroupName: aws.String("test-log-group"),
		},
	},
}

var listCloudWatchLogsTagResponse = cloudwatchlogs.ListTagsLogGroupOutput{
	Tags: map[string]*string{
		"Name":  aws.String("test-log-group"),
		"Owner": aws.String("mpostument"),
	},
}
