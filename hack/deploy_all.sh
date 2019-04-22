#/bin/env bash

source ./hack/common.sh

#minikube start --kubernetes-version=v1.11.0 --vm-driver=kvm2

#TODO: handle registry
#minikube addons enable registry
eval $(minikube docker-env)

wait_until minikube_up_and_running

wait_until kube-system_up_and_running

#./hack/start_mongo.sh
#Start mongo

TMPFILE=$(mktemp)
/usr/bin/openssl rand -base64 741 > $TMPFILE

#TODO: checks errs :) for example if resources already created
kubectl create secret generic shared-bootstrap-data --from-file=internal-auth-mongodb-keyfile=$TMPFILE
kubectl apply -f deployment/k8s/mongo.yaml

make image

wait_until mongo_up_and_running 2 60
rm -rf $TMPFILE

#Populate mongo...


kubectl create -f deployment/k8s/miniplanes

#TODO: create schedules
