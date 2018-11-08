# `miniapp`: a synthetic workload to test Federation-v2


# How to install it



## How to inject data in DB (locally)

`backend` exposes the REST API for the server DB


## To deploy it locally

TODO fill this

## To deploy in minikube

### Create a local registry

There're multiple way to deploy a local registry in minikube...
Here we use what is described [here](link here) so you should  https://gist.github.com/coco98/b750b3debc6d517308596c248daf3bb1

```shell
kubectl create -f https://gist.githubusercontent.com/coco98/b750b3debc6d517308596c248daf3bb1/raw/6efc11eb8c2dce167ba0a5e557833cc4ff38fa7c/kube-registry.yaml
```

if everything is you should obtain something similar

```shell
$ minikube ssh && curl localhost:5000
                         _             _
            _         _ ( )           ( )
  ___ ___  (_)  ___  (_)| |/')  _   _ | |_      __
/' _ ` _ `\| |/' _ `\| || , <  ( ) ( )| '_`\  /'__`\
| ( ) ( ) || || ( ) || || |\`\ | (_) || |_) )(  ___/
(_) (_) (_)(_)(_) (_)(_)(_) (_)`\___/'(_,__/'`\____)
$
```

### Deploy mongo

First of all you need `Mongo` in minikube

```shell
$ TMPFILE=$(mktemp)
$ /usr/bin/openssl rand -base64 741 > $TMPFILE
$ kubectl create secret generic shared-bootstrap-data --from-file=internal-auth-mongodb-keyfile=$TMPFILE
$ kubectl apply -f manifests/mongo.yaml
```


once you have deployed mongo it can be locally forwarded via something like


```shell
$ MYPORT=9999 #your unused port should be here
$ kubectl port-forward pod/$(kubectl get pods -lrole=mongo -o=jsonpath='{.items[0].metadata.name}') ${MYPORT}:27017
Forwarding from 127.0.0.1:9999 -> 27017
Forwarding from [::1]:9999 -> 27017
```

Now you can populate mongo using the `9999` port (for example). Evertime you run this command the port number may change.

```shell
$ mongoimport --port=${MYPORT} -d miniapp -c airports --type csv --file data/airports.dat --fieldFile=data/airports_schema.dat
$ mongoimport --port=${MYPORT} -d miniapp -c airlines --type csv --file data/airlines.dat --fieldFile=data/airlines_schema.dat
$ mongoimport --port=${MYPORT} -d miniapp -c routes --type csv --file data/routes.dat --fieldFile=data/routes_schema.dat
$ mongoimport --port=${MYPORT} -d miniapp -c schedules --type csv --file data/schedules.dat --fieldFile=data/schedules_schema.dat
```

Now your mongo db should be populated.

Otherwise you can populate using `jobs` defined in `.../manifests/import_*.yaml`.
First you need to build the docker image in `.../data/`

```shell
$ docker build . importer
$ kubectl port-forward --namespace kube-system $(kubectl get pods -n kube-system -l=k8s-app=kube-registry,version=v0 -o=jsonpath='{.items[0].metadata.name}') 5000:5000
$ docker push importer
```

and after 

```shell
$ kubectl create -f import_airports.yaml
$ kubectl create -f import_airlines.yaml
$ kubectl create -f import_routes.yaml
$ kubectl create -f import_schedules.yaml
```

the mongo db should be populated.