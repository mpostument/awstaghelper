package cloudFrontLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
)

// TagDistribution tag cloudfront distribution. Take as input data from csv file. Where first column Arn
func TagDistribution(csvData [][]string, client cloudfrontiface.CloudFrontAPI) {
	var tags cloudfront.Tags
	for r := 1; r < len(csvData); r++ {
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
