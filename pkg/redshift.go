package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/redshift/redshiftiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// getRedshiftInstances return all redshift instances from specified region
func getRedshiftInstances(client redshiftiface.RedshiftAPI) []*redshift.Cluster {
	input := &redshift.DescribeClustersInput{}

	var result []*redshift.Cluster

	err := client.DescribeClustersPages(input,
		func(page *redshift.DescribeClustersOutput, lastPage bool) bool {
			result = append(result, page.Clusters...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get redshift instances ", err)
		return nil
	}
	return result
}

// ParseRedshiftTags parse output from getRedshiftInstances and return arn and specified tags.
func ParseRedshiftTags(tagsToRead string, client redshiftiface.RedshiftAPI, stsClient stsiface.STSAPI, region string) [][]string {
	instancesOutput := getRedshiftInstances(client)

	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal("Not able to get account id", err)
	}
	rows := addHeadersToCsv(tagsToRead, "Arn")
	for _, redshiftInstances := range instancesOutput {
		clusterArn := fmt.Sprintf("arn:aws:redshift:%s:%s:cluster:%s",
			region, *callerIdentity.Account, *redshiftInstances.ClusterIdentifier)
		redshiftTags, err := client.DescribeTags(&redshift.DescribeTagsInput{ResourceName: &clusterArn})
		if err != nil {
			fmt.Println("Not able to get redshift tags ", err)
		}
		tags := map[string]string{}
		for _, tag := range redshiftTags.TaggedResources {
			tags[*tag.Tag.Key] = *tag.Tag.Value
		}
		rows = addTagsToCsv(tagsToRead, tags, rows, clusterArn)
	}
	return rows
}

// TagRedShift tag rds instances. Take as input data from csv file. Where first column arn
func TagRedShift(csvData [][]string, client redshiftiface.RedshiftAPI) {
	for r := 1; r < len(csvData); r++ {
		var tags []*redshift.Tag
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &redshift.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &redshift.CreateTagsInput{
			ResourceName: aws.String(csvData[r][0]),
			Tags:         tags,
		}

		_, err := client.CreateTags(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
