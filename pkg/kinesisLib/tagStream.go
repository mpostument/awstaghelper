package kinesisLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

// TagStream tag kinesis stream. Take as input data from csv file. Where first column id
func TagStream(csvData [][]string, client kinesisiface.KinesisAPI) {

	tags := make(map[string]*string)
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags[csvData[0][c]] = &csvData[r][c]
		}

		input := &kinesis.AddTagsToStreamInput{
			StreamName: aws.String(csvData[r][0]),
			Tags:       tags,
		}

		_, err := client.AddTagsToStream(input)
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
