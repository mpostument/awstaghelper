package elastiCacheHelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

// TagElasticache tag instances. Take as input data from csv file. Where first column id
func TagElasticache(csvData [][]string, session session.Session) {
	client := elasticache.New(&session)

	var tags []*elasticache.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &elasticache.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &elasticache.AddTagsToResourceInput{
			ResourceName: aws.String(csvData[r][0]),
			Tags:     tags,
		}

		_, err := client.AddTagsToResource(input)
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
