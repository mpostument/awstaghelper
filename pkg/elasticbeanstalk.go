package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk/elasticbeanstalkiface"
)

// getEBEnvironments return all elastic bean stalk environments from specified region
func getEBEnvironments(client elasticbeanstalkiface.ElasticBeanstalkAPI) *elasticbeanstalk.EnvironmentDescriptionsMessage {
	input := &elasticbeanstalk.DescribeEnvironmentsInput{}

	result, err := client.DescribeEnvironments(input)
	if err != nil {
		log.Fatal("Not able to get list of elastic bean stalk environments ", err)
	}

	return result
}

// ParseEBTags parse output from getEBInstances and return eb id and specified tags.
func ParseEBTags(tagsToRead string, client elasticbeanstalkiface.ElasticBeanstalkAPI) [][]string {
	instancesOutput := getEBEnvironments(client)
	rows := addHeadersToCsv(tagsToRead, "Arn")
	for _, ebEnv := range instancesOutput.Environments {
		ebTags, err := client.ListTagsForResource(&elasticbeanstalk.ListTagsForResourceInput{ResourceArn: ebEnv.EnvironmentArn})
		if err != nil {
			fmt.Println("Not able to get elastic bean stalk tags ", err)
		}
		tags := map[string]string{}
		for _, tag := range ebTags.ResourceTags {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *ebEnv.EnvironmentArn)
	}
	return rows
}

// TagEbEnvironments tag eb environments. Take as input data from csv file. Where first column is arn
func TagEbEnvironments(csvData [][]string, client elasticbeanstalkiface.ElasticBeanstalkAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags []*elasticbeanstalk.Tag
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &elasticbeanstalk.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &elasticbeanstalk.UpdateTagsForResourceInput{
			ResourceArn: aws.String(csvData[r][0]),
			TagsToAdd:   tags,
		}

		_, err := client.UpdateTagsForResource(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
