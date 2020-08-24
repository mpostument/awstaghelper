package pkg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"log"
	"strings"
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
		log.Fatal("Not able to get log groups", err)
		return nil
	}
	return result
}

// ParseCwAlarmTags parse output from getCWAlarm and return alarm arn and specified tags.
func ParseCwAlarmTags(tagsToRead string, client cloudwatchiface.CloudWatchAPI) [][]string {
	instancesOutput := getCWAlarm(client)
	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, alarm := range instancesOutput {

		input := &cloudwatch.ListTagsForResourceInput{
			ResourceARN: alarm.AlarmArn,
		}
		cwLogTags, err := client.ListTagsForResource(input)
		if err != nil {
			fmt.Println("Not able to get log group tags", err)
		}
		tags := map[string]string{}
		for _, tag := range cwLogTags.Tags {
			tags[*tag.Key] = *tag.Value
		}

		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*alarm.AlarmArn}, resultTags...))
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

// TagCloudWatchAlarm tag cloudwatch alarms. Take as input data from csv file. Where first column Arn
func TagCloudWatchAlarm(csvData [][]string, client cloudwatchiface.CloudWatchAPI) {
	var tags []*cloudwatch.Tag
	for r := 1; r < len(csvData); r++ {
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

// TagCloudWatchLogGroups tag cloudwatch log groups. Take as input data from csv file. Where first column LogGroupName
func TagCloudWatchLogGroups(csvData [][]string, client cloudwatchlogsiface.CloudWatchLogsAPI) {

	tags := make(map[string]*string)
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags[csvData[0][c]] = &csvData[r][c]
		}

		input := &cloudwatchlogs.TagLogGroupInput{
			LogGroupName: aws.String(csvData[r][0]),
			Tags:         tags,
		}

		_, err := client.TagLogGroup(input)
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
