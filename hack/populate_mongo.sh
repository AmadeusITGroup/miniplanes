#!/usr/bin/env bash

source ./hack/common.sh

cd ./data

docker build . -t importer

kubectl create -f $ROOTDIR/manifests/import_airports.yaml
kubectl create -f $ROOTDIR/manifests/import_airlines.yaml
kubectl create -f $ROOTDIR/manifests/import_courses.yaml
#kubectl create -f $ROOTDIR/manifests/import_schedules.yaml


#TODO automatize check to see if courses/airports/airlines have been imported
#TODO2: use workflow-controller to schedule in the right order

#kubectl create -f $ROOTDIR/manifests/generate_schedules.yaml
