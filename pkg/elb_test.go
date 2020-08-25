package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/stretchr/testify/assert"
)

type mockedElbv2 struct {
	elbv2iface.ELBV2API
	respDescribeLoadBalancers elbv2.DescribeLoadBalancersOutput
	respDescribeTags          elbv2.DescribeTagsOutput
}

func (m *mockedElbv2) DescribeLoadBalancersPages(input *elbv2.DescribeLoadBalancersInput, pageFunc func(*elbv2.DescribeLoadBalancersOutput, bool) bool) error {
	pageFunc(&m.respDescribeLoadBalancers, true)
	return nil
}

func (m *mockedElbv2) DescribeTags(*elbv2.DescribeTagsInput) (*elbv2.DescribeTagsOutput, error) {
	return &m.respDescribeTags, nil
}

func TestGetElbV2(t *testing.T) {
	cases := []*mockedElbv2{
		{
			respDescribeLoadBalancers: describeLoadBalancersResponse,
		},
	}

	expectedResult := describeLoadBalancersResponse.LoadBalancers

	for _, c := range cases {
		t.Run("getElbV2", func(t *testing.T) {
			result := getElbV2(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseElbV2Tags(t *testing.T) {
	cases := []*mockedElbv2{
		{
			respDescribeLoadBalancers: describeLoadBalancersResponse,
			respDescribeTags:          listElbV2TagsResponse,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:elasticloadbalancing:us-east-1:084888172679:loadbalancer/app/test-lb-1/4d1a0eb7f21dc6f6", "test-lb-1", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseElbV2Tags", func(t *testing.T) {
			result := ParseElbV2Tags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeLoadBalancersResponse = elbv2.DescribeLoadBalancersOutput{
	LoadBalancers: []*elbv2.LoadBalancer{
		{
			LoadBalancerArn: aws.String("arn:aws:elasticloadbalancing:us-east-1:084888172679:loadbalancer/app/test-lb-1/4d1a0eb7f21dc6f6"),
		},
	},
}

var listElbV2TagsResponse = elbv2.DescribeTagsOutput{
	TagDescriptions: []*elbv2.TagDescription{
		{
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("test-lb-1"),
				},
				{
					Key:   aws.String("Owner"),
					Value: aws.String("mpostument"),
				},
			},
		},
	},
}
