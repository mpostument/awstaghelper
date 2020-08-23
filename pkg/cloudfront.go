package pkg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"log"
	"strings"
)

// getDistributions return all cloudfront distributions from specified region
func getDistributions(client cloudfrontiface.CloudFrontAPI) *cloudfront.ListDistributionsOutput {
	input := &cloudfront.ListDistributionsInput{}

	result, err := client.ListDistributions(input)
	if err != nil {
		log.Fatal("Not able to get distributions", err)
		return nil
	}
	return result
}

// ParseDistributionsTags parse output from getDistributions and return distribution arn and specified tags.
func ParseDistributionsTags(tagsToRead string, client cloudfrontiface.CloudFrontAPI) [][]string {
	instancesOutput := getDistributions(client)
	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, distribution := range instancesOutput.DistributionList.Items {

		input := &cloudfront.ListTagsForResourceInput{
			Resource: distribution.ARN,
		}
		distributionTags, err := client.ListTagsForResource(input)
		if err != nil {
			fmt.Println("Not able to get log group tags", err)
		}
		tags := map[string]string{}
		for _, tag := range distributionTags.Tags.Items {
			tags[*tag.Key] = *tag.Value
		}

		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*distribution.ARN}, resultTags...))
	}
	return rows
}

// TagDistribution tag cloudfront distribution. Take as input data from csv file. Where first column Arn
func TagDistribution(csvData [][]string, client cloudfrontiface.CloudFrontAPI) {
	var tags cloudfront.Tags
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags.Items = append(tags.Items, &cloudfront.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &cloudfront.TagResourceInput{
			Resource: aws.String(csvData[r][0]),
			Tags:     &tags,
		}

		_, err := client.TagResource(input)
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
