package s3Helper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TagS3 tag instances. Take as input data from csv file. Where first column id
func TagS3(csvData [][]string, session session.Session) {
	client := s3.New(&session)
	var tags []*s3.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &s3.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &s3.PutBucketTaggingInput{
			Bucket: aws.String(csvData[r][0]),
			Tagging: &s3.Tagging{
				TagSet: tags,
			},
		}

		_, err := client.PutBucketTagging(input)
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
