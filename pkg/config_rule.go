package pkg

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
)

// getConfigRules return all config rules from specified region
func getConfigRules(client configserviceiface.ConfigServiceAPI) *configservice.DescribeConfigRulesOutput {
	input := &configservice.DescribeConfigRulesInput{}

	result, err := client.DescribeConfigRules(input)
	if err != nil {
		log.Fatal("Not able to get config rules", err)
		return nil
	}
	return result
}

// ParseConfigRuleTags parse output from getCWAlarm and return alarm arn and specified tags.
func ParseConfigRuleTags(tagsToRead string, client configserviceiface.ConfigServiceAPI) [][]string {
	instancesOutput := getConfigRules(client)
	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, rule := range instancesOutput.ConfigRules {

		input := &configservice.ListTagsForResourceInput{
			ResourceArn: rule.ConfigRuleArn,
		}
		configTags, err := client.ListTagsForResource(input)
		if err != nil {
			fmt.Println("Not able to get config rule tags", err)
		}
		tags := map[string]string{}
		for _, tag := range configTags.Tags {
			tags[*tag.Key] = *tag.Value
		}

		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*rule.ConfigRuleArn}, resultTags...))
	}
	return rows
}

// TagConfigRule tag config rules. Take as input data from csv file. Where first column Arn
func TagConfigRule(csvData [][]string, client configserviceiface.ConfigServiceAPI) {
	var tags []*configservice.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &configservice.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &configservice.TagResourceInput{
			ResourceArn: aws.String(csvData[r][0]),
			Tags:        tags,
		}

		_, err := client.TagResource(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}
	}
}
