package elasticSearchLib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice/elasticsearchserviceiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockedElasticSearchSts struct {
	elasticsearchserviceiface.ElasticsearchServiceAPI
	stsiface.STSAPI
	respListDomain        elasticsearchservice.ListDomainNamesOutput
	respGetCallerIdentity sts.GetCallerIdentityOutput
	respListTagsOutput    elasticsearchservice.ListTagsOutput
}

func (m *mockedElasticSearchSts) ListDomainNames(*elasticsearchservice.ListDomainNamesInput) (*elasticsearchservice.ListDomainNamesOutput, error) {
	return &m.respListDomain, nil
}

func (m *mockedElasticSearchSts) GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return &m.respGetCallerIdentity, nil
}

func (m *mockedElasticSearchSts) ListTags(*elasticsearchservice.ListTagsInput) (*elasticsearchservice.ListTagsOutput, error) {
	return &m.respListTagsOutput, nil
}

func TestGetInstances(t *testing.T) {
	cases := []*mockedElasticSearchSts{
		{
			respListDomain: listDomainsResponse,
		},
	}

	expectedResult := &listDomainsResponse

	for _, c := range cases {
		t.Run("GetInstances", func(t *testing.T) {
			result := getInstances(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseElastiCacheTags(t *testing.T) {
	cases := []*mockedElasticSearchSts{
		{
			respListDomain:        listDomainsResponse,
			respGetCallerIdentity: getCallerIdentityResponse,
			respListTagsOutput:    listTagsResponse,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:es:us-east-1:666666666:domain/test-cluster-1", "test-cluster-1", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseElastiCacheTags", func(t *testing.T) {
			result := ParseElasticSearchTags("Name,Owner", c, c, "us-east-1")
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var getCallerIdentityResponse = sts.GetCallerIdentityOutput{
	Account: aws.String("666666666"),
}

var listDomainsResponse = elasticsearchservice.ListDomainNamesOutput{
	DomainNames: []*elasticsearchservice.DomainInfo{
		{
			DomainName: aws.String("test-cluster-1"),
		},
	},
}

var listTagsResponse = elasticsearchservice.ListTagsOutput{
	TagList: []*elasticsearchservice.Tag{
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
