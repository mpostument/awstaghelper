package kinesisLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"log"
	"strings"
)

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
