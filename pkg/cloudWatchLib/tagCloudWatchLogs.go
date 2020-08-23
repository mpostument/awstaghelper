package cloudWatchLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

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
