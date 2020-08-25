package elastiCacheLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"log"
	"strings"
)

// getElastiCache return all ElastiCache from specified region
func getElastiCache(client elasticacheiface.ElastiCacheAPI) []*elasticache.CacheCluster {
	input := &elasticache.DescribeCacheClustersInput{}

	var result []*elasticache.CacheCluster

	err := client.DescribeCacheClustersPages(input,
		func(page *elasticache.DescribeCacheClustersOutput, lastPage bool) bool {
			result = append(result, page.CacheClusters...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get elasticache instances", err)
		return nil
	}
	return result
}

// ParseElastiCacheTags parse output from getElastiCache and return arn and specified tags.
func ParseElastiCacheTags(tagsToRead string, client elasticacheiface.ElastiCacheAPI, stsClient stsiface.STSAPI, region string) [][]string {
	instancesOutput := getElastiCache(client)
	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal("Not able to get account id", err)
	}
	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, elasticCacheInstance := range instancesOutput {

		clusterArn := fmt.Sprintf("arn:aws:elasticache:%s:%s:cluster:%s",
			region, *callerIdentity.Account, *elasticCacheInstance.CacheClusterId)

		input := &elasticache.ListTagsForResourceInput{
			ResourceName: aws.String(clusterArn),
		}
		elasticCacheTag, err := client.ListTagsForResource(input)
		if err != nil {
			fmt.Println("Not able to get elasticache tags", err)
		}
		tags := map[string]string{}
		for _, tag := range elasticCacheTag.TagList {
			tags[*tag.Key] = *tag.Value
		}

		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{clusterArn}, resultTags...))
	}
	return rows
}
