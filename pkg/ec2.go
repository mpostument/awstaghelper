package pkg

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// getEC2Instances return all ec2 instances from specified region
func getEC2Instances(client ec2iface.EC2API) []*ec2.Reservation {
	input := &ec2.DescribeInstancesInput{}

	var result []*ec2.Reservation

	err := client.DescribeInstancesPages(input,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			result = append(result, page.Reservations...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get EC2 instances ", err)
	}
	return result
}

// ParseEC2Tags parse output from getEC2Instances and return instances id and specified tags.
func ParseEC2Tags(tagsToRead string, client ec2iface.EC2API) [][]string {
	instancesOutput := getEC2Instances(client)
	rows := addHeadersToCsv(tagsToRead, "Id")
	for _, reservation := range instancesOutput {
		for _, instance := range reservation.Instances {
			tags := map[string]string{}
			for _, tag := range instance.Tags {
				tags[*tag.Key] = *tag.Value
			}
			rows = addTagsToCsv(tagsToRead, tags, rows, *instance.InstanceId)
		}
	}
	return rows
}

// TagEc2 tag instances. Take as input data from csv file. Where first column id
func TagEc2(csvData [][]string, client ec2iface.EC2API) {
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
