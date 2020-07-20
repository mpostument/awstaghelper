package ec2Helper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// TagEc2 tag instances. Take as input data from csv file. Where first column id
func TagEc2(csvData [][]string, session session.Session) {
	client := ec2.New(&session)
	var tags []*ec2.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &ec2.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &ec2.CreateTagsInput{
			Resources: []*string{
				aws.String(csvData[r][0]),
			},
			Tags: tags,
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
