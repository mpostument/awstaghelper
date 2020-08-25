package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/firehose/firehoseiface"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/stretchr/testify/assert"
)

type mockedFirehose struct {
	firehoseiface.FirehoseAPI
	respListFirehose firehose.ListDeliveryStreamsOutput
	respListTags     firehose.ListTagsForDeliveryStreamOutput
}

func (m *mockedFirehose) ListDeliveryStreams(*firehose.ListDeliveryStreamsInput) (*firehose.ListDeliveryStreamsOutput, error) {
	return &m.respListFirehose, nil
}

func (m *mockedFirehose) ListTagsForDeliveryStream(*firehose.ListTagsForDeliveryStreamInput) (*firehose.ListTagsForDeliveryStreamOutput, error) {
	return &m.respListTags, nil
}

func TestGetFirehose(t *testing.T) {
	cases := []*mockedFirehose{
		{
			respListFirehose: listFirehoseOutputResponse,
		},
	}

	expectedResult := &listFirehoseOutputResponse

	for _, c := range cases {
		t.Run("getFirehoses", func(t *testing.T) {
			result := getFirehoses(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseFirehoseTags(t *testing.T) {
	cases := []*mockedFirehose{
		{
			respListFirehose: listFirehoseOutputResponse,
			respListTags:     listFirehoseTagsResponse,
		},
	}

	expectedResult := [][]string{
		{"Name", "Environment", "Owner"},
		{"test-firehose-1", "test", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseFirehoseTags", func(t *testing.T) {
			result := ParseFirehoseTags("Environment,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var listFirehoseOutputResponse = firehose.ListDeliveryStreamsOutput{
	DeliveryStreamNames: []*string{
		aws.String("test-firehose-1"),
	},
}

var listFirehoseTagsResponse = firehose.ListTagsForDeliveryStreamOutput{
	Tags: []*firehose.Tag{
		{
			Key:   aws.String("Environment"),
			Value: aws.String("test"),
		},
		{
			Key:   aws.String("Owner"),
			Value: aws.String("mpostument"),
		},
	},
}

type mockedStream struct {
	kinesisiface.KinesisAPI
	respListStreams kinesis.ListStreamsOutput
	respListTags    kinesis.ListTagsForStreamOutput
}

func (m *mockedStream) ListStreamsPages(input *kinesis.ListStreamsInput, pageFunc func(*kinesis.ListStreamsOutput, bool) bool) error {
	pageFunc(&m.respListStreams, true)
	return nil
}

func (m *mockedStream) ListTagsForStream(*kinesis.ListTagsForStreamInput) (*kinesis.ListTagsForStreamOutput, error) {
	return &m.respListTags, nil
}

func TestGetStream(t *testing.T) {
	cases := []*mockedStream{
		{
			respListStreams: listStreamsOutputResponse,
		},
	}

	expectedResult := listStreamsOutputResponse.StreamNames

	for _, c := range cases {
		t.Run("getStreams", func(t *testing.T) {
			result := getStreams(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseKinesisTags(t *testing.T) {
	cases := []*mockedStream{
		{
			respListStreams: listStreamsOutputResponse,
			respListTags:    listTagsForStreamResponse,
		},
	}

	expectedResult := [][]string{
		{"Name", "Environment", "Owner"},
		{"test-stream-1", "test", "mpostument"},
	}

	for _, c := range cases {
		t.Run("ParseKinesisTags", func(t *testing.T) {
			result := ParseKinesisTags("Environment,Owner", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

var listStreamsOutputResponse = kinesis.ListStreamsOutput{
	StreamNames: []*string{
		aws.String("test-stream-1"),
	},
}

var listTagsForStreamResponse = kinesis.ListTagsForStreamOutput{
	Tags: []*kinesis.Tag{
		{
			Key:   aws.String("Environment"),
			Value: aws.String("test"),
		},
		{
			Key:   aws.String("Owner"),
			Value: aws.String("mpostument"),
		},
	},
}
