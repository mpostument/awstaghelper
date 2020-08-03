# Awstaghelper

### Painless AWS tagging

awstaghelper allow tagging hundreds of AWS resources in few commands

## Table of Contents

* [Usage Example](#usage-example)
* [Installing](#installing)
* [Getting Started](#getting-started)
  * [Ec2](#ec2)
  * [Rds](#rds)
  * [Lambda](#lambda)
  * [S3](#s3)
  * [ElastiCache](#elasticache)
  * [Redshift](#redshift)
  * [ElasticSearch](#elasticsearch)
  * [CloudWatch Logs](#cloudwatchlogs)
* [Global parameters](#global-parameters)
* [Contributing](#contributing)
* [License](#license)

## Usage example

![Usage](awstaghelper.gif)

## Installing

Download the latest binary from [releases](https://github.com/mpostument/awstaghelper/releases)

## Getting Started

### Ec2

#### Get ec2 tags

Get list of ec2 with required tags - `awstaghelper ec2 get-ec2-tags`  
Example:
 `awstaghelper ec2 get-ec2-tags --filename ec2Tag.csv --tags Name,Owner --profile main`

#### Tag ec2

Read csv and tag ec2 - `awstaghelper ec2 tag-ec2`  
Example:
 `awstaghelper ec2 tag-ec2 --filename ec2Tag.csv --profile main`  

### Rds

#### Get rds tags

Get list of rds with required tags - `awstaghelper rds get-rds-tags`  
Example:
 `awstaghelper rds get-rds-tags --filename rdsTag.csv --tags Name,Owner --profile main`

#### Tag rds

Read csv and tag rds - `awstaghelper rds tag-rds`  
Example:
 `awstaghelper rds tag-rds --filename rdsTag.csv --profile main`  

### Lambda

#### Get lambda tags

Get list of lambda with required tags - `awstaghelper lambda get-lambda-tags`  
Example:
 `awstaghelper lambda get-lambda-tags --filename lambdaTag.csv --tags Name,Owner --profile main`

#### Tag lambda

Read csv and tag lambda - `awstaghelper lambda tag-lambda`  
Example:
 `awstaghelper lambda tag-lambda --filename lambdaTag.csv --profile main`  

### S3

#### Get s3 tags

Get list of s3 with required tags - `awstaghelper s3 get-s3-tags`  
Example:
 `awstaghelper s3 get-s3-tags --filename s3Tag.csv --tags Name,Owner --profile main`

#### Tag s3

Read csv and tag s3 - `awstaghelper s3 tag-s3`  
Example:
 `awstaghelper s3 tag-s3 --filename ec2Tag.csv --profile main`  
 
### Elasticache

#### Get elasticache tags

Get list of s3 with required tags - `awstaghelper elasticache get-elasticache-tags`  
Example:
`awstaghelper elasticache get-elasticache-tags --filename elasticacheTag.csv --tags Name,Owner --profile main`

#### Tag elasticache

Read csv and tag s3 - `awstaghelper elasticache tag-elasticache`  
Example:
`awstaghelper elasticache tag-elasticache --filename elasticacheTag.csv --profile main`  

### Redshift

#### Get redshift tags

Get list of redshift with required tags - `awstaghelper redshift get-redshift-tags`  
Example:
`awstaghelper redshift get-redshift-tags --filename redshiftTag.csv --tags Name,Owner --profile main`

#### Tag redshift

Read csv and tag redshift - `awstaghelper redshift tag-redshift`  
Example:
`awstaghelper redshift tag-redshift --filename redshiftTag.csv --profile main`  

### ElasticSearch

#### Get elasticsearch tags

Get list of elasticsearch with required tags - `awstaghelper elasticsearch get-elasticsearch-tags`  
Example:
`awstaghelper elasticsearch get-elasticsearch-tags --filename elasticsearchTag.csv --tags Name,Owner --profile main`

#### Tag elasticsearch

Read csv and tag elasticsearch - `awstaghelper elasticsearch tag-elasticsearch`  
Example:
`awstaghelper elasticsearch tag-elasticsearch --filename elasticsearchTag.csv --profile main`  

### CloudWatch Logs

#### Get cloudwatch logs tags

Get list of cloudwatch logs with required tags - `awstaghelper cloudwatch get-cwlogs-tags`  
Example:
`awstaghelper cloudwatch get-cwlogs-tags --filename cwlTag.csv --tags Name,Owner --profile main`

#### Tag cloudwatch logs

Read csv and tag cloudwatch logs - `awstaghelper cloudwatch tag-cwlogs`  
Example:
`awstaghelper cloudwatch tag-cwlogs --filename cwlTag.csv --profile main`  

## Global parameters

`filename` - path where to write or read data. Supported by every option. Default `awsTags.csv`  
`tags` - list of tags to read. Supported only with `get-tags` option. Default `Name,Environment`  
`region` - aws region to use. Supported by every option. Default `us-east-1`  
`profile` - aws profile to use. Supported by every option. Default `default`  

## Contributing

1. Fork it
2. Download your fork to your PC ( `git clone https://github.com/your_username/awstaghelper && cd awstaghelper` )
3. Create your feature branch ( `git checkout -b my-new-feature` )
4. Make changes and add them ( `git add .` )
5. Commit your changes ( `git commit -m 'Add some feature'` )
6. Push to the branch ( `git push origin my-new-feature` )
7. Create new pull request

## License

Awstaghelper is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/mpostument/awstaghelper/blob/master/LICENSE)
