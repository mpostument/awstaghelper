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

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/cobra"
)

// iamCmd represents the iam command
var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Root command for interaction with AWS iam services",
	Long:  `Root command for interaction with AWS iam services.`,
}

var getIamUsersCmd = &cobra.Command{
	Use:   "get-user-tags",
	Short: "Write user name and required tags to csv",
	Long: `Write to csv data with user name and required tags to csv. 
This csv can be used with tag-user command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := iam.New(sess)
		pkg.WriteCsv(pkg.ParseIamUserTags(tags, client), filename)
	},
}

var tagIamUserCmd = &cobra.Command{
	Use:   "tag-user",
	Short: "Read csv and tag iam user with csv data",
	Long:  `Read csv generated with get-user-tags command and tag iam user with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := iam.New(sess)
		csvData := pkg.ReadCsv(filename)
		pkg.TagIamUser(csvData, client)
	},
}

var getIamRolesCmd = &cobra.Command{
	Use:   "get-role-tags",
	Short: "Write role name and required tags to csv",
	Long: `Write to csv data with role name and required tags to csv. 
This csv can be used with tag-role command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := iam.New(sess)
		pkg.WriteCsv(pkg.ParseIamRolesTags(tags, client), filename)
	},
}

var tagIamRoleCmd = &cobra.Command{
	Use:   "tag-role",
	Short: "Read csv and tag iam role with csv data",
	Long:  `Read csv generated with get-roles-tags command and tag iam role with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := iam.New(sess)
		csvData := pkg.ReadCsv(filename)
		pkg.TagIamRole(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(iamCmd)
	iamCmd.AddCommand(getIamUsersCmd)
	iamCmd.AddCommand(tagIamUserCmd)
	iamCmd.AddCommand(getIamRolesCmd)
	iamCmd.AddCommand(tagIamRoleCmd)
}
