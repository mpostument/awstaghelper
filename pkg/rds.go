package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
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
		log.Fatal("Not able to get rds instances", err)
		return nil
	}
	return result
}

// ParseRDSTags parse output from getRDSInstances and return arn and specified tags.
func ParseRDSTags(tagsToRead string, client rdsiface.RDSAPI) [][]string {
	instancesOutput := getRDSInstances(client)
	rows := addHeadersToCsv(tagsToRead, "Arn")
	for _, dbInstances := range instancesOutput {
		rdsTags, err := client.ListTagsForResource(&rds.ListTagsForResourceInput{ResourceName: dbInstances.DBInstanceArn})
		if err != nil {
			fmt.Println("Not able to get rds tags ", err)
		}
		tags := map[string]string{}
		for _, tag := range rdsTags.TagList {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *dbInstances.DBInstanceArn)
	}
	return rows
}

// TagRDS tag rds instances. Take as input data from csv file. Where first column arn
func TagRDS(csvData [][]string, client rdsiface.RDSAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags []*rds.Tag
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
		if awsErrorHandle(err) {
			return
		}
	}
}
