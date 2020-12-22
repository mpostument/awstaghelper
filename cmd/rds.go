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

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
)

// rdsCmd represents the rds command
var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Root command for interaction with AWS rds services",
	Long:  `Root command for interaction with AWS rds services.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("rds called")
	//},
}

var getRdsCmd = &cobra.Command{
	Use:   "get-rds-tags",
	Short: "Write rds arn and required tags to csv",
	Long: `Write to csv data with rds arn and required tags to csv. 
This csv can be used with tag-rds command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		client := rds.New(sess)
		pkg.WriteCsv(pkg.ParseRDSTags(tags, client), filename)
	},
}

var tagRdsCmd = &cobra.Command{
	Use:   "tag-rds",
	Short: "Read csv and tag rds with csv data",
	Long:  `Read csv generated with get-rds-tags command and tag rds instances with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := pkg.GetSession(region, profile)
		csvData := pkg.ReadCsv(filename)
		client := rds.New(sess)
		pkg.TagRDS(csvData, client)
	},
}

func init() {
	rootCmd.AddCommand(rdsCmd)
	rdsCmd.AddCommand(getRdsCmd)
	rdsCmd.AddCommand(tagRdsCmd)
}
