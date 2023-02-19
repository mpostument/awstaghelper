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
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// getEBSVolumes returns all EBS volumes from specified region
func getEBSVolumes(client ec2iface.EC2API) []*ec2.Volume {
	input := &ec2.DescribeVolumesInput{}

	var result []*ec2.Volume

	err := client.DescribeVolumesPages(input,
		func(page *ec2.DescribeVolumesOutput, lastPage bool) bool {
			result = append(result, page.Volumes...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get EBS volumes ", err)
	}
	return result
}

// ParseEBSVolumeTags parse output from getEBSVolumes and return volume ID and specified tags.
func ParseEBSVolumeTags(tagsToRead string, client ec2iface.EC2API) [][]string {
	volumesOutput := getEBSVolumes(client)
	rows := addHeadersToCsv(tagsToRead, "VolumeId")
	for _, volume := range volumesOutput {
		tags := map[string]string{}
		for _, tag := range volume.Tags {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *volume.VolumeId)
	}
	return rows
}

// TagEBSVolumes tag EBS volumes. Take as input data from csv file. Where first column is volume ID.
func TagEBSVolumes(csvData [][]string, client ec2iface.EC2API) {
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
