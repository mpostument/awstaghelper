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
package cmd

import (
	"awstaghelper/libs/cloudWatchLib"
	"awstaghelper/libs/commonLib"
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
		sess := commonLib.GetSession(region, profile)
		client := cloudwatchlogs.New(sess)
		commonLib.WriteCsv(cloudWatchLib.ParseCwLogGroupTags(tags, client), filename)
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
		sess := commonLib.GetSession(region, profile)
		client := cloudwatchlogs.New(sess)
		csvData := commonLib.ReadCsv(filename)
		cloudWatchLib.TagCloudWatchLogGroups(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(cloudwatchCmd)
	cloudwatchCmd.AddCommand(getCloudWatchLogsCmd)
	cloudwatchCmd.AddCommand(tagCloudWatchLogsCmd)
}
