package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/stretchr/testify/assert"
)

type mockedASG struct {
	autoscalingiface.AutoScalingAPI
	respDescribeAutoScalingGroups autoscaling.DescribeAutoScalingGroupsOutput
}

// AWS Mocks
func (m *mockedASG) DescribeAutoScalingGroupsPages(input *autoscaling.DescribeAutoScalingGroupsInput,
	pageFunc func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool) error {
	pageFunc(&m.respDescribeAutoScalingGroups, true)
	return nil
}

var parseASGTagsResponse = autoscaling.DescribeAutoScalingGroupsOutput{
	AutoScalingGroups: []*autoscaling.Group{
		{
			AutoScalingGroupName: aws.String("asg1"),
			Tags: []*autoscaling.TagDescription{
				{
					Key:               aws.String("Name"),
					Value:             aws.String("ASG1"),
					PropagateAtLaunch: aws.Bool(true),
					ResourceId:        aws.String("asg1"),
				},
				{
					Key:               aws.String("Environment"),
					Value:             aws.String("Test"),
					PropagateAtLaunch: aws.Bool(false),
					ResourceId:        aws.String("asg1"),
				},
			},
		},
		{
			AutoScalingGroupName: aws.String("asg2"),
			Tags: []*autoscaling.TagDescription{
				{
					Key:               aws.String("Name"),
					Value:             aws.String("ASG2"),
					PropagateAtLaunch: aws.Bool(true),
					ResourceId:        aws.String("asg2"),
				},
				{
					Key:               aws.String("Environment"),
					Value:             aws.String("Dev"),
					PropagateAtLaunch: aws.Bool(false),
					ResourceId:        aws.String("asg2"),
				},
			},
		},
	},
}

func Test_getASGs(t *testing.T) {
	cases := []*mockedASG{
		{
			respDescribeAutoScalingGroups: parseASGTagsResponse,
		},
	}

	expectedResult := parseASGTagsResponse.AutoScalingGroups
	for _, c := range cases {
		t.Run("GetASGs", func(t *testing.T) {
			result := getASGs(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseASGTags(t *testing.T) {
	cases := []*mockedASG{
		{
			respDescribeAutoScalingGroups: parseASGTagsResponse,
		},
	}
	expectedResult := [][]string{
		{"AutoScalingGroupName", "Name", "Environment"},
		{"asg1", "ASG1|Propagate=true", "Test|Propagate=false"},
		{"asg2", "ASG2|Propagate=true", "Dev|Propagate=false"},
	}
	for _, c := range cases {
		t.Run("ParseASGTags", func(t *testing.T) {
			result := ParseASGTags("Name,Environment", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})
	}
}
