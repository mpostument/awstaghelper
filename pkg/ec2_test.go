package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
)

type mockedEc2 struct {
	ec2iface.EC2API
	respDescribeInstances ec2.DescribeInstancesOutput
}

// AWS Mocks
func (m *mockedEc2) DescribeInstancesPages(input *ec2.DescribeInstancesInput, pageFunc func(*ec2.DescribeInstancesOutput, bool) bool) error {
	pageFunc(&m.respDescribeInstances, true)
	return nil
}

// Tests
func TestGetEC2Instances(t *testing.T) {
	cases := []*mockedEc2{
		{
			respDescribeInstances: parseEC2TagsResponse,
		},
	}

	expectedResult := parseEC2TagsResponse.Reservations
	for _, c := range cases {
		t.Run("GetInstances", func(t *testing.T) {
			result := getEC2Instances(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseEc2Tags(t *testing.T) {
	cases := []*mockedEc2{
		{
			respDescribeInstances: parseEC2TagsResponse,
		},
	}
	expectedResult := [][]string{
		{"Id", "Name", "Environment", "Owner"},
		{"i-666666", "TestInstance1", "Test", ""},
		{"i-777777", "TestInstance2", "Test", "mpostument"},
	}
	for _, c := range cases {
		t.Run("ParseEC2Tags", func(t *testing.T) {
			result := ParseEC2Tags("Name,Environment,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})
	}
}

var parseEC2TagsResponse = ec2.DescribeInstancesOutput{
	Reservations: []*ec2.Reservation{
		{
			Instances: []*ec2.Instance{
				{
					InstanceId: aws.String("i-666666"),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("TestInstance1"),
						},
						{
							Key:   aws.String("Environment"),
							Value: aws.String("Test"),
						},
					},
				},
				{
					InstanceId: aws.String("i-777777"),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("TestInstance2"),
						},
						{
							Key:   aws.String("Environment"),
							Value: aws.String("Test"),
						},
						{
							Key:   aws.String("Owner"),
							Value: aws.String("mpostument"),
						},
					},
				},
			},
		},
	},
}
