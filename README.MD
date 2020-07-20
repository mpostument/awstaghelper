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
By default result will be stored in file `ec2Tags.csv` . To change filename add `--filename newName.csv`
Tag `Name,Environment` will be added to csv. To use custom tags specify `--tags Name,Owner,CostCenter`

#### Tag ec2

Edit csv by adding new tags or changin value to existing one  
Run `awstaghelper ec2 tag-ec2` . All tags from csv will be applied to corresponding instances  
Custom filename can be specified with `--filename newName.csv`

### Rds

#### Get rds tags

Get list of rds with required tags - `awstaghelper rds get-rds-tags`
By default result will be stored in file `rdsTags.csv` . To change filename add `--filename newName.csv`
Tag `Name,Environment` will be added to csv. To use custom tags specify `--tags Name,Owner,CostCenter`

#### Tag rds

Edit csv by adding new tags or changin value to existing one  
Run `awstaghelper rds tag-rds` . All tags from csv will be applied to corresponding instances  
Custom filename can be specified with `--filename newName.csv`

### Lambda

#### Get lambda tags

Get list of lambda with required tags - `awstaghelper lambda get-lambda-tags`
By default result will be stored in file `lambdaTags.csv` . To change filename add `--filename newName.csv`
Tag `Name,Environment` will be added to csv. To use custom tags specify `--tags Name,Owner,CostCenter`

#### Tag lambda

Edit csv by adding new tags or changin value to existing one  
Run `awstaghelper lambda tag-lambda` . All tags from csv will be applied to corresponding instances  
Custom filename can be specified with `--filename newName.csv`

### S3

#### Get s3 tags

Get list of ec2 with required tags - `awstaghelper s3 get-s3-tags`
By default result will be stored in file `s3Tags.csv` . To change filename add `--filename newName.csv`
Tag `Name,Environment` will be added to csv. To use custom tags specify `--tags Name,Owner,CostCenter`

#### Tag s3

Edit csv by adding new tags or changin value to existing one  
Run `awstaghelper s3 tag-lambda` . All tags from csv will be applied to corresponding instances  
Custom filename can be specified with `--filename newName.csv`

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
