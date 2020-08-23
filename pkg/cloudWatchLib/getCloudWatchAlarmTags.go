package cloudWatchLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
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
