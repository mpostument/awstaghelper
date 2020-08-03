package cloudWatchLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"log"
	"strings"
)

// getCWLogsInstances return all CloudWatch Log groups from specified region
func getCWLogGroups(client cloudwatchlogsiface.CloudWatchLogsAPI) []*cloudwatchlogs.LogGroup {
	input := &cloudwatchlogs.DescribeLogGroupsInput{}

	var result []*cloudwatchlogs.LogGroup

	err := client.DescribeLogGroupsPages(input,
		func(page *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
			result = append(result, page.LogGroups...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get log groups", err)
		return nil
	}
	return result
}

// ParseCwLogGroupTags parse output from getInstances and return logGroupName and specified tags.
func ParseCwLogGroupTags(tagsToRead string, client cloudwatchlogsiface.CloudWatchLogsAPI) [][]string {
	instancesOutput := getCWLogGroups(client)
	var rows [][]string
	headers := []string{"LogGroupName"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, logGroup := range instancesOutput {

		input := &cloudwatchlogs.ListTagsLogGroupInput{
			LogGroupName: logGroup.LogGroupName,
		}
		cwLogTags, err := client.ListTagsLogGroup(input)
		if err != nil {
			fmt.Println("Not able to get log group tags", err)
		}
		tags := map[string]string{}
		for key, value := range cwLogTags.Tags {
			tags[key] = *value
		}

		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*logGroup.LogGroupName}, resultTags...))
	}
	return rows
}
