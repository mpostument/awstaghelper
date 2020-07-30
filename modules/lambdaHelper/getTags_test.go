package lambdaHelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockedLambda struct {
	lambdaiface.LambdaAPI
	respListFunctions lambda.ListFunctionsOutput
	respListTags      lambda.ListTagsOutput
}

func (m *mockedLambda) ListFunctionsPages(input *lambda.ListFunctionsInput, pageFunc func(*lambda.ListFunctionsOutput, bool) bool) error {
	pageFunc(&m.respListFunctions, true)
	return nil
}

func (m *mockedLambda) ListTags(*lambda.ListTagsInput) (*lambda.ListTagsOutput, error) {
	return &m.respListTags, nil
}

func TestGetInstances(t *testing.T) {
	cases := []*mockedLambda{
		{
			respListFunctions: describeListFunctionsResponse,
		},
	}

	expectedResult := describeListFunctionsResponse.Functions

	for _, c := range cases {
		t.Run("GetInstances", func(t *testing.T) {
			result := getInstances(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseLambdaTags(t *testing.T) {
	cases := []*mockedLambda{
		{
			respListFunctions: describeListFunctionsResponse,
			respListTags:      listTagsResponse,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:lambda:us-east-1:666666666666:function:test-function-1", "test-function-1", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseElastiCacheTags", func(t *testing.T) {
			result := ParseLambdasTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeListFunctionsResponse = lambda.ListFunctionsOutput{
	Functions: []*lambda.FunctionConfiguration{
		{
			FunctionArn: aws.String("arn:aws:lambda:us-east-1:666666666666:function:test-function-1"),
		},
	},
}

var listTagsResponse = lambda.ListTagsOutput{
	Tags: map[string]*string{
		"Name":  aws.String("test-function-1"),
		"Owner": aws.String("mpostument"),
	},
}
