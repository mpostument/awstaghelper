# Awstaghelper

## Table of Contents

* [Overview](#overview)
* [Installing](#installing)
* [Getting Started](#getting-started)
  + [Ec2](#ec2)
  + [Rds](#rds)
  + [Lambda](#lambda)
  + [S3](#se)
* [Contributing](#contributing)
* [License](#license)

## Overview

## Installing

## Getting Started

### Ec2

#### Get ec2 tags

Get list of ec2 with required tags - `awstaghelper ec2 get-ec2-tags`  
Default filename `ec2Tags.csv`  
Default tags list `Name,Environment`  
Example:
 `awstaghelper ec2 get-ec2-tags --filename ec2Tag.csv --tags Name,Owner --profile main`

#### Tag ec2

Read csv and tag ec2 - `awstaghelper ec2 tag-ec2`  
Default filename `ec2Tags.csv`  
Example:
 `awstaghelper ec2 tag-ec2 --filename ec2Tag.csv --profile main`  

### Rds

#### Get rds tags

Get list of rds with required tags - `awstaghelper rds get-rds-tags`  
Default filename `rdsTags.csv`  
Default tags list `Name,Environment`  
Example:
 `awstaghelper rds get-rds-tags --filename rdsTag.csv --tags Name,Owner --profile main`

#### Tag rds

Read csv and tag rds - `awstaghelper rds tag-rds`  
Default filename `rdsTags.csv`  
Example:
 `awstaghelper rds tag-rds --filename rdsTag.csv --profile main`  

### Lambda

#### Get lambda tags

Get list of lambda with required tags - `awstaghelper lambda get-lambda-tags`  
Default filename `lambdaTags.csv`  
Default tags list `Name,Environment`  
Example:
 `awstaghelper lambda get-lambda-tags --filename lambdaTag.csv --tags Name,Owner --profile main`

#### Tag lambda

Read csv and tag lambda - `awstaghelper lambda tag-lambda`  
Default filename `lambdaTags.csv`  
Example:
 `awstaghelper lambda tag-lambda --filename lambdaTag.csv --profile main`  

### S3

#### Get s3 tags

Get list of s3 with required tags - `awstaghelper s3 get-s3-tags`  
Default filename `s3Tags.csv`  
Default tags list `Environment,Project`  
Example:
 `awstaghelper s3 get-s3-tags --filename s3Tag.csv --tags Name,Owner --profile main`

#### Tag s3

Read csv and tag s3 - `awstaghelper s3 tag-s3`  
Default filename `ec2Tags.csv`  
Example:
 `awstaghelper s3 tag-s3 --filename ec2Tag.csv --profile main`  

### General parameters

`filename` - path where to write or read data. Supported by every option
`tags` - list of tags to read. Supported only with `get-tags` option
`region` - aws region to use. Supported by every option
`profile` - aws profile to use. Supported by every option

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
