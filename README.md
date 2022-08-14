# re-region-api
API to access demographic and cost of living estimates for US states and counties. The API accepts individual tax filer variables as parameters to provide individualized cost of living information for a given region. This app was built for educational purposes (both mine and hopefully others!).

link to API spec: https://www.reregion.com/

## Code
###Summary
The API is developed in Go. Also, SQL scripts are called by the data access layer to pull from the database, and shell scripts are used to intitiate the app on the server and also in the deployment pipeline. Some key points regarding the code developement:
* It is written using only standard libs (other than a Postrges driver dependency. I learned through this that Go's standard lib is very powerful!)
* Swagger UI Dist is embedded in the app, and it is used to serve a defined yaml API spec at the root path of the domain (the link above)
* Interfaces define all key services + data access layer
* A custom error structure is leveraged throughout the application, with all expected app errors defined through different constructors
* Unit tests are defined leveraging interfaces to be based around mocks

###Structure (in src folder)
apperrors: Package implementing custom error struct, holds public constructors for each type of app error <br>
controller: Handler functions for the core endpoints. Also includes utilities to process input and write responses <br>
dao: The data access layer. Holds Postgres implementation of the layer's interface and a "sql" folder holding all source SQL. <br>
logging: Package holds my implementation of an aggregated with public methods for different log levels that is used throughout the app <br>
model: Holds structures returned by core services and marshalled by the controller into JSON responses. Models hold methods tied to their behavior <br>
services: Interfaces and implementations of County, State, and Federal services. These services query/cache source data and return entities <br>
         in the model package to the controller <br>
static: Where the Swagger-UI dist and config is embedded <br>
main.go: Holds the server, where the mux registers all handler functions from the controller. Also includes handlers for the Swagger-UI and health <br> endpoint.


## Infrastructure
The API is deployed on the AWS cloud. The configuration is as follows:
 ### App
  * DNS records in Route53
  * AWS Systems Manager Parameter Store for database secrets
  * IAM for roles and policies required by all services used
  * S3 for artifacts used to run the app
  * VPC containing:
    * EC2 instance
       * Running containers for an nginx webserver and the API through Docker compose
       * nginx configured to enforce HTTPS, the server has an SSL certifricate + private key that autorefreshes using Let's Encrypt and Certbot
    * Postgres RDS instance for the API's relational database
    * Needed subnets, route tables, and security groups
 ### CI/CD
  * AWS CodePipeline with integration with this repository to deploy changes
  * AWS Codebuild 
     * Runs Unit tests defined in the repo's "test" folder
     * Builds Docker images from this repository's code and pushes to ECR
     * Deploys run_server.sh and docker-compose.yml files from repo to S3 used to start the application containers (these are pulled onto the EC2).

## Source Data
Data is sourced to the app's Postgres DB using a dockerized ETL CLI tool I developed. The tool sources taxation related data from excel files published by the Tax Foundation and survery statistics from the Census Bureau Data API and loads to the database. This project is not affiliated with either of those orgnaizations and the ETL does modify the intial source data through aggregation and fuzzy matching. The link to that repository and more information about the source data can be found here: https://github.com/Matthew-Curry/re-region-etl/tree/main

## Next steps
When I have time (and before my AWS free tier runs out) I plan to:
* Migrate the EC2 instance running Docker containers with docker-compose to an instance within an ECS cluster. I expect that seeing how docker-compose configuration maps to ECS service configuration will deepen my understanding of both Docker and ECS
* Specify Infrastructure as code using Cloudformation. Similiar to the ECS point, I think seeing how the console configuration maps to CloudFormation syntax will help solidify my understanding of the tool more so than directly jumping into reading the documentation on it.
* Add AWS CodeDeploy to the pipeline to directly deploy the code onto the instance. I realized after provisioning the instance that its OS (Ubuntu 22.04) actually does not currently have a compatible CodeDeploy agent, so for now I am deploying my startup files to S3 and then pulling them onto the instance. When I migrate to ECS, I will take care to choose an instance type that has a compatible CodeDeploy agent to have a more complete CI/CD pipeline.
