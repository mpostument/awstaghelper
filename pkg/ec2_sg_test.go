package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
)

type mockedSecurityGroup struct {
	ec2iface.EC2API
	respDescribeSecurityGroups ec2.DescribeSecurityGroupsOutput
}

// AWS Mocks
func (m *mockedSecurityGroup) DescribeSecurityGroupsPages(
	input *ec2.DescribeSecurityGroupsInput,
	pageFunc func(*ec2.DescribeSecurityGroupsOutput, bool) bool) error {
	pageFunc(&m.respDescribeSecurityGroups, true)
	return nil
}

// Tests
func TestGetSecurityGroups(t *testing.T) {
	cases := []*mockedSecurityGroup{
		{
			respDescribeSecurityGroups: parseSecurityGroupTagsResponse,
		},
	}

	expectedResult := parseSecurityGroupTagsResponse.SecurityGroups

	for _, c := range cases {
		t.Run("GetSecurityGroups", func(t *testing.T) {
			result := getSecurityGroups(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseSecurityGroupTags(t *testing.T) {
	cases := []*mockedSecurityGroup{
		{
			respDescribeSecurityGroups: parseSecurityGroupTagsResponse,
		},
	}
	expectedResult := [][]string{
		{"Id", "Name", "Environment", "Owner"},
		{"sg-test", "TestSecurityGroup1", "Test", ""},
	}
	for _, c := range cases {
		t.Run("ParseSecurityGroupTags", func(t *testing.T) {
			result := ParseSecurityGroupTags("Name,Environment,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})
	}
}

var parseSecurityGroupTagsResponse = ec2.DescribeSecurityGroupsOutput{
	SecurityGroups: []*ec2.SecurityGroup{
		{
			Description: aws.String("testSg"),
			GroupId:     aws.String("sg-test"),
			GroupName:   aws.String("testSg"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("TestSecurityGroup1"),
				},
				{
					Key:   aws.String("Environment"),
					Value: aws.String("Test"),
				},
			},
		},
	},
}
