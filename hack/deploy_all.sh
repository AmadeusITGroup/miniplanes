#/bin/env bash

source ./hack/common.sh

minikube start --kubernetes-version=v1.11.0 --vm-driver=kvm

wait_until minikube_up_and_running

kubectl create -f https://gist.githubusercontent.com/coco98/b750b3debc6d517308596c248daf3bb1/raw/6efc11eb8c2dce167ba0a5e557833cc4ff38fa7c/kube-registry.yaml

wait_until kube-system_up_and_running

./hack/start_mongo.sh
