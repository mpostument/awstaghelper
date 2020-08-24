package elasticSearchLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice/elasticsearchserviceiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"log"
	"strings"
)

// getElasticSearch return all elasticsearch from specified region
func getElasticSearch(client elasticsearchserviceiface.ElasticsearchServiceAPI) *elasticsearchservice.ListDomainNamesOutput {
	input := &elasticsearchservice.ListDomainNamesInput{}

	result, err := client.ListDomainNames(input)
	if err != nil {
		log.Fatal("Not able to get elasticsearch instances", err)
	}
	return result
}

// ParseElasticSearchTags parse output from getElasticSearch and return arn and specified tags.
func ParseElasticSearchTags(tagsToRead string, client elasticsearchserviceiface.ElasticsearchServiceAPI, stsClient stsiface.STSAPI, region string) [][]string {
	instancesOutput := getElasticSearch(client)
	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal("Not able to get account id", err)
	}
	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, elasticCacheInstance := range instancesOutput.DomainNames {

		clusterArn := fmt.Sprintf("arn:aws:es:%s:%s:domain/%s",
			region, *callerIdentity.Account, *elasticCacheInstance.DomainName)

		input := &elasticsearchservice.ListTagsInput{
			ARN: aws.String(clusterArn),
		}
		elasticSearchTags, err := client.ListTags(input)
		if err != nil {
			fmt.Println("Not able to get elasticsearch tags", err)
		}
		tags := map[string]string{}
		for _, tag := range elasticSearchTags.TagList {
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
