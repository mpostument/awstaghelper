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
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/mpostument/awstaghelper/pkg"

	"github.com/spf13/cobra"
)

// ebCmd represents the ecr command
var ebCmd = &cobra.Command{
	Use:   "eb",
	Short: "Root command for interaction with AWS elastic bean stalk services",
	Long:  `Root command for interaction with AWS elastic bean stalk services.`,
}

var getEbTagsCmd = &cobra.Command{
	Use:   "get-eb-tags",
	Short: "Write arn and required tags to csv",
	Long: `Write to csv data with arn and required tags to csv. 
This csv can be used with tag-eb command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := elasticbeanstalk.New(sess)
		pkg.WriteCsv(pkg.ParseEBTags(tags, client), filename)
	},
}

var tagEbCmd = &cobra.Command{
	Use:   "tag-eb",
	Short: "Read csv and tag eb e with csv data",
	Long:  `Read csv generated with get-eb-tags command and tag elastic bean stalk with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := elasticbeanstalk.New(sess)
		csvData := pkg.ReadCsv(filename)
		pkg.TagEbEnvironments(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(ebCmd)
	ebCmd.AddCommand(getEbTagsCmd)
	ebCmd.AddCommand(tagEbCmd)
}
