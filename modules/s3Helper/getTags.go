package s3Helper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"strings"
)

// getBuckets return all s3 buckets from specified region
func getBuckets(session session.Session) *s3.ListBucketsOutput {
	client := s3.New(&session)
	input := &s3.ListBucketsInput{}

	result, err := client.ListBuckets(input)
	if err != nil {
		log.Fatal("Not able to get list of buckets", err)
	}

	return result
}

// ParseS3Tags parse output from getBuckets and return instances id and specified tags.
func ParseS3Tags(tagsToRead string, session session.Session) [][]string {
	s3Output := getBuckets(session)
	client := s3.New(&session)
	var rows [][]string
	headers := []string{"Name"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, bucket := range s3Output.Buckets {
		s3Tags, err := client.GetBucketTagging(&s3.GetBucketTaggingInput{Bucket: bucket.Name})
		if err != nil {
			if err.(awserr.Error).Code() == "NoSuchTagSet" {
				fmt.Println("Tag set for bucket", *bucket.Name, "doesn't exist")
			} else if err.(awserr.Error).Code() == "AuthorizationHeaderMalformed" {
				fmt.Println("Bucket ", *bucket.Name, "is not in", *client.Config.Region, "region")
			} else {
				fmt.Println("Not able to get tags", err)
			}
		}
		tags := map[string]string{}
		for _, tag := range s3Tags.TagSet {
			tags[*tag.Key] = *tag.Value
		}
		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*bucket.Name}, resultTags...))
	}
	return rows
}
