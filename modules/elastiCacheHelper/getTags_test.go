package elastiCacheHelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockedElasticacheSts struct {
	elasticacheiface.ElastiCacheAPI
	stsiface.STSAPI
	respTagList               elasticache.TagListMessage
	respGetCallerIdentity     sts.GetCallerIdentityOutput
	respDescribeCacheClusters elasticache.DescribeCacheClustersOutput
}

func (m *mockedElasticacheSts) DescribeCacheClustersPages(input *elasticache.DescribeCacheClustersInput, pageFunc func(*elasticache.DescribeCacheClustersOutput, bool) bool) error {
	pageFunc(&m.respDescribeCacheClusters, true)
	return nil
}

func (m *mockedElasticacheSts) GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return &m.respGetCallerIdentity, nil
}

func (m *mockedElasticacheSts) ListTagsForResource(*elasticache.ListTagsForResourceInput) (*elasticache.TagListMessage, error) {
	return &m.respTagList, nil
}

func TestGetInstances(t *testing.T) {
	cases := []*mockedElasticacheSts{
		{
			respDescribeCacheClusters: describeCacheClustersResponse,
		},
	}

	expectedResult := describeCacheClustersResponse.CacheClusters

	for _, c := range cases {
		t.Run("GetInstances", func(t *testing.T) {
			result := getInstances(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseElastiCacheTags(t *testing.T) {
	cases := []*mockedElasticacheSts{
		{
			respTagList:               listTagsForResourceResponse,
			respGetCallerIdentity:     getCallerIdentityResponse,
			respDescribeCacheClusters: describeCacheClustersResponse,
		},
	}

	expectedResult := [][]string{
		{"Arn", "Name", "Owner"},
		{"arn:aws:elasticache:us-east-1:666666666:cluster:test-cluster-1", "test-cluster-1", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseElastiCacheTags", func(t *testing.T) {
			result := ParseElastiCacheTags("Name,Owner", c, c, "us-east-1")
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var getCallerIdentityResponse = sts.GetCallerIdentityOutput{
	Account: aws.String("666666666"),
}

var describeCacheClustersResponse = elasticache.DescribeCacheClustersOutput{
	CacheClusters: []*elasticache.CacheCluster{
		{
			CacheClusterId: aws.String("test-cluster-1"),
		},
	},
}

var listTagsForResourceResponse = elasticache.TagListMessage{
	TagList: []*elasticache.Tag{
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
