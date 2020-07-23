package elasticSearchHelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/sts"
	"log"
	"strings"
)

// getInstances return all elasticsearch from specified region
func getInstances(session session.Session) *elasticsearchservice.ListDomainNamesOutput {
	client := elasticsearchservice.New(&session)
	input := &elasticsearchservice.ListDomainNamesInput{}

	result, err := client.ListDomainNames(input)
	if err != nil {
		log.Fatal("Not able to get instances", err)
	}
	return result
}

// ParseElasticSearchTags parse output from getInstances and return arn and specified tags.
func ParseElasticSearchTags(tagsToRead string, session session.Session) [][]string {
	instancesOutput := getInstances(session)
	client := elasticsearchservice.New(&session)
	stsClient := sts.New(&session)
	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal("Not able to get elasticsearchservice tags", err)
	}
	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, elasticCacheInstance := range instancesOutput.DomainNames {

		clusterArn := fmt.Sprintf("arn:aws:es:%s:%s:domain/%s",
			*session.Config.Region, *callerIdentity.Account, *elasticCacheInstance.DomainName)

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
