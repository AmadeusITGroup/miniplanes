#!/bin/env bash

source ./hack/common.sh
kubectl delete -f $ROOTDIR/manifests/import_airports.yaml
kubectl delete -f $ROOTDIR/manifests/import_airlines.yaml
kubectl delete -f $ROOTDIR/manifests/import_routes.yaml
kubectl delete -f $ROOTDIR/manifests/import_schedules.yaml
