package ec2Lib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
	"testing"
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
func TestGetInstances(t *testing.T) {
	cases := []*mockedEc2{
		{
			respDescribeInstances: parseTagsResponse,
		},
	}

	expectedResult := parseTagsResponse.Reservations
	for _, c := range cases {
		t.Run("GetInstances", func(t *testing.T) {
			result := getEc2(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseEc2Tags(t *testing.T) {
	cases := []*mockedEc2{
		{
			respDescribeInstances: parseTagsResponse,
		},
	}
	expectedResult := [][]string{
		{"Id", "Name", "Environment", "Owner"},
		{"i-666666", "TestInstance1", "Test", ""},
		{"i-777777", "TestInstance2", "Test", "mpostument"},
	}
	for _, c := range cases {
		t.Run("ParseEc2Tags", func(t *testing.T) {
			result := ParseEc2Tags("Name,Environment,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})
	}
}

var parseTagsResponse = ec2.DescribeInstancesOutput{
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
