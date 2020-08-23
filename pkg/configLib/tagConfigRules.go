package configLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
)

// TagConfigRule tag config rules. Take as input data from csv file. Where first column Arn
func TagConfigRule(csvData [][]string, client configserviceiface.ConfigServiceAPI) {
	var tags []*configservice.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &configservice.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &configservice.TagResourceInput{
			ResourceArn: aws.String(csvData[r][0]),
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
