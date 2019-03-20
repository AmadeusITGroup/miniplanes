

## Working locally on your laptop (no container)

### Test locally

```shell
$ ./hack/test_e2e_local.sh
```


## Working with minikube

### Populating mongo with local forward
```shell
$ MYPORT=9999 #your unused port should be here
$ kubectl port-forward pod/$(kubectl get pods -lrole=mongo -o=jsonpath='{.items[0].metadata.name}') ${MYPORT}:27017
Forwarding from 127.0.0.1:9999 -> 27017
Forwarding from [::1]:9999 -> 27017
```

Now you can pouplate with usual `mongoimporter`

```shell
$ mongoimport --port=${MYPORT} -d miniplanes -c airports --type csv --file data/airports.dat --fieldFile=data/airports_schema.dat
$ mongoimport --port=${MYPORT} -d miniplanes -c airlines --type csv --file data/airlines.dat --fieldFile=data/airlines_schema.dat
$ mongoimport --port=${MYPORT} -d miniplanes -c routes --type csv --file data/routes.dat --fieldFile=data/routes_schema.dat
$ mongoimport --port=${MYPORT} -d miniplanes -c schedules --type csv --file data/schedules.dat --fieldFile=data/schedules_schema.dat
```

## Working with kind

TODO
