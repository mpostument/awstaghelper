package ec2Helper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
	"testing"
)

type stubEc2 struct {
	ec2iface.EC2API
	resp   ec2.DescribeInstancesOutput
	result [][]string
}

// AWS Mocks
func (m *stubEc2) DescribeInstancesPages(input *ec2.DescribeInstancesInput, pageFunc func(*ec2.DescribeInstancesOutput, bool) bool) error {
	pageFunc(&m.resp, true)
	return nil
}

// Tests
func TestParseEc2Tags(t *testing.T) {
	cases := []*stubEc2{
		{
			resp: ec2Resp,
			result: [][]string{
				{"Id", "Name", "Environment", "Owner"},
				{"i-666666", "TestInstance1", "Test", ""},
				{"i-777777", "TestInstance2", "Test", "mpostument"},
				{"i-888888", "TestInstance3", "Test", ""},
			},
		},
	}
	for _, c := range cases {
		t.Run("ParseEc2Tags", func(t *testing.T) {
			result := ParseEc2Tags("Name,Environment,Owner", c)
			assert := assert.New(t)
			assert.EqualValues(c.result, result)
		})

	}
}

var ec2Resp = ec2.DescribeInstancesOutput{
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
				{
					InstanceId: aws.String("i-888888"),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("TestInstance3"),
						},
						{
							Key:   aws.String("Environment"),
							Value: aws.String("Test"),
						},
						{
							Key:   aws.String("Project"),
							Value: aws.String("awstaghelper"),
						},
					},
				},
			},
		},
	},
}
