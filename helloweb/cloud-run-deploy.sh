#!/bin/bash
echo "-----------------------------------------------------------------------------"
echo "This script deploys a new container to Cloud Run"
echo
echo "Usage: <script> <Google artefact registry path> <region>"
echo 
echo "Example: $0 europe-docker.pkg.dev/alans-gcp-project/helloweb-repo/hello-app:\
v3.4.1 \
europe-west2"
echo "-----------------------------------------------------------------------------"
echo

# Exit if a version has not been supplied just leaving the help screen.
if [ "$#" -ne 2 ]; then exit; fi

IMAGE_VER=$1
REGION=$2
echo "Using image version: $IMAGE_VER"
echo "Using region: $REGION"

echo "Setting up PROJECT_ID ..."
export PROJECT_ID=alans-gcp-project
echo "glcoud config set project id: $PROJECT_ID"
gcloud config set project $PROJECT_ID

APP_NAME="hello-app"
echo "Deploying $APP_NAME ..."
gcloud run deploy $APP_NAME \
    --image $IMAGE_VER \
    --region $REGION \
    --project $PROJECT_ID

echo 
echo "Done."
echo
