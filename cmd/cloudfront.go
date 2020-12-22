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
	"github.com/mpostument/awstaghelper/pkg"

	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/spf13/cobra"
)

// cloudfrontCmd represents the cloudfront command
var cloudfrontCmd = &cobra.Command{
	Use:   "cloudfront",
	Short: "A root command for cloudfront",
	Long:  `A root command for cloudfront.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("cloudfront called")
	//},
}

var getDistributionsCmd = &cobra.Command{
	Use:   "get-distribution-tags",
	Short: "Write distribution arn and required tags to csv",
	Long: `Write to csv data with distribution arn and required tags to csv. 
This csv can be used with tag-distribution command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := cloudfront.New(sess)
		pkg.WriteCsv(pkg.ParseDistributionsTags(tags, client), filename)
	},
}

var tagDistributionsCmd = &cobra.Command{
	Use:   "tag-distribution",
	Short: "Read csv and tag distribution with csv data",
	Long:  `Read csv generated with get-distribution-tags command and tag distribution with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := cloudfront.New(sess)
		csvData := pkg.ReadCsv(filename)
		pkg.TagDistribution(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(cloudfrontCmd)
	cloudfrontCmd.AddCommand(getDistributionsCmd)
	cloudfrontCmd.AddCommand(tagDistributionsCmd)
}
