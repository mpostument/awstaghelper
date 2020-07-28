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

type mockedListTagsForResource struct {
	elasticacheiface.ElastiCacheAPI
	stsiface.STSAPI
	resp                    elasticache.TagListMessage
	respSts                 sts.GetCallerIdentityOutput
	respCacheClustersOutput elasticache.DescribeCacheClustersOutput
}

func (m *mockedListTagsForResource) DescribeCacheClustersPages(input *elasticache.DescribeCacheClustersInput, pageFunc func(*elasticache.DescribeCacheClustersOutput, bool) bool) error {
	pageFunc(&m.respCacheClustersOutput, true)
	return nil
}

func (m *mockedListTagsForResource) GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return &m.respSts, nil
}

func (m *mockedListTagsForResource) ListTagsForResource(*elasticache.ListTagsForResourceInput) (*elasticache.TagListMessage, error) {
	return &m.resp, nil
}

func TestGetInstances(t *testing.T) {
	cases := []*mockedListTagsForResource{
		{
			respCacheClustersOutput: describeCacheClustersResponse,
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
	cases := []*mockedListTagsForResource{
		{
			resp:                    listTagsForResourceResponse,
			respSts:                 getCallerIdentityResponse,
			respCacheClustersOutput: describeCacheClustersResponse,
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
