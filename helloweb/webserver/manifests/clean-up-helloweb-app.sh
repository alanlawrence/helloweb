#!/bin/bash
if [ -z $1 ]; then ARG1=none; else ARG1=$1; fi
if [ $ARG1 = "mutate" ]
then
  ECHO=
else
  ECHO=echo
fi
echo "This script cleans up the helloweb app as follows:"
echo "  * deletes the ingress and service for the 'hello' app"
echo "  * deletes the static ip address 'helloweb-ip'"
echo "  * deletes the deployment"
echo "  * checks the forwarding-rules are no longer listed"
echo "  * deletes the cluster 'hello-cluster'"
echo "Assumes the project ID is set and configured already."
echo "Assumes the script is run in the manifests directory."
echo "Running in $MODE mode. Modes: arg1 absent = echo; arg1 = mutate"
read -ep "Clean up? (y/n): " cleanup
if [ $cleanup != "y" ] 
then
    exit
fi
export PROJECT_ID=alans-gcp-project
echo "Using project id: $PROJECT_ID to set up ..."
$ECHO gcloud config set project $PROJECT_ID
$ECHO gcloud config set compute/zone europe-west2-a
echo "Starting clean up ..."
echo "Print current status ..."
$ECHO gcloud config get-value compute/zone
$ECHO gcloud compute instances list
$ECHO gcloud compute addresses describe helloweb-ip --global
echo
echo "*** Delete ingress and service ..."
$ECHO kubectl delete ingress,service -l app=hello
echo "*** Delete the static ip ..."
$ECHO gcloud compute addresses delete helloweb-ip --global
echo "*** Delete the deployment ..."
$ECHO kubectl delete -f helloweb-deployment.yaml
echo "Wait until the load balancer is deleted by watching the output of the forwarding-rule list command "
echo "    until a forwarding rule that contains "helloweb" in its name is no longer present) ..."
keepChecking=1
while [ $keepChecking -gt 0 ]
do
    keepChecking=`gcloud compute forwarding-rules list 2>&1 | grep -c helloweb`
    echo -n .
    if [ -n "$ECHO" ] && [ "echo" = $ECHO ]
    then
        keepChecking=0
    fi
    if [ $keepChecking -gt 0 ]
    then
        sleep 2s
    fi
done
echo
echo "   ... done, forwarding rule no longer present for "helloweb"."
echo "*** Delete the cluster ..."
$ECHO gcloud container clusters delete hello-cluster
echo "DONE!"
