#!/bin/env bash

function handle_int() {
  kill $REGPFPID
  echo "port-forward stopped"
  exit
}

source ./hack/common.sh

cd ./data

docker build . -t localhost:5000/importer

kubectl port-forward --namespace kube-system $(kubectl get pods -n kube-system -l=k8s-app=kube-registry,version=v0 -o=jsonpath='{.items[0].metadata.name}') 5000:5000 &
REGPFPID=$!

wait_until port_5000_forwaded

docker push localhost:5000/importer

kubectl create -f $ROOTDIR/manifests/import_airports.yaml
kubectl create -f $ROOTDIR/manifests/import_airlines.yaml
kubectl create -f $ROOTDIR/manifests/import_routes.yaml
kubectl create -f $ROOTDIR/manifests/import_schedules.yaml

kill $REGPFPID
