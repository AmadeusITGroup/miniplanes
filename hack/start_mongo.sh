#!/bin/env bash

source ./hack/common.sh

TMPFILE=$(mktemp)
/usr/bin/openssl rand -base64 741 > $TMPFILE
kubectl create secret generic shared-bootstrap-data --from-file=internal-auth-mongodb-keyfile=$TMPFILE
kubectl apply -f manifests/mongo.yaml

wait_until mongo_up_and_running

rm -rf $TMPFILE