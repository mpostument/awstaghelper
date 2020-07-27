package elastiCacheHelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockedDescribeCacheClustersPages struct {
	elasticacheiface.ElastiCacheAPI
	resp elasticache.DescribeCacheClustersOutput
}

func (m *mockedDescribeCacheClustersPages) DescribeCacheClustersPages(input *elasticache.DescribeCacheClustersInput, pageFunc func(*elasticache.DescribeCacheClustersOutput, bool) bool) error {
	pageFunc(&m.resp, true)
	return nil
}

func TestGetInstances(t *testing.T) {
	cases := []*mockedDescribeCacheClustersPages{
		{
			resp: describeCacheClustersResponse,
		},
	}

	expectedResult := describeCacheClustersResponse.CacheClusters

	for _, c := range cases {
		t.Run("GetInstances", func(t *testing.T) {
			result := getInstances(c)
			assert := assert.New(t)
			assert.EqualValues(expectedResult, result)
		})

	}
}

var describeCacheClustersResponse = elasticache.DescribeCacheClustersOutput{
	CacheClusters: []*elasticache.CacheCluster{
		{
			CacheClusterId: aws.String("test-cluster-1"),
		},
		{
			CacheClusterId: aws.String("test-cluster-2"),
		},
	},
}
