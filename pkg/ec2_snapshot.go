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
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// getSecurityGroups return all security groups from specified region
func getSnapshot(callerIdentity string, client ec2iface.EC2API) []*ec2.Snapshot {
	input := &ec2.DescribeSnapshotsInput{
		OwnerIds: []*string{aws.String(callerIdentity)},
	}

	var result []*ec2.Snapshot

	err := client.DescribeSnapshotsPages(input,
		func(page *ec2.DescribeSnapshotsOutput, lastPage bool) bool {
			result = append(result, page.Snapshots...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get snapshots ", err)
	}
	return result
}

// ParseSecurityGroupTags parse output from getSecurityGroups and return SG ids and specified tags.
func ParseSnapshotTags(tagsToRead string, client ec2iface.EC2API, stsClient stsiface.STSAPI) [][]string {
	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal("Not able to get account id", err)
	}
	snapshotsOutput := getSnapshot(*callerIdentity.Account, client)

	rows := addHeadersToCsv(tagsToRead, "Id")
	for _, snapshot := range snapshotsOutput {
		tags := map[string]string{}
		for _, tag := range snapshot.Tags {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *snapshot.SnapshotId)
	}
	return rows
}

// TagSecurityGroups tag security groups. Take as input data from csv file. Where first column id
func TagSnapshot(csvData [][]string, client ec2iface.EC2API) {
	for r := 1; r < len(csvData); r++ {
		var tags []*ec2.Tag
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &ec2.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &ec2.CreateTagsInput{
			Resources: []*string{
				aws.String(csvData[r][0]),
			},
			Tags: tags,
		}

		_, err := client.CreateTags(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
