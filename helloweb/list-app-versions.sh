#!/bin/bash
echo "---------------------------------------------------"
echo "This script list app versions in Artifact Registry."
echo
echo "Usage: <script> <app_name>"
echo
echo "Example: ./list-app-versions.sh hello-app"
echo
echo "Assumes that PROJECT_ID is exported in the env and set with"
echo "gcloud config set project <PROJECT_ID>"
echo
echo "Example: gcloud config set project alans-gcp-project"
echo "Example: export PROJECT_ID=alans-gcp-project"
echo "---------------------------------------------------"
echo

# Exit if an app name has not been supplied just leaving the help screen.
if [[ -z $1 ]]; then exit; fi

APP_NAME=$1
echo "App name is $APP_NAME"
echo "env project id: $PROJECT_ID"
echo -n "glcoud config project id: "
gcloud config get project
echo

# We assume that the Artifact Registry repo exists, e.g. at some point in the
# past, a creation command like this was run:
#   gcloud artifacts repositories create helloweb-repo
#      --repository-format=docker
#      --location=europe
#      --description=Repo-for-Alans-helloweb-container-images
#
# Repositories can be listed by running:
#   gcloud artifacts repositories list
#
# And for a particular repo's images with versions and most recently
# updated first:
echo "Checking the 'europe-docker.pkg.dev' Artifact Registry location"
echo "in the helloweb-repo directory:"
echo "First set up the meta data formatting parameters:"
FORMAT_PARAMETERS=(
    "table(package:label=IMAGE,"
    "version:label=DIGEST,"
    "tags:label=VERSION,"
    "metadata.imageSizeBytes.size():label=SIZE,"
    "create_time:label=CREATE,"
    "update_time:label=UPDATED)"
)
printf "    %s\n" "${FORMAT_PARAMETERS[@]}"
echo

gcloud artifacts docker images list \
europe-docker.pkg.dev/${PROJECT_ID}/helloweb-repo/${APP_NAME} \
--include-tags \
--sort-by="~update_time" \
--format="${FORMAT_PARAMETERS[*]}"

