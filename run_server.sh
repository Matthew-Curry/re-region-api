#!/bin/bash
# Script used to retrieve params and start the containers
# include argument "r" if docker images should be refreshed

# AWS vars
export AWS_REGION='us-east-2'
export AWS_ACCOUNT='643931054710'
# docker images
export API_IMAGE='re-region-api:latest'
export SERVER_IMAGE='re-region-nginx:latest'
# params to pull
PARAM_PATH='/re-region'

function logMessage() {
    echo "$(date) $1"
}

if [[ $1 == "r" ]]; then
    logMessage "Refresh flag given, retrieving updated docker images"
    # get ecr token
    aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com
    # pull both images needed
    docker pull $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/$SERVER_IMAGE
    docker pull $AWS_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/$API_IMAGE
fi

logMessage "Setting app server to listen on 8080"
export PORT="8080"

logMessage "Reading in database parameters from ssm"

params=`aws ssm get-parameters-by-path --path $PARAM_PATH --with-decryption --region $AWS_REGION`
logMessage "Successfully got params from ssm"

# process the parameters into variables with the values
for i in $(seq 0 4); do
    n=`echo $params | jq .Parameters[$i].Name | sed s/'"'/""/g`
    v=`echo $params | jq .Parameters[$i].Value | sed s/'"'/""/g`
    if [[ "$n" == '/re-region/api-db-password' ]]; then
        export RE_REGION_API_PASSWORD="$v"
        logMessage "Read in DB password"
    elif [[ "$n" == '/re-region/api-db-user' ]]; then
        export RE_REGION_API_USER="$v"
        logMessage "Read in DB user"
    elif [[ "$n" == '/re-region/db' ]]; then
        export RE_REGION_DB="$v"
        logMessage "Read in DB name"
    elif [[ "$n" == '/re-region/db-host' ]]; then
        export DB_HOST="$v"
        logMessage "Read in DB host"
    elif [[ "$n" == '/re-region/db-port' ]]; then
        export DB_PORT="$v"
        logMessage "Read in DB port"
    fi
done

# now that needed env vars are exported, start the container
logMessage "Starting the containers"
docker-compose up --remove-orphans
