package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

// getLambdaFunctions return all lambdas from specified region
func getLambdaFunctions(client lambdaiface.LambdaAPI) []*lambda.FunctionConfiguration {
	input := &lambda.ListFunctionsInput{}

	var result []*lambda.FunctionConfiguration

	err := client.ListFunctionsPages(input,
		func(page *lambda.ListFunctionsOutput, lastPage bool) bool {
			result = append(result, page.Functions...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get lambdas ", err)
		return nil
	}
	return result
}

// ParseLambdaFunctionTags parse output from getLambdaFunctions and return arn and specified tags.
func ParseLambdaFunctionTags(tagsToRead string, client lambdaiface.LambdaAPI) [][]string {
	instancesOutput := getLambdaFunctions(client)
	rows := addHeadersToCsv(tagsToRead, "Arn")
	for _, lambdaOutput := range instancesOutput {
		lambdaTags, err := client.ListTags(&lambda.ListTagsInput{Resource: lambdaOutput.FunctionArn})
		if err != nil {
			fmt.Println("Not able to get lambda tags ", err)
		}
		tags := map[string]string{}
		for key, value := range lambdaTags.Tags {
			tags[key] = *value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *lambdaOutput.FunctionArn)
	}
	return rows
}

// TagLambda tag instances. Take as input data from csv file. Where first column id
func TagLambda(csvData [][]string, client lambdaiface.LambdaAPI) {
	for r := 1; r < len(csvData); r++ {
		tags := make(map[string]*string)
		for c := 1; c < len(csvData[0]); c++ {
			tags[csvData[0][c]] = &csvData[r][c]
		}

		input := &lambda.TagResourceInput{
			Resource: aws.String(csvData[r][0]),
			Tags:     tags,
		}

		_, err := client.TagResource(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
