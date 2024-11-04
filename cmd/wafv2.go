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

// Package cmd is the package for the CLI of awstaghelper
package cmd

import (
	"github.com/mpostument/awstaghelper/pkg"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/spf13/cobra"
)

// wafv2Cmd represents the wafv2 command
var wafv2Cmd = &cobra.Command{
	Use:   "wafv2",
	Short: "Root command for interaction with AWS wafv2 services(cloudfront, regional)",
	Long:  `Root command for interaction with AWS wafv2 services(cloudfront, regional),`,
}

var getWebACLCmd = &cobra.Command{
	Use:   "get-webacl-tags",
	Short: "Write regional webacl arn and required tags to csv",
	Long: `Write to csv data with  regional webacl arn and required tags to csv.
This csv can be used with tag-webacl command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := wafv2.New(sess)
		pkg.WriteCsv(pkg.ParseWebACLTags(tags, "REGIONAL", client), filename)
	},
}

var tagWebACLCmd = &cobra.Command{
	Use:   "tag-webacl",
	Short: "Read csv and tag regional webacl with csv data",
	Long:  `Read csv generated with get-webacl-tags command and tag regional webacl with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := wafv2.New(sess, aws.NewConfig().WithRegion(region))
		csvData := pkg.ReadCsv(filename)
		pkg.TagWebACL(csvData, client)
	},
}

var getCFWebACLCmd = &cobra.Command{
	Use:   "get-cfwebacl-tags",
	Short: "Write cloudfront webacl arn and required tags to csv",
	Long: `Write to csv data with cloudfront webacl arn and required tags to csv.
This csv can be used with tag-cfwebacl command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		sess := pkg.GetSession("us-east-1", profile)
		client := wafv2.New(sess)
		pkg.WriteCsv(pkg.ParseWebACLTags(tags, "CLOUDFRONT", client), filename)
	},
}

var tagCFWebACLCmd = &cobra.Command{
	Use:   "tag-cfwebacl",
	Short: "Read csv and tag cloudfront webacl with csv data",
	Long:  `Read csv generated with get-cfwebacl-tags command and tag cloudfront webacl with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := wafv2.New(sess, aws.NewConfig().WithRegion("us-east-1"))
		csvData := pkg.ReadCsv(filename)
		pkg.TagWebACL(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(wafv2Cmd)
	wafv2Cmd.AddCommand(getWebACLCmd)
	wafv2Cmd.AddCommand(tagWebACLCmd)
	wafv2Cmd.AddCommand(getCFWebACLCmd)
	wafv2Cmd.AddCommand(tagCFWebACLCmd)
}
