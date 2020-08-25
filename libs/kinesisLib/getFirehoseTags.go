package kinesisLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/firehose/firehoseiface"
	"log"
	"strings"
)

// getFirehoses return all firehoses from specified region
func getFirehoses(client firehoseiface.FirehoseAPI) *firehose.ListDeliveryStreamsOutput {
	input := &firehose.ListDeliveryStreamsInput{Limit: aws.Int64(10000)}

	result, err := client.ListDeliveryStreams(input)
	if err != nil {
		log.Fatal("Not able to get list of firehose streams", err)
	}
	return result
}

// ParseKinesisTags parse output from getFirehoses and return firehose name and specified tags.
func ParseFirehoseTags(tagsToRead string, client firehoseiface.FirehoseAPI) [][]string {
	instancesOutput := getFirehoses(client)
	var rows [][]string
	headers := []string{"Name"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, stream := range instancesOutput.DeliveryStreamNames {

		input := &firehose.ListTagsForDeliveryStreamInput{
			DeliveryStreamName: stream,
		}
		distributionTags, err := client.ListTagsForDeliveryStream(input)
		if err != nil {
			fmt.Println("Not able to get firehose tags", err)
		}
		tags := map[string]string{}
		for _, tag := range distributionTags.Tags {
			tags[*tag.Key] = *tag.Value
		}

		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*stream}, resultTags...))
	}
	return rows
}
