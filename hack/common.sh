#!/bin/env bash

command -v kubectl >/dev/null 2>&1 || { echo >&2 "can't find kubectl.  Aborting."; exit 1; }
command -v minikube >/dev/null 2>&1 || { echo >&2 "can't find minikube. Aborting."; exit 1; }
command -v openssl >/dev/null 2>&1 || { echo >&2 "can't find openssl. Aborting."; exit 1; }
command -v curl >/dev/null 2>&1 || { echo >&2 "can't find curl. Aborting."; exit 1; }

ROOTDIR=$(git rev-parse --show-toplevel)

echo_red() {
  local RED='\033[0;31m'
  NC='\033[0m'
  printf "${RED}$1${NC}\n"
}

echo_yellow() {
  local YELLOW='\033[1;33m'
  NC='\033[0m'
  printf "${YELLOW}$1${NC}\n"
}

echo_green() {
  local GREEN='\033[0;32m'
  NC='\033[0m'
  printf "${GREEN}$1${NC}\n"
}

wait_until() {
  local script=$1
  local wait=${2:-.5}
  local timeout=${3:-10}
  local i

  script_pretty_name=$(echo "$script" | sed 's/_/ /g')
  times=$(echo "($(bc <<< "scale=2;$timeout/$wait")+0.5)/1" | bc)
  for i in $(seq 1 "$times"); do
    local out=$($script)
    if [ "$out" == "0" ]
    then
      echo_green "${script_pretty_name}: OK"
      return 0
    fi
    echo_yellow "${script_pretty_name}: Waiting..."
    sleep $wait
  done
  echo_red "${script_pretty_name}: ERROR"
  return 1
}

minikube_up_and_running() {
  clusterStatus=$(minikube status --format='{{ .ClusterStatus }}')
  minikubeStatus=$(minikube status --format='{{ .MinikubeStatus }}')
  if [[ "${clusterStatus}" == "Running" && "${minikubeStatus}" == "Running" ]]
  then
    echo "0"
    return
  fi
  echo "1"
}

kube-system_up_and_running() {
  nonrunning=$(kubectl get pods -n kube-system --field-selector=status.phase!=Running 2> /dev/null)
  if [ -z "$nonrunning" ]
  then
    echo 0 #No nonrunning found
    return
  fi
  echo 1
}


mongo_up_and_running() {
  mongopod=$(kubectl get pods -lrole=mongo --field-selector=status.phase=Running 2> /dev/null)
  if [ -z "$mongopod" ]
  then
    echo 1 #no mongo found
    return
  fi
  echo 0
}


local-registry_up_and_running() {
  howmanyregistry=$(kubectl get pods -lk8s-app=kube-registry --field-selector=status.phase=Running -n kube-system --no-headers=true | wc -l 2> /dev/null)
  if [[ ( "$howmanyregistry" != 2 ) ]] ;
  then
    echo 1 #no enough registry found
    return
  fi
  echo 0
}

port_5000_forwaded() {
  curl http://localhost:5000
  echo $?
}
