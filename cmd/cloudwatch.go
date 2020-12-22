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

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

// cloudwatchCmd represents the cloudwatch command
var cloudwatchCmd = &cobra.Command{
	Use:   "cloudwatch",
	Short: "A root command for cloudwatch (alarms, logs, events)",
	Long:  `A root command for cloudwatch (alarms, logs, events).`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("cloudwatch called")
	//},
}

var getCloudWatchLogsCmd = &cobra.Command{
	Use:   "get-cwlog-tags",
	Short: "Write log group arn and required tags to csv",
	Long: `Write to csv data with log group arn and required tags to csv. 
This csv can be used with tag-log command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := cloudwatchlogs.New(sess)
		pkg.WriteCsv(pkg.ParseCwLogGroupTags(tags, client), filename)
	},
}

var tagCloudWatchLogsCmd = &cobra.Command{
	Use:   "tag-cwlogs",
	Short: "Read csv and tag cloudwatch logs with csv data",
	Long:  `Read csv generated with get-cwlog-tags command and tag cloudwatch logGroup with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := cloudwatchlogs.New(sess)
		csvData := pkg.ReadCsv(filename)
		pkg.TagCloudWatchLogGroups(csvData, client)
	},
}

var getCloudWatchAlarmsCmd = &cobra.Command{
	Use:   "get-cwalarm-tags",
	Short: "Write alarm arn and required tags to csv",
	Long: `Write to csv data with alarm arn and required tags to csv. 
This csv can be used with tag-alarm command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := cloudwatch.New(sess)
		pkg.WriteCsv(pkg.ParseCwAlarmTags(tags, client), filename)
	},
}

var tagCloudWatchAlarmsCmd = &cobra.Command{
	Use:   "tag-cwalarms",
	Short: "Read csv and tag cloudwatch alarms with csv data",
	Long:  `Read csv generated with get-cwalarms-tags command and tag cloudwatch alarms with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := cloudwatch.New(sess)
		csvData := pkg.ReadCsv(filename)
		pkg.TagCloudWatchAlarm(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(cloudwatchCmd)
	cloudwatchCmd.AddCommand(getCloudWatchLogsCmd)
	cloudwatchCmd.AddCommand(tagCloudWatchLogsCmd)
	cloudwatchCmd.AddCommand(getCloudWatchAlarmsCmd)
	cloudwatchCmd.AddCommand(tagCloudWatchAlarmsCmd)
}
