#!/bin/env bash

source ./hack/common.sh

kubectl port-forward --namespace kube-system $(kubectl get pods -n kube-system -l=k8s-app=kube-registry,version=v0 -o=jsonpath='{.items[0].metadata.name}') 5000:5000 &
REGPFPID=$!

wait_until port_5000_forwaded

docker push localhost:5000/storage
docker push localhost:5000/itineraries-server
docker push localhost:5000/ui

#minikube ssh 'curl -X GET http://localhost:5000/v2/_catalog'
# TODO: some jq-ness to check if images have been loaded

kill $REGPFPID
