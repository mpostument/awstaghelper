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

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var getSecurityGroupsCmd = &cobra.Command{
	Use:   "get-sg-tags",
	Short: "Write sg id and required tags to csv",
	Long: `Write to csv data with sg id and required tags to csv. 
This csv can be used with tag-sg command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := ec2.New(sess)
		pkg.WriteCsv(pkg.ParseSecurityGroupTags(tags, client), filename)
	},
}

var tagSecurityGroupsCmd = &cobra.Command{
	Use:   "tag-sg",
	Short: "Read csv and tag sg with csv data",
	Long:  `Read csv generated with get-sg-tags command and tag sg resources with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		csvData := pkg.ReadCsv(filename)
		client := ec2.New(sess)
		pkg.TagSecurityGroups(csvData, client)
	},
}

func init() {
	ec2Cmd.AddCommand(getSecurityGroupsCmd)
	ec2Cmd.AddCommand(tagSecurityGroupsCmd)
}
