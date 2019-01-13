#!/bin/env bash

source ./hack/common.sh

kubectl delete secret shared-bootstrap-data

kubectl delete  -f manifests/mongo.yaml

#TODO: check whether  mongo is still running wait_until
