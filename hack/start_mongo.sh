#!/bin/bash

source ./hack/common.sh

TMPFILE=$(mktemp)
/usr/bin/openssl rand -base64 741 > $TMPFILE
#TODO: checks errs :) for example if resources already created

kubectl create secret generic shared-bootstrap-data --from-file=internal-auth-mongodb-keyfile=$TMPFILE
kubectl apply -f manifests/mongo.yaml

wait_until mongo_up_and_running 1 60
rm -rf $TMPFILE
