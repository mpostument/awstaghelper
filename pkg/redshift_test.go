package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/redshift/redshiftiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/stretchr/testify/assert"
)

type mockedRedshift struct {
	redshiftiface.RedshiftAPI
	stsiface.STSAPI
	respDescribeClustersPages redshift.DescribeClustersOutput
	respGetCallerIdentity     sts.GetCallerIdentityOutput
	respDescribeTags          redshift.DescribeTagsOutput
}

func (m *mockedRedshift) DescribeClustersPages(input *redshift.DescribeClustersInput, pageFunc func(*redshift.DescribeClustersOutput, bool) bool) error {
	pageFunc(&m.respDescribeClustersPages, true)
	return nil
}

func (m *mockedRedshift) GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return &m.respGetCallerIdentity, nil
}

func (m *mockedRedshift) DescribeTags(*redshift.DescribeTagsInput) (*redshift.DescribeTagsOutput, error) {
	return &m.respDescribeTags, nil
}

func TestGetInstances(t *testing.T) {
	cases := []*mockedRedshift{
		{
			respDescribeClustersPages: describeClustersPagesResponse,
		},
	}

	expectedResult := describeClustersPagesResponse.Clusters

	for _, c := range cases {
		t.Run("GetInstances", func(t *testing.T) {
			result := getRedshiftInstances(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeClustersPagesResponse = redshift.DescribeClustersOutput{
	Clusters: []*redshift.Cluster{
		{ClusterIdentifier: aws.String("test-cluster-1")},
	},
}

func TestParseRedshiftTags(t *testing.T) {
	cases := []*mockedRedshift{
		{
			respDescribeClustersPages: describeClustersPagesResponse,
			respGetCallerIdentity:     getRedShiftCallerIdentityResponse,
			respDescribeTags:          describeRedShiftTagsResponse,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:redshift:us-east-1:666666666:cluster:test-cluster-1", "test-cluster-1", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseRedshiftTags", func(t *testing.T) {
			result := ParseRedshiftTags("Name,Owner", c, c, "us-east-1")
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var getRedShiftCallerIdentityResponse = sts.GetCallerIdentityOutput{
	Account: aws.String("666666666"),
}

var describeRedShiftTagsResponse = redshift.DescribeTagsOutput{
	TaggedResources: []*redshift.TaggedResource{
		{
			Tag: &redshift.Tag{
				Key:   aws.String("Name"),
				Value: aws.String("test-cluster-1"),
			},
		},
		{
			Tag: &redshift.Tag{
				Key:   aws.String("Owner"),
				Value: aws.String("mpostument"),
			},
		},
	},
}
