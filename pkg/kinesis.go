package pkg

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/firehose/firehoseiface"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

// getFirehoses return all firehoses from specified region
func getFirehoses(client firehoseiface.FirehoseAPI) *firehose.ListDeliveryStreamsOutput {
	input := &firehose.ListDeliveryStreamsInput{Limit: aws.Int64(10000)}

	result, err := client.ListDeliveryStreams(input)
	if err != nil {
		log.Fatal("Not able to get list of buckets", err)
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
			fmt.Println("Not able to get kinesis tags", err)
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

// getStreams return all streams from specified region
func getStreams(client kinesisiface.KinesisAPI) []*string {
	input := &kinesis.ListStreamsInput{}

	var result []*string

	err := client.ListStreamsPages(input,
		func(page *kinesis.ListStreamsOutput, lastPage bool) bool {
			result = append(result, page.StreamNames...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get kinesis streams", err)
		return nil
	}
	return result
}

// ParseKinesisTags parse output from getStreams and return kinesis arn and specified tags.
func ParseKinesisTags(tagsToRead string, client kinesisiface.KinesisAPI) [][]string {
	instancesOutput := getStreams(client)
	var rows [][]string
	headers := []string{"Name"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, stream := range instancesOutput {

		input := &kinesis.ListTagsForStreamInput{
			StreamName: stream,
		}
		distributionTags, err := client.ListTagsForStream(input)
		if err != nil {
			fmt.Println("Not able to get kinesis tags", err)
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

// TagFirehose tag kinesis firehose. Take as input data from csv file. Where first column name
func TagFirehose(csvData [][]string, client firehoseiface.FirehoseAPI) {
	var tags []*firehose.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &firehose.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &firehose.TagDeliveryStreamInput{
			DeliveryStreamName: aws.String(csvData[r][0]),
			Tags:               tags,
		}

		_, err := client.TagDeliveryStream(input)
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

// TagKinesisStream tag kinesis stream. Take as input data from csv file. Where first column id
func TagKinesisStream(csvData [][]string, client kinesisiface.KinesisAPI) {

	tags := make(map[string]*string)
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags[csvData[0][c]] = &csvData[r][c]
		}

		input := &kinesis.AddTagsToStreamInput{
			StreamName: aws.String(csvData[r][0]),
			Tags:       tags,
		}

		_, err := client.AddTagsToStream(input)
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
