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

	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// redshiftCmd represents the redshift command
var redshiftCmd = &cobra.Command{
	Use:   "redshift",
	Short: "Root command for interaction with AWS redshfit services",
	Long:  `Root command for interaction with AWS redshfit services.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("redshift called")
	//},
}

var getRedshiftCmd = &cobra.Command{
	Use:   "get-redshift-tags",
	Short: "Write redshift id and required tags to csv",
	Long: `Write to csv data with redshift id and required tags to csv. 
This csv can be used with tag-redshift command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := redshift.New(sess)
		stsClient := sts.New(sess)
		pkg.WriteCsv(pkg.ParseRedshiftTags(tags, client, stsClient, region), filename)
	},
}

var tagRedshifCmd = &cobra.Command{
	Use:   "tag-redshift",
	Short: "Read csv and tag redshift with csv data",
	Long:  `Read csv generated with get-redshift-tags command and tag redshift instances with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := redshift.New(sess)
		csvData := pkg.ReadCsv(filename)
		pkg.TagRedShift(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(redshiftCmd)
	redshiftCmd.AddCommand(getRedshiftCmd)
	redshiftCmd.AddCommand(tagRedshifCmd)
}
