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
	"awstaghelper/modules/common"
	"awstaghelper/modules/ec2Helper"
	"github.com/spf13/cobra"
)

// ec2Cmd represents the ec2 command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Root command for interaction with AWS ec2 services",
	Long:  `Root command for interaction with AWS ec2 services.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("ec2 called")
	//},
}

var getEc2Cmd = &cobra.Command{
	Use:   "get-ec2-tags",
	Short: "Write ec2 id and required tags to csv",
	Long: `Write to csv data with ec2 id and required tags to csv. 
This csv can be used with tag-ec2 command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename. Default ec2Tags.csv`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := common.GetSession(region, profile)
		common.WriteCsv(ec2Helper.ParseEc2Tags(tags, *sess), filename)
	},
}

var tagEc2Cmd = &cobra.Command{
	Use:   "tag-ec2",
	Short: "Read csv and tag ec2 with csv data",
	Long:  `Read csv generated with get-ec2-tags command and tag ec2 instances with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := common.GetSession(region, profile)
		csvData := common.ReadCsv(filename)
		ec2Helper.TagEc2(csvData, *sess)
	},
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
	ec2Cmd.AddCommand(getEc2Cmd)
	ec2Cmd.AddCommand(tagEc2Cmd)
}
