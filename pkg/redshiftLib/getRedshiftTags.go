package redshiftLib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/redshift/redshiftiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"log"
	"strings"
)

// getInstances return all redshift instances from specified region
func getInstances(client redshiftiface.RedshiftAPI) []*redshift.Cluster {
	input := &redshift.DescribeClustersInput{}

	var result []*redshift.Cluster

	err := client.DescribeClustersPages(input,
		func(page *redshift.DescribeClustersOutput, lastPage bool) bool {
			result = append(result, page.Clusters...)
			return !lastPage
		})
	if err != nil {
		log.Fatal("Not able to get instances", err)
		return nil
	}
	return result
}

// ParseRedshiftTags parse output from getInstances and return arn and specified tags.
func ParseRedshiftTags(tagsToRead string, client redshiftiface.RedshiftAPI, stsClient stsiface.STSAPI, region string) [][]string {
	instancesOutput := getInstances(client)

	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal("Not able to get account id", err)
	}

	var rows [][]string
	headers := []string{"Arn"}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	for _, redshiftInstances := range instancesOutput {
		clusterArn := fmt.Sprintf("arn:aws:redshift:%s:%s:cluster:%s",
			region, *callerIdentity.Account, *redshiftInstances.ClusterIdentifier)
		redshiftTags, err := client.DescribeTags(&redshift.DescribeTagsInput{ResourceName: &clusterArn})
		if err != nil {
			fmt.Println("Not able to get redshift tags", err)
		}
		tags := map[string]string{}
		for _, tag := range redshiftTags.TaggedResources {
			tags[*tag.Tag.Key] = *tag.Tag.Value
		}
		var resultTags []string
		for _, key := range strings.Split(tagsToRead, ",") {
			resultTags = append(resultTags, tags[key])
		}
		rows = append(rows, append([]string{clusterArn}, resultTags...))
	}
	return rows
}
