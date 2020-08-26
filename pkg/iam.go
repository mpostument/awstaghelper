package pkg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"log"
	"strings"
)

func getIamUsers(client iamiface.IAMAPI) []*iam.User {
	input := &iam.ListUsersInput{}

	var result []*iam.User

	err := client.ListUsersPages(input,
		func(page *iam.ListUsersOutput, lastPage bool) bool {
			result = append(result, page.Users...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get iam users", err)
		return nil
	}
	return result
}

// ParseIamUserTags parse output from getIamUsers and return username and specified tags.
func ParseIamUserTags(tagsToRead string, client iamiface.IAMAPI) [][]string {
	usersList := getIamUsers(client)
	var rows [][]string
	headers := []string{"UserName"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, user := range usersList {
		userTags, err := client.ListUserTags(&iam.ListUserTagsInput{UserName: user.UserName})
		if err != nil {
			fmt.Println("Not able to get iam user tags", err)
		}
		tags := map[string]string{}
		for _, tag := range userTags.Tags {
			tags[*tag.Key] = *tag.Value
		}

		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{*user.UserName}, resultTags...))
	}
	return rows
}

// TagIamUser tag iam user. Take as input data from csv file. Where first column is name
func TagIamUser(csvData [][]string, client iamiface.IAMAPI) {
	var tags []*iam.Tag
	for r := 1; r < len(csvData); r++ {
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &iam.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &iam.TagUserInput{
			UserName: aws.String(csvData[r][0]),
			Tags:     tags,
		}

		_, err := client.TagUser(input)
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
