# Terrafrom for lamport on AWS

This is my development area for developing terraform to deploy a three node lamport cluster to AWS

Not sure what orchestration packacge I will use yet

This project is simply to setup and deploy the three nodes ande make them usable.


## Terraform
###Finding Terraform:
https://www.terraform.io/

###Setting up your environment:
https://www.terraform.io/intro/getting-started/install.html

Note that terraform will setup everything in the current directory.
Keep only files for one setup in the directory, e.g. lamportaws.tf 
Use "terraform plan" to see what it will set up

###Set up these env variables for logging
"export TF_LOG=TRACE"
"export TF_LOG_PATH=./terraform.log"



## AWS
###Setting up account on AWS:
http://docs.aws.amazon.com/AmazonSimpleDB/latest/DeveloperGuide/AboutAWSAccounts.html

Note: 
You can use a free tier

###Setting up AWS key pairs
http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html


### AWS security Credentials
https://console.aws.amazon.com/iam/home?#security_credential
