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

	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/spf13/cobra"
)

// kinesisCmd represents the kinesis command
var kinesisCmd = &cobra.Command{
	Use:   "kinesis",
	Short: "Root command for interaction with AWS kinesis services",
	Long:  `Root command for interaction with AWS kinesis services.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("kinesis called")
	//},
}

var getStreamCmd = &cobra.Command{
	Use:   "get-stream-tags",
	Short: "Write kinesis stream name and required tags to csv",
	Long: `Write to csv data with kinesis names and required tags to csv. 
This csv can be used with tag-stream command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := kinesis.New(sess)
		pkg.WriteCsv(pkg.ParseKinesisTags(tags, client), filename)
	},
}

var tagStreamCmd = &cobra.Command{
	Use:   "tag-stream",
	Short: "Read csv and tag kinesis stream with csv data",
	Long:  `Read csv generated with get-stream-tags command and tag kinesis stream with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := kinesis.New(sess)
		csvData := pkg.ReadCsv(filename)
		pkg.TagKinesisStream(csvData, client)
	},
}

var getFirehoseCmd = &cobra.Command{
	Use:   "get-firehose-tags",
	Short: "Write firehose stream name and required tags to csv",
	Long: `Write to csv data with firehose names and required tags to csv. 
This csv can be used with tag-stream command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := firehose.New(sess)
		pkg.WriteCsv(pkg.ParseFirehoseTags(tags, client), filename)
	},
}

var tagFirehoseCmd = &cobra.Command{
	Use:   "tag-firehose",
	Short: "Read csv and tag firehose stream with csv data",
	Long:  `Read csv generated with get-firehose-tags command and tag firehose stream with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := firehose.New(sess)
		csvData := pkg.ReadCsv(filename)
		pkg.TagFirehose(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(kinesisCmd)
	kinesisCmd.AddCommand(getStreamCmd)
	kinesisCmd.AddCommand(tagStreamCmd)
	kinesisCmd.AddCommand(getFirehoseCmd)
	kinesisCmd.AddCommand(tagFirehoseCmd)
}
