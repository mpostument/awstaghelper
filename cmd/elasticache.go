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
	"awstaghelper/modules/elastiCacheHelper"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// elasticacheCmd represents the elasticache command
var elasticacheCmd = &cobra.Command{
	Use:   "elasticache",
	Short: "Root command for interaction with AWS elasticache services",
	Long:  `Root command for interaction with AWS elasticache services.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("elasticache called")
	//},
}

var getElastiCacheCmd = &cobra.Command{
	Use:   "get-elasticache-tags",
	Short: "Write elasticache arn and required tags to csv",
	Long: `Write to csv data with elasticache arn and required tags to csv. 
This csv can be used with tag-elasticache command to tag aws environment.
Specify list of tags which should be read using tags flag: --tags Name,Env,Project.
Csv filename can be specified with flag filename.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := common.GetSession(region, profile)
		elastiCacheClient := elasticache.New(sess)
		stsClient := sts.New(sess)
		common.WriteCsv(elastiCacheHelper.ParseElastiCacheTags(tags, elastiCacheClient, stsClient, *sess.Config.Region), filename)
	},
}

var tagElastiCacheCmd = &cobra.Command{
	Use:   "tag-elasticache",
	Short: "Read csv and tag elasticache with csv data",
	Long:  `Read csv generated with get-elasticache-tags command and tag elasticache instances with tags from csv.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		sess := common.GetSession(region, profile)
		elastiCacheClient := elasticache.New(sess)
		csvData := common.ReadCsv(filename)
		elastiCacheHelper.TagElasticache(csvData, elastiCacheClient)
	},
}

func init() {
	rootCmd.AddCommand(elasticacheCmd)
	elasticacheCmd.AddCommand(getElastiCacheCmd)
	elasticacheCmd.AddCommand(tagElastiCacheCmd)
}
