package rdsHelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"log"
	"strings"
)

// getInstances return all rds instances from specified region
func getInstances(session session.Session) []*rds.DBInstance {
	client := rds.New(&session)
	input := &rds.DescribeDBInstancesInput{}

	var result []*rds.DBInstance

	err := client.DescribeDBInstancesPages(input,
		func(page *rds.DescribeDBInstancesOutput, lastPage bool) bool {
			result = append(result, page.DBInstances...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get instances", err)
		return nil
	}
	return result
}

// ParseRdsTags parse output from getInstances and return arn and specified tags.
func ParseRdsTags(tagsToRead string, session session.Session) [][]string {
	instancesOutput := getInstances(session)
	client := rds.New(&session)
	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, dbInstances := range instancesOutput {
		rdsTags, err := client.ListTagsForResource(&rds.ListTagsForResourceInput{ResourceName: dbInstances.DBInstanceArn})
		if err != nil {
			fmt.Println("Not able to get rds tags", err)
		}
		tags := map[string]string{}
		for _, tag := range rdsTags.TagList {
			tags[*tag.Key] = *tag.Value
		}
		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*dbInstances.DBInstanceArn}, resultTags...))
	}
	return rows
}
