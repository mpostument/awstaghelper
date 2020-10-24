package pkg

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// getSecurityGroups return all security groups from specified region
func getSecurityGroups(client ec2iface.EC2API) []*ec2.SecurityGroup {
	input := &ec2.DescribeSecurityGroupsInput{}

	var result []*ec2.SecurityGroup

	err := client.DescribeSecurityGroupsPages(input,
		func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			result = append(result, page.SecurityGroups...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get security groups ", err)
	}
	return result
}

// ParseSecurityGroupTags parse output from getSecurityGroups and return SG ids and specified tags.
func ParseSecurityGroupTags(tagsToRead string, client ec2iface.EC2API) [][]string {
	groupsOutput := getSecurityGroups(client)
	rows := addHeadersToCsv(tagsToRead, "Id")
	for _, sg := range groupsOutput {
		tags := map[string]string{}
		for _, tag := range sg.Tags {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *sg.GroupId)
	}
	return rows
}

// TagSecurityGroups tag security groups. Take as input data from csv file. Where first column id
func TagSecurityGroups(csvData [][]string, client ec2iface.EC2API) {
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
