#!/bin/bash
# Script used to retrieve params and start the containers

function logMessage() {
    echo "$(date) $1"
}

logMessage "Setting app server to listen on 8080"
export PORT="8080"

logMessage "Reading in database parameters from ssm"

params=`aws ssm get-parameters-by-path --path /re-region --with-decryption --region us-east-2`
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
