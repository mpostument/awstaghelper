package pkg

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
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
		log.Fatal("Not able to get IAM users ", err)
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
			fmt.Println("Not able to get IAM user tags ", err)
		}
		tags := map[string]string{}
		for _, tag := range userTags.Tags {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *user.UserName)
	}
	return rows
}

// TagIamUser tag iam user. Take as input data from csv file. Where first column is name
func TagIamUser(csvData [][]string, client iamiface.IAMAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags []*iam.Tag
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
		if awsErrorHandle(err) {
			return
		}
	}
}

func getIamRoles(client iamiface.IAMAPI) []*iam.Role {
	input := &iam.ListRolesInput{}

	var result []*iam.Role

	err := client.ListRolesPages(input,
		func(page *iam.ListRolesOutput, lastPage bool) bool {
			result = append(result, page.Roles...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get IAM roles ", err)
		return nil
	}
	return result
}

// ParseIamRolesTags parse output from getIamRoles and return roles and specified tags.
func ParseIamRolesTags(tagsToRead string, client iamiface.IAMAPI) [][]string {
	usersList := getIamRoles(client)
	rows := addHeadersToCsv(tagsToRead, "RoleName")
	for _, role := range usersList {
		roleTags, err := client.ListRoleTags(&iam.ListRoleTagsInput{RoleName: role.RoleName})
		if err != nil {
			fmt.Println("Not able to get iam roles tags", err)
		}
		tags := map[string]string{}
		for _, tag := range roleTags.Tags {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *role.RoleName)
	}
	return rows
}

// TagIamRole tag iam user. Take as input data from csv file. Where first column is name
func TagIamRole(csvData [][]string, client iamiface.IAMAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags []*iam.Tag
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &iam.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &iam.TagRoleInput{
			RoleName: aws.String(csvData[r][0]),
			Tags:     tags,
		}

		_, err := client.TagRole(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
