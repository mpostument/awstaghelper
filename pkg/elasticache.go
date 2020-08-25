package pkg

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/elasticache/elasticacheiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// getElastiCacheClusters return all ElastiCache from specified region
func getElastiCacheClusters(client elasticacheiface.ElastiCacheAPI) []*elasticache.CacheCluster {
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

// ParseElastiCacheClusterTags parse output from getElastiCacheClusters and return arn and specified tags.
func ParseElastiCacheClusterTags(tagsToRead string, client elasticacheiface.ElastiCacheAPI, stsClient stsiface.STSAPI, region string) [][]string {
	instancesOutput := getElastiCacheClusters(client)
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

// TagElastiCache tag instances. Take as input data from csv file. Where first column id
func TagElastiCache(csvData [][]string, client elasticacheiface.ElastiCacheAPI) {

	var tags []*elasticache.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &elasticache.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &elasticache.AddTagsToResourceInput{
			ResourceName: aws.String(csvData[r][0]),
			Tags:         tags,
		}

		_, err := client.AddTagsToResource(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}
	}
}
