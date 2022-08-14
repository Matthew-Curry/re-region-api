# re-region-api
API to access demographic and taxation estimates for US states and counties. This app was built for educational purposes (both mine and hopefully others!).

link to API spec: https://www.reregion.com/

## Architecture
The API is deployed on the AWS cloud. The configuration is as follows:
 ### App
  * DNS records in Route53
  * Containers for both the http server and the API running on an EC2 instance
  * Postgres RDS instance for the API's relational database
 ### CI/CD
  * AWS CodePipeline with integration with this repository to deploy changes
  * AWS Codebuild to build Docker images from this repository and deploy to ECR,
    as well as to deploy a shell script and docker-compose.yml file to S3 used to start the application.

## Source Data
Taxation information is sourced to the app's database from datasets published by the Tax Foundation. This application is in no way affiliated or endorsed by the Tax Foundation. This data has been transformed and cleaned for storage in the database so it does not match the orginial form. For instance, the linkage of counties to local tax jurisdictions is performed by this application and is not a part of the Tax Foundation's orginial datasets. Also, the application performs the taxation estimates based on this data.

Tax foundation works are licensed under a Creative Commons Attribution NonCommercial 4.0 International License.

https://taxfoundation.org/copyright-notice/

Links to the original source data sets:

https://taxfoundation.org/publications/federal-tax-rates-and-tax-brackets/

https://taxfoundation.org/publications/state-individual-income-tax-rates-and-brackets/

https://taxfoundation.org/local-income-taxes-2019/

This application uses the Census Bureau Data API to source lifestyle and demographic data on regions, but is not endorsed or certified by the Census Bureau. This is accessed from the Census API at the county level, this applicaiton does the aggregation of those metrics to the state level.

## Next steps
When I have time (and before my free tier runs out) I would like to:
* Migrate the EC2 instance running Docker containers with docker-compose to an instance within and ECS cluster. I expect that seeing how docker-compose configuration maps to ECS service configuration will deepen my understanding of both Docker and ECS
* Specify Infrastructure as code using Cloudformation. Similiar to the ECS point, I think seeing how the console configuration maps to CloudFormation syntax will help solidify my understanding of the tool more so than directly jumping into reading the documentation on it.
* Add AWS CodeDeploy to the pipeline to directly deploy the code onto the instance. I realized after provisioning the instance that its OS (the latest version of Ubuntu) actually does not currently have a compatible CodeDeploy agent, so for now I am deploying my startup files to S3 and then pulling them onto the instance. When I migrate to ECS, I will take care to choose an instance type that has a compatible CodeDeploy agent to have a more complete CI/CD pipeline.
