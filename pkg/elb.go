package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
)

// getElbV2 return all elbv2 (application and network) instances from specified region
func getElbV2(client elbv2iface.ELBV2API) []*elbv2.LoadBalancer {
	input := &elbv2.DescribeLoadBalancersInput{}

	var result []*elbv2.LoadBalancer

	err := client.DescribeLoadBalancersPages(input,
		func(page *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			result = append(result, page.LoadBalancers...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get load balancers ", err)
		return nil
	}
	return result
}

// ParseElbV2Tags parse output from getInstances and return instances id and specified tags.
func ParseElbV2Tags(tagsToRead string, client elbv2iface.ELBV2API) [][]string {
	instancesOutput := getElbV2(client)
	rows := addHeadersToCsv(tagsToRead, "Arn")
	for _, elb := range instancesOutput {
		elbTags, err := client.DescribeTags(&elbv2.DescribeTagsInput{ResourceArns: []*string{elb.LoadBalancerArn}})
		if err != nil {
			fmt.Println("Not able to get load balancer tags ", err)
		}
		tags := map[string]string{}
		for _, tagsToWrite := range elbTags.TagDescriptions {
			for _, tag := range tagsToWrite.Tags {
				tags[*tag.Key] = *tag.Value
			}
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *elb.LoadBalancerArn)
	}
	return rows
}

// TagElbV2 tag elbv2(application and network). Take as input data from csv file. Where first column id
func TagElbV2(csvData [][]string, client elbv2iface.ELBV2API) {
	for r := 1; r < len(csvData); r++ {
		var tags []*elbv2.Tag
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &elbv2.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &elbv2.AddTagsInput{
			ResourceArns: []*string{aws.String(csvData[r][0])},
			Tags:         tags,
		}

		_, err := client.AddTags(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
