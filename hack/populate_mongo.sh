#!/bin/env bash

function handle_int() {
  kill $REGPFPID
  echo "port-forward stopped"
  exit
}

source ./hack/common.sh

cd ./data

docker build . -t importer

kubectl create -f $ROOTDIR/manifests/import_airports.yaml
kubectl create -f $ROOTDIR/manifests/import_airlines.yaml
kubectl create -f $ROOTDIR/manifests/import_routes.yaml
kubectl create -f $ROOTDIR/manifests/import_schedules.yaml

kill $REGPFPID
