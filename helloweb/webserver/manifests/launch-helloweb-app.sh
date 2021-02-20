#!/bin/bash
if [ -z $1 ]; then ARG1=none; else ARG1=$1; fi
if [ $ARG1 = "mutate" ]
then
  ECHO=
  MODE=mutate
else
  ECHO=echo
  MODE=echo
fi
echo "This script launches the helloweb app as follows:"
echo "   * Assumes the service, ingress and static-ip yamls are in the cwd"
echo "   * Assumes the container image is built and registered in gcr.io"
echo "     - check helloweb-deployment.yaml references the desired image"
echo "   * https only"
echo "   * 1 pod"
echo "   * static ip, ensure it is conigured on your domain: alanjl.info"
echo "You will be prompted to check your domains configuration."
echo "WARNING: The script is not designed to be idempotent."
echo "---------------------------------------------------------------"
echo "Running in $MODE mode. Modes: arg1 absent = echo; arg1 = mutate"
read -ep "Launch? (y/n): " launch
if [ $launch != "y" ] 
then
    exit
fi
echo "Starting launch ..."
export PROJECT_ID=alans-gcp-project
echo "Using project id: $PROJECT_ID"
$ECHO gcloud config set project $PROJECT_ID
$ECHO gcloud config set compute/zone europe-west2-a
$ECHO gcloud config get-value compute/zone
$ECHO gcloud container clusters create hello-cluster --num-nodes=1 --release-channel=rapid
$ECHO gcloud compute instances list
$ECHO gcloud compute addresses create helloweb-ip --global
$ECHO gcloud compute addresses describe helloweb-ip --global
echo "Check that the domain alanjl.info is configured with the static ip shown."
echo "   * Visit https://domains.google.com/m/registrar/alanjl.info/dns"
echo "   * and edit the type A entry to match the static ip shown."
echo "Look for: Created [...]. "
echo "          address: <static ip> <-----"
echo "          addressType: EXTERNAL."
read -ep "Domain configured with static IP address? (y/n): " configured
if [ $configured != "y" ] 
then
    exit
fi
$ECHO kubectl apply -f helloweb-deployment.yaml
$ECHO kubectl get pods
$ECHO kubectl apply -f helloweb-ingress-static-ip.yaml
$ECHO kubectl get ingress helloweb
$ECHO kubectl expose deployment helloweb --type=LoadBalancer --port 80 --target-port 8080
echo "..."
echo "It might take 5 mins or so for the application to be accessible."
echo "This command:  curl https://alanjl.info/"
echo "    should be accessible after 5 mins or so."
echo 
echo "This command: curl http://alanjl.info/"
echo "    should return a 404 Not Found error after 5 mins or so."
echo "---"
echo

