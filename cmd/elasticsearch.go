/*
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

// Package cmd is the package for the CLI of awstaghelper
package cmd

import (
	"fmt"

	"github.com/mpostument/awstaghelper/pkg"

	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/spf13/cobra"
)

// elasticsearchCmd represents the elasticsearch command
var elasticsearchCmd = &cobra.Command{
	Use:   "elasticsearch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("elasticsearch called")
	},
}

var getElasticSearchCmd = &cobra.Command{
	Use:   "get-elasticsearch-tags",
	Short: "Write elasticsearch arn and required tags to csv",
	Long: `Write to csv data with elasticsearch arn and required tags to csv. 
This csv can be used with tag-elasticsearch command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := elasticsearchservice.New(sess)
		stsClient := sts.New(sess)
		pkg.WriteCsv(pkg.ParseElasticSearchTags(tags, client, stsClient, region), filename)
	},
}

var tagElasticSearchCmd = &cobra.Command{
	Use:   "tag-elasticsearch",
	Short: "Read csv and tag elasticsearch with csv data",
	Long:  `Read csv generated with get-elasticsearch-tags command and tag elasticsearch instances with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		csvData := pkg.ReadCsv(filename)
		client := elasticsearchservice.New(sess)
		pkg.TagElasticSearch(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(elasticsearchCmd)
	elasticsearchCmd.AddCommand(getElasticSearchCmd)
	elasticsearchCmd.AddCommand(tagElasticSearchCmd)
}
