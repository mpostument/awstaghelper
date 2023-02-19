/*
Copyright © 2023 Cristian Magherusan-Stanciu cristi@leanercloud.com
Copyright © 2020 Maksym Postument 777rip777@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
)

type mockedEBS struct {
	ec2iface.EC2API
	respDescribeVolumes ec2.DescribeVolumesOutput
}

// AWS Mocks
func (m *mockedEBS) DescribeVolumesPages(input *ec2.DescribeVolumesInput,
	pageFunc func(*ec2.DescribeVolumesOutput, bool) bool) error {
	pageFunc(&m.respDescribeVolumes, true)
	return nil
}

var ParseEBSVolumeTagsResponse = ec2.DescribeVolumesOutput{
	Volumes: []*ec2.Volume{
		{
			VolumeId: aws.String("vol-1"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("Volume1"),
				},
				{
					Key:   aws.String("Environment"),
					Value: aws.String("Test"),
				},
			},
		},
		{
			VolumeId: aws.String("vol-2"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("Volume2"),
				},
				{
					Key:   aws.String("Environment"),
					Value: aws.String("Dev"),
				},
			},
		},
	},
}

func Test_getEBSVolumes(t *testing.T) {
	cases := []*mockedEBS{
		{
			respDescribeVolumes: ParseEBSVolumeTagsResponse,
		},
	}

	expectedResult := ParseEBSVolumeTagsResponse.Volumes
	for _, c := range cases {
		t.Run("getEBSVolumes", func(t *testing.T) {
			result := getEBSVolumes(c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseEBSVolumeTags(t *testing.T) {
	cases := []*mockedEBS{
		{
			respDescribeVolumes: ParseEBSVolumeTagsResponse,
		},
	}
	expectedResult := [][]string{
		{"VolumeId", "Name", "Environment"},
		{"vol-1", "Volume1", "Test"},
		{"vol-2", "Volume2", "Dev"},
	}
	for _, c := range cases {
		t.Run("ParseEBSVolumeTags", func(t *testing.T) {
			result := ParseEBSVolumeTags("Name,Environment", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})
	}
}
