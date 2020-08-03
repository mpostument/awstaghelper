package elasticSearchLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice/elasticsearchserviceiface"
)

// TagElasticSearch tag instances. Take as input data from csv file. Where first column id
func TagElasticSearch(csvData [][]string, client elasticsearchserviceiface.ElasticsearchServiceAPI) {

	var tags []*elasticsearchservice.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &elasticsearchservice.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &elasticsearchservice.AddTagsInput{
			ARN:     aws.String(csvData[r][0]),
			TagList: tags,
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
