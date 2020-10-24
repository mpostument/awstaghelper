package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// getBuckets return all s3 buckets from specified region
func getBuckets(client s3iface.S3API) *s3.ListBucketsOutput {
	input := &s3.ListBucketsInput{}

	result, err := client.ListBuckets(input)
	if err != nil {
		log.Fatal("Not able to get list of S3 buckets ", err)
	}

	return result
}

// ParseS3Tags parse output from getBuckets and return instances id and specified tags.
func ParseS3Tags(tagsToRead string, client s3iface.S3API) [][]string {
	s3Output := getBuckets(client)
	rows := addHeadersToCsv(tagsToRead, "Name")
	for _, bucket := range s3Output.Buckets {
		s3Tags, err := client.GetBucketTagging(&s3.GetBucketTaggingInput{Bucket: bucket.Name})
		if err != nil {
			if err.(awserr.Error).Code() == "NoSuchTagSet" {
				fmt.Println("Tag set for bucket", *bucket.Name, "doesn't exist")
			} else if err.(awserr.Error).Code() == "AuthorizationHeaderMalformed" {
				fmt.Println("Bucket ", *bucket.Name, "is not in your region", "region")
			} else {
				fmt.Println("Not able to get tags ", err)
			}
		}
		tags := map[string]string{}
		for _, tag := range s3Tags.TagSet {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *bucket.Name)
	}
	return rows
}

// TagS3 tag instances. Take as input data from csv file. Where first column id
func TagS3(csvData [][]string, client s3iface.S3API) {
	for r := 1; r < len(csvData); r++ {
		var tags []*s3.Tag
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
		if awsErrorHandle(err) {
			return
		}
	}
}
