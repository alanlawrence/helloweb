#!/bin/bash
echo "This script builds a new container and pushes it to Artifact Registry."
echo "Run in same directory as the Dockerfile."
echo
echo "Usage: <script> <app_name> <ver>"
echo 
echo "Example: ./build.sh hello-app 3.1"
echo
echo "Assumes that PROJECT_ID is exported in the env and with "
echo "gcloud config set project <PROJECT_ID>"
echo

# Exit if a version has not been supplied just leaving the help screen.
if [[ ( -z $1 ) || ( -z $2 ) ]]; then exit; fi

APP_NAME=$1
NEW_VER=$2
echo "App name is $APP_NAME"
echo "New version is v$NEW_VER"
echo "env project id: $PROJECT_ID"
echo "glcoud config project id: "
gcloud config get project

# We assume that the Artifact Registry repo exists, e.g. at some point in the past,
# a creation command like this was run:
#   gcloud artifacts repositories create helloweb-repo
#      --repository-format=docker
#      --location=europe
#      --description=Repo-for-Alans-helloweb-container-images
#
# Repositories can be listed by running:
#   gcloud artifacts repositories list
#
# And for a particular repo's images:
#  gcloud artifacts docker images list europe-docker.pkg.dev/alans-gcp-project/helloweb-repo

echo
echo "Make sure the registry is populated with our container image."
echo "Running docker build and push commands for europe-docker.pkg.dev/${PROJECT_ID}/$APP_NAME:v$NEW_VER"
docker build -t europe-docker.pkg.dev/${PROJECT_ID}/helloweb-repo/$APP_NAME:v$NEW_VER .
read -ep "Push to Artifact Registry? (y/n): " push
if [ $push != "y" ]
then 
    echo -n "Exiting build/push after building but without pushing to "
    echo    "Artifact Registry."
    exit
fi
docker push europe-docker.pkg.dev/${PROJECT_ID}/helloweb-repo/$APP_NAME:v$NEW_VER
echo "... build and push commands have been run."
