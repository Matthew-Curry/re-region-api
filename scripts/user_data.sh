# user_data script for re-region

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
sudo apt-get update
sudo apt-get install awscli

# install docker
logMessage "Installing Docker"
sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin

# postgres install
logMessage "Installing Postgres"
sudo apt install postgresql-client-common
sudo apt-get install postgresql-client

# to run docker without sudo
logMessage "Adding Ubuntu user to the docker group"
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker 

# pull latest docker images
logMessage "Pulling latest docker images for the application"

# pull needed files to run app
logMessage "Pulling files needed to run the app from S3"
cd $CMD_DIR
sudo aws s3 cp s3://re-region/cmd/docker-compose.yml .
sudo aws s3 cp s3://re-region/cmd/run.sh .

# make the directories
logMessage "Making directories to store the code"
sudo mkdir $CODE_DIR
sudo mkdir $APP_DIR
sudo mkdir $CMD_DIR
sudo mkdir $CERT_DIR

# start the app
logMessage "Starting the app"

# setup certs
logMessage "Setting up the certs"
sudo certbot certonly --webroot -w certs \ 
              --preferred-challenges http \
              -d reregion.com -d www.reregion.com \ 
              --non-interactive --agree-tos -m $EMAIL

logMessage "User data script has run successfully."