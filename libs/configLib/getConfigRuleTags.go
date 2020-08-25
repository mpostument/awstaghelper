package configLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
	"log"
	"strings"
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
