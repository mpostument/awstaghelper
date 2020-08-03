package ec2Lib

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"log"
	"strings"
)

// getInstances return all ec2 instances from specified region
func getInstances(client ec2iface.EC2API) []*ec2.Reservation {
	input := &ec2.DescribeInstancesInput{}

	var result []*ec2.Reservation

	err := client.DescribeInstancesPages(input,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			result = append(result, page.Reservations...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get instances", err)
	}
	return result
}

// ParseEc2Tags parse output from getInstances and return instances id and specified tags.
func ParseEc2Tags(tagsToRead string, client ec2iface.EC2API) [][]string {
	instancesOutput := getInstances(client)
	var rows [][]string
	headers := []string{"Id"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, reservation := range instancesOutput {
		for _, instance := range reservation.Instances {
			tags := map[string]string{}
			for _, tag := range instance.Tags {
				tags[*tag.Key] = *tag.Value
			}
			var resultTags []string
			for _, key := range strings.Split(tagsToRead, ",") {
				resultTags = append(resultTags, tags[key])
			}
			rows = append(rows, append([]string{*instance.InstanceId}, resultTags...))
		}
	}
	return rows
}
