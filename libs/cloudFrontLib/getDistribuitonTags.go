package cloudFrontLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"log"
	"strings"
)

// getDistributions return all config rules from specified region
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
