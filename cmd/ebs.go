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
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mpostument/awstaghelper/pkg"

	"github.com/spf13/cobra"
)

// ebsCmd represents the ebs command
var ebsCmd = &cobra.Command{
	Use:   "ebs",
	Short: "Root command for interaction with AWS EBS volumes",
	Long:  `Root command for interaction with AWS EBS volumes,`,
}

var getEBSVolumeCmd = &cobra.Command{
	Use:   "get-ebs-tags",
	Short: "Write EBS volume IDs and required tags to CSV",
	Long:  `Write to csv data with EBS volume IDs and required tags to CSV. This CSV can be used with tag-ebs command to tag AWS environment. Specify list of tags which should be read using tags flag: --tags Name,Env,Project. Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := ec2.New(sess)
		pkg.WriteCsv(pkg.ParseEBSVolumeTags(tags, client), filename)
	},
}

var tagEBSVolumeCmd = &cobra.Command{
	Use:   "tag-ebs",
	Short: "Read CSV and tag EBS volumes with CSV data",
	Long:  `Read CSV generated with get-ebs-tags command and tag EBS volumes with tags from CSV.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		csvData := pkg.ReadCsv(filename)
		client := ec2.New(sess)
		pkg.TagEBSVolumes(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(ebsCmd)
	ebsCmd.AddCommand(getEBSVolumeCmd)
	ebsCmd.AddCommand(tagEBSVolumeCmd)
}
