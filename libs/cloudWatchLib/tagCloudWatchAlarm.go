package cloudWatchLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

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
