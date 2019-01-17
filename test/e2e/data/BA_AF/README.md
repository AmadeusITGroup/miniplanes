
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

with something like
```shell
cd .../miniapp/schedules-generator/cmd
go run main.go --port=$MYPORT
```

which generates a file `schedules.csv` that you can import in the usual way

```shell
mongoimport --port=$MYPORT -d miniapp -c schedules --type csv --file schedules.dat --fieldFile=../../../../data/schedules_schema.dat
```

with a

```shell
mongo --port=$MYPORT
```

you should be able to see `schedules` collection populated

## minikube
TODO
