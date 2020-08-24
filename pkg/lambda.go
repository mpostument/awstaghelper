package pkg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"log"
	"strings"
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
		log.Fatal("Not able to get instances", err)
		return nil
	}
	return result
}

// ParseLambdaFunctionTags parse output from getLambdaFunctions and return arn and specified tags.
func ParseLambdaFunctionTags(tagsToRead string, client lambdaiface.LambdaAPI) [][]string {
	instancesOutput := getLambdaFunctions(client)
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

// TagLambda tag instances. Take as input data from csv file. Where first column id
func TagLambda(csvData [][]string, client lambdaiface.LambdaAPI) {

	tags := make(map[string]*string)
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags[csvData[0][c]] = &csvData[r][c]
		}

		input := &lambda.TagResourceInput{
			Resource: aws.String(csvData[r][0]),
			Tags:     tags,
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
