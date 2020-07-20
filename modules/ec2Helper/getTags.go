package ec2Helper

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"strings"
)

// getInstances return all ec2 instances from specified region
func getInstances(session session.Session) []*ec2.Reservation {
	client := ec2.New(&session)
	input := &ec2.DescribeInstancesInput{}
	//result, err := client.DescribeInstances(input)

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
func ParseEc2Tags(tagsToRead string, session session.Session) [][]string {
	instancesOutput := getInstances(session)
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
