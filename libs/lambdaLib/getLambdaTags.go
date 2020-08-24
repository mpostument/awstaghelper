package lambdaLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"log"
	"strings"
)

// getLambdas return all lambdas from specified region
func getLambdas(client lambdaiface.LambdaAPI) []*lambda.FunctionConfiguration {
	input := &lambda.ListFunctionsInput{}

	var result []*lambda.FunctionConfiguration

	err := client.ListFunctionsPages(input,
		func(page *lambda.ListFunctionsOutput, lastPage bool) bool {
			result = append(result, page.Functions...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get lambdas", err)
		return nil
	}
	return result
}

// ParseLambdasTags parse output from getLambdas and return arn and specified tags.
func ParseLambdasTags(tagsToRead string, client lambdaiface.LambdaAPI) [][]string {
	instancesOutput := getLambdas(client)
	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, lambdaOutput := range instancesOutput {
		lambdaTags, err := client.ListTags(&lambda.ListTagsInput{Resource: lambdaOutput.FunctionArn})
		if err != nil {
			fmt.Println("Not able to get lambda tags", err)
		}
		tags := map[string]string{}
		for key, value := range lambdaTags.Tags {
			tags[key] = *value
		}

		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*lambdaOutput.FunctionArn}, resultTags...))
	}
	return rows
}
