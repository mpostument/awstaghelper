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

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Root command for interaction with AWS s3 services",
	Long:  `Root command for interaction with AWS s3 services.`,
}

var getS3Cmd = &cobra.Command{
	Use:   "get-s3-tags",
	Short: "Write rds arn and required tags to csv",
	Long: `Write to csv data with rds arn and required tags to csv. 
This csv can be used with tag-s3 command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := s3.New(sess)
		pkg.WriteCsv(pkg.ParseS3Tags(tags, client), filename)
	},
}

var tagS3Cmd = &cobra.Command{
	Use:   "tag-s3",
	Short: "Read csv and tag s3 with csv data",
	Long:  `Read csv generated with get-s3-tags command and tag s3 instances with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := s3.New(sess)
		csvData := pkg.ReadCsv(filename)
		pkg.TagS3(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)
	s3Cmd.AddCommand(getS3Cmd)
	s3Cmd.AddCommand(tagS3Cmd)
}
