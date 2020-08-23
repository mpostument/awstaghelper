package kinesisLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/firehose/firehoseiface"
)

// TagFirehose tag kinesis firehose. Take as input data from csv file. Where first column name
func TagFirehose(csvData [][]string, client firehoseiface.FirehoseAPI) {
	var tags []*firehose.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &firehose.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &firehose.TagDeliveryStreamInput{
			DeliveryStreamName: aws.String(csvData[r][0]),
			Tags:               tags,
		}

		_, err := client.TagDeliveryStream(input)
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
