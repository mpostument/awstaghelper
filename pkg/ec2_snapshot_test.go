/*
Copyright © 2024 Jaemok Hong jaemokhong@lguplus.co.kr
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
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/stretchr/testify/assert"
)

type mockedSnapshot struct {
	ec2iface.EC2API
	stsiface.STSAPI
	respGetCallerIdentity sts.GetCallerIdentityOutput
	respDescribeSnapshots ec2.DescribeSnapshotsOutput
}

// AWS Mocks
func (m *mockedSnapshot) DescribeSnapshotsPages(
	input *ec2.DescribeSnapshotsInput,
	pageFunc func(*ec2.DescribeSnapshotsOutput, bool) bool) error {
	pageFunc(&m.respDescribeSnapshots, true)
	return nil
}

func (m *mockedSnapshot) GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return &m.respGetCallerIdentity, nil
}

// Tests
func TestGetSnapshots(t *testing.T) {
	cases := []*mockedSnapshot{
		{
			respDescribeSnapshots: parseSnapshotTagsResponse,
		},
	}

	expectedResult := parseSnapshotTagsResponse.Snapshots

	for _, c := range cases {
		t.Run("getSnapshot", func(t *testing.T) {
			result := getSnapshot("callerIdentity", c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})

	}
}

func TestParseSnapshotTags(t *testing.T) {
	cases := []*mockedSnapshot{
		{
			respGetCallerIdentity: getSnapshotCallerIdentityResponse,
			respDescribeSnapshots: parseSnapshotTagsResponse,
		},
	}
	expectedResult := [][]string{
		{"Id", "Name", "Environment", "Owner"},
		{"snapshot-test", "TestSnapshot1", "Test", ""},
	}
	for _, c := range cases {
		t.Run("ParseSnapshotTags", func(t *testing.T) {
			result := ParseSnapshotTags("Name,Environment,Owner", c, c)
			assertions := assert.New(t)
			assertions.EqualValues(expectedResult, result)
		})
	}
}

var getSnapshotCallerIdentityResponse = sts.GetCallerIdentityOutput{
	Account: aws.String("666666666"),
}

var parseSnapshotTagsResponse = ec2.DescribeSnapshotsOutput{
	Snapshots: []*ec2.Snapshot{
		{
			Description: aws.String("testSg"),
			SnapshotId:  aws.String("snapshot-test"),
			VolumeId:    aws.String("testVolumId"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("TestSnapshot1"),
				},
				{
					Key:   aws.String("Environment"),
					Value: aws.String("Test"),
				},
			},
		},
	},
}
