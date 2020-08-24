package pkg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"log"
	"strings"
)

// getRDSInstances return all rds instances from specified region
func getRDSInstances(client rdsiface.RDSAPI) []*rds.DBInstance {
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

// ParseRDSTags parse output from getRDSInstances and return arn and specified tags.
func ParseRDSTags(tagsToRead string, client rdsiface.RDSAPI) [][]string {
	instancesOutput := getRDSInstances(client)
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

// TagRDS tag rds instances. Take as input data from csv file. Where first column arn
func TagRDS(csvData [][]string, client rdsiface.RDSAPI) {
	var tags []*rds.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &rds.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &rds.AddTagsToResourceInput{
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
