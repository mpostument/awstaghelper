package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
)

// getDistributions return all cloudfront distributions from specified region
func getDistributions(client cloudfrontiface.CloudFrontAPI) *cloudfront.ListDistributionsOutput {
	input := &cloudfront.ListDistributionsInput{}

	result, err := client.ListDistributions(input)
	if err != nil {
		log.Fatal("Not able to get distributions ", err)
		return nil
	}
	return result
}

// ParseDistributionsTags parse output from getDistributions and return distribution arn and specified tags.
func ParseDistributionsTags(tagsToRead string, client cloudfrontiface.CloudFrontAPI) [][]string {
	instancesOutput := getDistributions(client)
	rows := addHeadersToCsv(tagsToRead, "Arn")
	for _, distribution := range instancesOutput.DistributionList.Items {

		input := &cloudfront.ListTagsForResourceInput{
			Resource: distribution.ARN,
		}
		distributionTags, err := client.ListTagsForResource(input)
		if err != nil {
			fmt.Println("Not able to get distributions tags ", err)
		}
		tags := map[string]string{}
		for _, tag := range distributionTags.Tags.Items {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *distribution.ARN)
	}
	return rows
}

// TagDistribution tag cloudfront distribution. Take as input data from csv file. Where first column Arn
func TagDistribution(csvData [][]string, client cloudfrontiface.CloudFrontAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags cloudfront.Tags
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
		if awsErrorHandle(err) {
			return
		}
	}
}
