package pkg

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
)

// getASGs return all ASGs from specified region
func getASGs(client autoscalingiface.AutoScalingAPI) []*autoscaling.Group {
	input := &autoscaling.DescribeAutoScalingGroupsInput{}

	var result []*autoscaling.Group

	err := client.DescribeAutoScalingGroupsPages(input,
		func(page *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
			result = append(result, page.AutoScalingGroups...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get ASGs ", err)
	}
	return result
}

// ParseASGTags parse output from getASGs and return ASG name and specified tags.
func ParseASGTags(tagsToRead string, client autoscalingiface.AutoScalingAPI) [][]string {
	asgsOutput := getASGs(client)
	rows := addHeadersToCsv(tagsToRead, "AutoScalingGroupName")
	for _, asg := range asgsOutput {
		tags := map[string]string{}
		for _, tag := range asg.Tags {
			tags[*tag.Key] = *tag.Value + fmt.Sprintf("|Propagate=%t", *tag.PropagateAtLaunch)
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *asg.AutoScalingGroupName)
	}
	return rows
}

// TagASG tag ASGs. Take as input data from csv file.
func TagASG(csvData [][]string, client autoscalingiface.AutoScalingAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags []*autoscaling.Tag
		for c := 1; c < len(csvData[0]); c++ {
			val := strings.Split(csvData[r][c], "|Propagate=")
			var propagate bool
			var err error

			switch len(val) {
			case 1:
				{
					propagate = true
				}
			case 2:
				{
					propagate, err = strconv.ParseBool(val[1])
					if err != nil {
						propagate = true
					}
				}
			default:
				{
					log.Printf("Invalid CSV format: %v", csvData[r][c])
				}
			}

			tags = append(tags, &autoscaling.Tag{
				Key:               &csvData[0][c],
				Value:             &val[0],
				PropagateAtLaunch: &propagate,
				ResourceId:        &csvData[r][0],
				ResourceType:      aws.String("auto-scaling-group"),
			})
		}

		input := &autoscaling.CreateOrUpdateTagsInput{
			Tags: tags,
		}

		_, err := client.CreateOrUpdateTags(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
