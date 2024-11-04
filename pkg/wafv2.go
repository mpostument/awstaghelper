/*
Copyright © 2024 Jaemok Hong jaemokhong@lguplus.co.kr
Copyright © 2020 Maksym Postument 777rip777@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/aws/aws-sdk-go/service/wafv2/wafv2iface"
)

// getWebACL return all webacls(regional, cloudfront) from specified region
func getWebACL(scope string, client wafv2iface.WAFV2API) *wafv2.ListWebACLsOutput {
	input := &wafv2.ListWebACLsInput{
		Scope: aws.String(scope),
	}

	var allWebACLs []*wafv2.WebACLSummary
	var nextMarker *string

	for {
		result, err := client.ListWebACLs(input)
		if err != nil {
			log.Fatal("Not able to get webacls ", err)
			return nil
		}
		allWebACLs = append(allWebACLs, result.WebACLs...)
		nextMarker = result.NextMarker
		if nextMarker == nil || *nextMarker == "" {
			break
		}
		input.NextMarker = nextMarker
	}
	return &wafv2.ListWebACLsOutput{WebACLs: allWebACLs}
}

// ParseWebACLTags parse output from getWebACL and return webacl arn and specified tags.
func ParseWebACLTags(tagsToRead string, scope string, client wafv2iface.WAFV2API) [][]string {
	wafv2Output := getWebACL(scope, client)
	rows := addHeadersToCsv(tagsToRead, "Arn")

	for _, webACL := range wafv2Output.WebACLs {
		var nextMarker *string

		for {
			input := &wafv2.ListTagsForResourceInput{
				ResourceARN: webACL.ARN,
				Limit:       aws.Int64(5),
			}
			if nextMarker != nil {
				input.NextMarker = nextMarker
			}
			webACLTags, err := client.ListTagsForResource(input)
			if err != nil {
				fmt.Println("Not able to get webACL tags ", err)
				break
			}
			tags := map[string]string{}
			for _, tag := range webACLTags.TagInfoForResource.TagList {
				tags[*tag.Key] = *tag.Value
			}
			rows = addTagsToCsv(tagsToRead, tags, rows, *webACL.ARN)
			nextMarker = webACLTags.NextMarker
			if nextMarker == nil || *nextMarker == "" {
				break
			}
		}
	}
	return rows
}

// TagWebACL tag webacl. Take as input data from csv file. Where first column Arn
func TagWebACL(csvData [][]string, client wafv2iface.WAFV2API) {
	for r := 1; r < len(csvData); r++ {
		var tags []*wafv2.Tag
		for c := 1; c < len(csvData[0]); c++ {
			tags = append(tags, &wafv2.Tag{
				Key:   &csvData[0][c],
				Value: &csvData[r][c],
			})
		}

		input := &wafv2.TagResourceInput{
			ResourceARN: aws.String(csvData[r][0]),
			Tags:        tags,
		}

		_, err := client.TagResource(input)
		if awsErrorHandle(err) {
			return
		}
	}
}
