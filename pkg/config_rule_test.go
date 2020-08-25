package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
	"github.com/stretchr/testify/assert"
)

type mockedConfigRule struct {
	configserviceiface.ConfigServiceAPI
	respDescribeConfigRules configservice.DescribeConfigRulesOutput
	respListTagsForResource configservice.ListTagsForResourceOutput
}

func (m *mockedConfigRule) DescribeConfigRules(*configservice.DescribeConfigRulesInput) (*configservice.DescribeConfigRulesOutput, error) {
	return &m.respDescribeConfigRules, nil
}

func (m *mockedConfigRule) ListTagsForResource(*configservice.ListTagsForResourceInput) (*configservice.ListTagsForResourceOutput, error) {
	return &m.respListTagsForResource, nil
}

func TestGetConfigRules(t *testing.T) {
	cases := []*mockedConfigRule{
		{
			respDescribeConfigRules: describeConfigRuleResponse,
		},
	}

	expectedResult := &describeConfigRuleResponse

	for _, c := range cases {
		t.Run("getConfigRules", func(t *testing.T) {
			result := getConfigRules(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseConfigRuleTags(t *testing.T) {
	cases := []*mockedConfigRule{
		{
			respDescribeConfigRules: describeConfigRuleResponse,
			respListTagsForResource: listConfigRuleTags,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:config:us-east-1:6666666666:config-rule/aws-service-rule/fms.amazonaws.com/test-rule", "test-rule", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseConfigRuleTags", func(t *testing.T) {
			result := ParseConfigRuleTags("Name,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var describeConfigRuleResponse = configservice.DescribeConfigRulesOutput{
	ConfigRules: []*configservice.ConfigRule{
		{
			ConfigRuleArn: aws.String("arn:aws:config:us-east-1:6666666666:config-rule/aws-service-rule/fms.amazonaws.com/test-rule"),
		},
	},
}

var listConfigRuleTags = configservice.ListTagsForResourceOutput{
	Tags: []*configservice.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("test-rule"),
		},
		{
			Key:   aws.String("Owner"),
			Value: aws.String("mpostument"),
		},
	},
}
