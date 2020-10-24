package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice/elasticsearchserviceiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// getElasticSearchDomains return all elasticsearch from specified region
func getElasticSearchDomains(client elasticsearchserviceiface.ElasticsearchServiceAPI) *elasticsearchservice.ListDomainNamesOutput {
	input := &elasticsearchservice.ListDomainNamesInput{}

	result, err := client.ListDomainNames(input)
	if err != nil {
		log.Fatal("Not able to get elasticsearch instances ", err)
	}
	return result
}

// ParseElasticSearchTags parse output from getElasticSearchDomains and return arn and specified tags.
func ParseElasticSearchTags(tagsToRead string, client elasticsearchserviceiface.ElasticsearchServiceAPI, stsClient stsiface.STSAPI, region string) [][]string {
	instancesOutput := getElasticSearchDomains(client)
	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal("Not able to get account id ", err)
	}
	rows := addHeadersToCsv(tagsToRead, "Arn")
	for _, elasticCacheInstance := range instancesOutput.DomainNames {

		clusterArn := fmt.Sprintf("arn:aws:es:%s:%s:domain/%s",
			region, *callerIdentity.Account, *elasticCacheInstance.DomainName)

		input := &elasticsearchservice.ListTagsInput{
			ARN: aws.String(clusterArn),
		}
		elasticSearchTags, err := client.ListTags(input)
		if err != nil {
			fmt.Println("Not able to get elasticsearch tags ", err)
		}
		tags := map[string]string{}
		for _, tag := range elasticSearchTags.TagList {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, clusterArn)
	}
	return rows
}

// TagElasticSearch tag instances. Take as input data from csv file. Where first column id
func TagElasticSearch(csvData [][]string, client elasticsearchserviceiface.ElasticsearchServiceAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags []*elasticsearchservice.Tag
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &elasticsearchservice.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &elasticsearchservice.AddTagsInput{
			ARN:     aws.String(csvData[r][0]),
			TagList: tags,
		}

		_, err := client.AddTags(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
