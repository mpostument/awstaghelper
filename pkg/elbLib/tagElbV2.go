package elbLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
)

// TagElbV2 tag elbv2(application and network). Take as input data from csv file. Where first column id
func TagElbV2(csvData [][]string, client elbv2iface.ELBV2API) {

	var tags []*elbv2.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &elbv2.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &elbv2.AddTagsInput{
			ResourceArns: []*string{aws.String(csvData[r][0])},
			Tags:         tags,
		}

		_, err := client.AddTags(input)
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
