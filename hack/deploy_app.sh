#!/bin/env bash

source ./hack/common.sh

make images

./hack/push_images_in_local_registry.sh

kubectl create -f deployment/k8s
