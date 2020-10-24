package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

// getCWAlarm return all CloudWatch alarms from specified region
func getCWAlarm(client cloudwatchiface.CloudWatchAPI) []*cloudwatch.MetricAlarm {
	input := &cloudwatch.DescribeAlarmsInput{}

	var result []*cloudwatch.MetricAlarm

	err := client.DescribeAlarmsPages(input,
		func(page *cloudwatch.DescribeAlarmsOutput, lastPage bool) bool {
			result = append(result, page.MetricAlarms...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get alarms ", err)
		return nil
	}
	return result
}

// ParseCwAlarmTags parse output from getCWAlarm and return alarm arn and specified tags.
func ParseCwAlarmTags(tagsToRead string, client cloudwatchiface.CloudWatchAPI) [][]string {
	instancesOutput := getCWAlarm(client)
	rows := addHeadersToCsv(tagsToRead, "Arn")
	for _, alarm := range instancesOutput {

		input := &cloudwatch.ListTagsForResourceInput{
			ResourceARN: alarm.AlarmArn,
		}
		cwLogTags, err := client.ListTagsForResource(input)
		if err != nil {
			fmt.Println("Not able to get alarm tags ", err)
		}
		tags := map[string]string{}
		for _, tag := range cwLogTags.Tags {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *alarm.AlarmArn)
	}
	return rows
}

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
		log.Fatal("Not able to get log groups ", err)
		return nil
	}
	return result
}

// ParseCwLogGroupTags parse output from getInstances and return logGroupName and specified tags.
func ParseCwLogGroupTags(tagsToRead string, client cloudwatchlogsiface.CloudWatchLogsAPI) [][]string {
	instancesOutput := getCWLogGroups(client)
	rows := addHeadersToCsv(tagsToRead, "LogGroupName")
	for _, logGroup := range instancesOutput {

		input := &cloudwatchlogs.ListTagsLogGroupInput{
			LogGroupName: logGroup.LogGroupName,
		}
		cwLogTags, err := client.ListTagsLogGroup(input)
		if err != nil {
			fmt.Println("Not able to get log group tags ", err)
		}
		tags := map[string]string{}
		for key, value := range cwLogTags.Tags {
			tags[key] = *value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *logGroup.LogGroupName)
	}
	return rows
}

// TagCloudWatchAlarm tag cloudwatch alarms. Take as input data from csv file. Where first column Arn
func TagCloudWatchAlarm(csvData [][]string, client cloudwatchiface.CloudWatchAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags []*cloudwatch.Tag
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &cloudwatch.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &cloudwatch.TagResourceInput{
			ResourceARN: aws.String(csvData[r][0]),
			Tags:        tags,
		}

		_, err := client.TagResource(input)
		if awsErrorHandle(err) {
			return
		}
	}
}

// TagCloudWatchLogGroups tag cloudwatch log groups. Take as input data from csv file. Where first column LogGroupName
func TagCloudWatchLogGroups(csvData [][]string, client cloudwatchlogsiface.CloudWatchLogsAPI) {
	for r := 1; r < len(csvData); r++ {
		tags := make(map[string]*string)
		for c := 1; c < len(csvData[0]); c++ {
			tags[csvData[0][c]] = &csvData[r][c]
		}

		input := &cloudwatchlogs.TagLogGroupInput{
			LogGroupName: aws.String(csvData[r][0]),
			Tags:         tags,
		}

		_, err := client.TagLogGroup(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
