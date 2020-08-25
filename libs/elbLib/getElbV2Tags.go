package elbLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"log"
	"strings"
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
		log.Fatal("Not able to get load balancers", err)
		return nil
	}
	return result
}

// ParseElbV2Tags parse output from getInstances and return instances id and specified tags.
func ParseElbV2Tags(tagsToRead string, client elbv2iface.ELBV2API) [][]string {
	instancesOutput := getElbV2(client)
	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, elb := range instancesOutput {
		elbTags, err := client.DescribeTags(&elbv2.DescribeTagsInput{ResourceArns: []*string{elb.LoadBalancerArn}})
		if err != nil {
			fmt.Println("Not able to get load balancer tags", err)
		}
		tags := map[string]string{}
		for _, tagsToWrite := range elbTags.TagDescriptions {
			for _, tag := range tagsToWrite.Tags {
				tags[*tag.Key] = *tag.Value
			}
		}

		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*elb.LoadBalancerArn}, resultTags...))
	}
	return rows
}
