# user_data script for re-region ec2 instance

function logMessage() {
    echo "$(date) $1"
}

logMessage "Setting script vars"
# email for certbot
EMAIL="matt.curry56@gmail.com"
# directories to store source
CODE_DIR="/code"
APP_DIR="$CODE_DIR/re-region"
CMD_DIR="$APP_DIR/cmd"
CERT_DIR="$CERT_DIR/cmd"


# aws cli
logMessage "Installing AWS CLI"
apt-get update
apt-get install awscli

# install docker
logMessage "Installing Docker"
apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update
apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin

# postgres install
logMessage "Installing Postgres"
apt install postgresql-client-common
apt-get install postgresql-client

# to run docker without sudo
logMessage "Adding Ubuntu user to the docker group"
groupadd docker
usermod -aG docker $USER
newgrp docker 

# pull latest docker images
logMessage "Pulling latest docker images for the application"

# pull needed files to run app
logMessage "Pulling files needed to run the app from S3"
cd $CMD_DIR
touch pull_files.sh
chmod 700 pull_files.sh
echo 'sudo aws s3 cp s3://re-region/cmd . --rec' >> pull_files.sh
./pull_files.sh

# start the app
logMessage "Starting the app"
chmod 700 run_server.sh
./run_server.sh


# setup certs
logMessage "Setting up the certs"
certbot certonly --webroot -w certs \ 
              --preferred-challenges http \
              -d reregion.com -d www.reregion.com \ 
              --non-interactive --agree-tos -m $EMAIL

logMessage "User data script has run successfully."