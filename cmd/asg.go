/*
Copyright © 2023 Cristian Magherusan-Stanciu cristi@leanercloud.com
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
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/mpostument/awstaghelper/pkg"

	"github.com/spf13/cobra"
)

// asgCmd represents the asg command
var asgCmd = &cobra.Command{
	Use:   "asg",
	Short: "Root command for interaction with AWS autoscaling groups",
	Long:  `Root command for interaction with AWS autoscaling groups.`,
}

var getASGCmd = &cobra.Command{
	Use:   "get-asg-tags",
	Short: "Write ASG names and required tags to CSV",
	Long: `Write to csv data with ASG name and required tags to CSV.
This CSV can be used with tag-asg command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := autoscaling.New(sess)
		pkg.WriteCsv(pkg.ParseASGTags(tags, client), filename)
	},
}

var tagASGCmd = &cobra.Command{
	Use:   "tag-asg",
	Short: "Read CSV and tag ASGs with CSV data",
	Long:  `Read CSV generated with get-asg-tags command and tag ASGs with tags from CSV.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		csvData := pkg.ReadCsv(filename)
		client := autoscaling.New(sess)
		pkg.TagASG(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(asgCmd)
	asgCmd.AddCommand(getASGCmd)
	asgCmd.AddCommand(tagASGCmd)
}
