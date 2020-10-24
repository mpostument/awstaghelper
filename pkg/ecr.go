package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

func getEcrRepositories(client ecriface.ECRAPI) []*ecr.Repository {
	input := &ecr.DescribeRepositoriesInput{}

	var result []*ecr.Repository

	err := client.DescribeRepositoriesPages(input,
		func(page *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
			result = append(result, page.Repositories...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get ecr repositories ", err)
		return nil
	}
	return result
}

// ParseEcrRepositoriesTags parse output from getEcrRepositories and return repo arn and specified tags.
func ParseEcrRepositoriesTags(tagsToRead string, client ecriface.ECRAPI) [][]string {
	repoList := getEcrRepositories(client)
	rows := addHeadersToCsv(tagsToRead, "Arn")
	for _, repo := range repoList {
		repoTags, err := client.ListTagsForResource(&ecr.ListTagsForResourceInput{ResourceArn: repo.RepositoryArn})
		if err != nil {
			fmt.Println("Not able to get ecr tags ", err)
		}
		tags := map[string]string{}
		for _, tag := range repoTags.Tags {
			tags[*tag.Key] = *tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, *repo.RepositoryArn)
	}
	return rows
}

// TagEcrRepo tag ecr repo. Take as input data from csv file. Where first column is name
func TagEcrRepo(csvData [][]string, client ecriface.ECRAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags []*ecr.Tag
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &ecr.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &ecr.TagResourceInput{
			ResourceArn: aws.String(csvData[r][0]),
			Tags:        tags,
		}

		_, err := client.TagResource(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
