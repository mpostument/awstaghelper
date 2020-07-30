package lambdaHelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

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
