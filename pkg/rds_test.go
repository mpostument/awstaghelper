package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/stretchr/testify/assert"
)

type mockedRds struct {
	rdsiface.RDSAPI
	respDescribeDBInstances rds.DescribeDBInstancesOutput
	respListTagsForResource rds.ListTagsForResourceOutput
}

func (m *mockedRds) DescribeDBInstancesPages(input *rds.DescribeDBInstancesInput, pageFunc func(*rds.DescribeDBInstancesOutput, bool) bool) error {
	pageFunc(&m.respDescribeDBInstances, true)
	return nil
}

func (m *mockedRds) ListTagsForResource(*rds.ListTagsForResourceInput) (*rds.ListTagsForResourceOutput, error) {
	return &m.respListTagsForResource, nil
}

func TestGetRDSInstances(t *testing.T) {
	cases := []*mockedRds{
		{
			respDescribeDBInstances: describeDbInstancesResponse,
		},
	}

	expectedResult := describeDbInstancesResponse.DBInstances

	for _, c := range cases {
		t.Run("GetInstances", func(t *testing.T) {
			result := getRDSInstances(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseRDSTags(t *testing.T) {
	cases := []*mockedRds{
		{
			respDescribeDBInstances: describeDbInstancesResponse,
			respListTagsForResource: listTagsForResourceResponse,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:rds:us-east-1:66666666:db:test-db-1", "test-cluster-1", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseRDSTags", func(t *testing.T) {
			result := ParseRDSTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeDbInstancesResponse = rds.DescribeDBInstancesOutput{
	DBInstances: []*rds.DBInstance{
		{
			DBInstanceArn: aws.String("arn:aws:rds:us-east-1:66666666:db:test-db-1"),
		},
	},
}

var listTagsForResourceResponse = rds.ListTagsForResourceOutput{
	TagList: []*rds.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("test-cluster-1"),
		},
		{
			Key:   aws.String("Owner"),
			Value: aws.String("mpostument"),
		},
	},
}
