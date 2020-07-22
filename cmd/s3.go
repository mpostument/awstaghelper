/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"awstaghelper/modules/common"
	"awstaghelper/modules/s3Helper"
	"github.com/spf13/cobra"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Root command for interaction with AWS s3 services",
	Long:  `Root command for interaction with AWS s3 services.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("s3 called")
	//},
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
		sess := common.GetSession(region, profile)
		common.WriteCsv(s3Helper.ParseS3Tags(tags, *sess), filename)
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
		sess := common.GetSession(region, profile)
		csvData := common.ReadCsv(filename)
		s3Helper.TagS3(csvData, *sess)
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)
	s3Cmd.AddCommand(getS3Cmd)
	s3Cmd.AddCommand(tagS3Cmd)
}
