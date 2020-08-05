package cloudWatchLib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/stretchr/testify/assert"
	"testing"
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
			respListTagsForResource:  listAlarmTags,
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

var listAlarmTags = cloudwatch.ListTagsForResourceOutput{
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
