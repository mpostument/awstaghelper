package kinesisLib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/firehose/firehoseiface"
	"github.com/stretchr/testify/assert"
	"testing"
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
