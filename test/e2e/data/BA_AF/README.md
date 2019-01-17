
# How to generate schedules data

As soon you've `storage` && `mongo` that is running in your env you can import data.


## locally

```
MYPORT=9999
kubectl port-forward pod/$(kubectl get pods -lrole=mongo -o=jsonpath='{.items[0].metadata.name}') ${MYPORT}:27017
```

and in another shell you can execute mongoimporter

```shell
mongoimport --port=$MYPORT -d miniapp -c airports --type csv --file airports.dat --fieldFile=../../../../data/airports_schema.dat
mongoimport --port=$MYPORT -d miniapp -c airlines --type csv --file airlines.dat --fieldFile=../../../../data/airlines_schema.dat
mongoimport --port=$MYPORT -d miniapp -c courses --type csv --file courses.dat --fieldFile=../../../../data/courses_schema.dat
```

Now you can create the schedules....


## minikube
