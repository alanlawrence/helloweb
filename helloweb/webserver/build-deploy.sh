#!/bin/bash
echo "This script builds and deploys a new container in 3 steps after editing main.go:"
echo " 1. docker build"
echo " 2. docker push // to the registry"
echo " 3. update the helloweb-deployment.yaml with the new version using sed"
echo "   3.1 sed yaml > yaml.new"
echo "   3.2 mv yaml yaml.<timestamp>"
echo "   3.3 mv yaml.new yaml"
echo " 4. kubectl apply deployment"
echo " 5. kubectl get service // to see it has been done"
echo
echo "Usage: <script> <ver>"
echo 
echo "Example: ./build-deploy.sh 3.1"
echo "WARNING: assumes the app has been launched with ./launch-helloweb-app.sh"
echo "Run in same directory as main.go and Dockerfile. Assumes yaml in ./manifests"
echo

APP_NAME=hello-app
NEW_VER=$1
echo "New version is v$NEW_VER"
echo
echo "Setting up PROJECT_ID ..."
export PROJECT_ID=alans-gcp-project
echo "glcoud config set project id: $PROJECT_ID"
gcloud config set project $PROJECT_ID


echo
echo "Make sure the registry is populated with our container image."
echo "Running docker build and push commands for gcr.io/${PROJECT_ID}/$APP_NAME:v$NEW_VER"
docker build -t gcr.io/${PROJECT_ID}/$APP_NAME:v$NEW_VER .
docker push gcr.io/${PROJECT_ID}/$APP_NAME:v$NEW_VER
echo
echo "Now edit the deployment yaml file ..."
cd manifests
pwd
YAML_FILE="helloweb-deployment.yaml"
echo "yaml filename=\"$YAML_FILE\""
echo
sed "s/\(image: gcr.io\/alans-gcp-project\/$APP_NAME:\)v.*$/\1v$NEW_VER/" $YAML_FILE > $YAML_FILE.new
echo "Created $YAML_FILE.new"

mv $YAML_FILE $YAML_FILE.`date -I'seconds'`
echo "Backed up old file with date-time suffix"

mv $YAML_FILE.new $YAML_FILE
echo "Renamed $YAML_FILE.new file to be $YAML_FILE"
echo
echo "yaml file now ready."
echo
echo "Now kubectl apply the deployment and get the service and image for confirmation ..."
kubectl apply -f $YAML_FILE
kubectl get service
kubectl get pods -o custom-columns=IMAGE:.spec.containers[0].image
echo 
echo "Done."
echo
