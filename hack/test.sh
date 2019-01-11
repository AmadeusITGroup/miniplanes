#!/bin/env bash

function handle_int() {
  kill $REGPFPID
  echo "port-forward stopped"
  exit
}

source ./hack/common.sh

ui_pod_name=$(kubectl get pods -l app=ui -o=jsonpath='{.items[0].metadata.name}')
ui_svc_target_port=$(kubectl get svc -l name=ui -o=jsonpath='{.items[0].spec.ports[0].targetPort}')
ui_svc_port=$(kubectl get svc -l name=ui -o=jsonpath='{.items[0].spec.ports[0].port}')
kubectl port-forward ${ui_pod_name} ${ui_svc_target_port}:${ui_svc_port} &

UISVCPID=$!

echo "RUN TEST TEST E2E.. But for the moment... Just see if ui is there"

STATUS=$(http -h  localhost:8090/airports | grep HTTP | cut -d ' ' -f 2)
if  [ "$STATUS" = "200" ]; then
  echo "OK"
else
  echo "KO"
fi

kill $UISVCPID
