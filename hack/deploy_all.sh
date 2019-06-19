#/bin/env bash

source ./hack/common.sh

#TODO check if kind is already running

kind create cluster --name miniplanes
export KUBECONFIG="$(kind get kubeconfig-path --name="miniplanes")"
#TODO: check if cluster is running
wait_until kube-system_up_and_running


kubectl create -f deployment/k8s/kube-registry.yaml

wait_until kube-system_up_and_running

TMPFILE=$(mktemp)
/usr/bin/openssl rand -base64 741 > $TMPFILE

#TODO: checks errs :) for example if resources already created
kubectl create secret generic shared-bootstrap-data --from-file=internal-auth-mongodb-keyfile=$TMPFILE
kubectl apply -f deployment/k8s/mongo.yaml
wait_until mongo_up_and_running 2 60
rm -rf $TMPFILE

make images

#forward registy svc
kubectl port-forward --namespace kube-system $(kubectl get pods -n kube-system -l=k8s-app=kube-registry,version=v0 -o=jsonpath='{.items[0].metadata.name}') 5000:5000 &
REGPFPID=$!
wait_until port_5000_forwaded

# and push images
docker tag miniplanes localhost:5000/miniplanes
docker push localhost:5000/miniplanes

#check whether images have been pushed
ok=$(curl -X GET http://localhost:5000/v2/_catalog  | jq '.repositories | index("miniplanes")')

kill $REGPFPID


#Create miniplanes
kubectl create -f deployment/k8s/miniplanes

#Create Airports
#Create Airlines
#Create Schedules

#Expose UI (service)
#Test
