package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockedEcrRepo struct {
	ecriface.ECRAPI
	respDescribeRepositories ecr.DescribeRepositoriesOutput
	respListTagsForResource  ecr.ListTagsForResourceOutput
}

func (m *mockedEcrRepo) DescribeRepositoriesPages(input *ecr.DescribeRepositoriesInput, pageFunc func(*ecr.DescribeRepositoriesOutput, bool) bool) error {
	pageFunc(&m.respDescribeRepositories, true)
	return nil
}

func (m *mockedEcrRepo) ListTagsForResource(*ecr.ListTagsForResourceInput) (*ecr.ListTagsForResourceOutput, error) {
	return &m.respListTagsForResource, nil
}

func TestGetEcrRepositories(t *testing.T) {
	cases := []*mockedEcrRepo{
		{
			respDescribeRepositories: describeRepositoriesResponse,
		},
	}

	expectedResult := describeRepositoriesResponse.Repositories

	for _, c := range cases {
		t.Run("getEcrRepositories", func(t *testing.T) {
			result := getEcrRepositories(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseEcrRepositoriesTags(t *testing.T) {
	cases := []*mockedEcrRepo{
		{
			respDescribeRepositories: describeRepositoriesResponse,
			respListTagsForResource:  listTagsForEcrResourceResponse,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:ecr:region:012345678910:repository/test-repo1", "test-repo1", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseEcrRepositoriesTags", func(t *testing.T) {
			result := ParseEcrRepositoriesTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeRepositoriesResponse = ecr.DescribeRepositoriesOutput{
	Repositories: []*ecr.Repository{
		{
			RepositoryArn: aws.String("arn:aws:ecr:region:012345678910:repository/test-repo1"),
		},
	},
}

var listTagsForEcrResourceResponse = ecr.ListTagsForResourceOutput{
	Tags: []*ecr.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("test-repo1"),
		},
		{
			Key:   aws.String("Owner"),
			Value: aws.String("mpostument"),
		},
	},
}
