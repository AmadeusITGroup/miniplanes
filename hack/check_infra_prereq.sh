#!/bin/env bash

source ./hack/common.sh

wait_until minikube_up_and_running

wait_until kube-system_up_and_running

wait_until local-registry_up_and_running
