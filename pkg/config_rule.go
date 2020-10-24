package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
)

// getConfigRules return all config rules from specified region
func getConfigRules(client configserviceiface.ConfigServiceAPI) *configservice.DescribeConfigRulesOutput {
	input := &configservice.DescribeConfigRulesInput{}

	result, err := client.DescribeConfigRules(input)
	if err != nil {
		log.Fatal("Not able to get config rules ", err)
		return nil
	}
	return result
}

// ParseConfigRuleTags parse output from getCWAlarm and return alarm arn and specified tags.
func ParseConfigRuleTags(tagsToRead string, client configserviceiface.ConfigServiceAPI) [][]string {
	instancesOutput := getConfigRules(client)
	rows := addHeadersToCsv(tagsToRead, "Arn")
	for _, rule := range instancesOutput.ConfigRules {

		input := &configservice.ListTagsForResourceInput{
			ResourceArn: rule.ConfigRuleArn,
		}
		configTags, err := client.ListTagsForResource(input)
		if err != nil {
			fmt.Println("Not able to get config rule tags ", err)
		}
		tags := map[string]string{}
		for _, tag := range configTags.Tags {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *rule.ConfigRuleArn)
	}
	return rows
}

// TagConfigRule tag config rules. Take as input data from csv file. Where first column Arn
func TagConfigRule(csvData [][]string, client configserviceiface.ConfigServiceAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags []*configservice.Tag
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
		if awsErrorHandle(err) {
			return
		}
	}
}
