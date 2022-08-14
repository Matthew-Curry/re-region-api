# re-region-api
API to access demographic and cost of living estimates for US states and counties. The API accepts individual tax filer variables as parameters to provide individualized cost of living information for a given region. This app was built for educational purposes (both mine and hopefully others!).

link to API spec: https://www.reregion.com/

## Architecture
The API is deployed on the AWS cloud. The configuration is as follows:
 ### App
  * DNS records in Route53
  * AWS Systems Manager Parameter Store for database secrets
  * IAM for roles and policies required by all services used
  * VPC containing:
    * Containers for both the http server and the API running on an EC2 instance
    * Postgres RDS instance for the API's relational database
    * Needed subnets, route tables, and security groups
 ### CI/CD
  * AWS CodePipeline with integration with this repository to deploy changes
  * AWS Codebuild to build Docker images from this repository and deploy to ECR,
    as well as to deploy a shell script and docker-compose.yml file to S3 used to start the application containers.

## Source Data
Data is sourced to the app's Postgres DB using a Dockerized ETL tool I developed. The tool sources data from excel files published by the Tax Foundation and the Census Bureau Data API. The link to that repository and more information about the source data can be found here: https://github.com/Matthew-Curry/re-region-etl/tree/main

## Next steps
When I have time (and before my free tier runs out) I would like to:
* Migrate the EC2 instance running Docker containers with docker-compose to an instance within an ECS cluster. I expect that seeing how docker-compose configuration maps to ECS service configuration will deepen my understanding of both Docker and ECS
* Specify Infrastructure as code using Cloudformation. Similiar to the ECS point, I think seeing how the console configuration maps to CloudFormation syntax will help solidify my understanding of the tool more so than directly jumping into reading the documentation on it.
* Add AWS CodeDeploy to the pipeline to directly deploy the code onto the instance. I realized after provisioning the instance that its OS (the latest version of Ubuntu) actually does not currently have a compatible CodeDeploy agent, so for now I am deploying my startup files to S3 and then pulling them onto the instance. When I migrate to ECS, I will take care to choose an instance type that has a compatible CodeDeploy agent to have a more complete CI/CD pipeline.
