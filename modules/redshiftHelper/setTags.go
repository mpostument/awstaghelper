package redshiftHelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/redshift/redshiftiface"
)

// TagRedsfhit tag rds instances. Take as input data from csv file. Where first column arn
func TagRedsfhit(csvData [][]string, client redshiftiface.RedshiftAPI) {
	var tags []*redshift.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &redshift.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &redshift.CreateTagsInput{
			ResourceName: aws.String(csvData[r][0]),
			Tags:         tags,
		}

		_, err := client.CreateTags(input)
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
